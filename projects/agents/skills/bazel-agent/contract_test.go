package bazel_agent

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var legacyBazelAgentInvocation = regexp.MustCompile(
	`(?m)^\s*bazel_agent (build|test|run|query)(\s|$)`,
)

func TestActiveSkillGuidanceUsesRunnerContract(t *testing.T) {
	root := filepath.Join(
		os.Getenv("TEST_SRCDIR"),
		os.Getenv("TEST_WORKSPACE"),
	)
	if os.Getenv("TEST_SRCDIR") == "" || os.Getenv("TEST_WORKSPACE") == "" {
		t.Fatal("runfiles root is unavailable")
	}
	skillsDirectory := filepath.Join(root, "projects/agents/skills")
	entries, err := os.ReadDir(skillsDirectory)
	if err != nil {
		t.Fatalf("os.ReadDir(%q) error = %v", skillsDirectory, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		path := filepath.Join(skillsDirectory, entry.Name(), "SKILL.md")
		contents, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("os.ReadFile(%q) error = %v", path, err)
		}
		if matches := legacyBazelAgentInvocation.FindAll(contents, -1); matches != nil {
			relativePath, err := filepath.Rel(root, path)
			if err != nil {
				t.Fatalf("filepath.Rel(%q) error = %v", root, err)
			}
			t.Errorf(
				"%s uses legacy bazel_agent invocation: use bazel_agent bazel <command>",
				relativePath,
			)
		}
	}
}
