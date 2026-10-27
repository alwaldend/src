package graph

import (
	"fmt"
	"strings"
)

// writeMindmap lists selected targets by ownership without inventing dependency edges.
func writeMindmap(output *strings.Builder, nodes map[string]node, labels []string, ids map[string]string, root string, grouped bool) error {
	output.WriteString("%%{init: {\"layout\": \"tidy-tree\"}}%%\nmindmap\n")
	fmt.Fprintf(output, "  o0[\"Build input inventory<br/>%s\"]\n", escapeLabel(root))
	if !grouped {
		for _, label := range labels {
			if !strings.HasPrefix(label, "@") {
				fmt.Fprintf(output, "    %s[\"%s<br/>%s\"]\n", ids[label], escapeLabel(label), escapeLabel(nodes[label].kind))
			}
		}
		externalOpen := false
		for _, label := range labels {
			if strings.HasPrefix(label, "@") {
				if !externalOpen {
					output.WriteString("    x0[\"External\"]\n")
					externalOpen = true
				}
				fmt.Fprintf(output, "      %s[\"%s<br/>%s\"]\n", ids[label], escapeLabel(label), escapeLabel(nodes[label].kind))
			}
		}
		return nil
	}
	tree, err := groupLabels(labels)
	if err != nil {
		return fmt.Errorf("organize ownership branches: %w", err)
	}
	packageIDs := make(map[string]string, len(tree.packageNames))
	for index, name := range tree.packageNames {
		packageIDs[name] = fmt.Sprintf("p%d", index)
	}
	repositoryIndex := 0
	externalOpen := false
	for _, repository := range tree.repositoryNames {
		indent := "      "
		if repository == "" {
			output.WriteString("    l0[\"Local packages\"]\n")
		} else {
			if !externalOpen {
				output.WriteString("    x0[\"External\"]\n")
				externalOpen = true
			}
			fmt.Fprintf(output, "      r%d[\"%s\"]\n", repositoryIndex, escapeLabel(repository))
			repositoryIndex++
			indent = "        "
		}
		for _, packageName := range tree.repositories[repository] {
			fmt.Fprintf(output, "%s%s[\"%s\"]\n", indent, packageIDs[packageName], escapeLabel(packageName))
			for _, label := range tree.packages[packageName] {
				fmt.Fprintf(output, "%s  %s[\"%s<br/>%s\"]\n", indent, ids[label], escapeLabel(tree.parts[label].targetName), escapeLabel(nodes[label].kind))
			}
		}
	}
	return nil
}
