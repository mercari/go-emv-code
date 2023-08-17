//go:build tools
// +build tools

package tools

// see https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "honnef.co/go/tools/cmd/staticcheck"
)
