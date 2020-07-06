// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Bool parses given environment variables as a boolean, or returns the default if the environment variable is empty/unset.
// If the value is empty or unset it will return the first value of def or false if none is given.
// Evaluates true if the value case-insensitive matches 1|t|true|y|yes.
func Bool(name string, def ...bool) bool {
	if v := os.Getenv(name); v != "" {
		v = strings.ToLower(v)
		if v == "1" || v == "t" || v == "true" || v == "y" || v == "yes" {
			return true
		}
		return false
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

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

// URL parses given environment variable as a URL, or returns the default if the environment variable is empty/unset.
// URL will panic if it fails to parse the value.
func URL(name string, def ...string) *url.URL {
	v := ""
	if len(def) > 0 {
		v = def[0]
	}

	value := Get(name, v)
	u, err := url.Parse(value)
	if err != nil {
		err = errors.Wrap(err, "failed to parse URL from env var")
		panic(err)
	}
	return u
}

// Duration parses given environment variable as a time.Duration, or returns the default if the environment variable is empty/unset.
// Duration will panic if it fails to parse the value.
func Duration(name string, def ...time.Duration) time.Duration {
	var v time.Duration
	if len(def) > 0 {
		v = def[0]
	}

	value := Get(name, v.String())
	d, err := time.ParseDuration(value)
	if err != nil {
		err = errors.Wrap(err, "failed to parse duration from env var")
		panic(err)
	}
	return d
}
