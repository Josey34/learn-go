package middleware

import "net/http"

type RateLimiter struct {
	queue chan struct{}
	limit int
}

func NewRateLimiter(limit int) *RateLimiter {
	queue := make(chan struct{}, limit)

	return &RateLimiter{
		queue: queue,
		limit: limit,
	}
}

func (rl *RateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.queue <- struct{}{}
		next.ServeHTTP(w, r)
		<-rl.queue
	})
}
