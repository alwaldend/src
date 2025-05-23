package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"git.alwaldend.com/src/proto/bazel-worker/worker_protocol"
)

// worker
type Worker struct {
	protocol *WorkerProtocol
}

type WorkerRequestArguments struct {
	FlagfilePath string `json:"flagfile_path"`
	// flagfile contents
	Flagfile *WorkerFlagfile `json:"flagfile"`
	// current working directory
	WorkDir string `json:"work_dir"`
}

// flagfile contents passed from the rule
type WorkerFlagfile struct {
	// Cmd passed from the rule
	Cmd string `json:"cmd"`
	// args for set comand
	SetArgs []string `json:"set_args"`
	// shell executable to use
	Shell string `json:"shell"`
	// Make vars
	Var map[string]string `json:"var"`
}

type WorkerRequestCmd struct {
	// full command to execute
	Cmd []string `json:"cmd"`
	// comand output (both stdout and stderr)
	Output string `json:"stdout"`
}

// data request to handle the request
type WorkerRequestData struct {
	Request   *worker_protocol.WorkRequest `json:"request"`
	Arguments *WorkerRequestArguments      `json:"arguments"`
	Cmd       *WorkerRequestCmd            `json:"cmd"`
}

// context for the cmd template
type WorkerTemplateContext struct {
	Ctx       *WorkerFlagfile
	Request   *worker_protocol.WorkRequest
	Arguments *WorkerRequestArguments
}

// data to create an error response
type WorkerErrorResponse struct {
	Error string             `json:"error"`
	Data  *WorkerRequestData `json:"data"`
}

func NewWorker(protocol *WorkerProtocol) *Worker {
	return &Worker{protocol: protocol}
}

func (self *Worker) Run(persist bool) {
	if !persist {
		self.RunOnce()
	}
	for {
		self.RunOnce()
	}
}

func (self *Worker) RunOnce() {
	request, err := self.protocol.ReadRequest()
	if err != nil {
		panic(fmt.Sprintf("could not read request: %s", err))
	}
	args, err := self.parseArguments(request)
	if err != nil {
		self.protocol.WriteResponse(self.errorResponse(request, nil, fmt.Errorf("could not parse arguments: %w", err)))
	}
	data := &WorkerRequestData{Request: request, Arguments: args, Cmd: &WorkerRequestCmd{}}
	response, err := self.handleRequest(data)
	if err != nil {
		self.protocol.WriteResponse(self.errorResponse(request, data, err))
		return
	}
	self.protocol.WriteResponse(response)
}

func (self *Worker) parseArguments(request *worker_protocol.WorkRequest) (*WorkerRequestArguments, error) {
	flags := flag.FlagSet{}
	flagfileFlag := flags.String("flagfile", "", "argfile path")
	err := flags.Parse(request.Arguments)
	if err != nil {
		return nil, fmt.Errorf("could not parse arguments: %w", err)
	}
	flagfileByts, err := os.ReadFile(*flagfileFlag)
	if err != nil {
		return nil, fmt.Errorf("could not read flag file: %w", err)
	}
	flagfile := &WorkerFlagfile{}
	err = json.Unmarshal(flagfileByts, flagfile)
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
	return &worker_protocol.WorkResponse{
		ExitCode:  1,
		Output:    fmt.Sprintf("Data:%s\nOutput:\n%s\nError:\n%s", string(responseBytes), data.Cmd.Output, err.Error()),
		RequestId: request.RequestId,
	}
}

func (self *Worker) constructCmd(data *WorkerRequestData) ([]string, error) {
	var buffer bytes.Buffer
	t, err := template.New("cmd").Parse(data.Arguments.Flagfile.Cmd)
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %w", err)
	}

	if err := t.Execute(&buffer, WorkerTemplateContext{
		Ctx:       data.Arguments.Flagfile,
		Request:   data.Request,
		Arguments: data.Arguments,
	}); err != nil {
		return nil, fmt.Errorf("could not execute template: %w", err)
	}
	cmdText := fmt.Sprintf("set %s\n%s", strings.Join(data.Arguments.Flagfile.SetArgs, " "), buffer.String())
	cmdArgs := []string{
		data.Arguments.Flagfile.Shell, "-c", cmdText,
	}
	return cmdArgs, nil
}

func (self *Worker) handleRequest(data *WorkerRequestData) (*worker_protocol.WorkResponse, error) {
	var output bytes.Buffer
	cmdArgs, err := self.constructCmd(data)
	if err != nil {
		return nil, err
	}
	data.Cmd.Cmd = cmdArgs
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	data.Cmd.Output = output.String()
	if err != nil {
		return nil, fmt.Errorf("could not run command: %w", err)
	}
	return &worker_protocol.WorkResponse{RequestId: data.Request.RequestId}, nil
}
