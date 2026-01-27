package middleware

import "context"

func GetCorrelationID(ctx context.Context) string {
	id, ok := ctx.Value("correlation-id").(string)
	if !ok {
		return ""
	}

	return id
}
