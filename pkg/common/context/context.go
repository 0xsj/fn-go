package context

import (
	"context"
)

type contextKey string

const (
    UserIDKey contextKey = "user_id"
    RoleKey contextKey = "role"
    RequestIDKey contextKey = "request_id"
)

func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
    v, ok := ctx.Value(UserIDKey).(string)
    return v, ok
}

func WithRole(ctx context.Context, role string) context.Context {
    return context.WithValue(ctx, RoleKey, role)
}

func GetRole(ctx context.Context) (string, bool) {
    v, ok := ctx.Value(RoleKey).(string)
    return v, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
    v, ok := ctx.Value(RequestIDKey).(string)
    return v, ok
}