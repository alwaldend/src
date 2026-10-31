package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// wrapper identifies the built command in the test runfiles.
var wrapper = flag.String("wrapper", "", "Scan wrapper runfile")

// TestMain also supplies a controlled scanner process for command-boundary cases.
func TestMain(m *testing.M) {
	if os.Getenv("SCAN_FIXTURE_MODE") == "1" {
		for _, required := range []string{"--json", "--no-verification", "--no-update", "--fail", "--fail-on-scan-errors"} {
			if !slices.Contains(os.Args, required) {
				os.Exit(97)
			}
		}
		allRefs := os.Getenv("SCAN_FIXTURE_ALL_REFS") == "1"
		if slices.Contains(os.Args, "--skip-additional-refs") == allRefs || slices.Contains(os.Args, "--branch") == allRefs {
			os.Exit(98)
		}
		if os.Getenv("SCAN_FIXTURE_WAIT") == "1" {
			time.Sleep(10 * time.Second)
		}
		if _, err := os.Stdout.WriteString(os.Getenv("SCAN_FIXTURE_OUTPUT")); err != nil {
			os.Exit(96)
		}
		status, err := strconv.Atoi(os.Getenv("SCAN_FIXTURE_STATUS"))
		if err != nil {
			os.Exit(95)
		}
		os.Exit(status)
	}
	os.Exit(m.Run())
}

// coordinates returns one accepted immutable source location.
func coordinates() map[string]any {
	return map[string]any{"detector": "URI", "path": "fixture.txt", "commit": strings.Repeat("a", 40), "line": 3, "reason": "Synthetic fixture"}
}

// policy serializes entries using the public policy contract.
func policy(entries ...map[string]any) string {
	if entries == nil {
		entries = []map[string]any{}
	}
	data, err := json.Marshal(map[string]any{"schema": "trufflehog_allowlist/v1", "findings": entries})
	if err != nil {
		panic(err)
	}
	return string(data)
}

// finding emits the pinned TruffleHog Git-result shape with a visible fixture value.
func finding(entry map[string]any) string {
	data, err := json.Marshal(map[string]any{
		"DetectorName": entry["detector"], "Verified": false, "Raw": "original-fixture-value",
		"SourceMetadata": map[string]any{"Data": map[string]any{"Git": map[string]any{
			"file": entry["path"], "commit": entry["commit"], "line": entry["line"],
		}}},
	})
	if err != nil {
		panic(err)
	}
	return string(data) + "\n"
}

// scenario describes a complete wrapper invocation and its expected result.
type scenario struct {
	// name identifies the behavioral case in the retained report.
	name string
	// policy is the complete on-disk JSON allowlist.
	policy string
	// output is the scanner's result stream.
	output string
	// scannerStatus is the controlled scanner process exit code.
	scannerStatus int
	// wantStatus is the wrapper's expected process exit code.
	wantStatus int
	// suppressed is the expected approved count when nonnegative.
	suppressed int
	// unmatched is the expected actionable count when nonnegative.
	unmatched int
	// allRefs selects the all-local-refs command contract.
	allRefs bool
	// wait exercises scanner cancellation.
	wait bool
	// missingScanner exercises launch failure.
	missingScanner bool
}

// TestScanCommand verifies filtering, process errors, and scan-only flags end to end.
func TestScanCommand(t *testing.T) {
	binary, err := runfiles.Rlocation(*wrapper)
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	accepted := coordinates()
	cases := []scenario{
		{name: "empty", policy: policy(), scannerStatus: 0, wantStatus: 0},
		{name: "accepted", policy: policy(accepted), output: finding(accepted), scannerStatus: 183, wantStatus: 0, suppressed: 1},
		{name: "all_refs", policy: policy(accepted), output: finding(accepted), scannerStatus: 183, wantStatus: 0, suppressed: 1, allRefs: true},
		{name: "unmatched", policy: policy(), output: finding(accepted), scannerStatus: 183, wantStatus: 183, unmatched: 1},
		{name: "scanner_error_after_accepted", policy: policy(accepted), output: finding(accepted), scannerStatus: 1, wantStatus: 1, suppressed: -1, unmatched: -1},
		{name: "truncated_json", policy: policy(), output: `{"DetectorName":`, scannerStatus: 183, wantStatus: 1, suppressed: -1, unmatched: -1},
		{name: "missing_metadata", policy: policy(), output: `{"DetectorName":"URI","Verified":false,"Raw":"original-fixture-value"}`, scannerStatus: 183, wantStatus: 1, suppressed: -1, unmatched: -1},
		{name: "findings_with_success_code", policy: policy(accepted), output: finding(accepted), scannerStatus: 0, wantStatus: 1, suppressed: -1, unmatched: -1},
		{name: "finding_code_without_findings", policy: policy(), scannerStatus: 183, wantStatus: 1, suppressed: -1, unmatched: -1},
		{name: "timeout", policy: policy(), scannerStatus: 0, wantStatus: 1, wait: true, suppressed: -1, unmatched: -1},
		{name: "launch_failure", policy: policy(), wantStatus: 1, missingScanner: true, suppressed: -1, unmatched: -1},
		{name: "duplicate_policy", policy: policy(accepted, accepted), wantStatus: 1, suppressed: -1, unmatched: -1},
		{name: "invalid_json_policy", policy: "{", wantStatus: 1, suppressed: -1, unmatched: -1},
	}
	for key, value := range map[string]any{"detector": "OtherDetector", "path": "other.txt", "commit": strings.Repeat("b", 40), "line": 4} {
		changed := coordinates()
		changed[key] = value
		cases = append(cases, scenario{name: "different_" + key, policy: policy(accepted), output: finding(changed), scannerStatus: 183, wantStatus: 183, unmatched: 1})
	}
	for key, value := range map[string]any{"commit": "HEAD", "path": "../fixture.txt", "line": 0, "reason": "", "value_sha256": "unsupported"} {
		invalid := coordinates()
		invalid[key] = value
		cases = append(cases, scenario{name: "invalid_" + key, policy: policy(invalid), wantStatus: 1, suppressed: -1, unmatched: -1})
	}
	changed := coordinates()
	changed["line"] = 7
	cases = append(cases, scenario{name: "mixed_findings", policy: policy(accepted), output: finding(accepted) + finding(changed), scannerStatus: 183, wantStatus: 183, suppressed: 1, unmatched: 1})
	slices.SortFunc(cases, func(a, b scenario) int { return strings.Compare(a.name, b.name) })
	report := []map[string]any{}
	for _, tc := range cases {
		passed := t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			policyPath := filepath.Join(dir, "allowlist.json")
			if err := os.WriteFile(policyPath, []byte(tc.policy), 0o600); err != nil {
				t.Fatal(err)
			}
			scanner := helper
			if tc.missingScanner {
				scanner = filepath.Join(dir, "missing-scanner")
			}
			args := []string{"--scanner", scanner, "--repository", dir, "--allowlist", policyPath, "--timeout", "2s"}
			if tc.allRefs {
				args = append(args, "--all-refs")
			}
			if tc.wait {
				args = append(args, "--timeout", "30ms")
			}
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, binary, args...)
			cmd.Env = append(os.Environ(), "SCAN_FIXTURE_MODE=1", "SCAN_FIXTURE_OUTPUT="+tc.output, "SCAN_FIXTURE_STATUS="+strconv.Itoa(tc.scannerStatus), "SCAN_FIXTURE_ALL_REFS="+strconv.Itoa(boolInt(tc.allRefs)), "SCAN_FIXTURE_WAIT="+strconv.Itoa(boolInt(tc.wait)))
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			status := 0
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok {
					status = exit.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			if status != tc.wantStatus {
				t.Fatalf("status %d, want %d; stdout %s; stderr %s", status, tc.wantStatus, stdout.String(), stderr.String())
			}
			if tc.name == "missing_metadata" && !strings.Contains(stdout.String(), "original-fixture-value") {
				t.Fatal("a malformed result lost its original evidence")
			}
			if tc.suppressed >= 0 {
				var summary map[string]any
				for _, line := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
					var event map[string]any
					if err := json.Unmarshal([]byte(line), &event); err != nil {
						t.Fatal(err)
					}
					if event["type"] == "summary" {
						summary = event
						continue
					}
					if !strings.Contains(line, "original-fixture-value") {
						t.Fatal("original finding value was not preserved")
					}
				}
				if summary["suppressed"] != float64(tc.suppressed) || summary["unmatched"] != float64(tc.unmatched) {
					t.Fatalf("unexpected summary: %v", summary)
				}
			}
		})
		report = append(report, map[string]any{"case": tc.name, "passed": passed, "expected_exit": tc.wantStatus})
	}
	artifactDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if artifactDir == "" {
		artifactDir = t.TempDir()
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "e2e-report.json"), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

// boolInt adapts a fixture switch to its process environment representation.
func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
