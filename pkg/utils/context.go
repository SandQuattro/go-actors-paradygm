package utils

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type scopeCtxKey struct{}

func GetCtxWithScope(ctx context.Context, scope string) context.Context {
	return context.WithValue(ctx, scopeCtxKey{}, scope)
}

type traceIdCtxKey struct{}

func GetCtxWithTraceId(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceIdCtxKey{}, generateTraceId())
}

func GetScopeFromCtx(ctx context.Context) (string, bool) {
	if scope, ok := ctx.Value(scopeCtxKey{}).(string); ok {
		return scope, true
	}
	return "", false
}

func GetTraceIdFromCtx(ctx context.Context) (string, bool) {
	if traceId, ok := ctx.Value(traceIdCtxKey{}).(string); ok {
		return traceId, true
	}
	return "", false
}

func generateTraceId() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal().Err(err)
	}

	return newUUID.String()
}
