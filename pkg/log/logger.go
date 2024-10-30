package logger

import (
	"context"
	"fmt"
	"go-actors/pkg/utils"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	scopeFieldName   = "scope"
	traceIdFieldName = "trace_id"
)

var logger zerolog.Logger

func GetCtxLogger(ctx context.Context) zerolog.Logger {
	return logger.With().Ctx(ctx).Logger()
}

func InitLogger(scope string, debug bool) (context.Context, context.CancelFunc) {
	ctx, cancelFn := context.WithCancel(context.Background())

	partsOrder := []string{
		zerolog.LevelFieldName,
		zerolog.TimestampFieldName,
		traceIdFieldName,
		scopeFieldName,
		zerolog.MessageFieldName,
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		PartsOrder: partsOrder,
		FormatPrepare: func(m map[string]any) error {
			formatFieldValue[string](m, "%s", traceIdFieldName)
			formatFieldValue[string](m, "[%s]", scopeFieldName)
			return nil
		},
		FieldsExclude: []string{traceIdFieldName, scopeFieldName},
	}

	logger = zerolog.New(consoleWriter).Hook(ctxHook{})
	if debug {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}
	logger = logger.With().Timestamp().Logger()

	ctx = utils.GetCtxWithScope(ctx, scope)
	return ctx, cancelFn
}

func formatFieldValue[T any](vs map[string]any, format string, field string) {
	if v, ok := vs[field].(T); ok {
		vs[field] = fmt.Sprintf(format, v)
	} else {
		vs[field] = ""
	}
}

type ctxHook struct{}

func (h ctxHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	if scope, ok := utils.GetScopeFromCtx(e.GetCtx()); ok {
		e.Str(scopeFieldName, scope)
	}
	if traceId, ok := utils.GetTraceIdFromCtx(e.GetCtx()); ok {
		e.Str(traceIdFieldName, traceId)
	}
}
