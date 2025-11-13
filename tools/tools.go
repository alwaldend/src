//go:build tools
// +build tools

package tool

import (
	_ "github.com/StackExchange/dnscontrol/v4"
	_ "github.com/bazelbuild/buildtools/buildifier"
	_ "mvdan.cc/sh/v3/cmd/shfmt"
	// _ "github.com/gohugoio/hugo"
)
