package limit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pmorjan/limit"
	"github.com/tomasen/realip"
)

func TestAllowed(t *testing.T) {
	limiter, _ := limit.New(1, 2)
	if !limiter.Allowed("foo") {
		t.Fatalf("expected allowed")
	}
	if !limiter.Allowed("foo") {
		t.Fatalf("expected allowed")
	}
	if limiter.Allowed("foo") {
		t.Fatalf("expected not allowed")
	}

	if !limiter.Allowed("bar") {
		t.Fatalf("expected  allowed")
	}

	time.Sleep(time.Second)
	if !limiter.Allowed("foo") {
		t.Fatalf("expected allowed")
	}
}

func TestInvalidRate(t *testing.T) {
	_, err := limit.New(float64(1/limit.MaxKeepRecords), 1)
	if err != limit.ErrInvalidRate {
		t.Fatalf("expected Error: %v", limit.ErrInvalidRate)
	}
}

func ExampleLimit() {

	// middleware function
	limitRequests := func(next http.Handler) http.Handler {
		userLimit, _ := limit.New(1, 3) // 1 request per second, burst 3

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := realip.FromRequest(r)
			if !userLimit.Allowed(ip) {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	srv := httptest.NewServer(limitRequests(handler))
	defer srv.Close()

	// 4 requests
	for i := 0; i < 4; i++ {
		res, _ := http.Get(srv.URL)
		fmt.Println(res.Status)
	}

	// next request after 1 sec delay
	time.Sleep(time.Second)
	res, _ := http.Get(srv.URL)
	fmt.Println(res.Status)

	// Output:
	// 200 OK
	// 200 OK
	// 200 OK
	// 429 Too Many Requests
	// 200 OK
}
