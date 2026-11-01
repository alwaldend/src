package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/api/contracts"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// main runs the independent consumer outside Bazel actions and retains verified outputs.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var source, runner, output string
	command := &cobra.Command{Use: "consumer-check", SilenceUsage: true, Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if err := checkConsumer(cmd.Context(), source, runner, output); err != nil {
			return fmt.Errorf("validate independent consumer: %w", err)
		}
		return nil
	}}
	command.Flags().StringVar(&source, "source_root", "", "Repository root containing the first-party modules")
	command.Flags().StringVar(&runner, "al", "", "Native AL executable override (default: the source-built runfile)")
	command.Flags().StringVar(&output, "output", "", "Absent or empty directory for the consumer and result.json")
	command.SetContext(ctx)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

// checkConsumer builds with independent configuration and verifies the public output contract.
func checkConsumer(ctx context.Context, source, runner, output string) error {
	if runner == "" {
		resolved, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
		if err != nil {
			return fmt.Errorf("resolve source-built AL: %w", err)
		}
		runner = resolved
	}
	if !filepath.IsAbs(source) || !filepath.IsAbs(runner) || !filepath.IsAbs(output) {
		return fmt.Errorf("--source_root, --al, and --output must be absolute paths")
	}
	if err := os.MkdirAll(output, 0o700); err != nil {
		return fmt.Errorf("create consumer output: %w", err)
	}
	entries, err := os.ReadDir(output)
	if err != nil {
		return fmt.Errorf("inspect consumer output: %w", err)
	}
	if len(entries) != 0 {
		return fmt.Errorf("consumer output must be empty")
	}
	module := filepath.Join(source, "tools/rules_al")
	inputs := map[string]string{
		"MODULE.bazel": filepath.Join(module, "tests/consumer/MODULE.bazel.in"),
		"BUILD.bazel":  filepath.Join(module, "tests/consumer/BUILD.bazel.in"),
		"payload.txt":  filepath.Join(module, "tests/payload.txt"),
		"status.txt":   filepath.Join(module, "tests/status.txt"),
		".bazeliskrc":  filepath.Join(source, ".bazeliskrc"),
	}
	for name, input := range inputs {
		data, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("read consumer input %q: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(output, name), data, 0o600); err != nil {
			return fmt.Errorf("write consumer input %q: %w", name, err)
		}
	}
	if err := os.WriteFile(filepath.Join(output, ".bazelrc"), []byte("common:agent --color=no\n"), 0o600); err != nil {
		return fmt.Errorf("write consumer configuration: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	args := []string{"bazel", "build", "//:bundle", "@rules_al//:al"}
	for _, name := range []string{"rules_al", "rules_docs", "rules_skills", "rules_promptfoo", "rules_openspec"} {
		args = append(args, "--override_module="+name+"="+filepath.Join(source, "tools", name))
	}
	build := exec.CommandContext(ctx, runner, args...)
	build.Dir = output
	build.Stdout, build.Stderr = os.Stdout, os.Stderr
	// Bazel run variables identify the caller, whereas this command owns a new workspace.
	for _, value := range os.Environ() {
		if len(value) >= 6 && value[:6] == "BUILD_" {
			continue
		}
		build.Env = append(build.Env, value)
	}
	if err := build.Run(); err != nil {
		return fmt.Errorf("build external consumer: %w", err)
	}
	bundle := filepath.Join(output, "bazel-bin/bundle")
	manifest, err := os.ReadFile(filepath.Join(bundle, "release.json"))
	if err != nil {
		return fmt.Errorf("read consumer manifest: %w", err)
	}
	release := &contracts.Release{}
	if err := protojson.Unmarshal(manifest, release); err != nil {
		return fmt.Errorf("decode consumer manifest: %w", err)
	}
	if release.Name != "2026.44.0" || release.SourceCommit != "0123456789012345678901234567890123456789" || release.Channel != "release" || release.TreeState != "clean" || release.GetProject().GetSubdir() != "projects/consumer" || len(release.Items) != 1 {
		return fmt.Errorf("external consumer emitted an incorrect source identity")
	}
	file := release.Items[0].GetFile()
	if file.GetName() != "payload.txt" {
		return fmt.Errorf("external consumer emitted an incorrect payload set")
	}
	payload, err := os.ReadFile(filepath.Join(bundle, "files/payload.txt"))
	if err != nil {
		return fmt.Errorf("read consumer payload: %w", err)
	}
	expected, err := os.ReadFile(filepath.Join(output, "payload.txt"))
	if err != nil {
		return fmt.Errorf("read expected payload: %w", err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	verified := false
	for _, hash := range file.Hashes {
		if hash.Algo == "SHA-256" && hash.Content == digest {
			verified = true
		}
	}
	if !bytes.Equal(payload, expected) || file.Size != int64(len(payload)) || !verified {
		return fmt.Errorf("consumer payload bytes, size, or digest do not match")
	}
	result, err := json.MarshalIndent(map[string]any{
		"schema": "release/consumer_check/v1", "verified": true,
		"manifest": json.RawMessage(manifest), "payload_sha256": digest, "payload_size": len(payload),
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode consumer result: %w", err)
	}
	if err := os.WriteFile(filepath.Join(output, "result.json"), append(result, '\n'), 0o600); err != nil {
		return fmt.Errorf("write consumer result: %w", err)
	}
	fmt.Println("Verified external consumer: " + strconv.Itoa(len(payload)) + " bytes, SHA-256 " + digest)
	return nil
}
