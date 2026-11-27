package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"google.golang.org/protobuf/encoding/protodelim"

	"git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/worker_protocol"
)

// serve pairs each complete request with one response and terminates at a frame boundary.
func serve(ctx context.Context, input io.Reader, output io.Writer, persist bool) error {
	reader := bufio.NewReader(input)
	for {
		request := &worker_protocol.WorkRequest{}
		if err := protodelim.UnmarshalFrom(reader, request); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("decode shell work request: %w", err)
		}
		if !request.Cancel {
			response := perform(ctx, request)
			if _, err := protodelim.MarshalTo(output, response); err != nil {
				return fmt.Errorf("write shell response %d: %w", request.RequestId, err)
			}
		}
		if !persist {
			return nil
		}
	}
}
