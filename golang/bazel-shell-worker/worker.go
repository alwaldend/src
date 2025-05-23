package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/encoding/protodelim"

	"git.alwaldend.com/src/proto/bazel-worker/worker_protocol"
)

type Worker struct {
	stdout io.Writer
	stdin  protodelim.Reader
}

type WorkerRequestArguments struct {
	FlagfilePath string          `json:"flagfile_path"`
	Flagfile     *WorkerFlagfile `json:"flagfile"`
	WorkDir      string          `json:"work_dir"`
}

type WorkerLabel struct {
	Label string `json:"label"`
}

type WorkerTarget struct {
	Label string `json:"label"`
}

type WorkerFlagfileAttr struct {
	Cmd                  string            `json:"cmd"`
	Shell                string            `json:"shell"`
	CompatibleWith       []string          `json:"compatible_with"`
	Deprecation          string            `json:"deprecation"`
	ExecCompatibleWith   string            `json:"exec_oompatible_with"`
	ExecProperties       map[string]string `json:"exec_properties"`
	ExpectFailure        string            `json:"expect_failure"`
	Features             []string          `json:"features"`
	GeneratorFunction    string            `json:"generator_function"`
	GeneratorLocation    string            `json:"generator_location"`
	GeneratorName        string            `json:"generator_name"`
	Name                 string            `json:"name"`
	Outs                 []WorkerLabel     `json:"outs"`
	PackageMetadata      []WorkerTarget    `json:"package_metadata"`
	RestrictedTo         []string          `json:"restricted_to"`
	Srcs                 []WorkerLabel     `json:"srcs"`
	Tags                 []string          `json:"tags"`
	TargetCompatibleWith []string          `json:"target_compatible_with"`
	Testonly             bool              `json:"testonly"`
	Toolchains           []string          `json:"toolchains"`
	TransitiveConfigs    []string          `json:"transitive_configs"`
	Visibility           []string          `json:"visibility"`
	Worker               WorkerLabel       `json:"worker"`
	SetFlags             []string          `json:"set_flags"`
}

type WorkerFlagfile struct {
	Attr WorkerFlagfileAttr `json:"attr"`
}

type WorkerRequestCmd struct {
	Cmd    []string `json:"cmd"`
	Output string   `json:"stdout"`
}

type WorkerRequestData struct {
	Request   *worker_protocol.WorkRequest `json:"reqquest"`
	Arguments *WorkerRequestArguments      `json:"arguments"`
	Cmd       *WorkerRequestCmd            `json:"cmd"`
}

type WorkerErrorResponse struct {
	Error string             `json:"error"`
	Data  *WorkerRequestData `json:"data"`
}

func NewWorker(stdout io.Writer, stdin io.Reader) *Worker {
	return &Worker{stdout: stdout, stdin: bufio.NewReader(stdin)}
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
	request, err := self.readRequest()
	if err != nil {
		panic(fmt.Sprintf("could not read request: %s", err))
	}
	args, err := self.parseArguments(request)
	if err != nil {
		self.writeResponse(self.errorResponse(request, nil, fmt.Errorf("could not parse arguments: %w", err)))
	}
	data := &WorkerRequestData{Request: request, Arguments: args}
	response, err := self.handleRequest(data)
	if err != nil {
		self.writeResponse(self.errorResponse(request, data, err))
		return
	}
	self.writeResponse(response)
}

func (self *Worker) readRequest() (*worker_protocol.WorkRequest, error) {
	request := &worker_protocol.WorkRequest{}
	err := protodelim.UnmarshalFrom(self.stdin, request)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal work request: %w", err)
	}
	return request, nil
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

func (self *Worker) writeResponse(response *worker_protocol.WorkResponse) {
	_, err := protodelim.MarshalTo(self.stdout, response)
	if err != nil {
		panic(fmt.Sprintf("could not marshal response to stdout: %v", err))
	}
}

func (self *Worker) errorResponse(request *worker_protocol.WorkRequest, data *WorkerRequestData, err error) *worker_protocol.WorkResponse {
	response := WorkerErrorResponse{Data: data, Error: err.Error()}
	responseBytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("could not marshal response: %s", err))
	}
	return &worker_protocol.WorkResponse{
		ExitCode:  1,
		Output:    fmt.Sprintf("Data:%s\nOutput:\n%s", string(responseBytes), data.Cmd.Output),
		RequestId: request.RequestId,
	}
}

func (self *Worker) constructCmd(data *WorkerRequestData) ([]string, error) {
	var buffer bytes.Buffer
	t, err := template.New("cmd").Parse(data.Arguments.Flagfile.Attr.Cmd)
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %w", err)
	}

	if err := t.Execute(&buffer, data); err != nil {
		return nil, fmt.Errorf("could not execute template: %w", err)
	}
	cmdText := fmt.Sprintf("set %s\n%s", strings.Join(data.Arguments.Flagfile.Attr.SetFlags, " "), buffer.String())
	cmdArgs := []string{
		data.Arguments.Flagfile.Attr.Shell, "-c", cmdText,
	}
	return cmdArgs, nil
}

func (self *Worker) handleRequest(data *WorkerRequestData) (*worker_protocol.WorkResponse, error) {
	var output bytes.Buffer
	cmdArgs, err := self.constructCmd(data)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	data.Cmd = &WorkerRequestCmd{Cmd: cmdArgs, Output: output.String()}
	if err != nil {
		return nil, fmt.Errorf("could not run command: %w", err)
	}
	return &worker_protocol.WorkResponse{
		RequestId: data.Request.RequestId,
		ExitCode:  0,
		Output:    "success",
	}, nil
}
