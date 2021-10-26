// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package rollbar

import (
	"os"

	"github.com/packethost/pkg/env"
	rollbar "github.com/rollbar/rollbar-go"
	rollbarerrors "github.com/rollbar/rollbar-go/errors"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Setup(l *zap.SugaredLogger, service string) func() {
	log = l

	token := os.Getenv("ROLLBAR_TOKEN")
	if token == "" {
		log.Panicw("required envvar(ROLLBAR_TOKEN) is unset", "envvar", "ROLLBAR_TOKEN")
	}
	rollbar.SetToken(token)

	pkgEnv := getEnvironment()
	rollbar.SetEnvironment(pkgEnv)

	v := getVersion()
	rollbar.SetCodeVersion(v)

	enable := true
	if os.Getenv("ROLLBAR_DISABLE") != "" {
		enable = false
	}
	rollbar.SetEnabled(enable)
	rollbar.SetStackTracer(rollbarerrors.StackTracer)

	return rollbar.Wait
}

func Notify(err error) {
	rollbar.Error(err)
}

func getEnvironment() string {
	pkgEnv := env.Get("ENV", env.Get("EQUINIX_ENV"))
	if pkgEnv == "" {
		// EQUINIX_ENV was not set! Checking for PACKET_ENV - Packet no longer exists, please switch
		pkgEnv = env.Get("ENV", env.Get("PACKET_ENV"))
	}

	if pkgEnv == "" {
		log.Panicw("required envvar(ENV) is unset", "envvar", "ENV")
	}
	return pkgEnv
}

func getVersion() string {
	version := env.Get("VERSION", env.Get("EQUINIX_VERSION"))
	if version == "" {
		// EQUINIX_VERSION was not set! Checking for PACKET_VERSION - Packet no longer exists, please switch
		version = env.Get("VERSION", env.Get("PACKET_VERSION"))
	}

	if version == "" {
		log.Panicw("required envvar(VERSION) is unset", "envvar", "VERSION")
	}
	return version
}
