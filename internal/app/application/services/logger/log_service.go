package logger

import (
	"context"
	"encoding/json"
	"github.com/renderview-inc/backend/internal/app/application/services/logger/option"
	"github.com/renderview-inc/backend/internal/app/infrastructure/repositories/logger"
	"github.com/renderview-inc/backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
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
	repo   LogRepository
	level  zapcore.Level
}

func (l *LogService) newConsoleLogger() *zap.Logger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		l.level,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func NewLogService(cfg *config.LogConfig, repo LogRepository) (*LogService, error) {
	var logService LogService

	logService.level = zapcore.InfoLevel
	err := logService.level.Set(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}

	logService.repo = repo
	logService.logger = logService.newConsoleLogger()

	return &logService, nil
}

func (l *LogService) log(ctx context.Context, level zapcore.Level, msg string, opts ...option.LogOption) {
	additional := make(map[string]any)
	for _, opt := range opts {
		opt(additional)
	}

	correlationID, ok := ctx.Value(CorrelationID).(string)
	if !ok {
		return
	}

	fields := make([]zap.Field, 0, len(additional)+1)
	for k, v := range additional {
		fields = append(fields, zap.Any(k, v))
	}
	if correlationID != "" {
		fields = append(fields, zap.String(logger.CorrelationIDKey, correlationID))
	}

	l.logger.Log(level, msg, fields...)

	if level == zapcore.DebugLevel {
		return
	}

	jsonFields, err := json.Marshal(additional)
	if err != nil {
		l.logger.Error("failed to marshal fields", zap.Error(err))
		jsonFields = []byte("{}")
	}

	log := map[string]any{
		logger.TimestampKey:     time.Now(),
		logger.LevelKey:         level.String(),
		logger.MsgKey:           msg,
		logger.FieldsKey:        string(jsonFields),
		logger.CorrelationIDKey: correlationID,
	}

	if err = l.repo.Save(ctx, log); err != nil {
		l.logger.Error(err.Error())
	}
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
