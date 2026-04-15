module github.com/packethost/pkg/log/logr/v2

go 1.16

require (
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.2
	github.com/jacobweinstock/rollzap v0.1.3
	// parent module imported based on
	// https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository
	github.com/packethost/pkg/log/logr v0.0.0-20210106215246-8e2e62dc8f0c
	github.com/pkg/errors v0.9.1
	github.com/rollbar/rollbar-go v1.2.0
	go.uber.org/zap v1.19.0
)
