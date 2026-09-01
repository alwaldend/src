// Command control_status renders an offline runtime control status for one
// task namespace: package lifecycle states with deadlines, cross-process
// task lock/lease state, desired and observed revisions, contract hashes,
// and the published asset manifest.
//
// It performs no network or stateful operation beyond reading the task
// control root, and never mutates source. A missing optional control root or
// lease is a structured unavailable field.
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
	workspaceRoot string
	controlRoot   string
	namespace     string
	markdown      bool
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
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		opts.workspaceRoot = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	return opts, nil
}

type status struct {
	Schema      string                  `json:"schema"`
	Namespace   string                  `json:"namespace"`
	RuntimeID   string                  `json:"runtimeId"`
	Healthy     bool                    `json:"healthy"`
	Packages    []control.PackageStatus `json:"packages"`
	Asset       *control.AssetState     `json:"asset,omitempty"`
	Unavailable []string                `json:"unavailable,omitempty"`
	ObservedAt  string                  `json:"observedAt"`
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
		Schema:     "agents.alwaldend.com/control-status/v1alpha1",
		Namespace:  opts.namespace,
		ObservedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}
	kernel, err := control.New(control.KernelOptions{
		Root:      controlRoot,
		Namespace: opts.namespace,
		RuntimeID: "control-status-reader",
	})
	if err != nil {
		result.Unavailable = append(result.Unavailable, "control root unavailable: "+bounded(err))
		return encode(result, opts, stdout)
	}
	result.RuntimeID = kernel.RuntimeID()
	if snapshot, err := kernel.ReadSnapshot(); err == nil && len(snapshot) > 0 {
		result.Packages = snapshot
	} else {
		result.Packages = kernel.Status()
	}
	result.Healthy = kernel.Health()
	if asset, err := kernel.ReadAsset(); err == nil {
		value := asset
		result.Asset = &value
	} else {
		result.Unavailable = append(result.Unavailable, "asset unavailable: "+bounded(err))
	}
	return encode(result, opts, stdout)
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
	fmt.Fprintf(&builder, "- Runtime: `%s`\n", result.RuntimeID)
	fmt.Fprintf(&builder, "- Healthy: %t\n", result.Healthy)
	builder.WriteString("\n## Packages\n\n")
	for _, pkg := range result.Packages {
		fmt.Fprintf(&builder, "- `%s`: %s (deadline %s, observed %s)\n",
			pkg.ID, pkg.State, pkg.Deadline, pkg.ObservedRevision)
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
