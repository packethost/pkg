// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"fmt"
	"os"
)

func ExampleGet() {
	name := "some_environment_variable_that_is_not_set"
	os.Unsetenv(name)
	fmt.Println(Get(name))
	fmt.Println(Get(name, "this is the default"))
	fmt.Println(Get(name, "this is the default", "this one is ignored"))
	fmt.Println(Get(name, "", "this one is ignored"))
	os.Setenv(name, "this is the value set")
	fmt.Println(Get(name))
	fmt.Println(Get(name, "this is the default"))
	// Output:
	//
	// this is the default
	// this is the default
	//
	// this is the value set
	// this is the value set
}
