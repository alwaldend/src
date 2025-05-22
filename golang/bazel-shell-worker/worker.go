package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

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
}

type WorkerLabel struct {
	Label string `json:"label"`
}

type WorkerTarget struct {
	Label string `json:"label"`
}

type WorkerFlagfileAttr struct {
	Cmd                  string            `json:"cmd"`
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
	Srcs                 []string          `json:"srcs"`
	Tags                 []string          `json:"tags"`
	TargetCompatibleWith []string          `json:"target_compatible_with"`
	Testonly             bool              `json:"testonly"`
	Toolchains           []string          `json:"toolchains"`
	TransitiveConfigs    []string          `json:"transitive_configs"`
	Visibility           []string          `json:"visibility"`
	Worker               WorkerLabel       `json:"worker"`
}

type WorkerFlagfile struct {
	Attr WorkerFlagfileAttr `json:"attr"`
}

type WorkerErrorResponse struct {
	Request   *worker_protocol.WorkRequest `json:"request"`
	Error     string                       `json:"error"`
	Arguments *WorkerRequestArguments      `json:"arguments"`
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
	response, err := self.handleRequest(request, args)
	if err != nil {
		self.writeResponse(self.errorResponse(request, args, err))
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
	return &WorkerRequestArguments{FlagfilePath: *flagfileFlag, Flagfile: flagfile}, nil
}

func (self *Worker) writeResponse(response *worker_protocol.WorkResponse) {
	_, err := protodelim.MarshalTo(self.stdout, response)
	if err != nil {
		panic(fmt.Sprintf("could not marshal response to stdout: %v", err))
	}
}

func (self *Worker) errorResponse(request *worker_protocol.WorkRequest, args *WorkerRequestArguments, err error) *worker_protocol.WorkResponse {
	response := WorkerErrorResponse{
		Request:   request,
		Arguments: args,
		Error:     err.Error(),
	}
	responseBytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("could not marshal response: %s", err))
	}
	return &worker_protocol.WorkResponse{
		ExitCode:  1,
		Output:    string(responseBytes),
		RequestId: request.RequestId,
	}
}

func (self *Worker) handleRequest(request *worker_protocol.WorkRequest, args *WorkerRequestArguments) (*worker_protocol.WorkResponse, error) {
	return nil, fmt.Errorf("request error")
}
