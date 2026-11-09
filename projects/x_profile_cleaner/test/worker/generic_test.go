package worker_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestGenericArchiveScript evaluates X-operation plans without cleanup declarations.
func TestGenericArchiveScript(t *testing.T) {
	socket, _ := startWorker(t, "2s", "5m")
	entries := map[string]string{
		"data/account.js":         `window.YTD.account.part0 = [{"account":{"accountId":"100"}}];`,
		"data/tweets.js":          `window.YTD.tweets.part0 = [];`,
		"data/like.js":            `window.YTD.like.part0 = [{"like":{"tweetId":"9007199254740993","fullText":"liked"}}];`,
		"data/like-part1.js":      `window.YTD.like.part1 = [{"like":{"tweetId":"9223372036854775807"}}];`,
		"data/following.js":       `window.YTD.following.part0 = [{"following":{"accountId":"200"}}];`,
		"data/broken.js":          `window.YTD.broken.part0 = []; execute();`,
		"data/inexact.js":         `window.YTD.inexact.part0 = [{"id":9007199254740993}];`,
		"data/direct-messages.js": "unrequested invalid data must remain unread",
	}
	var data bytes.Buffer
	w := zip.NewWriter(&data)
	for name, content := range entries {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(socket), "generic.zip"), data.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	rpc := workerpb.NewArchiveWorkerClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	script := `assert(os==nil and io==nil and require==nil and debug==nil)
for _, row in ipairs(archive.likes()) do
    assert(plan.unlike(row.post_id)==nil)
end
local rows = archive.following()
rows[1].account_id="changed"
assert(archive.following()[1].account_id=="200")
plan.unfollow("200")
plan.create_post("new")
plan.read_post("1")
local profile={description="updated"}
plan.update_profile(profile)
profile.description="changed"
assert(context.account_id=="100" and context.evaluated_at==1794225600)
plan.update_profile({description=""})
local options={forced=true}
plan.delete_post("1", options)
options.forced=false
plan.unrepost("2", nil)
`
	request := &workerpb.Evaluation{OperationId: "generic", AccountId: "100", EvaluatedAt: timestamppb.New(time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)), SourceName: "exports", ArchiveKey: "generic.zip", Script: script}
	if _, err := rpc.SubmitEvaluation(ctx, request); err != nil {
		t.Fatal(err)
	}
	ref := &workerpb.EvaluationRef{JobId: request.OperationId, AccountId: request.AccountId}
	job := waitStatus(t, ctx, rpc, ref)
	if job.State != workerpb.EvaluationState_SUCCEEDED {
		t.Fatalf("generic script rejected: %v", job.Failure)
	}
	want := &workerpb.Plan{Operations: []*workerpb.Operation{
		{Action: &workerpb.Operation_Unlike{Unlike: &workerpb.Unlike{PostId: "9007199254740993"}}},
		{Action: &workerpb.Operation_Unlike{Unlike: &workerpb.Unlike{PostId: "9223372036854775807"}}},
		{Action: &workerpb.Operation_Unfollow{Unfollow: &workerpb.Unfollow{UserId: "200"}}},
		{Action: &workerpb.Operation_CreatePost{CreatePost: &workerpb.CreatePost{Text: "new"}}},
		{Action: &workerpb.Operation_ReadPost{ReadPost: &workerpb.ReadPost{PostId: "1"}}},
		{Action: &workerpb.Operation_UpdateProfile{UpdateProfile: &workerpb.UpdateProfile{Description: proto.String("updated")}}},
		{Action: &workerpb.Operation_UpdateProfile{UpdateProfile: &workerpb.UpdateProfile{Description: proto.String("")}}},
		{Action: &workerpb.Operation_DeletePost{DeletePost: &workerpb.DeletePost{PostId: "1", Forced: true}}},
		{Action: &workerpb.Operation_Unrepost{Unrepost: &workerpb.Unrepost{PostId: "2"}}},
	}}
	if !proto.Equal(job.Plan, want) {
		t.Fatalf("typed plan differs: %v", job.Plan)
	}
	raw, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(job.Plan)
	if err != nil {
		t.Fatal(err)
	}
	writeAsyncArtifact(t, "typed-plan.json", json.RawMessage(raw))
	if _, err := rpc.DeleteEvaluation(ctx, ref); err != nil {
		t.Fatal(err)
	}
	for name, source := range map[string]string{
		"legacy-plan-result":       `return {operations={{unlike={post_id="1"}}}}`,
		"scalar-result":            `return "text"`,
		"nil-result":               `return nil`,
		"returned-after-operation": `plan.unlike("1"); return {operations={}}`,
		"large-operation":          `pcall(plan.create_post,string.rep("x",17*1024*1024))`,
		"aggregate-bytes":          `local text=string.rep("x",1024*1024);for i=1,17 do pcall(plan.create_post,text) end`,
		"missing-dataset":          `archive.read_raw("absent")`,
		"caught-read":              `pcall(archive.read_raw,"broken")`,
		"caught-projection":        `pcall(archive.read_raw,"inexact")`,
		"missing-profile-field":    `plan.update_profile({})`,
		"unknown-operation":        `pcall(function() plan.launch_process("anything") end)`,
		"bad-target":               `pcall(plan.unlike,"not-an-id")`,
		"numeric-target":           `pcall(plan.unlike,123)`,
		"unknown-option":           `plan.delete_post("1",{endpoint="arbitrary"})`,
		"duplicate-target-option":  `plan.delete_post("1",{post_id="2"})`,
		"invalid-flag":             `pcall(plan.delete_post,"1",{forced="true"})`,
		"invalid-options":          `plan.delete_post("1",true)`,
		"missing-target":           `pcall(plan.unlike)`,
		"extra-argument":           `pcall(plan.unlike,"1",true)`,
		"method-call":              `pcall(function() plan:unlike("1") end)`,
		"overwrite-method":         `pcall(function() plan.unlike=function() end end)`,
		"partial-operation":        `plan.unlike("1"); pcall(plan.unlike,false)`,
		"partial-script":           `plan.unlike("1"); error("failure")`,
		"returned-cleanup-result":  `local result=cleaner.plan(""); return result`,
		"caught-cleanup-arguments": `pcall(cleaner.plan,123)`,
		"caught-cleanup-error":     `pcall(cleaner.plan,"error('failure')")`,
		"multiple-cleanup-reports": `cleaner.plan(""); pcall(cleaner.plan,"")`,
		"operation-count":          `for i=1,50001 do pcall(plan.unlike,"1") end`,
		"invalid-text":             `pcall(plan.create_post,true)`,
	} {
		request.OperationId, request.Script = name, source
		if _, err := rpc.SubmitEvaluation(ctx, request); err != nil {
			t.Fatal(err)
		}
		ref.JobId = name
		job = waitStatus(t, ctx, rpc, ref)
		if job.State != workerpb.EvaluationState_FAILED || job.Plan != nil {
			t.Fatalf("%s published a result: %v", name, job.State)
		}
		if _, err := rpc.DeleteEvaluation(ctx, ref); err != nil {
			t.Fatal(err)
		}
	}
	request.OperationId, request.Script = "recovery", `assert(plan.unlike~=nil); return`
	if _, err := rpc.SubmitEvaluation(ctx, request); err != nil {
		t.Fatal(err)
	}
	ref.JobId = request.OperationId
	if result := waitStatus(t, ctx, rpc, ref); result.State != workerpb.EvaluationState_SUCCEEDED || len(result.Plan.GetOperations()) != 0 || result.Plan.GetCleanupReport() != nil {
		t.Fatalf("generic worker did not recover: %v", result.Failure)
	}
}
