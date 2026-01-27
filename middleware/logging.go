package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		correlationID := GetCorrelationID(r.Context())
		log.Printf("[%s]Started %s %s", correlationID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[%s]Completed %s %s in %v", correlationID, r.Method, r.URL.Path, duration)
	})
}
