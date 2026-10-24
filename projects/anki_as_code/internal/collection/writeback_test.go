package collection

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/ankitects/anki/proto/anki/notes"
)

// TestStoreWriteBackPolicy verifies archive persistence follows the selected application policy.
func TestStoreWriteBackPolicy(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(strconv.FormatBool(enabled), func(t *testing.T) {
			_, archive := nativeSource(t)
			before, err := os.ReadFile(archive)
			if err != nil {
				t.Fatal(err)
			}
			store := fixtureStore(t)
			client, err := store.Open(t.Context(), ankiascodeapi.ColOpenOptions{Path: archive, WriteBack: enabled})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = client.Stop(context.Background()) })
			connection, err := client.Connection()
			if err != nil {
				t.Fatal(err)
			}
			value := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"publication": {Id: 1}}}}
			if _, err := state.New(connection, "").Replace(t.Context(), value); err != nil {
				t.Fatal(err)
			}
			if _, err := notes.NewNotesServiceClient(connection).GetNote(t.Context(), &notes.NoteId{Nid: -1}); err == nil {
				t.Fatal("invalid note unexpectedly exists")
			}
			if err := client.Stop(t.Context()); err != nil {
				t.Fatalf("native error suppressed collection cleanup: %v", err)
			}
			if !enabled {
				after, err := os.ReadFile(archive)
				if err != nil || !bytes.Equal(before, after) {
					t.Fatalf("default write-back policy changed source: %v", err)
				}
			}
			actual, err := readState(t.Context(), store, archive, "")
			if err != nil {
				t.Fatal(err)
			}
			if persisted := actual.GetResources().GetPresets()["publication"].GetId() == 1; persisted != enabled {
				t.Fatalf("archive persistence=%v, enabled=%v", persisted, enabled)
			}
			emitArtifact(t, archive, "write-back-"+strconv.FormatBool(enabled)+".colpkg")
		})
	}
}

// TestStorePublicationFailure retains recoverable native storage and the original close error.
func TestStorePublicationFailure(t *testing.T) {
	_, archive := nativeSource(t)
	store := fixtureStore(t)
	client, err := store.Open(t.Context(), ankiascodeapi.ColOpenOptions{Path: archive, WriteBack: true})
	if err != nil {
		t.Fatal(err)
	}
	var restored string
	t.Cleanup(func() {
		_ = client.Stop(context.Background())
		if restored != "" {
			if err := os.RemoveAll(restored); err != nil {
				t.Error(err)
			}
		}
	})
	connection, err := client.Connection()
	if err != nil {
		t.Fatal(err)
	}
	restored = networkRestoreDirectory(t, connection)
	value := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"recoverable": {Id: 1}}}}
	if _, err := state.New(connection, "").Replace(t.Context(), value); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(archive, archive+".original"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(archive, 0o700); err != nil {
		t.Fatal(err)
	}
	closeCause := client.Stop(t.Context())
	if closeCause == nil {
		t.Fatal("write-back to a directory unexpectedly succeeded")
	}
	for errors.Unwrap(closeCause) != nil {
		closeCause = errors.Unwrap(closeCause)
	}
	if err := client.Stop(t.Context()); !errors.Is(err, closeCause) {
		t.Fatalf("repeated cleanup lost its publication failure: %v", err)
	}
	database := filepath.Join(restored, "collection.anki2")
	actual, err := readState(t.Context(), store, database, "")
	if err != nil || actual.GetResources().GetPresets()["recoverable"].GetId() != 1 {
		t.Fatalf("failed publication lost recoverable state: %v", err)
	}
	emitArtifact(t, database, "publication-recovery.anki2")
}
