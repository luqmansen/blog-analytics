package api

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

const ThrottleRequest = 2 * time.Second

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.RWMutex

func init() {
	go cleanUpVisitors()
}

func RequestThrottleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				logrus.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			limiter := getVisitorLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(1, 1)
		visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanUpVisitors() {
	for {
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > ThrottleRequest {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
