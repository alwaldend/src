// dump writes the repository's generated DNS declaration pages.
//
// It reads the checked-in declarations and either reports or rewrites the
// generated pages, mirroring the repository's other generated-projection
// writers.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/infra/dns/internal/dump"
)

func main() {
	if err := run(os.Args[1:], os.Getenv("BUILD_WORKSPACE_DIRECTORY")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, defaultWorkspace string) error {
	flags := flag.NewFlagSet("dump", flag.ContinueOnError)
	if defaultWorkspace == "" {
		defaultWorkspace = "."
	}
	workspace := flags.String("workspace", defaultWorkspace, "workspace root to read (default: BUILD_WORKSPACE_DIRECTORY or current directory)")
	write := flags.Bool("write", false, "rewrite the generated pages instead of checking them")
	check := flags.Bool("check", false, "fail when a generated page is out of date")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("positional arguments are not supported")
	}
	if *write && *check {
		return fmt.Errorf("--write and --check are mutually exclusive")
	}

	pages, err := dump.RenderPages(os.DirFS(*workspace))
	if err != nil {
		return err
	}
	var stale []string
	for _, page := range pages {
		path := filepath.Join(*workspace, filepath.FromSlash(page.Path))
		switch {
		case *write:
			if err := os.WriteFile(path, []byte(page.Body), 0o644); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "wrote %s\n", page.Path)
		case *check:
			existing, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if string(existing) != page.Body {
				stale = append(stale, page.Path)
			}
		default:
			fmt.Printf("== %s ==\n%s", page.Path, page.Body)
		}
	}
	if len(stale) > 0 {
		return fmt.Errorf("generated DNS pages are out of date: %v; run bazel run //infra/dns/cmd/dump -- --write", stale)
	}
	return nil
}
