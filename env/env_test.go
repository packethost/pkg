// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func ExampleBool() {
	name := "some_bool_environment_variable_that_is_not_set"
	os.Unsetenv(name)

	fmt.Println(Bool(name))
	fmt.Println(Bool(name, true))
	fmt.Println(Bool(name, true, false))
	fmt.Println(Bool(name, false, true))

	os.Setenv(name, "true")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, false))
	os.Setenv(name, "false")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, true))
	os.Setenv(name, "t")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, false))
	os.Setenv(name, "f")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, true))
	os.Setenv(name, "1")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, false))
	os.Setenv(name, "0")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, true))
	os.Setenv(name, "random-value")
	fmt.Println(Bool(name))
	fmt.Println(Bool(name, true))
	// Output:
	// false
	// true
	// true
	// false
	// true
	// true
	// false
	// false
	// true
	// true
	// false
	// false
	// true
	// true
	// false
	// false
	// false
	// false
}

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

func ExampleDuration() {
	name := "some_duration_environment_variable_that_is_not_set"
	os.Unsetenv(name)
	fmt.Println(Duration(name))
	fmt.Println(Duration(name, 15*time.Minute))
	fmt.Println(Duration(name, 15*time.Minute, 2*time.Hour))
	fmt.Println(Duration(name, 0, 2*time.Hour))
	os.Setenv(name, "14m2s")
	fmt.Println(Duration(name))
	fmt.Println(Duration(name).Seconds())
	fmt.Println(Duration(name, 30*time.Minute))
	// Output:
	// 0s
	// 15m0s
	// 15m0s
	// 0s
	// 14m2s
	// 842
	// 14m2s
}
