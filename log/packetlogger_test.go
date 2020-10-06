package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestPacketLogger(t *testing.T) {
	expectedLogMsg := "new logger test message"
	capturedOutput := captureOutput(func() {
		l, _, err := NewPacketLogger()
		if err != nil {
			t.Fatal(err)
		}
		l.V(0).Info(expectedLogMsg)
	})
	if !strings.Contains(capturedOutput, expectedLogMsg) {
		t.Fatalf("expected to contain: %v, got: %v", expectedLogMsg, capturedOutput)
	}
}

func TestPacketLoggerRollbarEnabled(t *testing.T) {
	expectedLogMsg := "new logger test message"

	capturedOutput := captureOutput(func() {
		l, _, err := NewPacketLogger(
			WithLogLevel("debug"),
			WithEnableRollbar(true),
			WithRollbarConfig(RollbarConfig{
				Token:   "badtoken",
				Env:     "production",
				Version: "v2",
			}),
		)
		if err != nil {
			t.Fatal(err)
		}
		l.V(0).Error(errors.New("V0 testing error"), expectedLogMsg)

	})
	if !strings.Contains(capturedOutput, expectedLogMsg) {
		t.Fatalf("expected to contain: %v, got: %v", expectedLogMsg, capturedOutput)
	}
	fmt.Println(capturedOutput)

}

func TestPacketLoggerWithOptions(t *testing.T) {
	expectedLogMsg := "new logger test message"
	capturedOutput := captureOutput(func() {
		l, _, err := NewPacketLogger(
			WithLogLevel("debug"),
			WithOutputPaths([]string{"stdout"}),
			WithServiceName("testing"),
			WithKeysAndValues(map[string]interface{}{"hello": "world"}),
			WithEnableErrLogsToStderr(true),
		)
		if err != nil {
			t.Fatal(err)
		}
		l.V(0).Info(expectedLogMsg)
	})
	if !strings.Contains(capturedOutput, expectedLogMsg) {
		t.Fatalf("expected to contain: %v, got: %v", expectedLogMsg, capturedOutput)
	}
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
