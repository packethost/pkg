// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package rollbar

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	rollbarerr "github.com/rollbar/rollbar-go/errors"
)

func TestStack(t *testing.T) {
	pc := make([]uintptr, 10)

	err := errors.New("hi")
	n := runtime.Callers(1, pc)

	stack, ok := rollbarerr.StackTracer(err)
	if !ok {
		t.Fatalf("rollbarerr.StackTracer returned false")
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
			t.Fatalf("filename mismatch: i=%d\nwant=%s\n got=%s\n", i, want.File, got.File)
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

func TestPkgErrorsCompat(t *testing.T) {
	for name, err := range map[string]error{
		"Errorf":                errors.Errorf("Errorf"),
		"New":                   errors.New("New"),
		"WithStack(fmt.Errorf)": errors.WithStack(fmt.Errorf("fmt.Errorf")),
		"WithStack(errors.New)": errors.WithStack(errors.New("New")),
		"Wrap(fmt.Errorf)":      errors.Wrap(fmt.Errorf("fmt.Errorf"), "Wrap fmt.Errrof"),
		"Wrap(errors.New)":      errors.Wrap(errors.New("New"), "Wrap errors.New"),
	} {
		t.Run(name, func(t *testing.T) {
			stack, ok := rollbarerr.StackTracer(err)
			if !ok {
				t.Fatal("rollbarerr.StackTracer returned false")
			}
			if stack == nil {
				t.Fatalf("expected err to implement rollbar.StackTracer but does not: %v", err)
			}
		})
	}
}
