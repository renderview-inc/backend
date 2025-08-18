package logger

import (
	"context"
	"github.com/renderview-inc/backend/internal/app/application/services/logger/option"
	"github.com/renderview-inc/backend/internal/app/infrastructure/repositories/logger"
	"github.com/renderview-inc/backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ctxKey string

const (
	CorrelationID ctxKey = "correlation_id"
)

const fileMode = 0644

type LogRepository interface {
	Save(ctx context.Context, log map[string]any) error
}

type LogService struct {
	logger *zap.Logger
	level  zapcore.Level
}

func (l *LogService) newConsoleCore() zapcore.Core {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l.level
	})

	return zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLevel)
}

func (l *LogService) newJSONCore() (zapcore.Core, error) {
	jsonCfg := zap.NewProductionEncoderConfig()
	jsonCfg.TimeKey = "@timestamp"
	jsonCfg.LevelKey = "log.level"
	jsonCfg.MessageKey = "message"
	jsonCfg.CallerKey = "caller"
	jsonCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	jsonEncoder := zapcore.NewJSONEncoder(jsonCfg)

	file, err := os.OpenFile("/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
	if err != nil {
		return nil, err
	}

	jsonLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl != zapcore.DebugLevel
	})

	return zapcore.NewCore(jsonEncoder, zapcore.AddSync(file), jsonLevel), nil
}

func (l *LogService) newDualLogger() (*zap.Logger, error) {
	consoleCore := l.newConsoleCore()

	jsonCore, err := l.newJSONCore()
	if err != nil {
		return nil, err
	}

	core := zapcore.NewTee(consoleCore, jsonCore)

	return zap.New(core, zap.AddCaller()), nil
}

func (l *LogService) Sync() error {
	return l.logger.Sync()
}

func NewLogService(cfg *config.LogConfig) (*LogService, error) {
	var logService LogService

	logService.level = zapcore.InfoLevel
	err := logService.level.Set(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}

	logService.logger, err = logService.newDualLogger()
	if err != nil {
		return nil, err
	}

	return &logService, nil
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
