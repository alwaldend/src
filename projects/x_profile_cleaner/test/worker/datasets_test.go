package worker_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/script"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// datasetEntries supplies standard multipart exports with nullable optional fields.
func datasetEntries() map[string]string {
	return map[string]string{
		"data/account.js":         `window.YTD.account.part0 = [{"account":{"accountId":"100","username":"fixture","accountDisplayName":null}}];`,
		"data/tweets.js":          `window.YTD.tweets.part0 = [{"tweet":{"id_str":"101","created_at":"Mon Nov 09 12:00:00 +0000 2020","full_text":"root","favorite_count":"9007199254740993","reply_count":"0","entities":{"hashtags":[{"indices":[0,4],"text":"root"}]},"extended_entities":{"media":[{"id_str":"123"}]}}}];`,
		"data/like.js":            `window.YTD.like.part0 = [{"like":{"tweetId":"9007199254740993","fullText":"liked","expandedUrl":null}}];`,
		"data/like-part1.js":      `window.YTD.like.part1 = [{"like":{"tweetId":"9223372036854775807"}}];`,
		"data/following.js":       `window.YTD.following.part0 = [{"following":{"accountId":"200","userLink":null}}];`,
		"data/following-part1.js": `window.YTD.following.part1 = [{"following":{"accountId":"9007199254740993","userLink":"https://x.com/fixture"}}];`,
		"data/custom.js":          `window.YTD.custom.part0 = [{"value":null}];`,
	}
}

// datasetJob evaluates a real source-owned ZIP through the worker process.
func datasetJob(t *testing.T, rpc workerpb.ArchiveWorkerClient, socket, id, script string, changes map[string]string, expected workerpb.EvaluationState) *workerpb.Plan {
	t.Helper()
	entries := datasetEntries()
	for name, value := range changes {
		entries[name] = value
	}
	data := datasetZIP(t, entries)
	key := id + ".zip"
	if err := os.WriteFile(filepath.Join(filepath.Dir(socket), key), data, 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request := &workerpb.Evaluation{OperationId: id, AccountId: "100", SourceName: "exports", ArchiveKey: key, EvaluatedAt: timestamppb.New(time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)), Script: script}
	if _, err := rpc.SubmitEvaluation(ctx, request); err != nil {
		t.Fatal(err)
	}
	ref := &workerpb.EvaluationRef{JobId: id, AccountId: "100"}
	job := waitStatus(t, ctx, rpc, ref)
	if job.State != expected {
		t.Fatalf("job %s: state %v, want %v; failure %v", id, job.State, expected, job.Failure)
	}
	if expected != workerpb.EvaluationState_SUCCEEDED && job.Plan != nil {
		t.Fatal("failed typed read published a plan")
	}
	if _, err := rpc.DeleteEvaluation(ctx, ref); err != nil {
		t.Fatal(err)
	}
	return job.Plan
}

// TestPolicyEvaluatorOptions verifies the public options with real archive/Lua evaluation.
// Failure cases are ignored option deadlines, extended/cancelled parent contexts,
// negative timeouts and interpreter state leaking between uses of the same evaluator.
func TestPolicyEvaluatorOptions(t *testing.T) {
	data := datasetZIP(t, datasetEntries())
	open := func() *archive.Export {
		value, err := archive.Open(context.Background(), bytes.NewReader(data), int64(len(data)))
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	parent, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	evaluator := script.Evaluator{Options: script.Options{Timeout: 20 * time.Millisecond}}
	if _, err := evaluator.Evaluate(parent, open(), `while true do end`, now); !errors.Is(err, context.DeadlineExceeded) || parent.Err() != nil {
		t.Fatalf("option deadline not enforced: %v", err)
	}
	short, stop := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer stop()
	if _, err := (script.Evaluator{Options: script.Options{Timeout: time.Hour}}).Evaluate(short, open(), `while true do end`, now); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("parent deadline not preserved: %v", err)
	}
	cancelled, stopCancelled := context.WithCancel(context.Background())
	stopCancelled()
	if _, err := evaluator.Evaluate(cancelled, open(), `return`, now); !errors.Is(err, context.Canceled) {
		t.Fatalf("parent cancellation lost: %v", err)
	}
	if _, err := (script.Evaluator{Options: script.Options{Timeout: -time.Second}}).Evaluate(parent, open(), `return`, now); err == nil {
		t.Fatal("negative timeout accepted")
	}
	shared := script.Evaluator{}
	for _, source := range []string{`shared=42;return`, `assert(shared==nil);return`} {
		if _, err := shared.Evaluate(parent, open(), source, now); err != nil {
			t.Fatal(err)
		}
	}
	writeAsyncArtifact(t, "evaluator-options.json", map[string]any{"option_timeout_enforced": true, "parent_deadline_preserved": true, "parent_cancellation_preserved": true, "negative_timeout_rejected": true, "state_isolated": true})
}

// TestTypedArchiveDatasets exercises validated records, raw escape hatch and sticky failures.
// Failure cases include absent/null/wrong-type identifiers or envelopes, malformed optional
// fields, duplicate keys, invalid tweets, swallowed read errors and mutable projections.
func TestTypedArchiveDatasets(t *testing.T) {
	socket, _ := startWorker(t, "2s", "5m")
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	rpc := workerpb.NewArchiveWorkerClient(connection)
	t.Run("validated-records", func(t *testing.T) {
		plan := datasetJob(t, rpc, socket, "typed-records", `
local account = archive.account()
assert(account.id=="100" and account.username=="fixture" and account.display_name==nil)
local following = archive.following()
assert(#following==2 and following[1].account_id=="200" and following[1].user_link==nil)
following[1].account_id="changed"
assert(archive.following()[1].account_id=="200")
local likes=archive.likes()
assert(#likes==2 and likes[1].text=="liked" and likes[1].url==nil)
local tweets=archive.tweets()
assert(#tweets==1 and tweets[1].id=="101" and tweets[1].text=="root")
assert(tweets[1].author_id=="100" and tweets[1].kind=="POST")
assert(tweets[1].created_at.seconds==1604923200 and tweets[1].created_at.nanos==0)
assert(tweets[1].metrics.like_count=="9007199254740993" and tweets[1].metrics.reply_count==0 and tweets[1].metrics.quote_count==nil)
assert(tweets[1].attachments.media_keys[1]=="123" and #tweets[1].attachments.poll_ids==0)
assert(tweets[1].entities.hashtags[1].tag=="root" and tweets[1].entities.hashtags[1].start==0)
tweets[1].entities.hashtags[1].tag="changed"
assert(archive.tweets()[1].entities.hashtags[1].tag=="root")
tweets[1].text="changed"
assert(archive.tweets()[1].text=="root")
assert(archive.read_raw("custom")[1].value==json.null)
plan.unfollow(archive.following()[2].account_id)
plan.unlike(likes[1].post_id)
plan.unlike(likes[2].post_id)
plan.read_post(tweets[1].id)
`, nil, workerpb.EvaluationState_SUCCEEDED)
		want := &workerpb.Plan{Operations: []*workerpb.Operation{
			{Action: &workerpb.Operation_Unfollow{Unfollow: &workerpb.Unfollow{UserId: "9007199254740993"}}},
			{Action: &workerpb.Operation_Unlike{Unlike: &workerpb.Unlike{PostId: "9007199254740993"}}},
			{Action: &workerpb.Operation_Unlike{Unlike: &workerpb.Unlike{PostId: "9223372036854775807"}}},
			{Action: &workerpb.Operation_ReadPost{ReadPost: &workerpb.ReadPost{PostId: "101"}}},
		}}
		if !proto.Equal(plan, want) {
			t.Fatalf("typed dataset plan differs: %v", plan)
		}
		data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(plan)
		if err != nil {
			t.Fatal(err)
		}
		writeAsyncArtifact(t, "typed-datasets-plan.json", json.RawMessage(data))
	})
	for name, row := range map[string]string{
		"missing-id":     `{"following":{}}`,
		"null-id":        `{"following":{"accountId":null}}`,
		"numeric-id":     `{"following":{"accountId":200}}`,
		"invalid-id":     `{"following":{"accountId":"invalid"}}`,
		"null-envelope":  `{"following":null}`,
		"wrong-envelope": `{"user":{"accountId":"200"}}`,
		"wrong-optional": `{"following":{"accountId":"200","userLink":12}}`,
		"duplicate-id":   `{"following":{"accountId":"200","accountId":"201"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			datasetJob(t, rpc, socket, name, `pcall(archive.following); return`, map[string]string{"data/following.js": "window.YTD.following.part0 = [" + row + "];"}, workerpb.EvaluationState_FAILED)
		})
	}
	t.Run("bad-like", func(t *testing.T) {
		datasetJob(t, rpc, socket, "bad-like", `pcall(archive.likes); return`, map[string]string{"data/like.js": `window.YTD.like.part0 = [{"like":{"tweetId":"0"}}];`}, workerpb.EvaluationState_FAILED)
	})
	t.Run("bad-tweet", func(t *testing.T) {
		datasetJob(t, rpc, socket, "bad-tweet", `pcall(archive.tweets); return`, map[string]string{"data/tweets.js": `window.YTD.tweets.part0 = [{"tweet":{"id_str":"101","created_at":"bad","full_text":"root"}}];`}, workerpb.EvaluationState_FAILED)
	})
	t.Run("bad-account", func(t *testing.T) {
		datasetJob(t, rpc, socket, "bad-account", `pcall(archive.account); return`, map[string]string{"data/account.js": `window.YTD.account.part0 = [{"account":{"accountId":"100","username":12}}];`}, workerpb.EvaluationState_FAILED)
	})
	for _, counter := range []string{"favorite_count", "retweet_count", "reply_count", "quote_count", "bookmark_count", "impression_count"} {
		t.Run("null-"+counter, func(t *testing.T) {
			tweet := fmt.Sprintf(`window.YTD.tweets.part0 = [{"tweet":{"id_str":"101","created_at":"Mon Nov 09 12:00:00 +0000 2020","full_text":"root","%s":null}}];`, counter)
			plan := datasetJob(t, rpc, socket, "null-"+counter, `
local tweet=archive.tweets()[1]
for _, metric in ipairs({"like_count","repost_count","reply_count","quote_count","bookmark_count","impression_count"}) do
    assert(tweet.metrics[metric]==nil)
end
cleaner.plan('delete_post_when(function(post) return post.metrics.like_count==nil end)')
`, map[string]string{"data/tweets.js": tweet}, workerpb.EvaluationState_SUCCEEDED)
			if len(plan.Operations) != 1 || plan.Operations[0].GetDeletePost().GetPostId() != "101" {
				t.Fatalf("null counter changed cleanup selection: %v", plan)
			}
			writeAsyncArtifact(t, "null-"+counter+".json", map[string]any{"counter": counter, "absent_metrics": true, "deleted_post": "101"})
		})
	}
	for name, value := range map[string]string{"negative": "-1", "fractional": "1.5", "boolean": "true", "overflow": "9223372036854775808"} {
		t.Run("invalid-counter-"+name, func(t *testing.T) {
			tweet := fmt.Sprintf(`window.YTD.tweets.part0 = [{"tweet":{"id_str":"101","created_at":"Mon Nov 09 12:00:00 +0000 2020","full_text":"root","reply_count":%s}}];`, value)
			datasetJob(t, rpc, socket, "invalid-counter-"+name, `pcall(archive.tweets); return`, map[string]string{"data/tweets.js": tweet}, workerpb.EvaluationState_FAILED)
		})
	}
	t.Run("explicit-raw", func(t *testing.T) {
		datasetJob(t, rpc, socket, "explicit-raw", `assert(type(archive.read_raw("following")[1].following.accountId)=="number"); return`, map[string]string{"data/following.js": `window.YTD.following.part0 = [{"following":{"accountId":200}}];`}, workerpb.EvaluationState_SUCCEEDED)
	})
	for name, changes := range map[string]map[string]string{
		"canonical-account-input":   {"data/account.js": `window.YTD.account.part0 = [{"account":{"id":"100"}}];`},
		"canonical-like-input":      {"data/like.js": `window.YTD.like.part0 = [{"like":{"post_id":"123"}}];`},
		"canonical-following-input": {"data/following.js": `window.YTD.following.part0 = [{"following":{"account_id":"200"}}];`},
		"case-fold-tweet-input":     {"data/tweets.js": `window.YTD.tweets.part0 = [{"tweet":{"ID_STR":"101","created_at":"Mon Nov 09 12:00:00 +0000 2020","full_text":"root"}}];`},
	} {
		t.Run(name, func(t *testing.T) {
			datasetJob(t, rpc, socket, name, `pcall(archive.likes); pcall(archive.following); pcall(archive.tweets); return`, changes, workerpb.EvaluationState_FAILED)
		})
	}
	t.Run("recovery", func(t *testing.T) {
		datasetJob(t, rpc, socket, "typed-recovery", `assert(#archive.following()==2); return`, nil, workerpb.EvaluationState_SUCCEEDED)
	})
}
