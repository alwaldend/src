package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/worker_protocol"
)

// Worker runs Hugo builds on behalf of Bazel persistent workers.
type Worker struct {
	protocol *WorkerProtocol
}

// WorkerRequestArguments carries the parsed flagfile for one work request.
type WorkerRequestArguments struct {
	FlagfilePath string          `json:"flagfile_path"`
	Flagfile     *WorkerFlagfile `json:"flagfile"`
	WorkDir      string          `json:"work_dir"`
}

// WorkerFlagfile is the JSON content passed from the Hugo rule.
type WorkerFlagfile struct {
	// SiteArchive is the path (relative to the execution root) of the site
	// source archive.
	SiteArchive string `json:"site_archive"`
	// Arguments are the Hugo arguments. --destination and --source are
	// appended by the worker.
	Arguments []string `json:"arguments"`
	// EnvScript is a shell snippet that exports the environment variables
	// Hugo needs.
	EnvScript string `json:"env_script"`
	// Hugo is the path (relative to the execution root) of the Hugo binary.
	Hugo string `json:"hugo"`
	// EnvFile is the path of the generated Hugo environment file.
	EnvFile string `json:"env_file"`
	// Postcss is the path (relative to the execution root) of the PostCSS
	// binary.
	Postcss string `json:"postcss"`
	// PostcssBindir is the Bazel output directory of the PostCSS binary.
	PostcssBindir string `json:"postcss_bindir"`
	// OutputDir is the Bazel declared output directory for the built site.
	OutputDir string `json:"output_dir"`
	// Shell is the shell used to run the environment and Hugo commands.
	Shell string `json:"shell"`
}

// WorkerRequestData is the per-request runtime state.
type WorkerRequestData struct {
	Request   *worker_protocol.WorkRequest `json:"request"`
	Arguments *WorkerRequestArguments      `json:"arguments"`
	Output    string                       `json:"output"`
}

// WorkerErrorResponse is returned on failure.
type WorkerErrorResponse struct {
	Error string             `json:"error"`
	Data  *WorkerRequestData `json:"data"`
}

func NewWorker(protocol *WorkerProtocol) *Worker {
	return &Worker{protocol: protocol}
}

func (self *Worker) Run(persist bool) {
	for {
		finished := self.RunOnce()
		if !persist || !finished {
			return
		}
	}
}

// RunOnce handles a single work request and returns whether it completed.
func (self *Worker) RunOnce() bool {
	request, err := self.protocol.ReadRequest()
	if err != nil {
		return false
	}
	args, err := self.parseArguments(request)
	if err != nil {
		self.protocol.WriteResponse(self.errorResponse(request, nil, fmt.Errorf("could not parse arguments: %w", err)))
		return true
	}
	data := &WorkerRequestData{Request: request, Arguments: args}
	response, err := self.handleRequest(data)
	if err != nil {
		self.protocol.WriteResponse(self.errorResponse(request, data, err))
		return true
	}
	self.protocol.WriteResponse(response)
	return true
}

func (self *Worker) parseArguments(request *worker_protocol.WorkRequest) (*WorkerRequestArguments, error) {
	flags := flag.FlagSet{}
	flagfileFlag := flags.String("flagfile", "", "argfile path")
	err := flags.Parse(request.Arguments)
	if err != nil {
		return nil, fmt.Errorf("could not parse arguments: %w", err)
	}
	flagfileBytes, err := os.ReadFile(*flagfileFlag)
	if err != nil {
		return nil, fmt.Errorf("could not read flag file: %w", err)
	}
	flagfile := &WorkerFlagfile{}
	err = json.Unmarshal(flagfileBytes, flagfile)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal flagfile contents: %s", err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &WorkerRequestArguments{
		FlagfilePath: *flagfileFlag,
		Flagfile:     flagfile,
		WorkDir:      workDir,
	}, nil
}

func (self *Worker) errorResponse(request *worker_protocol.WorkRequest, data *WorkerRequestData, err error) *worker_protocol.WorkResponse {
	response := WorkerErrorResponse{Data: data, Error: err.Error()}
	responseBytes, marshalErr := json.MarshalIndent(response, "", "    ")
	if marshalErr != nil {
		panic(fmt.Sprintf("could not marshal response: %s: %s", marshalErr, err))
	}
	output := ""
	if data != nil {
		output = data.Output
	}
	return &worker_protocol.WorkResponse{
		ExitCode:  1,
		Output:    fmt.Sprintf("Data:%s\nOutput:\n%s\nError:\n%s", string(responseBytes), output, err.Error()),
		RequestId: request.RequestId,
	}
}

// execDir returns the directory that mirrors the execution root for this
// request. Bazel starts sandboxed workers with a working directory whose
// layout mirrors the execroot: inputs are symlinked in and outputs must be
// written there, exec-relative, so Bazel can move them to the real execroot
// after the request. Multiplex sandboxed workers receive that directory
// relative to the worker working directory in WorkRequest.sandbox_dir;
// singleplex sandboxed workers use the worker working directory directly.
func (self *Worker) execDir(data *WorkerRequestData) string {
	if data.Request.SandboxDir != "" {
		return filepath.Join(data.Arguments.WorkDir, data.Request.SandboxDir)
	}
	return data.Arguments.WorkDir
}

// siteDir returns a per-request directory where the site archive is
// extracted, inside the request's exec-mirror directory.
func (self *Worker) siteDir(data *WorkerRequestData) string {
	return filepath.Join(self.execDir(data), fmt.Sprintf("hugo-site-%d", data.Request.RequestId))
}

func (self *Worker) handleRequest(data *WorkerRequestData) (*worker_protocol.WorkResponse, error) {
	flagfile := data.Arguments.Flagfile
	execDir := self.execDir(data)
	siteDir := self.siteDir(data)
	err := os.RemoveAll(siteDir)
	if err != nil {
		return nil, fmt.Errorf("could not remove site dir: %w", err)
	}
	err = os.MkdirAll(siteDir, 0o755)
	if err != nil {
		return nil, fmt.Errorf("could not create site dir: %w", err)
	}

	// Build the shell script that mirrors the non-worker Hugo action: export
	// the site environment, extract the site archive, create the project-local
	// PostCSS shim Hugo prefers, copy the generated Hugo environment file, and
	// run Hugo into the declared output directory.
	var script strings.Builder
	script.WriteString("set -eu\n")
	if flagfile.EnvScript != "" {
		script.WriteString(flagfile.EnvScript)
		script.WriteString("\n")
	}
	fmt.Fprintf(&script, "tar -xf '%s' -C '%s'\n", filepath.Join(execDir, flagfile.SiteArchive), siteDir)
	fmt.Fprintf(&script, "mkdir -p '%s/node_modules/.bin'\n", siteDir)
	postcssPath := filepath.Join(execDir, flagfile.Postcss)
	postcssBindir := filepath.Join(execDir, flagfile.PostcssBindir)
	fmt.Fprintf(&script, "cat >'%s/node_modules/.bin/postcss' <<'EOF'\n", siteDir)
	script.WriteString("#!/usr/bin/env sh\n")
	script.WriteString("set -eu\n")
	fmt.Fprintf(&script, "export BAZEL_BINDIR='%s'\n", postcssBindir)
	fmt.Fprintf(&script, "export NODE_PATH='%s/node_modules'\n", postcssBindir)
	fmt.Fprintf(&script, "exec '%s' \"$@\"\n", postcssPath)
	script.WriteString("EOF\n")
	fmt.Fprintf(&script, "chmod +x '%s/node_modules/.bin/postcss'\n", siteDir)
	fmt.Fprintf(&script, "mkdir -p '%s/static'\n", siteDir)
	fmt.Fprintf(&script, "cp '%s' '%s/static/hugo_env.txt'\n", filepath.Join(execDir, flagfile.EnvFile), siteDir)
	fmt.Fprintf(&script, "exec '%s' %s --destination '%s' --source '%s'\n",
		filepath.Join(execDir, flagfile.Hugo),
		quoteArgs(flagfile.Arguments),
		filepath.Join(execDir, flagfile.OutputDir),
		siteDir,
	)

	cmd := exec.Command(flagfile.Shell, "-c", script.String())
	cmd.Dir = execDir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	data.Output = output.String()
	if err != nil {
		return nil, fmt.Errorf("hugo build failed: %w", err)
	}
	return &worker_protocol.WorkResponse{RequestId: data.Request.RequestId}, nil
}

func quoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, "'"+strings.ReplaceAll(arg, "'", "'\\''")+"'")
	}
	return strings.Join(quoted, " ")
}
