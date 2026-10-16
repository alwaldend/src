package renders_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

func TestReimuContinuationPreservesEvidence(t *testing.T) {
	const (
		continuationRoot = "projects/renders/openspec/changes/reimu-fumo-finish/provenance/continuation-20260921"
		originalRoot     = "projects/renders/goals/reimu-fumo-finish/"
		originalHead     = "c7601f0fc80e0a94b585910459639ffccbdbdbd4"
		stashOID         = "f7766dcb29654f85087bd7759d552f89975833de"
		untrackedOID     = "44f303b09581f13f41f098e5ddfc5421337d4509"
	)
	runfilesSet, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	workspace := os.Getenv("TEST_WORKSPACE")
	if workspace == "" {
		t.Fatal("TEST_WORKSPACE is unavailable")
	}
	locate := func(relative string) string {
		t.Helper()
		if !fs.ValidPath(relative) || strings.Contains(relative, "\\") {
			t.Fatalf("unsafe continuation path %q", relative)
		}
		filename, err := runfilesSet.Rlocation(path.Join(workspace, continuationRoot, relative))
		if err != nil {
			t.Fatal(err)
		}
		return filename
	}
	read := func(relative string) []byte {
		t.Helper()
		data, err := os.ReadFile(locate(relative))
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	type artifact struct {
		Path         string `json:"path"`
		OriginalPath string `json:"original_path"`
		SourceOID    string `json:"source_oid"`
		Bytes        int    `json:"bytes"`
		SHA256       string `json:"sha256"`
	}
	var manifest struct {
		Schema              string     `json:"schema"`
		OriginalHead        string     `json:"original_head"`
		StashOID            string     `json:"stash_oid"`
		UntrackedTreeCommit string     `json:"untracked_tree_commit"`
		ResourceVersion     int        `json:"resource_version"`
		ActiveAttempt       string     `json:"active_attempt"`
		State               string     `json:"state"`
		Files               []artifact `json:"files"`
		RetiredCodePatch    struct {
			Path    string `json:"path"`
			Bytes   int    `json:"bytes"`
			SHA256  string `json:"sha256"`
			BaseOID string `json:"base_oid"`
			HeadOID string `json:"head_oid"`
		} `json:"retired_code_patch"`
	}
	if err := json.Unmarshal(read("manifest.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Schema != "reimu-fumo-continuation/v1" || manifest.OriginalHead != originalHead ||
		manifest.StashOID != stashOID || manifest.UntrackedTreeCommit != untrackedOID {
		t.Fatal("continuation schema or source identities changed")
	}
	if manifest.ResourceVersion != 375 || manifest.ActiveAttempt != "independent-original-macro-115" ||
		manifest.State != "historical open/active; no accepted candidate" {
		t.Fatal("continuation changed the preserved work or acceptance state")
	}
	if len(manifest.Files) != 629 {
		t.Fatalf("preserved %d files, want 629", len(manifest.Files))
	}
	verify := func(relative string, size int, expectedDigest string) {
		t.Helper()
		data := read(relative)
		digest := sha256.Sum256(data)
		if len(data) != size || hex.EncodeToString(digest[:]) != expectedDigest {
			t.Errorf("snapshot differs from preserved evidence: %s", relative)
		}
	}
	seen := make(map[string]bool)
	sourceCounts := make(map[string]int)
	totalBytes := 0
	for _, file := range manifest.Files {
		if !strings.HasPrefix(file.Path, "source/") || seen[file.Path] ||
			file.OriginalPath != originalRoot+strings.TrimPrefix(file.Path, "source/") {
			t.Fatalf("unexpected or duplicate continuation entry %+v", file)
		}
		if file.SourceOID != stashOID && file.SourceOID != untrackedOID {
			t.Fatalf("unrecognized source identity for %s: %s", file.Path, file.SourceOID)
		}
		seen[file.Path] = true
		sourceCounts[file.SourceOID]++
		totalBytes += file.Bytes
		verify(file.Path, file.Bytes, file.SHA256)
	}
	if totalBytes != 3137401 || sourceCounts[stashOID] != 348 || sourceCounts[untrackedOID] != 281 {
		t.Fatalf("continuation inventory changed: %d bytes, source counts %v", totalBytes, sourceCounts)
	}
	for _, required := range []string{
		"source/goal.yaml",
		"source/criteria.yaml",
		"source/criteria-revisions/4.yaml",
		"source/attempts/independent-original-macro-115/attempt.yaml",
		"source/attempts/independent-original-macro-115/plan.md",
		"source/attempts/independent-original-macro-115/result.md",
	} {
		if !seen[required] {
			t.Errorf("required continuation history is absent: %s", required)
		}
	}
	patch := manifest.RetiredCodePatch
	if patch.Path != "retired-goal-fix.patch" || patch.HeadOID != originalHead ||
		patch.BaseOID != "9e9f6148bf0cecb55a84261a230bf5c0565bb7dd" {
		t.Fatal("retired code patch source identities changed")
	}
	verify(patch.Path, patch.Bytes, patch.SHA256)

	// Runfiles directories contain symlinked files; visit each leaf without
	// following it so an unlisted packaged file cannot escape the inventory.
	sourceDir := filepath.Dir(locate("source/goal.yaml"))
	if err := filepath.WalkDir(sourceDir, func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(sourceDir, filename)
		if err != nil {
			return err
		}
		if !seen["source/"+filepath.ToSlash(relative)] {
			t.Errorf("unlisted continuation evidence: %s", relative)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
