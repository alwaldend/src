package graph

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

// PackageOptions is the command boundary for capture-wide ownership aggregation.
type PackageOptions struct {
	// BuildDepsOnly selects explicit deps/srcs inputs before package aggregation.
	BuildDepsOnly bool
	// PackageDepth is the positive number of local package path segments retained.
	PackageDepth int
	// Format selects Mermaid presentation or the shared JSON aggregate.
	Format string
	// Direction specifies the Mermaid flowchart orientation.
	Direction string
	// IncludeFiles retains source and generated files in the selected graph.
	IncludeFiles bool
	// IncludeExternal retains all external repositories in one aggregate.
	IncludeExternal bool
	// MergeReciprocal combines opposite Mermaid arrows without changing JSON.
	MergeReciprocal bool
}

// packageGraph is the single aggregate serialized by Mermaid and JSON consumers.
type packageGraph struct {
	// Groups exports the sorted owner nodes and membership contract.
	Groups []packageSummary `json:"groups"`
	// Edges exports every distinct directed inter-group relationship.
	Edges []packageEdge `json:"edges"`
	// Meta exports selection and capture-boundary counts.
	Meta packageMetadata `json:"meta"`
}

// packageSummary serializes one disjoint owner and its captured target members.
type packageSummary struct {
	// Caption is root, a qualified local package prefix, or External.
	Caption string `json:"caption"`
	// TargetCount is the number of unique selected targets in the owner.
	TargetCount int `json:"targetCount"`
	// InternalEdges counts selected target edges within the same owner.
	InternalEdges int `json:"internalEdges"`
	// Members exports sorted normalized labels for consumers such as the viewer.
	Members []string `json:"members"`
}

// packageEdge serializes a directed pair of group-array indices and its weight.
type packageEdge struct {
	// Source is the zero-based index of the depending group.
	Source int `json:"source"`
	// Dependency is the zero-based index of the dependency group.
	Dependency int `json:"dependency"`
	// TargetEdges counts distinct captured target edges represented by this pair.
	TargetEdges int `json:"targetEdges"`
}

// packageMetadata serializes the aggregate's selection and known capture limits.
type packageMetadata struct {
	// BuildDepsOnly identifies structural filtering without changing default JSON.
	BuildDepsOnly bool `json:"buildDepsOnly,omitempty"`
	// PackageDepth records the selected local hierarchy depth.
	PackageDepth int `json:"packageDepth"`
	// SelectedTargets counts included unique captured records.
	SelectedTargets int `json:"selectedTargets"`
	// SelectedEdges counts included internal and inter-group target edges.
	SelectedEdges int `json:"selectedEdges"`
	// UnresolvedReferences counts distinct absent labels referenced by selected sources.
	UnresolvedReferences int `json:"unresolvedReferences"`
	// IncludeFiles records the source/generated-file selection.
	IncludeFiles bool `json:"includeFiles"`
	// IncludeExternal records the external-repository selection.
	IncludeExternal bool `json:"includeExternal"`
}

// ConvertPackages constructs one capture-wide aggregate for either serializer.
// It does not write files; callers retain input and output lifetime ownership.
func ConvertPackages(input io.Reader, options PackageOptions) ([]byte, error) {
	if options.PackageDepth < 1 {
		return nil, errors.New("package depth must be a positive integer")
	}
	if options.Format != "mermaid" && options.Format != "json" {
		return nil, fmt.Errorf("unsupported package output format %q", options.Format)
	}
	switch options.Direction {
	case "LR", "RL", "TB", "BT":
	default:
		return nil, fmt.Errorf("unsupported direction %q", options.Direction)
	}
	nodes, err := readGraph(input, options.BuildDepsOnly)
	if err != nil {
		return nil, fmt.Errorf("read package target graph: %w", err)
	}
	aggregate, err := aggregatePackages(nodes, options)
	if err != nil {
		return nil, fmt.Errorf("aggregate captured packages: %w", err)
	}
	if options.Format == "json" {
		output, err := json.MarshalIndent(aggregate, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("encode package graph JSON: %w", err)
		}
		return append(output, '\n'), nil
	}
	return packageMermaid(aggregate, options), nil
}

// aggregatePackages counts the selected induced graph before any presentation.
func aggregatePackages(nodes map[string]node, options PackageOptions) (packageGraph, error) {
	result := packageGraph{
		Groups: []packageSummary{},
		Edges:  []packageEdge{},
		Meta: packageMetadata{
			BuildDepsOnly:   options.BuildDepsOnly,
			PackageDepth:    options.PackageDepth,
			IncludeFiles:    options.IncludeFiles,
			IncludeExternal: options.IncludeExternal,
		},
	}
	owners := make(map[string]string, len(nodes))
	groups := map[string]*packageSummary{}
	for label, item := range nodes {
		owner, err := packageOwner(label, options.PackageDepth)
		if err != nil {
			return packageGraph{}, fmt.Errorf("identify owner of %q: %w", label, err)
		}
		if (item.file && !options.IncludeFiles) || (owner == "External" && !options.IncludeExternal) {
			continue
		}
		owners[label] = owner
		if groups[owner] == nil {
			groups[owner] = &packageSummary{Caption: owner}
		}
		groups[owner].Members = append(groups[owner].Members, label)
	}
	result.Meta.SelectedTargets = len(owners)
	edges := map[[2]string]int{}
	missing := map[string]bool{}
	for label, owner := range owners {
		for _, dependency := range nodes[label].inputs {
			if _, exists := nodes[dependency]; !exists {
				missing[dependency] = true
				continue
			}
			dependencyOwner, included := owners[dependency]
			if !included {
				continue
			}
			result.Meta.SelectedEdges++
			if owner == dependencyOwner {
				groups[owner].InternalEdges++
			} else {
				edges[[2]string{owner, dependencyOwner}]++
			}
		}
	}
	result.Meta.UnresolvedReferences = len(missing)
	for _, group := range groups {
		slices.Sort(group.Members)
		group.TargetCount = len(group.Members)
		result.Groups = append(result.Groups, *group)
	}
	slices.SortFunc(result.Groups, func(a, b packageSummary) int {
		return strings.Compare(a.Caption, b.Caption)
	})
	indices := make(map[string]int, len(result.Groups))
	for index, group := range result.Groups {
		indices[group.Caption] = index
	}
	for pair, count := range edges {
		result.Edges = append(result.Edges, packageEdge{
			Source: indices[pair[0]], Dependency: indices[pair[1]], TargetEdges: count,
		})
	}
	slices.SortFunc(result.Edges, func(a, b packageEdge) int {
		if order := cmp.Compare(a.Source, b.Source); order != 0 {
			return order
		}
		return cmp.Compare(a.Dependency, b.Dependency)
	})
	return result, nil
}

// packageOwner derives one disjoint caption at the requested local hierarchy depth.
func packageOwner(label string, depth int) (string, error) {
	parts, err := parseTargetGroup(label)
	if err != nil {
		return "", fmt.Errorf("parse package identity: %w", err)
	}
	if parts.repository != "" {
		return "External", nil
	}
	path := strings.TrimPrefix(parts.packageName, "//")
	if path == "" {
		return "root", nil
	}
	segments := strings.Split(path, "/")
	return "//" + strings.Join(segments[:min(depth, len(segments))], "/"), nil
}

// packageMermaid presents an aggregate without edge captions or target self-loops.
func packageMermaid(aggregate packageGraph, options PackageOptions) []byte {
	var output strings.Builder
	output.WriteString("%% Captured package aggregate; arrows point to dependencies.\n")
	output.WriteString("%% Aggregate paths do not establish target-level reachability.\n")
	writeBuildDepsComment(&output, options.BuildDepsOnly)
	fmt.Fprintf(&output, "%%%% package_depth=%d include_files=%t include_external=%t merge_reciprocal=%t\n", options.PackageDepth, options.IncludeFiles, options.IncludeExternal, options.MergeReciprocal)
	fmt.Fprintf(&output, "%%%% selected_targets=%d selected_edges=%d missing_references=%d\n", aggregate.Meta.SelectedTargets, aggregate.Meta.SelectedEdges, aggregate.Meta.UnresolvedReferences)
	output.WriteString(packageLayoutDirective)
	fmt.Fprintf(&output, "flowchart %s\n", options.Direction)
	for index, group := range aggregate.Groups {
		fmt.Fprintf(&output, "    n%d[\"%s<br/>%d targets<br/>%d internal edges\"]\n", index, escapeLabel(group.Caption), group.TargetCount, group.InternalEdges)
	}
	connections := make(map[[2]int]bool, len(aggregate.Edges))
	for _, edge := range aggregate.Edges {
		connections[[2]int{edge.Source, edge.Dependency}] = true
	}
	for _, edge := range aggregate.Edges {
		arrow := "-->"
		if options.MergeReciprocal && connections[[2]int{edge.Dependency, edge.Source}] {
			if edge.Source > edge.Dependency {
				continue
			}
			arrow = "<-->"
		}
		fmt.Fprintf(&output, "    n%d %s n%d\n", edge.Source, arrow, edge.Dependency)
	}
	return []byte(output.String())
}
