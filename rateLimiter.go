package RateLimiter

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/burhon94/RateLimiter/limiter/limit"

	"golang.org/x/time/rate"
)

// set limmiter allow 1 request every 20 second on from IP
var limiter = SetParam(20, 1)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	addr := "192.168.8.101:9999"
	log.Printf("listining: %s\n", addr)
	if err := http.ListenAndServe(addr, LimitMiddleware(mux)); err != nil {
		log.Fatalf("can' t start: %s", err.Error())
	}
}

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getLimiter := limiter.GetLimiter(getIPFromRemoteAddr(r.RemoteAddr))

		if !getLimiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIPFromRemoteAddr(r.RemoteAddr)
	log.Printf("ip: %s\n", ip)

	w.Write([]byte("status OK"))
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
