// package log sets up a shared zap.Logger that can be used by all packages.
package log

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	logger   *zap.Logger
	logLevel = zap.LevelFlag("log-level", zap.InfoLevel, "Log level, one of FATAL, PANIC, DPANIC, ERROR, WARN, INFO, or DEBUG")
)

// Logger is a wrapper around zap.SugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

// Init initializes the logging system and sets the "service" key to the provided argument.
// This func should only be called once and after flag.Parse() has been called otherwise leveled logging will not be configured correctly.
func Init(service string) (Logger, func() error, error) {
	var config zap.Config
	if os.Getenv("DEBUG") != "" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Level = zap.NewAtomicLevelAt(*logLevel)

	l, err := config.Build()
	if err != nil {
		return Logger{}, nil, errors.Wrap(err, "failed to build logger config")
	}

	logger = l.With(zap.String("service", service))

	return Logger{logger.Sugar()}, logger.Sync, nil
}

func (l Logger) Notice(args ...interface{}) {
	l.Info(args)
}

func (l Logger) Trace(args ...interface{}) {
	l.Debug(args)
}

func (l Logger) Warning(args ...interface{}) {
	l.Warn(args)
}

func (l Logger) With(args ...interface{}) Logger {
	return Logger{l.SugaredLogger.With(args)}
}

func (l Logger) Package(pkg string) Logger {
	return Logger{l.SugaredLogger.With("pkg", pkg)}
}
