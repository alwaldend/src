package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	goalcatalog "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func fixtureRegistry() string {
	return `{
		"schema": "agents.alwaldend.com/phase1-registry/v1alpha1",
		"authorities": [
			{
				"id": "repository.goals",
				"kind": "goals",
				"source": "projects/agents/goals"
			}
		]
	}`
}

func fixtureGoalDir() map[string]string {
	return map[string]string{
		"tools/agents/declarations/registry.json": fixtureRegistry(),
		"projects/agents/goals/repo-agent-system/goal.yaml": `apiVersion: goals.alwaldend.com/v1alpha1
kind: Goal
metadata:
  name: repo-agent-system
  resourceVersion: "8"
  generation: 1
  creationTimestamp: "2026-08-31T13:55:01.640303113Z"
  annotations:
    goals.alwaldend.com/local-owner-root: projects/agents
spec:
  title: Make the repository a coherent agent system
  scope: project
  retention:
    policy: durable
  relationships:
    dependsOnGoalRefs: []
    supersedesGoalRefs: []
status:
  lifecycleGeneration: 5
  outcome: achieved
  execution: paused
  acceptedAttemptID: system-design-001
  acceptedResultDigest: sha256:` + strings.Repeat("0", 64) + `
  criteriaRevision: 2
  observedAt: "2026-08-31T16:32:21.398457663Z"
`,
		"projects/agents/goals/repo-agent-system/criteria.yaml": `apiVersion: goals.alwaldend.com/v1alpha1
kind: GoalCriteria
metadata:
  name: repo-agent-system
  resourceVersion: "2"
  generation: 2
  creationTimestamp: "2026-08-31T13:55:01.640303113Z"
spec:
  goalRef:
    name: repo-agent-system
  revision: 2
  items:
  - criterionID: audit-current-state
    revision: 1
    required: true
    statement: A repository-wide evidence audit identifies the current system.
    evidenceMethod: Inspect the committed current-state analysis.
`,
		"projects/agents/goals/repo-agent-system/criteria-revisions/2.yaml": `apiVersion: goals.alwaldend.com/v1alpha1
kind: GoalCriteria
metadata:
  name: repo-agent-system
  resourceVersion: "2"
  generation: 2
  creationTimestamp: "2026-08-31T13:55:01.640303113Z"
spec:
  goalRef:
    name: repo-agent-system
  revision: 2
  items:
  - criterionID: audit-current-state
    revision: 1
    required: true
    statement: A repository-wide evidence audit identifies the current system.
    evidenceMethod: Inspect the committed current-state analysis.
`,
		"projects/agents/goals/repo-agent-system/attempts/system-design-001/attempt.yaml": `apiVersion: goals.alwaldend.com/v1alpha1
kind: GoalAttempt
metadata:
  name: system-design-001
  resourceVersion: "2"
  generation: 1
  creationTimestamp: "2026-08-31T14:42:06.530697381Z"
spec:
  goalRef:
    name: repo-agent-system
  goalGeneration: 1
  lifecycleGeneration: 4
  criteriaRevision: 2
  criteriaDigest: sha256:` + strings.Repeat("a", 64) + `
  goalStateDigest: sha256:` + strings.Repeat("1", 64) + `
  workType: change
status:
  state: closed
  closedAt: "2026-08-31T16:32:21.39804838Z"
  artifacts:
    planDigest: sha256:` + strings.Repeat("2", 64) + `
    resultDigest: sha256:` + strings.Repeat("3", 64) + `
    evidence:
    - path: evidence/design-review.md
      digest: sha256:` + strings.Repeat("4", 64) + `
  review:
    decision: accept
    criteria:
    - criterionID: audit-current-state
      criterionRevision: 1
      verdict: pass
      evidenceRefs:
      - evidence/design-review.md
      - result.md
  observedAt: "2026-08-31T16:32:21.39804838Z"
`,
		"projects/agents/goals/repo-agent-system/attempts/system-design-001/plan.md":                   "# Plan\n",
		"projects/agents/goals/repo-agent-system/attempts/system-design-001/result.md":                 "# Result\n",
		"projects/agents/goals/repo-agent-system/attempts/system-design-001/evidence/design-review.md": "# Review\n",
	}
}

func TestGoalCompileComplete(t *testing.T) {
	root := writeFiles(t, fixtureGoalDir())
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/goal.json",
		"--markdown", "out/goal.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/goal.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := goalcatalog.DecodeGoalStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != goalcatalog.CompletenessComplete {
		t.Fatalf("expected complete, got %s: %v",
			catalog.Completeness, catalog.Limitations)
	}
	if len(catalog.Goals) != 1 {
		t.Fatalf("unexpected goals: %#v", catalog.Goals)
	}
	record := catalog.Goals[0]
	if record.Availability != "available" ||
		record.Identity == nil ||
		record.Identity.Name != "repo-agent-system" ||
		record.Identity.OwnerRoot != "projects/agents" ||
		record.CoarseStatus == nil ||
		record.CoarseStatus.Outcome != "achieved" {
		t.Fatalf("unexpected record: %#v", record)
	}
	tracked, err := os.ReadFile(filepath.Join(root, "out/goal.md"))
	if err != nil || !strings.Contains(string(tracked), "# Goal catalog") {
		t.Fatalf("markdown render missing: %v", err)
	}
}

func TestGoalCompileIncompleteOnMissingManifest(t *testing.T) {
	files := fixtureGoalDir()
	delete(files, "projects/agents/goals/repo-agent-system/goal.yaml")
	root := writeFiles(t, files)
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/goal.json",
		"--markdown", "out/goal.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/goal.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := goalcatalog.DecodeGoalStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != goalcatalog.CompletenessPartial {
		t.Fatalf("expected partial, got %s", catalog.Completeness)
	}
	if len(catalog.Goals) != 1 ||
		catalog.Goals[0].Availability != "unavailable" ||
		catalog.Goals[0].Reason == "" {
		t.Fatalf("unexpected goals: %#v", catalog.Goals)
	}
}

func TestGoalCompileRecordsInvalidAsUnavailable(t *testing.T) {
	files := fixtureGoalDir()
	files["projects/agents/goals/repo-agent-system/goal.yaml"] =
		"apiVersion: goals.alwaldend.com/v1alpha1\n" +
			"kind: Goal\n" +
			"metadata:\n" +
			"  name: repo-agent-system\n" +
			"  resourceVersion: \"8\"\n" +
			"  generation: 1\n" +
			"  creationTimestamp: \"2026-08-31T13:55:01.640303113Z\"\n" +
			"  annotations:\n" +
			"    goals.alwaldend.com/local-owner-root: projects/agents\n" +
			"spec:\n" +
			"  title: Make the repository a coherent agent system\n" +
			"  scope: project\n" +
			"  retention:\n" +
			"    policy: durable\n" +
			"  relationships:\n" +
			"    dependsOnGoalRefs: []\n" +
			"    supersedesGoalRefs: []\n" +
			"status:\n" +
			"  lifecycleGeneration: 5\n" +
			"  outcome: bogus\n" +
			"  execution: paused\n" +
			"  observedAt: \"2026-08-31T16:32:21.398457663Z\"\n" +
			"  criteriaRevision: 2\n"
	root := writeFiles(t, files)
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/goal.json",
		"--markdown", "out/goal.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/goal.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := goalcatalog.DecodeGoalStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != goalcatalog.CompletenessPartial {
		t.Fatalf("expected partial, got %s", catalog.Completeness)
	}
	if len(catalog.Goals) != 1 ||
		catalog.Goals[0].Availability != "unavailable" ||
		!strings.Contains(catalog.Goals[0].Reason, "invalid record") {
		t.Fatalf("unexpected goals: %#v", catalog.Goals)
	}
}
