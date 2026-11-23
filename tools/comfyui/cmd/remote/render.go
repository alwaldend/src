package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// renderCommand runs the packaged CLI stages through the invocation's existing HTTPS bridge.
func renderCommand(cli, profile string, env []string) *cobra.Command {
	var workflow, outputDir string
	var slots, images []string
	cmd := &cobra.Command{
		Use:           "render",
		Short:         "Run a rental workflow and download its output files",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !filepath.IsAbs(workflow) || !filepath.IsAbs(outputDir) {
				return errors.New("--workflow and --output-dir must be absolute paths")
			}
			for _, slot := range slots {
				if !validAssignment(slot) {
					return errors.New("--slot requires ADDRESS=VALUE")
				}
			}
			for _, image := range images {
				_, path, _ := strings.Cut(image, "=")
				if !validAssignment(image) || !filepath.IsAbs(path) {
					return errors.New("--image requires ADDRESS=/absolute/path")
				}
			}
			if _, err := os.Stat(workflow); err != nil {
				return fmt.Errorf("open workflow: %w", err)
			}
			if err := os.MkdirAll(outputDir, 0o700); err != nil {
				return fmt.Errorf("prepare output directory: %w", err)
			}
			ctx := cmd.Context()
			for _, image := range images {
				address, path, _ := strings.Cut(image, "=")
				upload, err := renderCLI(ctx, cli, env, "upload", path)
				if err != nil {
					return fmt.Errorf("stage reference image: %w", err)
				}
				if len(upload.Data.Uploads) != 1 || upload.Data.Uploads[0].Name == "" {
					return errors.New("upload did not return one reference filename")
				}
				file := upload.Data.Uploads[0]
				slots = append(slots, address+"="+filepath.Join(file.Subfolder, file.Name))
			}
			if len(slots) > 0 {
				args := append([]string{"workflow", "set-slot", "--stdout", workflow}, slots...)
				edited, err := renderCLI(ctx, cli, env, args...)
				if err != nil {
					return fmt.Errorf("set workflow inputs: %w", err)
				}
				if len(edited.Data.Workflow) == 0 || string(edited.Data.Workflow) == "null" {
					return errors.New("set-slot did not return a workflow")
				}
				workflow = filepath.Join(profile, "workflow.json")
				if err := os.WriteFile(workflow, edited.Data.Workflow, 0o600); err != nil {
					return fmt.Errorf("write temporary workflow: %w", err)
				}
			}
			job, err := renderCLI(ctx, cli, env, "run", "--workflow", workflow, "--wait")
			if err != nil {
				return fmt.Errorf("run workflow: %w", err)
			}
			promptID := job.Data.PromptID
			if job.Data.Status != "completed" || promptID == "" {
				return fmt.Errorf("workflow did not report completion (prompt ID %q)", promptID)
			}
			download, err := renderCLI(ctx, cli, env, "download", "-o", outputDir, "--", promptID)
			if err != nil {
				return fmt.Errorf("download completed prompt %q: %w", promptID, err)
			}
			paths := make([]string, 0, len(download.Data.Files))
			for _, file := range download.Data.Files {
				paths = append(paths, file.Path)
			}
			if len(paths) == 0 {
				return fmt.Errorf("download for prompt %q returned no files", promptID)
			}
			result := map[string]any{"prompt_id": promptID, "files": paths}
			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(result); err != nil {
				return fmt.Errorf("write render result: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&workflow, "workflow", "", "Absolute path to the workflow JSON")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Absolute directory for downloaded files")
	cmd.Flags().StringArrayVar(&slots, "slot", nil, "Override a native workflow input (ADDRESS=VALUE, repeatable)")
	cmd.Flags().StringArrayVar(&images, "image", nil, "Upload a local image and set its native workflow input (ADDRESS=/absolute/path, repeatable)")
	return cmd
}

// validAssignment checks the shape of a slot assignment before any upload or generation.
func validAssignment(assignment string) bool {
	address, _, ok := strings.Cut(assignment, "=")
	return ok && address != "" && !strings.HasPrefix(address, "-")
}

// renderEnvelope contains the upstream responses consumed by the render stages.
type renderEnvelope struct {
	// OK indicates the operation succeeded.
	OK bool `json:"ok"`
	// Error contains the upstream diagnostic for a failed operation.
	Error struct {
		Message string `json:"message"`
		Hint    string `json:"hint"`
	} `json:"error"`
	// Data contains command-specific output.
	Data struct {
		PromptID string          `json:"prompt_id"`
		Status   string          `json:"status"`
		Workflow json.RawMessage `json:"workflow_json"`
		Uploads  []struct {
			Name      string `json:"cloud_name"`
			Subfolder string `json:"subfolder"`
		} `json:"uploads"`
		Files []struct {
			Path string `json:"path"`
		} `json:"files"`
	} `json:"data"`
}

// renderCLI runs a single pinned CLI command and consumes its JSON response.
func renderCLI(ctx context.Context, cli string, env []string, args ...string) (renderEnvelope, error) {
	child := exec.CommandContext(ctx, cli, append([]string{"--json"}, args...)...)
	child.Env = env
	child.Stderr = os.Stderr
	output, runErr := child.Output()
	var response renderEnvelope
	decodeErr := json.Unmarshal(output, &response)
	if runErr != nil {
		if decodeErr == nil && response.Error.Message != "" {
			return response, fmt.Errorf("%s %s: %w", response.Error.Message, response.Error.Hint, runErr)
		}
		return response, fmt.Errorf("run packaged ComfyUI CLI: %w", runErr)
	}
	if decodeErr != nil {
		return response, fmt.Errorf("decode packaged CLI response: %w", decodeErr)
	}
	if !response.OK {
		return response, fmt.Errorf("packaged CLI reported failure: %s", response.Error.Message)
	}
	return response, nil
}
