package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// converter names the built command in the runfiles tree.
var converter = flag.String("converter", "", "Converter runfile")

// diagram names the Mermaid output generated from the real Bazel fixture.
var diagram = flag.String("diagram", "", "Fixture Mermaid runfile")

// raster names the rendered fixture image.
var raster = flag.String("raster", "", "Fixture WebP runfile")

// buildInputs names the graph captured from inputs, tools, and visibility metadata.
var buildInputs = flag.String("build-inputs", "", "Build-input fixture Mermaid runfile")

// groupedInput names the captured fixture with repeated target names and repositories.
var groupedInput = flag.String("grouped-input", "", "Grouped fixture JSON runfile")

// groupedDiagram names the rule-generated grouped Mermaid fixture.
var groupedDiagram = flag.String("grouped-diagram", "", "Grouped fixture Mermaid runfile")

// groupedRaster names the native dark render of the grouped fixture.
var groupedRaster = flag.String("grouped-raster", "", "Grouped fixture WebP runfile")

// mindmapDiagram names the ownership inventory captured from the real Bazel graph.
var mindmapDiagram = flag.String("mindmap-diagram", "", "Native-query mindmap runfile")

// mindmapRaster names the rendered native-query ownership inventory.
var mindmapRaster = flag.String("mindmap-raster", "", "Native-query mindmap WebP runfile")

// mindmapGrouped names the JSON-input ownership inventory with external packages.
var mindmapGrouped = flag.String("mindmap-grouped", "", "Grouped mindmap runfile")

// mindmapGroupedRaster names the dark render of the external ownership hierarchy.
var mindmapGroupedRaster = flag.String("mindmap-grouped-raster", "", "Grouped mindmap WebP runfile")

// nodePattern finds generated identifiers and their complete quoted labels.
var nodePattern = regexp.MustCompile(`(?m)^[ \t]+(n[0-9]+)\["([^"]*)"\]$`)

// groupPattern finds generated package or repository groups and their captions.
var groupPattern = regexp.MustCompile(`^subgraph ([prx][0-9]+)\["([^"]*)"\]$`)

// mindmapPattern identifies ownership branches and target captions with indentation.
var mindmapPattern = regexp.MustCompile(`^( +)([olxprn][0-9]+)\["([^"]*)"\]$`)

// command executes the public CLI with bounded execution and returns its output.
func command(t *testing.T, input string, args ...string) ([]byte, error) {
	t.Helper()
	binary, err := runfiles.Rlocation(*converter)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("converter: %w: %s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}

// targetCommand executes the target converter with the public subcommand.
func targetCommand(t *testing.T, input string, args ...string) ([]byte, error) {
	t.Helper()
	output, err := command(t, input, append([]string{"analyze", "target"}, args...)...)
	if err != nil {
		return output, fmt.Errorf("analyze target command: %w", err)
	}
	return output, nil
}

// graph checks semantic labels and edges without depending on identifier order.
func graph(t *testing.T, output []byte, labels []string, edges [][2]string) {
	t.Helper()
	ids := map[string]string{}
	for _, match := range nodePattern.FindAllStringSubmatch(string(output), -1) {
		label := strings.Split(match[2], "<br/>")[0]
		if _, exists := ids[label]; exists {
			t.Fatalf("duplicate node %q", label)
		}
		ids[label] = match[1]
	}
	if len(ids) != len(labels) {
		t.Fatalf("nodes %v, want %v; diagram:\n%s", ids, labels, output)
	}
	for _, label := range labels {
		if ids[label] == "" {
			t.Fatalf("missing node %q in %s", label, output)
		}
	}
	if strings.Count(string(output), " --> ") != len(edges) {
		t.Fatalf("incorrect edge count in %s", output)
	}
	for _, edge := range edges {
		want := "    " + ids[edge[0]] + " --> " + ids[edge[1]] + "\n"
		if !bytes.Contains(output, []byte(want)) {
			t.Fatalf("missing edge %v in %s", edge, output)
		}
	}
}

// TestCommand verifies graph selection, canonical output, and CLI failures.
func TestCommand(t *testing.T) {
	records := []string{
		`{"type":"RULE","rule":{"name":"//fixture:root","ruleClass":"filegroup","ruleInput":["//fixture:b","//fixture:a","@ext//:lib","//fixture:source","//fixture:output","//fixture:absent"]}}`,
		`{"type":"RULE","rule":{"name":"//fixture:a","ruleClass":"filegroup","ruleInput":["//fixture:shared","//fixture:shared"]}}`,
		`{"type":"RULE","rule":{"name":"//fixture:b","ruleClass":"filegroup","ruleInput":["//fixture:shared"]}}`,
		`{"type":"RULE","rule":{"name":"//fixture:shared","ruleClass":"filegroup","ruleInput":["//fixture:root"]}}`,
		`{"type":"RULE","rule":{"name":"@ext//:lib","ruleClass":"go_library"}}`,
		`{"type":"SOURCE_FILE","sourceFile":{"name":"//fixture:source"}}`,
		`{"type":"GENERATED_FILE","generatedFile":{"name":"//fixture:output","generatingRule":"//fixture:producer"}}`,
		`{"type":"RULE","rule":{"name":"//fixture:producer","ruleClass":"write_file"}}`,
	}
	input := strings.Join(records, "\n")
	args := []string{"--root=//fixture:root"}
	output, err := targetCommand(t, input, args...)
	if err != nil {
		t.Fatal(err)
	}
	graph(t, output, []string{"//fixture:root", "//fixture:a", "//fixture:b", "//fixture:shared", "//fixture:producer"}, [][2]string{
		{"//fixture:root", "//fixture:a"},
		{"//fixture:root", "//fixture:b"},
		{"//fixture:a", "//fixture:shared"},
		{"//fixture:b", "//fixture:shared"},
		{"//fixture:shared", "//fixture:root"},
	})
	slices.Reverse(records)
	reordered, err := targetCommand(t, strings.Join(records, "\n"), args...)
	if err != nil || !bytes.Equal(output, reordered) {
		t.Fatalf("record ordering changed output: %v", err)
	}
	for _, direction := range []string{"LR", "RL", "TB", "BT"} {
		limited, err := targetCommand(t, input, "--root=@@//fixture:root", "--depth=0", "--direction="+direction)
		if err != nil {
			t.Fatal(err)
		}
		graph(t, limited, []string{"//fixture:root"}, nil)
		if !bytes.Contains(limited, []byte("flowchart "+direction+"\n")) {
			t.Fatalf("direction %s missing: %s", direction, limited)
		}
	}
	limited, err := targetCommand(t, input, "--root=//fixture:root", "--depth=1")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, limited, []string{"//fixture:root", "//fixture:a", "//fixture:b"}, [][2]string{
		{"//fixture:root", "//fixture:a"}, {"//fixture:root", "//fixture:b"},
	})
	complete, err := targetCommand(t, input, "--root=//fixture:root", "--include-files", "--include-external")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(complete, []byte("@ext//:lib")) || !bytes.Contains(complete, []byte("//fixture:output")) || strings.Count(string(complete), " --> ") != 9 {
		t.Fatalf("file/external edges missing: %s", complete)
	}
	generated, err := targetCommand(t, input, "--root=//fixture:output", "--include-files")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, generated, []string{"//fixture:output", "//fixture:producer"}, [][2]string{{"//fixture:output", "//fixture:producer"}})
	groups := `{"type":"PACKAGE_GROUP","packageGroup":{"name":"//fixture:group","includedPackageGroup":["//fixture:included"]}}
{"type":"PACKAGE_GROUP","packageGroup":{"name":"//fixture:included"}}
{"type":"ENVIRONMENT_GROUP","environmentGroup":{"name":"//fixture:environments","environment":["//fixture:environment"]}}
{"type":"RULE","rule":{"name":"//fixture:environment","ruleClass":"environment"}}`
	for _, pair := range [][2]string{{"//fixture:group", "//fixture:included"}, {"//fixture:environments", "//fixture:environment"}} {
		group, err := targetCommand(t, groups, "--root="+pair[0])
		if err != nil {
			t.Fatal(err)
		}
		graph(t, group, []string{pair[0], pair[1]}, [][2]string{pair})
	}
	for _, invalid := range []string{"", `{`, `[]`, `{}`, `null`, `{"type":"RULE","rule":{}}`, input + ` {`, input + ` {"type":"RULE","rule":{"name":"//fixture:root","ruleClass":"other"}}`} {
		result, err := targetCommand(t, invalid, args...)
		if err == nil || len(result) != 0 {
			t.Fatalf("invalid input succeeded or wrote partial output: %q: %v %s", invalid, err, result)
		}
	}
	for _, invalidArgs := range [][]string{{"--root=//missing:root"}, {"--root=//fixture:root", "--depth=-2"}, {"--root=//fixture:root", "--direction=invalid"}, {}} {
		if result, err := targetCommand(t, input, invalidArgs...); err == nil || len(result) != 0 {
			t.Fatalf("invalid arguments succeeded: %v: %v %s", invalidArgs, err, result)
		}
	}
	resultPath := filepath.Join(t.TempDir(), "output.mmd")
	if _, err := targetCommand(t, `{`, "--root=//fixture:root", "--output="+resultPath); err == nil {
		t.Fatal("malformed input accepted")
	}
	if _, err := os.Stat(resultPath); !os.IsNotExist(err) {
		t.Fatalf("partial output created: %v", err)
	}
	label := `//fixture:quote"<b>&#35;[x]%%{init}`
	escaped, err := targetCommand(t, `{"type":"RULE","rule":{"name":"//fixture:quote\"<b>&#35;[x]%%{init}","ruleClass":"filegroup"}}`, "--root="+label)
	if err != nil || bytes.Contains(escaped, []byte(`<b>`)) || bytes.Contains(escaped, []byte(`%%{init}`)) || !bytes.Contains(escaped, []byte(`#34;`)) || !bytes.Contains(escaped, []byte(`#38;#35;`)) {
		t.Fatalf("unsafe label output: %v: %s", err, escaped)
	}
}

// TestExternalLeaves checks the opt-in boundary without changing existing modes.
func TestExternalLeaves(t *testing.T) {
	input := `{"type":"RULE","rule":{"name":"//fixture:root","ruleClass":"filegroup","ruleInput":["//fixture:local","@@ext+//:first","@other//:first"]}}
{"type":"RULE","rule":{"name":"//fixture:local","ruleClass":"filegroup","ruleInput":["@@ext+//:first","//fixture:shared"]}}
{"type":"RULE","rule":{"name":"@@ext+//:first","ruleClass":"filegroup","ruleInput":["@other//:first","@@ext+//:transitive","//fixture:external_only","//fixture:shared"]}}
{"type":"RULE","rule":{"name":"@other//:first","ruleClass":"filegroup","ruleInput":["//fixture:root"]}}
{"type":"RULE","rule":{"name":"@@ext+//:transitive","ruleClass":"filegroup"}}
{"type":"RULE","rule":{"name":"//fixture:external_only","ruleClass":"filegroup"}}
{"type":"RULE","rule":{"name":"//fixture:shared","ruleClass":"filegroup"}}`
	leaves, err := targetCommand(t, input, "--root=//fixture:root", "--external-leaves")
	if err != nil {
		t.Fatal(err)
	}
	labels := []string{"//fixture:root", "//fixture:local", "//fixture:shared", "@@ext+//:first", "@other//:first"}
	edges := [][2]string{
		{"//fixture:root", "//fixture:local"},
		{"//fixture:root", "@@ext+//:first"},
		{"//fixture:root", "@other//:first"},
		{"//fixture:local", "@@ext+//:first"},
		{"//fixture:local", "//fixture:shared"},
	}
	graph(t, leaves, labels, edges)
	externalContainer(t, leaves)
	combined, err := targetCommand(t, input, "--root=//fixture:root", "--include-external", "--external-leaves")
	if err != nil || !bytes.Equal(leaves, combined) {
		t.Fatalf("explicit inclusion changed external-leaf behavior: %v", err)
	}
	limited, err := targetCommand(t, input, "--root=//fixture:root", "--external-leaves", "--depth=1")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, limited, []string{"//fixture:root", "//fixture:local", "@@ext+//:first", "@other//:first"}, edges[:4])
	root, err := targetCommand(t, input, "--root=@@ext+//:first", "--external-leaves")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, root, []string{"@@ext+//:first"}, nil)
	externalContainer(t, root)
	defaults, err := targetCommand(t, input, "--root=//fixture:root")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, defaults, []string{"//fixture:root", "//fixture:local", "//fixture:shared", "//fixture:external_only"}, [][2]string{
		{"//fixture:root", "//fixture:local"}, {"//fixture:local", "//fixture:shared"},
	})
	if bytes.Contains(defaults, []byte(`subgraph x0["External"]`)) {
		t.Fatal("hidden external targets left an empty container")
	}
	expanded, err := targetCommand(t, input, "--root=//fixture:root", "--include-external")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, expanded, append(labels, "@@ext+//:transitive", "//fixture:external_only"), append(edges,
		[2]string{"@@ext+//:first", "@other//:first"},
		[2]string{"@@ext+//:first", "@@ext+//:transitive"},
		[2]string{"@@ext+//:first", "//fixture:external_only"},
		[2]string{"@@ext+//:first", "//fixture:shared"},
		[2]string{"@other//:first", "//fixture:root"},
	))
	externalContainer(t, expanded)
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "external-leaves.mmd"), leaves, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// externalContainer checks the semantic external boundary in a flat diagram.
func externalContainer(t *testing.T, output []byte) {
	t.Helper()
	var stack []string
	containers := 0
	for _, line := range strings.Split(string(output), "\n") {
		trimmed := strings.TrimSpace(line)
		if group := groupPattern.FindStringSubmatch(trimmed); group != nil {
			stack = append(stack, group[2])
			if group[2] == "External" {
				containers++
			}
		} else if trimmed == "end" {
			if len(stack) == 0 {
				t.Fatal("unmatched container end")
			}
			stack = stack[:len(stack)-1]
		} else if node := nodePattern.FindStringSubmatch(line); node != nil {
			external := strings.HasPrefix(node[2], "@")
			if external != slices.Contains(stack, "External") {
				t.Fatalf("incorrect external classification: %s in %v", node[2], stack)
			}
		}
	}
	if containers != 1 || len(stack) != 0 {
		t.Fatalf("invalid external containers: count=%d stack=%v", containers, stack)
	}
}

// groupedGraph resolves package membership before checking full identities and edges.
func groupedGraph(t *testing.T, output []byte) {
	t.Helper()
	var stack, flat []string
	externalGroups := 0
	for _, line := range strings.Split(string(output), "\n") {
		trimmed := strings.TrimSpace(line)
		if match := groupPattern.FindStringSubmatch(trimmed); match != nil {
			if match[2] == "External" {
				externalGroups++
			}
			stack = append(stack, match[2])
			continue
		}
		if trimmed == "end" {
			if len(stack) == 0 {
				t.Fatal("unmatched group end")
			}
			stack = stack[:len(stack)-1]
			continue
		}
		if match := nodePattern.FindStringSubmatch(line); match != nil {
			if len(stack) == 0 {
				t.Fatalf("ungrouped node %s", line)
			}
			pkg := stack[len(stack)-1]
			repository, _, qualified := strings.Cut(pkg, "//")
			if !qualified || (repository == "" && len(stack) != 1) || (repository != "" && (len(stack) != 3 || stack[0] != "External" || stack[1] != repository)) {
				t.Fatalf("incorrect package nesting: %v", stack)
			}
			flat = append(flat, "    "+match[1]+`["`+pkg+":"+match[2]+`"]`)
		} else if strings.Contains(line, " --> ") {
			flat = append(flat, line)
		}
	}
	if len(stack) != 0 || externalGroups != 1 {
		t.Fatalf("invalid package groups: stack=%v external=%d", stack, externalGroups)
	}
	graph(t, []byte(strings.Join(flat, "\n")+"\n"), []string{
		"//app:root", "//app:left", "//lib:shared", "@@ext_one+//client:shared",
		"@@ext_one+//:shared", "@ext_two//client:shared", "@@ext_one+//client:api.proto",
	}, [][2]string{
		{"//app:root", "//app:left"},
		{"//app:root", "//lib:shared"},
		{"//app:root", "@@ext_one+//client:shared"},
		{"//app:root", "@@ext_one+//:shared"},
		{"//app:root", "@ext_two//client:shared"},
		{"//app:root", "@@ext_one+//client:api.proto"},
		{"//app:left", "//lib:shared"},
	})
}

// TestPackageGrouping checks nested namespaces, short captions, and stable identities.
func TestPackageGrouping(t *testing.T) {
	path, err := runfiles.Rlocation(*groupedInput)
	if err != nil {
		t.Fatal(err)
	}
	input, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"--root=//app:root", "--external-leaves", "--include-files", "--group-by-package"}
	output, err := targetCommand(t, string(input), args...)
	if err != nil {
		t.Fatal(err)
	}
	groupedGraph(t, output)
	records := strings.Split(strings.TrimSpace(string(input)), "\n")
	slices.Reverse(records)
	reordered, err := targetCommand(t, strings.Join(records, "\n"), args...)
	if err != nil || !bytes.Equal(output, reordered) {
		t.Fatalf("group ordering changed: %v", err)
	}
	flatArgs := args[:len(args)-1]
	flat, err := targetCommand(t, string(input), flatArgs...)
	if err != nil || bytes.Contains(flat, []byte("subgraph p")) || bytes.Contains(flat, []byte("subgraph r")) || !bytes.Contains(flat, []byte("//app:root<br/>")) {
		t.Fatalf("flat default changed: %v: %s", err, flat)
	}
	externalContainer(t, flat)
	explicitFlat, err := targetCommand(t, string(input), append(flatArgs, "--group-by-package=false")...)
	if err != nil || !bytes.Equal(flat, explicitFlat) {
		t.Fatalf("explicit flat output differs: %v", err)
	}
	hostile := `{"type":"RULE","rule":{"name":"//quote\"<b>[pkg]:end","ruleClass":"filegroup"}}`
	escaped, err := targetCommand(t, hostile, `--root=//quote"<b>[pkg]:end`, "--group-by-package")
	if err != nil || bytes.Contains(escaped, []byte("<b>")) || !bytes.Contains(escaped, []byte(`#34;#60;b#62;#91;pkg#93;`)) {
		t.Fatalf("unsafe package caption: %v: %s", err, escaped)
	}
}

// mindmap checks unique target identities and their actual ownership ancestry.
func mindmap(t *testing.T, output []byte, grouped bool, want []string) {
	t.Helper()
	if !bytes.Contains(output, []byte("\nmindmap\n")) || !bytes.Contains(output, []byte("Build input inventory")) || bytes.Contains(output, []byte(" --> ")) {
		t.Fatalf("output is not an ownership inventory: %s", output)
	}
	var stack, labels []string
	externalBranches, externalTargets := 0, 0
	for _, line := range strings.Split(string(output), "\n") {
		match := mindmapPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		depth := len(match[1]) / 2
		if len(match[1])%2 != 0 || depth < 1 || depth > len(stack)+1 {
			t.Fatalf("invalid ownership indentation: %q", line)
		}
		stack = append(stack[:depth-1], strings.Split(match[3], "<br/>")[0])
		if match[2] == "x0" {
			externalBranches++
			if len(stack) != 2 || stack[1] != "External" {
				t.Fatalf("incorrect External branch: %v", stack)
			}
		}
		if !strings.HasPrefix(match[2], "n") {
			continue
		}
		label := stack[len(stack)-1]
		if grouped {
			if len(stack) < 4 {
				t.Fatalf("missing package ancestry: %v", stack)
			}
			pkg := stack[len(stack)-2]
			repository, _, qualified := strings.Cut(pkg, "//")
			if !qualified || (repository == "" && (len(stack) != 4 || stack[1] != "Local packages")) || (repository != "" && (len(stack) != 5 || stack[1] != "External" || stack[2] != repository)) {
				t.Fatalf("incorrect ownership ancestry: %v", stack)
			}
			label = pkg + ":" + label
		} else if (strings.HasPrefix(label, "@") && (len(stack) != 3 || stack[1] != "External")) || (!strings.HasPrefix(label, "@") && len(stack) != 2) {
			t.Fatalf("incorrect flat ownership ancestry: %v", stack)
		}
		if strings.HasPrefix(label, "@") {
			externalTargets++
		}
		labels = append(labels, label)
	}
	slices.Sort(labels)
	want = slices.Clone(want)
	slices.Sort(want)
	if !slices.Equal(labels, want) || externalBranches > 1 || (externalTargets > 0) != (externalBranches == 1) {
		t.Fatalf("inventory targets=%v, want=%v; external branches=%d", labels, want, externalBranches)
	}
}

// TestMindmap verifies selection, ownership, determinism, escaping, and safe failures.
func TestMindmap(t *testing.T) {
	path, err := runfiles.Rlocation(*groupedInput)
	if err != nil {
		t.Fatal(err)
	}
	input, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"--root=//app:root", "--external-leaves", "--include-files"}
	defaults, err := targetCommand(t, string(input), args...)
	if err != nil {
		t.Fatal(err)
	}
	explicit, err := targetCommand(t, string(input), append(slices.Clone(args), "--diagram-type=flowchart")...)
	if err != nil || !bytes.Equal(defaults, explicit) {
		t.Fatalf("explicit flowchart differs from default: %v", err)
	}
	args = append(args, "--diagram-type=mindmap")
	labels := []string{"//app:root", "//app:left", "//lib:shared", "@@ext_one+//client:shared", "@@ext_one+//:shared", "@ext_two//client:shared", "@@ext_one+//client:api.proto"}
	for _, grouped := range []bool{false, true} {
		modeArgs := append(slices.Clone(args), fmt.Sprintf("--group-by-package=%t", grouped))
		output, err := targetCommand(t, string(input), modeArgs...)
		if err != nil {
			t.Fatal(err)
		}
		mindmap(t, output, grouped, labels)
		records := strings.Split(strings.TrimSpace(string(input)), "\n")
		slices.Reverse(records)
		reordered, err := targetCommand(t, strings.Join(records, "\n"), modeArgs...)
		if err != nil || !bytes.Equal(output, reordered) {
			t.Fatalf("mindmap changed with record order: %v", err)
		}
	}
	for _, selection := range []struct {
		// args chooses a root and display boundary.
		args []string
		// labels identifies the targets that remain visible.
		labels []string
	}{
		{[]string{"--root=//app:root", "--depth=0"}, []string{"//app:root"}},
		{[]string{"--root=//app:root"}, labels[:3]},
		{[]string{"--root=//app:root", "--external-leaves"}, labels[:6]},
		{[]string{"--root=//app:root", "--include-files", "--include-external"}, append(slices.Clone(labels), "@@ext_one+//internal:hidden")},
		{[]string{"--root=@@ext_one+//client:shared", "--external-leaves"}, []string{"@@ext_one+//client:shared"}},
	} {
		output, err := targetCommand(t, string(input), append(selection.args, "--diagram-type=mindmap", "--group-by-package")...)
		if err != nil {
			t.Fatal(err)
		}
		mindmap(t, output, true, selection.labels)
	}
	cycle := `{"type":"RULE","rule":{"name":"//cycle:a","ruleClass":"filegroup","ruleInput":["//cycle:b","//cycle:b"]}}
{"type":"RULE","rule":{"name":"//cycle:b","ruleClass":"filegroup","ruleInput":["//cycle:a"]}}`
	output, err := targetCommand(t, cycle, "--root=//cycle:a", "--diagram-type=mindmap")
	if err != nil {
		t.Fatal(err)
	}
	mindmap(t, output, false, []string{"//cycle:a", "//cycle:b"})
	hostile := `{"type":"RULE","rule":{"name":"//quote\"<b>[pkg]:end","ruleClass":"filegroup"}}`
	escaped, err := targetCommand(t, hostile, `--root=//quote"<b>[pkg]:end`, "--diagram-type=mindmap", "--group-by-package")
	if err != nil || bytes.Contains(escaped, []byte("<b>")) || !bytes.Contains(escaped, []byte(`#34;#60;b#62;#91;pkg#93;`)) {
		t.Fatalf("unsafe mindmap caption: %v: %s", err, escaped)
	}
	outputPath := filepath.Join(t.TempDir(), "invalid.mmd")
	invalidArgs := []string{"--root=//app:root", "--diagram-type=invalid", "--output=" + outputPath}
	if result, err := targetCommand(t, string(input), invalidArgs...); err == nil || len(result) != 0 {
		t.Fatalf("unknown diagram type succeeded: %v", err)
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Fatalf("invalid type created output: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("existing diagram"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := targetCommand(t, string(input), invalidArgs...); err == nil {
		t.Fatal("unknown diagram type accepted")
	}
	unchanged, err := os.ReadFile(outputPath)
	if err != nil || string(unchanged) != "existing diagram" {
		t.Fatalf("invalid type replaced existing output: %v", err)
	}
}

// TestBuildPipeline verifies a real query fixture and retains repeatable artifacts.
func TestBuildPipeline(t *testing.T) {
	artifacts := map[string]string{
		"fixture.mmd": *diagram, "fixture.webp": *raster, "build-inputs.mmd": *buildInputs,
		"grouped.mmd": *groupedDiagram, "grouped.dark.webp": *groupedRaster,
		"mindmap.mmd": *mindmapDiagram, "mindmap.dark.webp": *mindmapRaster,
		"mindmap-grouped.mmd": *mindmapGrouped, "mindmap-grouped.dark.webp": *mindmapGroupedRaster,
	}
	graphs := map[string]map[string]int{}
	for name, location := range artifacts {
		path, err := runfiles.Rlocation(location)
		if err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if name == "mindmap-grouped.mmd" {
			mindmap(t, content, true, []string{"//app:root", "//app:left", "//lib:shared", "@@ext_one+//client:shared", "@@ext_one+//:shared", "@ext_two//client:shared", "@@ext_one+//client:api.proto"})
		} else if name == "mindmap.mmd" {
			prefix := "//projects/bazel_deps_analyzer/test/e2e:"
			mindmap(t, content, true, []string{prefix + "root", prefix + "left", prefix + "right", prefix + "shared"})
		} else if name == "grouped.mmd" {
			groupedGraph(t, content)
		} else if name == "build-inputs.mmd" {
			prefix := "//projects/bazel_deps_analyzer/test/e2e:"
			graph(t, content, []string{prefix + "build_inputs", prefix + "explicit_input", prefix + "implicit_input"}, [][2]string{
				{prefix + "build_inputs", prefix + "explicit_input"},
				{prefix + "build_inputs", prefix + "implicit_input"},
			})
		} else if strings.HasSuffix(name, ".mmd") {
			prefix := "//projects/bazel_deps_analyzer/test/e2e:"
			graph(t, content, []string{prefix + "root", prefix + "left", prefix + "right", prefix + "shared"}, [][2]string{
				{prefix + "root", prefix + "left"},
				{prefix + "root", prefix + "right"},
				{prefix + "left", prefix + "shared"},
				{prefix + "right", prefix + "shared"},
			})
		} else if len(content) < 1000 || string(content[:4]) != "RIFF" || string(content[8:12]) != "WEBP" {
			t.Fatalf("invalid WebP artifact %q", name)
		}
		if strings.HasSuffix(name, ".mmd") {
			graphs[name] = map[string]int{"nodes": len(nodePattern.FindAllSubmatch(content, -1)), "edges": strings.Count(string(content), " --> ")}
		}
		if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
			if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		names := make([]string, 0, len(artifacts))
		for name := range artifacts {
			names = append(names, name)
		}
		slices.Sort(names)
		report, err := json.MarshalIndent(map[string]any{
			"passed": true, "nodes": graphs["fixture.mmd"]["nodes"], "edges": graphs["fixture.mmd"]["edges"],
			"graphs": graphs, "artifacts": names,
		}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "report.json"), append(report, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
