package log

import (
	"context"
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestGetLoggerNoop(t *testing.T) {
	ctx := context.Background()

	GetLogger(ctx).Info("this message should not be visible")
}

func TestContextEmbedding(t *testing.T) {
	bar := func(ctx context.Context) {
		GetLogger(ctx).Info("bar called")
	}

	foo := func(ctx context.Context) {
		logger := GetLogger(ctx)

		logger.Info("foo called")

		ctx = ContextWithLogger(ctx, logger.With("baz", "quux"))
		bar(ctx)
	}

	enabler := zap.NewAtomicLevelAt(zap.InfoLevel)
	core, logs := observer.New(enabler)

	logger, err := configureLogger(zap.New(core), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	foo(ContextWithLogger(context.Background(), logger))

	if logs.Len() != 2 {
		t.Errorf("unexpected entry count: %d", logs.Len())
		t.FailNow()
	}

	entries := logs.TakeAll()

	for i, expected := range []struct {
		fields  []zapcore.Field
		message string
	}{
		{
			fields:  []zapcore.Field{zap.String("service", "test")},
			message: "foo called",
		},
		{
			fields: []zapcore.Field{
				zap.String("service", "test"),
				zap.String("baz", "quux"),
			},
			message: "bar called",
		},
	} {
		if !reflect.DeepEqual(expected.fields, entries[i].Context) {
			t.Errorf("fields don't match, expected: %v, got: %v", expected.fields, entries[i].Context)
		}

		if expected.message != entries[i].Message {
			t.Errorf("messages don't match, expected: %s, got: %s", expected.message, entries[i].Message)
		}
	}
}
