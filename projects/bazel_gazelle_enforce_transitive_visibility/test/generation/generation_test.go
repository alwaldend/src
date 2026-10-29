package generation

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	// gazellePath selects the plugin-only executable provided by Bazel.
	gazellePath = flag.String("gazelle", "", "Gazelle executable runfile")
	// goGazellePath selects the executable that also creates Go packages.
	goGazellePath = flag.String("go-gazelle", "", "Go Gazelle executable runfile")
)

// TestGeneration verifies path selection, ordering, preservation, and idempotence.
func TestGeneration(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"MODULE.bazel":                "module(name = \"fixture\")\n",
		"BUILD.bazel":                 "# Workspace policy.\npackage_group(name = \"consumers\", packages = [\"//...\"])\n",
		"infra/BUILD.bazel":           "package_group(name = \"consumers\", packages = [\"//infra/...\"])\n",
		"infra/child/BUILD":           "load(\"//:defs.bzl\", \"artifact\")\nartifact(name = \"before\")\n# Keep the package comment.\npackage(\n    default_visibility = [\"//visibility:public\"], # keep\n    default_testonly = True,\n    transitive_visibility = \"//:obsolete\",\n)\nfilegroup(name = \"certificate\", srcs = [\"certificate.txt\"])\n",
		"infra/child/certificate.txt": "public fixture certificate\n",
		"infra/override/BUILD.bazel":  "filegroup(name = \"public\")\n",
		"infra/sibling/BUILD.bazel":   "filegroup(name = \"sibling\")\n",
		"projects/BUILD.bazel":        "filegroup(name = \"project\")\n",
		"infrastructure/BUILD.bazel":  "filegroup(name = \"sibling_prefix\")\n",
		"tools/transitive_visibility/transitive_visibility.toml": "[paths]\n\".\" = \"//:consumers\"\ninfra = \"//infra:consumers\"\n\"infra/override\" = \"//:consumers\"\n",
		"infra/empty/README.txt":                                 "No BUILD file should be generated here.\n",
	}
	writeFiles(t, root, files)
	runGazelle(t, root, *gazellePath, false)
	want := map[string]string{
		"BUILD.bazel":                "//:consumers",
		"infra/BUILD.bazel":          "//infra:consumers",
		"infra/child/BUILD":          "//infra:consumers",
		"infra/override/BUILD.bazel": "//:consumers",
		"infra/sibling/BUILD.bazel":  "//infra:consumers",
		"projects/BUILD.bazel":       "//:consumers",
		"infrastructure/BUILD.bazel": "//:consumers",
	}
	before := make(map[string]string)
	for name, policy := range want {
		got := readFile(t, root, name)
		if !strings.Contains(got, "transitive_visibility = \""+policy+"\"") {
			t.Errorf("%s: missing policy %q in:\n%s", name, policy, got)
		}
		if strings.Count(got, "package(") != 1 {
			t.Errorf("%s: expected one package declaration:\n%s", name, got)
		}
		assertPackageOrder(t, name, got)
		before[name] = got
	}
	child := before["infra/child/BUILD"]
	for _, preserved := range []string{
		"# Keep the package comment.",
		"default_visibility = [\"//visibility:public\"],  # keep",
		"default_testonly = True",
		"name = \"certificate\"",
		"srcs = [\"certificate.txt\"]",
	} {
		if !strings.Contains(child, preserved) {
			t.Errorf("child package lost %q:\n%s", preserved, child)
		}
	}
	for _, name := range []string{"infra/empty/BUILD", "infra/empty/BUILD.bazel"} {
		if _, err := os.Stat(filepath.Join(root, name)); !os.IsNotExist(err) {
			t.Errorf("unexpected generated package %s: %v", name, err)
		}
	}
	runGazelle(t, root, *gazellePath, false)
	for name, original := range before {
		if got := readFile(t, root, name); got != original {
			t.Errorf("second run changed %s:\n%s", name, got)
		}
	}
	retainArtifacts(t, "inheritance", before)
}

// TestNoConfiguration preserves existing policy without configuration.
func TestNoConfiguration(t *testing.T) {
	root := t.TempDir()
	input := "package(transitive_visibility = \"//:manual\")\n"
	writeFiles(t, root, map[string]string{
		"MODULE.bazel": "module(name = \"fixture\")\n",
		"BUILD.bazel":  input,
	})
	runGazelle(t, root, *gazellePath, false)
	if got := readFile(t, root, "BUILD.bazel"); got != input {
		t.Fatalf("unmanaged package changed:\n%s", got)
	}
}

// TestNewPackage covers packages produced by an earlier language extension.
func TestNewPackage(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"MODULE.bazel":       "module(name = \"fixture\")\n",
		"BUILD.bazel":        "# gazelle:prefix example.com/fixture\npackage_group(name = \"consumers\", packages = [\"//...\"])\n",
		"library/library.go": "package library\n",
		"tools/transitive_visibility/transitive_visibility.toml": "[paths]\n\".\" = \"//:consumers\"\n",
	})
	runGazelle(t, root, *goGazellePath, false)
	got := readFile(t, root, "library/BUILD.bazel")
	for _, fragment := range []string{"go_library(", "transitive_visibility = \"//:consumers\""} {
		if !strings.Contains(got, fragment) {
			t.Errorf("generated package lacks %q:\n%s", fragment, got)
		}
	}
	// Upstream Gazelle appends the declaration in a newly created BUILD file.
	// A later run has a file to reposition; subsequent runs must be stable.
	runGazelle(t, root, *goGazellePath, false)
	ordered := readFile(t, root, "library/BUILD.bazel")
	assertPackageOrder(t, "library/BUILD.bazel", ordered)
	runGazelle(t, root, *goGazellePath, false)
	if final := readFile(t, root, "library/BUILD.bazel"); final != ordered {
		t.Fatalf("repeat run changed existing package:\n%s", final)
	}
	retainArtifacts(t, "new-package", map[string]string{
		"first/BUILD.bazel":   got,
		"ordered/BUILD.bazel": ordered,
	})
}

// TestInvalidConfiguration rejects every invalid mapping before writing files.
func TestInvalidConfiguration(t *testing.T) {
	for name, configuration := range map[string]string{
		"syntax":            "[paths\n",
		"unknown-field":     "typo = true\n[paths]\ninfra = \"//:consumers\"\n",
		"empty":             "[paths]\n",
		"relative-label":    "[paths]\ninfra = \":consumers\"\n",
		"empty-label":       "[paths]\ninfra = \"\"\n",
		"multiple-labels":   "[paths]\ninfra = \"//:one //:two\"\n",
		"parent-path":       "[paths]\n\"../infra\" = \"//:consumers\"\n",
		"absolute-path":     "[paths]\n\"/infra\" = \"//:consumers\"\n",
		"noncanonical-path": "[paths]\n\"infra/../projects\" = \"//:consumers\"\n",
		"duplicate-path":    "[paths]\ninfra = \"//:one\"\ninfra = \"//:two\"\n",
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			input := "filegroup(name = \"artifact\")\n"
			writeFiles(t, root, map[string]string{
				"MODULE.bazel": "module(name = \"fixture\")\n",
				"BUILD.bazel":  input,
				"tools/transitive_visibility/transitive_visibility.toml": configuration,
			})
			runGazelle(t, root, *gazellePath, true)
			if got := readFile(t, root, "BUILD.bazel"); got != input {
				t.Fatalf("invalid configuration changed BUILD file:\n%s", got)
			}
		})
	}
}

// TestSelectedConfiguration verifies explicit selection and missing-file errors.
func TestSelectedConfiguration(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"MODULE.bazel":      "module(name = \"fixture\")\n",
		"BUILD.bazel":       "package(transitive_visibility = \"//:manual\")\n",
		"infra/BUILD.bazel": "filegroup(name = \"artifact\")\n",
		"policy.toml":       "[paths]\ninfra = \"//:selected\"\n",
	})
	runGazelle(t, root, *gazellePath, false, "-transitive_visibility_config=policy.toml")
	if got := readFile(t, root, "infra/BUILD.bazel"); !strings.Contains(got, "//:selected") {
		t.Fatalf("explicit configuration was not used:\n%s", got)
	}
	before := readFile(t, root, "BUILD.bazel")
	if !strings.Contains(before, "//:manual") {
		t.Fatalf("unmatched manual policy changed:\n%s", before)
	}
	runGazelle(t, root, *gazellePath, true, "-transitive_visibility_config=missing.toml")
	if got := readFile(t, root, "BUILD.bazel"); got != before {
		t.Fatalf("missing configuration changed BUILD file:\n%s", got)
	}
}

// TestObsoleteDirective requires policy to have a single configuration owner.
func TestObsoleteDirective(t *testing.T) {
	root := t.TempDir()
	input := "# gazelle:transitive_visibility //:old\nfilegroup(name = \"artifact\")\n"
	writeFiles(t, root, map[string]string{
		"MODULE.bazel": "module(name = \"fixture\")\n",
		"BUILD.bazel":  input,
	})
	runGazelle(t, root, *gazellePath, true)
	if got := readFile(t, root, "BUILD.bazel"); got != input {
		t.Fatalf("obsolete directive changed BUILD file:\n%s", got)
	}
}

// TestNestedWorkspace does not read configuration above the selected repository root.
func TestNestedWorkspace(t *testing.T) {
	outer := t.TempDir()
	writeFiles(t, outer, map[string]string{
		"tools/transitive_visibility/transitive_visibility.toml": "[paths]\n\".\" = \"//:outer\"\n",
		"nested/MODULE.bazel": "module(name = \"nested\")\n",
		"nested/BUILD.bazel":  "filegroup(name = \"artifact\")\n",
	})
	root := filepath.Join(outer, "nested")
	runGazelle(t, root, *gazellePath, false)
	if got := readFile(t, root, "BUILD.bazel"); strings.Contains(got, "transitive_visibility") {
		t.Fatalf("nested workspace inherited outer policy:\n%s", got)
	}
}

// runGazelle invokes the real executable in the isolated fixture workspace.
func runGazelle(t *testing.T, root, runfile string, wantFailure bool, options ...string) {
	t.Helper()
	binary, err := runfiles.Rlocation(runfile)
	if err != nil {
		t.Fatalf("locate Gazelle executable %q: %v", runfile, err)
	}
	args := append([]string{"-repo_root=" + root, "-index=false"}, options...)
	command := exec.Command(binary, append(args, root)...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if wantFailure {
		if err == nil || !strings.Contains(string(output), "transitive_visibility") {
			t.Fatalf("expected configuration failure, got error %v and output:\n%s", err, output)
		}
		return
	}
	if err != nil {
		t.Fatalf("generate fixture BUILD files: %v\n%s", err, output)
	}
}

// writeFiles initializes fixture inputs beneath the supplied test directory.
func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		filename := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
			t.Fatalf("create fixture directory for %q: %v", name, err)
		}
		if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
			t.Fatalf("write fixture %q: %v", name, err)
		}
	}
}

// readFile reads a generated fixture file for observable output assertions.
func readFile(t *testing.T, root, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("read fixture %q: %v", name, err)
	}
	return string(content)
}

// retainArtifacts preserves generated BUILD files in Bazel's test outputs.
func retainArtifacts(t *testing.T, name string, files map[string]string) {
	t.Helper()
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		return
	}
	root := filepath.Join(output, name)
	writeFiles(t, root, files)
	writeFiles(t, root, map[string]string{"result.txt": "Generated BUILD output verified by the Gazelle integration test.\n"})
}

// assertPackageOrder checks the observable ordering without depending on the plugin.
func assertPackageOrder(t *testing.T, name, content string) {
	t.Helper()
	position := strings.Index(content, "package(")
	if position < 0 {
		t.Fatalf("%s: package declaration missing", name)
	}
	for _, call := range []string{"load(", "filegroup(", "package_group(", "artifact(", "go_library("} {
		other := strings.Index(content, call)
		if other >= 0 && ((call == "load(" && other > position) || (call != "load(" && other < position)) {
			t.Errorf("%s: package declaration misplaced relative to %s:\n%s", name, call, content)
		}
	}
}
