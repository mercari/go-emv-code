// +build tools

package tools

// see https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/kisielk/errcheck"
	_ "github.com/reviewdog/reviewdog/cmd/reviewdog"
	_ "golang.org/x/lint"
	_ "gotest.tools/gotestsum"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
