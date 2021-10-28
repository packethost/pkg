//go:build tools
// +build tools

package tools

import (
	_ "github.com/tinkerbell/lint-install"
	_ "golang.org/x/tools/cmd/goimports"
	_ "mvdan.cc/gofumpt"
)
