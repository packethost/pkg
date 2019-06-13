package main

import (
	"github.com/packethost/pkg/log"
	"github.com/pkg/errors"
)

// helpfulWrapper is just here to demonstrate use of AddCallerSkip
// we don't really care to see helpfulWrapper as the caller when logging, we want to know the code that called helpfulWrapper instead
func helpfulWrapper(l log.Logger, message string) {
	l.AddCallerSkip(1).Info(message)
}

func main() {
	l, err := log.Init("github.com/packethost/pkg")
	if err != nil {
		panic(err)
	}

	ll := l.Package("log")
	ll.With("debug", true).Debug("hello this is a debug message")
	ll.With("debug", false).Info("hello this is a Info message")

	err = errors.New("the transducer has overheated")
	ll.With("error", err).Info("just an info level message about an error")
	ll.Error(err, "this is an actual error! Will even go to rollbar where we can ignore it or not")

	helpfulWrapper(ll, "this is being called via helpfulWrapper")
	l.Close()
}
