package e2e

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// writeStandardArchive creates synthetic export files without using production parsing.
func writeStandardArchive(t *testing.T, directory string, files map[string]string) string {
	t.Helper()
	var data bytes.Buffer
	writer := zip.NewWriter(&data)
	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(entry, content); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(directory, "export.zip")
	if err := os.WriteFile(output, data.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return output
}

// startArchiveWorker runs a credential-free source-owning process for CLI fixtures.
func startArchiveWorker(t *testing.T, binary, directory string) string {
	t.Helper()
	socketDirectory, err := os.MkdirTemp("", "cleaner-offline-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(socketDirectory); err != nil {
			t.Error(err)
		}
	})
	socket := filepath.Join(socketDirectory, "worker.sock")
	process := exec.Command(binary, "worker", "serve", "--socket", socket, "--source", "exports="+directory, "--evaluation-timeout", "10s")
	process.Env = []string{"PATH=/usr/bin:/bin"}
	var logs bytes.Buffer
	process.Stdout, process.Stderr = &logs, &logs
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() { finished <- process.Wait() }()
	stopped := false
	t.Cleanup(func() {
		if stopped {
			return
		}
		if err := process.Process.Signal(os.Interrupt); err != nil {
			t.Error(err)
		}
		select {
		case err := <-finished:
			if err != nil {
				t.Errorf("worker shutdown: %v %s", err, logs.String())
			}
		case <-time.After(6 * time.Second):
			_ = process.Process.Kill()
			<-finished
			t.Error("worker shutdown timeout")
		}
	})
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(socket); err == nil {
			return socket
		}
		select {
		case err := <-finished:
			stopped = true
			t.Fatalf("archive worker did not start: %v %s", err, logs.String())
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("archive worker readiness timeout")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// TestArchiveOnlyCLI proves offline planning and mutation-only credential use.
func TestArchiveOnlyCLI(t *testing.T) {
	binary, err := runfiles.Rlocation(*cleanerBinary)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	writeStandardArchive(t, directory, map[string]string{
		"data/account.js": `window.YTD.account.part0 = [{"account":{"accountId":"100","username":"fixture"}}];`,
		"data/tweets.js":  `window.YTD.tweets.part0 = [{"tweet":{"id_str":"101","full_text":"owned archive post","created_at":"Mon Nov 09 12:00:00 +0000 2020"}}];`,
	})
	socket := startArchiveWorker(t, binary, directory)
	policy := filepath.Join(directory, "cleanup.lua")
	if err := os.WriteFile(policy, []byte("delete_post_when(function(p) return true end)"), 0o600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(directory, "plan.json")
	var reads, mutations atomic.Int32
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			reads.Add(1)
			http.Error(w, "X reads are forbidden", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/2/tweets/101" || r.Header.Get("Authorization") == "" {
			http.Error(w, "unexpected mutation", http.StatusForbidden)
			return
		}
		mutations.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":{"deleted":true}}`)
	}))
	t.Cleanup(endpoint.Close)
	invoke := func(credentials bool, args ...string) (string, error) {
		t.Helper()
		args = append(args, "--worker-socket", socket, "--archive-source", "exports", "--archive-key", "export.zip", "--account", "100", "--api-base-url", endpoint.URL+"/2")
		command := exec.Command(binary, args...)
		command.Env = []string{"PATH=/usr/bin:/bin"}
		if credentials {
			for name, value := range fixtureCredentials {
				command.Env = append(command.Env, name+"="+value)
			}
		}
		data, err := command.CombinedOutput()
		assertPrivateDiagnostics(t, string(data))
		return string(data), err
	}
	if text, err := invoke(false, "plan", "--config", policy, "--out", output); err != nil {
		t.Fatalf("offline archive plan: %v %s", err, text)
	}
	document := readPlanDocument(t, output)
	actions := document["plan"].(map[string]any)["operations"].([]any)
	if len(actions) != 1 || actions[0].(map[string]any)["delete_post"].(map[string]any)["post_id"] != "101" || reads.Load() != 0 || mutations.Load() != 0 {
		t.Fatalf("archive planning effects: %+v reads=%d mutations=%d", document, reads.Load(), mutations.Load())
	}
	if text, err := invoke(false, "apply", "--plan", output); err == nil || !strings.Contains(text, "credentials") {
		t.Fatalf("apply without authorization: %v %s", err, text)
	}
	if text, err := invoke(true, "apply", "--plan", output); err != nil {
		t.Fatalf("authorized exact apply: %v %s", err, text)
	}
	if reads.Load() != 0 || mutations.Load() != 1 {
		t.Fatalf("unexpected API dispatch: reads=%d mutations=%d", reads.Load(), mutations.Load())
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(map[string]any{"archive_account": "100", "target": "101", "x_reads": reads.Load(), "mutations": mutations.Load(), "offline_plan": true}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "archive-cli.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// exportFixtures links each isolated mutation endpoint to its synthetic archive input.
var exportFixtures sync.Map

// writeFixtureArchive projects legacy graph fixtures into independently encoded YTD files.
func writeFixtureArchive(t *testing.T, fixture *fixture, directory string) {
	t.Helper()
	rows := []any{}
	for _, page := range fixture.pages {
		for _, record := range page {
			if record["author_id"] != "100" || record["article"] != nil {
				continue
			}
			encoded, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			var source map[string]any
			if err := json.Unmarshal(encoded, &source); err != nil {
				t.Fatal(err)
			}
			tweet := map[string]any{"id_str": source["id"], "user_id_str": source["author_id"], "full_text": source["text"]}
			if stamp, ok := source["created_at"].(string); ok {
				value, err := time.Parse(time.RFC3339, stamp)
				if err != nil {
					t.Fatal(err)
				}
				tweet["created_at"] = value.Format("Mon Jan 02 15:04:05 -0700 2006")
			}
			if fixture.malformed {
				delete(tweet, "created_at")
			}
			for input, output := range map[string]string{"lang": "lang", "edit_history_post_ids": "edit_history_tweet_ids"} {
				if value, ok := source[input]; ok {
					tweet[output] = value
				}
			}
			if metrics, ok := source["public_metrics"].(map[string]any); ok {
				for input, output := range map[string]string{"like_count": "favorite_count", "repost_count": "retweet_count", "reply_count": "reply_count", "quote_count": "quote_count", "bookmark_count": "bookmark_count", "impression_count": "impression_count"} {
					if value, ok := metrics[input]; ok {
						tweet[output] = value
					}
				}
			}
			references, _ := source["referenced_posts"].([]any)
			for _, value := range references {
				reference := value.(map[string]any)
				key := map[string]string{"replied_to": "in_reply_to_status_id_str", "quoted": "quoted_status_id_str", "retweeted": "retweeted_status_id_str"}[reference["type"].(string)]
				tweet[key] = reference["id"]
			}
			entities := source["entities"]
			if note, ok := source["note_post"].(map[string]any); ok {
				tweet["full_text"], entities = note["text"], note["entities"]
			}
			if groups, ok := entities.(map[string]any); ok {
				result := map[string]any{}
				for input, output := range map[string]string{"hashtags": "hashtags", "cashtags": "symbols", "mentions": "user_mentions", "urls": "urls"} {
					if entries, ok := groups[input].([]any); ok {
						converted := []any{}
						for _, entry := range entries {
							original := entry.(map[string]any)
							value := map[string]any{"indices": []any{original["start"], original["end"]}}
							for from, to := range map[string]string{"tag": "text", "id": "id_str", "username": "screen_name", "url": "url", "expanded_url": "expanded_url", "display_url": "display_url"} {
								if field, ok := original[from]; ok {
									value[to] = field
								}
							}
							converted = append(converted, value)
						}
						result[output] = converted
					}
				}
				tweet["entities"] = result
			}
			if attachments, ok := source["attachments"].(map[string]any); ok {
				media := []any{}
				keys, _ := attachments["media_keys"].([]any)
				sources, _ := attachments["media_source_tweet_id"].([]any)
				for index, key := range keys {
					id := key.(string)
					if _, suffix, ok := strings.Cut(id, "_"); ok {
						id = suffix
					}
					item := map[string]any{"id_str": id, "media_key": key}
					if index < len(sources) {
						item["source_status_id_str"] = sources[index]
					}
					media = append(media, item)
				}
				tweet["extended_entities"] = map[string]any{"media": media}
			}
			rows = append(rows, map[string]any{"tweet": tweet})
		}
	}
	encoded, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	writeStandardArchive(t, directory, map[string]string{
		"data/account.js": `window.YTD.account.part0 = [{"account":{"accountId":"100"}}];`,
		"data/tweets.js":  "window.YTD.tweets.part0 = " + string(encoded) + ";",
	})
}
