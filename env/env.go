// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Get retrieves the value of the environment variable named by the key.
// If the value is empty or unset it will return the first value of def or "" if none is given
func Get(name string, def ...string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// Int parses given environment variable as an int, or returns the default if the environment variable is empty/unset.
// Int will panic if it fails to parse the value.
func Int(name string, def ...int) int {
	v := os.Getenv(name)
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			err = errors.Wrap(err, "failed to parse int from env var")
			panic(err)
		}
		return i
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
