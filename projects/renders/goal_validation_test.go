package renders_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestGoalRecordsValidate(t *testing.T) {
	goalBinary, found := bazel.FindBinary("projects/goal/cmd/goal", "goal")
	if !found {
		t.Fatal("goal validator binary is unavailable")
	}

	runfilesRoot := filepath.Join(
		os.Getenv("TEST_SRCDIR"),
		os.Getenv("TEST_WORKSPACE"),
	)
	if os.Getenv("TEST_SRCDIR") == "" || os.Getenv("TEST_WORKSPACE") == "" {
		t.Fatal("runfiles root is unavailable")
	}

	workspaceRoot := filepath.Join(t.TempDir(), "workspace")
	goalsRoot := filepath.Join(workspaceRoot, "projects", "renders", "goals")
	if err := copyTree(
		filepath.Join(runfilesRoot, "projects", "renders", "goals"),
		goalsRoot,
	); err != nil {
		t.Fatalf("copy goal records: %v", err)
	}

	runtimeRoot := filepath.Join(t.TempDir(), "runtime")
	if err := os.Mkdir(runtimeRoot, 0o700); err != nil {
		t.Fatalf("create runtime directory: %v", err)
	}

	command := exec.Command(
		goalBinary,
		"--workspace-root", workspaceRoot,
		"validate",
		"--goals-root", "projects/renders/goals",
	)
	command.Env = []string{"XDG_RUNTIME_DIR=" + runtimeRoot}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("validate goal records: %v\n%s", err, output)
	}

	var result struct {
		Kind  string `json:"kind"`
		Valid bool   `json:"valid"`
		Count int    `json:"count"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode validation result: %v\n%s", err, output)
	}
	if result.Kind != "GoalValidation" || !result.Valid || result.Count == 0 {
		t.Fatalf("unexpected validation result: %+v", result)
	}
}

func copyTree(sourceRoot string, destinationRoot string) error {
	return filepath.WalkDir(sourceRoot, func(
		path string,
		entry fs.DirEntry,
		walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		destination := filepath.Join(destinationRoot, relative)
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported goal record entry %s", relative)
		}
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}
		source, err := os.Open(path)
		if err != nil {
			return err
		}
		target, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			source.Close()
			return err
		}
		_, copyErr := io.Copy(target, source)
		sourceCloseErr := source.Close()
		closeErr := target.Close()
		if copyErr != nil {
			return copyErr
		}
		if sourceCloseErr != nil {
			return sourceCloseErr
		}
		return closeErr
	})
}
