// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"fmt"
	"os"
	"strconv"
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

func ExampleInt() {
	name := "some_int_environment_variable_that_is_not_set"
	os.Unsetenv(name)

	fmt.Println(Int(name))
	fmt.Println(Int(name, 42))
	fmt.Println(Int(name, 42, 21))
	fmt.Println(Int(name, 0, 48))

	os.Setenv(name, strconv.Itoa(9))
	fmt.Println(Int(name))
	fmt.Println(Int(name, 42))
	// Output:
	// 0
	// 42
	// 42
	// 0
	// 9
	// 9
}

func ExampleURL() {
	name := "some_url_environment_variable_that_is_not_set"
	os.Unsetenv(name)
	fmt.Println(URL(name))
	fmt.Println(URL(name, "https://packet.com"))
	fmt.Println(URL(name, "https://packet.com", "https://www.equinix.com"))
	fmt.Println(URL(name, "", "https://tinkerbell.org"))
	os.Setenv(name, "https://tinkerbell.org")
	fmt.Println(URL(name))
	fmt.Println(URL(name).Host)
	fmt.Println(URL(name, "https://www.equinix.com/"))
	// Output:
	//
	// https://packet.com
	// https://packet.com
	//
	// https://tinkerbell.org
	// tinkerbell.org
	// https://tinkerbell.org
}
