package worker_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/cleanup"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestAsyncLifecycle exercises detached submission and explicit job management.
func TestAsyncLifecycle(t *testing.T) {
	socket, _ := startWorker(t, "2s", "500ms")
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := connection.Close(); err != nil {
			t.Error(err)
		}
	})
	rpc := workerpb.NewArchiveWorkerClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ref := &workerpb.EvaluationRef{JobId: "detached", AccountId: "100"}
	started := time.Now()
	ack, err := rpc.SubmitEvaluation(ctx, &workerpb.Evaluation{OperationId: ref.JobId, AccountId: ref.AccountId, EvaluatedAt: timestamppb.Now(), Script: "while true do end", SourceName: "exports", ArchiveKey: "empty.zip"})
	if err != nil {
		t.Fatal(err)
	}
	if ack.JobId != ref.JobId || time.Since(started) > time.Second {
		t.Fatalf("submission waited for evaluation: %+v", ack)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	connection, err = grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	rpc = workerpb.NewArchiveWorkerClient(connection)
	job, err := rpc.GetEvaluation(ctx, ref)
	if err != nil {
		t.Fatal(err)
	}
	if job.State != workerpb.EvaluationState_RUNNING || job.Plan != nil {
		t.Fatalf("running status: %+v", job)
	}
	foreign := &workerpb.EvaluationRef{JobId: ref.JobId, AccountId: "200"}
	if _, err := rpc.GetEvaluation(ctx, foreign); status.Code(err) != codes.NotFound {
		t.Fatalf("foreign read: %v", err)
	}
	if _, err := rpc.CancelEvaluation(ctx, foreign); status.Code(err) != codes.NotFound {
		t.Fatalf("foreign cancel: %v", err)
	}
	if _, err := rpc.DeleteEvaluation(ctx, ref); status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("delete running: %v", err)
	}
	if _, err := rpc.CancelEvaluation(ctx, ref); err != nil {
		t.Fatal(err)
	}
	job = waitStatus(t, ctx, rpc, ref)
	if job.State != workerpb.EvaluationState_CANCELLED || job.Plan != nil {
		t.Fatalf("cancelled status: %+v", job)
	}
	if _, err := rpc.CancelEvaluation(ctx, ref); err != nil {
		t.Fatal(err)
	}
	if _, err := rpc.DeleteEvaluation(ctx, ref); err != nil {
		t.Fatal(err)
	}
	if _, err := rpc.GetEvaluation(ctx, ref); status.Code(err) != codes.NotFound {
		t.Fatalf("deleted result: %v", err)
	}
	// Retention is fixed from completion, even while a caller polls repeatedly.
	ack, err = submitEmpty(ctx, rpc, "expiry", "100")
	if err != nil {
		t.Fatal(err)
	}
	ref = &workerpb.EvaluationRef{JobId: ack.JobId, AccountId: "100"}
	job = waitStatus(t, ctx, rpc, ref)
	if job.State != workerpb.EvaluationState_SUCCEEDED {
		t.Fatalf("empty policy result: %+v", job)
	}
	again, err := rpc.GetEvaluation(ctx, ref)
	if err != nil || !proto.Equal(again, job) {
		t.Fatalf("repeated result changed: %+v %v", again, err)
	}
	for {
		_, err := rpc.GetEvaluation(ctx, ref)
		if status.Code(err) == codes.NotFound {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		select {
		case <-ctx.Done():
			t.Fatal("polling prevented expiry")
		case <-time.After(20 * time.Millisecond):
		}
	}
	writeAsyncArtifact(t, "async-lifecycle.json", map[string]any{"job_id": "detached", "submission_detached": true, "cancelled": true, "result_expired": true, "repeatable_result": true})
}

// submitEmpty accepts a no-data policy to exercise job lifecycle independently.
func submitEmpty(ctx context.Context, rpc workerpb.ArchiveWorkerClient, id, account string) (*workerpb.SubmitEvaluationResponse, error) {
	ack, err := rpc.SubmitEvaluation(ctx, &workerpb.Evaluation{OperationId: id, AccountId: account, EvaluatedAt: timestamppb.Now(), SourceName: "exports", ArchiveKey: "empty.zip", Script: ""})
	if err != nil {
		return nil, fmt.Errorf("submit empty archive policy: %w", err)
	}
	return ack, nil
}

// TestRetainedJobLimit checks duplicate protection, bounded records, and deletion.
func TestRetainedJobLimit(t *testing.T) {
	socket, _ := startWorker(t, "2s", "20s")
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	rpc := workerpb.NewArchiveWorkerClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for index := range 16 {
		id := fmt.Sprint("retained-", index)
		if _, err := submitEmpty(ctx, rpc, id, "100"); err != nil {
			t.Fatal(err)
		}
		if job := waitStatus(t, ctx, rpc, &workerpb.EvaluationRef{JobId: id, AccountId: "100"}); job.State != workerpb.EvaluationState_SUCCEEDED {
			t.Fatalf("retained job: %+v", job)
		}
	}
	if _, err := submitEmpty(ctx, rpc, "retained-0", "100"); status.Code(err) != codes.AlreadyExists {
		t.Fatalf("duplicate ID: %v", err)
	}
	if _, err := submitEmpty(ctx, rpc, "overflow", "100"); status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("retained job limit: %v", err)
	}
	if _, err := rpc.DeleteEvaluation(ctx, &workerpb.EvaluationRef{JobId: "retained-0", AccountId: "200"}); status.Code(err) != codes.NotFound {
		t.Fatalf("foreign delete: %v", err)
	}
	if _, err := rpc.DeleteEvaluation(ctx, &workerpb.EvaluationRef{JobId: "retained-0", AccountId: "100"}); err != nil {
		t.Fatal(err)
	}
	if _, err := submitEmpty(ctx, rpc, "overflow", "100"); err != nil {
		t.Fatal(err)
	}
}

// TestRetainedResultBytes rejects excess result retention without losing old jobs.
func TestRetainedResultBytes(t *testing.T) {
	socket, _ := startWorker(t, "5s", "20s")
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16*1024*1024)))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	rpc := workerpb.NewArchiveWorkerClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	payload := strings.Repeat("a", 800*1024)
	for index := range 4 {
		id := fmt.Sprint("large-", index)
		items := []model.Item{}
		for record := range 16 {
			items = append(items, model.Item{ID: fmt.Sprint(record + 1), AuthorID: "100", Kind: model.ItemPost, CreatedAt: time.Now().UTC(), Text: payload})
		}
		key := id + ".zip"
		writeWorkerArchive(t, socket, key, "100", items)
		if _, err := rpc.SubmitEvaluation(ctx, &workerpb.Evaluation{OperationId: id, AccountId: "100", EvaluatedAt: timestamppb.Now(), Script: cleanup.Script("delete_post_when(function(p) return true end)"), SourceName: "exports", ArchiveKey: key}); err != nil {
			t.Fatal(err)
		}
		job := waitStatus(t, ctx, rpc, &workerpb.EvaluationRef{JobId: id, AccountId: "100"})
		if index == 2 {
			if job.State != workerpb.EvaluationState_FAILED || job.Plan != nil || job.Failure == nil || job.Failure.Code != workerpb.EvaluationFailure_RESULT_CAPACITY {
				t.Fatalf("result retention limit: state=%v failure=%v has_plan=%t", job.State, job.Failure, job.Plan != nil)
			}
			if _, err := rpc.DeleteEvaluation(ctx, &workerpb.EvaluationRef{JobId: "large-0", AccountId: "100"}); err != nil {
				t.Fatal(err)
			}
		} else if job.State != workerpb.EvaluationState_SUCCEEDED || len(job.Plan.Operations) != 16 {
			t.Fatalf("retained result state: %v", job.State)
		}
	}
	writeAsyncArtifact(t, "async-capacity.json", map[string]any{"retained_result_limit_bytes": 32 * 1024 * 1024, "overflow_failed_without_plan": true, "deletion_released_capacity": true})
}

// TestWorkerRestart forgets transient jobs without resuming their computation.
func TestWorkerRestart(t *testing.T) {
	socket, stop := startWorker(t, "2s", "20s")
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	rpc := workerpb.NewArchiveWorkerClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := submitEmpty(ctx, rpc, "restart", "100"); err != nil {
		t.Fatal(err)
	}
	ref := &workerpb.EvaluationRef{JobId: "restart", AccountId: "100"}
	_ = waitStatus(t, ctx, rpc, ref)
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	stop()
	startWorkerAt(t, socket, "2s", "20s")
	connection, err = grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	rpc = workerpb.NewArchiveWorkerClient(connection)
	if _, err := rpc.GetEvaluation(ctx, ref); status.Code(err) != codes.NotFound {
		t.Fatalf("restart retained execution: %v", err)
	}
}
