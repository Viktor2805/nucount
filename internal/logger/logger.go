package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *zap.Logger
)

// Init initializes the global logger
func Init() *zap.Logger {
	once.Do(func() {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		cfg.Encoding = "json"
		cfg.OutputPaths = []string{"stdout"}

		l, err := cfg.Build()

		if err != nil {
			panic(err)
		}

		logger = l
	})

	return logger
}

func L() *zap.Logger {
	if logger == nil {
		return Init()
	}

	return logger
}
