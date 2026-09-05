// Command control_status renders persisted package observations and publication
// state for one task namespace. It cannot establish current runtime liveness.
//
// It performs no network or stateful operation beyond reading the task
// control root, and never mutates source. A missing snapshot or asset is a
// structured unavailable field.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/agents/control"
)

type options struct {
	workspaceRoot  string
	controlRoot    string
	namespace      string
	markdown       bool
	maxSnapshotAge time.Duration
}

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("control_status", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.controlRoot, "control-root",
		"out/agent-system/control", "task control root relative to workspace root")
	flags.StringVar(&opts.namespace, "namespace", "default",
		"task namespace to inspect")
	flags.BoolVar(&opts.markdown, "markdown", false,
		"emit human Markdown instead of JSON")
	flags.DurationVar(&opts.maxSnapshotAge, "max-snapshot-age", 0,
		"maximum accepted package observation age (0 leaves freshness unknown)")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		opts.workspaceRoot = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	if opts.namespace == "" {
		return options{}, fmt.Errorf("--namespace is required")
	}
	if opts.maxSnapshotAge < 0 {
		return options{}, fmt.Errorf("--max-snapshot-age must not be negative")
	}
	return opts, nil
}

type status struct {
	Schema         string                  `json:"schema"`
	Namespace      string                  `json:"namespace"`
	RuntimeID      string                  `json:"runtimeId"`
	Healthy        *bool                   `json:"healthy"`
	HealthBasis    string                  `json:"healthBasis"`
	Freshness      string                  `json:"freshness"`
	MaxSnapshotAge string                  `json:"maxSnapshotAge,omitempty"`
	SnapshotPath   string                  `json:"snapshotPath"`
	Packages       []control.PackageStatus `json:"packages"`
	Asset          *control.AssetState     `json:"asset,omitempty"`
	Unavailable    []string                `json:"unavailable,omitempty"`
	ObservedAt     string                  `json:"observedAt"`
}

func run(args []string, stdout io.Writer) error {
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}
	root, err := filepath.Abs(opts.workspaceRoot)
	if err != nil {
		return err
	}
	controlRoot, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(opts.controlRoot)))
	if err != nil {
		return err
	}
	result := status{
		Schema:       "agents.alwaldend.com/control-status/v1alpha1",
		Namespace:    opts.namespace,
		HealthBasis:  "persisted-package-snapshot",
		Freshness:    "unavailable",
		SnapshotPath: filepath.Join(controlRoot, "assets", opts.namespace, control.SnapshotFile),
	}
	snapshot, snapshotErr := control.ReadSnapshot(controlRoot, opts.namespace)
	now := time.Now().UTC()
	result.ObservedAt = now.Format(time.RFC3339Nano)
	if opts.maxSnapshotAge > 0 {
		result.MaxSnapshotAge = opts.maxSnapshotAge.String()
	}
	result.Unavailable = append(result.Unavailable,
		"runtime identity unavailable: package snapshots do not record their writer",
		"live runtime health unavailable: this command reads persisted observations")
	if snapshotErr == nil {
		result.Packages = snapshot
		result.Healthy, result.Freshness, err = snapshotHealth(snapshot, now, opts.maxSnapshotAge)
		if err != nil {
			result.Unavailable = append(result.Unavailable, "snapshot health unavailable: "+bounded(err))
		} else if result.Freshness == "unknown" {
			result.Unavailable = append(result.Unavailable,
				"snapshot freshness unknown: no maximum observation age was supplied")
		}
	} else {
		result.Unavailable = append(result.Unavailable, "package snapshot unavailable: "+bounded(snapshotErr))
	}
	if asset, err := control.ReadAsset(controlRoot, opts.namespace); err == nil {
		value := asset
		result.Asset = &value
	} else {
		result.Unavailable = append(result.Unavailable, "asset unavailable: "+bounded(err))
	}
	return encode(result, opts, stdout)
}

// snapshotHealth reports the recorded aggregate using the kernel's healthy
// states. An age bound is a caller policy, not evidence of runtime liveness.
func snapshotHealth(packages []control.PackageStatus, now time.Time, maxAge time.Duration) (*bool, string, error) {
	if len(packages) == 0 {
		return nil, "unavailable", fmt.Errorf("snapshot contains no package observations")
	}
	healthy := true
	stale := false
	seen := make(map[string]bool, len(packages))
	for _, pkg := range packages {
		if pkg.ID == "" || seen[pkg.ID] {
			return nil, "unavailable", fmt.Errorf("missing or duplicate package id %q", pkg.ID)
		}
		seen[pkg.ID] = true
		switch pkg.State {
		case control.PackageReady, control.PackageDegraded, control.PackageDisabled:
		case control.PackageLoading, control.PackageFailed, control.PackageTimeout, control.PackageDraining:
			healthy = false
		default:
			return nil, "unavailable", fmt.Errorf("package %s has unknown state %q", pkg.ID, pkg.State)
		}
		observed, err := time.Parse(time.RFC3339Nano, pkg.ObservationTime)
		if err != nil || observed.IsZero() || observed.After(now) {
			return nil, "unavailable", fmt.Errorf("package %s has invalid or future observation time", pkg.ID)
		}
		if pkg.Deadline != "" {
			if _, err := time.Parse(time.RFC3339Nano, pkg.Deadline); err != nil {
				return nil, "unavailable", fmt.Errorf("package %s has invalid deadline", pkg.ID)
			}
		}
		if maxAge > 0 && now.Sub(observed) > maxAge {
			stale = true
		}
	}
	if stale {
		return nil, "stale", fmt.Errorf("package observation exceeds maximum age %s", maxAge)
	}
	if maxAge == 0 {
		return &healthy, "unknown", nil
	}
	return &healthy, "within-age-bound", nil
}

func encode(result status, opts options, stdout io.Writer) error {
	sort.Strings(result.Unavailable)
	if opts.markdown {
		_, err := io.WriteString(stdout, renderMarkdown(result))
		return err
	}
	content, err := jsonIndent(result)
	if err != nil {
		return err
	}
	_, err = stdout.Write(content)
	return err
}

func renderMarkdown(result status) string {
	var builder strings.Builder
	builder.WriteString("# Control status\n\n")
	fmt.Fprintf(&builder, "- Namespace: `%s`\n", result.Namespace)
	builder.WriteString("- Runtime: unknown\n")
	if result.Healthy == nil {
		builder.WriteString("- Snapshot healthy: unavailable\n")
	} else {
		fmt.Fprintf(&builder, "- Snapshot healthy: %t\n", *result.Healthy)
	}
	fmt.Fprintf(&builder, "- Health basis: %s\n", result.HealthBasis)
	fmt.Fprintf(&builder, "- Snapshot freshness: %s\n", result.Freshness)
	if result.MaxSnapshotAge != "" {
		fmt.Fprintf(&builder, "- Maximum observation age: %s\n", result.MaxSnapshotAge)
	}
	fmt.Fprintf(&builder, "- Snapshot: `%s`\n", result.SnapshotPath)
	fmt.Fprintf(&builder, "- Read at: %s\n", result.ObservedAt)
	builder.WriteString("\n## Packages\n\n")
	for _, pkg := range result.Packages {
		fmt.Fprintf(&builder, "- `%s`: %s (deadline %s, revision %s, observed at %s)\n",
			pkg.ID, pkg.State, pkg.Deadline, pkg.ObservedRevision, pkg.ObservationTime)
		if pkg.Error != "" {
			fmt.Fprintf(&builder, "  - error: %s\n", pkg.Error)
		}
	}
	builder.WriteString("\n## Publication\n\n")
	if result.Asset == nil {
		builder.WriteString("None.\n")
	} else {
		asset := result.Asset
		fmt.Fprintf(&builder, "- Revision: `%s`\n", asset.Revision)
		fmt.Fprintf(&builder, "- Contract hash: `%s`\n", asset.ContractHash)
		fmt.Fprintf(&builder, "- Manifest: `%s`\n", asset.ManifestPath)
	}
	if len(result.Unavailable) > 0 {
		builder.WriteString("\n## Unavailable\n\n")
		for _, item := range result.Unavailable {
			fmt.Fprintf(&builder, "- %s\n", item)
		}
	}
	return builder.String()
}

func bounded(err error) string {
	if err == nil {
		return "unavailable"
	}
	message := strings.TrimSpace(err.Error())
	if len(message) > 200 {
		message = message[:200]
	}
	return message
}

func jsonIndent(value any) ([]byte, error) {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(content, '\n'), nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "control_status:", err)
		os.Exit(1)
	}
}
