// Package catalogv1alpha1 defines shared deterministic repository-agent
// catalog contracts.
//
// Catalogs are replaceable projections, never authorities. Each catalog is a
// bounded offline derivation from owner-local facts. The bytes are
// deterministic: they contain no generation timestamp, no absolute checkout
// path, and no observation time. The digest covers the canonical schema
// bytes with this digest field omitted.
package catalogv1alpha1

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

const (
	// APIVersion is the catalog API version.
	APIVersion = "agents.alwaldend.com/catalog/v1alpha1"
)

var (
	identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._/-][a-z0-9]+)*$`)
	digestPattern     = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)
	versionPattern    = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
)

// Completeness describes how completely the catalog covers its eligible
// universe.
type Completeness string

const (
	CompletenessComplete  Completeness = "complete"
	CompletenessPartial   Completeness = "partial"
	CompletenessTruncated Completeness = "truncated"
	CompletenessUnknown   Completeness = "unknown"
)

// CatalogInput is one owner-local source consumed by the compiler.
type CatalogInput struct {
	Path   string `json:"path"`
	Role   string `json:"role"`
	Digest string `json:"digest"`
}

// CatalogBounds records explicit cardinality and byte bounds.
type CatalogBounds struct {
	Eligible       int  `json:"eligible"`
	Emitted        int  `json:"emitted"`
	Unavailable    int  `json:"unavailable"`
	MaxItems       int  `json:"maxItems"`
	MaxInputBytes  int  `json:"maxInputBytes"`
	MaxOutputBytes int  `json:"maxOutputBytes"`
	Truncated      bool `json:"truncated"`
}

// CatalogConflict reports a deterministic conflict between owner-local facts.
type CatalogConflict struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	SourcePaths []string `json:"sourcePaths"`
}

// CatalogEnvelope is the shared deterministic catalog shape. Catalogs are
// complete structural documents; limitations are required whenever
// completeness is not complete.
type CatalogEnvelope struct {
	Schema            string            `json:"schema"`
	Kind              string            `json:"kind"`
	ID                string            `json:"id"`
	DerivationVersion string            `json:"derivationVersion"`
	ProducerRef       string            `json:"producerRef"`
	SourceRevision    string            `json:"sourceRevision"`
	Inputs            []CatalogInput    `json:"inputs"`
	Bounds            CatalogBounds     `json:"bounds"`
	Completeness      Completeness      `json:"completeness"`
	Limitations       []string          `json:"limitations"`
	Conflicts         []CatalogConflict `json:"conflicts"`
	Digest            string            `json:"digest"`
}

// Validate checks the envelope invariants. itemIDs carries the identities of
// the catalog-specific items for identity and sort validation.
func (envelope CatalogEnvelope) Validate(itemIDs []string) error {
	if envelope.Schema != APIVersion+"/"+envelope.Kind {
		return fmt.Errorf("schema %q does not match kind %q", envelope.Schema, envelope.Kind)
	}
	if !identifierPattern.MatchString(envelope.Kind) ||
		!identifierPattern.MatchString(envelope.ID) {
		return fmt.Errorf("malformed catalog identity kind=%q id=%q", envelope.Kind, envelope.ID)
	}
	if !versionPattern.MatchString(envelope.DerivationVersion) {
		return fmt.Errorf("malformed derivation version %q", envelope.DerivationVersion)
	}
	if envelope.SourceRevision == "" || len(envelope.SourceRevision) > 64 {
		return fmt.Errorf("malformed source revision %q", envelope.SourceRevision)
	}
	if envelope.Inputs == nil {
		return fmt.Errorf("catalog inputs must be a non-null array")
	}
	if err := validateInputs(envelope.Inputs); err != nil {
		return err
	}
	if envelope.Bounds.MaxItems <= 0 ||
		envelope.Bounds.MaxInputBytes <= 0 ||
		envelope.Bounds.MaxOutputBytes <= 0 {
		return fmt.Errorf("catalog bounds must be positive")
	}
	if envelope.Bounds.Eligible < 0 || envelope.Bounds.Emitted < 0 ||
		envelope.Bounds.Unavailable < 0 {
		return fmt.Errorf("catalog bound counts cannot be negative")
	}
	if envelope.Bounds.Emitted > envelope.Bounds.MaxItems ||
		envelope.Bounds.Unavailable > envelope.Bounds.MaxItems {
		return fmt.Errorf("emitted or unavailable items exceed maxItems")
	}
	if envelope.Bounds.Truncated && envelope.Completeness != CompletenessTruncated {
		return fmt.Errorf("truncated bounds require truncated completeness")
	}
	if !oneOfCompleteness(envelope.Completeness) {
		return fmt.Errorf("unknown completeness %q", envelope.Completeness)
	}
	if envelope.Completeness != CompletenessComplete &&
		len(envelope.Limitations) == 0 {
		return fmt.Errorf("non-complete catalog requires at least one limitation")
	}
	if envelope.Conflicts == nil {
		return fmt.Errorf("catalog conflicts must be a non-null array")
	}
	if err := validateConflicts(envelope.Conflicts); err != nil {
		return err
	}
	if envelope.Digest != "" && !digestPattern.MatchString(envelope.Digest) {
		return fmt.Errorf("malformed catalog digest %q", envelope.Digest)
	}
	seen := map[string]bool{}
	for _, id := range itemIDs {
		if !identifierPattern.MatchString(id) || strings.Contains(id, "..") {
			return fmt.Errorf("malformed item identity %q", id)
		}
		if seen[id] {
			return fmt.Errorf("duplicate item identity %q", id)
		}
		seen[id] = true
	}
	return nil
}

func oneOfCompleteness(value Completeness) bool {
	switch value {
	case CompletenessComplete, CompletenessPartial, CompletenessTruncated,
		CompletenessUnknown:
		return true
	}
	return false
}

func validateInputs(inputs []CatalogInput) error {
	seen := map[string]bool{}
	for _, input := range inputs {
		if input.Path == "" || absPath(input.Path) || escapes(input.Path) ||
			strings.HasPrefix(input.Path, "../") {
			return fmt.Errorf("malformed input path %q", input.Path)
		}
		if input.Role == "" {
			return fmt.Errorf("input %s lacks role", input.Path)
		}
		if !digestPattern.MatchString(input.Digest) {
			return fmt.Errorf("input %s has malformed digest", input.Path)
		}
		key := input.Path + "\x00" + input.Role
		if seen[key] {
			return fmt.Errorf("duplicate input %s (role %s)", input.Path, input.Role)
		}
		seen[key] = true
	}
	return nil
}

func validateConflicts(conflicts []CatalogConflict) error {
	seen := map[string]bool{}
	for _, conflict := range conflicts {
		if conflict.ID == "" || conflict.Code == "" ||
			!identifierPattern.MatchString(conflict.ID) ||
			!identifierPattern.MatchString(conflict.Code) {
			return fmt.Errorf("malformed conflict id=%q code=%q", conflict.ID, conflict.Code)
		}
		if len(conflict.SourcePaths) == 0 {
			return fmt.Errorf("conflict %s lacks source paths", conflict.ID)
		}
		for _, path := range conflict.SourcePaths {
			if path == "" || absPath(path) || escapes(path) {
				return fmt.Errorf("conflict %s has malformed source path %q", conflict.ID, path)
			}
		}
		if seen[conflict.ID] {
			return fmt.Errorf("duplicate conflict identity %q", conflict.ID)
		}
		seen[conflict.ID] = true
	}
	return nil
}

func absPath(path string) bool {
	return strings.HasPrefix(path, "/")
}

func escapes(path string) bool {
	return path == ".." || strings.HasPrefix(path, "../") ||
		strings.Contains(path, "/../") || strings.HasSuffix(path, "/..")
}

// CanonicalJSON encodes the catalog deterministically with the digest field
// computed over the canonical bytes with digest omitted.
func CanonicalJSON(envelope CatalogEnvelope, itemIDs []string) ([]byte, error) {
	if err := envelope.Validate(itemIDs); err != nil {
		return nil, err
	}
	withoutDigest := envelope
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

func digest(content []byte) string {
	value := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(value[:])
}

// verifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of a concrete catalog with the digest field omitted and compares it to
// the stored digest. An empty stored digest is rejected because the strict
// catalog contract is content-addressed. The caller must pass the concrete
// catalog value (with its catalog-specific fields) with Digest cleared, not
// the standalone envelope, so the recomputed digest covers the full document.
func verifySelfDigest(stored string, withoutDigest any) error {
	if stored == "" {
		return fmt.Errorf("catalog digest is required")
	}
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return fmt.Errorf("encode catalog for digest verification: %w", err)
	}
	expected := digest(content)
	if stored != expected {
		return fmt.Errorf(
			"catalog self-digest mismatch: stored %s, recomputed %s",
			stored,
			expected,
		)
	}
	return nil
}

// DecodeStrict decodes a catalog JSON document with unknown-field rejection
// and trailing-JSON rejection.
func DecodeStrict(content []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode catalog: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("catalog contains multiple JSON values")
		}
		return fmt.Errorf("decode trailing JSON: %w", err)
	}
	return nil
}

// SortedUnique returns a sorted, de-duplicated copy of a set-like slice.
func SortedUnique(values []string) []string {
	copied := append([]string(nil), values...)
	sort.Strings(copied)
	result := copied[:0]
	for _, value := range copied {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}

// DecodeTopologyStrict decodes and validates a complete topology catalog
// document with unknown-field and trailing-JSON rejection.
func DecodeTopologyStrict(content []byte) (TopologyCatalog, error) {
	var catalog TopologyCatalog
	if err := DecodeStrict(content, &catalog); err != nil {
		return TopologyCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return TopologyCatalog{}, err
	}
	if err := catalog.VerifySelfDigest(); err != nil {
		return TopologyCatalog{}, err
	}
	return catalog, nil
}
