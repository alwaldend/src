//go:build tools
// +build tools

package tools

import (
	_ "github.com/StackExchange/dnscontrol/v4"
	_ "github.com/bazelbuild/buildtools/buildifier"
	_ "mvdan.cc/sh/v3/cmd/shfmt"
)
