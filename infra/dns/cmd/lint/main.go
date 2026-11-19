// lint loads DNS files at runtime, checks ownership, and prints their domains.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"git.alwaldend.com/alwaldend/src/infra/dns/internal/lint"
)

func main() {
	if err := run(os.Args[1:], os.Getenv("BUILD_WORKSPACE_DIRECTORY"), os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, defaultWorkspace string, output io.Writer) error {
	flags := flag.NewFlagSet("lint", flag.ContinueOnError)
	if defaultWorkspace == "" {
		defaultWorkspace = "."
	}
	workspace := flags.String("workspace", defaultWorkspace, "workspace to scan (default: BUILD_WORKSPACE_DIRECTORY or current directory)")
	zone := flags.String("zone", "alwaldend.com", "DNS zone for relative names")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("positional arguments are not supported")
	}
	report, err := lint.Scan(os.DirFS(*workspace), *zone)
	if err != nil {
		return err
	}
	return lint.WriteMarkdown(output, report)
}
