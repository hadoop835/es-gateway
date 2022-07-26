package log

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func Log() *zap.Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger
}

type ctxLogKeyType struct{}

// CtxLogKey indicates the context key for logger
// public for test usage.
var CtxLogKey = ctxLogKeyType{}

// Logger gets a contextual logger from current context.
// contextual logger will output common fields from context.
func Logger(ctx context.Context) *zap.Logger {
	if ctxlogger, ok := ctx.Value(CtxLogKey).(*zap.Logger); ok {
		return ctxlogger
	}
	return Log()
}

// WithConnID attaches connId to context.
func WithConnID(ctx context.Context, connID uint64) context.Context {
	var logger *zap.Logger
	if ctxLogger, ok := ctx.Value(CtxLogKey).(*zap.Logger); ok {
		logger = ctxLogger
	} else {
		logger = Log()
	}
	return context.WithValue(ctx, CtxLogKey, logger.With(zap.Uint64("conn", connID)))
}
