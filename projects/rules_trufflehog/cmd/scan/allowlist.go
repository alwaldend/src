package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/source_metadatapb"
)

// commitPattern requires a complete SHA-1 or SHA-256 Git object ID, never a ref.
var commitPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

// location is the complete suppression identity, independent of credential values.
type location struct {
	// Detector identifies the scanner's detector name.
	Detector string `json:"detector"`
	// Path is a canonical repository-relative Git path.
	Path string `json:"path"`
	// Commit is a full immutable Git commit ID.
	Commit string `json:"commit"`
	// Line is the positive source line reported by the scanner.
	Line int64 `json:"line"`
}

// validate rejects coordinates that could broaden or ambiguously match a finding.
func (l location) validate() error {
	if strings.TrimSpace(l.Detector) == "" || l.Detector != strings.TrimSpace(l.Detector) {
		return errors.New("detector must be a nonempty exact name")
	}
	if l.Path == "" || l.Path == "." || path.IsAbs(l.Path) || path.Clean(l.Path) != l.Path ||
		strings.HasPrefix(l.Path, "../") || strings.ContainsAny(l.Path, "\\\x00") {
		return errors.New("path must be a canonical repository-relative Git path")
	}
	if !commitPattern.MatchString(l.Commit) {
		return errors.New("commit must be a full lowercase Git commit ID")
	}
	if l.Line <= 0 {
		return errors.New("line must be positive")
	}
	return nil
}

// entry records the approved source tuple and its review rationale.
type entry struct {
	// location contains the serialized exact coordinates.
	location
	// Reason explains why the finding is accepted.
	Reason string `json:"reason"`
}

// allowlistDocument is the strict versioned public policy format.
type allowlistDocument struct {
	// Schema selects the supported policy contract.
	Schema string `json:"schema"`
	// Findings lists approved historical findings.
	Findings []entry `json:"findings"`
}

// loadAllowlist validates the whole policy before any scanner process starts.
func loadAllowlist(filename string) (map[location]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read allowlist: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var document allowlistDocument
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode allowlist: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("allowlist must contain exactly one JSON document")
	}
	if document.Schema != "trufflehog_allowlist/v1" || document.Findings == nil {
		return nil, errors.New("allowlist requires schema trufflehog_allowlist/v1 and a findings array")
	}
	approved := make(map[location]string, len(document.Findings))
	for index, item := range document.Findings {
		if err := item.location.validate(); err != nil {
			return nil, fmt.Errorf("allowlist entry %d: %w", index+1, err)
		}
		if strings.TrimSpace(item.Reason) == "" {
			return nil, fmt.Errorf("allowlist entry %d requires a reason", index+1)
		}
		if _, exists := approved[item.location]; exists {
			return nil, fmt.Errorf("allowlist entry %d duplicates a finding location", index+1)
		}
		approved[item.location] = item.Reason
	}
	return approved, nil
}

// findingSource selects the Git metadata from the scanner's source envelope.
type findingSource struct {
	// Data contains source-specific metadata keyed by its source kind.
	Data map[string]json.RawMessage `json:"Data"`
}

// scannerFinding contains only the fields needed to validate a result's identity.
type scannerFinding struct {
	// DetectorName identifies the detector whose result is being filtered.
	DetectorName string `json:"DetectorName"`
	// Verified must be present even though verification is disabled by the wrapper.
	Verified *bool `json:"Verified"`
	// SourceMetadata contains the immutable Git source coordinates.
	SourceMetadata findingSource `json:"SourceMetadata"`
}

// findingLocation validates result metadata without retaining a second raw value.
func findingLocation(raw json.RawMessage) (location, error) {
	var result scannerFinding
	if err := json.Unmarshal(raw, &result); err != nil {
		return location{}, fmt.Errorf("decode finding envelope: %w", err)
	}
	var source source_metadatapb.Git
	if err := json.Unmarshal(result.SourceMetadata.Data["Git"], &source); err != nil {
		return location{}, fmt.Errorf("decode finding Git metadata: %w", err)
	}
	if result.Verified == nil {
		return location{}, errors.New("finding is missing verification status")
	}
	key := location{Detector: result.DetectorName, Path: source.File, Commit: source.Commit, Line: source.Line}
	if err := key.validate(); err != nil {
		return location{}, fmt.Errorf("invalid finding coordinates: %w", err)
	}
	return key, nil
}
