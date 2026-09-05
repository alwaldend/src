package agent_system

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func TestRunEmitsHonestOfflineJSON(t *testing.T) {
	var stdout bytes.Buffer
	if err := Run([]string{
		"--workspace-root", "testdata/root",
		"--path", "projects/agents",
		"--revision", "abcdef0123456789",
		"--json",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	var capsule v1alpha1.ContextCapsule
	if err := json.Unmarshal(stdout.Bytes(), &capsule); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	decoded, err := v1alpha1.DecodeContextCapsule(stdout.Bytes())
	if err != nil {
		t.Fatalf("capsule does not validate: %v", err)
	}
	if decoded.Provenance.Completeness != v1alpha1.CompletenessPartial {
		t.Fatalf("completeness = %s, want partial offline observations", decoded.Provenance.Completeness)
	}
	if decoded.Identity.Revision != "abcdef0123456789" ||
		!strings.HasPrefix(decoded.Identity.SourceDigest, "sha256:") {
		t.Fatalf("caller revision and source digest conflated: %+v", decoded.Identity)
	}
	if decoded.Component.Workspace == "" {
		t.Fatalf("applicable workspace is empty")
	}
	if len(decoded.Provenance.NextActions) == 0 {
		t.Fatalf("expected safe next discovery actions")
	}
}

func TestRunEmitsStructuredUnavailableOnMissingCatalogs(t *testing.T) {
	var stdout bytes.Buffer
	if err := Run([]string{
		"--workspace-root", "testdata/missing",
		"--path", ".",
		"--revision", "abcdef0123456789",
		"--json",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	capsule, err := v1alpha1.DecodeContextCapsule(stdout.Bytes())
	if err != nil {
		t.Fatalf("capsule does not validate: %v", err)
	}
	if capsule.Provenance.Completeness != v1alpha1.CompletenessPartial {
		t.Fatalf("completeness = %s, want partial with structured unavailable", capsule.Provenance.Completeness)
	}
	if len(capsule.Provenance.Limitations) == 0 {
		t.Fatalf("expected structured unavailable limitations")
	}
}

func TestMarkdownRenderUsesSameData(t *testing.T) {
	var jsonOut, markdownOut bytes.Buffer
	now := func() time.Time { return time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC) }
	if err := runWithClock([]string{
		"--workspace-root", "testdata/root",
		"--path", "projects/agents",
		"--revision", "abcdef0123456789",
		"--json",
	}, &jsonOut, now); err != nil {
		t.Fatal(err)
	}
	if err := runWithClock([]string{
		"--workspace-root", "testdata/root",
		"--path", "projects/agents",
		"--revision", "abcdef0123456789",
		"--markdown",
		"--json=false",
	}, &markdownOut, now); err != nil {
		t.Fatal(err)
	}
	var capsule v1alpha1.ContextCapsule
	if err := json.Unmarshal(jsonOut.Bytes(), &capsule); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(markdownOut.String(), capsule.ID) {
		t.Fatalf("markdown does not state the same capsule ID")
	}
}

func TestCapsuleReportsReadTimeWithoutInventingAuthority(t *testing.T) {
	var output bytes.Buffer
	before := time.Now()
	if err := Run([]string{
		"--workspace-root", "testdata/root", "--task", "unrelated-task",
	}, &output); err != nil {
		t.Fatal(err)
	}
	after := time.Now()
	capsule, err := v1alpha1.DecodeContextCapsule(output.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	observed, err := time.Parse(time.RFC3339Nano, capsule.Provenance.ObservedAt)
	if err != nil || observed.Before(before) || observed.After(after) {
		t.Fatalf("read time %q is outside invocation: %v", capsule.Provenance.ObservedAt, err)
	}
	if capsule.Provenance.Freshness != "unknown" ||
		capsule.Provenance.AuthorizedBy != "" ||
		capsule.Outcome.Authority != "" || capsule.Outcome.GoalBinding != "" ||
		capsule.Component.ReviewOwners != "" {
		t.Fatalf("offline capsule invented an observation or authority: %+v", capsule)
	}
	for _, provider := range capsule.Providers {
		if provider.State != "unavailable" || provider.ObservedRevision != "" ||
			provider.ObservationTime != "" || provider.Unavailable == "" {
			t.Fatalf("offline provider invented a live observation: %+v", provider)
		}
	}
	if capsule.Identity.ByteSize != int64(output.Len()) {
		t.Fatalf("byteSize = %d, actual = %d", capsule.Identity.ByteSize, output.Len())
	}
}

func TestOutputLimitFailsBeforeWriting(t *testing.T) {
	for _, format := range []string{"--json", "--markdown"} {
		var output bytes.Buffer
		err := Run([]string{
			"--workspace-root", "testdata/root", "--max-bytes", "1024", format,
		}, &output)
		if err == nil || !strings.Contains(err.Error(), "exceeding --max-bytes") {
			t.Fatalf("%s error = %v, want output limit refusal", format, err)
		}
		if output.Len() != 0 {
			t.Fatalf("%s emitted partial output", format)
		}
	}
}

func TestContextPathAndFormatFlags(t *testing.T) {
	for _, format := range []string{"--markdown", "--json=false"} {
		opts, err := parseFlags([]string{
			"--workspace-root", "testdata/root", "--path", "//projects/agents:agents", format,
		})
		if err != nil || opts.path != "projects/agents" || !opts.markdown {
			t.Fatalf("label/format flags = %+v, error = %v", opts, err)
		}
	}
	for _, path := range []string{"../escape", "/absolute", "@external//package:target"} {
		if _, err := parseFlags([]string{"--workspace-root", "testdata/root", "--path", path}); err == nil {
			t.Fatalf("accepted out-of-workspace path %q", path)
		}
	}
}

func TestContextDiscoversApplicableDocumentsAndDeepestOwner(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"AGENTS.md":                            "root policy\n",
		"CODEOWNERS":                           "* @actual-owner\n",
		"projects/README.md":                   "tree boundary\n",
		"projects/sample/README.md":            "component purpose\n",
		"projects/sample/AGENTS.md":            "component policy\n",
		"projects/sample/src/AGENTS.md":        "nearest policy\n",
		"projects/sample/src/main.go":          "package sample\n",
		"MODULE.bazel":                         "module(name = \"test\")\n",
		"projects/sample/include.MODULE.bazel": "# component dependencies\n",
		"projects/sample/src/BUILD.bazel":      "# narrow package targets\n",
	}
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot := &catalogSnapshot{
		topology: &catalogv1alpha1.TopologyCatalog{
			Components: []catalogv1alpha1.TopologyComponent{
				{ID: "parent", Path: "projects", Lifecycle: "active"},
				{ID: "sample", Path: "projects/sample", Lifecycle: "experimental"},
			},
		},
		policy: &catalogv1alpha1.PolicyCatalog{
			Policies: []catalogv1alpha1.PolicyRecord{
				{ID: "tree", PathPrefix: "/projects", OwnerBoundaryRef: "projects/README.md"},
			},
		},
	}
	build := func() v1alpha1.ContextCapsule {
		t.Helper()
		builder := &capsuleBuilder{root: root, opts: options{path: "projects/sample/src/main.go", repository: "test/repo", dirty: true}}
		capsule, err := builder.build(snapshot, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		return capsule
	}
	capsule := build()
	if capsule.Component.ComponentID != "sample" || capsule.Component.OwnerReadme != "projects/sample/README.md" {
		t.Fatalf("incorrect closest owner: %+v", capsule.Component)
	}
	want := []string{
		"AGENTS.md", "projects/sample/AGENTS.md", "projects/sample/src/AGENTS.md",
		"projects/sample/README.md", "projects/sample/src/BUILD.bazel",
		"projects/sample/include.MODULE.bazel", "MODULE.bazel", "projects/README.md", "CODEOWNERS",
	}
	if len(capsule.Documents) != len(want) {
		t.Fatalf("documents = %+v, want %v", capsule.Documents, want)
	}
	for index, path := range want {
		if capsule.Documents[index].Path != path || capsule.Documents[index].Digest == "" {
			t.Fatalf("document %d = %+v, want %s", index, capsule.Documents[index], path)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "projects/sample/src/AGENTS.md"), []byte("changed policy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed := build(); changed.Identity.InputDigest == capsule.Identity.InputDigest {
		t.Fatal("applicable policy edit did not change the capsule input digest")
	}
	for _, action := range capsule.Provenance.NextActions {
		if strings.Contains(action, "//...") {
			t.Fatalf("missing workspace catalog recommended a broad check: %s", action)
		}
	}
}

func TestNestedWorkspaceKeepsLocalAndRootCandidates(t *testing.T) {
	snapshot := &catalogSnapshot{
		workspaceCheck: &catalogv1alpha1.WorkspaceCheckCatalog{
			Workspaces: []catalogv1alpha1.WorkspaceRecord{
				{ID: "root", Path: "."},
				{ID: "sample", Path: "projects/sample"},
				{ID: "other", Path: "projects/sample_extra"},
			},
		},
		capability: &catalogv1alpha1.CapabilityCatalog{
			Skills: []catalogv1alpha1.CapabilitySkill{
				{ID: "global", Owner: "projects/agents", CapabilityRefs: []string{"an.action.reference"}},
				{ID: "local", Owner: "projects/sample"},
				{ID: "other", Owner: "projects/sample_extra"},
			},
		},
	}
	builder := &capsuleBuilder{opts: options{path: "projects/sample/src/main.go"}}
	workspace := builder.applicableWorkspace(snapshot)
	if workspace.ID != "sample" {
		t.Fatalf("workspace = %s", workspace.ID)
	}
	capabilities := builder.capabilities(snapshot, workspace)
	if len(capabilities) != 2 || capabilities[0].ID != "global" || capabilities[1].ID != "local" {
		t.Fatalf("nested candidates = %+v", capabilities)
	}
	if len(capabilities[0].Effects) != 0 {
		t.Fatal("capability references were mislabeled as atomic effects")
	}
}

func TestContextFollowsSkillDiscoverySymlink(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "workspace")
	files := map[string]string{
		"AGENTS.md":                 "repository policy\n",
		"README.md":                 "repository overview\n",
		"projects/agents/README.md": "agent component\n",
		"projects/agents/AGENTS.md": "agent component policy\n",
		"projects/agents/skills/repo-bazel/SKILL.md": "skill content\n",
	}
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	discovery := filepath.Join(root, ".agents/skills")
	if err := os.MkdirAll(discovery, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../../projects/agents/skills/repo-bazel", filepath.Join(discovery, "repo-bazel")); err != nil {
		t.Fatal(err)
	}
	workspaceLink := filepath.Join(parent, "workspace-link")
	if err := os.Symlink(root, workspaceLink); err != nil {
		t.Fatal(err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, suffix := range []string{"", "/SKILL.md", "/new-reference.md"} {
		var output bytes.Buffer
		if err := Run([]string{
			"--workspace-root", workspaceLink,
			"--path", ".agents/skills/repo-bazel" + suffix,
		}, &output); err != nil {
			t.Fatal(err)
		}
		capsule, err := v1alpha1.DecodeContextCapsule(output.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if capsule.Component.Path != "projects/agents/skills/repo-bazel"+suffix ||
			capsule.Identity.WorkspaceRoot != canonicalRoot ||
			capsule.Component.OwnerReadme != "projects/agents/README.md" {
			t.Fatalf("discovery link used lexical ownership: %+v", capsule)
		}
		foundPolicy := false
		for _, document := range capsule.Documents {
			if document.Path == "projects/agents/AGENTS.md" {
				foundPolicy = true
			}
		}
		if !foundPolicy {
			t.Fatal("resolved skill omitted its owning component policy")
		}
	}
}

func TestContextRejectsEscapingAndUnresolvedSymlinks(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "workspace")
	outside := filepath.Join(parent, "outside")
	for _, path := range []string{root, outside} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	links := map[string]string{
		"escape":   outside,
		"dangling": "missing-target",
		"cycle":    "cycle",
	}
	for name, target := range links {
		if err := os.Symlink(target, filepath.Join(root, name)); err != nil {
			t.Fatal(err)
		}
		for _, suffix := range []string{"", "/new-file.md"} {
			var output bytes.Buffer
			err := Run([]string{
				"--workspace-root", root, "--path", name + suffix,
			}, &output)
			if err == nil || output.Len() != 0 {
				t.Fatalf("%s%s error = %v, output bytes = %d", name, suffix, err, output.Len())
			}
		}
	}
}

func TestContextReportsUnavailableApplicablePolicy(t *testing.T) {
	for _, condition := range []string{"absent", "dangling", "escaping", "cyclic", "directory", "unreadable"} {
		t.Run(condition, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("root policy\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(root, "projects/sample"), 0o755); err != nil {
				t.Fatal(err)
			}
			policyPath := "projects/sample/AGENTS.md"
			full := filepath.Join(root, filepath.FromSlash(policyPath))
			var err error
			switch condition {
			case "dangling":
				err = os.Symlink("missing-policy", full)
			case "escaping":
				outside := filepath.Join(t.TempDir(), "policy.md")
				if err := os.WriteFile(outside, []byte("outside policy\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				err = os.Symlink(outside, full)
			case "cyclic":
				err = os.Symlink("AGENTS.md", full)
			case "directory":
				err = os.Mkdir(full, 0o755)
			case "unreadable":
				err = os.WriteFile(full, []byte("scoped policy\n"), 0o000)
				if err == nil {
					if _, readErr := os.ReadFile(full); readErr == nil {
						t.Skip("process can read mode-000 files")
					}
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			builder := &capsuleBuilder{root: root, opts: options{path: "projects/sample", repository: "test/repo", dirty: true}}
			capsule, err := builder.build(&catalogSnapshot{}, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			count := 0
			for _, limitation := range capsule.Provenance.Limitations {
				if limitation == "document unavailable: "+policyPath {
					count++
				}
				if limitation == "document unavailable: projects/AGENTS.md" {
					t.Fatal("confirmed absent intermediate policy was marked unavailable")
				}
			}
			want := 1
			if condition == "absent" {
				want = 0
			}
			if count != want {
				t.Fatalf("unavailable policy count = %d, want %d; limitations %v", count, want, capsule.Provenance.Limitations)
			}
			for _, document := range capsule.Documents {
				if document.Path == policyPath {
					t.Fatalf("unavailable policy was accepted: %+v", document)
				}
			}
		})
	}
}

func TestContextBindsDistinctApplicablePolicyAxisSources(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "policy"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		"AGENTS.md": "root policy\n", "policy/consumers.md": "allowed consumers\n",
	} {
		if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	build := func() v1alpha1.ContextCapsule {
		t.Helper()
		builder := &capsuleBuilder{root: root, opts: options{path: "projects/sample", repository: "test/repo", dirty: true}}
		capsule, err := builder.build(&catalogSnapshot{
			policy: &catalogv1alpha1.PolicyCatalog{Policies: []catalogv1alpha1.PolicyRecord{
				{ID: "applicable", PathPrefix: "/projects/sample", AgentPolicySource: "AGENTS.md", Axes: []catalogv1alpha1.PolicyAxis{
					{Name: "consumers", Source: "policy/consumers.md"},
					{Name: "publication", Source: "policy/consumers.md"},
				}},
				{ID: "unrelated", PathPrefix: "/projects/other", Axes: []catalogv1alpha1.PolicyAxis{
					{Name: "consumers", Source: "policy/unrelated-missing.md"},
				}},
			}},
		}, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		return capsule
	}
	first := build()
	count := 0
	for _, document := range first.Documents {
		if document.Path == "policy/consumers.md" && document.Digest != "" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("axis authority should be bound once: %+v", first.Documents)
	}
	for _, limitation := range first.Provenance.Limitations {
		if strings.Contains(limitation, "unrelated-missing") {
			t.Fatalf("unrelated policy axis affected capsule: %s", limitation)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "policy/consumers.md"), []byte("changed consumers\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if second := build(); second.Identity.InputDigest == first.Identity.InputDigest {
		t.Fatal("distinct axis source edit did not change capsule input digest")
	}
}

func TestContextEscalatesOptionalMissingDocumentWhenRequiredByAxis(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "projects/sample"), 0o755); err != nil {
		t.Fatal(err)
	}
	builder := &capsuleBuilder{root: root, opts: options{path: "projects/sample", repository: "test/repo", dirty: true}}
	capsule, err := builder.build(&catalogSnapshot{
		policy: &catalogv1alpha1.PolicyCatalog{Policies: []catalogv1alpha1.PolicyRecord{
			{ID: "applicable", PathPrefix: "/projects/sample", Axes: []catalogv1alpha1.PolicyAxis{
				{Name: "consumers", Source: "projects/sample/README.md"},
				{Name: "publication", Source: "projects/sample/README.md"},
				{Name: "instructions", Source: "AGENTS.md"},
			}},
		}},
	}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"projects/sample/README.md", "AGENTS.md"} {
		count := 0
		for _, limitation := range capsule.Provenance.Limitations {
			if limitation == "document unavailable: "+path {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("required %s unavailable count = %d, want 1; limitations %v", path, count, capsule.Provenance.Limitations)
		}
	}
}

func TestContextPreservesNearestUnavailableDiscoverySource(t *testing.T) {
	for _, name := range []string{"README.md", "BUILD.bazel", "include.MODULE.bazel", "MODULE.bazel", "CODEOWNERS"} {
		for _, condition := range []string{"absent", "dangling", "escaping", "cyclic", "unreadable"} {
			t.Run(name+"/"+condition, func(t *testing.T) {
				root := t.TempDir()
				if err := os.MkdirAll(filepath.Join(root, "projects/sample"), 0o755); err != nil {
					t.Fatal(err)
				}
				for path, content := range map[string]string{
					"AGENTS.md": "root policy\n", "README.md": "root boundary\n",
				} {
					if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o644); err != nil {
						t.Fatal(err)
					}
				}
				nearest := "projects/sample/" + name
				fallback := name
				if name == "CODEOWNERS" {
					nearest, fallback = name, ""
				} else {
					if err := os.WriteFile(filepath.Join(root, name), []byte("ancestor source\n"), 0o644); err != nil {
						t.Fatal(err)
					}
				}
				if name == "BUILD.bazel" {
					fallback = "projects/sample/BUILD"
					if err := os.WriteFile(filepath.Join(root, fallback), []byte("alternate build source\n"), 0o644); err != nil {
						t.Fatal(err)
					}
				}
				full := filepath.Join(root, filepath.FromSlash(nearest))
				var err error
				switch condition {
				case "dangling":
					err = os.Symlink("missing-source", full)
				case "escaping":
					outside := filepath.Join(t.TempDir(), "source")
					if err := os.WriteFile(outside, []byte("outside source\n"), 0o644); err != nil {
						t.Fatal(err)
					}
					err = os.Symlink(outside, full)
				case "cyclic":
					err = os.Symlink(name, full)
				case "unreadable":
					err = os.WriteFile(full, []byte("local source\n"), 0o000)
					if err == nil {
						if _, readErr := os.ReadFile(full); readErr == nil {
							t.Skip("process can read mode-000 files")
						}
					}
				}
				if err != nil {
					t.Fatal(err)
				}
				builder := &capsuleBuilder{root: root, opts: options{path: "projects/sample", repository: "test/repo", dirty: true}}
				capsule, err := builder.build(&catalogSnapshot{}, time.Now())
				if err != nil {
					t.Fatal(err)
				}
				if name == "README.md" {
					want := nearest
					if condition == "absent" {
						want = fallback
					}
					if capsule.Component.OwnerReadme != want {
						t.Fatalf("owner README = %s, want %s", capsule.Component.OwnerReadme, want)
					}
				}
				foundFallback := false
				for _, document := range capsule.Documents {
					if document.Path == nearest {
						t.Fatalf("unavailable source was accepted: %+v", document)
					}
					if document.Path == fallback {
						foundFallback = true
					}
					if condition != "absent" && name == "BUILD.bazel" && document.Path == "BUILD.bazel" {
						t.Fatal("unavailable local BUILD.bazel was replaced with ancestor BUILD.bazel")
					}
				}
				if fallback != "" && foundFallback != (condition == "absent") {
					t.Fatalf("fallback included = %t under %s", foundFallback, condition)
				}
				unavailable := 0
				for _, limitation := range capsule.Provenance.Limitations {
					if limitation == "document unavailable: "+nearest {
						unavailable++
					}
				}
				want := 1
				if condition == "absent" {
					want = 0
				}
				if unavailable != want {
					t.Fatalf("unavailable count = %d, want %d", unavailable, want)
				}
			})
		}
	}
}
