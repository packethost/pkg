package logr

import (
	"github.com/jacobweinstock/rollzap"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// rollbarConfig if enabled
type rollbarConfig struct {
	token   string
	env     string
	version string
}

// rollbarLogger for implementing the rollbar client logger
type rollbarLogger struct {
	*zap.Logger
}

func (c rollbarConfig) setupRollbar(service string, logger *zap.Logger) zap.Option {
	rollbar.SetToken(c.token)
	rollbar.SetEnvironment(c.env)
	rollbar.SetCodeVersion(c.version)
	rollbar.SetServerRoot(service)
	rollbar.SetLogger(rollbarLogger{logger})

	rollbarCore := rollzap.NewRollbarCore(zapcore.ErrorLevel)
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, rollbarCore)
	})
}

// Printf for internal rollbar errors
func (r rollbarLogger) Printf(format string, args ...interface{}) {
	r.Sugar().Infof(format, args...)
}
