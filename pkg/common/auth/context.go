package auth

import (
	"context"
)

type contextKey string

const userContextKey contextKey = "user"

func ContextWithUser(ctx context.Context, userContext map[string]interface{}) context.Context {
	return context.WithValue(ctx, userContextKey, userContext)
}

func UserFromContext(ctx context.Context) map[string]interface{} {
	userContext, ok := ctx.Value(userContextKey).(map[string]interface{})
	if !ok {
		return nil
	}
	return userContext
}
