package rollbar

import (
	"fmt"
	"os"
	"strconv"
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
	err error
}

func (e rError) Error() string {
	return e.err.Error()
}

func (e rError) Cause() error {
	return errors.Cause(error(e.err))
}

// Stack converts a github.com/pkg/errors Error stack into a rollbar stack
func (e rError) Stack() rollbar.Stack {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	cause := e.Cause()
	err, ok := cause.(stackTracer)
	if !ok {
		log.Errorw("error does not implement StackTracer", "error", cause)
		return nil
	}

	stack := err.StackTrace()
	rStack := rollbar.Stack(make([]rollbar.Frame, len(stack)))

	var b strings.Builder
	for i := range stack {
		fmt.Fprintf(&b, "%s", stack[i])
		rStack[i].Filename = b.String()

		fmt.Fprintf(&b, "%n", stack[i])
		rStack[i].Method = b.String()

		fmt.Fprintf(&b, "%d", stack[i])
		d, err := strconv.Atoi(b.String())
		if err != nil {
			log.Errorw("failed to convert frame line number to int", "lineString", b.String())
			return nil
		}
		rStack[i].Line = d
		b.Reset()
	}

	return rStack
}

func Notify(err error, args ...interface{}) {
	rErr := rError{err: err}
	rollbar.Error(rErr)
}
