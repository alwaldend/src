package renders_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

func TestReimuMigrationPreservesEvidence(t *testing.T) {
	const changeRoot = "projects/renders/openspec/changes/reimu-fumo-finish"
	runfilesSet, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	workspace := os.Getenv("TEST_WORKSPACE")
	if workspace == "" {
		t.Fatal("TEST_WORKSPACE is unavailable")
	}
	read := func(relative string) []byte {
		t.Helper()
		if path.IsAbs(relative) || path.Clean(relative) != relative || strings.HasPrefix(relative, "../") {
			t.Fatalf("unsafe migration path %q", relative)
		}
		filename, err := runfilesSet.Rlocation(path.Join(workspace, changeRoot, relative))
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	var manifest struct {
		SchemaVersion int `json:"schema_version"`
		Files         []struct {
			SourcePath   string `json:"source_path"`
			SnapshotPath string `json:"snapshot_path"`
			SHA256       string `json:"sha256"`
			SizeBytes    int    `json:"size_bytes"`
		} `json:"files"`
	}
	if err := json.Unmarshal(read("provenance/manifest.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != 1 || len(manifest.Files) == 0 {
		t.Fatal("missing versioned provenance inventory")
	}
	seen := make(map[string]bool)
	for _, file := range manifest.Files {
		if !strings.HasPrefix(file.SourcePath, "projects/renders/goals/reimu-fumo-finish/") ||
			!strings.HasPrefix(file.SnapshotPath, "provenance/source/") || seen[file.SourcePath] {
			t.Fatalf("unexpected or duplicate provenance entry %+v", file)
		}
		seen[file.SourcePath] = true
		data := read(file.SnapshotPath)
		digest := sha256.Sum256(data)
		if len(data) != file.SizeBytes || hex.EncodeToString(digest[:]) != file.SHA256 {
			t.Errorf("snapshot differs from imported evidence: %s", file.SourcePath)
		}
	}
	for _, required := range []string{"goal.yaml", "criteria.yaml", "criteria-revisions/4.yaml", "attempts/flatten-dose-response-018/result.md"} {
		if !seen["projects/renders/goals/reimu-fumo-finish/"+required] {
			t.Errorf("required history is absent: %s", required)
		}
	}

	var migration struct {
		LegacyStatus struct {
			Outcome              string  `json:"outcome"`
			Execution            string  `json:"execution"`
			AcceptedAttemptID    *string `json:"accepted_attempt_id"`
			AcceptedResultDigest *string `json:"accepted_result_digest"`
		} `json:"legacy_status"`
		AcceptanceMapping []struct {
			CriterionID string `json:"criterion_id"`
			Checked     bool   `json:"checked"`
		} `json:"acceptance_mapping"`
	}
	if err := json.Unmarshal(read("migration.json"), &migration); err != nil {
		t.Fatal(err)
	}
	status := migration.LegacyStatus
	if status.Outcome != "open" || status.Execution != "blocked" || status.AcceptedAttemptID != nil || status.AcceptedResultDigest != nil {
		t.Fatalf("migration changed the legacy acceptance state: %+v", status)
	}
	if len(migration.AcceptanceMapping) != 8 {
		t.Fatalf("migrated %d acceptance criteria, want 8", len(migration.AcceptanceMapping))
	}
	for _, criterion := range migration.AcceptanceMapping {
		if criterion.Checked {
			t.Errorf("migration falsely completed %s", criterion.CriterionID)
		}
	}
}
