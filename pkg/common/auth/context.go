package auth

import (
	"context"
)

type contextKey string

const userContextKey contextKey = "user"

// ContextWithUser menambahkan user ke dalam context
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext mengambil user dari context
func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}
