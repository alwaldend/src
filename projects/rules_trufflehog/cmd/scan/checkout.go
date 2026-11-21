package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// checkoutGit runs Git at the checked worktree with an explicit argument vector.
func checkoutGit(ctx context.Context, opts options, diagnostics io.Writer, arguments ...string) ([]byte, error) {
	git := opts.git
	if git == "" {
		git = "git"
	}
	cmd := exec.CommandContext(ctx, git, arguments...)
	cmd.Dir = opts.repository
	cmd.Stderr = diagnostics
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run git %s: %w", arguments[0], err)
	}
	return output, nil
}

// prepareCheckout verifies the checkout and writes the scanner exclusion list.
func prepareCheckout(ctx context.Context, opts options, exclusions string, diagnostics io.Writer) (string, error) {
	runGit := func(arguments ...string) ([]byte, error) {
		output, err := checkoutGit(ctx, opts, diagnostics, arguments...)
		if err != nil {
			return nil, fmt.Errorf("inspect checkout: %w", err)
		}
		return output, nil
	}
	head, err := runGit("rev-parse", "--verify", "HEAD")
	if err != nil {
		return "", fmt.Errorf("identify checkout revision: %w", err)
	}
	commit := strings.TrimSpace(string(head))
	if !commitPattern.MatchString(commit) {
		return "", errors.New("invalid checkout commit ID")
	}
	if _, err := runGit("diff", "--no-ext-diff", "--quiet", commit, "--"); err != nil {
		return "", fmt.Errorf("checkout scan requires tracked files to match %s: %w", commit, err)
	}
	lfs, err := runGit("lfs", "ls-files", "--json")
	if err != nil {
		return "", fmt.Errorf("list Git LFS files: %w", err)
	}
	var manifest struct {
		// Files is the Git LFS inventory for the checkout.
		Files []struct {
			// Name is the tracked path relative to the worktree.
			Name string `json:"name"`
			// Checkout reports whether the worktree contains the payload.
			Checkout bool `json:"checkout"`
		} `json:"files"`
	}
	if err := json.Unmarshal(lfs, &manifest); err != nil {
		return "", fmt.Errorf("decode Git LFS inventory: %w", err)
	}
	for _, file := range manifest.Files {
		if !file.Checkout {
			return "", fmt.Errorf("Git LFS payload %q is not materialized", file.Name)
		}
	}
	untracked, err := runGit("ls-files", "--others", "--directory", "-z")
	if err != nil {
		return "", fmt.Errorf("list untracked paths: %w", err)
	}
	patterns := []string{`^\.git(/|$)`}
	for _, name := range strings.Split(string(untracked), "\x00") {
		if name == "" {
			continue
		}
		path := regexp.QuoteMeta(strings.TrimSuffix(name, "/"))
		path = strings.NewReplacer("\n", `\n`, "\r", `\r`).Replace(path)
		patterns = append(patterns, "^"+path+"(/|$)")
	}
	if err := os.WriteFile(exclusions, []byte(strings.Join(patterns, "\n")+"\n"), 0o600); err != nil {
		return "", fmt.Errorf("write checkout scan exclusions: %w", err)
	}
	return commit, nil
}
