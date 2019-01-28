// package log sets up a shared zap.Logger that can be used by all packages.
package log

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	logLevel = zap.LevelFlag("log-level", zap.InfoLevel, "Log level, one of ERROR, INFO, or DEBUG")
)

// Logger is a wrapper around zap.SugaredLogger
type Logger struct {
	s *zap.SugaredLogger
}

// Init initializes the logging system and sets the "service" key to the provided argument.
// This func should only be called once and after flag.Parse() has been called otherwise leveled logging will not be configured correctly.
func Init(service string) (Logger, func(), error) {
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

	l = l.With(zap.String("service", service))


	cleanup := func() {
		l.Sync()
	}

	return Logger{s: l.Sugar()}, cleanup, nil
}

func (l Logger) Error(err error, args ...interface{}) {
	l.s.With("error", err).Error(args)
}
func (l Logger) Info(args ...interface{}) {
	l.s.Info(args)
}
func (l Logger) Debug(args ...interface{}) {
	l.s.Debug(args)
}
func (l Logger) With(args ...interface{}) Logger {
	return Logger{s: l.s.With(args)}
}

func (l Logger) Package(pkg string) Logger {
	return Logger{s: l.s.With("pkg", pkg)}
}
