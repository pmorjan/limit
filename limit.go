// Package limit provides a simple rate limiter for concurrent access.
package limit

import (
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	// MaxKeepRecords defines how long records are keept.
	MaxKeepRecords time.Duration = time.Minute * 5
	// ErrInvalidRate is returned to indicate an invalid rate.
	ErrInvalidRate error = errors.New("rate value too low")
)

// Limit defines the limiter.
type Limit struct {
	mu    sync.Mutex
	users map[string]*user
	rate  float64
	burst int
}

type user struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// New creates a limiter with a "token bucket" of size burst, initially full and
// refilled at rate per second.
func New(rate float64, burst int) (*Limit, error) {
	minRate := float64(1 / (MaxKeepRecords / time.Second))
	if rate <= minRate {
		return nil, ErrInvalidRate
	}
	l := &Limit{
		burst: burst,
		rate:  rate,
		users: make(map[string]*user),
	}
	go l.cleanup()
	return l, nil
}

// Allow reports whether an event may happen now.
func (l *Limit) Allowed(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, exists := l.users[id]
	if !exists {
		l.users[id] = &user{
			limiter: rate.NewLimiter(rate.Limit(l.rate), l.burst),
		}
	}
	l.users[id].lastSeen = time.Now()
	return l.users[id].limiter.Allow()
}

func (l *Limit) cleanup() {
	for {
		time.Sleep(time.Minute)
		l.mu.Lock()
		for id, v := range l.users {
			if time.Now().Sub(v.lastSeen) > MaxKeepRecords {
				log.Printf("deleting user %s", id)
				delete(l.users, id)
			}
		}
		l.mu.Unlock()
	}
}
