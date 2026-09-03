package catalogs

import (
	"os"
	"path/filepath"
	"testing"

	agent_system "git.alwaldend.com/alwaldend/src/tools/agents/cmd/agent_system"
)

func TestSkillCoverageArtifactsAreCurrent(t *testing.T) {
	root := t.TempDir()
	args := []string{
		"aggregate",
		"--catalog", "capability.json",
		"--input", "skill-cases.json",
		"--output", filepath.Join(root, "skill-coverage.json"),
		"--markdown", filepath.Join(root, "skill-coverage.md"),
	}
	if err := agent_system.Run(args, os.Stdout); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"skill-coverage.json", "skill-coverage.md"} {
		generated, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		tracked, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if string(generated) != string(tracked) {
			t.Fatalf("stale: tools/agents/catalogs/%s", name)
		}
	}
}
