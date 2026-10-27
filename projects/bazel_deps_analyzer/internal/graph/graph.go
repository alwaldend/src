package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
)

// Options defines graph selection, ownership grouping, and the output view.
type Options struct {
	// Root is the complete label of the target whose dependencies are selected.
	Root string
	// Depth is the maximum dependency distance, or -1 for the full closure.
	Depth int
	// DiagramType selects a flowchart (the default) or mindmap ownership inventory.
	DiagramType string
	// Direction is the flowchart's LR, RL, TB, or BT orientation.
	Direction string
	// IncludeFiles retains source and generated file nodes.
	IncludeFiles bool
	// IncludeExternal retains targets outside the main repository.
	IncludeExternal bool
	// ExternalLeaves shows external targets and stops traversal at their boundary.
	ExternalLeaves bool
	// GroupByPackage groups short target captions within full package namespaces.
	GroupByPackage bool
}

// record decodes the graph-bearing subset of Bazel's Target JSON message.
type record struct {
	// Type identifies which target payload Bazel emitted.
	Type string `json:"type"`
	// Rule describes a rule and its dependency labels.
	Rule *ruleRecord `json:"rule"`
	// SourceFile describes a source file target.
	SourceFile *namedRecord `json:"sourceFile"`
	// GeneratedFile links an output file to its generating rule.
	GeneratedFile *generatedRecord `json:"generatedFile"`
	// PackageGroup describes a visibility package group and its inclusions.
	PackageGroup *packageGroupRecord `json:"packageGroup"`
	// EnvironmentGroup describes a Bazel environment group and its members.
	EnvironmentGroup *environmentGroupRecord `json:"environmentGroup"`
}

// namedRecord decodes target messages whose graph needs only their label.
type namedRecord struct {
	// Name is the complete Bazel target label.
	Name string `json:"name"`
}

// packageGroupRecord decodes visibility groups and their group dependencies.
type packageGroupRecord struct {
	// Name is the complete group label.
	Name string `json:"name"`
	// IncludedPackageGroup names the groups included by this group.
	IncludedPackageGroup []string `json:"includedPackageGroup"`
}

// environmentGroupRecord decodes environment groups and their member labels.
type environmentGroupRecord struct {
	// Name is the complete environment group label.
	Name string `json:"name"`
	// Environment names the environments belonging to this group.
	Environment []string `json:"environment"`
}

// ruleRecord decodes rule identity and dependency edges without attribute data.
type ruleRecord struct {
	// Name is the complete Bazel rule label.
	Name string `json:"name"`
	// RuleClass is the rule's kind, such as go_binary or filegroup.
	RuleClass string `json:"ruleClass"`
	// RuleInput contains target labels on which this rule depends.
	RuleInput []string `json:"ruleInput"`
}

// generatedRecord decodes a generated file and its producer dependency.
type generatedRecord struct {
	// Name is the complete generated-file label.
	Name string `json:"name"`
	// GeneratingRule names the rule responsible for this file.
	GeneratingRule string `json:"generatingRule"`
}

// node stores normalized identity and edges for one target.
type node struct {
	// label is the complete normalized target label.
	label string
	// kind is the rule class or target category shown under the label.
	kind string
	// inputs contains unique normalized dependency labels in lexical order.
	inputs []string
	// file identifies source and generated file targets for display filtering.
	file bool
}

// Convert reads a stream of Bazel Target JSON messages and returns a complete
// Mermaid document. It returns no document on failure and never writes files.
func Convert(input io.Reader, options Options) ([]byte, error) {
	if options.Root == "" {
		return nil, errors.New("root target is required")
	}
	if options.Depth < -1 {
		return nil, errors.New("depth must be -1 or nonnegative")
	}
	if options.DiagramType == "" {
		options.DiagramType = "flowchart"
	}
	if options.DiagramType != "flowchart" && options.DiagramType != "mindmap" {
		return nil, fmt.Errorf("unsupported diagram type %q", options.DiagramType)
	}
	switch options.Direction {
	case "LR", "RL", "TB", "BT":
	default:
		return nil, fmt.Errorf("unsupported direction %q", options.Direction)
	}
	nodes, err := readGraph(input)
	if err != nil {
		return nil, fmt.Errorf("read target graph: %w", err)
	}
	root := normalizeLabel(options.Root)
	if _, exists := nodes[root]; !exists {
		return nil, fmt.Errorf("root target %q is absent from the graph", root)
	}
	selected := selectNodes(nodes, root, options)
	labels := make([]string, 0, len(selected))
	for label := range selected {
		labels = append(labels, label)
	}
	slices.Sort(labels)
	ids := make(map[string]string, len(labels))
	var output strings.Builder
	if options.DiagramType == "mindmap" {
		output.WriteString("%% Bazel build-input inventory; branches show ownership, not dependency edges.\n")
	} else {
		output.WriteString("%% Bazel dependency graph; arrows point to dependencies.\n")
	}
	fmt.Fprintf(&output, "%%%% depth=%d include_files=%t include_external=%t\n", options.Depth, options.IncludeFiles, options.IncludeExternal || options.ExternalLeaves)
	if options.ExternalLeaves {
		output.WriteString("%% external_leaves=true\n")
	}
	if options.GroupByPackage {
		output.WriteString("%% group_by_package=true\n")
	}
	for index, label := range labels {
		ids[label] = fmt.Sprintf("n%d", index)
	}
	if options.DiagramType == "mindmap" {
		if err := writeMindmap(&output, nodes, labels, ids, root, options.GroupByPackage); err != nil {
			return nil, fmt.Errorf("write mindmap: %w", err)
		}
		return []byte(output.String()), nil
	}
	if options.GroupByPackage {
		output.WriteString("%%{init: {\"layout\": \"elk\", \"elk\": {\"nodePlacementStrategy\": \"NETWORK_SIMPLEX\", \"mergeEdges\": true}}}%%\n")
	}
	fmt.Fprintf(&output, "flowchart %s\n", options.Direction)
	if options.GroupByPackage {
		if err := writeGroupedNodes(&output, nodes, labels, ids); err != nil {
			return nil, fmt.Errorf("write package groups: %w", err)
		}
	} else {
		writeFlatNodes(&output, nodes, labels, ids)
	}
	for _, label := range labels {
		if options.ExternalLeaves && strings.HasPrefix(label, "@") {
			continue
		}
		for _, dependency := range nodes[label].inputs {
			if id, exists := ids[dependency]; exists {
				fmt.Fprintf(&output, "    %s --> %s\n", ids[label], id)
			}
		}
	}
	return []byte(output.String()), nil
}

// writeFlatNodes marks external targets without adding package groups.
func writeFlatNodes(output *strings.Builder, nodes map[string]node, labels []string, ids map[string]string) {
	for _, label := range labels {
		if !strings.HasPrefix(label, "@") {
			fmt.Fprintf(output, "    %s[\"%s<br/>%s\"]\n", ids[label], escapeLabel(label), escapeLabel(nodes[label].kind))
		}
	}
	externalOpen := false
	for _, label := range labels {
		if strings.HasPrefix(label, "@") {
			if !externalOpen {
				output.WriteString("    subgraph x0[\"External\"]\n")
				externalOpen = true
			}
			fmt.Fprintf(output, "        %s[\"%s<br/>%s\"]\n", ids[label], escapeLabel(label), escapeLabel(nodes[label].kind))
		}
	}
	if externalOpen {
		output.WriteString("    end\n")
	}
}

// targetGroup identifies a target's repository, package, and short caption.
type targetGroup struct {
	// repository is the external repository prefix, or empty for the main repository.
	repository string
	// packageName is the full repository-qualified package label without a target.
	packageName string
	// targetName is the target component, including any package-relative subpath.
	targetName string
}

// parseTargetGroup splits an absolute Bazel label without conflating repositories.
func parseTargetGroup(label string) (targetGroup, error) {
	repository, remainder, absolute := strings.Cut(label, "//")
	if !absolute || (repository != "" && !strings.HasPrefix(repository, "@")) {
		return targetGroup{}, errors.New("package grouping requires an absolute Bazel label")
	}
	packagePath, targetName, explicit := strings.Cut(remainder, ":")
	if !explicit {
		targetName = packagePath[strings.LastIndex(packagePath, "/")+1:]
	}
	if targetName == "" {
		return targetGroup{}, errors.New("package grouping requires a nonempty target name")
	}
	return targetGroup{repository: repository, packageName: repository + "//" + packagePath, targetName: targetName}, nil
}

// writeGroupedNodes nests external repositories and packages inside External.
func writeGroupedNodes(output *strings.Builder, nodes map[string]node, labels []string, ids map[string]string) error {
	tree, err := groupLabels(labels)
	if err != nil {
		return fmt.Errorf("organize package groups: %w", err)
	}
	packageIDs := make(map[string]string, len(tree.packageNames))
	for index, name := range tree.packageNames {
		packageIDs[name] = fmt.Sprintf("p%d", index)
	}
	repositoryIndex := 0
	externalOpen := false
	for _, repository := range tree.repositoryNames {
		indent := "    "
		if repository != "" {
			if !externalOpen {
				output.WriteString("    subgraph x0[\"External\"]\n")
				externalOpen = true
			}
			fmt.Fprintf(output, "        subgraph r%d[\"%s\"]\n", repositoryIndex, escapeLabel(repository))
			repositoryIndex++
			indent = "            "
		}
		for _, packageName := range tree.repositories[repository] {
			fmt.Fprintf(output, "%ssubgraph %s[\"%s\"]\n", indent, packageIDs[packageName], escapeLabel(packageName))
			for _, label := range tree.packages[packageName] {
				fmt.Fprintf(output, "%s    %s[\"%s<br/>%s\"]\n", indent, ids[label], escapeLabel(tree.parts[label].targetName), escapeLabel(nodes[label].kind))
			}
			fmt.Fprintf(output, "%send\n", indent)
		}
		if repository != "" {
			output.WriteString("        end\n")
		}
	}
	if externalOpen {
		output.WriteString("    end\n")
	}
	return nil
}

// readGraph decodes the entire stream before any output can be returned.
func readGraph(input io.Reader) (map[string]node, error) {
	decoder := json.NewDecoder(input)
	nodes := map[string]node{}
	for index := 1; ; index++ {
		var target record
		if err := decoder.Decode(&target); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode target %d at byte %d: %w", index, decoder.InputOffset(), err)
		}
		item, err := targetNode(target)
		if err != nil {
			return nil, fmt.Errorf("target %d: %w", index, err)
		}
		if previous, exists := nodes[item.label]; exists && (previous.kind != item.kind || previous.file != item.file || !slices.Equal(previous.inputs, item.inputs)) {
			return nil, fmt.Errorf("conflicting records for target %q", item.label)
		}
		nodes[item.label] = item
	}
	if len(nodes) == 0 {
		return nil, errors.New("no target records in input")
	}
	return nodes, nil
}

// targetNode validates a target's discriminant and normalizes its graph fields.
func targetNode(target record) (node, error) {
	var item node
	switch {
	case target.Type == "RULE" && target.Rule != nil:
		item = node{label: target.Rule.Name, kind: target.Rule.RuleClass, inputs: target.Rule.RuleInput}
	case target.Type == "SOURCE_FILE" && target.SourceFile != nil:
		item = node{label: target.SourceFile.Name, kind: "source file", file: true}
	case target.Type == "GENERATED_FILE" && target.GeneratedFile != nil:
		item = node{label: target.GeneratedFile.Name, kind: "generated file", file: true, inputs: []string{target.GeneratedFile.GeneratingRule}}
	case target.Type == "PACKAGE_GROUP" && target.PackageGroup != nil:
		item = node{label: target.PackageGroup.Name, kind: "package group", inputs: target.PackageGroup.IncludedPackageGroup}
	case target.Type == "ENVIRONMENT_GROUP" && target.EnvironmentGroup != nil:
		item = node{label: target.EnvironmentGroup.Name, kind: "environment group", inputs: target.EnvironmentGroup.Environment}
	default:
		return node{}, fmt.Errorf("unsupported or missing target payload for type %q", target.Type)
	}
	item.label = normalizeLabel(item.label)
	if item.label == "" || item.kind == "" {
		return node{}, errors.New("target name and kind must be nonempty")
	}
	for index, label := range item.inputs {
		if label == "" {
			return node{}, fmt.Errorf("empty dependency label on %q", item.label)
		}
		item.inputs[index] = normalizeLabel(label)
	}
	slices.Sort(item.inputs)
	item.inputs = slices.Compact(item.inputs)
	return item, nil
}

// selectNodes computes shortest dependency distances, then applies display filters.
func selectNodes(nodes map[string]node, root string, options Options) map[string]bool {
	distances := map[string]int{root: 0}
	queue := []string{root}
	selected := map[string]bool{}
	for head := 0; head < len(queue); head++ {
		label := queue[head]
		item := nodes[label]
		external := strings.HasPrefix(label, "@")
		if label == root || ((!item.file || options.IncludeFiles) && (!external || options.IncludeExternal || options.ExternalLeaves)) {
			selected[label] = true
		}
		if (options.ExternalLeaves && external) || (options.Depth >= 0 && distances[label] >= options.Depth) {
			continue
		}
		for _, dependency := range item.inputs {
			if _, exists := nodes[dependency]; !exists {
				continue
			}
			if _, visited := distances[dependency]; visited {
				continue
			}
			distances[dependency] = distances[label] + 1
			queue = append(queue, dependency)
		}
	}
	return selected
}

// normalizeLabel equates the consistent and ordinary spellings of main labels.
func normalizeLabel(label string) string {
	if strings.HasPrefix(label, "@@//") {
		return strings.TrimPrefix(label, "@@")
	}
	if strings.HasPrefix(label, "@//") {
		return strings.TrimPrefix(label, "@")
	}
	return label
}

// escapeLabel encodes syntax and HTML characters using Mermaid numeric entities.
func escapeLabel(label string) string {
	var output strings.Builder
	for _, character := range label {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || strings.ContainsRune(" /@._:+-", character) {
			output.WriteRune(character)
		} else {
			fmt.Fprintf(&output, "#%d;", character)
		}
	}
	return output.String()
}
