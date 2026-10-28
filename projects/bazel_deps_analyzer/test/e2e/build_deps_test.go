package e2e

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// buildDepsInput is the hand-authored metadata and dependency fixture.
var buildDepsInput = flag.String("build-deps-input", "", "Explicit dependency capture runfile")

// buildDepsDiagram is the filtered target JSON-input rule artifact.
var buildDepsDiagram = flag.String("build-deps-diagram", "", "Filtered target diagram runfile")

// buildDepsJSON is the filtered package JSON-input rule artifact.
var buildDepsJSON = flag.String("build-deps-json", "", "Filtered package JSON runfile")

// buildDepsNative is the filtered native capture macro artifact.
var buildDepsNative = flag.String("build-deps-native", "", "Filtered native diagram runfile")

// buildDepsCapture is the metadata stream emitted by the native capture macro.
var buildDepsCapture = flag.String("build-deps-capture", "", "Filtered native capture runfile")

// buildDepsStringCapture contains a native non-label selector with a None default.
var buildDepsStringCapture = flag.String("build-deps-string-capture", "", "Native string-selector capture runfile")

// TestBuildDepsOnly verifies filtering before traversal and complete aggregation.
func TestBuildDepsOnly(t *testing.T) {
	_, input := packageArtifact(t, *buildDepsInput)
	args := []string{"--root=//app:root", "--include-files", "--build-deps-only"}
	output, err := targetCommand(t, string(input), args...)
	if err != nil {
		t.Fatal(err)
	}
	labels := []string{"//app:root", "//lib:kept", "//lib:alias", "//files:source", "//files:generated", "//producer:make"}
	edges := [][2]string{{"//app:root", "//lib:kept"}, {"//app:root", "//lib:alias"}, {"//app:root", "//files:source"}, {"//app:root", "//files:generated"}, {"//files:generated", "//producer:make"}}
	graph(t, output, labels, edges)
	if !bytes.Contains(output, []byte("%% build_deps_only=true\n")) {
		t.Fatal("filtered target output lacks mode comment")
	}
	limited, err := targetCommand(t, string(input), append(slices.Clone(args), "--depth=1")...)
	if err != nil {
		t.Fatal(err)
	}
	graph(t, limited, labels[:5], edges[:4])
	inventory, err := targetCommand(t, string(input), append(slices.Clone(args), "--diagram-type=mindmap")...)
	if err != nil {
		t.Fatal(err)
	}
	mindmap(t, inventory, false, labels)
	withoutFiles, err := targetCommand(t, string(input), "--root=//app:root", "--build-deps-only")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, withoutFiles, []string{"//app:root", "//lib:kept", "//lib:alias", "//producer:make"}, edges[:2])
	packages, err := command(t, string(input), "analyze", "package", "--build-deps-only")
	if err != nil {
		t.Fatal(err)
	}
	nodes := map[string][2]int{"//app": {1, 0}, "//lib": {6, 0}, "//runtime": {2, 0}, "//deep": {1, 0}, "//files": {2, 0}, "//producer": {1, 0}, "//meta": {3, 0}, "//metadata": {2, 0}}
	packageGraph(t, packages, nodes, []string{"//app -> //files", "//app -> //lib", "//files -> //producer", "//runtime -> //deep"}, 0, "selected_targets=18 selected_edges=6 missing_references=0 build_deps_only=true")
	jsonOutput, err := command(t, string(input), "analyze", "package", "--build-deps-only", "--format=json")
	if err != nil {
		t.Fatal(err)
	}
	var aggregate map[string]any
	if err := json.Unmarshal(jsonOutput, &aggregate); err != nil {
		t.Fatal(err)
	}
	meta := aggregate["meta"].(map[string]any)
	if meta["buildDepsOnly"] != true || meta["selectedTargets"] != float64(18) || meta["selectedEdges"] != float64(6) || meta["unresolvedReferences"] != float64(0) {
		t.Fatalf("incorrect filtered JSON metadata: %s", jsonOutput)
	}
	if got := fmt.Sprint(aggregate["edges"]); got != "[map[dependency:2 source:0 targetEdges:2] map[dependency:3 source:0 targetEdges:2] map[dependency:6 source:2 targetEdges:1] map[dependency:1 source:7 targetEdges:1]]" {
		t.Fatalf("incorrect JSON directed weights: %s", got)
	}
	for _, mode := range []string{"target", "package"} {
		flags := []string{"analyze", mode}
		if mode == "target" {
			flags = append(flags, "--root=//app:root", "--include-files")
		} else {
			flags = append(flags, "--format=json")
		}
		baseline, err := command(t, string(input), flags...)
		if err != nil {
			t.Fatal(err)
		}
		explicitDefault, err := command(t, string(input), append(slices.Clone(flags), "--build-deps-only=false")...)
		if err != nil || !bytes.Equal(baseline, explicitDefault) || bytes.Contains(baseline, []byte("buildDepsOnly")) || bytes.Contains(baseline, []byte("build_deps_only")) {
			t.Fatalf("changed default %s: %v", mode, err)
		}
		records := strings.Split(strings.TrimSpace(string(input)), "\n")
		records = append(records, records[0])
		slices.Reverse(records)
		reordered, err := command(t, strings.Join(records, "\n"), append(flags, "--build-deps-only")...)
		want := output
		if mode == "package" {
			want = jsonOutput
		}
		if err != nil || !bytes.Equal(want, reordered) {
			t.Fatalf("reordered/duplicate metadata changed %s: %v", mode, err)
		}
	}
	for location, want := range map[string][]byte{*buildDepsDiagram: output, *buildDepsJSON: jsonOutput} {
		_, actual := packageArtifact(t, location)
		if !bytes.Equal(actual, want) {
			t.Fatalf("native JSON-input rule differs: %s", location)
		}
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		for name, data := range map[string][]byte{"build-deps.ndjson": input, "build-deps.mmd": output, "build-deps-packages.mmd": packages, "build-deps-packages.json": jsonOutput, "build-deps-report.json": []byte("{\"passed\":true,\"selectedTargets\":18,\"selectedEdges\":6,\"targetClosureNodes\":6,\"targetClosureEdges\":5}\n")} {
			if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// dependencyRecord wraps selected metadata in a rule and one captured dependency.
func dependencyRecord(attribute string) string {
	return `{"type":"RULE","rule":{"name":"//app:root","ruleClass":"fixture","ruleInput":["//lib:kept","//condition:flag"],"attribute":[{"name":"name"},` + attribute + `]}}
{"type":"RULE","rule":{"name":"//lib:kept","ruleClass":"filegroup","attribute":[{"name":"name"}]}}`
}

// TestBuildDepsEncodings checks label-bearing values without condition/key guessing.
func TestBuildDepsEncodings(t *testing.T) {
	empty := `{"type":"RULE","rule":{"name":"//app:root","ruleClass":"filegroup","attribute":[{"name":"name"},{"name":"srcs","type":"LABEL_LIST","explicitlySpecified":true}]}}`
	explicitEmpty := `{"type":"RULE","rule":{"name":"//app:root","ruleClass":"filegroup","ruleInput":[],"attribute":[{"name":"name"},{"name":"srcs","type":"LABEL_LIST","explicitlySpecified":true,"stringListValue":[]}]}}`
	for _, mode := range []string{"target", "package"} {
		args := []string{"analyze", mode, "--build-deps-only"}
		if mode == "target" {
			args = append(args, "--root=//app:root")
		}
		baseline, err := command(t, empty, args...)
		if err != nil {
			t.Fatal(err)
		}
		duplicate, err := command(t, empty+"\n"+explicitEmpty, args...)
		if err != nil || !bytes.Equal(baseline, duplicate) {
			t.Fatalf("omitted/empty inputs or attribute values conflict: %v", err)
		}
	}
	for _, encoding := range []struct {
		// kind identifies the protobuf attribute discriminator.
		kind string
		// value supplies the corresponding typed JSON value.
		value string
	}{
		{"LABEL", `"stringValue":"@@//lib:kept"`},
		{"LABEL_LIST", `"stringListValue":["@//lib:kept","//lib:kept"]`},
		{"LABEL_DICT_UNARY", `"labelDictUnaryValue":[{"key":"//not:a_label","value":"//lib:kept"}]`},
		{"LABEL_LIST_DICT", `"stringListDictValue":[{"key":"//not:a_label","value":["//lib:kept"]}]`},
		{"LABEL_KEYED_STRING_DICT", `"labelKeyedStringDictValue":[{"key":"//lib:kept","value":"//not:a_label"}]`},
	} {
		for _, selector := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/selector=%t", encoding.kind, selector), func(t *testing.T) {
				kind, value := encoding.kind, encoding.value
				if selector {
					kind = "SELECTOR_LIST"
					value = `"selectorList":{"type":"` + encoding.kind + `","elements":[{"entries":[{"label":"//condition:flag",` + encoding.value + `},{"label":"//conditions:default","isDefaultValue":true}]}]}`
				}
				input := dependencyRecord(`{"name":"deps","type":"` + kind + `","explicitlySpecified":true,` + value + `}`)
				output, err := targetCommand(t, input, "--root=//app:root", "--build-deps-only")
				if err != nil {
					t.Fatal(err)
				}
				graph(t, output, []string{"//app:root", "//lib:kept"}, [][2]string{{"//app:root", "//lib:kept"}})
			})
		}
	}
	for _, attribute := range []string{
		`{"name":"deps","type":"STRING","explicitlySpecified":true,"stringValue":"//lib:kept"}`,
		`{"name":"srcs","type":"STRING_LIST","explicitlySpecified":true,"stringListValue":["//lib:kept"]}`,
		`{"name":"srcs","type":"OUTPUT","explicitlySpecified":true,"stringValue":"//lib:kept"}`,
		`{"name":"srcs","type":"OUTPUT_LIST","explicitlySpecified":true,"stringListValue":["//lib:kept"]}`,
		`{"name":"deps","type":"NODEP_LABEL","explicitlySpecified":true,"stringValue":"//lib:kept"}`,
		`{"name":"deps","type":"NODEP_LABEL_LIST","explicitlySpecified":true,"stringListValue":["//lib:kept"]}`,
		`{"name":"deps","type":"LABEL","explicitlySpecified":false}`,
		`{"name":"deps","type":"LABEL","explicitlySpecified":true}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"sourceAspectName":"aspect","stringListValue":["//lib:kept"]}`,
	} {
		output, err := targetCommand(t, dependencyRecord(attribute), "--root=//app:root", "--build-deps-only")
		if err != nil {
			t.Fatal(err)
		}
		graph(t, output, []string{"//app:root"}, nil)
	}
}

// TestBuildDepsFailures verifies malformed metadata and old captures preserve output.
func TestBuildDepsFailures(t *testing.T) {
	valid := dependencyRecord(`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"stringListValue":["//lib:kept"]}`)
	inputs := []string{
		`{"type":"RULE","rule":{"name":"//app:root","ruleClass":"filegroup"}}`,
		strings.ReplaceAll(valid, `{"name":"name"},`, ""),
		valid + "\n" + `{"type":"RULE","rule":{"name":"//disconnected:old","ruleClass":"filegroup"}}`,
		valid + "\n" + strings.Replace(valid, `"explicitlySpecified":true`, `"explicitlySpecified":false`, 1),
		valid + "\n" + strings.Replace(valid, `"//condition:flag"`, `"//different:input"`, 1),
	}
	for _, attribute := range []string{
		`{"name":"deps","type":"LABEL_LIST","stringListValue":["//lib:kept"]}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"stringListValue":"//lib:kept"}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"stringListValue":[""]}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"stringListValue":[null]}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"stringListValue":["relative"]}`,
		`{"name":"deps","type":"LABEL","explicitlySpecified":true,"stringValue":null}`,
		`{"name":"deps","type":"FUTURE_LABEL","explicitlySpecified":true}`,
		`{"name":"deps","explicitlySpecified":true}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"nodep":"false"}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"nodep":null}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"sourceAspectName":null}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true,"stringValue":"//lib:kept"}`,
		`{"name":"deps","type":"LABEL_LIST_DICT","explicitlySpecified":true,"labelListDictValue":[{"key":"name","value":["//lib:kept"]}]}`,
		`{"name":"deps","type":"LABEL_DICT_UNARY","explicitlySpecified":true,"labelDictUnaryValue":[{"key":"name"}]}`,
		`{"name":"deps","type":"LABEL_LIST_DICT","explicitlySpecified":true,"stringListDictValue":[{"value":["//lib:kept"]}]}`,
		`{"name":"deps","type":"LABEL_KEYED_STRING_DICT","explicitlySpecified":true,"labelKeyedStringDictValue":[{"key":"//lib:kept"}]}`,
		`{"name":"deps","type":"SELECTOR_LIST","explicitlySpecified":true}`,
		`{"name":"deps","type":"SELECTOR_LIST","explicitlySpecified":true,"selectorList":{"type":"LABEL","elements":[{"entries":[{"label":"//condition:flag"}]}]}}`,
		`{"name":"deps","type":"SELECTOR_LIST","explicitlySpecified":true,"selectorList":{"type":"LABEL","elements":[{"entries":[{"stringValue":"//lib:kept"}]}]}}`,
		`{"name":"deps","type":"LABEL_LIST","explicitlySpecified":true},{"name":"deps","type":"LABEL_LIST","explicitlySpecified":false}`,
	} {
		inputs = append(inputs, dependencyRecord(attribute))
	}
	for index, input := range inputs {
		for _, mode := range []string{"target", "package"} {
			t.Run(fmt.Sprintf("%d/%s", index, mode), func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "existing.out")
				if err := os.WriteFile(path, []byte("preserved"), 0o644); err != nil {
					t.Fatal(err)
				}
				args := []string{"analyze", mode, "--build-deps-only", "--output=" + path}
				if mode == "target" {
					args = append(args, "--root=//app:root", "--depth=0")
				}
				stdout, err := command(t, input, args...)
				if err == nil || len(stdout) != 0 {
					t.Fatalf("invalid metadata accepted: %s", input)
				}
				if index < 3 {
					for _, flag := range []string{"--noimplicit_deps", "--notool_deps", "--nonodep_deps", "--proto:output_rule_attrs=name,deps,srcs", "--proto:default_values=true", "--proto:flatten_selects=false", "--proto:rule_inputs_and_outputs=true"} {
						if !strings.Contains(err.Error(), flag) {
							t.Fatalf("missing recapture flag %s: %v", flag, err)
						}
					}
				}
				data, readErr := os.ReadFile(path)
				if readErr != nil || string(data) != "preserved" {
					t.Fatalf("failure replaced output: %v", readErr)
				}
			})
		}
	}
}

// TestBuildDepsNative verifies real Bazel name/default/selector metadata and flags.
func TestBuildDepsNative(t *testing.T) {
	_, capture := packageArtifact(t, *buildDepsCapture)
	_, output := packageArtifact(t, *buildDepsNative)
	prefix := "//projects/bazel_deps_analyzer/test/e2e:"
	graph(t, output, []string{prefix + "filter_inputs", prefix + "explicit_input", prefix + "shared"}, [][2]string{{prefix + "filter_inputs", prefix + "explicit_input"}, {prefix + "filter_inputs", prefix + "shared"}})
	for _, marker := range []string{`"name":"name"`, `"explicitlySpecified":false`, `"type":"SELECTOR_LIST"`} {
		if !bytes.Contains(capture, []byte(marker)) {
			t.Fatalf("native metadata missing %s", marker)
		}
	}
	for _, excluded := range []string{"tool_input", "implicit_input", "fixture_visibility"} {
		if bytes.Contains(capture, []byte(`"name":"`+prefix+excluded+`"`)) || bytes.Contains(capture, []byte(`"name":"@@`+prefix+excluded+`"`)) {
			t.Fatalf("capture flags retained excluded record %s", excluded)
		}
	}
	converted, err := targetCommand(t, string(capture), "--root="+prefix+"filter_inputs", "--build-deps-only")
	if err != nil || !bytes.Equal(output, converted) {
		t.Fatalf("native capture differs from command: %v", err)
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		for name, data := range map[string][]byte{"build-deps-native.ndjson": capture, "build-deps-native.mmd": output} {
			if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	_, stringCapture := packageArtifact(t, *buildDepsStringCapture)
	for _, marker := range []string{`"type":"SELECTOR_LIST"`, `"type":"STRING"`, `"stringValue":""`, `"isDefaultValue":true`} {
		if !bytes.Contains(stringCapture, []byte(marker)) {
			t.Fatalf("native None selector fixture missing %s: %s", marker, stringCapture)
		}
	}
	stringOutput, err := targetCommand(t, string(stringCapture), "--root="+prefix+"filter_string_inputs", "--build-deps-only")
	if err != nil {
		t.Fatal(err)
	}
	graph(t, stringOutput, []string{prefix + "filter_string_inputs"}, nil)
	for _, mode := range []string{"target", "package"} {
		args := []string{"analyze", mode, "--build-deps-only"}
		if mode == "target" {
			args = append(args, "--root="+prefix+"filter_inputs")
		}
		if _, err := command(t, string(capture)+"\n"+string(stringCapture), args...); err != nil {
			t.Fatalf("disconnected non-label selector breaks %s: %v", mode, err)
		}
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "build-deps-string-native.ndjson"), stringCapture, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
