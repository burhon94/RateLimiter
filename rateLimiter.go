package RateLimiter

import (
	"net/http"
	"strings"
	"time"

	"github.com/burhon94/RateLimiter/limiter/limit"

	"golang.org/x/time/rate"
)

func LimitMiddleware(next http.Handler, limiter *limit.IPRateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getLimiter := limiter.GetLimiter(getIPFromRemoteAddr(r.RemoteAddr))

		if !getLimiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIPFromRemoteAddr(remoteAddr string) (ip string) {
	trim := strings.Split(remoteAddr, ":")
	ip = trim[0]
	return
}

func SetParam(timer time.Duration, count int) *limit.IPRateLimiter {
	limiterFromParam := limit.NewIPRateLimiter(rate.Every(timer*time.Second), count)

	return limiterFromParam
}
