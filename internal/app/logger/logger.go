package logger

import "go.uber.org/zap"

var log *zap.Logger = zap.NewNop()

func Log() *zap.Logger {
	return log
}

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	defer zl.Sync()

	log = zl
	return nil
}
