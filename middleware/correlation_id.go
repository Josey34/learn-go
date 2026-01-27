package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New().String()

		ctx = context.WithValue(ctx, "correlation-id", id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
