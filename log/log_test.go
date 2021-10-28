// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"os"
	"testing"

	"github.com/packethost/pkg/internal/testenv"
	assert "github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogging(t *testing.T) {
	errorMessage := "the flobnarm overheated"
	tests := []struct {
		level    zapcore.Level
		levels   []zapcore.Level
		messages []string
	}{
		{zap.DebugLevel, []zapcore.Level{zap.DebugLevel, zap.InfoLevel, zap.ErrorLevel, zap.ErrorLevel}, []string{"debug", "info", "oh no an error", errorMessage}},
		{zap.InfoLevel, []zapcore.Level{zap.InfoLevel, zap.ErrorLevel, zap.ErrorLevel}, []string{"info", "oh no an error", errorMessage}},
		{zap.ErrorLevel, []zapcore.Level{zap.ErrorLevel, zap.ErrorLevel}, []string{"oh no an error", errorMessage}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.level), func(t *testing.T) {
			enabler := zap.NewAtomicLevelAt(tt.level)
			core, logs := observer.New(enabler)

			service := fmt.Sprintf("testing-%v", tt.level)
			logger, err := configureLogger(zap.New(core), service)
			if err != nil {
				t.Fatal(err)
			}

			defer logger.Close()

			logger.Debug("debug")
			logger.Info("info")
			logger.Error(fmt.Errorf(errorMessage), "oh no an error")
			logger.Error(fmt.Errorf(errorMessage))

			if logs.Len() != len(tt.messages) {
				t.Fatalf("unexpected number of messages: want=%d, got=%d", len(tt.messages), logs.Len())
			}

			for i, log := range logs.All() {
				if log.Level != tt.levels[i] {
					t.Fatalf("unexpected log level: want=%v, got=%v", tt.levels[i], log.Level)
				}

				msg := tt.messages[i]
				got := log.Message
				if got != msg {
					t.Fatalf("unexpected message:\nwant=|%s|\n got=|%s|", msg, got)
				}

				contexts := map[string]string{
					"service": service,
				}
				if log.Level == zap.ErrorLevel {
					contexts["error"] = errorMessage
				}

				ctx := log.ContextMap()
				if len(ctx) != len(contexts) {
					t.Fatalf("unexpected number of contexts: want=%d, got=%d", len(contexts), len(ctx))
				}

				for k, wantV := range contexts {
					gotV, ok := ctx[k]
					if !ok {
						t.Fatalf("missing key in context: key=%s contexts:%v", k, ctx)
					}
					if gotV != wantV {
						t.Fatalf("unexpected value for service context:\nwant=%s\n got=%s", wantV, gotV)
					}
				}
			}
		})
	}
}

func TestContext(t *testing.T) {
	enabler := zap.NewAtomicLevelAt(zap.InfoLevel)
	core, logs := observer.New(enabler)

	service := fmt.Sprintf("testing-%v", zap.InfoLevel)
	logger1, err := configureLogger(zap.New(core), service)
	if err != nil {
		t.Fatal(err)
	}
	defer logger1.Close()

	assertMapsEqual := func(want, got map[string]interface{}) {
		if len(want) != len(got) {
			t.Fatalf("unexpected number of contexts: want=%d, got=%d", len(want), len(got))
		}
		for k := range want {
			vW := want[k]
			vG, ok := got[k]
			if !ok {
				t.Fatalf("missing key in context: key=%s contexts:%v", k, got)
			}
			if vW != vG {
				t.Fatalf("unexpected value for service context: want=%s, got=%s", vW, vG)
			}
		}
	}

	contexts := map[string]interface{}{
		"service": service,
	}

	logger1.Info("logger1 info")
	msgs := logs.All()

	want := 1
	if len(msgs) != want {
		t.Fatalf("unexpected number of messages: want=%d, got=%d", want, len(msgs))
	}

	assertMapsEqual(contexts, msgs[0].ContextMap())

	logger2 := logger1.With("foo", "bar")
	logger1.Info("logger1 info2")
	logger2.Info("logger2 info")
	logger1.Package("logger1").Info("packaged1 info")
	logger2.Package("logger2").Info("packaged2 info")
	logger1.Info("logger1 info3")

	msgs = logs.All()
	want = 6
	if len(msgs) != want {
		t.Fatalf("unexpected number of messages: want=%d, got=%d", want, len(msgs))
	}

	assertMapsEqual(contexts, msgs[0].ContextMap()) // hasn't changed
	assertMapsEqual(contexts, msgs[1].ContextMap())
	contexts["foo"] = "bar"
	assertMapsEqual(contexts, msgs[2].ContextMap())
	delete(contexts, "foo")
	contexts["pkg"] = "logger1"
	assertMapsEqual(contexts, msgs[3].ContextMap())
	contexts["foo"] = "bar"
	contexts["pkg"] = "logger2"
	assertMapsEqual(contexts, msgs[4].ContextMap())
	delete(contexts, "foo")
	delete(contexts, "pkg")
	assertMapsEqual(contexts, msgs[5].ContextMap())

	for i, msg := range []string{"logger1 info", "logger1 info2", "logger2 info", "packaged1 info", "packaged2 info", "logger1 info3"} {
		got := msgs[i].Message
		if got != msg {
			t.Fatalf("unexpected message: want=%s, got=%s", msg, got)
		}
	}
}

func TestInit(t *testing.T) {
	defer testenv.Clear().Restore()

	t.Run("defaults", func(t *testing.T) {
		_, _ = Init("defaults")
	})

	os.Setenv("LOG_DISCARD_LOGS", "true")

	t.Run("debug mode", func(t *testing.T) {
		l, _ := Init("debug")
		l.s.DPanic("this should not panic")

		os.Setenv("DEBUG", "1")
		defer os.Unsetenv("DEBUG")

		l, _ = Init("debug")

		assert.Panics(t, func() {
			// we don't actually care about DPanic, but its a nice/easy indicator that zap was setup for debugging
			l.s.DPanic("this should panic")
		})

	})

	// Rest of the tests will test the rollbar setup
	os.Setenv("ROLLBAR_TOKEN", "TEST-TOKEN")
	os.Setenv("ROLLBAR_DISABLE", "1")

	t.Run("RollbarMissingVersionPanics", func(t *testing.T) {
		// ensure Init fails if none of the *ENV vars are set
		// we need one of VERSION to exist
		os.Setenv("VERSION", "VERSION")
		defer os.Unsetenv("VERSION")

		for _, env := range []string{
			"ENV",
			"EQUINIX_ENV",
			"PACKET_ENV",
		} {
			t.Run(env, func(t *testing.T) {
				os.Setenv(env, env)
				_, _ = Init(env + " set: should-not-fail")

				os.Unsetenv(env)
				assert.PanicsWithValue(t, "required envvar(ENV) is unset", func() {
					_, _ = Init(env + " unset: should-fail")
				})
			})
		}
	})

	t.Run("RollbarMissingVersionPanics", func(t *testing.T) {
		// ensure Init fails if none of the *VERSION vars are set
		// we need one of ENV to exist
		os.Setenv("ENV", "ENV")
		defer os.Unsetenv("ENV")

		for _, env := range []string{
			"VERSION",
			"EQUINIX_VERSION",
			"PACKET_VERSION",
		} {
			t.Run(env, func(t *testing.T) {
				os.Setenv(env, env)
				_, _ = Init(env + " set: should-not-fail")

				os.Unsetenv(env)
				assert.PanicsWithValue(t, "required envvar(VERSION) is unset", func() {
					_, _ = Init(env + " unset: should-fail")
				})
			})
		}
	})
}

func TestFatal(t *testing.T) {
	enabler := zap.NewAtomicLevelAt(zap.InfoLevel)
	core, logs := observer.New(enabler)

	logger, _ := configureLogger(zap.New(core), "TestFatal")
	defer logger.Close()

	msg := "an error"
	want := fmt.Errorf(msg)
	defer func() {
		iface := recover()
		if iface == nil {
			t.Fatal("expected a non-nil return from recover()")
		}
		err, ok := iface.(error)
		if !ok {
			t.Fatalf("unexpected return from recover() want: error, got:%T", iface)
		}
		if err != want {
			t.Fatalf("error mismatch, want: %v, got: %v", want, err)
		}
		if logs.Len() != 1 {
			t.Fatalf("log message mismatch, want: %v, got: %v", 1, logs.Len())
		}
		log := logs.All()[0]
		level := zapcore.ErrorLevel
		if log.Level != level {
			t.Fatalf("log level mismatch want: %v, got: %v", level, log.Level)
		}
		if log.Message != msg {
			t.Fatalf("log message mismatch want: %s, got: %s", msg, log.Message)
		}
	}()
	logger.Fatal(want)
	t.Fatal("should have panic'ed before getting here")
}
