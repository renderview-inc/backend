package logger

import (
	"context"
	"github.com/renderview-inc/backend/internal/app/application/services/logger/core"
	"github.com/renderview-inc/backend/internal/app/application/services/logger/option"
	"github.com/renderview-inc/backend/internal/app/infrastructure/repositories/logger"
	"github.com/renderview-inc/backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	CorrelationID ctxKey = "correlation_id"
)

type LogRepository interface {
	Save(ctx context.Context, log map[string]any) error
}

type LogService struct {
	logger *zap.Logger
}

func (l *LogService) Sync() error {
	return l.logger.Sync()
}

func NewLogService(cfg *config.LogConfig) (*LogService, error) {
	level := zapcore.InfoLevel
	if err := level.Set(cfg.Logger.Level); err != nil {
		return nil, err
	}

	builder := core.NewCoreBuilder(level)

	newLogger, err := builder.DualLogger()
	if err != nil {
		return nil, err
	}

	return &LogService{
		logger: newLogger,
	}, nil
}

func (l *LogService) log(ctx context.Context, level zapcore.Level, msg string, opts ...option.LogOption) {
	additional := make(map[string]any)
	for _, opt := range opts {
		opt(additional)
	}

	correlationID, _ := ctx.Value(CorrelationID).(string)

	fields := make([]zap.Field, 0, len(additional)+1)
	for k, v := range additional {
		fields = append(fields, zap.Any(k, v))
	}
	if correlationID != "" {
		fields = append(fields, zap.String(logger.CorrelationIDKey, correlationID))
	}

	l.logger.Log(level, msg, fields...)
}

func (l *LogService) Debug(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.DebugLevel, msg, opts...)
}

func (l *LogService) Info(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.InfoLevel, msg, opts...)
}

func (l *LogService) Warn(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.WarnLevel, msg, opts...)
}

func (l *LogService) Error(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.ErrorLevel, msg, opts...)
}

func (l *LogService) DPanic(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.DPanicLevel, msg, opts...)
}

func (l *LogService) Panic(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.PanicLevel, msg, opts...)
}

func (l *LogService) Fatal(ctx context.Context, msg string, opts ...option.LogOption) {
	l.log(ctx, zapcore.FatalLevel, msg, opts...)
}
