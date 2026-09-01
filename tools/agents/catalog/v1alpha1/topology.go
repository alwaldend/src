package catalogv1alpha1

import (
	"encoding/json"
	"fmt"
)

const (
	// KindTopologyCatalog identifies the topology catalog kind.
	KindTopologyCatalog = "topology-catalog"
)

// TreeBoundaryClass classifies a top-level repository boundary.
type TreeBoundaryClass string

const (
	TreeBoundaryProduct    TreeBoundaryClass = "product"
	TreeBoundaryRepository TreeBoundaryClass = "repository_internal"
	TreeBoundaryTool       TreeBoundaryClass = "tool"
	TreeBoundaryData       TreeBoundaryClass = "data"
	TreeBoundaryThirdParty TreeBoundaryClass = "third_party"
	TreeBoundaryUser       TreeBoundaryClass = "user"
)

// TopologyTree is one top-level repository boundary tree.
type TopologyTree struct {
	ID         string            `json:"id"`
	Path       string            `json:"path"`
	ReadmePath string            `json:"readmePath"`
	Boundary   TreeBoundaryClass `json:"boundaryClass"`
}

// TopologyComponent is one registered maintained component.
type TopologyComponent struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	OwnerReadme string `json:"ownerReadmePath"`
	BuildPath   string `json:"buildPath"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Lifecycle   string `json:"lifecycle"`
	DocsState   string `json:"docsState"`
	DocsTarget  string `json:"docsTarget,omitempty"`
}

// TopologyWorkspace is one tracked Bzlmod workspace root.
type TopologyWorkspace struct {
	ID         string `json:"id"`
	Path       string `json:"path"`
	ModulePath string `json:"modulePath"`
	ModuleName string `json:"moduleName"`
}

// TopologyCatalog is the bounded topology projection over the registered
// projects universe and tracked workspaces.
type TopologyCatalog struct {
	CatalogEnvelope
	Trees      []TopologyTree      `json:"trees"`
	Components []TopologyComponent `json:"components"`
	Workspaces []TopologyWorkspace `json:"workspaces"`
}

// Validate checks the topology-specific invariants and the shared envelope.
func (catalog TopologyCatalog) Validate() error {
	if catalog.Trees == nil {
		return fmt.Errorf("trees must be a non-null array")
	}
	if catalog.Components == nil {
		return fmt.Errorf("components must be a non-null array")
	}
	if catalog.Workspaces == nil {
		return fmt.Errorf("workspaces must be a non-null array")
	}
	itemIDs := make([]string, 0,
		len(catalog.Trees)+len(catalog.Components)+len(catalog.Workspaces))
	for _, tree := range catalog.Trees {
		if !validTreeBoundary(tree.Boundary) {
			return fmt.Errorf("tree %s has unknown boundary %q", tree.ID, tree.Boundary)
		}
		itemIDs = append(itemIDs, "tree."+tree.ID)
	}
	for _, component := range catalog.Components {
		if component.Lifecycle == "" {
			return fmt.Errorf("component %s lacks lifecycle", component.ID)
		}
		if !oneOf(component.Lifecycle,
			"abandoned", "active", "experimental", "finished",
			"in_progress", "maintenance") {
			return fmt.Errorf("component %s has unknown lifecycle %q",
				component.ID, component.Lifecycle)
		}
		if !oneOf(component.DocsState, "owned", "absent", "unavailable") {
			return fmt.Errorf("component %s has unknown docsState %q",
				component.ID, component.DocsState)
		}
		itemIDs = append(itemIDs, "component."+component.ID)
	}
	for _, workspace := range catalog.Workspaces {
		if workspace.ModuleName == "" {
			return fmt.Errorf("workspace %s lacks module name", workspace.ID)
		}
		itemIDs = append(itemIDs, "workspace."+workspace.ID)
	}
	return catalog.CatalogEnvelope.Validate(itemIDs)
}

func validTreeBoundary(value TreeBoundaryClass) bool {
	switch value {
	case TreeBoundaryProduct, TreeBoundaryRepository, TreeBoundaryTool,
		TreeBoundaryData, TreeBoundaryThirdParty, TreeBoundaryUser:
		return true
	}
	return false
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

// ItemIDs returns the sorted set-like identity list for a catalog.
func (catalog TopologyCatalog) ItemIDs() []string {
	var ids []string
	for _, tree := range catalog.Trees {
		ids = append(ids, "tree."+tree.ID)
	}
	for _, component := range catalog.Components {
		ids = append(ids, "component."+component.ID)
	}
	for _, workspace := range catalog.Workspaces {
		ids = append(ids, "workspace."+workspace.ID)
	}
	return ids
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete topology document and compares it to the stored
// digest, rejecting an empty or mismatched digest.
func (catalog TopologyCatalog) VerifySelfDigest() error {
	withoutDigest := catalog
	withoutDigest.Digest = ""
	return verifySelfDigest(catalog.Digest, withoutDigest)
}

// CanonicalJSONTopology encodes a complete topology catalog with a digest
// computed over the canonical bytes with the digest field omitted.
func CanonicalJSONTopology(catalog TopologyCatalog) ([]byte, error) {
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := catalog
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode topology catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode topology catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// HasBoundaryClass reports whether a tree record with the given boundary
// class exists.
func (catalog TopologyCatalog) HasBoundaryClass(value TreeBoundaryClass) bool {
	for _, tree := range catalog.Trees {
		if tree.Boundary == value {
			return true
		}
	}
	return false
}
