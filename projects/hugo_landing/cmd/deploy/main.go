// Command deploy publishes rendered project landing sites to GitHub Pages.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type sites map[string]string

func (s sites) String() string { return "project=rendered/site" }
func (s sites) Set(value string) error {
	project, path, ok := strings.Cut(value, "=")
	if !ok {
		return errors.New("site must be project=rendered/site")
	}
	if _, exists := s[project]; exists {
		return fmt.Errorf("duplicate project %q", project)
	}
	s[project] = path
	return nil
}

type options struct {
	scratch string
	owner   string
	git     string
	dryRun  bool
}

var (
	projectName = regexp.MustCompile(`^[a-z0-9]+(?:_[a-z0-9]+)*$`)
	ownerName   = regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`)
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) error {
	flags := flag.NewFlagSet("deploy", flag.ContinueOnError)
	var opts options
	var manifest string
	var selected string
	all := sites{}
	flags.Var(all, "site", "repeatable project=rendered/site (relative to the working directory)")
	flags.StringVar(&manifest, "manifest", "", "JSON object mapping projects to site directories relative to the working directory")
	flags.StringVar(&selected, "project", "", "deploy only this project from the supplied sites")
	flags.StringVar(&opts.scratch, "scratch", "", "absolute task-owned scratch directory")
	flags.StringVar(&opts.owner, "owner", "alwaldend", "GitHub repository owner")
	flags.StringVar(&opts.git, "git", "git", "Git executable")
	flags.BoolVar(&opts.dryRun, "dry-run", false, "stage changes without committing or pushing")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional arguments")
	}
	if manifest != "" {
		data, err := os.ReadFile(manifest)
		if err != nil {
			return err
		}
		var entries map[string]string
		if err := json.Unmarshal(data, &entries); err != nil {
			return fmt.Errorf("decode manifest: %w", err)
		}
		for project, path := range entries {
			if err := all.Set(project + "=" + path); err != nil {
				return err
			}
		}
	}
	if !filepath.IsAbs(opts.scratch) || !ownerName.MatchString(opts.owner) || len(all) == 0 {
		return errors.New("absolute --scratch, valid --owner and at least one site are required")
	}
	all, err := normalizeSites(all, selected)
	if err != nil {
		return err
	}
	projects := make([]string, 0, len(all))
	for project := range all {
		projects = append(projects, project)
	}
	sort.Strings(projects)
	for _, project := range projects {
		slug := strings.ReplaceAll(project, "_", "-")
		remote := "git@github.com:" + opts.owner + "/" + slug + "-landing.git"
		result, err := deploy(opts, project, all[project], remote)
		if err != nil {
			return fmt.Errorf("%s: %w", project, err)
		}
		fmt.Fprintf(output, "%s: %s https://%s.alwaldend.com/\n", project, result, slug)
	}
	return nil
}

func normalizeSites(all sites, selected string) (sites, error) {
	if selected != "" {
		path, exists := all[selected]
		if !exists {
			return nil, fmt.Errorf("project %q is absent from supplied sites", selected)
		}
		all = sites{selected: path}
	}
	resolved := sites{}
	for project, path := range all {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve site directory: %w", project, err)
		}
		// Bazel runfiles may link the root directory to a tree artifact.
		// Only resolve the root; validateSite still rejects interior symlinks.
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve site directory: %w", project, err)
		}
		if err := validateSite(project, path); err != nil {
			return nil, fmt.Errorf("%s: %w", project, err)
		}
		resolved[project] = path
	}
	return resolved, nil
}

func validateSite(project, site string) error {
	if !projectName.MatchString(project) || len(strings.ReplaceAll(project, "_", "-")) > 63 {
		return errors.New("invalid project name")
	}
	if !filepath.IsAbs(site) {
		return errors.New("site directory must be absolute")
	}
	info, err := os.Stat(filepath.Join(site, "index.html"))
	if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
		return errors.New("site requires a nonempty index.html")
	}
	return filepath.WalkDir(site, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.EqualFold(entry.Name(), ".git") || (!entry.IsDir() && !entry.Type().IsRegular()) {
			return fmt.Errorf("unsupported site entry %q", path)
		}
		return nil
	})
}

// gitCommand intentionally excludes subprocess output from errors: authentication
// helpers and remote responses can contain credentials. Only the operation and
// exit status are reported.
func gitCommand(opts options, dir string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, opts.git, append([]string{"-c", "core.hooksPath=/dev/null"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	result, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s failed: %w", args[0], err)
	}
	return result, nil
}

func deploy(opts options, project, site, remote string) (result string, err error) {
	if err := validateSite(project, site); err != nil {
		return "", err
	}
	if err := os.MkdirAll(opts.scratch, 0o700); err != nil {
		return "", err
	}
	work, err := os.MkdirTemp(opts.scratch, project+"-")
	if err != nil {
		return "", err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(work)) }()
	clone := filepath.Join(work, "repo")
	if _, err := gitCommand(opts, work, "clone", "--no-checkout", "--", remote, clone); err != nil {
		return "", err
	}
	refs, err := gitCommand(opts, clone, "for-each-ref", "--format=%(refname)", "refs/remotes/origin/pages")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(refs)) != "" {
		// A clone already creates pages locally when it is the remote default.
		// Reset only this fresh clone's branch to the fetched deployment head.
		if _, err := gitCommand(opts, clone, "checkout", "-B", "pages", "origin/pages"); err != nil {
			return "", err
		}
	} else {
		if _, err := gitCommand(opts, clone, "checkout", "--orphan", "pages"); err != nil {
			return "", err
		}
	}
	if _, err := gitCommand(opts, clone, "rm", "-r", "-f", "--ignore-unmatch", "--", "."); err != nil {
		return "", err
	}
	if err := copySite(site, clone); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(clone, "CNAME"), []byte(strings.ReplaceAll(project, "_", "-")+".alwaldend.com\n"), 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(clone, ".nojekyll"), nil, 0o644); err != nil {
		return "", err
	}
	if _, err := gitCommand(opts, clone, "add", "--all", "--force", "--", "."); err != nil {
		return "", err
	}
	status, err := gitCommand(opts, clone, "status", "--porcelain")
	if err != nil {
		return "", err
	}
	if len(status) == 0 {
		head, err := gitCommand(opts, clone, "rev-parse", "HEAD")
		if err != nil {
			return "", err
		}
		return "unchanged " + strings.TrimSpace(string(head)), nil
	}
	if opts.dryRun {
		return "would publish", nil
	}
	if _, err := gitCommand(opts, clone, "-c", "user.name=Project landing publisher", "-c", "user.email=landing-publisher@users.noreply.github.com", "-c", "commit.gpgsign=false", "commit", "-m", "Publish project landing page"); err != nil {
		return "", err
	}
	if _, err := gitCommand(opts, clone, "push", "origin", "HEAD:refs/heads/pages"); err != nil {
		return "", err
	}
	head, err := gitCommand(opts, clone, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return "published " + strings.TrimSpace(string(head)), nil
}

func copySite(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !entry.Type().IsRegular() || strings.EqualFold(entry.Name(), ".git") {
			return fmt.Errorf("unsupported site entry %q", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
