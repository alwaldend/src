package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandReportsToolVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute(
		context.Background(),
		[]string{"--version"},
		func(string) string { return "" },
		&stdout,
		&stderr,
	); err != nil {
		t.Fatalf("Execute(--version): %v; stderr=%s", err, stderr.String())
	}
	if got, want := stdout.String(), "goal version 0.0.1\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestCommandInitAttachAndSessionShow(t *testing.T) {
	root := t.TempDir()
	runtimeRoot := t.TempDir()
	if err := os.Chmod(runtimeRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	getenv := func(name string) string {
		switch name {
		case "BUILD_WORKSPACE_DIRECTORY":
			return root
		case "XDG_RUNTIME_DIR":
			return runtimeRoot
		}
		return ""
	}
	run := func(args ...string) map[string]any {
		t.Helper()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		if err := Execute(context.Background(), args, getenv, &stdout, &stderr); err != nil {
			t.Fatalf("Execute(%v): %v; stderr=%s", args, err, stderr.String())
		}
		var result map[string]any
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			t.Fatalf("decode output %q: %v", stdout.String(), err)
		}
		return result
	}
	initialized := run(
		"init",
		"--goals-root", "out/cli-task/goals",
		"--goal-id", "cli-goal",
		"--title", "CLI goal",
		"--criterion", "The CLI works.",
	)
	if initialized["goalID"] != "cli-goal" {
		t.Fatalf("unexpected init output: %+v", initialized)
	}
	goalDir := filepath.Join(root, "out", "cli-task", "goals", "cli-goal")
	attached := run(
		"attach",
		"--session-root", "out/cli-task/goal-sessions",
		"--session-id", "cli-session",
		"--goal-dir", goalDir,
	)
	if attached["kind"] != "GoalSessionBinding" {
		t.Fatalf("unexpected attach output: %+v", attached)
	}
	shown := run(
		"show",
		"--session-root", "out/cli-task/goal-sessions",
		"--session-id", "cli-session",
	)
	if shown["kind"] != "GoalView" {
		t.Fatalf("unexpected show output: %+v", shown)
	}
}

func TestCommandRequiresTaskSpecificSessionRoot(t *testing.T) {
	root := t.TempDir()
	runtimeRoot := t.TempDir()
	if err := os.Chmod(runtimeRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	getenv := func(name string) string {
		switch name {
		case "BUILD_WORKSPACE_DIRECTORY":
			return root
		case "XDG_RUNTIME_DIR":
			return runtimeRoot
		}
		return ""
	}
	var output bytes.Buffer
	err := Execute(
		context.Background(),
		[]string{"show", "--session-id", "missing"},
		getenv,
		&output,
		&output,
	)
	if err == nil {
		t.Fatal("expected missing --session-root to fail")
	}
}

func TestCommandGraphAndSetRelationships(t *testing.T) {
	root := t.TempDir()
	runtimeRoot := t.TempDir()
	if err := os.Chmod(runtimeRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	getenv := func(name string) string {
		switch name {
		case "BUILD_WORKSPACE_DIRECTORY":
			return root
		case "XDG_RUNTIME_DIR":
			return runtimeRoot
		}
		return ""
	}
	run := func(args ...string) map[string]any {
		t.Helper()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		if err := Execute(
			context.Background(),
			args,
			getenv,
			&stdout,
			&stderr,
		); err != nil {
			t.Fatalf("Execute(%v): %v; stderr=%s", args, err, stderr.String())
		}
		var result map[string]any
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			t.Fatalf("decode output %q: %v", stdout.String(), err)
		}
		return result
	}
	goalsRoot := "out/cli-graph/goals"
	for _, name := range []string{"gamma", "alpha", "beta"} {
		run(
			"init",
			"--goals-root", goalsRoot,
			"--goal-id", name,
			"--title", name,
		)
	}
	alphaDir := filepath.Join(root, "out", "cli-graph", "goals", "alpha")
	updated := run(
		"set-relationships",
		"--goal-dir", alphaDir,
		"--expected-resource-version", "1",
		"--parent-goal", "not-created-yet",
		"--depends-on", "gamma",
		"--depends-on", "beta",
	)
	if updated["resourceVersion"] != "2" {
		t.Fatalf("unexpected relationship output: %#v", updated)
	}

	analysis := run("graph", "--goals-root", goalsRoot)
	if analysis["state"] != "Unknown" {
		t.Fatalf("graph state = %#v, want Unknown", analysis["state"])
	}
	assertCommandGraphNodeState(t, analysis, "alpha", "Waiting")
	if got := commandGraphEdgeCount(t, analysis, "Dependency"); got != 2 {
		t.Fatalf("dependency edge count = %d, want 2", got)
	}

	updated = run(
		"set-relationships",
		"--goal-dir", alphaDir,
		"--expected-resource-version", "2",
		"--clear-parent",
		"--depends-on", "beta,gamma",
	)
	if updated["resourceVersion"] != "3" {
		t.Fatalf("unexpected clear-parent output: %#v", updated)
	}
	analysis = run("graph", "--goals-root", goalsRoot)
	if analysis["state"] != "Valid" {
		t.Fatalf("graph state = %#v, want Valid", analysis["state"])
	}

	var output bytes.Buffer
	err := Execute(
		context.Background(),
		[]string{
			"set-relationships",
			"--goal-dir", alphaDir,
			"--expected-resource-version", "3",
			"--parent-goal", "beta",
			"--clear-parent",
		},
		getenv,
		&output,
		&output,
	)
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("parent flag conflict error = %v", err)
	}
}

func assertCommandGraphNodeState(
	t *testing.T,
	analysis map[string]any,
	name string,
	want string,
) {
	t.Helper()
	nodes, ok := analysis["nodes"].([]any)
	if !ok {
		t.Fatalf("graph nodes have unexpected shape: %#v", analysis["nodes"])
	}
	for _, value := range nodes {
		node, ok := value.(map[string]any)
		if !ok {
			continue
		}
		reference, ok := node["goalRef"].(map[string]any)
		if ok && reference["name"] == name {
			if node["dependencyState"] != want {
				t.Fatalf(
					"%s dependencyState = %#v, want %s",
					name,
					node["dependencyState"],
					want,
				)
			}
			return
		}
	}
	t.Fatalf("graph output has no node %q: %#v", name, analysis)
}

func commandGraphEdgeCount(
	t *testing.T,
	analysis map[string]any,
	relationshipName string,
) int {
	t.Helper()
	relationships, ok := analysis["relationships"].([]any)
	if !ok {
		t.Fatalf(
			"graph relationships have unexpected shape: %#v",
			analysis["relationships"],
		)
	}
	for _, value := range relationships {
		relationship, ok := value.(map[string]any)
		if !ok || relationship["relationship"] != relationshipName {
			continue
		}
		edges, ok := relationship["edges"].([]any)
		if !ok {
			t.Fatalf("graph edges have unexpected shape: %#v", relationship["edges"])
		}
		return len(edges)
	}
	t.Fatalf("graph output has no %q relationship", relationshipName)
	return 0
}

func TestCommandMigrateUsesDistinctSourceAndDestination(t *testing.T) {
	root := t.TempDir()
	runtimeRoot := t.TempDir()
	if err := os.Chmod(runtimeRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	getenv := func(name string) string {
		switch name {
		case "BUILD_WORKSPACE_DIRECTORY":
			return root
		case "XDG_RUNTIME_DIR":
			return runtimeRoot
		}
		return ""
	}
	source := filepath.Join(root, "out", "cli-migrate", "legacy", "legacy-goal")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := []byte("# Legacy goal\n\nOriginal prose.\n")
	if err := os.WriteFile(filepath.Join(source, "README.md"), legacy, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute(
		context.Background(),
		[]string{
			"migrate",
			"--source-goal-dir", source,
			"--destination-goals-root", "out/cli-migrate/imported/goals",
		},
		getenv,
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf("migrate command failed: %v; stderr=%s", err, stderr.String())
	}
	var result map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["goalID"] != "legacy-goal" || result["goalRef"] != "legacy-goal" {
		t.Fatalf("unexpected migration result: %#v", result)
	}
	unchanged, err := os.ReadFile(filepath.Join(source, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(unchanged, legacy) ||
		pathExistsForCommandTest(filepath.Join(source, "goal.yaml")) {
		t.Fatal("migrate command changed the source directory")
	}
	target := filepath.Join(
		root,
		"out",
		"cli-migrate",
		"imported",
		"goals",
		"legacy-goal",
		"goal.yaml",
	)
	if !pathExistsForCommandTest(target) {
		t.Fatal("migrate command did not publish the destination record")
	}

	stdout.Reset()
	stderr.Reset()
	err = Execute(
		context.Background(),
		[]string{
			"migrate",
			"--goal-dir", source,
			"--destination-goals-root", "out/cli-migrate/other/goals",
		},
		getenv,
		&stdout,
		&stderr,
	)
	if err == nil || !strings.Contains(err.Error(), "unknown flag") {
		t.Fatalf("deprecated in-place flag error = %v", err)
	}
	for _, args := range [][]string{
		{"migrate", "--source-goal-dir", source},
		{"migrate", "--destination-goals-root", "out/cli-migrate/missing/goals"},
	} {
		stdout.Reset()
		stderr.Reset()
		err = Execute(
			context.Background(),
			args,
			getenv,
			&stdout,
			&stderr,
		)
		if err == nil || !strings.Contains(
			err.Error(),
			"--source-goal-dir and --destination-goals-root are required",
		) {
			t.Fatalf("missing migration path error for %v = %v", args, err)
		}
	}
}

func pathExistsForCommandTest(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
