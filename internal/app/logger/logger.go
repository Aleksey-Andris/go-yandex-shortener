// The configs package is designed to configure the logger.
package logger

import "go.uber.org/zap"

var log *zap.Logger = zap.NewNop()

// Log - returns the global application logger.
func Log() *zap.Logger {
	return log
}

// Initialize - initializes the logger and puts it in the global variable of this package.
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
