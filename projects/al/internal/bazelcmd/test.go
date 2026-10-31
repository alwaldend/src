package bazelcmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// testOptions selects the Git comparison and Bazel workspace without changing either.
type testOptions struct {
	// base identifies the explicit Git comparison reference.
	base string
	// workspace identifies the Bazel workspace to query and test.
	workspace string
	// dryRun requests selection without executing tests.
	dryRun bool
}

// testSelection describes requested test targets, not proof that tests passed or executed.
type testSelection struct {
	// Schema identifies the public selection report format.
	Schema string `json:"schema"`
	// Workspace is the absolute Bazel workspace whose graph was queried.
	Workspace string `json:"workspace"`
	// BaseReference is the operator's explicit comparison reference.
	BaseReference string `json:"base_reference"`
	// BaseCommit records the resolved comparison commit.
	BaseCommit string `json:"base_commit"`
	// MergeBase is the common ancestor used for changed-path collection.
	MergeBase string `json:"merge_base"`
	// Scope contains the target patterns limiting candidate tests.
	Scope []string `json:"scope"`
	// ChangedPaths contains repository-relative committed and local changes.
	ChangedPaths []string `json:"changed_paths"`
	// Selection distinguishes affected, full_scope, and no_changes selection.
	Selection string `json:"selection"`
	// Reasons explains why the complete requested scope was selected.
	Reasons []string `json:"reasons"`
	// Tests contains the sorted labels requested from Bazel, subject to its test filters.
	Tests []string `json:"tests"`
	// DryRun records whether test execution was explicitly disabled.
	DryRun bool `json:"dry_run"`
}

// NewCommand exposes graph-based affected tests without a general Bazel passthrough.
func NewCommand() *cobra.Command {
	options := testOptions{}
	command := &cobra.Command{Use: "bazel", Short: "Select tests using Git changes and the Bazel graph"}
	test := &cobra.Command{
		Use:   "test --base REF [TARGET_PATTERN ...]",
		Short: "Run affected tests within a Bazel workspace",
		Long:  "Compare Git changes since the merge base with REF, including staged, unstaged, and untracked files. Query Bazel for affected tests within TARGET_PATTERN (default //...). Changes that cannot be mapped safely select the entire requested scope. Wildcard patterns exclude manual-tagged roots; explicit tests and suites retain their members. Bazel configuration and caching control execution.",
		RunE: func(cmd *cobra.Command, patterns []string) error {
			selection, err := selectTests(cmd.Context(), options, patterns, cmd.ErrOrStderr())
			if err != nil {
				return fmt.Errorf("select affected Bazel tests: %w", err)
			}
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(selection); err != nil {
				return fmt.Errorf("write Bazel test selection: %w", err)
			}
			if options.dryRun || len(selection.Tests) == 0 {
				return nil
			}
			child := exec.CommandContext(cmd.Context(), "bazel", append([]string{"test", "--"}, selection.Tests...)...)
			child.Dir = selection.Workspace
			child.Stdin = cmd.InOrStdin()
			child.Stdout, child.Stderr = cmd.ErrOrStderr(), cmd.ErrOrStderr()
			if err := child.Run(); err != nil {
				return fmt.Errorf("run selected Bazel tests in %q: %w", selection.Workspace, err)
			}
			return nil
		},
	}
	test.Flags().StringVar(&options.base, "base", "", "Git reference whose merge base with HEAD selects changed paths (required)")
	test.Flags().StringVar(&options.workspace, "workspace", "", "Bazel workspace directory (default: current enclosing workspace)")
	test.Flags().BoolVar(&options.dryRun, "dry-run", false, "Print the selection without running tests")
	_ = test.MarkFlagRequired("base")
	command.AddCommand(test)
	return command
}

// selectTests derives a bounded test selection from Git changes and Bazel's dependency graph.
func selectTests(ctx context.Context, options testOptions, patterns []string, stderr io.Writer) (testSelection, error) {
	selection := testSelection{
		Schema: "bazel/test_selection/v1", BaseReference: options.base, DryRun: options.dryRun,
		ChangedPaths: []string{}, Reasons: []string{}, Tests: []string{},
	}
	if options.base == "" {
		return selection, fmt.Errorf("--base requires a nonempty Git reference")
	}
	if len(patterns) == 0 {
		patterns = []string{"//..."}
	}
	selection.Scope = append([]string(nil), patterns...)
	scope, err := testScope(patterns)
	if err != nil {
		return selection, fmt.Errorf("read test scope: %w", err)
	}
	selection.Workspace, err = testWorkspace(options.workspace)
	if err != nil {
		return selection, fmt.Errorf("locate Bazel workspace: %w", err)
	}
	root, err := gitOutput(ctx, selection.Workspace, "rev-parse", "--show-toplevel")
	if err != nil {
		return selection, fmt.Errorf("locate Git repository: %w", err)
	}
	repository := strings.TrimSpace(string(root))
	resolved, err := gitOutput(ctx, repository, "rev-parse", "--verify", "--end-of-options", options.base+"^{commit}")
	if err != nil {
		return selection, fmt.Errorf("resolve comparison reference %q: %w", options.base, err)
	}
	selection.BaseCommit = strings.TrimSpace(string(resolved))
	ancestor, err := gitOutput(ctx, repository, "merge-base", "--all", selection.BaseCommit, "HEAD")
	if err != nil {
		return selection, fmt.Errorf("find comparison merge base: %w", err)
	}
	bases := strings.Fields(string(ancestor))
	if len(bases) != 1 {
		return selection, fmt.Errorf("comparison requires one merge base; found %d", len(bases))
	}
	selection.MergeBase = bases[0]
	selection.ChangedPaths, err = changedTestPaths(ctx, repository, selection.MergeBase)
	if err != nil {
		return selection, fmt.Errorf("collect changed paths: %w", err)
	}
	if len(selection.ChangedPaths) == 0 {
		selection.Selection = "no_changes"
		return selection, nil
	}
	selection.Selection = "affected"
	for _, path := range selection.ChangedPaths {
		absolute := filepath.Join(repository, filepath.FromSlash(path))
		relative, err := filepath.Rel(selection.Workspace, absolute)
		if err != nil {
			return selection, fmt.Errorf("resolve changed path %q: %w", path, err)
		}
		if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			selection.Reasons = append(selection.Reasons, "outside_workspace: "+path)
		} else if buildConfiguration(path) {
			selection.Reasons = append(selection.Reasons, "build_configuration: "+path)
		} else if _, err := os.Stat(absolute); os.IsNotExist(err) {
			selection.Reasons = append(selection.Reasons, "deleted_path: "+path)
		} else if err != nil {
			return selection, fmt.Errorf("inspect changed path %q: %w", path, err)
		}
	}
	expression := "tests(" + scope + ")"
	if len(selection.Reasons) == 0 {
		sources, err := bazelQuery(ctx, selection.Workspace, "kind(\"source file\", deps("+scope+"))", "label", stderr)
		if err != nil {
			return selection, fmt.Errorf("query source dependencies: %w", err)
		}
		labels := map[string]string{}
		for _, label := range strings.Split(strings.TrimSpace(string(sources)), "\n") {
			if !strings.HasPrefix(label, "//") {
				continue // External repository paths are not inferred from label names.
			}
			pkg, name, ok := strings.Cut(strings.TrimPrefix(label, "//"), ":")
			if ok {
				labels[filepath.Clean(filepath.Join(selection.Workspace, filepath.FromSlash(pkg), filepath.FromSlash(name)))] = label
			}
		}
		var affected []string
		for _, path := range selection.ChangedPaths {
			label, ok := labels[filepath.Join(repository, filepath.FromSlash(path))]
			if !ok {
				selection.Reasons = append(selection.Reasons, "unmapped_path: "+path)
				continue
			}
			if _, err := queryWord(label); err != nil {
				selection.Reasons = append(selection.Reasons, "unmapped_path: "+path)
				continue
			}
			affected = append(affected, label)
		}
		if len(selection.Reasons) == 0 {
			changed, err := patternSet(affected)
			if err != nil {
				return selection, fmt.Errorf("construct changed-source query: %w", err)
			}
			expression = "(" + expression + " intersect rdeps(" + scope + ", " + changed + "))"
		}
	}
	if len(selection.Reasons) > 0 {
		selection.Selection = "full_scope"
	}
	result, err := bazelQuery(ctx, selection.Workspace, expression, "label", stderr)
	if err != nil {
		return selection, fmt.Errorf("query selected tests: %w", err)
	}
	for _, label := range strings.Split(strings.TrimSpace(string(result)), "\n") {
		if label == "" {
			continue
		}
		if _, err := patternSet([]string{label}); err != nil {
			return selection, fmt.Errorf("Bazel returned an invalid test label: %w", err)
		}
		selection.Tests = append(selection.Tests, label)
	}
	sort.Strings(selection.Tests)
	selection.Tests = slices.Compact(selection.Tests)
	return selection, nil
}

// testWorkspace finds the closest Bazel workspace above the selected caller directory.
func testWorkspace(selected string) (string, error) {
	if selected == "" {
		selected = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if selected == "" {
		selected = "."
	}
	directory, err := filepath.Abs(selected)
	if err != nil {
		return "", fmt.Errorf("resolve workspace directory: %w", err)
	}
	directory, err = filepath.EvalSymlinks(directory)
	if err != nil {
		return "", fmt.Errorf("resolve workspace links: %w", err)
	}
	for {
		for _, name := range []string{"MODULE.bazel", "WORKSPACE.bazel", "WORKSPACE"} {
			if info, err := os.Stat(filepath.Join(directory, name)); err == nil && info.Mode().IsRegular() {
				return directory, nil
			} else if err != nil && !os.IsNotExist(err) {
				return "", fmt.Errorf("inspect workspace marker %q: %w", name, err)
			}
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("no Bazel workspace contains %q", selected)
		}
		directory = parent
	}
}

// gitOutput executes a read-only Git command with explicit directory and cancellation.
func gitOutput(ctx context.Context, repository string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", repository}, args...)...)
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(diagnostic.String()))
	}
	return output, nil
}

// changedTestPaths unions the working tree, index, and untracked paths without splitting names.
func changedTestPaths(ctx context.Context, repository, base string) ([]string, error) {
	paths := map[string]bool{}
	for _, args := range [][]string{
		{"diff", "--name-only", "--no-renames", "-z", base, "--"},
		{"diff", "--cached", "--name-only", "--no-renames", "-z", base, "--"},
		{"ls-files", "--others", "--exclude-standard", "-z"},
	} {
		output, err := gitOutput(ctx, repository, args...)
		if err != nil {
			return nil, fmt.Errorf("read changed Git paths: %w", err)
		}
		for _, path := range bytes.Split(output, []byte{0}) {
			if len(path) > 0 {
				paths[string(path)] = true
			}
		}
	}
	result := make([]string, 0, len(paths))
	for path := range paths {
		result = append(result, path)
	}
	sort.Strings(result)
	return result, nil
}

// buildConfiguration identifies inputs whose effects extend beyond ordinary source edges.
func buildConfiguration(path string) bool {
	name := filepath.Base(path)
	return slices.Contains([]string{"BUILD", "BUILD.bazel", "MODULE.bazel", "MODULE.bazel.lock", "WORKSPACE", "WORKSPACE.bazel", ".bazelignore", ".bazelversion", ".bazeliskrc"}, name) ||
		strings.HasSuffix(name, ".bzl") || strings.HasSuffix(name, ".bazelrc")
}

// queryWord quotes a literal Bazel query word without interpreting operator input as syntax.
func queryWord(value string) (string, error) {
	if value == "" || strings.ContainsAny(value, "\r\n\x00\\") {
		return "", fmt.Errorf("query word contains unsupported characters")
	}
	if !strings.Contains(value, "\"") {
		return "\"" + value + "\"", nil
	}
	if !strings.Contains(value, "'") {
		return "'" + value + "'", nil
	}
	return "", fmt.Errorf("query word contains both quote delimiters")
}

// patternSet quotes absolute patterns as one Bazel query set expression.
func patternSet(patterns []string) (string, error) {
	words := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		if !strings.HasPrefix(pattern, "//") && !(strings.HasPrefix(pattern, "@") && strings.Contains(pattern, "//")) {
			return "", fmt.Errorf("target pattern %q must be absolute (//... or @repository//...)", pattern)
		}
		word, err := queryWord(pattern)
		if err != nil {
			return "", fmt.Errorf("quote target pattern %q: %w", pattern, err)
		}
		words = append(words, word)
	}
	return "set(" + strings.Join(words, " ") + ")", nil
}

// testScope excludes manual wildcard roots before test-suite expansion while preserving named roots.
func testScope(patterns []string) (string, error) {
	scope, err := patternSet(patterns)
	if err != nil {
		return "", fmt.Errorf("quote test scope: %w", err)
	}
	// Bazel renders list attributes as [value, value]; these boundaries match the
	// complete manual tag rather than tags such as not-manual or manual-test.
	eligible := "(" + scope + ` except attr("tags", "[\[ ]manual[,\]]", ` + scope + "))"
	var named []string
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "/...") || strings.HasSuffix(pattern, ":all") ||
			strings.HasSuffix(pattern, ":all-targets") || strings.HasSuffix(pattern, ":*") {
			continue
		}
		named = append(named, pattern)
	}
	if len(named) == 0 {
		return eligible, nil
	}
	explicit, err := patternSet(named)
	if err != nil {
		return "", fmt.Errorf("quote explicitly named tests: %w", err)
	}
	return "(" + eligible + " union " + explicit + ")", nil
}

// bazelQuery captures a complete query result and rejects nonzero partial results.
func bazelQuery(ctx context.Context, workspace, expression, format string, stderr io.Writer) ([]byte, error) {
	command := exec.CommandContext(ctx, "bazel", "query", expression, "--output="+format)
	command.Dir = workspace
	command.Stderr = stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("run Bazel query in %q: %w", workspace, err)
	}
	return output, nil
}
