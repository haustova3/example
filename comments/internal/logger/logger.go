package logger

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.SugaredLogger
)

func init() {
	c := zap.NewProductionConfig()
	c.Level.SetLevel(zapcore.InfoLevel)
	c.ErrorOutputPaths = []string{"stdout"}
	c.OutputPaths = []string{"stdout"}
	logger, err := NewLogger(c)
	if err != nil {
		panic(err)
	}
	SetLogger(logger)
}

func NewLogger(config zap.Config) (*zap.SugaredLogger, error) {
	l, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}

func SetLogger(l *zap.SugaredLogger) {
	globalLogger = l
}

func Logger() *zap.SugaredLogger {
	return globalLogger
}

type loggerCtxKeyType string

const (
	loggerCtxKey loggerCtxKeyType = "logger"
)

func ToContext(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "time", time.Now().String())
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		keysAndValues = append(keysAndValues, "trace_id", span.TraceID())
	}
	if span.HasSpanID() {
		keysAndValues = append(keysAndValues, "span_id", span.SpanID())
	}
	if l, ok := ctx.Value(loggerCtxKey).(*zap.SugaredLogger); ok {
		l.Infow(msg, keysAndValues...)

		return
	}
	Logger().Infow(msg, keysAndValues...)
}

func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "time", time.Now().String())
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		keysAndValues = append(keysAndValues, "trace_id", span.TraceID())
	}
	if span.HasSpanID() {
		keysAndValues = append(keysAndValues, "span_id", span.SpanID())
	}
	if l, ok := ctx.Value(loggerCtxKey).(*zap.SugaredLogger); ok {
		l.Warnw(msg, keysAndValues...)

		return
	}

	Logger().Warnw(msg, "time", time.Now().String(), keysAndValues)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "time", time.Now().String())
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		keysAndValues = append(keysAndValues, "trace_id", span.TraceID())
	}
	if span.HasSpanID() {
		keysAndValues = append(keysAndValues, "span_id", span.SpanID())
	}
	if l, ok := ctx.Value(loggerCtxKey).(*zap.SugaredLogger); ok {
		l.Errorw(msg, keysAndValues...)

		return
	}

	Logger().Errorw(msg, keysAndValues...)
}
