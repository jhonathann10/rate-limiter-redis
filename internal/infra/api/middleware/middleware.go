package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jhonathann10/rate-limiter-redis/internal/infra/internalerrors"
)

type RateLimiter struct {
	visitorsToken map[string]*VisitorToken
	visitorsIp    map[string]*VisitorIp
	mu            sync.Mutex
	rateToken     time.Duration
	rateIp        time.Duration
	burstToken    int
	burstIp       int
}

type VisitorToken struct {
	limiter  *time.Ticker
	lastSeen time.Time
	count    int
}

type VisitorIp struct {
	limiter  *time.Ticker
	lastSeen time.Time
	count    int
}

func NewRateLimiter(rateToken, rateIp time.Duration, burstToken, burstIp int) *RateLimiter {
	rl := &RateLimiter{
		visitorsToken: make(map[string]*VisitorToken),
		visitorsIp:    make(map[string]*VisitorIp),
		rateToken:     rateToken,
		rateIp:        rateIp,
		burstToken:    burstToken,
		burstIp:       burstIp,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) getVisitor(token, ip string) (*VisitorToken, *VisitorIp) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	vToken, existsToken := rl.visitorsToken[token]
	vIp, existsIp := rl.visitorsIp[ip]

	if !existsToken {
		limiter := time.NewTicker(rl.rateToken)
		rl.visitorsToken[token] = &VisitorToken{limiter: limiter, lastSeen: time.Now(), count: 0}
		vToken = rl.visitorsToken[token]
	}

	if !existsIp {
		limiter := time.NewTicker(rl.rateIp)
		rl.visitorsIp[ip] = &VisitorIp{limiter: limiter, lastSeen: time.Now(), count: 0}
		vIp = rl.visitorsIp[ip]
	}

	return vToken, vIp
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		rl.mu.Lock()
		for token, v := range rl.visitorsToken {
			if time.Since(v.lastSeen) > rl.rateToken {
				log.Printf("deleting token %s", token)
				delete(rl.visitorsToken, token)
			}
		}
		for ip, v := range rl.visitorsIp {
			if time.Since(v.lastSeen) > rl.rateIp {
				log.Printf("deleting ip %s", ip)
				delete(rl.visitorsIp, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token := r.Header.Get("Authorization")
		ip := r.RemoteAddr

		visitorToken, visitorIp := rl.getVisitor(token, ip)
		visitorToken.count++
		if visitorToken.count > rl.burstToken {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(getToManyRequestsError("Token rate limit exceeded"))
			return
		}
		visitorIp.count++
		if visitorIp.count > rl.burstIp {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(getToManyRequestsError("IP rate limit exceeded"))
			return
		}

		select {
		case <-visitorToken.limiter.C:
			visitorToken.count = 0
		case <-visitorIp.limiter.C:
			visitorIp.count = 0
		default:
		}

		next.ServeHTTP(w, r)
	})
}

func getToManyRequestsError(errMsg string) *internalerrors.InternalError {
	return &internalerrors.InternalError{
		Message: errMsg,
		Status:  http.StatusTooManyRequests,
	}
}
