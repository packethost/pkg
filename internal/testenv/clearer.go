// Package testenv provides functionality for clearing and restoring back the environment.
// This is mainly intended to be used in tests.
package testenv

import (
	"os"
	"strings"
)

// Restorer is the interface that wraps the Restore method.
//
// Restore will fully restore the environment variables to what it was
// before cleared with Clear.
type Restorer interface {
	Restore()
}

type restorer func()

func (r restorer) Restore() {
	r()

}

// Clear will clear the environment and return a Restorer, whose
// Restore function should be called to return the environment
// back to what it was before Clear was called.
//
// A nice way to call this would be `defer Clear().Restore()`
func Clear() Restorer {
	old := os.Environ()
	fn := func() {
		os.Clearenv()
		for _, e := range old {
			kv := strings.SplitN(e, "=", 2)
			os.Setenv(kv[0], kv[1])
		}
	}
	os.Clearenv()

	return restorer(fn)
}
