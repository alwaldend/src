package wasm

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	wasmadapter "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/wasm"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/scheduler"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

// ankiWorkerRunfile selects the production worker containing the real Anki guest.
var ankiWorkerRunfile = flag.String("anki-worker", "", "offline Anki WASM worker runfile")

// offlineClient selects the production guest's process or in-memory transport.
func offlineClient(t *testing.T) ankiclient.Client {
	t.Helper()
	if os.Getenv("ANKI_INPROCESS_WASM") == "true" {
		return ankiclient.NewClient(wasmadapter.New())
	}
	path, err := runfiles.Rlocation(*ankiWorkerRunfile)
	if err != nil {
		t.Fatalf("resolve offline Anki worker: %v", err)
	}
	return ankiclient.NewClient(ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: path, Arguments: []string{"serve"}}))
}

// TestOfflineAnki rejects network operations without damaging local collection access.
func TestOfflineAnki(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	worker := offlineClient(t)
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("start offline Anki guest: %v", err)
	}
	t.Cleanup(func() {
		stop, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer stopCancel()
		if err := worker.Stop(stop); err != nil {
			t.Errorf("reap offline Anki worker: %v", err)
		}
	})
	connection, err := worker.Connection()
	if err != nil {
		t.Fatalf("get offline connection: %v", err)
	}
	for _, method := range []string{
		"/anki.sync.BackendSyncService/SyncCollection",
		"/anki.sync.BackendSyncService/SyncMedia",
		"/anki.sync.BackendSyncService/SyncLogin",
		"/anki.ankihub.BackendAnkiHubService/AnkihubLogin",
		"/anki.github.BackendGithubService/DownloadRelease",
		"/anki.media.BackendMediaService/AddMediaFromUrl",
	} {
		if err := connection.Invoke(ctx, method, &generic.Empty{}, &generic.Empty{}); err == nil || !strings.Contains(strings.ToLower(err.Error()), "unsupported") {
			t.Fatalf("network operation %q did not report unsupported: %v", method, err)
		}
	}
	directory := t.TempDir()
	collections := collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath:  filepath.Join(directory, "collection.anki2"),
		MediaFolderPath: filepath.Join(directory, "collection.media"),
		MediaDbPath:     filepath.Join(directory, "collection.media.db2"),
	}); err != nil {
		t.Fatalf("open collection after rejected network calls: %v", err)
	}
	source := filepath.Join(directory, "local-media.txt")
	if err := os.WriteFile(source, []byte("local WASI media"), 0o600); err != nil {
		t.Fatalf("write local media fixture: %v", err)
	}
	added, err := media.NewBackendMediaServiceClient(connection).AddMediaFromUrl(ctx, &media.AddMediaFromUrlRequest{Url: "file://" + source})
	if err != nil || added.GetError() != "" || added.GetFilename() != "local-media.txt" {
		t.Fatalf("local file URL media import: %v; %v", added, err)
	}
	data, err := os.ReadFile(filepath.Join(directory, "collection.media", added.GetFilename()))
	if err != nil || string(data) != "local WASI media" {
		t.Fatalf("guest-written local media changed: %q; %v", data, err)
	}
	backup := filepath.Join(directory, "backups")
	if err := os.Mkdir(backup, 0o700); err != nil {
		t.Fatalf("create backup directory: %v", err)
	}
	created, err := collections.CreateBackup(ctx, &collection.CreateBackupRequest{BackupFolder: backup, Force: true})
	if err != nil || !created.GetVal() {
		t.Fatalf("offline collection backup: %v; %v", created, err)
	}
	if _, err := collections.AwaitBackupCompletion(ctx, &generic.Empty{}); err != nil {
		t.Fatalf("await offline collection backup: %v", err)
	}
	packages, err := filepath.Glob(filepath.Join(backup, "*.colpkg"))
	if err != nil || len(packages) != 1 {
		t.Fatalf("backup artifact missing: %v; %v", packages, err)
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatalf("close offline collection: %v", err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := os.ReadFile(packages[0])
		if err != nil {
			t.Fatalf("read offline backup: %v", err)
		}
		if err := os.WriteFile(filepath.Join(output, "offline-backup.colpkg"), data, 0o600); err != nil {
			t.Fatalf("retain offline backup: %v", err)
		}
		if err := os.WriteFile(filepath.Join(output, "offline-capabilities.txt"), []byte("network calls rejected; local collection open/close passed\n"), 0o600); err != nil {
			t.Fatalf("retain offline capability artifact: %v", err)
		}
	}
}

// TestOfflineSchedulingTimezone preserves the host's non-UTC collection offset.
func TestOfflineSchedulingTimezone(t *testing.T) {
	t.Setenv("TZ", "America/New_York")
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load test timezone: %v", err)
	}
	_, secondsEast := time.Now().In(location).Zone()
	minutesWest := -secondsEast / 60
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	worker := offlineClient(t)
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("start Anki timezone worker: %v", err)
	}
	t.Cleanup(func() {
		stop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := worker.Stop(stop); err != nil {
			t.Errorf("reap Anki timezone worker: %v", err)
		}
	})
	connection, err := worker.Connection()
	if err != nil {
		t.Fatalf("get Anki timezone connection: %v", err)
	}
	directory := t.TempDir()
	collections := collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath:  filepath.Join(directory, "collection.anki2"),
		MediaFolderPath: filepath.Join(directory, "collection.media"),
		MediaDbPath:     filepath.Join(directory, "collection.media.db2"),
	}); err != nil {
		t.Fatalf("open timezone collection: %v", err)
	}
	defer func() {
		if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
			t.Errorf("close timezone collection: %v", err)
		}
	}()
	configs := config.NewConfigServiceClient(connection)
	if _, err := configs.SetConfigJson(ctx, &config.SetConfigJsonRequest{Key: "localOffset", ValueJson: []byte(strconv.Itoa(minutesWest))}); err != nil {
		t.Fatalf("seed local collection offset: %v", err)
	}
	if _, err := scheduler.NewSchedulerServiceClient(connection).SchedTimingToday(ctx, &generic.Empty{}); err != nil {
		t.Fatalf("compute local scheduling day: %v", err)
	}
	got, err := configs.GetConfigJson(ctx, &generic.String{Val: "localOffset"})
	if err != nil || string(got.GetJson()) != strconv.Itoa(minutesWest) {
		t.Fatalf("local offset changed from host %d minutes west to %s: %v", minutesWest, got.GetJson(), err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "anki-timezone-verification.txt"), got.GetJson(), 0o600); err != nil {
			t.Fatalf("retain Anki timezone verification: %v", err)
		}
	}
}
