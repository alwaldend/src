package renders_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

const a202PacketRunfileRoot = "projects/renders/assets/reimu_fumo/donors/a202"

type a202PacketArtifact struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
}

type a202PacketCandidate struct {
	a202PacketArtifact
	Derivation      string `json:"derivation"`
	SourceSHA256    string `json:"source_sha256"`
	SourceSizeBytes int64  `json:"source_size_bytes"`
}

type a202PacketManifest struct {
	SchemaVersion int                           `json:"schema_version"`
	Status        string                        `json:"status"`
	Candidate     a202PacketCandidate           `json:"candidate"`
	Renders       map[string]a202PacketArtifact `json:"renders"`
	Sanitization  struct {
		UniqueWeakReferencesCleared int `json:"unique_library_weak_references_cleared"`
		EmbeddedTextBodiesCleared   int `json:"embedded_text_bodies_cleared"`
		OtherContentChanges         int `json:"other_geometry_rig_material_or_inventory_changes"`
	} `json:"sanitization"`
	Audit struct {
		LinkedLibraries                     int  `json:"linked_libraries"`
		MissingRuntimeDependencies          int  `json:"missing_runtime_dependencies"`
		EmbeddedTextBodiesNonempty          int  `json:"embedded_text_bodies_nonempty"`
		WeakReferenceFilepathsNonempty      int  `json:"library_weak_reference_filepaths_nonempty"`
		UnreadableProperties                int  `json:"unreadable_properties"`
		PrivacyFindings                     int  `json:"privacy_findings"`
		PNGPathTimestampOrEXIFChunks        int  `json:"png_path_timestamp_or_exif_chunks"`
		CandidatePreservedBeforeAfterRender bool `json:"candidate_preserved_before_and_after_render"`
	} `json:"audit"`
	Acceptance struct {
		InheritedCriterionPasses *[]string `json:"inherited_criterion_passes"`
	} `json:"acceptance"`
}

func TestA202PacketIntegrity(t *testing.T) {
	runfilesSet, err := runfiles.New()
	if err != nil {
		t.Fatalf("open runfiles: %v", err)
	}
	workspace := os.Getenv("TEST_WORKSPACE")
	if workspace == "" {
		t.Fatal("TEST_WORKSPACE is unavailable")
	}
	locate := func(relative string) string {
		runfilePath := path.Join(workspace, a202PacketRunfileRoot, relative)
		resolved, err := runfilesSet.Rlocation(runfilePath)
		if err != nil {
			t.Fatalf("locate %s: %v", relative, err)
		}
		return resolved
	}

	manifestFile, err := os.Open(locate("manifest.json"))
	if err != nil {
		t.Fatalf("open manifest: %v", err)
	}
	defer manifestFile.Close()
	var manifest a202PacketManifest
	decoder := json.NewDecoder(manifestFile)
	if err := decoder.Decode(&manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("manifest has trailing JSON: %v", err)
	}

	if manifest.SchemaVersion != 1 {
		t.Fatalf("schema_version = %d, want 1", manifest.SchemaVersion)
	}
	if manifest.Status != "rejected_technical_donor" {
		t.Fatalf("status = %q, want rejected_technical_donor", manifest.Status)
	}
	if manifest.Acceptance.InheritedCriterionPasses == nil {
		t.Fatal("acceptance.inherited_criterion_passes is missing")
	}
	if passes := *manifest.Acceptance.InheritedCriterionPasses; len(passes) != 0 {
		t.Fatalf("rejected donor inherits criterion passes: %v", passes)
	}

	if manifest.Candidate.Path != "model.blend" {
		t.Fatalf("candidate path = %q, want model.blend", manifest.Candidate.Path)
	}
	if manifest.Candidate.Derivation != "privacy_sanitized_from_recovered_a202" ||
		manifest.Candidate.SourceSHA256 != "6a9f3757facba526550e78817dc85f1d23cf85bcdad360228e113bb60d5f3aa0" ||
		manifest.Candidate.SourceSizeBytes != 1623719 {
		t.Fatalf("unexpected source/derivation contract: %+v", manifest.Candidate)
	}
	if manifest.Sanitization.UniqueWeakReferencesCleared != 113 ||
		manifest.Sanitization.EmbeddedTextBodiesCleared != 1 ||
		manifest.Sanitization.OtherContentChanges != 0 {
		t.Fatalf("unexpected sanitization contract: %+v", manifest.Sanitization)
	}
	if manifest.Audit.LinkedLibraries != 0 || manifest.Audit.MissingRuntimeDependencies != 0 ||
		manifest.Audit.EmbeddedTextBodiesNonempty != 0 ||
		manifest.Audit.WeakReferenceFilepathsNonempty != 0 || manifest.Audit.UnreadableProperties != 0 ||
		manifest.Audit.PrivacyFindings != 0 || manifest.Audit.PNGPathTimestampOrEXIFChunks != 0 ||
		!manifest.Audit.CandidatePreservedBeforeAfterRender {
		t.Fatalf("donor audit contract is not clean: %+v", manifest.Audit)
	}
	verifyA202PacketArtifact(t, locate, "candidate", manifest.Candidate.a202PacketArtifact)

	expectedViews := map[string]string{
		"front":                "renders/front.png",
		"rear":                 "renders/rear.png",
		"side":                 "renders/side.png",
		"three_quarter":        "renders/three_quarter.png",
		"three_quarter_mirror": "renders/three_quarter_mirror.png",
	}
	if len(manifest.Renders) != len(expectedViews) {
		t.Errorf(
			"render view count = %d, want %d",
			len(manifest.Renders),
			len(expectedViews),
		)
	}
	for view := range manifest.Renders {
		if _, expected := expectedViews[view]; !expected {
			t.Errorf("unexpected render view %q", view)
		}
	}
	for view, expectedPath := range expectedViews {
		artifact, found := manifest.Renders[view]
		if !found {
			t.Errorf("missing render view %q", view)
			continue
		}
		if artifact.Path != expectedPath {
			t.Errorf(
				"render %s path = %q, want %q",
				view,
				artifact.Path,
				expectedPath,
			)
			continue
		}
		verifyA202PacketArtifact(t, locate, "render "+view, artifact)
	}
}

func verifyA202PacketArtifact(
	t *testing.T,
	locate func(string) string,
	name string,
	artifact a202PacketArtifact,
) {
	t.Helper()
	file, err := os.Open(locate(artifact.Path))
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}
	defer file.Close()
	digest := sha256.New()
	size, err := io.Copy(digest, file)
	if err != nil {
		t.Fatalf("hash %s: %v", name, err)
	}
	actualSHA256 := hex.EncodeToString(digest.Sum(nil))
	if actualSHA256 != artifact.SHA256 {
		t.Errorf("%s SHA-256 = %s, want %s", name, actualSHA256, artifact.SHA256)
	}
	if size != artifact.SizeBytes {
		t.Errorf("%s size = %d, want %d", name, size, artifact.SizeBytes)
	}
}
