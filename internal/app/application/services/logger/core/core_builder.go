package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

const (
	LogsDirName = "/usr/share/filebeat/logs"
	LogFileName = "app.log"
)

const fileMode = 0644

type CoreBuilder struct {
	level zapcore.Level
}

func NewCoreBuilder(level zapcore.Level) *CoreBuilder {
	return &CoreBuilder{
		level: level,
	}
}

func (c *CoreBuilder) ConsoleCore() zapcore.Core {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= c.level
	})

	return zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLevel)
}

func (c *CoreBuilder) JSONCore() (zapcore.Core, error) {
	jsonCfg := zap.NewProductionEncoderConfig()
	jsonCfg.TimeKey = "@timestamp"
	jsonCfg.LevelKey = "log.level"
	jsonCfg.MessageKey = "message"
	jsonCfg.CallerKey = "caller"
	jsonCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	jsonEncoder := zapcore.NewJSONEncoder(jsonCfg)

	if err := os.MkdirAll(LogsDirName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create logs dir err: %w", err)
	}

	logFile := filepath.Join(LogsDirName, LogFileName)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, fileMode)
	if err != nil {
		return nil, fmt.Errorf("create log file failed: %w", err)
	}

	jsonLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl != zapcore.DebugLevel
	})

	return zapcore.NewCore(jsonEncoder, zapcore.AddSync(file), jsonLevel), nil
}

func (c *CoreBuilder) DualLogger() (*zap.Logger, error) {
	consoleCore := c.ConsoleCore()

	jsonCore, err := c.JSONCore()
	if err != nil {
		return nil, err
	}

	core := zapcore.NewTee(consoleCore, jsonCore)

	return zap.New(core, zap.AddCaller()), nil
}
