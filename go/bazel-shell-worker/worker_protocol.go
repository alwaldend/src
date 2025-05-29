package main

import (
	"bufio"
	"fmt"
	"io"

	"google.golang.org/protobuf/encoding/protodelim"

	"git.alwaldend.com/src/proto/bazel-worker/worker_protocol"
)

type WorkerProtocol struct {
	stdout io.Writer
	stdin  protodelim.Reader
}

func NewWorkerProtocol(stdout io.Writer, stdin io.Reader) *WorkerProtocol {
	return &WorkerProtocol{stdout: stdout, stdin: bufio.NewReader(stdin)}
}

func (self *WorkerProtocol) ReadRequest() (*worker_protocol.WorkRequest, error) {
	request := &worker_protocol.WorkRequest{}
	err := protodelim.UnmarshalFrom(self.stdin, request)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal work request: %w", err)
	}
	return request, nil
}

func (self *WorkerProtocol) WriteResponse(response *worker_protocol.WorkResponse) {
	_, err := protodelim.MarshalTo(self.stdout, response)
	if err != nil {
		panic(fmt.Sprintf("could not marshal response to stdout: %v", err))
	}
}
