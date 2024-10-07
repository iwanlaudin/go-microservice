package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
)

type contextKey string

const userContextKey contextKey = "userContext"

type UserContext struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func ContextWithUser(ctx context.Context, userContext *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, userContext)
}

func UserFromContext(ctx context.Context) *UserContext {
	userContext, ok := ctx.Value(userContextKey).(*UserContext)
	if !ok {
		panic("Context value is nil")
	}
	return userContext
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	user, ok := ctx.Value(userContextKey).(*UserContext)
	if !ok {
		return uuid.Nil
	}
	id, _ := helpers.ConvertStringToUUID(user.ID)
	return id
}

func UserEmailFromContext(ctx context.Context) string {
	user, ok := ctx.Value(userContextKey).(*UserContext)
	if !ok {
		return ""
	}
	return user.Email
}
