package release

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestGitInfoLinkedCheckoutCLI exercises metadata generation from a linked
// checkout, including objects and remote configuration in its common Git directory.
func TestGitInfoLinkedCheckoutCLI(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	runGit := func(args ...string) string {
		t.Helper()
		command := exec.CommandContext(ctx, gitPath, args...)
		command.Env = []string{"PATH=" + os.Getenv("PATH"), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("Git fixture: %v\n%s", err, output)
		}
		return strings.TrimSpace(string(output))
	}
	directory := t.TempDir()
	repository := filepath.Join(directory, "repository")
	runGit("init", "--quiet", "--initial-branch=master", repository)
	runGit("-C", repository, "config", "user.name", "Release Fixture")
	runGit("-C", repository, "config", "user.email", "fixture@example.invalid")
	runGit("-C", repository, "config", "commit.gpgsign", "false")
	writeBundleFixture(t, filepath.Join(repository, "release.txt"), "Linked checkout payload.\n")
	runGit("-C", repository, "add", "release.txt")
	runGit("-C", repository, "commit", "--quiet", "-m", "Release metadata fixture")
	runGit("-C", repository, "remote", "add", "origin", "https://example.invalid/repository.git")
	head := runGit("-C", repository, "rev-parse", "HEAD")
	checkout := filepath.Join(directory, "checkout")
	runGit("-C", repository, "worktree", "add", "--quiet", "--detach", checkout, head)
	output := filepath.Join(directory, "git-info.json")
	if data, err := releaseBundleCommand("git-info", "--git_root", checkout, "--output", output).CombinedOutput(); err != nil {
		t.Fatalf("generate linked metadata: %v\n%s", err, data)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	var info struct {
		// Commits records history read through the common Git object directory.
		Commits map[string]json.RawMessage `json:"commits"`
		// Remotes records configuration read through the common Git directory.
		Remotes []struct {
			// Name identifies the fixture remote.
			Name string `json:"name"`
		} `json:"remotes"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
		t.Fatal(err)
	}
	if len(info.Commits) != 1 || info.Commits[head] == nil || len(info.Remotes) != 1 || info.Remotes[0].Name != "origin" {
		t.Fatalf("incomplete linked checkout metadata: %s", data)
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "linked-git-info.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
