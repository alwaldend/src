package worker_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/cleanup"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	operationplan "git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/worker"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// evaluateClient waits only inside the E2E driver; the production API is detached.
func evaluateClient(ctx context.Context, client *worker.Client, request worker.Request) (*workerpb.Plan, error) {
	request.Script = cleanup.Script(request.Script)
	if _, err := client.Submit(ctx, request); err != nil {
		return nil, fmt.Errorf("submit fixture policy: %w", err)
	}
	for {
		job, err := client.Get(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("get fixture result: %w", err)
		}
		if job.State == workerpb.EvaluationState_RUNNING {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("wait for fixture result: %w", ctx.Err())
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}
		if err := client.Delete(ctx, request.AccountID, request.OperationID); err != nil {
			return nil, fmt.Errorf("release fixture result: %w", err)
		}
		if err := jobError(job.State, job.Failure); err != nil {
			return nil, fmt.Errorf("fixture policy failed: %w", err)
		}
		if job.Plan == nil {
			return nil, fmt.Errorf("fixture completed without plan")
		}
		err = cleanup.Validate(job.Plan)
		if err != nil {
			return nil, fmt.Errorf("validate fixture cleanup plan: %w", err)
		}
		return job.Plan, nil
	}
}

// jobError maps terminal job failures to the existing E2E rejection assertions.
func jobError(state workerpb.EvaluationState, failure *workerpb.EvaluationFailure) error {
	if state == workerpb.EvaluationState_SUCCEEDED {
		return nil
	}
	if state == workerpb.EvaluationState_CANCELLED {
		return status.Error(codes.Canceled, "job cancelled")
	}
	if failure == nil {
		return status.Error(codes.Internal, "missing job failure")
	}
	code := codes.Internal
	switch failure.Code {
	case workerpb.EvaluationFailure_INVALID_SCRIPT, workerpb.EvaluationFailure_INVALID_ARCHIVE:
		code = codes.InvalidArgument
	case workerpb.EvaluationFailure_DEADLINE_EXCEEDED:
		code = codes.DeadlineExceeded
	case workerpb.EvaluationFailure_OUTPUT_LIMIT, workerpb.EvaluationFailure_RESULT_CAPACITY:
		code = codes.ResourceExhausted
	}
	return status.Error(code, failure.Message)
}

// fixtureOperation constructs untrusted fixture payloads independently of production codecs.
func fixtureOperation(kind, target string, forced bool) *workerpb.Operation {
	if kind == "delete_post" {
		return &workerpb.Operation{Action: &workerpb.Operation_DeletePost{DeletePost: &workerpb.DeletePost{PostId: target, Forced: forced}}}
	}
	return &workerpb.Operation{Action: &workerpb.Operation_Unrepost{Unrepost: &workerpb.Unrepost{PostId: target, Forced: forced}}}
}

// TestWorkerProcess verifies archive-backed policy semantics and bounded execution.
func TestWorkerProcess(t *testing.T) {
	socket, _ := startWorker(t, "2s", "5m")
	client, err := worker.Dial(socket)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	items := []model.Item{
		{ID: "101", EditHistoryIDs: []string{"1", "101"}, AuthorID: "100", Kind: model.ItemPost, CreatedAt: now.Add(-48 * time.Hour), Text: "root"},
		{ID: "2", AuthorID: "100", Kind: model.ItemPost, CreatedAt: now.Add(-24 * time.Hour), Text: "reply", References: []model.Reference{{ID: "1", Kind: model.ReferenceReply}}},
		{ID: "3", AuthorID: "100", Kind: model.ItemPost, CreatedAt: now.Add(-time.Hour), Text: "recent", Metrics: map[string]int64{"like_count": 7, "reply_count": 0}},
		{ID: "4", AuthorID: "100", Kind: model.ItemRepost, SourceID: "900", CreatedAt: now.Add(-48 * time.Hour), Text: "repost"},
	}
	empty, username := "", "author"
	items[2].Language = &empty
	items[2].Attachments = &model.Attachments{MediaKeys: []string{"123"}, MediaSourceIDs: []string{"99"}}
	items[2].Entities = &model.TextEntities{Hashtags: []model.TextEntity{}, Mentions: []model.TextEntity{{Start: 0, End: 1, Username: &username}}}
	writeWorkerArchive(t, socket, "cleanup.zip", "100", items)
	request := worker.Request{OperationID: "archive-policy", AccountID: "100", EvaluatedAt: now, SourceName: "exports", ArchiveKey: "cleanup.zip", Script: `assert(os==nil and io==nil and package==nil and require==nil)
local root=get_post("101")
assert(root.language==nil and #root.attachments.media_keys==0 and #root.entities.hashtags==0)
local p=get_post("3")
assert(p.language=="" and p.metrics.like_count==7 and p.metrics.reply_count==0 and p.metrics.quote_count==nil)
assert(p.attachments.media_keys[1]=="123" and p.attachments.media_source_ids[1]=="99" and #p.attachments.poll_ids==0)
assert(p.entities.mentions[1].username=="author" and p.entities.mentions[1].tag==nil)
p.entities.mentions[1].username="changed"
assert(get_post("3").entities.mentions[1].username=="author")
delete_post_when(function(p) return p.age_days>=1 end)
preserve_thread("1"); force_delete_post("2"); delete_repost("4")`}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	plan, err := evaluateClient(ctx, client, request)
	if err != nil {
		t.Fatal(err)
	}
	want := []*workerpb.Operation{fixtureOperation("unrepost", "900", false), fixtureOperation("delete_post", "2", true)}
	if !proto.Equal(&workerpb.Plan{Operations: plan.Operations}, &workerpb.Plan{Operations: want}) {
		t.Fatalf("archive actions: %+v", plan.Operations)
	}
	if len(plan.CleanupReport.Decisions) != 4 || plan.CleanupReport.Decisions[0].Outcome != "preserved" {
		t.Fatalf("archive decisions: %+v", plan.CleanupReport.Decisions)
	}
	request.Script = `set_time_window(days_ago(0.5),nil);delete_post_when(function(p) return true end)`
	if result, err := evaluateClient(ctx, client, request); err != nil || !proto.Equal(&workerpb.Plan{Operations: result.Operations}, &workerpb.Plan{Operations: []*workerpb.Operation{fixtureOperation("delete_post", "3", false)}}) {
		t.Fatalf("archive window: %+v %v", result.Operations, err)
	}
	for _, source := range []string{"this is not Lua", "while true do end"} {
		request.Script = source
		if _, err := evaluateClient(ctx, client, request); err == nil {
			t.Fatal("unsafe policy succeeded")
		}
	}
	request.Script = ""
	if result, err := evaluateClient(ctx, client, request); err != nil || len(result.Operations) != 0 {
		t.Fatalf("worker did not recover: %+v %v", result, err)
	}
	request.ArchiveKey = "../cleanup.zip"
	if _, err := client.Submit(ctx, request); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("source traversal: %v", err)
	}
	request.ArchiveKey = "cleanup.zip"
	request.SourceName = "unknown"
	if _, err := client.Submit(ctx, request); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("unknown source: %v", err)
	}
	items[3].SourceID = ""
	writeWorkerArchive(t, socket, "unresolved.zip", "100", items)
	request.SourceName = "exports"
	request.ArchiveKey = "unresolved.zip"
	request.Script = `delete_post_when(function(p) return false end)`
	if _, err := evaluateClient(ctx, client, request); err != nil {
		t.Fatalf("own posts with unresolved repost: %v", err)
	}
	request.Script = `delete_repost_when(function(p) return true end)`
	if _, err := evaluateClient(ctx, client, request); err == nil {
		t.Fatal("unknown repost source was guessed")
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		data, err := (protojson.MarshalOptions{UseProtoNames: true, Indent: "  ", EmitDefaultValues: true}).Marshal(plan)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "worker-plan.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// forgedWorker deliberately returns invalid worker output to the backend adapter.
type forgedWorker struct {
	// UnsafeArchiveWorkerServer requires implementation of the complete current service.
	workerpb.UnsafeArchiveWorkerServer
	// mu synchronizes response replacement with the RPC handler.
	mu sync.Mutex
	// response is replaced between sequential calls.
	response *workerpb.EvaluationStatus
}

// GetEvaluation returns the configured untrusted status without computation.
func (s *forgedWorker) GetEvaluation(context.Context, *workerpb.EvaluationRef) (*workerpb.EvaluationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.response, nil
}

// SubmitEvaluation is unsupported by the result-validation fixture.
func (s *forgedWorker) SubmitEvaluation(context.Context, *workerpb.Evaluation) (*workerpb.SubmitEvaluationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "fixture supports status reads only")
}

// CancelEvaluation is unsupported by the result-validation fixture.
func (s *forgedWorker) CancelEvaluation(context.Context, *workerpb.EvaluationRef) (*workerpb.CancelEvaluationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "fixture supports status reads only")
}

// DeleteEvaluation is unsupported by the result-validation fixture.
func (s *forgedWorker) DeleteEvaluation(context.Context, *workerpb.EvaluationRef) (*workerpb.DeleteEvaluationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "fixture supports status reads only")
}

// TestUntrustedOutput checks that worker transport success cannot authorize a plan.
func TestUntrustedOutput(t *testing.T) {
	directory, err := os.MkdirTemp("", "cleaner-output-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(directory); err != nil {
			t.Error(err)
		}
	})
	socket := filepath.Join(directory, "worker.sock")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	service := &forgedWorker{}
	workerpb.RegisterArchiveWorkerServer(server, service)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	client, err := worker.Dial(socket)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	now := time.Now().UTC()
	post := model.Item{ID: "123", AuthorID: "100", Kind: model.ItemPost, CreatedAt: now, Text: "post"}
	repost := model.Item{ID: "124", AuthorID: "100", Kind: model.ItemRepost, SourceID: "900", CreatedAt: now, Text: "repost"}
	request := worker.Request{OperationID: "output-check", AccountID: "100", EvaluatedAt: now, SourceName: "exports", ArchiveKey: "cleanup.zip"}
	valid := &workerpb.Plan{CleanupReport: &workerpb.CleanupReport{}}
	for _, test := range []string{"operation", "account", "time", "source", "archive key", "envelope time", "unknown plan field", "unknown operation", "missing operation payload", "invalid profile update", "invalid UTF-8", "target", "missing plan", "size", "running with plan", "cancelled with plan", "failed without error", "missing expiry", "zero expiry", "unknown state", "unknown failure", "oversized failure", "invalid reason post", "invalid time post", "preserved post", "unchanged post", "preserved repost", "unchanged repost", "unexplained post", "unexplained repost", "missing action for post", "missing action for repost", "unforced post", "unforced repost", "forced post", "forced repost", "wrong repost target", "valid post", "valid forced post", "valid repost", "valid forced repost"} {
		t.Run(test, func(t *testing.T) {
			plan := proto.Clone(valid).(*workerpb.Plan)
			accepted := strings.HasPrefix(test, "valid ")
			response := &workerpb.EvaluationStatus{JobId: request.OperationID, AccountId: request.AccountID, State: workerpb.EvaluationState_SUCCEEDED, ExpiresAt: timestamppb.New(time.Now().Add(time.Minute)), SourceName: request.SourceName, ArchiveKey: request.ArchiveKey, EvaluatedAt: timestamppb.New(request.EvaluatedAt)}
			switch test {
			case "operation":
				response.JobId = "other"
			case "account":
				response.AccountId = "200"
			case "time":
				response.EvaluatedAt = timestamppb.New(now.Add(time.Second))
			case "target":
				plan.Operations = []*workerpb.Operation{fixtureOperation("delete_post", "777", false)}
			default:
				if strings.HasSuffix(test, "post") || test == "wrong repost target" {
					item := post
					forced := strings.HasPrefix(test, "forced ") || strings.HasPrefix(test, "valid forced ")
					action := fixtureOperation("delete_post", post.ID, forced)
					if strings.Contains(test, "repost") {
						item = repost
						action = fixtureOperation("unrepost", repost.SourceID, forced)
					}
					decision := &workerpb.CleanupDecision{ItemId: item.ID, SourceId: item.SourceID, Kind: string(item.Kind), CreatedAt: timestamppb.New(item.CreatedAt), Text: item.Text, Outcome: "delete"}
					if test == "invalid reason post" {
						decision.Reasons = []*workerpb.CleanupReason{{RuleId: 0, Operation: "delete_post"}}
					}
					if test == "invalid time post" {
						decision.CreatedAt = timestamppb.New(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
					}
					switch {
					case strings.HasPrefix(test, "preserved "):
						decision.Outcome = "preserved"
					case strings.HasPrefix(test, "unchanged "):
						decision.Outcome = "unchanged"
					case strings.HasPrefix(test, "unforced "):
						decision.Outcome = "force_delete"
					case strings.HasPrefix(test, "forced "):
						// forced payload was constructed above
					case strings.HasPrefix(test, "valid forced "):
						decision.Outcome = "force_delete"
					}
					plan.Operations, plan.CleanupReport.Decisions = []*workerpb.Operation{action}, []*workerpb.CleanupDecision{decision}
					if strings.HasPrefix(test, "unexplained ") {
						plan.CleanupReport.Decisions = nil
					}
					if strings.HasPrefix(test, "missing ") {
						plan.Operations = []*workerpb.Operation{}
					}
					if test == "wrong repost target" {
						plan.Operations[0].GetUnrepost().PostId = repost.ID
					}
				}
			}
			response.Plan = plan
			switch test {
			case "source":
				response.SourceName = "other"
			case "archive key":
				response.ArchiveKey = "other.zip"
			case "envelope time":
				response.EvaluatedAt = timestamppb.New(now.Add(time.Second))
			case "unknown plan field":
				response.Plan.ProtoReflect().SetUnknown([]byte{0xf8, 0x07, 0x01})
			case "unknown operation":
				operation := &workerpb.Operation{}
				operation.ProtoReflect().SetUnknown([]byte{0xf8, 0x07, 0x01})
				response.Plan.Operations = []*workerpb.Operation{operation}
			case "missing operation payload":
				response.Plan.Operations = []*workerpb.Operation{{}}
			case "invalid profile update":
				response.Plan.Operations = []*workerpb.Operation{{Action: &workerpb.Operation_UpdateProfile{UpdateProfile: &workerpb.UpdateProfile{}}}}
			case "invalid UTF-8":
				response.Plan.Operations = []*workerpb.Operation{{Action: &workerpb.Operation_CreatePost{CreatePost: &workerpb.CreatePost{Text: string([]byte{0xff})}}}}

			}
			if test == "missing plan" {
				response.Plan = nil
			}
			if test == "size" {
				response.Plan.CleanupReport.Decisions = []*workerpb.CleanupDecision{{Text: strings.Repeat("a", 16*1024*1024+1)}}
			}
			switch test {
			case "running with plan":
				response.State = workerpb.EvaluationState_RUNNING
				response.ExpiresAt = nil
			case "cancelled with plan":
				response.State = workerpb.EvaluationState_CANCELLED
			case "failed without error":
				response.State = workerpb.EvaluationState_FAILED
				response.Plan = nil
			case "missing expiry":
				response.ExpiresAt = nil
			case "zero expiry":
				response.ExpiresAt = timestamppb.New(time.Time{})
			case "unknown state":
				response.State = workerpb.EvaluationState(99)
			case "unknown failure":
				response.State = workerpb.EvaluationState_FAILED
				response.Plan = nil
				response.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_Code(99)}
			case "oversized failure":
				response.State = workerpb.EvaluationState_FAILED
				response.Plan = nil
				response.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_INVALID_SCRIPT, Message: strings.Repeat("a", 257)}
			}
			service.mu.Lock()
			service.response = response
			service.mu.Unlock()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			job, err := client.Get(ctx, request)
			var result *workerpb.Plan
			if err == nil && job.Plan != nil {
				result = job.Plan
				err = cleanup.Validate(result)
			}
			if accepted {
				if err != nil || !proto.Equal(result, plan) {
					t.Fatalf("valid worker plan rejected: %+v %v", result, err)
				}
			} else if err == nil {
				t.Fatal("untrusted output accepted")
			}
		})
	}
	service.mu.Lock()
	service.response = &workerpb.EvaluationStatus{JobId: request.OperationID, AccountId: request.AccountID, SourceName: request.SourceName, ArchiveKey: request.ArchiveKey, EvaluatedAt: timestamppb.New(now), State: workerpb.EvaluationState_SUCCEEDED, ExpiresAt: timestamppb.New(time.Now().Add(time.Minute)), Plan: &workerpb.Plan{Operations: []*workerpb.Operation{{Action: &workerpb.Operation_UpdateProfile{UpdateProfile: &workerpb.UpdateProfile{Description: proto.String("changed")}}}}}}
	service.mu.Unlock()
	job, err := client.Get(context.Background(), request)
	if err != nil {
		t.Fatalf("typed operation plan rejected by transport: %v", err)
	}
	if err := cleanup.Validate(job.Plan); err == nil {
		t.Fatal("cleanup backend accepted an unsupported operation plan")
	}
}

// TestStableReferencesAndDiagnostics checks user-visible policy behavior through a worker.
func TestStableReferencesAndDiagnostics(t *testing.T) {
	socket, _ := startWorker(t, "2s", "5m")
	client, err := worker.Dial(socket)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	now := time.Now().UTC().Truncate(time.Second)
	writeWorkerArchive(t, socket, "references.zip", "100", []model.Item{{ID: "2", AuthorID: "100", Kind: model.ItemPost, CreatedAt: now, Text: "private-account-text", References: []model.Reference{{ID: "1", Kind: model.ReferenceReply}, {ID: "3", Kind: model.ReferenceQuote}}}})
	request := worker.Request{OperationID: "stable-reference", AccountID: "100", EvaluatedAt: now, SourceName: "exports", ArchiveKey: "references.zip", Script: `local r=get_post("2").references; assert(#r==2 and r[1].kind=="replied_to" and r[2].kind=="quoted"); delete_post("2")`}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	for range 64 {
		plan, err := evaluateClient(ctx, client, request)
		if err != nil || len(plan.Operations) != 1 || plan.Operations[0].GetDeletePost().PostId != "2" {
			t.Fatalf("unstable reference order: %v %v", plan.Operations, err)
		}
	}
	for _, test := range []struct {
		// source supplies malformed or failing Lua without disclosing its error object.
		source string
		// stage is the safe diagnostic category.
		stage string
		// line is the source location expected without source contents.
		line string
	}{{"\nlocal 123 = 1", "syntax", "line 2"}, {"delete_post_when(function(p)\n error('PRIVATE_DIAGNOSTIC_MARKER')\nend)", "evaluation", "line 2"}} {
		request.Script = test.source
		_, err := evaluateClient(ctx, client, request)
		if err == nil || !strings.Contains(err.Error(), test.stage) || !strings.Contains(err.Error(), test.line) {
			t.Fatalf("missing actionable diagnostic: %v", err)
		}
		if strings.Contains(err.Error(), "PRIVATE_DIAGNOSTIC_MARKER") || strings.Contains(err.Error(), "private-account-text") {
			t.Fatal("diagnostic disclosed policy or account contents")
		}
	}
	writeAsyncArtifact(t, "policy-diagnostics.json", map[string]any{"stable_reference_runs": 64, "syntax_line": 2, "callback_line": 2, "content_disclosed": false})
}

// cleanupWorker returns a complete result but rejects transient result deletion.
type cleanupWorker struct {
	// UnsafeArchiveWorkerServer requires all current job RPCs.
	workerpb.UnsafeArchiveWorkerServer
	// mu guards the submitted request and cleanup observation.
	mu sync.Mutex
	// request is the exact operation accepted from the CLI process.
	request *workerpb.Evaluation
	// deletes counts attempted terminal-result cleanup calls.
	deletes int
}

// SubmitEvaluation captures the backend-generated bindings.
func (s *cleanupWorker) SubmitEvaluation(_ context.Context, r *workerpb.Evaluation) (*workerpb.SubmitEvaluationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.request = r
	return &workerpb.SubmitEvaluationResponse{JobId: r.OperationId}, nil
}

// GetEvaluation returns a complete, internally consistent empty plan.
func (s *cleanupWorker) GetEvaluation(context.Context, *workerpb.EvaluationRef) (*workerpb.EvaluationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.request
	return &workerpb.EvaluationStatus{JobId: r.OperationId, AccountId: r.AccountId, SourceName: r.SourceName, ArchiveKey: r.ArchiveKey, State: workerpb.EvaluationState_SUCCEEDED, ExpiresAt: timestamppb.New(time.Now().Add(time.Minute)), EvaluatedAt: r.EvaluatedAt, Plan: &workerpb.Plan{CleanupReport: &workerpb.CleanupReport{}}}, nil
}

// CancelEvaluation acknowledges an unnecessary cancellation for fixture completeness.
func (s *cleanupWorker) CancelEvaluation(context.Context, *workerpb.EvaluationRef) (*workerpb.CancelEvaluationResponse, error) {
	return &workerpb.CancelEvaluationResponse{}, nil
}

// DeleteEvaluation simulates loss of availability after the result was delivered.
func (s *cleanupWorker) DeleteEvaluation(context.Context, *workerpb.EvaluationRef) (*workerpb.DeleteEvaluationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deletes++
	return nil, status.Error(codes.Unavailable, "synthetic cleanup failure")
}

// TestCLIResultCleanupFailure preserves a validated result despite unavailable cleanup.
func TestCLIResultCleanupFailure(t *testing.T) {
	directory, err := os.MkdirTemp("", "cleaner-cleanup-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(directory); err != nil {
			t.Error(err)
		}
	})
	socket := filepath.Join(directory, "w.sock")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	service := &cleanupWorker{}
	server := grpc.NewServer()
	workerpb.RegisterArchiveWorkerServer(server, service)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	binary, err := runfiles.Rlocation(*cleanerBinary)
	if err != nil {
		t.Fatal(err)
	}
	policy := filepath.Join(directory, "cleanup.lua")
	if err := os.WriteFile(policy, []byte("-- empty policy"), 0o600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, "plan", "--config", policy, "--worker-socket", socket, "--archive-source", "exports", "--archive-key", "export.zip", "--account", "100")
	command.Env = []string{"PATH=/usr/bin:/bin"}
	var output, diagnostics strings.Builder
	command.Stdout = &output
	command.Stderr = &diagnostics
	if err := command.Run(); err != nil {
		t.Fatalf("cleanup discarded successful result: %v %s", err, diagnostics.String())
	}
	var plan operationplan.Document
	if err := json.Unmarshal([]byte(output.String()), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.AccountID != "100" || len(plan.Plan.Operations) != 0 {
		t.Fatalf("unexpected cleanup plan: %+v", plan)
	}
	service.mu.Lock()
	deletes := service.deletes
	service.mu.Unlock()
	if deletes != 1 {
		t.Fatalf("terminal cleanup attempts: %d", deletes)
	}
	writeAsyncArtifact(t, "cleanup-result.json", map[string]any{"plan_published": true, "delete_attempts": deletes, "cleanup_failure_fatal": false})
}
