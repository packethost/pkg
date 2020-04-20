// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package env

import "os"

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
