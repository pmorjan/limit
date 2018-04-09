[![GoDoc](https://godoc.org/github.com/pmorjan/limit?status.svg)](https://godoc.org/github.com/pmorjan/limit)

# limit

Go package limit provides a simple rate limiter for concurrent access.

Example middleware to limit HTTP requests:
```go
import (
	"github.com/pmorjan/limit"
	"github.com/tomasen/realip"
)

func limitRequests(next http.Handler) http.Handler {
    userLimit,_ := limit.New(1, 10) // 1 request per second, burst 10

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)
		if !userLimit.Allowed(ip) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```
