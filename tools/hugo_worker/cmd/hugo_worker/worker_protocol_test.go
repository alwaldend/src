package main

import (
	"bytes"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/worker_protocol"
	"google.golang.org/protobuf/encoding/protodelim"
)

// TestProtocolRoundTrip proves that a correctly framed WorkRequest (varint
// request_id, as produced by protodelim.MarshalTo) is answered by the worker.
// A manual probe that encoded request_id as a length-delimited field instead
// of a varint is rejected by protobuf-go and correctly produces no response.
func TestProtocolRoundTrip(t *testing.T) {
	var input bytes.Buffer
	_, err := protodelim.MarshalTo(&input, &worker_protocol.WorkRequest{
		Arguments: []string{"--flagfile=does-not-matter.json"},
		RequestId: 7,
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	var output bytes.Buffer
	protocol := NewWorkerProtocol(&output, &input)
	request, err := protocol.ReadRequest()
	if err != nil {
		t.Fatalf("read request: %v", err)
	}
	if request.RequestId != 7 {
		t.Fatalf("request id = %d, want 7", request.RequestId)
	}
	protocol.WriteResponse(&worker_protocol.WorkResponse{
		ExitCode:  0,
		Output:    "ok",
		RequestId: request.RequestId,
	})
	var response worker_protocol.WorkResponse
	err = protodelim.UnmarshalFrom(&output, &response)
	if err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if response.RequestId != 7 {
		t.Fatalf("response request id = %d, want 7", response.RequestId)
	}
	if response.ExitCode != 0 || response.Output != "ok" {
		t.Fatalf("unexpected response: exit=%d output=%q", response.ExitCode, response.Output)
	}
}

// TestProtocolRejectsLengthDelimitedID reproduces the manual probe failure:
// encoding field 3 (request_id, varint) as a length-delimited field makes
// protobuf-go fail to parse the request, so the worker must not respond.
func TestProtocolRejectsLengthDelimitedID(t *testing.T) {
	// The manual probe (request.bin) encoded request_id as field 3 with a
	// length-delimited wire type instead of a varint:
	//   22 <len> <message>  (the outer 22 is field 3, wire type 2)
	probe := []byte{
		0x22, 0x0a, 0x01, 0x07, 0x12, 0x1d,
		0x2d, 0x2d, 0x66, 0x6c, 0x61, 0x67, 0x66, 0x69, 0x6c, 0x65,
		0x3d, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x66, 0x6c, 0x61, 0x67,
		0x66, 0x69, 0x6c, 0x65, 0x2e, 0x6a, 0x73, 0x6f, 0x6e,
	}
	var output bytes.Buffer
	protocol := NewWorkerProtocol(&output, bytes.NewReader(probe))
	_, err := protocol.ReadRequest()
	if err == nil {
		t.Fatal("expected length-delimited request_id to be rejected")
	}
	if !strings.Contains(err.Error(), "cannot parse") &&
		!strings.Contains(err.Error(), "mismatching") {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Len() != 0 {
		t.Fatalf("worker must not respond when the request cannot be parsed; got %d bytes", output.Len())
	}
}
