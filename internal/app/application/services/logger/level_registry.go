package logger

import "go.uber.org/zap/zapcore"

type LevelRegistry struct {
	levels map[string]zapcore.Level
}

func NewLevelRegistry() *LevelRegistry {
	return &LevelRegistry{
		levels: map[string]zapcore.Level{
			"debug": zapcore.DebugLevel,
			"info":  zapcore.InfoLevel,
			"warn":  zapcore.WarnLevel,
			"error": zapcore.ErrorLevel,
		},
	}
}

func (lr *LevelRegistry) Get(name string) zapcore.Level {
	if l, ok := lr.levels[name]; ok {
		return l
	}

	return zapcore.InfoLevel
}
