package graph

import (
	"fmt"
	"slices"
)

// packageLayoutDirective selects the shared automatic grouped-graph layout.
const packageLayoutDirective = "%%{init: {\"layout\": \"elk\", \"elk\": {\"nodePlacementStrategy\": \"NETWORK_SIMPLEX\", \"mergeEdges\": true}}}%%\n"

// packageTree organizes selected target identities for either diagram serializer.
type packageTree struct {
	// parts holds each target's parsed repository, package, and short name.
	parts map[string]targetGroup
	// packages maps a fully qualified package to its selected target labels.
	packages map[string][]string
	// repositories maps a repository to its sorted full package names.
	repositories map[string][]string
	// packageNames contains sorted package identities for stable identifiers.
	packageNames []string
	// repositoryNames contains sorted repository identities, including the main one.
	repositoryNames []string
}

// groupLabels derives ownership once from labels already sorted by identity.
func groupLabels(labels []string) (packageTree, error) {
	tree := packageTree{
		parts:        make(map[string]targetGroup, len(labels)),
		packages:     map[string][]string{},
		repositories: map[string][]string{},
	}
	for _, label := range labels {
		group, err := parseTargetGroup(label)
		if err != nil {
			return packageTree{}, fmt.Errorf("group target %q: %w", label, err)
		}
		tree.parts[label] = group
		if _, exists := tree.packages[group.packageName]; !exists {
			tree.packageNames = append(tree.packageNames, group.packageName)
			tree.repositories[group.repository] = append(tree.repositories[group.repository], group.packageName)
		}
		tree.packages[group.packageName] = append(tree.packages[group.packageName], label)
	}
	slices.Sort(tree.packageNames)
	for name, packages := range tree.repositories {
		tree.repositoryNames = append(tree.repositoryNames, name)
		slices.Sort(packages)
	}
	slices.Sort(tree.repositoryNames)
	return tree, nil
}
