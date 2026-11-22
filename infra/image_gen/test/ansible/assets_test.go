package ansible

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// ansiblePath identifies the actual packaged deployment without its live AL wrapper.
var ansiblePath = flag.String("ansible", "", "Packaged Ansible launcher runfile")

// TestProfileWorkflows installs each real profile's native workflows with local-only inventory.
func TestProfileWorkflows(t *testing.T) {
	t.Setenv("TMPDIR", os.Getenv("TEST_TMPDIR"))
	launcher, err := runfiles.Rlocation(*ansiblePath)
	if err != nil {
		t.Fatal(err)
	}
	files, err := runfiles.Rlocation("_main/files/configuration.json")
	if err != nil {
		t.Fatal(err)
	}
	packagedRoot := filepath.Dir(filepath.Dir(files))
	configuration, err := os.ReadFile(files)
	if err != nil {
		t.Fatal(err)
	}
	var registry map[string]any
	if err := json.Unmarshal(configuration, &registry); err != nil {
		t.Fatal(err)
	}
	profiles := registry["model_profiles"].(map[string]any)
	installed := map[string]int{}
	for name, raw := range profiles {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			workflows := filepath.Join(root, "ComfyUI", "user", "default", "workflows")
			if err := os.MkdirAll(workflows, 0o700); err != nil {
				t.Fatal(err)
			}
			var selected map[string]any
			if err := json.Unmarshal(configuration, &selected); err != nil {
				t.Fatal(err)
			}
			selected["model_profile"] = name
			extra := map[string]any{"imagegen_configuration": selected, "comfyui_tools_directory": filepath.Join(root, "ComfyUI"), "ansible_python_interpreter": "{{ ansible_playbook_python }}"}
			arguments := []string{"--tags", "workflows"}
			playbook := filepath.Join(packagedRoot, "playbook.yaml")
			if name == "noobai" {
				// Exercise the actual playbook JSON translation with local substitutes.
				if err := os.WriteFile(filepath.Join(root, "ComfyUI", "main.py"), []byte("# fixture application\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					switch r.URL.Path {
					case "/prompt":
						fmt.Fprint(w, `{"exec_info":{"queue_remaining":0}}`)
					case "/object_info":
						classes := map[string]any{"ImagegenDanbooruTagSearch": map[string]any{}}
						for _, rawNode := range raw.(map[string]any)["first_party_nodes"].([]any) {
							for _, class := range selected["first_party_nodes"].(map[string]any)[rawNode.(string)].(map[string]any)["node_classes"].([]any) {
								classes[class.(string)] = map[string]any{}
							}
						}
						if err := json.NewEncoder(w).Encode(classes); err != nil {
							t.Error(err)
						}
					default:
						fmt.Fprint(w, "public profile model fixture")
					}
				}))
				defer server.Close()
				checksum := fmt.Sprintf("%x", sha256.Sum256([]byte("public profile model fixture")))
				for _, filename := range raw.(map[string]any)["downloads"].([]any) {
					asset := selected["downloads"].(map[string]any)[filename.(string)].(map[string]any)
					asset["url"], asset["sha256"] = server.URL+"/model", checksum
				}
				repository := filepath.Join(root, "node-fixture")
				if err := os.MkdirAll(repository, 0o700); err != nil {
					t.Fatal(err)
				}
				git := func(args ...string) string {
					gitContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()
					command := exec.CommandContext(gitContext, "git", args...)
					command.Dir = repository
					command.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
					output, err := command.CombinedOutput()
					if err != nil {
						t.Fatalf("fixture git: %v: %s", err, shortDiagnostics(output))
					}
					return strings.TrimSpace(string(output))
				}
				git("init", "--quiet")
				if err := os.WriteFile(filepath.Join(repository, "__init__.py"), []byte("NODE_CLASS_MAPPINGS = {}\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				git("add", "__init__.py")
				git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "--quiet", "-m", "Fixture node")
				revision := git("rev-parse", "HEAD")
				for _, node := range selected["custom_nodes"].(map[string]any) {
					entry := node.(map[string]any)
					entry["url"], entry["ref"] = (&url.URL{Scheme: "file", Path: repository}).String(), revision
				}
				touch, err := exec.LookPath("touch")
				if err != nil {
					t.Fatal(err)
				}
				extra["comfyui_tools_api_url"] = server.URL
				extra["comfyui_tools_state_directory"] = filepath.Join(root, "state")
				extra["comfyui_tools_environment_file"] = filepath.Join(root, "environment")
				extra["comfyui_tools_restart_command"] = []string{touch, filepath.Join(root, "restarted")}
				extra["comfyui_tools_python_requirements"] = ""
				extra["comfyui_tools_interpolation_source"] = ""
				arguments = nil
			}
			variables, err := json.Marshal(extra)
			if err != nil {
				t.Fatal(err)
			}
			variablePath := filepath.Join(root, "variables.json")
			if err := os.WriteFile(variablePath, variables, 0o600); err != nil {
				t.Fatal(err)
			}
			inventory := filepath.Join(root, "inventory.yaml")
			if err := os.WriteFile(inventory, []byte("imagegen:\n  hosts:\n    fixture:\n      ansible_connection: local\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			arguments = append([]string{"-i", inventory, "--limit", "fixture"}, arguments...)
			arguments = append(arguments, "-e", "@"+variablePath, playbook)
			command := exec.CommandContext(ctx, launcher, arguments...)
			command.Dir = packagedRoot
			command.Env = append(os.Environ(), "ANSIBLE_CONFIG="+filepath.Join(packagedRoot, "ansible.cfg"), "ANSIBLE_COLLECTIONS_PATH="+filepath.Join(packagedRoot, "collections"), "ANSIBLE_LOCAL_TEMP="+filepath.Join(root, "controller"), "ANSIBLE_REMOTE_TEMP="+filepath.Join(root, "remote"), "TMPDIR="+root, "ANSIBLE_NOCOLOR=1")
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("local profile installation: %v\n%s", err, shortDiagnostics(output))
			}
			if name == "noobai" {
				if _, err := os.Stat(filepath.Join(root, "restarted")); err != nil {
					t.Fatal("full installation did not restart the fixture")
				}
				for _, filename := range []string{"noobai-text-to-image.app.json", "danbooru-tag-search.app.json", "danbooru-tag-search.json"} {
					if _, err := os.Stat(filepath.Join(workflows, filename)); err != nil {
						t.Fatalf("required NoobAI workflow %s was not installed: %v", filename, err)
					}
				}
				for _, filename := range []string{"__init__.py", "tags.csv"} {
					path := filepath.Join(root, "ComfyUI", "custom_nodes", "danbooru_tag_search", filename)
					info, err := os.Stat(path)
					if err != nil || info.Size() == 0 {
						t.Fatalf("tag search asset %s was not installed: %v", filename, err)
					}
				}
				for _, node := range raw.(map[string]any)["first_party_nodes"].([]any) {
					sourceName := selected["first_party_nodes"].(map[string]any)[node.(string)].(map[string]any)["source"].(string)
					expected, err := os.ReadFile(filepath.Join(packagedRoot, "files", sourceName))
					if err != nil {
						t.Fatal(err)
					}
					actual, err := os.ReadFile(filepath.Join(root, "ComfyUI", "custom_nodes", node.(string), "__init__.py"))
					if err != nil || !bytes.Equal(expected, actual) {
						t.Fatalf("selected first-party node %s differs: %v", node, err)
					}
				}
				expected, err := os.ReadFile(filepath.Join(root, "ComfyUI", "models", "checkpoints", "NoobAI-XL-v1.1.safetensors"))
				if err != nil || string(expected) != "public profile model fixture" {
					t.Fatalf("selected checkpoint was not installed: %v", err)
				}
			}
			for _, rawName := range raw.(map[string]any)["workflows"].([]any) {
				filename := rawName.(string) + ".json"
				expected, err := os.ReadFile(filepath.Join(packagedRoot, "files", "workflows", name, filename))
				if err != nil {
					t.Fatal(err)
				}
				actual, err := os.ReadFile(filepath.Join(workflows, filename))
				if err != nil || !bytes.Equal(actual, expected) {
					t.Fatalf("workflow %s differs: %v", filename, err)
				}
				installed[name]++
			}
		})
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		evidence, err := json.Marshal(map[string]any{"installed_workflows": installed, "all_contents_preserved": true, "full_profile_with_local_substitutes": "noobai", "json_translation": "playbook", "live_operations": 0})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "profile-workflows.json"), evidence, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// shortDiagnostics bounds fixture failures and omits embedded Ansible module archives.
func shortDiagnostics(output []byte) string {
	var result strings.Builder
	for _, line := range strings.Split(string(output), "\n") {
		if len(line) > 400 {
			continue
		}
		if result.Len()+len(line) > 12000 {
			break
		}
		result.WriteString(line + "\n")
	}
	return result.String()
}
