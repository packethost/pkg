package log

import (
	"github.com/jacobweinstock/rollzap"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RollbarConfig if enabled
type RollbarConfig struct {
	Token   string
	Env     string
	Version string
}

// RollbarLogger for implementing the rollbar client logger
type RollbarLogger struct {
	*zap.Logger
}

func (c RollbarConfig) setupRollbar(service string, logger *zap.Logger) zap.Option {
	rollbar.SetToken(c.Token)
	rollbar.SetEnvironment(c.Env)
	rollbar.SetCodeVersion(c.Version)
	rollbar.SetServerRoot("/" + service)
	rollbar.SetLogger(RollbarLogger{logger})

	rollbarCore := rollzap.NewRollbarCore(zapcore.ErrorLevel)
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, rollbarCore)
	})
}

// Printf for internal rollbar errors
func (r RollbarLogger) Printf(format string, args ...interface{}) {
	r.Sugar().Infof(format, args...)
}
