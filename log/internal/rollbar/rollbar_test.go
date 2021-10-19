// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package rollbar

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	rollbar "github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestStack(t *testing.T) {
	pc := make([]uintptr, 10)

	err := errors.New("hi")
	n := runtime.Callers(1, pc)

	stack := []rollbar.Frame(rError{err: err}.Stack())
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
		t.Logf(" got: method=%v %v:%v", got.Method, got.Filename, got.Line)
		if want.File != got.Filename {
			t.Fatalf("filename mismatch: i=%d\nwant=%s\n got=%s\n", i, want.File, got.Filename)
		}
		if want.Func.Name() != got.Method {
			t.Fatalf("func name mismatch: i=%d\nwant=%s\n got=%s\n", i, want.Func.Name(), got.Method)
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
	for name, err := range map[string]error{
		"Errorf":                errors.Errorf("Errorf"),
		"New":                   errors.New("New"),
		"WithMessage":           errors.WithMessage(errors.New("New"), "WithMessage errors.New"),
		"WithStack(fmt.Errorf)": errors.WithStack(fmt.Errorf("fmt.Errorf")),
		"WithStack(errors.New)": errors.WithStack(errors.New("New")),
		"Wrap(fmt.Errorf)":      errors.Wrap(fmt.Errorf("fmt.Errorf"), "Wrap fmt.Errrof"),
		"Wrap(errors.New)":      errors.Wrap(errors.New("New"), "Wrap errors.New"),
	} {
		t.Run(name, func(t *testing.T) {
			rErr := rError{err: err}
			if rErr.Stack() == nil {
				t.Fatalf("expected err to implement stackTracer but does not: %v", err)
			}
		})
	}
}
