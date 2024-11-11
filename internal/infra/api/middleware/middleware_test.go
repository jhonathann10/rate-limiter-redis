package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jhonathann10/rate-limiter-redis/internal/infra/api/middleware"
)

var (
	rateToken     = time.Duration(1) * time.Second
	rateIp        = time.Duration(1) * time.Second
	burstToken    = 100
	burstIp       = 100
	totalRequests = 1000000
)

func TestNewRateLimiter__Success(t *testing.T) {
	rl := middleware.NewRateLimiter(rateToken, rateIp, burstToken, burstIp)
	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/", nil)

	start := time.Now()
	successCount := 0
	for i := 0; i <= totalRequests; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code == http.StatusOK {
			successCount++
		}
	}
	duration := time.Since(start)
	t.Logf("Processed %d requests in %v with %d successful responses", totalRequests, duration, successCount)
}
