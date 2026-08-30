package fsstore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckpointRejectsBinaryMarkdownInputs(t *testing.T) {
	for _, encoding := range []struct {
		name    string
		content []byte
		want    string
	}{
		{name: "invalid UTF-8", content: []byte{0xff}, want: "valid UTF-8"},
		{name: "NUL", content: []byte("# Text\n\x00"), want: "NUL"},
	} {
		for _, artifact := range []string{"plan", "result", "evidence"} {
			t.Run(encoding.name+" "+artifact, func(t *testing.T) {
				store, root := newTestStore(t)
				goalID := fmt.Sprintf(
					"binary-%s-%s",
					strings.ToLower(strings.ReplaceAll(encoding.name, " ", "-")),
					artifact,
				)
				goalDir := initTestGoal(t, store, root, goalID)
				input := filepath.Join(root, "out", "task", artifact+".md")
				if err := os.WriteFile(input, encoding.content, 0o644); err != nil {
					t.Fatal(err)
				}
				options := CheckpointOptions{
					GoalDir:                 goalDir,
					ExpectedResourceVersion: "1",
					AttemptID:               "attempt-1",
				}
				switch artifact {
				case "plan":
					options.PlanFile = input
				case "result":
					options.ResultFile = input
				case "evidence":
					options.EvidenceFiles = []string{input}
				}
				reference, err := store.Checkpoint(options)
				if err == nil || !strings.Contains(err.Error(), encoding.want) {
					t.Fatalf("Checkpoint() error = %v, want %q", err, encoding.want)
				}
				if reference != (GoalReference{}) {
					t.Fatalf("Checkpoint() reference = %+v, want no commit", reference)
				}
				goal, err := store.readGoalManifest(goalDir)
				if err != nil {
					t.Fatal(err)
				}
				if goal.Metadata.ResourceVersion != "1" {
					t.Fatalf("resourceVersion = %s, want 1", goal.Metadata.ResourceVersion)
				}
				if pathExists(filepath.Join(goalDir, "attempts", "attempt-1")) {
					t.Fatal("rejected Markdown input left a canonical attempt")
				}
			})
		}
	}
}

func TestValidateGoalRejectsBinaryCanonicalMarkdown(t *testing.T) {
	for _, test := range []struct {
		name     string
		relative string
		content  []byte
		want     string
	}{
		{name: "plan invalid UTF-8", relative: "plan.md", content: []byte{0xff}, want: "valid UTF-8"},
		{name: "result NUL", relative: "result.md", content: []byte("result\x00"), want: "NUL"},
		{
			name:     "evidence invalid UTF-8",
			relative: filepath.Join("evidence", "proof.md"),
			content:  []byte{0xff},
			want:     "valid UTF-8",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalDir := initTestGoal(t, store, root, "binary-canonical")
			evidence := filepath.Join(root, "out", "task", "proof.md")
			writeTestFile(t, evidence, "# Evidence\n\nVerified.\n")
			if _, err := store.Checkpoint(CheckpointOptions{
				GoalDir:                 goalDir,
				ExpectedResourceVersion: "1",
				AttemptID:               "attempt-1",
				EvidenceFiles:           []string{evidence},
			}); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(goalDir, "attempts", "attempt-1", test.relative)
			if err := os.WriteFile(target, test.content, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := store.ValidateGoal(goalDir); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidateGoal() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestMigrationRejectsBinaryLegacyMarkdown(t *testing.T) {
	for _, test := range []struct {
		name    string
		content []byte
		want    string
	}{
		{name: "invalid UTF-8", content: []byte{0xff}, want: "valid UTF-8"},
		{name: "NUL", content: []byte("evidence\x00"), want: "NUL"},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, root := newTestStore(t)
			source := filepath.Join(root, "legacy", "binary-import")
			writeLegacyGoal(t, source, "Binary import")
			if err := os.WriteFile(filepath.Join(source, "evidence.md"), test.content, 0o644); err != nil {
				t.Fatal(err)
			}
			destinationRoot := filepath.Join(root, "project", "goals")
			_, err := store.Migrate(MigrateOptions{
				SourceGoalDir:        source,
				DestinationGoalsRoot: destinationRoot,
			})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Migrate() error = %v, want %q", err, test.want)
			}
			if pathExists(filepath.Join(destinationRoot, "binary-import")) {
				t.Fatal("rejected legacy Markdown published a migrated goal")
			}
		})
	}
}
