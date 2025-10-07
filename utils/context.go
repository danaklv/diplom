package utils

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey = contextKey("userID")

// Сохраняем userID в контексте
func ContextWithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// Достаём userID из контекста
func UserIDFromContext(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, errors.New("userID not found in context")
	}
	return id, nil
}
