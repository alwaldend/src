package collection

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
)

// TestNetworkScopeCleanup verifies confirmed cleanup releases application reservations.
func TestNetworkScopeCleanup(t *testing.T) {
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			_, archive := nativeSource(t)
			address := networkFixtureServer(t, transport)
			store := New(Options{NewClient: ankigrpc.NewFactory(address).NewClient})
			for _, scenario := range []string{"success", "native-error", "filesystem-cleanup"} {
				t.Run(scenario, func(t *testing.T) {
					if scenario == "filesystem-cleanup" && os.Geteuid() == 0 {
						t.Skip("root bypasses directory write permissions")
					}
					client, err := store.Open(t.Context(), ankiascodeapi.ColOpenOptions{Path: archive})
					if err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() { _ = client.Stop(context.Background()) })
					connection, err := client.Connection()
					if err != nil {
						t.Fatal(err)
					}
					restored := networkRestoreDirectory(t, connection)
					if scenario == "native-error" {
						if _, err := notes.NewNotesServiceClient(connection).GetNote(t.Context(), &notes.NoteId{Nid: -1}); err == nil {
							t.Fatal("missing note did not return a native error")
						}
					}
					if scenario == "filesystem-cleanup" {
						if _, err := ankicol.NewBackendCollectionServiceClient(connection).CloseCollection(t.Context(), &ankicol.CloseCollectionRequest{}); err != nil {
							t.Fatal(err)
						}
						if err := os.Chmod(restored, 0o500); err != nil {
							t.Fatal(err)
						}
						t.Cleanup(func() {
							if err := os.Chmod(restored, 0o700); err != nil {
								t.Error(err)
							}
							if err := os.RemoveAll(restored); err != nil {
								t.Error(err)
							}
						})
					}
					cleanup, cancel := context.WithTimeout(t.Context(), 5*time.Second)
					defer cancel()
					err = client.Stop(cleanup)
					if scenario == "filesystem-cleanup" {
						if !errors.Is(err, os.ErrPermission) {
							t.Fatalf("cleanup lost the filesystem failure: %v", err)
						}
						if _, err := os.Stat(restored); err != nil {
							t.Fatalf("failed cleanup discarded the restore: %v", err)
						}
					} else if err != nil {
						t.Fatal(err)
					}
					if _, err := decks.NewDecksServiceClient(connection).GetDeck(cleanup, &decks.DeckId{Did: 1}); err == nil {
						t.Fatal("stopped collection retained native access")
					}
					if err := store.WithCollection(cleanup, ankiascodeapi.ColOpenOptions{Path: archive}, func(grpc.ClientConnInterface) error { return nil }); err != nil {
						t.Fatalf("confirmed native close retained its source: %v", err)
					}
					emitNetworkReport(t, "network-cleanup-"+transport+"-"+scenario+".json", map[string]any{"transport": transport, "scenario": scenario, "reservation_released": true, "external_server_retained": true})
				})
			}
		})
	}
}

// TestScopeConcurrentStop drains an active native call while each closing caller has its own deadline.
func TestScopeConcurrentStop(t *testing.T) {
	_, archive := nativeSource(t)
	store := fixtureStore(t)
	client, err := store.Open(t.Context(), ankiascodeapi.ColOpenOptions{Path: archive})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Stop(context.Background()) })
	connection, err := client.Connection()
	if err != nil {
		t.Fatal(err)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		close(entered)
		<-release
		if _, err := response.Write([]byte("completed scope media")); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	defer releaseOnce.Do(func() { close(release) })
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	downloaded := make(chan error, 1)
	go func() {
		_, err := media.NewBackendMediaServiceClient(connection).AddMediaFromUrl(ctx, &media.AddMediaFromUrlRequest{Url: server.URL + "/completed.txt"})
		downloaded <- err
	}()
	select {
	case <-entered:
	case err := <-downloaded:
		t.Fatalf("native download returned before the fixture: %v", err)
	case <-ctx.Done():
		t.Fatal("native fixture did not start")
	}
	stopped := make(chan error, 1)
	go func() { stopped <- client.Stop(ctx) }()
	joining, cancelJoining := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancelJoining()
	if err := client.Stop(joining); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("closing caller did not observe its own deadline: %v", err)
	}
	releaseOnce.Do(func() { close(release) })
	if err := <-downloaded; err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-stopped:
		if err != nil {
			t.Fatalf("a closing caller cancelled owned cleanup: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("cleanup did not finish after the native call")
	}
	if err := store.WithCollection(ctx, ankiascodeapi.ColOpenOptions{Path: archive}, func(grpc.ClientConnInterface) error { return nil }); err != nil {
		t.Fatalf("completed cleanup retained its source reservation: %v", err)
	}
	emitNetworkReport(t, "concurrent-scope-stop.json", map[string]bool{"active_call_drained": true, "caller_deadline_independent": true, "source_released": true})
}

// TestNetworkScopeCancellation retains restored files and reservations after unconfirmed remote work.
func TestNetworkScopeCancellation(t *testing.T) {
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			_, archive := nativeSource(t)
			var restored string
			t.Cleanup(func() {
				if restored != "" {
					if err := os.RemoveAll(restored); err != nil {
						t.Error(err)
					}
				}
			})
			address := networkFixtureServer(t, transport)
			store := New(Options{NewClient: ankigrpc.NewFactory(address).NewClient, CloseTimeout: 100 * time.Millisecond})
			client, err := store.Open(t.Context(), ankiascodeapi.ColOpenOptions{Path: archive})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = client.Stop(context.Background()) })
			connection, err := client.Connection()
			if err != nil {
				t.Fatal(err)
			}
			restored = networkRestoreDirectory(t, connection)
			entered, release := make(chan struct{}), make(chan struct{})
			var releaseOnce sync.Once
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
				close(entered)
				<-release
				response.Header().Set("Content-Type", "text/plain")
				if _, err := response.Write([]byte("remote collection media")); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			defer releaseOnce.Do(func() { close(release) })
			download, cancelDownload := context.WithCancel(t.Context())
			defer cancelDownload()
			downloaded := make(chan error, 1)
			go func() {
				_, err := media.NewBackendMediaServiceClient(connection).AddMediaFromUrl(download, &media.AddMediaFromUrlRequest{Url: server.URL + "/pending.txt"})
				downloaded <- err
			}()
			startup, cancelStartup := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancelStartup()
			select {
			case <-entered:
			case err := <-downloaded:
				t.Fatalf("native download returned before reaching its fixture: %v", err)
			case <-startup.Done():
				t.Fatal("native download did not start")
			}
			cancelDownload()
			if err := <-downloaded; !errors.Is(err, context.Canceled) {
				t.Fatalf("cancel native download: %v", err)
			}
			cleanup, cancelCleanup := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancelCleanup()
			if err := client.Stop(cleanup); err == nil || errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("unconfirmed cleanup did not report bounded failure: %v", err)
			}
			if _, err := os.Stat(filepath.Join(restored, "collection.anki2")); err != nil {
				t.Fatalf("unconfirmed cleanup removed native storage: %v", err)
			}
			assertNetworkReservation(t, &store, archive)
			releaseOnce.Do(func() { close(release) })
			observer := networkFixtureClient(t, address)
			observed, err := observer.Connection()
			if err != nil {
				t.Fatal(err)
			}
			verification, cancelVerification := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancelVerification()
			if _, err := decks.NewDecksServiceClient(observed).GetDeck(verification, &decks.DeckId{Did: 1}); err != nil {
				t.Fatalf("remote work did not survive disconnection: %v", err)
			}
			if _, err := ankicol.NewBackendCollectionServiceClient(observed).CloseCollection(verification, &ankicol.CloseCollectionRequest{}); err != nil {
				t.Fatal(err)
			}
			assertNetworkReservation(t, &store, archive)
			if err := os.RemoveAll(restored); err != nil {
				t.Fatal(err)
			}
			restored = ""
			emitNetworkReport(t, "network-cancellation-"+transport+".json", map[string]any{"transport": transport, "bounded_stop": true, "active_restore_retained": true, "reservation_retained": true, "external_cleanup_verified": true})
		})
	}
}

// assertNetworkReservation verifies an unconfirmed scope cannot be reopened through its Store.
func assertNetworkReservation(t *testing.T, store *Store, source string) {
	t.Helper()
	client, err := store.Open(t.Context(), ankiascodeapi.ColOpenOptions{Path: source})
	if err == nil {
		_ = client.Stop(context.Background())
	}
	if err == nil || !strings.Contains(err.Error(), "already open") {
		t.Fatalf("unconfirmed cleanup released its source reservation: %v", err)
	}
}

// networkRestoreDirectory observes the native restore location through the public media service.
func networkRestoreDirectory(t *testing.T, connection grpc.ClientConnInterface) string {
	t.Helper()
	path, err := media.NewMediaServiceClient(connection).GetAbsoluteMediaPath(t.Context(), &generic.String{Val: "example.txt"})
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(filepath.Dir(path.Val))
}

// networkFixtureClient connects a raw observer without owning the remote server or collection.
func networkFixtureClient(t *testing.T, address string) *ankiclient.ClientGRPC {
	t.Helper()
	client, err := ankiclient.NewGRPC(address)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := client.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	startup, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	if err := client.Start(startup); err != nil {
		t.Fatal(err)
	}
	return client
}

// networkFixtureServer starts an externally owned native listener and reaps it after its clients.
func networkFixtureServer(t *testing.T, transport string) string {
	t.Helper()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	binary, err = filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())
	listen := "127.0.0.1:0"
	if transport == "unix" {
		listen = "collection.sock"
	}
	command := exec.CommandContext(t.Context(), binary, "serve", "--transport", transport, "--listen", listen)
	command.Stdout = io.Discard
	diagnostics, err := command.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := command.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			t.Error(err)
		}
		var exit *exec.ExitError
		if err := command.Wait(); err != nil && !errors.As(err, &exit) {
			t.Error(err)
		}
	})
	ready := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(diagnostics)
		if scanner.Scan() {
			ready <- scanner.Text()
		} else {
			ready <- ""
		}
		for scanner.Scan() {
		}
	}()
	startup, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	select {
	case line := <-ready:
		if !strings.HasPrefix(line, "Listening on "+transport+"://") {
			t.Fatalf("external collection server did not start: %q", line)
		}
		return strings.TrimPrefix(line, "Listening on ")
	case <-startup.Done():
		t.Fatalf("external collection server startup: %v", startup.Err())
		return ""
	}
}

// emitNetworkReport retains completion evidence beside native collection artifacts.
func emitNetworkReport(t *testing.T, name string, report any) {
	t.Helper()
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	emitArtifact(t, path, name)
}
