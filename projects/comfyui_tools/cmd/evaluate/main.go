package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/comfyui_tools/internal/evaluation"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
)

// main propagates ordinary interruption so the runner executes its final sweep.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// resolve accepts an existing caller path or a declared Bazel runfile path.
func resolve(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	if _, err := os.Stat(path); err == nil {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("normalize existing input path: %w", err)
		}
		return absolute, nil
	}
	resolved, err := runfiles.Rlocation(path)
	if err != nil {
		return "", fmt.Errorf("resolve packaged path %q: %w", path, err)
	}
	return resolved, nil
}

// command adapts declarative infrastructure arguments to reusable operations.
func command() *cobra.Command {
	root := &cobra.Command{Use: "evaluate", Short: "Run and collect bounded ComfyUI evaluations", SilenceUsage: true, SilenceErrors: true}
	for _, operation := range []string{"run", "cleanup", "check"} {
		var options evaluation.Options
		cmd := &cobra.Command{Use: operation, RunE: func(cmd *cobra.Command, args []string) error {
			options.Progress = cmd.OutOrStdout()
			if operation == "run" {
				enabled, err := cmd.Flags().GetBool("cleanup")
				if err != nil {
					return fmt.Errorf("read cleanup option: %w", err)
				}
				options.DisableCleanup = !enabled
			}
			paths := []*string{&options.Configuration, &options.Profiles, &options.Cases, &options.Workflows, &options.Terraform, &options.StartupTemplate, &options.RuntimeDefaults, &options.Requirements, &options.SequenceSource, &options.GIFSource, &options.QwenGIFSource, &options.Bridge, &options.CLI, &options.CAFile, &options.Worker, &options.StartupFile}
			for _, path := range paths {
				value, err := resolve(*path)
				if err != nil {
					return fmt.Errorf("resolve evaluator inputs: %w", err)
				}
				*path = value
			}
			if value, err := resolve(options.OfferDefaults); err != nil {
				return fmt.Errorf("resolve offer defaults: %w", err)
			} else {
				options.OfferDefaults = value
			}
			if options.Cases != "" {
				if info, err := os.Stat(options.Cases); err == nil && !info.IsDir() {
					options.Cases = filepath.Dir(options.Cases)
				}
			}
			if options.Workflows != "" {
				if info, err := os.Stat(options.Workflows); err == nil && !info.IsDir() {
					options.Workflows = filepath.Dir(filepath.Dir(options.Workflows))
				}
			}
			if workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY"); workspace != "" && !filepath.IsAbs(options.Output) {
				options.Output = filepath.Join(workspace, options.Output)
			}
			if operation == "cleanup" {
				if err := evaluation.Cleanup(cmd.Context(), options); err != nil {
					return fmt.Errorf("cleanup evaluation namespace: %w", err)
				}
				return nil
			}
			if operation == "check" {
				if err := evaluation.Check(cmd.Context(), options); err != nil {
					return fmt.Errorf("check packaged evaluation inputs: %w", err)
				}
				return nil
			}
			if err := evaluation.Run(cmd.Context(), options); err != nil {
				return fmt.Errorf("run evaluation suite: %w", err)
			}
			return nil
		}}
		flags := cmd.Flags()
		flags.StringSliceVar(&options.Groups, "group", nil, "Select named matrix groups; cleanup still covers the whole namespace")
		if operation == "run" {
			flags.Bool("cleanup", true, "Destroy rentals and sweep the namespace when the run ends")
			flags.Int64SliceVar(&options.ExcludedOffers, "exclude-offer", nil, "Skip known failed offer IDs for this run only")
		}
		flags.StringVar(&options.Configuration, "config", "", "Evaluation matrix JSON")
		flags.StringVar(&options.Profiles, "profiles", "", "Existing public model/hardware registry")
		flags.StringVar(&options.Endpoint, "endpoint", "https://console.vast.ai/api/v0", "Vast API root")
		flags.StringVar(&options.Output, "output", "out/imagegen/evaluations/results", "Private run journals and results")
		flags.StringVar(&options.Cases, "cases", "", "Reusable case directory")
		flags.StringVar(&options.Workflows, "workflows", "", "Generated native workflow directory")
		flags.StringVar(&options.Terraform, "terraform", "", "Pinned pure-template compiler")
		flags.StringVar(&options.StartupTemplate, "startup-template", "", "Canonical Vast startup template")
		flags.StringVar(&options.RuntimeDefaults, "runtime-defaults", "", "Shared ComfyUI runtime JSON")
		flags.StringVar(&options.OfferDefaults, "offer-defaults", "", "Shared image-compatible offer filters")
		flags.StringVar(&options.Requirements, "requirements", "", "Existing approved node requirements")
		flags.StringVar(&options.SequenceSource, "sequence-source", "", "Existing Prompt Sequence source")
		flags.StringVar(&options.GIFSource, "gif-source", "", "Existing chained GIF source")
		flags.StringVar(&options.QwenGIFSource, "qwen-gif-source", "", "Qwen image-edit GIF source")
		flags.StringVar(&options.Bridge, "bridge", "", "Existing authenticated CLI bridge")
		flags.StringVar(&options.CLI, "cli", "", "Pinned ComfyUI CLI")
		flags.StringVar(&options.CAFile, "ca-file", "", "Pinned Vast root certificate")
		flags.StringVar(&options.Worker, "worker", "", "Optional configured workload executable")
		flags.StringSliceVar(&options.WorkerArgs, "worker-arg", nil, "Literal optional workload arguments")
		flags.StringVar(&options.StartupFile, "startup-file", "", "Optional explicitly supplied startup recipe")
		cmd.MarkFlagRequired("config")
		root.AddCommand(cmd)
	}
	return root
}
