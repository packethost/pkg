package logr

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithLogLevel sets the log level
func WithLogLevel(level string) LoggerOption {
	return func(args *PacketLogr) { args.logLevel = level }
}

// WithOutputPaths adds output paths
func WithOutputPaths(paths []string) LoggerOption {
	return func(args *PacketLogr) { args.outputPaths = paths }
}

// WithServiceName adds a service name a logged field
func WithServiceName(name string) LoggerOption {
	return func(args *PacketLogr) { args.serviceName = name }
}

// WithKeysAndValues adds extra key/value fields
func WithKeysAndValues(kvs []interface{}) LoggerOption {
	return func(args *PacketLogr) { args.keysAndValues = append(args.keysAndValues, kvs...) }
}

// WithEnableErrLogsToStderr sends .Error logs to stderr
func WithEnableErrLogsToStderr(enable bool) LoggerOption {
	return func(args *PacketLogr) { args.enableErrLogsToStderr = enable }
}

// WithEnableRollbar sends error logs to Rollbar service
func WithEnableRollbar(enable bool) LoggerOption {
	return func(args *PacketLogr) { args.enableRollbar = enable }
}

// WithRollbarConfig customizes the Rollbar details
func WithRollbarConfig(config rollbarConfig) LoggerOption {
	return func(args *PacketLogr) { args.rollbarConfig = config }
}

// PacketLogr is a wrapper around zap.SugaredLogger
type PacketLogr struct {
	logr.Logger
	logLevel              string
	outputPaths           []string
	serviceName           string
	keysAndValues         []interface{}
	enableErrLogsToStderr bool
	enableRollbar         bool
	rollbarConfig         rollbarConfig
}

// LoggerOption for setting optional values
type LoggerOption func(*PacketLogr)

// NewPacketLogr is the opionated packet logger setup
func NewPacketLogr(opts ...LoggerOption) (logr.Logger, *zap.Logger, error) {
	// defaults
	const (
		defaultLogLevel    = "info"
		defaultServiceName = "not/set"
	)
	var (
		defaultOutputPaths   = []string{"stdout"}
		defaultKeysAndValues = []interface{}{}
		zapConfig            = zap.NewProductionConfig()
		zLevel               = zap.InfoLevel
		defaultZapOpts       = []zap.Option{}
		rollbarOptions       zap.Option
		defaultRollbarConfig = rollbarConfig{
			token:   "123",
			env:     "production",
			version: "1",
		}
	)

	pl := &PacketLogr{
		Logger:        nil,
		logLevel:      defaultLogLevel,
		outputPaths:   defaultOutputPaths,
		serviceName:   defaultServiceName,
		keysAndValues: defaultKeysAndValues,
		enableRollbar: false,
		rollbarConfig: defaultRollbarConfig,
	}

	for _, opt := range opts {
		opt(pl)
	}

	switch pl.logLevel {
	case "debug":
		zLevel = zap.DebugLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(zLevel)
	zapConfig.OutputPaths = sliceDedupe(pl.outputPaths)

	if pl.enableErrLogsToStderr {
		defaultZapOpts = append(defaultZapOpts, errLogsToStderr(zapConfig))
	}

	zapLogger, err := zapConfig.Build(defaultZapOpts...)
	if err != nil {
		return pl, zapLogger, errors.Wrap(err, "failed to build logger config")
	}
	if pl.enableRollbar {
		rollbarOptions = pl.rollbarConfig.setupRollbar(pl.serviceName, zapLogger)
		zapLogger = zapLogger.WithOptions(rollbarOptions)
	}
	keysAndValues := append(pl.keysAndValues, "service", pl.serviceName)
	zapLogger = zapLogger.With(handleFields(zapLogger, keysAndValues)...)
	pl.Logger = zapr.NewLogger(zapLogger)
	return pl, zapLogger, err
}

func sliceDedupe(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func errLogsToStderr(c zap.Config) zap.Option {
	errorLogs := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	nonErrorLogs := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return !errorLogs(lvl)
	})
	console := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	encoder := zapcore.NewJSONEncoder(c.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, console, nonErrorLogs),
		zapcore.NewCore(encoder, consoleErrors, errorLogs),
	)
	splitLogger := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core

	})
	return splitLogger
}

// handleFields converts a bunch of arbitrary key-value pairs into Zap fields.  It takes
// additional pre-converted Zap fields, for use with automatically attached fields, like
// `error`. copy/paste from https://github.com/go-logr/zapr/blob/146009e52d528183a25bf1a1e3cf56d1ff3919b5/zapr.go#L79
func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}
