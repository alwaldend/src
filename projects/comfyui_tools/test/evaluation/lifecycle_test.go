package evaluation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// fixture owns only simulated evaluation contracts and their ordered events.
type fixture struct {
	// mu serializes HTTP requests from independent evaluator workers.
	mu sync.Mutex
	// labels maps simulated contract IDs to exact ownership labels.
	labels map[int]string
	// events records non-secret lifecycle observations for the E2E receipt.
	events []string
	// active tracks the maximum simultaneous simulated rentals.
	active int
	// maximum records the observed concurrency peak.
	maximum int
	// rejectDelete simulates a provider-side cleanup failure.
	rejectDelete bool
	// oneOffer limits the fixture to one reusable offer.
	oneOffer bool
	// queries preserves search input for shared-filter acceptance.
	queries []map[string]any
	// createdOffers records which eligible capacity was actually rented.
	createdOffers []int64
	// exitAfterReady stops the simulated provider after its first running observation.
	exitAfterReady bool
	// statusReads counts running observations for the single-rental failure case.
	statusReads int
	// ambiguousCreate records an allocation while losing its acknowledgment.
	ambiguousCreate bool
	// providerStatus optionally overrides the simulated running state.
	providerStatus string
}

// ServeHTTP supplies the supported Vast protocol without contacting an account.
func (f *fixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.URL.Path == "/bundles/":
		var query map[string]any
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			http.Error(w, "invalid query", 400)
			return
		}
		f.queries = append(f.queries, query)
		offers := []map[string]any{{"id": 1, "hosting_type": 1, "verification": "verified", "gpu_name": "fixture", "num_gpus": 1, "gpu_ram": 49152, "dph_total": 0.5}}
		if !f.oneOffer {
			offers = append(offers, map[string]any{"id": 2, "hosting_type": 1, "verification": "verified", "gpu_name": "fixture", "num_gpus": 1, "gpu_ram": 49152, "dph_total": 0.5})
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"offers": offers}); err != nil {
			return
		}
	case strings.HasPrefix(r.URL.Path, "/asks/"):
		offer, err := strconv.ParseInt(strings.Trim(strings.TrimPrefix(r.URL.Path, "/asks/"), "/"), 10, 64)
		if err != nil {
			http.Error(w, "invalid fixture offer", 400)
			return
		}
		f.createdOffers = append(f.createdOffers, offer)
		var body struct {
			// Label is the only instance creation value retained by the fixture.
			Label string `json:"label"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid fixture creation", 400)
			return
		}
		id := len(f.events) + 100
		f.labels[id] = body.Label
		f.events = append(f.events, "create:"+body.Label)
		f.active++
		f.maximum = max(f.maximum, f.active)
		if f.ambiguousCreate {
			http.Error(w, "fixture acknowledgment lost", 502)
			return
		}
		fmt.Fprintf(w, `{"success":true,"new_contract":%d}`, id)
	case r.URL.Path == "/instances/":
		http.Error(w, "deprecated instance listing", 410)
	case r.URL.Path == "/api/v1/instances/":
		var rows []map[string]any
		for id, label := range f.labels {
			rows = append(rows, map[string]any{"id": id, "label": label})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i]["id"].(int) < rows[j]["id"].(int) })
		if r.URL.Query().Get("after_token") == "" && len(rows) > 1 {
			if err := json.NewEncoder(w).Encode(map[string]any{"instances": rows[:1], "next_token": "fixture-next"}); err != nil {
				return
			}
			return
		}
		if r.URL.Query().Get("after_token") != "" {
			rows = rows[1:]
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"instances": rows, "next_token": nil}); err != nil {
			return
		}
	case strings.HasPrefix(r.URL.Path, "/instances/"):
		id, err := strconv.Atoi(strings.Trim(r.URL.Path[len("/instances/"):], "/"))
		if err != nil {
			http.Error(w, "invalid fixture ID", 400)
			return
		}
		label, found := f.labels[id]
		if !found {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			if f.rejectDelete {
				http.Error(w, "fixture deletion unavailable", 503)
				return
			}
			f.events = append(f.events, "delete:"+label)
			delete(f.labels, id)
			f.active--
			fmt.Fprint(w, `{"success":true}`)
			return
		}
		status := "running"
		if f.providerStatus != "" {
			status = f.providerStatus
		}
		if f.exitAfterReady {
			f.statusReads++
			if f.statusReads > 1 {
				status = "exited"
			}
		}
		fmt.Fprintf(w, `{"instances":{"id":%d,"label":%q,"actual_status":%q,"public_ipaddr":"127.0.0.1","ports":{"7860/tcp":[{"HostPort":"443"},{"HostPort":"443"}]}}}`, id, label, status)
	default:
		http.NotFound(w, r)
	}
}

// TestRunOfferExclusion skips known failed capacity without weakening shared eligibility.
func TestRunOfferExclusion(t *testing.T) {
	for _, id := range []string{"1", "0"} {
		t.Run(id, func(t *testing.T) {
			t.Setenv("XDG_CACHE_HOME", t.TempDir())
			f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}}
			server := httptest.NewServer(f)
			defer server.Close()
			command := fixtureCommand(t, server.URL, filepath.Join(t.TempDir(), "results"), []map[string]any{fixtureGroup("exclusion")})
			command.Args = append(command.Args, "--exclude-offer", id)
			output, err := command.CombinedOutput()
			f.mu.Lock()
			defer f.mu.Unlock()
			if id == "0" {
				if err == nil || len(f.createdOffers) != 0 {
					t.Fatal("invalid excluded identity did not fail before allocation")
				}
				return
			}
			if err != nil {
				t.Fatalf("run with excluded capacity: %v\n%s", err, output)
			}
			if len(f.createdOffers) != 1 || f.createdOffers[0] != 2 || len(f.labels) != 1 {
				t.Fatal("excluded cheapest offer was rented or test capacity was left behind")
			}
		})
	}
}

// TestCleanupNamespace verifies all matching ages/runs and failed-delete recovery.
func TestCleanupNamespace(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f := &fixture{labels: map[int]string{99: "src_infra_image_gen", 100: "imagegen-eval:fixture:old:job:", 101: "imagegen-eval:fixture:active:job:", 102: "imagegen-eval:another_owner:run:job:"}, active: 2, rejectDelete: true}
	server := httptest.NewServer(f)
	defer server.Close()
	binary, err := runfiles.Rlocation(os.Getenv("EVALUATOR_BINARY"))
	if err != nil {
		t.Fatalf("resolve cleanup executable: %v", err)
	}
	config := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(config, []byte(`{"owner":"fixture","max_parallel":2,"max_lifetime_seconds":30,"boot_timeout_seconds":5,"workload_timeout_seconds":5,"cleanup_timeout_seconds":0.1,"poll_seconds":0.02,"matrix":[]}`), 0o600); err != nil {
		t.Fatalf("write cleanup matrix: %v", err)
	}
	run := func() error {
		command := exec.Command(binary, "cleanup", "--config", config, "--endpoint", server.URL)
		command.Env = append(os.Environ(), "VAST_API_KEY=fixture-key")
		if err := command.Run(); err != nil {
			return fmt.Errorf("execute cleanup fixture: %w", err)
		}
		return nil
	}
	if err := run(); err == nil {
		t.Fatal("failed deletion was reported as successful cleanup")
	}
	f.mu.Lock()
	f.rejectDelete = false
	f.mu.Unlock()
	if err := run(); err != nil {
		t.Fatalf("recover failed cleanup: %v", err)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.labels) != 2 || f.labels[99] != "src_infra_image_gen" || f.labels[102] == "" {
		t.Fatal("cleanup did not remove all matching tests while preserving other namespaces")
	}
}

// TestParallelLifecycle runs the actual command with staggered workload completion.
func TestParallelLifecycle(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}}
	server := httptest.NewServer(f)
	defer server.Close()
	binary, err := runfiles.Rlocation(os.Getenv("EVALUATOR_BINARY"))
	if err != nil {
		t.Fatalf("resolve evaluator: %v", err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve fixture worker: %v", err)
	}
	directory := t.TempDir()
	profiles := filepath.Join(directory, "profiles.json")
	if err := os.WriteFile(profiles, []byte(`{"image":"fixture","hardware_profiles":{"test":{"min_gpu_ram_gb":48,"disk_gb":100,"max_hourly_price":1},"other":{"min_gpu_ram_gb":48,"disk_gb":100,"max_hourly_price":1}},"model_profiles":{"test":{"downloads":[],"workflows":[]}},"downloads":{},"custom_nodes":{}}`), 0o600); err != nil {
		t.Fatalf("write profiles: %v", err)
	}
	configuration := map[string]any{"owner": "fixture", "max_parallel": 2, "max_lifetime_seconds": 30, "boot_timeout_seconds": 5, "workload_timeout_seconds": 5, "cleanup_timeout_seconds": 5, "poll_seconds": 0.02, "matrix": []map[string]any{fixtureGroup("fast"), fixtureGroup("slow")}}
	configuration["matrix"].([]map[string]any)[1]["hardware_profiles"] = []string{"other"}
	data, err := json.Marshal(configuration)
	if err != nil {
		t.Fatalf("encode configuration: %v", err)
	}
	config := filepath.Join(directory, "config.json")
	if err := os.WriteFile(config, data, 0o600); err != nil {
		t.Fatalf("write configuration: %v", err)
	}
	startup := filepath.Join(directory, "startup.sh")
	if err := os.WriteFile(startup, []byte("set -eu\ntrue\n"), 0o600); err != nil {
		t.Fatalf("write fixture startup: %v", err)
	}
	cmd := exec.Command(binary, "run", "--config", config, "--profiles", profiles, "--endpoint", server.URL, "--output", filepath.Join(directory, "results"), "--worker", self, "--worker-arg=-test.run=^TestWorkloadProcess$", "--startup-file", startup, "--offer-defaults", fixtureOffers(t, directory))
	cmd.Env = append(os.Environ(), "VAST_API_KEY=fixture-key", "EVALUATION_WORKER_FIXTURE=1", "GORACE=atexit_sleep_ms=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run evaluator: %v\n%s", err, output)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.labels) != 1 || f.labels[99] != "src_infra_image_gen" {
		t.Fatal("evaluation cleanup changed a personal rental or left evaluation contracts")
	}
	if f.maximum != 2 {
		t.Fatalf("expected independent parallel rentals, observed %d", f.maximum)
	}
	var fastDelete, slowDelete int = -1, -1
	for index, event := range f.events {
		if strings.HasPrefix(event, "delete:") && strings.Contains(event, ":fast_001:") {
			fastDelete = index
		}
		if strings.HasPrefix(event, "delete:") && strings.Contains(event, ":slow_001:") {
			slowDelete = index
		}
	}
	if fastDelete < 0 || slowDelete <= fastDelete {
		t.Fatal("fast rental waited for the slow worker before deletion")
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		receipt, err := json.MarshalIndent(map[string]any{"events": f.events, "maximum_parallel": f.maximum, "personal_instance_preserved": true, "evaluation_contracts_remaining": 0, "live_operations": 0}, "", "  ")
		if err != nil {
			t.Fatalf("encode lifecycle artifact: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "evaluation-lifecycle.json"), receipt, 0o600); err != nil {
			t.Fatalf("write lifecycle artifact: %v", err)
		}
	}
}

// TestWorkloadProcess supplies a controlled subprocess at the workload boundary.
func TestWorkloadProcess(t *testing.T) {
	if os.Getenv("EVALUATION_WORKER_FIXTURE") != "1" {
		return
	}
	if os.Getenv("VAST_API_KEY") != "" {
		t.Fatal("rental account key reached a workload process")
	}
	if strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "failure" {
		t.Fatal("configured workload failed")
	}
	if strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "missing" {
		return
	}
	if strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "provider_exit" || strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "hold" || strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "deadline" || strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "cancel" {
		time.Sleep(5 * time.Second)
	}
	if strings.TrimSuffix(os.Getenv("EVALUATION_JOB_ID"), "_001") == "slow" {
		time.Sleep(300 * time.Millisecond)
	}
	directory := os.Getenv("EVALUATION_OUTPUT")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatalf("create fixture output: %v", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "result.json"), []byte(`{"download_complete":true,"review_scores":null}`), 0o600); err != nil {
		t.Fatalf("write fixture result: %v", err)
	}
}

// TestFailureCleanup exercises real command failure, deadline, cancellation and lost acknowledgment.
func TestFailureCleanup(t *testing.T) {
	for _, scenario := range []string{"failure", "missing", "deadline", "cancel", "ambiguous"} {
		t.Run(scenario, func(t *testing.T) {
			t.Setenv("XDG_CACHE_HOME", t.TempDir())
			f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}, ambiguousCreate: scenario == "ambiguous"}
			server := httptest.NewServer(f)
			defer server.Close()
			binary, err := runfiles.Rlocation(os.Getenv("EVALUATOR_BINARY"))
			if err != nil {
				t.Fatalf("resolve evaluator: %v", err)
			}
			self, err := os.Executable()
			if err != nil {
				t.Fatalf("resolve fixture process: %v", err)
			}
			directory := t.TempDir()
			profiles := filepath.Join(directory, "profiles.json")
			if err := os.WriteFile(profiles, []byte(`{"image":"fixture","hardware_profiles":{"test":{"min_gpu_ram_gb":48,"disk_gb":100,"max_hourly_price":1}},"model_profiles":{"test":{"downloads":[]}},"downloads":{},"custom_nodes":{}}`), 0o600); err != nil {
				t.Fatalf("write profile: %v", err)
			}
			workloadTimeout := 0.1
			if scenario == "cancel" {
				workloadTimeout = 5
			}
			configuration := map[string]any{"owner": "fixture", "max_parallel": 1, "max_lifetime_seconds": 10, "boot_timeout_seconds": 2, "workload_timeout_seconds": workloadTimeout, "cleanup_timeout_seconds": 2, "poll_seconds": 0.02, "matrix": []map[string]any{fixtureGroup(scenario)}}
			data, err := json.Marshal(configuration)
			if err != nil {
				t.Fatalf("encode failure configuration: %v", err)
			}
			config := filepath.Join(directory, "config.json")
			if err := os.WriteFile(config, data, 0o600); err != nil {
				t.Fatalf("write failure configuration: %v", err)
			}
			startup := filepath.Join(directory, "startup.sh")
			if err := os.WriteFile(startup, []byte("true\n"), 0o600); err != nil {
				t.Fatalf("write startup: %v", err)
			}
			command := exec.Command(binary, "run", "--config", config, "--profiles", profiles, "--endpoint", server.URL, "--output", filepath.Join(directory, "results"), "--worker", self, "--worker-arg=-test.run=^TestWorkloadProcess$", "--startup-file", startup, "--offer-defaults", fixtureOffers(t, directory))
			command.Env = append(os.Environ(), "VAST_API_KEY=fixture-key", "EVALUATION_WORKER_FIXTURE=1", "GORACE=atexit_sleep_ms=0")
			if err := command.Start(); err != nil {
				t.Fatalf("start failure fixture: %v", err)
			}
			if scenario == "cancel" {
				deadline := time.Now().Add(2 * time.Second)
				for time.Now().Before(deadline) {
					f.mu.Lock()
					allocated := len(f.labels) > 1
					f.mu.Unlock()
					if allocated {
						break
					}
					time.Sleep(5 * time.Millisecond)
				}
				if err := command.Process.Signal(syscall.SIGTERM); err != nil {
					t.Fatalf("interrupt evaluator: %v", err)
				}
			}
			if err := command.Wait(); err == nil {
				t.Fatal("failed or cancelled workload was reported as successful")
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			if len(f.labels) != 1 || f.labels[99] != "src_infra_image_gen" {
				t.Fatal("failure left evaluation instances or changed the personal rental")
			}
		})
	}
}

// fixtureGroup defines one compatible axis combination at the executable boundary.
func fixtureGroup(id string) map[string]any {
	return map[string]any{"id": id, "hardware_profiles": []string{"test"}, "model_profiles": []string{"test"}, "workflows": []string{"fixture"}, "case_ids": []string{"fixture"}, "seeds": []int{42}}
}

// fixtureOffers supplies a mutable public declaration whose predicates must be retained.
func fixtureOffers(t *testing.T, directory string) string {
	t.Helper()
	path := filepath.Join(directory, "offer-defaults.json")
	if err := os.WriteFile(path, []byte(`{"verified":{"eq":true},"datacenter":{"eq":true},"rentable":{"eq":false},"reliability":{"gte":0.99},"rented":{"eq":false},"num_gpus":{"eq":1},"type":"on-demand","order":[["dph_total","desc"]]}`), 0o600); err != nil {
		t.Fatalf("write fixture offer declaration: %v", err)
	}
	return path
}

// fixtureCommand composes public local inputs for the real controller executable.
func fixtureCommand(t *testing.T, endpoint, output string, groups []map[string]any) *exec.Cmd {
	t.Helper()
	binary, err := runfiles.Rlocation(os.Getenv("EVALUATOR_BINARY"))
	if err != nil {
		t.Fatalf("resolve evaluator: %v", err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve fixture workload: %v", err)
	}
	directory := t.TempDir()
	profiles := filepath.Join(directory, "profiles.json")
	if err := os.WriteFile(profiles, []byte(`{"image":"fixture","hardware_profiles":{"test":{"min_gpu_ram_gb":48,"disk_gb":100,"max_hourly_price":1},"other":{"min_gpu_ram_gb":48,"disk_gb":100,"max_hourly_price":1}},"model_profiles":{"test":{"downloads":[]},"other":{"downloads":[]}},"downloads":{},"custom_nodes":{}}`), 0o600); err != nil {
		t.Fatalf("write controller profiles: %v", err)
	}
	config := filepath.Join(directory, "config.json")
	value := map[string]any{"owner": "fixture", "max_parallel": 1, "max_lifetime_seconds": 15, "boot_timeout_seconds": 2, "workload_timeout_seconds": 10, "cleanup_timeout_seconds": 2, "poll_seconds": 0.02, "matrix": groups, "minimum_reliability": 0}
	if len(groups) > 0 && groups[0]["id"] == "matrix" {
		value["max_parallel"] = 2
	}
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode controller configuration: %v", err)
	}
	if err := os.WriteFile(config, data, 0o600); err != nil {
		t.Fatalf("write controller configuration: %v", err)
	}
	startup := filepath.Join(directory, "startup.sh")
	if err := os.WriteFile(startup, []byte("true\n"), 0o600); err != nil {
		t.Fatalf("write controller startup: %v", err)
	}
	command := exec.Command(binary, "run", "--config", config, "--profiles", profiles, "--endpoint", endpoint, "--output", output, "--worker", self, "--worker-arg=-test.run=^TestWorkloadProcess$", "--startup-file", startup, "--offer-defaults", fixtureOffers(t, directory))
	command.Env = append(os.Environ(), "VAST_API_KEY=fixture-key", "EVALUATION_WORKER_FIXTURE=1", "GORACE=atexit_sleep_ms=0")
	return command
}

// TestControllerCorrections exercises matrix expansion, offer reuse and early rejection.
func TestControllerCorrections(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	observations := map[string]any{}
	for _, scenario := range []string{"reuse", "matrix", "duplicate_seed", "duplicate_axis", "expansion_limit", "insecure_filters"} {
		t.Run(scenario, func(t *testing.T) {
			f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}, oneOffer: true}
			server := httptest.NewServer(f)
			defer server.Close()
			group := fixtureGroup(scenario)
			expected := 1
			expectedCells := 1
			switch scenario {
			case "insecure_filters":
				expected = 0
				expectedCells = 0
			case "reuse":
				group["hardware_profiles"] = []string{"test", "other"}
				expected = 2
				expectedCells = 2
			case "matrix":
				group["hardware_profiles"] = []string{"test", "other"}
				group["model_profiles"] = []string{"test", "other"}
				group["workflows"] = []string{"first", "second"}
				group["case_ids"] = []string{"first", "second"}
				expected = 4
				expectedCells = 16
			case "duplicate_seed":
				group["seeds"] = []int{42, 42}
				expected = 0
				expectedCells = 0
			case "duplicate_axis":
				group["hardware_profiles"] = []string{"test", "test"}
				expected = 0
				expectedCells = 0
			case "expansion_limit":
				axis := make([]string, 129)
				for i := range axis {
					axis[i] = fmt.Sprintf("choice_%d", i)
				}
				group["hardware_profiles"] = axis
				expected = 0
				expectedCells = 0
			}
			output := filepath.Join(t.TempDir(), "results")
			command := fixtureCommand(t, server.URL, output, []map[string]any{group})
			if scenario == "insecure_filters" {
				for i, arg := range command.Args {
					if arg == "--offer-defaults" {
						path := command.Args[i+1]
						data, err := os.ReadFile(path)
						if err != nil {
							t.Fatalf("read mutable offer fixture: %v", err)
						}
						data = []byte(strings.Replace(string(data), `"datacenter":{"eq":true}`, `"datacenter":{"eq":false}`, 1))
						if err := os.WriteFile(path, data, 0o600); err != nil {
							t.Fatalf("write immutable-filter fixture: %v", err)
						}
					}
				}
			}
			diagnostics, err := command.CombinedOutput()
			if expected == 0 && err == nil {
				t.Fatal("invalid matrix allocated successfully")
			}
			if expected > 0 && err != nil {
				t.Fatalf("execute controller correction: %v\n%s", err, diagnostics)
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			creations := 0
			for _, event := range f.events {
				if strings.HasPrefix(event, "create:") {
					creations++
				}
			}
			if creations != expected || len(f.labels) != 1 {
				t.Fatalf("expected %d independent rentals and no leftovers, observed %d and %d labels", expected, creations, len(f.labels))
			}
			if expected > 0 {
				if len(f.queries) == 0 || f.queries[0]["rentable"].(map[string]any)["eq"] != false {
					t.Fatal("controller overwrote the fixed shared offer predicate")
				}
				if _, present := f.queries[0]["reliability"]; present {
					t.Fatal("explicit zero reliability floor did not remove the inherited threshold")
				}
				files, err := filepath.Glob(filepath.Join(output, "*", "run.json"))
				if err != nil || len(files) != 1 {
					t.Fatalf("locate matrix report: %v", err)
				}
				data, err := os.ReadFile(files[0])
				if err != nil {
					t.Fatalf("read matrix report: %v", err)
				}
				var report struct {
					// Jobs records each expanded independent matrix cell.
					Jobs []map[string]any `json:"jobs"`
				}
				if err := json.Unmarshal(data, &report); err != nil {
					t.Fatalf("decode matrix report: %v", err)
				}
				coordinates := map[string]bool{}
				for _, job := range report.Jobs {
					key := fmt.Sprintf("%s/%s/%s/%s", job["hardware_profile"], job["model_profile"], job["workflow"], job["case_id"])
					coordinates[key] = true
				}
				if len(coordinates) != expectedCells {
					t.Fatal("matrix lost or duplicated a rental/model/workflow/case coordinate")
				}
			}
			if expectedCells > 0 {
				receipts, err := filepath.Glob(filepath.Join(output, "*", "*", "*", "*", "result.json"))
				if err != nil || len(receipts) != expectedCells {
					t.Fatalf("not every expanded workflow produced a retrieved result: %d, %v", len(receipts), err)
				}
			}
			observations[scenario] = map[string]any{"rentals": creations, "workloads": expectedCells, "remaining": 0, "live_operations": 0}
		})
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatalf("encode correction evidence: %v", err)
		}
		if err := os.WriteFile(filepath.Join(output, "controller-corrections.json"), data, 0o600); err != nil {
			t.Fatalf("save correction evidence: %v", err)
		}
	}
}

// TestControllerLockAcrossOutputs prevents two result directories sharing one cleanup namespace.
func TestControllerLockAcrossOutputs(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}}
	server := httptest.NewServer(f)
	defer server.Close()
	first := fixtureCommand(t, server.URL, filepath.Join(t.TempDir(), "first"), []map[string]any{fixtureGroup("hold")})
	if err := first.Start(); err != nil {
		t.Fatalf("start first controller: %v", err)
	}
	defer func() {
		if err := first.Process.Signal(syscall.SIGTERM); err != nil {
			t.Errorf("cancel first controller: %v", err)
		}
		if err := first.Wait(); err == nil {
			t.Error("cancelled controller succeeded")
		}
	}()
	deadline := time.Now().Add(3 * time.Second)
	allocated := false
	for time.Now().Before(deadline) {
		f.mu.Lock()
		allocated = len(f.labels) > 1
		f.mu.Unlock()
		if allocated {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if !allocated {
		t.Fatal("first controller did not allocate its controlled workload")
	}
	second := fixtureCommand(t, server.URL, filepath.Join(t.TempDir(), "second"), []map[string]any{fixtureGroup("fast")})
	if diagnostics, err := second.CombinedOutput(); err == nil || !strings.Contains(string(diagnostics), "already has a local controller") {
		t.Fatalf("second results directory bypassed namespace exclusion: %v\n%s", err, diagnostics)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	creations := 0
	for _, event := range f.events {
		if strings.HasPrefix(event, "create:") {
			creations++
		}
	}
	if creations != 1 || len(f.labels) != 2 {
		t.Fatal("second controller created a rental or swept the first controller's active workload")
	}
}

// TestProviderExitCleanup stops queued work promptly when the provider container exits.
func TestProviderExitCleanup(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}, oneOffer: true, exitAfterReady: true}
	server := httptest.NewServer(f)
	defer server.Close()
	command := fixtureCommand(t, server.URL, filepath.Join(t.TempDir(), "results"), []map[string]any{fixtureGroup("provider_exit")})
	started := time.Now()
	output, err := command.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "exited during workload") {
		t.Fatalf("provider exit did not stop queued work: %v\n%s", err, output)
	}
	if time.Since(started) > 2*time.Second {
		t.Fatal("provider exit waited for the full workload timeout")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.labels) != 1 || f.labels[99] != "src_infra_image_gen" {
		t.Fatal("provider exit left evaluation rentals or touched a personal rental")
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		data, err := json.Marshal(map[string]any{"provider_exit_cancelled_work": true, "remaining_evaluation_rentals": 0, "live_operations": 0})
		if err != nil {
			t.Fatalf("encode provider-exit evidence: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "provider-exit-cleanup.json"), data, 0o600); err != nil {
			t.Fatalf("save provider-exit evidence: %v", err)
		}
	}
}

// TestCleanupDisabled retains rentals on success and failure, then verifies explicit cleanup.
func TestCleanupDisabled(t *testing.T) {
	for _, outcome := range []string{"success", "failure"} {
		t.Run(outcome, func(t *testing.T) {
			t.Setenv("XDG_CACHE_HOME", t.TempDir())
			f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}}
			server := httptest.NewServer(f)
			defer server.Close()
			output := filepath.Join(t.TempDir(), "results")
			command := fixtureCommand(t, server.URL, output, []map[string]any{fixtureGroup(outcome)})
			command.Args = append(command.Args, "--cleanup=false")
			diagnostics, err := command.CombinedOutput()
			if (err != nil) != (outcome == "failure") {
				t.Fatalf("cleanup-disabled execution: %v\n%s", err, diagnostics)
			}
			f.mu.Lock()
			retained := len(f.labels)
			f.mu.Unlock()
			if retained != 2 {
				t.Fatalf("cleanup-disabled run retained %d labels, expected personal plus eval", retained)
			}
			reports, err := filepath.Glob(filepath.Join(output, "*", "run.json"))
			if err != nil || len(reports) != 1 {
				t.Fatalf("find retention report: %v", err)
			}
			data, err := os.ReadFile(reports[0])
			if err != nil {
				t.Fatal(err)
			}
			var report map[string]any
			if err := json.Unmarshal(data, &report); err != nil {
				t.Fatal(err)
			}
			if report["cleanup_confirmed"] != false || report["cleanup_skipped"] != true {
				t.Fatal("retention report incorrectly confirms destruction")
			}
			var config string
			for i, arg := range command.Args {
				if arg == "--config" {
					config = command.Args[i+1]
				}
			}
			cleanup := exec.Command(command.Path, "cleanup", "--config", config, "--endpoint", server.URL)
			cleanup.Env = command.Env
			if data, err := cleanup.CombinedOutput(); err != nil {
				t.Fatalf("explicit cleanup: %v\n%s", err, data)
			}
			f.mu.Lock()
			remaining := len(f.labels)
			f.mu.Unlock()
			if remaining != 1 {
				t.Fatal("explicit cleanup failed to destroy retained rental")
			}
			if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
				if err := os.WriteFile(filepath.Join(artifact, "cleanup-disabled-"+outcome+".json"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestRetainedProvisioningInputs rejects stale code and environments before reuse.
func TestRetainedProvisioningInputs(t *testing.T) {
	observations := map[string]any{}
	for _, scenario := range []string{"unchanged", "unused_qwen", "sequence-source", "gif-source", "qwen-gif-source", "requirements", "runtime-defaults", "startup-template", "recipe", "missing_fingerprint"} {
		t.Run(scenario, func(t *testing.T) {
			t.Setenv("XDG_CACHE_HOME", t.TempDir())
			f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}}
			server := httptest.NewServer(f)
			defer server.Close()
			directory := t.TempDir()
			// configure uses public fixture files and a deterministic renderer, while
			// exercising the real provisioning assembly and controller boundary.
			configure := func(command *exec.Cmd, name string) map[string]string {
				t.Helper()
				paths := map[string]string{}
				if scenario == "recipe" {
					for i, arg := range command.Args {
						if arg == "--startup-file" {
							paths["recipe"] = command.Args[i+1]
						}
					}
					return paths
				}
				for i, arg := range command.Args {
					if arg == "--startup-file" {
						command.Args = append(command.Args[:i], command.Args[i+2:]...)
						break
					}
				}
				for flag, content := range map[string]string{
					"sequence-source": "# public sequence source\n", "gif-source": "# public GIF source\n", "qwen-gif-source": "# public Qwen source\n",
					"requirements": "# public requirements\n", "runtime-defaults": `{"model_root":"/workspace/ComfyUI","model_directories":{},"arguments":"--listen 127.0.0.1"}`,
					"startup-template": "# public startup template\n", "terraform": "#!/bin/sh\nprintf '\"dHJ1ZQo=\"\\n'\n",
				} {
					path := filepath.Join(directory, name+"-"+flag)
					if err := os.WriteFile(path, []byte(content), 0o700); err != nil {
						t.Fatal(err)
					}
					paths[flag] = path
					command.Args = append(command.Args, "--"+flag, path)
				}
				for i, arg := range command.Args {
					if arg != "--profiles" {
						continue
					}
					data, err := os.ReadFile(command.Args[i+1])
					if err != nil {
						t.Fatal(err)
					}
					// The Qwen scenario includes its model-specific first-party node.
					if scenario == "qwen-gif-source" {
						data = []byte(strings.ReplaceAll(string(data), `"test":{"downloads":[]}`, `"qwen":{"downloads":[]}`))
					}
					if scenario == "requirements" {
						data = []byte(strings.ReplaceAll(string(data), `"test":{"downloads":[]}`, `"test":{"downloads":[],"python_requirements":"fixture"}`))
					}
					if err := os.WriteFile(command.Args[i+1], data, 0o600); err != nil {
						t.Fatal(err)
					}
				}
				command.Env = append(command.Env, "IMAGEGEN_WEB_USERNAME=fixture", "IMAGEGEN_WEB_PASSWORD="+strings.Repeat("x", 48))
				return paths
			}
			group := fixtureGroup("original")
			if scenario == "qwen-gif-source" {
				group["model_profiles"] = []string{"qwen"}
			}
			first := fixtureCommand(t, server.URL, filepath.Join(directory, "first"), []map[string]any{group})
			configure(first, "first")
			first.Args = append(first.Args, "--cleanup=false")
			if output, err := first.CombinedOutput(); err != nil {
				t.Fatalf("create retained fixture: %v\n%s", err, output)
			}
			reports, err := filepath.Glob(filepath.Join(directory, "first", "*", "run.json"))
			if err != nil || len(reports) != 1 {
				t.Fatalf("find retained report: %v", err)
			}
			source := filepath.Dir(reports[0])
			next := fixtureCommand(t, server.URL, filepath.Join(directory, "second"), []map[string]any{group})
			paths := configure(next, "second")
			next.Args = append(next.Args, "--cleanup=false", "--reuse-run", source)
			if scenario == "missing_fingerprint" {
				journals, err := filepath.Glob(filepath.Join(source, "*", "journal.json"))
				if err != nil || len(journals) != 1 {
					t.Fatalf("find retained journal: %v", err)
				}
				data, err := os.ReadFile(journals[0])
				if err != nil {
					t.Fatal(err)
				}
				var state map[string]any
				if err := json.Unmarshal(data, &state); err != nil {
					t.Fatal(err)
				}
				delete(state, "installation_fingerprint")
				data, err = json.Marshal(state)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(journals[0], data, 0o600); err != nil {
					t.Fatal(err)
				}
			} else if scenario == "runtime-defaults" {
				if err := os.WriteFile(paths[scenario], []byte(`{"model_root":"/workspace/ComfyUI","model_directories":{},"arguments":"--listen 127.0.0.1 --disable-metadata"}`), 0o600); err != nil {
					t.Fatal(err)
				}
			} else if scenario != "unchanged" {
				key := scenario
				if scenario == "unused_qwen" {
					key = "qwen-gif-source"
				}
				file, err := os.OpenFile(paths[key], os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := file.WriteString("\n "); err != nil {
					t.Fatal(err)
				}
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			f.mu.Lock()
			queries, creates := len(f.queries), len(f.createdOffers)
			labels, err := json.Marshal(f.labels)
			f.mu.Unlock()
			if err != nil {
				t.Fatal(err)
			}
			output, runErr := next.CombinedOutput()
			accepted := scenario == "unchanged" || scenario == "unused_qwen"
			if (runErr == nil) != accepted {
				t.Fatalf("provisioning reuse %s: %v\n%s", scenario, runErr, output)
			}
			f.mu.Lock()
			after, err := json.Marshal(f.labels)
			changed := len(f.queries) != queries || len(f.createdOffers) != creates || string(after) != string(labels)
			f.mu.Unlock()
			if err != nil || changed {
				t.Fatal("provisioning reuse changed account resources")
			}
			results, err := filepath.Glob(filepath.Join(directory, "second", "*", "*", "*", "*", "result.json"))
			if err != nil || (len(results) > 0) != accepted {
				t.Fatal("provisioning mismatch executed work or valid reuse lost its output")
			}
			observations[scenario] = map[string]any{"accepted": accepted, "new_allocations": 0, "live_operations": 0}
		})
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "retained-provisioning-inputs.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestRetainedRunReuse exercises repeat command invocations without new allocations.
func TestRetainedRunReuse(t *testing.T) {
	observations := map[string]any{}
	for _, scenario := range []string{"revised", "failure", "interrupted", "cleanup", "foreign", "missing", "stopped", "installation", "released", "duplicate", "endpoint", "missing_model"} {
		t.Run(scenario, func(t *testing.T) {
			t.Setenv("XDG_CACHE_HOME", t.TempDir())
			f := &fixture{labels: map[int]string{99: "src_infra_image_gen"}}
			server := httptest.NewServer(f)
			defer server.Close()
			directory := t.TempDir()
			groups := []map[string]any{fixtureGroup("original")}
			if scenario == "duplicate" {
				other := fixtureGroup("other")
				other["model_profiles"] = []string{"other"}
				groups = append(groups, other)
			}
			first := fixtureCommand(t, server.URL, filepath.Join(directory, "first"), groups)
			first.Args = append(first.Args, "--cleanup=false")
			if output, err := first.CombinedOutput(); err != nil {
				t.Fatalf("create retained fixture: %v\n%s", err, output)
			}
			reports, err := filepath.Glob(filepath.Join(directory, "first", "*", "run.json"))
			if err != nil || len(reports) != 1 {
				t.Fatalf("find retained run: %v", err)
			}
			source := filepath.Dir(reports[0])
			group := fixtureGroup("revision")
			group["seeds"] = []int{123, 456}
			group["sampling"] = map[string]any{"cfg": 4.5}
			if scenario == "failure" {
				group = fixtureGroup("failure")
			}
			if scenario == "interrupted" {
				group = fixtureGroup("cancel")
			}
			if scenario == "missing_model" {
				group["model_profiles"] = []string{"other"}
			}
			nextGroups := []map[string]any{group}
			if scenario == "duplicate" {
				nextGroups = append(nextGroups, groups[1])
			}
			next := fixtureCommand(t, server.URL, filepath.Join(directory, "second"), nextGroups)
			next.Args = append(next.Args, "--reuse-run", source)
			if scenario != "cleanup" {
				next.Args = append(next.Args, "--cleanup=false")
			}
			journals, err := filepath.Glob(filepath.Join(source, "*", "journal.json"))
			if err != nil {
				t.Fatal(err)
			}
			var firstIdentity int
			for _, path := range journals {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				var state map[string]any
				if err := json.Unmarshal(data, &state); err != nil {
					t.Fatal(err)
				}
				identity := int(state["instance_id"].(float64))
				if firstIdentity == 0 {
					firstIdentity = identity
				}
				if scenario == "duplicate" {
					state["instance_id"] = firstIdentity
				}
				if scenario == "released" {
					state["released_at"] = time.Now().UTC().Format(time.RFC3339Nano)
				}
				if scenario == "duplicate" || scenario == "released" {
					data, err := json.Marshal(state)
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(path, data, 0o600); err != nil {
						t.Fatal(err)
					}
				}
			}
			f.mu.Lock()
			switch scenario {
			case "foreign":
				f.labels[firstIdentity] = "src_infra_image_gen"
			case "missing":
				delete(f.labels, firstIdentity)
			case "stopped":
				f.providerStatus = "exited"
			}
			beforeLabels, err := json.Marshal(f.labels)
			beforeQueries, beforeCreates := len(f.queries), len(f.createdOffers)
			f.mu.Unlock()
			if err != nil {
				t.Fatal(err)
			}
			if scenario == "installation" {
				for i, arg := range next.Args {
					if arg != "--profiles" {
						continue
					}
					data, err := os.ReadFile(next.Args[i+1])
					if err != nil {
						t.Fatal(err)
					}
					data = []byte(strings.Replace(string(data), `"image":"fixture"`, `"image":"different"`, 1))
					if err := os.WriteFile(next.Args[i+1], data, 0o600); err != nil {
						t.Fatal(err)
					}
				}
			}
			if scenario == "endpoint" {
				data, err := os.ReadFile(reports[0])
				if err != nil {
					t.Fatal(err)
				}
				var report map[string]any
				if err := json.Unmarshal(data, &report); err != nil {
					t.Fatal(err)
				}
				report["provider_endpoint"] = "https://different.invalid/api/v0"
				data, err = json.Marshal(report)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(reports[0], data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			original, err := os.ReadFile(journals[0])
			if err != nil {
				t.Fatal(err)
			}
			var diagnostics []byte
			if scenario == "interrupted" {
				if err = next.Start(); err != nil {
					t.Fatal(err)
				}
				time.Sleep(100 * time.Millisecond)
				if err = next.Process.Signal(syscall.SIGTERM); err != nil {
					t.Fatal(err)
				}
				err = next.Wait()
			} else {
				diagnostics, err = next.CombinedOutput()
			}
			if (err == nil) != (scenario == "revised") {
				t.Fatalf("retained reuse %s: %v\n%s", scenario, err, diagnostics)
			}
			f.mu.Lock()
			afterLabels, encodeErr := json.Marshal(f.labels)
			changed := len(f.queries) != beforeQueries || len(f.createdOffers) != beforeCreates || string(beforeLabels) != string(afterLabels)
			f.mu.Unlock()
			if encodeErr != nil {
				t.Fatal(encodeErr)
			}
			if changed {
				t.Fatal("reuse searched offers, allocated or changed account rentals")
			}
			if scenario != "revised" && scenario != "failure" && scenario != "interrupted" {
				results, err := filepath.Glob(filepath.Join(directory, "second", "*", "*", "*", "*", "result.json"))
				if err != nil || len(results) != 0 {
					t.Fatal("invalid reuse executed a workload before rejecting its inputs")
				}
			}
			after, err := os.ReadFile(journals[0])
			if err != nil || string(after) != string(original) {
				t.Fatal("reuse overwrote its source journal")
			}
			if scenario == "revised" {
				results, err := filepath.Glob(filepath.Join(directory, "second", "*", "*", "*", "*", "result.json"))
				if err != nil || len(results) != 2 {
					t.Fatal("revised seed trials were not independently retrieved")
				}
				newJournals, err := filepath.Glob(filepath.Join(directory, "second", "*", "*", "journal.json"))
				if err != nil || len(newJournals) != 1 {
					t.Fatal("reuse journal absent")
				}
				data, err := os.ReadFile(newJournals[0])
				if err != nil {
					t.Fatal(err)
				}
				var state map[string]any
				if err := json.Unmarshal(data, &state); err != nil {
					t.Fatal(err)
				}
				if int(state["instance_id"].(float64)) != firstIdentity || state["reused_from"] != source {
					t.Fatal("reuse lost rental identity or provenance")
				}
				third := fixtureCommand(t, server.URL, filepath.Join(directory, "third"), []map[string]any{fixtureGroup("followup")})
				third.Args = append(third.Args, "--cleanup=false", "--reuse-run", filepath.Dir(filepath.Dir(newJournals[0])))
				if output, err := third.CombinedOutput(); err != nil {
					t.Fatalf("reuse an already reused rental: %v\n%s", err, output)
				}
				f.mu.Lock()
				changed := len(f.queries) != beforeQueries || len(f.createdOffers) != beforeCreates
				f.mu.Unlock()
				if changed {
					t.Fatal("third iteration allocated instead of preserving rental provenance")
				}
			}
			observations[scenario] = map[string]any{"account_unchanged": true, "new_allocations": 0, "new_offer_searches": 0, "live_operations": 0}
		})
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "retained-run-reuse.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
