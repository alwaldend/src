package versioncmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/versioning"
	"github.com/spf13/cobra"
)

// newShowCommand prints the selected source version as text or JSON.
func newShowCommand(options *versionOptions) *cobra.Command {
	format := "text"
	cmd := &cobra.Command{
		Use: "show", Short: "Print the calculated repository version", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if format != "text" && format != "json" {
				return fmt.Errorf("unknown format %q", format)
			}
			state, err := options.repositoryFor(cmd).Inspect()
			if err != nil {
				return fmt.Errorf("inspect repository version: %w", err)
			}
			if format == "json" {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(state); err != nil {
					return fmt.Errorf("write version state: %w", err)
				}
			} else if _, err := fmt.Fprintln(cmd.OutOrStdout(), state.Version); err != nil {
				return fmt.Errorf("write version: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	return cmd
}

// newMutationCommand prepares one kind of local release ref without publishing it.
func newMutationCommand(options *versionOptions, now func() time.Time, operation string) *cobra.Command {
	date, dryRun, switchBranch := now().UTC().Format("2006-01-02"), false, true
	cmd := &cobra.Command{
		Use: operation, Short: "Prepare local " + operation + " refs", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireAutomaticMutationContext(options.channel, options.release); err != nil {
				return fmt.Errorf("prepare refs: %w", err)
			}
			repository := options.repositoryFor(cmd)
			var value string
			var err error
			if operation == "release-tag" {
				value, err = repository.TagRelease(dryRun)
			} else {
				parsed, parseErr := versioning.ParseDate(date)
				if parseErr != nil {
					return fmt.Errorf("parse preparation date: %w", parseErr)
				}
				if operation == "nightly-tag" {
					value, err = repository.CreateNightly(parsed, dryRun)
				} else {
					value, err = repository.StartRelease(parsed, switchBranch, dryRun)
				}
			}
			if err != nil {
				return fmt.Errorf("prepare %s: %w", operation, err)
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), value); err != nil {
				return fmt.Errorf("write prepared ref: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate without changing Git refs")
	if operation != "release-tag" {
		cmd.Flags().StringVar(&date, "date", date, "Calendar date (YYYY-MM-DD)")
	}
	if operation == "release-start" {
		cmd.Flags().BoolVar(&switchBranch, "switch", true, "Switch to the created release branch")
	}
	return cmd
}

// newPlanCommand writes deterministic source-ref publication instructions.
func newPlanCommand(options *versionOptions) *cobra.Command {
	planPath, lease, atomic := "", "", true
	cmd := &cobra.Command{
		Use: "release-plan", Short: "Write a reviewed release-ref plan", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireAutomaticMutationContext(options.channel, options.release); err != nil {
				return fmt.Errorf("plan refs: %w", err)
			}
			if strings.TrimSpace(planPath) == "" {
				return fmt.Errorf("release-plan requires --plan")
			}
			repository := options.repositoryFor(cmd)
			state, err := repository.Inspect()
			if err != nil {
				return fmt.Errorf("inspect planned source: %w", err)
			}
			plan, err := versioning.BuildReleaseRefPlan(state, repository.TrunkBranch)
			if err != nil {
				return fmt.Errorf("build release-ref plan: %w", err)
			}
			plan.Atomic = atomic
			if atomic {
				plan.Lease = lease
			}
			if err := plan.Validate(); err != nil {
				return fmt.Errorf("validate release-ref plan: %w", err)
			}
			content, err := json.MarshalIndent(plan, "", "  ")
			if err != nil {
				return fmt.Errorf("encode release-ref plan: %w", err)
			}
			content = append(content, '\n')
			if err := writeAtomicFile(planPath, content, 0o600); err != nil {
				return fmt.Errorf("write release-ref plan: %w", err)
			}
			if _, err := cmd.OutOrStdout().Write(content); err != nil {
				return fmt.Errorf("print release-ref plan: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&planPath, "plan", "", "Output path for the reviewed release-ref plan")
	cmd.Flags().StringVar(&lease, "lease", "", "Explicit release-refs lease")
	cmd.Flags().BoolVar(&atomic, "atomic", true, "Require atomic multi-ref publication")
	return cmd
}

// newPublishCommand consumes the exact reviewed plan and emits verified publication evidence.
func newPublishCommand(options *versionOptions) *cobra.Command {
	planPath, receiptPath, remote, atomic := "", "", "origin", true
	cmd := &cobra.Command{
		Use: "release-publish", Short: "Publish guarded release refs from a reviewed plan", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireAutomaticMutationContext(options.channel, options.release); err != nil {
				return fmt.Errorf("publish refs: %w", err)
			}
			if strings.TrimSpace(planPath) == "" {
				return fmt.Errorf("release-publish requires --plan")
			}
			content, err := os.ReadFile(planPath)
			if err != nil {
				return fmt.Errorf("read release-ref plan: %w", err)
			}
			var plan versioning.ReleaseRefPlan
			decoder := json.NewDecoder(bytes.NewReader(content))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&plan); err != nil {
				return fmt.Errorf("decode release-ref plan: %w", err)
			}
			if err := plan.Validate(); err != nil {
				return fmt.Errorf("release-ref plan is invalid: %w", err)
			}
			if plan.Atomic != atomic {
				return fmt.Errorf("--atomic must match the reviewed plan; plan requires %v", plan.Atomic)
			}
			publisher := &versioning.LocalGitReleasePublisher{Git: versioning.ExecRunner{Directory: options.repository, Context: cmd.Context()}, Remote: remote}
			if err := publisher.Preflight(cmd.Context()); err != nil {
				return fmt.Errorf("release-ref publisher preflight: %w", err)
			}
			receipt, err := versioning.PublishReleaseRefs(cmd.Context(), publisher, plan)
			if err != nil {
				return fmt.Errorf("release-ref publication refused: %w", err)
			}
			data, err := json.MarshalIndent(receipt, "", "  ")
			if err != nil {
				return fmt.Errorf("encode release-ref receipt: %w", err)
			}
			data = append(data, '\n')
			if strings.TrimSpace(receiptPath) != "" {
				if err := writeAtomicFile(receiptPath, data, 0o600); err != nil {
					return fmt.Errorf("write release-ref receipt: %w", err)
				}
			}
			if _, err := cmd.OutOrStdout().Write(data); err != nil {
				return fmt.Errorf("print release-ref receipt: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&planPath, "plan", "", "Reviewed release-ref plan path")
	cmd.Flags().StringVar(&receiptPath, "receipt", "", "Output release-ref receipt path")
	cmd.Flags().StringVar(&remote, "remote", "origin", "Git remote whose refs are guarded")
	cmd.Flags().BoolVar(&atomic, "atomic", true, "Require atomic multi-ref publication")
	return cmd
}

// requireAutomaticMutationContext rejects read-only version selectors for ref preparation.
func requireAutomaticMutationContext(channel string, release string) error {
	if channel != "auto" || release != "" {
		return fmt.Errorf("--channel and --release select read-only or stamping context and cannot be used with ref mutations")
	}
	return nil
}

// writeAtomicFile installs content at path atomically with restrictive
// permissions, creating parent directories when needed.
func writeAtomicFile(path string, content []byte, mode os.FileMode) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("output path is required")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(absolute), ".plan.tmp-")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}
	temporaryPath := temporary.Name()
	cleanup := true
	defer func() {
		_ = temporary.Close()
		if cleanup {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(mode); err != nil {
		return fmt.Errorf("set output permissions: %w", err)
	}
	if _, err := temporary.Write(content); err != nil {
		return fmt.Errorf("write temporary output: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync temporary output: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary output: %w", err)
	}
	if err := os.Rename(temporaryPath, absolute); err != nil {
		return fmt.Errorf("install output atomically: %w", err)
	}
	cleanup = false
	return nil
}
