package logger

import (
	"context"
	"github.com/renderview-inc/backend/pkg/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type LogRepository interface {
	Save(ctx context.Context, log map[string]any) error
}

type LogService struct {
	logger *zap.Logger
	cfg    *config.LogConfig
	repo   LogRepository
}

func NewLogService(repo LogRepository) (*LogService, error) {
	var service LogService
	cfg, err := service.loadConfig()
	if err != nil {
		return nil, err
	}
	service.cfg = cfg

	zapCfg := zap.NewProductionConfig()
	lr := NewLevelRegistry()
	zapCfg.Level = zap.NewAtomicLevelAt(lr.Get(cfg.Logger.Level))

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	service.logger = logger
	service.repo = repo

	return &service, nil
}

func (l *LogService) loadConfig() (*config.LogConfig, error) {
	v := viper.New()

	var cfg config.LogConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
