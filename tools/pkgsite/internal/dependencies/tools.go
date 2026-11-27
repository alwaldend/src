package main

import (
	_ "github.com/google/licensecheck"
	_ "github.com/google/safehtml"
	_ "golang.org/x/mod/modfile"
	_ "golang.org/x/net/context"
	_ "golang.org/x/pkgsite/static"
	_ "golang.org/x/tools/imports"
	_ "rsc.io/markdown"
)

// main supplies the entry point for the dependency-only executable.
func main() {
}
