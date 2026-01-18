package middleware

import "net/http"

// Chain composes multiple middlewares into a single middleware
// Middlewares are applied in order (left to right)
// Example: Chain(LoggingMiddleware, RecoveryMiddleware)(handler)
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		// Apply middlewares in reverse order (right to left execution)
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
