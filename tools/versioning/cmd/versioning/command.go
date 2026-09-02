package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/versioning/internal/versioning"
)

const usage = `usage: versioning [--repo PATH] [--trunk BRANCH] [--channel CHANNEL] [--release YYYY.W] COMMAND [OPTIONS]

Commands:
  show             print the calculated repository version
  nightly-tag      create today's nightly tag on a clean trunk commit
  release-start    atomically create the weekly release branch and patch-zero tag
  release-tag      tag the current clean release-branch commit
  release-plan     write a reviewed, deterministic release-ref plan
  release-publish  consume a reviewed plan and publish guarded release refs
  bazel-status     emit deterministic Bazel workspace-status keys
  bazel -- ARGS    run Bazel with this binary as its workspace-status command
`

const globalUsage = "usage: versioning [--repo PATH] [--trunk BRANCH] " +
	"[--channel CHANNEL] [--release YYYY.W] "

var commandUsage = map[string]string{
	"show":            globalUsage + "show [--format text|json]\n",
	"nightly-tag":     globalUsage + "nightly-tag [--date YYYY-MM-DD] [--dry-run]\n",
	"release-start":   globalUsage + "release-start [--date YYYY-MM-DD] [--dry-run] [--switch=BOOL]\n",
	"release-tag":     globalUsage + "release-tag [--dry-run]\n",
	"release-plan":    globalUsage + "release-plan --plan PATH [--atomic=BOOL] [--lease LEASE]\n",
	"release-publish": globalUsage + "release-publish --plan PATH --receipt PATH [--atomic=BOOL]\n",
	"bazel-status":    globalUsage + "bazel-status\n",
	"bazel":           globalUsage + "bazel -- ARGS\n",
}

func execute(args []string, now func() time.Time, stdout io.Writer) error {
	global := flag.NewFlagSet("versioning", flag.ContinueOnError)
	global.SetOutput(io.Discard)
	repoPath := global.String("repo", ".", "Git repository path")
	trunk := global.String("trunk", "master", "trunk branch name")
	channel := global.String(
		"channel",
		"auto",
		"tag channel for an ambiguous exact commit: auto, release, or nightly",
	)
	release := global.String(
		"release",
		"",
		"release line for a detached release-branch checkout (YYYY.W)",
	)
	if err := global.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			_, _ = io.WriteString(stdout, usage)
			return nil
		}
		return err
	}
	if *channel != "auto" && *channel != "release" && *channel != "nightly" {
		return fmt.Errorf(
			"invalid channel %q; want auto, release, or nightly",
			*channel,
		)
	}
	if *release != "" {
		if _, err := versioning.ParseCalendar(*release); err != nil {
			return fmt.Errorf("invalid release context: %w", err)
		}
	}
	remaining := global.Args()
	if len(remaining) == 0 {
		_, _ = io.WriteString(stdout, usage)
		return nil
	}
	if len(remaining) == 2 && (remaining[1] == "--help" || remaining[1] == "-h") {
		if help, ok := commandUsage[remaining[0]]; ok {
			_, err := io.WriteString(stdout, help)
			return err
		}
	}
	repository := versioning.Repository{
		Git:              versioning.ExecRunner{Directory: *repoPath},
		TrunkBranch:      *trunk,
		RequestedChannel: *channel,
		RequestedRelease: *release,
	}

	switch remaining[0] {
	case "show":
		flags := newFlags("show")
		format := flags.String("format", "text", "output format: text or json")
		if done, err := parseFlags(flags, remaining[1:], "show", stdout); done || err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("show accepts no positional arguments")
		}
		state, err := repository.Inspect()
		if err != nil {
			return err
		}
		switch *format {
		case "text":
			_, err = fmt.Fprintln(stdout, state.Version)
		case "json":
			encoder := json.NewEncoder(stdout)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(state)
		default:
			return fmt.Errorf("unknown format %q", *format)
		}
		return err

	case "nightly-tag":
		if err := requireAutomaticMutationContext(*channel, *release); err != nil {
			return err
		}
		flags, date, dryRun := mutationFlags("nightly-tag", now)
		if done, err := parseFlags(flags, remaining[1:], "nightly-tag", stdout); done || err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("nightly-tag accepts no positional arguments")
		}
		parsed, err := versioning.ParseDate(*date)
		if err != nil {
			return err
		}
		tag, err := repository.CreateNightly(parsed, *dryRun)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdout, tag)
		return err

	case "release-start":
		if err := requireAutomaticMutationContext(*channel, *release); err != nil {
			return err
		}
		flags, date, dryRun := mutationFlags("release-start", now)
		switchBranch := flags.Bool("switch", true, "switch to the created release branch")
		if done, err := parseFlags(flags, remaining[1:], "release-start", stdout); done || err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("release-start accepts no positional arguments")
		}
		parsed, err := versioning.ParseDate(*date)
		if err != nil {
			return err
		}
		branch, err := repository.StartRelease(parsed, *switchBranch, *dryRun)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdout, branch)
		return err

	case "release-tag":
		if err := requireAutomaticMutationContext(*channel, *release); err != nil {
			return err
		}
		flags := newFlags("release-tag")
		dryRun := flags.Bool("dry-run", false, "validate and print without creating the tag")
		if done, err := parseFlags(flags, remaining[1:], "release-tag", stdout); done || err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("release-tag accepts no positional arguments")
		}
		tag, err := repository.TagRelease(*dryRun)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdout, tag)
		return err

	case "release-plan":
		if err := requireAutomaticMutationContext(*channel, *release); err != nil {
			return err
		}
		flags := newFlags("release-plan")
		planPath := flags.String("plan", "", "output path for the reviewed release-ref plan")
		atomic := flags.Bool("atomic", true, "require atomic multi-ref publication")
		lease := flags.String("lease", "", "explicit release-refs lease (empty acquires one)")
		if done, err := parseFlags(flags, remaining[1:], "release-plan", stdout); done || err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("release-plan accepts no positional arguments")
		}
		if strings.TrimSpace(*planPath) == "" {
			return fmt.Errorf("release-plan requires --plan")
		}
		state, err := repository.Inspect()
		if err != nil {
			return err
		}
		plan, err := versioning.BuildReleaseRefPlan(state, repository.TrunkBranch)
		if err != nil {
			return err
		}
		plan.Atomic = *atomic
		if !plan.Atomic {
			plan.Lease = ""
		} else {
			plan.Lease = *lease
		}
		if err := plan.Validate(); err != nil {
			return err
		}
		content, err := json.MarshalIndent(plan, "", "  ")
		if err != nil {
			return err
		}
		content = append(content, '\n')
		if err := writeAtomicFile(*planPath, content, 0o600); err != nil {
			return fmt.Errorf("write release-ref plan: %w", err)
		}
		_, err = stdout.Write(content)
		return err

	case "release-publish":
		if err := requireAutomaticMutationContext(*channel, *release); err != nil {
			return err
		}
		flags := newFlags("release-publish")
		planPath := flags.String("plan", "", "reviewed release-ref plan path")
		receiptPath := flags.String("receipt", "", "output release-ref receipt path")
		atomic := flags.Bool("atomic", true, "require atomic multi-ref publication")
		remote := flags.String("remote", "origin", "Git remote whose refs are guarded")
		if done, err := parseFlags(flags, remaining[1:], "release-publish", stdout); done || err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("release-publish accepts no positional arguments")
		}
		if strings.TrimSpace(*planPath) == "" {
			return fmt.Errorf("release-publish requires --plan")
		}
		content, err := os.ReadFile(*planPath)
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
		if plan.Atomic != *atomic {
			return fmt.Errorf(
				"--atomic must match the reviewed plan; plan requires %v",
				plan.Atomic,
			)
		}
		publisher := &versioning.LocalGitReleasePublisher{
			Git:    versioning.ExecRunner{Directory: *repoPath},
			Remote: *remote,
		}
		if err := publisher.Preflight(ctxVersioning()); err != nil {
			return fmt.Errorf("release-ref publisher preflight: %w", err)
		}
		receipt, err := versioning.PublishReleaseRefs(
			ctxVersioning(),
			publisher,
			plan,
		)
		if err != nil {
			return fmt.Errorf("release-ref publication refused: %w", err)
		}
		receiptContent, err := json.MarshalIndent(receipt, "", "  ")
		if err != nil {
			return fmt.Errorf("encode release-ref receipt: %w", err)
		}
		receiptContent = append(receiptContent, '\n')
		if strings.TrimSpace(*receiptPath) != "" {
			if err := writeAtomicFile(*receiptPath, receiptContent, 0o600); err != nil {
				return fmt.Errorf("write release-ref receipt: %w", err)
			}
		}
		_, err = stdout.Write(receiptContent)
		return err

	case "bazel-status":
		if len(remaining) != 1 {
			return fmt.Errorf("bazel-status accepts no arguments")
		}
		status, err := repository.BazelStatus()
		if err != nil {
			return err
		}
		_, err = io.WriteString(stdout, status)
		return err

	case "bazel":
		if len(remaining) < 3 || remaining[1] != "--" {
			return fmt.Errorf("bazel requires -- followed by a Bazel command")
		}
		return launchBazel(*repoPath, *trunk, *channel, *release, remaining[2:])

	case "help", "--help", "-h":
		_, err := io.WriteString(stdout, usage)
		return err
	default:
		return fmt.Errorf("unknown command %q\n%s", remaining[0], usage)
	}
}

func newFlags(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	return flags
}

func parseFlags(
	flags *flag.FlagSet,
	args []string,
	name string,
	stdout io.Writer,
) (bool, error) {
	err := flags.Parse(args)
	if errors.Is(err, flag.ErrHelp) {
		_, writeErr := io.WriteString(stdout, commandUsage[name])
		return true, writeErr
	}
	return false, err
}

func requireAutomaticMutationContext(channel string, release string) error {
	if channel != "auto" || release != "" {
		return fmt.Errorf(
			"--channel and --release select read-only or stamping context and cannot be used with ref mutations",
		)
	}
	return nil
}

func mutationFlags(name string, now func() time.Time) (*flag.FlagSet, *string, *bool) {
	flags := newFlags(name)
	date := flags.String("date", now().UTC().Format("2006-01-02"), "calendar date in YYYY-MM-DD")
	dryRun := flags.Bool("dry-run", false, "validate and print without changing Git refs")
	return flags, date, dryRun
}

func ctxVersioning() context.Context {
	return context.Background()
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
