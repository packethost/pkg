package rollbar

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	rollbar "github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Setup(l *zap.SugaredLogger, service string) func() {
	log = l

	token := os.Getenv("ROLLBAR_TOKEN")
	if token == "" {
		log.Panicw("required envvar is unset", "envvar", "ROLLBAR_TOKEN")
	}
	rollbar.SetToken(token)

	env := os.Getenv("PACKET_ENV")
	if env == "" {
		log.Panicw("required envvar is unset", "envvar", "PACKET_ENV")
	}
	rollbar.SetEnvironment(env)

	v := os.Getenv("PACKET_VERSION")
	if v == "" {
		log.Panicw("required envvar is unset", "envvar", "PACKET_VERSION")
	}
	rollbar.SetCodeVersion(v)
	rollbar.SetServerRoot(service)

	enable := true
	if os.Getenv("ROLLBAR_DISABLE") != "" {
		enable = false
	}
	rollbar.SetEnabled(enable)

	return rollbar.Wait
}

// rError exists to implement rollbar.CauseStacker so that rollbar can have stack info.
// see https://github.com/rollbar/rollbar-go/blob/v1.0.2/doc.go#L64
type rError struct {
	service string
	err     error
}

func (e rError) Error() string {
	return e.err.Error()
}

func (e rError) Cause() error {
	return e.err
}

// logInternalError is a helper to log errors through zap and to rollbar if we run into an error while logging a client's error.
// We can use rollbar.ErrorWithExtras here because the stack trace rollbar collects will be of where error is.
// This handles the so called "error while logging error" case.
func logInternalError(err error, ctx map[string]interface{}) {
	l := log.With("error", err)
	if len(ctx) != 0 {
		fields := make([]interface{}, 0, len(ctx)*2)
		for k, v := range ctx {
			fields = append(fields, k)
			fields = append(fields, v)
		}
		l = l.With(fields...)
	}
	ctx["errorVerbose"] = fmt.Sprintf("%+v", err)
	l.Error(err)
	// 1 level of stack frames are skipped, because we don't want care to have logInternalError show up
	rollbar.ErrorWithStackSkipWithExtras(rollbar.ERR, err, 1, ctx)
}

// shortenFilePath removes un-needed information from the source file path.
// This makes them shorter in Rollbar UI as well as making them the same, regardless of the machine the code was compiled on.
func shortenFilePath(service, s string) string {
	idx := strings.Index(s, service)
	if idx != -1 {
		return s[idx:]
	}
	return s
}

// Stack converts a github.com/pkg/errors Error stack into a rollbar stack
func (e rError) Stack() rollbar.Stack {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	ctx := map[string]interface{}{}

	cause := e.Cause()
	st, ok := cause.(stackTracer)
	if !ok {
		ctx["cause"] = cause
		logInternalError(errors.New("cause does not implement StackTracer"), ctx)
		return nil
	}

	stack := st.StackTrace()
	rStack := rollbar.Stack(make([]rollbar.Frame, len(stack)))

	var b strings.Builder
	for i := range stack {
		b.Reset()
		fmt.Fprintf(&b, "%+s", stack[i])
		var filename string
		n, err := fmt.Sscanf(b.String(), "%s\n\t%s", &rStack[i].Method, &filename)
		rStack[i].Filename = shortenFilePath(e.service, filename)

		if err != nil {
			ctx["lineString"] = b.String()
			logInternalError(errors.Wrap(err, "failed to scan stack frame"), ctx)
			return nil
		}
		if n != 2 {
			ctx["lineString"] = b.String()
			ctx["count"] = n
			logInternalError(errors.Wrap(err, "unexpected number of values scanned when scanning for stack frame func and file names"), ctx)
			return nil
		}

		b.Reset()
		fmt.Fprintf(&b, "%d", stack[i])
		n, err = fmt.Sscanf(b.String(), "%d", &rStack[i].Line)
		if err != nil {
			ctx["lineString"] = b.String()
			logInternalError(errors.Wrap(err, "failed to scan stack frame line number"), ctx)
			return nil
		}
		if n != 1 {
			ctx["lineString"] = b.String()
			ctx["count"] = n
			logInternalError(errors.Wrap(err, "unexpected number of values scanned when scanning for stack frame line number"), ctx)
			return nil
		}
	}

	return rStack
}

func Notify(service string, err error, args ...interface{}) {
	rErr := rError{service: service, err: err}
	rollbar.Error(rErr)
}
