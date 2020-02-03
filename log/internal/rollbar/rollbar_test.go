package rollbar

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestStack(t *testing.T) {
	pc := make([]uintptr, 10)

	err := errors.New("hi")
	n := runtime.Callers(1, pc)

	stack := []runtime.Frame(rError{err: err}.Stack())
	if n == 0 {
		t.Fatal("unable to recover stack information")
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	i := 0
	for {
		want, more := frames.Next()
		if i >= len(stack) {
			t.Fatalf("stack frame count mismatch, index=%d, gotLength=%d", i, len(stack))
		}
		got := stack[i]
		if i == 0 {
			want.Line -= 1 // account for calling runtime.Callers on line after errors.New
		}

		t.Logf("want: method=%v %v:%v", want.Function, want.File, want.Line)
		t.Logf(" got: method=%v %v:%v", got.Function, got.File, got.Line)
		if want.File != got.File {
			t.Fatalf("filename mismatch: i=%d\nwant=%s\n got=%s\n", i, want.File, got.Function)
		}
		if want.Func.Name() != got.Function {
			t.Fatalf("func name mismatch: i=%d\nwant=%s\n got=%s\n", i, want.Func.Name(), got.Function)
		}
		if want.Line != got.Line {
			t.Fatalf("line number mismatch: i=%d\nwant=%d\n got=%d\n", i, want.Line, got.Line)
		}

		if !more {
			break
		}
		i++
	}

	if i < len(stack)-1 {
		t.Fatalf("stack frame count mismatch, index=%d, got=%d", i, len(stack))
	}
}

func TestError(t *testing.T) {
	err := errors.New("hi")
	rErr := rError{err: err}
	if rErr.Error() != err.Error() {
		t.Fatalf("Error() mismatch:\nwant=%s\n got=%s\n", err.Error(), rErr.Error())
	}
}

func TestCause(t *testing.T) {
	err := errors.New("hi")
	rErr := rError{err: err}
	if rErr.Cause() != err {
		t.Fatalf("Cause() mismatch:\nwant=%v\n got=%v\n", err, rErr.Cause())
	}
}

func TestPkgErrorsCompat(t *testing.T) {
	log = zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel)).Sugar()
	for _, err := range []error{
		errors.Errorf("Errorf"),
		errors.New("New"),
		errors.WithMessage(errors.New("New"), "WithMessage errors.New"),
		errors.WithStack(fmt.Errorf("fmt.Errorf")),
		errors.WithStack(errors.New("New")),
		errors.Wrap(fmt.Errorf("fmt.Errorf"), "Wrap fmt.Errrof"),
		errors.Wrap(errors.New("New"), "Wrap errors.New"),
	} {
		rErr := rError{err: err}
		if rErr.Stack() == nil {
			t.Fatalf("expected err to implement stackTracer but does not: %v", err)
		}
	}
}
