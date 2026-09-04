package renders_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"os"
	"path"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

const (
	reimuAssetRunfileRoot             = "projects/renders/assets/reimu_fumo"
	reimuGoalCriteriaRunfile          = "projects/renders/goals/reimu-fumo-finish/criteria.yaml"
	expectedReimuReviewContractSHA256 = "4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4"
)

type reviewArtifact struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
}

type reviewReference struct {
	reviewArtifact
	Authority     string `json:"authority"`
	SourceStatus  string `json:"source_status"`
	SourceLocator string `json:"source_locator"`
	LicenseStatus string `json:"license_status"`
}

type fixedView struct {
	Camera        string     `json:"camera"`
	Location      [3]float64 `json:"location_m"`
	RotationEuler [3]float64 `json:"rotation_euler_rad"`
}

type reviewContract struct {
	SchemaVersion int                        `json:"schema_version"`
	ContractID    string                     `json:"contract_id"`
	Landmarks     reviewArtifact             `json:"landmarks"`
	References    map[string]reviewReference `json:"references"`
	FixedViews    map[string]fixedView       `json:"fixed_views"`
	Camera        struct {
		Projection               string  `json:"projection"`
		OrthoScale               float64 `json:"ortho_scale_m"`
		Resolution               [3]int  `json:"resolution"`
		CharacterUpAxis          string  `json:"character_up_axis"`
		CharacterFrontCameraAxis string  `json:"character_front_camera_axis"`
	} `json:"camera"`
	Reviewers struct {
		RetainMinimum        int    `json:"retain_minimum"`
		StageAndFinalMinimum int    `json:"stage_and_final_minimum"`
		ImplementationBlind  bool   `json:"implementation_blind"`
		IdentityRequirement  string `json:"identity_requirement"`
	} `json:"reviewers"`
}

func TestReimuReviewContract(t *testing.T) {
	locate := newReimuRunfileLocator(t)
	contractPath := locate("review_contract.json")
	verifyExactFileSHA256(t, contractPath, expectedReimuReviewContractSHA256)

	criteria, err := os.ReadFile(newWorkspaceRunfileLocator(t)(reimuGoalCriteriaRunfile))
	if err != nil {
		t.Fatalf("read active Reimu criteria: %v", err)
	}
	contractMarker := []byte("review contract sha256:")
	expectedMarker := append(contractMarker, expectedReimuReviewContractSHA256...)
	if markers, bindings := bytes.Count(criteria, contractMarker), bytes.Count(criteria, expectedMarker); markers != 2 || bindings != 2 {
		t.Fatalf("active criteria have %d review-contract markers and %d exact bindings, want 2 and 2", markers, bindings)
	}

	var contract reviewContract
	decodeReimuJSON(t, contractPath, &contract)

	if contract.SchemaVersion != 1 || contract.ContractID != "reimu-fumo-review-v1" {
		t.Fatalf("unexpected review contract identity: v%d %q", contract.SchemaVersion, contract.ContractID)
	}
	verifyReviewArtifact(t, locate, contract.Landmarks)
	if contract.Landmarks.Path != "LANDMARKS.md" {
		t.Fatalf("landmark path = %q, want LANDMARKS.md", contract.Landmarks.Path)
	}

	expectedReferences := map[string]string{
		"canonical_front": "references/canonical_front_25cm.png",
		"canonical_turn":  "references/canonical_turn_180.gif",
		"clean_front":     "references/clean_front.png",
		"physical_front":  "references/physical_front.png",
		"physical_side":   "references/physical_side.png",
		"sofa":            "references/sofa.gif",
		"turn":            "references/turn.gif",
	}
	if len(contract.References) != len(expectedReferences) {
		t.Fatalf("reference count = %d, want %d", len(contract.References), len(expectedReferences))
	}
	for name, expectedPath := range expectedReferences {
		reference, found := contract.References[name]
		if !found {
			t.Errorf("missing reference %q", name)
			continue
		}
		if reference.Path != expectedPath {
			t.Errorf("reference %s path = %q, want %q", name, reference.Path, expectedPath)
		}
		if reference.Authority == "" || reference.SourceStatus == "" || reference.SourceLocator == "" ||
			reference.LicenseStatus == "" {
			t.Errorf("reference %s has incomplete authority/provenance fields", name)
		}
		verifyReviewArtifact(t, locate, reference.reviewArtifact)
	}
	for name := range contract.References {
		if _, found := expectedReferences[name]; !found {
			t.Errorf("unexpected reference %q", name)
		}
	}

	expectedViews := map[string]fixedView{
		"front": {
			Camera: "Review_front_Camera", Location: [3]float64{0, -0.8, 0.13},
			RotationEuler: [3]float64{1.570796, 0, 0},
		},
		"rear": {
			Camera: "Review_rear_Camera", Location: [3]float64{0, 0.8, 0.13},
			RotationEuler: [3]float64{-1.570796, 3.141593, 0},
		},
		"side": {
			Camera: "Review_side_Camera", Location: [3]float64{0.8, 0, 0.13},
			RotationEuler: [3]float64{1.570796, 0, 1.570796},
		},
		"three_quarter": {
			Camera: "Review_three_quarter_Camera", Location: [3]float64{0.52, -0.52, 0.135},
			RotationEuler: [3]float64{1.563997, 0, 0.785398},
		},
		"three_quarter_mirror": {
			Camera: "Review_three_quarter_mirror_Camera", Location: [3]float64{-0.52, -0.52, 0.135},
			RotationEuler: [3]float64{1.563997, 0, -0.785398},
		},
	}
	if len(contract.FixedViews) != len(expectedViews) {
		t.Fatalf("fixed view count = %d, want %d", len(contract.FixedViews), len(expectedViews))
	}
	for name, expected := range expectedViews {
		actual, found := contract.FixedViews[name]
		if !found {
			t.Errorf("missing fixed view %q", name)
			continue
		}
		if actual.Camera != expected.Camera || !closeTriplet(actual.Location, expected.Location) ||
			!closeTriplet(actual.RotationEuler, expected.RotationEuler) {
			t.Errorf("fixed view %s = %+v, want %+v", name, actual, expected)
		}
	}
	if contract.Camera.Projection != "ORTHO" || math.Abs(contract.Camera.OrthoScale-0.292) > 1e-9 ||
		contract.Camera.Resolution != [3]int{512, 512, 100} || contract.Camera.CharacterUpAxis != "+Z" ||
		contract.Camera.CharacterFrontCameraAxis != "-Y" {
		t.Fatalf("unexpected fixed camera contract: %+v", contract.Camera)
	}
	if contract.Reviewers.RetainMinimum != 1 || contract.Reviewers.StageAndFinalMinimum != 2 ||
		!contract.Reviewers.ImplementationBlind || contract.Reviewers.IdentityRequirement == "" {
		t.Fatalf("unexpected reviewer contract: %+v", contract.Reviewers)
	}
}

type failureArtifact struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
}

type failureEvidenceManifest struct {
	SchemaVersion  int    `json:"schema_version"`
	Status         string `json:"status"`
	AcceptanceNote string `json:"acceptance_note"`
	Privacy        struct {
		ForbiddenMetadataChunksPresent bool `json:"forbidden_metadata_chunks_present"`
		SensitiveStringMatches         int  `json:"sensitive_string_matches"`
	} `json:"privacy"`
	Attempts []struct {
		AttemptID        string `json:"attempt_id"`
		AttemptStatus    string `json:"attempt_status"`
		AcceptanceStatus string `json:"acceptance_status"`
		Candidate        struct {
			BytesPublished             bool   `json:"bytes_published"`
			ReproducibleFromThisPacket bool   `json:"reproducible_from_this_packet"`
			SHA256                     string `json:"sha256"`
		} `json:"candidate"`
		RepresentationFamily string                     `json:"representation_family"`
		Views                map[string]failureArtifact `json:"views"`
	} `json:"attempts"`
}

func TestReimuFailureEvidence(t *testing.T) {
	locate := newReimuRunfileLocator(t)
	var manifest failureEvidenceManifest
	decodeReimuJSON(t, locate("failure_evidence/manifest.json"), &manifest)
	if manifest.SchemaVersion != 1 || manifest.Status != "historical_failure_evidence_only" ||
		manifest.AcceptanceNote == "" || manifest.Privacy.ForbiddenMetadataChunksPresent ||
		manifest.Privacy.SensitiveStringMatches != 0 {
		t.Fatalf("unexpected failure-evidence envelope: %+v", manifest)
	}

	expectedViews := map[string]map[string]bool{
		"head-hair-rebuild-003": {"front": true, "side": true},
		"curved-cap-004":        {"front": true, "side": true},
		"integrated-hair-005": {
			"front": true, "rear": true, "side": true, "three_quarter": true,
			"three_quarter_mirror": true,
		},
		"shingled-hair-006":   {"front": true, "side": true, "three_quarter": true},
		"evaluated-lobes-007": {"front": true, "side": true, "three_quarter": true},
	}
	if len(manifest.Attempts) != len(expectedViews) {
		t.Fatalf("failure attempt count = %d, want %d", len(manifest.Attempts), len(expectedViews))
	}
	seen := map[string]bool{}
	for _, attempt := range manifest.Attempts {
		expected, found := expectedViews[attempt.AttemptID]
		if !found || seen[attempt.AttemptID] {
			t.Errorf("unexpected or duplicate failure attempt %q", attempt.AttemptID)
			continue
		}
		seen[attempt.AttemptID] = true
		if attempt.AttemptStatus != "rejected" || attempt.AcceptanceStatus != "not_accepted" ||
			attempt.Candidate.BytesPublished || attempt.Candidate.ReproducibleFromThisPacket ||
			len(attempt.Candidate.SHA256) != 64 || attempt.RepresentationFamily == "" {
			t.Errorf("attempt %s has unsafe acceptance/provenance state", attempt.AttemptID)
		}
		if len(attempt.Views) != len(expected) {
			t.Errorf("attempt %s view count = %d, want %d", attempt.AttemptID, len(attempt.Views), len(expected))
		}
		for view, artifact := range attempt.Views {
			if !expected[view] {
				t.Errorf("attempt %s has unexpected view %q", attempt.AttemptID, view)
			}
			artifact.Path = path.Join("failure_evidence", artifact.Path)
			fullPath := verifyReviewArtifact(t, locate, reviewArtifact(artifact))
			verifyPrivacyCleanPNG(t, fullPath)
		}
	}
}

func newReimuRunfileLocator(t *testing.T) func(string) string {
	t.Helper()
	locate := newWorkspaceRunfileLocator(t)
	return func(relative string) string {
		if relative == "" || path.IsAbs(relative) || path.Clean(relative) != relative || relative == ".." ||
			len(relative) >= 3 && relative[:3] == "../" {
			t.Fatalf("unsafe artifact path %q", relative)
		}
		return locate(path.Join(reimuAssetRunfileRoot, relative))
	}
}

func newWorkspaceRunfileLocator(t *testing.T) func(string) string {
	t.Helper()
	runfilesSet, err := runfiles.New()
	if err != nil {
		t.Fatalf("open runfiles: %v", err)
	}
	workspace := os.Getenv("TEST_WORKSPACE")
	if workspace == "" {
		t.Fatal("TEST_WORKSPACE is unavailable")
	}
	return func(relative string) string {
		if relative == "" || path.IsAbs(relative) || path.Clean(relative) != relative || relative == ".." ||
			len(relative) >= 3 && relative[:3] == "../" {
			t.Fatalf("unsafe artifact path %q", relative)
		}
		resolved, err := runfilesSet.Rlocation(path.Join(workspace, relative))
		if err != nil {
			t.Fatalf("locate %s: %v", relative, err)
		}
		return resolved
	}
}

func verifyExactFileSHA256(t *testing.T, filename, expected string) {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read %s: %v", filename, err)
	}
	digest := sha256.Sum256(data)
	if actual := hex.EncodeToString(digest[:]); actual != expected {
		t.Fatalf("%s = sha256:%s, want sha256:%s", filename, actual, expected)
	}
}

func decodeReimuJSON(t *testing.T, filename string, value any) {
	t.Helper()
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("open %s: %v", filename, err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(value); err != nil {
		t.Fatalf("decode %s: %v", filename, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("%s has trailing JSON: %v", filename, err)
	}
}

func verifyReviewArtifact(t *testing.T, locate func(string) string, artifact reviewArtifact) string {
	t.Helper()
	fullPath := locate(artifact.Path)
	file, err := os.Open(fullPath)
	if err != nil {
		t.Fatalf("open %s: %v", artifact.Path, err)
	}
	defer file.Close()
	digest := sha256.New()
	size, err := io.Copy(digest, file)
	if err != nil {
		t.Fatalf("hash %s: %v", artifact.Path, err)
	}
	actualSHA256 := hex.EncodeToString(digest.Sum(nil))
	if actualSHA256 != artifact.SHA256 || size != artifact.SizeBytes {
		t.Errorf("%s = sha256:%s size:%d, want sha256:%s size:%d", artifact.Path, actualSHA256, size, artifact.SHA256, artifact.SizeBytes)
	}
	return fullPath
}

func verifyPrivacyCleanPNG(t *testing.T, filename string) {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read PNG %s: %v", filename, err)
	}
	for _, forbidden := range [][]byte{[]byte("/var/"), []byte("/home/"), []byte("simeon"), []byte("t3code-"), []byte(".blend")} {
		if bytes.Contains(bytes.ToLower(data), bytes.ToLower(forbidden)) {
			t.Errorf("PNG %s contains forbidden metadata string %q", filename, forbidden)
		}
	}
	if len(data) < 8 || !bytes.Equal(data[:8], []byte("\x89PNG\r\n\x1a\n")) {
		t.Fatalf("%s is not a PNG", filename)
	}
	allowed := map[string]bool{
		"IDAT": true, "IEND": true, "IHDR": true, "PLTE": true, "cHRM": true,
		"gAMA": true, "pHYs": true, "sRGB": true, "tRNS": true,
	}
	for offset := 8; offset < len(data); {
		if offset+12 > len(data) {
			t.Fatalf("truncated PNG chunk header in %s", filename)
		}
		length := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		chunkType := string(data[offset+4 : offset+8])
		if length < 0 || offset+12+length > len(data) {
			t.Fatalf("invalid PNG chunk %s in %s", chunkType, filename)
		}
		if !allowed[chunkType] {
			t.Errorf("PNG %s contains non-whitelisted chunk %s", filename, chunkType)
		}
		offset += 12 + length
	}
}

func closeTriplet(left, right [3]float64) bool {
	for index := range left {
		if math.Abs(left[index]-right[index]) > 1e-6 {
			return false
		}
	}
	return true
}
