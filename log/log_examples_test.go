// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"

	"go.uber.org/zap"
)

func setupForExamples(example string) Logger {
	service := "github.com/packethost/pkg"
	c := setupConfig(service)
	c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	c.OutputPaths = []string{"stdout"}
	c.ErrorOutputPaths = c.OutputPaths
	c.EncoderConfig.TimeKey = ""
	z, err := buildConfig(c)
	if err != nil {
		panic(err)
	}
	logger, err := configureLogger(z, service)
	if err != nil {
		panic(err)
	}

	l := logger.Package(example)
	return l
}

func ExampleLogger_Debug() {
	l := setupForExamples("debug")
	defer l.Close()

	l.Debug("debug message")
	//Output:
	//{"level":"debug","caller":"log/log_examples_test.go:36","msg":"debug message","service":"github.com/packethost/pkg","pkg":"debug"}

}

func ExampleLogger_Info() {
	l := setupForExamples("info")
	defer l.Close()

	defer func() {
		_ = recover()
	}()
	l.Info("info message")
	//Output:
	//{"level":"info","caller":"log/log_examples_test.go:49","msg":"info message","service":"github.com/packethost/pkg","pkg":"info"}

}

func ExampleLogger_Error() {
	l := setupForExamples("error")
	defer l.Close()

	l.Error(fmt.Errorf("oh no an error"))
	//Output:
	//{"level":"error","caller":"log/log_examples_test.go:59","msg":"oh no an error","service":"github.com/packethost/pkg","pkg":"error","error":"oh no an error"}

}

func ExampleLogger_Fatal() {
	l := setupForExamples("fatal")
	defer l.Close()

	defer func() {
		_ = recover()
	}()
	l.Fatal(fmt.Errorf("oh no an error"))
	//Output:
	//{"level":"error","caller":"log/log_examples_test.go:72","msg":"oh no an error","service":"github.com/packethost/pkg","pkg":"fatal","error":"oh no an error"}

}

func ExampleLogger_With() {
	l := setupForExamples("with")
	defer l.Close()

	l.With("true", true).Info("info message")
	//Output:
	//{"level":"info","caller":"log/log_examples_test.go:82","msg":"info message","service":"github.com/packethost/pkg","pkg":"with","true":true}

}

func ExampleLogger_Package() {
	l := setupForExamples("info")
	defer l.Close()

	l.Info("info message")
	l = l.Package("package")
	l.Info("info message")
	//Output:
	//{"level":"info","caller":"log/log_examples_test.go:92","msg":"info message","service":"github.com/packethost/pkg","pkg":"info"}
	//{"level":"info","caller":"log/log_examples_test.go:94","msg":"info message","service":"github.com/packethost/pkg","pkg":"info","pkg":"package"}
}
