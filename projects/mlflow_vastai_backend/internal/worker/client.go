package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// maxRequestSize matches gRPC's default maximum inbound message size.
const maxRequestSize = 4 << 20

// call forwards one bounded protobuf request and writes one Job response.
func call(ctx context.Context, directory string, operation machine.Event, input io.Reader, output io.Writer) (err error) {
	request, err := io.ReadAll(io.LimitReader(input, maxRequestSize+1))
	if err != nil {
		return fmt.Errorf("read worker request: %w", err)
	}
	if len(request) > maxRequestSize {
		return errors.New("worker request exceeds 4 MiB")
	}
	connection, err := grpc.NewClient("unix://"+socketPath(directory), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create local worker client: %w", err)
	}
	defer func() { err = errors.Join(err, closeResource(connection, "worker gRPC connection")) }()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	response, err := invoke(ctx, jobpb.NewWorkerServiceClient(connection), operation, request)
	if err != nil {
		return fmt.Errorf("invoke local worker RPC: %w", err)
	}
	encoded, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("encode worker response: %w", err)
	}
	if n, err := output.Write(encoded); err != nil {
		return fmt.Errorf("write worker response: %w", err)
	} else if n != len(encoded) {
		return fmt.Errorf("write worker response: %w", io.ErrShortWrite)
	}
	return nil
}

// invoke translates the binary command protocol to its typed gRPC method.
func invoke(ctx context.Context, client jobpb.WorkerServiceClient, operation machine.Event, request []byte) (*jobpb.Job, error) {
	if operation == submitEvent {
		message := &jobpb.SubmitJobRequest{}
		if err := proto.Unmarshal(request, message); err != nil {
			return nil, fmt.Errorf("decode submission: %w", err)
		}
		job, err := client.SubmitJob(ctx, message)
		if err != nil {
			return nil, fmt.Errorf("submit to local service: %w", err)
		}
		return job, nil
	}
	reference := &jobpb.JobReference{}
	if err := proto.Unmarshal(request, reference); err != nil {
		return nil, fmt.Errorf("decode job reference: %w", err)
	}
	var job *jobpb.Job
	var err error
	switch operation {
	case statusEvent:
		job, err = client.GetJob(ctx, reference)
	case cancelEvent:
		job, err = client.CancelJob(ctx, reference)
	default:
		return nil, fmt.Errorf("unsupported worker operation %q", operation)
	}
	if err != nil {
		return nil, fmt.Errorf("send job reference to local service: %w", err)
	}
	return job, nil
}
