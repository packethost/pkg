package log

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
	return func(args *PacketLogger) { args.LogLevel = level }
}

// WithOutputPaths adds output paths
func WithOutputPaths(paths []string) LoggerOption {
	return func(args *PacketLogger) { args.OutputPaths = paths }
}

// WithServiceName adds a service name a logged field
func WithServiceName(name string) LoggerOption {
	return func(args *PacketLogger) { args.ServiceName = name }
}

// WithKeysAndValues adds extra key/value fields
func WithKeysAndValues(kvs map[string]interface{}) LoggerOption {
	return func(args *PacketLogger) { args.KeysAndValues = kvs }
}

// WithEnableErrLogsToStderr sends .Error logs to stderr
func WithEnableErrLogsToStderr(enable bool) LoggerOption {
	return func(args *PacketLogger) { args.EnableErrLogsToStderr = enable }
}

// WithEnableRollbar sends error logs to Rollbar service
func WithEnableRollbar(enable bool) LoggerOption {
	return func(args *PacketLogger) { args.EnableRollbar = enable }
}

// WithRollbarConfig customizes the Rollbar details
func WithRollbarConfig(config RollbarConfig) LoggerOption {
	return func(args *PacketLogger) { args.RollbarConfig = config }
}

// PacketLogger is a wrapper around zap.SugaredLogger
type PacketLogger struct {
	logr.Logger
	LogLevel              string
	OutputPaths           []string
	ServiceName           string
	KeysAndValues         map[string]interface{}
	EnableErrLogsToStderr bool
	EnableRollbar         bool
	RollbarConfig         RollbarConfig
}

// LoggerOption for setting optional values
type LoggerOption func(*PacketLogger)

// NewPacketLogger is the opionated packet logger setup
func NewPacketLogger(opts ...LoggerOption) (logr.Logger, *zap.Logger, error) {
	// defaults
	const (
		defaultLogLevel    = "info"
		defaultServiceName = "not/set"
	)
	var (
		defaultOutputPaths   = []string{"stdout"}
		defaultKeysAndValues = map[string]interface{}{"service": defaultServiceName}
		zapConfig            = zap.NewProductionConfig()
		zLevel               = zap.InfoLevel
		rollbarOptions       zap.Option
		defaultRollbarConfig = RollbarConfig{
			Token:   "123",
			Env:     "production",
			Version: "1",
		}
	)

	pl := &PacketLogger{
		Logger:        nil,
		LogLevel:      defaultLogLevel,
		OutputPaths:   defaultOutputPaths,
		ServiceName:   defaultServiceName,
		KeysAndValues: defaultKeysAndValues,
		EnableRollbar: false,
		RollbarConfig: defaultRollbarConfig,
	}

	for _, opt := range opts {
		opt(pl)
	}

	switch pl.LogLevel {
	case "debug":
		zLevel = zap.DebugLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(zLevel)
	zapConfig.OutputPaths = pl.OutputPaths
	zapConfig.OutputPaths = sliceDedupe(append(zapConfig.OutputPaths, "stdout"))
	zapConfig.InitialFields = pl.KeysAndValues

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return pl, zapLogger, errors.Wrap(err, "failed to build logger config")
	}

	if pl.EnableRollbar {
		rollbarOptions = pl.RollbarConfig.setupRollbar(pl.ServiceName, zapLogger)
		zapLogger = zapLogger.WithOptions(rollbarOptions)
	}
	if pl.EnableErrLogsToStderr {
		zapLogger = zapLogger.WithOptions(errLogsToStderr(zapConfig))
	}

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
