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
	"os/exec"
	"strings"
	"text/template"
	"time"

	"git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/worker_protocol"
)

// maxShellText bounds each flagfile and rendered shell command.
const maxShellText = 4 << 20

// maxOutput bounds the bytes retained for a command failure diagnostic.
const maxOutput = 1 << 20

// workerArguments exposes the action paths and configuration to command templates.
type workerArguments struct {
	// FlagfilePath is the path Bazel supplied for the action configuration.
	FlagfilePath string
	// Flagfile contains the decoded JSON rule attributes.
	Flagfile *workerFlagfile
	// WorkDir is the worker's current execution directory.
	WorkDir string
}

// workerFlagfile describes the shell template and its rule-defined inputs.
type workerFlagfile struct {
	// Cmd is the Go template expanded into the shell script.
	Cmd string `json:"cmd"`
	// SetArgs contains arguments to the shell's set builtin.
	SetArgs []string `json:"set_args"`
	// Shell identifies the interpreter to execute.
	Shell string `json:"shell"`
	// Var exposes Bazel make variables to the template.
	Var map[string]string `json:"var"`
}

// workerTemplateContext provides the documented field names used by rule templates.
type workerTemplateContext struct {
	// Ctx holds the rule attributes.
	Ctx *workerFlagfile
	// Request holds the original Bazel request.
	Request *worker_protocol.WorkRequest
	// Arguments holds the resolved action context.
	Arguments *workerArguments
}

// boundedBuffer retains a prefix of a stream or rejects writes exceeding its limit.
type boundedBuffer struct {
	// buffer contains the retained bytes.
	buffer bytes.Buffer
	// limit is the maximum number of bytes retained.
	limit int
	// truncate allows excess child output to be drained without retaining it.
	truncate bool
	// truncated records whether diagnostic output was discarded.
	truncated bool
}

// Write accepts a complete write only when it fits or truncation was requested.
func (b *boundedBuffer) Write(p []byte) (int, error) {
	room := b.limit - b.buffer.Len()
	if len(p) <= room {
		b.buffer.Write(p)
		return len(p), nil
	}
	if !b.truncate {
		return 0, fmt.Errorf("rendered command exceeds %d bytes", b.limit)
	}
	b.buffer.Write(p[:room])
	b.truncated = true
	return len(p), nil
}

// perform converts an action failure into a response without stopping subsequent requests.
func perform(ctx context.Context, request *worker_protocol.WorkRequest) *worker_protocol.WorkResponse {
	response := &worker_protocol.WorkResponse{RequestId: request.RequestId}
	if err := runAction(ctx, request); err != nil {
		response.ExitCode = 1
		// Protobuf strings require UTF-8. Child output may contain arbitrary
		// bytes, including a partial rune at the diagnostic size limit.
		response.Output = strings.ToValidUTF8(err.Error(), "\uFFFD")
	}
	return response
}

// runAction expands the rule's template and executes its shell with bounded diagnostic output.
func runAction(ctx context.Context, request *worker_protocol.WorkRequest) error {
	flags := flag.NewFlagSet("shell action", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	path := flags.String("flagfile", "", "action configuration")
	if err := flags.Parse(request.Arguments); err != nil {
		return fmt.Errorf("parse shell request flags: %w", err)
	}
	if *path == "" {
		return fmt.Errorf("shell request requires --flagfile")
	}
	file, err := os.Open(*path)
	if err != nil {
		return fmt.Errorf("open shell flagfile %q: %w", *path, err)
	}
	data, readErr := io.ReadAll(io.LimitReader(file, maxShellText+1))
	if err := errors.Join(readErr, file.Close()); err != nil {
		return fmt.Errorf("read shell flagfile %q: %w", *path, err)
	}
	if len(data) > maxShellText {
		return fmt.Errorf("shell flagfile exceeds %d bytes", maxShellText)
	}
	var config workerFlagfile
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("decode shell flagfile %q: %w", *path, err)
	}
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve action working directory: %w", err)
	}
	tmpl, err := template.New("cmd").Option("missingkey=error").Parse(config.Cmd)
	if err != nil {
		return fmt.Errorf("parse command template: %w", err)
	}
	script := boundedBuffer{limit: maxShellText}
	if _, err := fmt.Fprintf(&script, "set %s\n", strings.Join(config.SetArgs, " ")); err != nil {
		return fmt.Errorf("write shell options: %w", err)
	}
	args := workerArguments{FlagfilePath: *path, Flagfile: &config, WorkDir: dir}
	if err := tmpl.Execute(&script, workerTemplateContext{Ctx: &config, Request: request, Arguments: &args}); err != nil {
		return fmt.Errorf("render command template %q: %w", *path, err)
	}
	output := boundedBuffer{limit: maxOutput, truncate: true}
	cmd := exec.CommandContext(ctx, config.Shell, "-c", script.buffer.String())
	// Bound the wait for descendants that keep the shell's output pipe open.
	cmd.WaitDelay = 5 * time.Second
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		text := output.buffer.String()
		if output.truncated {
			text += "\n[output truncated]"
		}
		return fmt.Errorf("execute shell action: %w\n%s", err, text)
	}
	return nil
}
