//go:build js && wasm

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"syscall/js"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/browser"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/scheduler"
	"github.com/spf13/afero"
)

// storageVerification records persisted values after two independent Rust lifetimes.
type storageVerification struct {
	// NoteID is allocated by Anki and read again after backend replacement.
	NoteID int64 `json:"noteId"`
	// Fields retains the exact Unicode content read from the reopened collection.
	Fields []string `json:"fields"`
	// MediaSHA256 identifies the binary media written by Rust and read by Go.
	MediaSHA256 string `json:"mediaSHA256"`
	// ConcurrentMediaSHA256 identifies files written by simultaneous generated calls.
	ConcurrentMediaSHA256 []string `json:"concurrentMediaSHA256"`
	// CollectionBytes records the closed SQLite file visible in the shared store.
	CollectionBytes int `json:"collectionBytes"`
	// CollectionSHA256 identifies the final closed database retained in shared provider storage.
	CollectionSHA256 string `json:"collectionSHA256"`
}

// verifyStorage uses one caller-owned JavaScript filesystem across independently owned browser clients.
func verifyStorage(report *verification) error {
	root := js.Global().Get("ankiBrowserBackendRoot").String()
	if !filepath.IsAbs(root) {
		return fmt.Errorf("expected an absolute caller-supplied filesystem root, got %q", root)
	}
	storage := afero.NewBasePathFs(afero.NewOsFs(), root)
	if err := storage.MkdirAll("/workspace", 0o700); err != nil {
		return fmt.Errorf("prepare shared workspace: %w", err)
	}
	if err := afero.WriteFile(storage, "/blocked", []byte("not a directory"), 0o600); err != nil {
		return fmt.Errorf("prepare invalid collection parent: %w", err)
	}
	report.Storage = &storageVerification{ConcurrentMediaSHA256: make([]string, 2)}
	for session := range 2 {
		if err := verifyStorageSession(storage, session == 0, report); err != nil {
			return fmt.Errorf("verify browser storage session %d: %w", session, err)
		}
		data, err := afero.ReadFile(storage, "/workspace/collection.anki2")
		if err != nil {
			return fmt.Errorf("read closed SQLite collection from Go: %w", err)
		}
		if !bytes.HasPrefix(data, []byte("SQLite format 3\x00")) {
			return fmt.Errorf("shared collection is not a SQLite database")
		}
		report.Storage.CollectionBytes = len(data)
		report.Storage.CollectionSHA256 = fmt.Sprintf("%x", sha256.Sum256(data))
		data, err = afero.ReadFile(storage, "/workspace/collection.media/rust-media.bin")
		if err != nil {
			return fmt.Errorf("read Rust-created media from Go after shutdown: %w", err)
		}
		if !bytes.Equal(data, storageMedia()) {
			return fmt.Errorf("binary media changed after session %d", session)
		}
		report.Storage.MediaSHA256 = fmt.Sprintf("%x", sha256.Sum256(data))
		for index := range report.Storage.ConcurrentMediaSHA256 {
			name := fmt.Sprintf("/workspace/collection.media/concurrent-%d.bin", index)
			data, err := afero.ReadFile(storage, name)
			if err != nil {
				return fmt.Errorf("read concurrent media %q after session %d: %w", name, session, err)
			}
			if !bytes.Equal(data, []byte{byte(index), 0, 255, 128}) {
				return fmt.Errorf("concurrent media %q changed after session %d", name, session)
			}
			report.Storage.ConcurrentMediaSHA256[index] = fmt.Sprintf("%x", sha256.Sum256(data))
		}
		if session == 0 {
			if err := afero.WriteFile(storage, "/workspace/collection.media/go-media.bin", []byte{0, 254, 3, 128}, 0o600); err != nil {
				return fmt.Errorf("write Go-owned media before reopen: %w", err)
			}
		}
	}
	report.Cases = append(report.Cases, "closed SQLite and binary media survive backend replacement in one caller-owned JavaScript filesystem")
	return nil
}

// storageMedia crosses the WASI host's 64 KiB copy boundary with non-text data.
func storageMedia() []byte {
	return bytes.Repeat([]byte{0, 255, 128, 13, 10, 1, 254}, 16384)
}

// writeConcurrentMedia exercises shared-reactor entry through the shared synchronous JavaScript filesystem.
func writeConcurrentMedia(ctx context.Context, client media.MediaServiceClient) error {
	start := make(chan struct{})
	results := make(chan error, 2)
	for index := range 2 {
		go func(index int) {
			<-start
			name := fmt.Sprintf("concurrent-%d.bin", index)
			added, err := client.AddMediaFile(ctx, &media.AddMediaFileRequest{
				DesiredName: name, Data: []byte{byte(index), 0, 255, 128},
			})
			if err != nil {
				results <- fmt.Errorf("write concurrent media %q: %w", name, err)
			} else if added.GetVal() != name {
				results <- fmt.Errorf("concurrent media returned filename %q instead of %q", added.GetVal(), name)
			} else {
				results <- nil
			}
		}(index)
	}
	close(start)
	var failure error
	for range 2 {
		failure = errors.Join(failure, <-results)
	}
	if failure != nil {
		return fmt.Errorf("complete concurrent browser media writes: %w", failure)
	}
	return nil
}

// verifyStorageSession owns one backend and collection, including failure cleanup.
func verifyStorageSession(storage afero.Fs, create bool, report *verification) (result error) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	client := browser.NewClient(js.Global().Get("ankiBrowserBackendFactory"))
	defer func() {
		cleanup, finish := context.WithTimeout(context.Background(), 10*time.Second)
		defer finish()
		if err := client.Stop(cleanup); err != nil {
			result = errors.Join(result, fmt.Errorf("release browser storage client: %w", err))
		}
		if err := waitDone(cleanup, client); err != nil {
			result = errors.Join(result, fmt.Errorf("join browser storage cleanup: %w", err))
		}
	}()
	if err := client.Start(ctx); err != nil {
		return fmt.Errorf("start browser storage client: %w", err)
	}
	connection, err := client.Connection()
	if err != nil {
		return fmt.Errorf("borrow browser storage connection: %w", err)
	}
	collections := collection.NewBackendCollectionServiceClient(connection)
	if create {
		if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
			CollectionPath:  "/blocked/collection.anki2",
			MediaFolderPath: "/blocked/collection.media", MediaDbPath: "/blocked/collection.media.db2",
		}); err == nil {
			return fmt.Errorf("collection open accepted a file as its parent directory")
		}
		report.Cases = append(report.Cases, "invalid filesystem path rejects collection open")
	}
	if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath:  "/workspace/collection.anki2",
		MediaFolderPath: "/workspace/collection.media", MediaDbPath: "/workspace/collection.media.db2",
	}); err != nil {
		return fmt.Errorf("open shared collection: %w", err)
	}
	opened := true
	defer func() {
		if opened {
			cleanup, finish := context.WithTimeout(context.Background(), 10*time.Second)
			defer finish()
			if _, err := collections.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
				result = errors.Join(result, fmt.Errorf("close shared collection after failure: %w", err))
			}
		}
	}()
	items := notes.NewNotesServiceClient(connection)
	mediaClient := media.NewMediaServiceClient(connection)
	fields := []string{"Browser persistence: 你好, café", "Shared JavaScript storage"}
	if create {
		model, err := notetypes.NewNotetypesServiceClient(connection).GetNotetypeIdByName(ctx, &generic.String{Val: "Basic"})
		if err != nil {
			return fmt.Errorf("find built-in Basic note type: %w", err)
		}
		added, err := items.AddNote(ctx, &notes.AddNoteRequest{
			Note: &notes.Note{NotetypeId: model.GetNtid(), Fields: fields}, DeckId: 1,
		})
		if err != nil {
			return fmt.Errorf("write note through generated browser client: %w", err)
		}
		if added.GetNoteId() == 0 {
			return fmt.Errorf("Anki did not allocate a note ID")
		}
		report.Storage.NoteID = added.GetNoteId()
		addedMedia, err := mediaClient.AddMediaFile(ctx, &media.AddMediaFileRequest{DesiredName: "rust-media.bin", Data: storageMedia()})
		if err != nil {
			return fmt.Errorf("write binary media through Rust: %w", err)
		}
		if addedMedia.GetVal() != "rust-media.bin" {
			return fmt.Errorf("unexpected media filename %q", addedMedia.GetVal())
		}
		report.Cases = append(report.Cases, "valid collection opens after a filesystem error and stores a note and binary media")
		if err := writeConcurrentMedia(ctx, mediaClient); err != nil {
			return fmt.Errorf("verify simultaneous generated media calls: %w", err)
		}
		report.Cases = append(report.Cases, "concurrent media commands complete without reentering the Rust backend")
	} else {
		if _, err := items.GetNote(ctx, &notes.NoteId{Nid: -1}); err == nil {
			return fmt.Errorf("missing note lookup unexpectedly succeeded")
		}
		note, err := items.GetNote(ctx, &notes.NoteId{Nid: report.Storage.NoteID})
		if err != nil {
			return fmt.Errorf("read persisted note after backend replacement: %w", err)
		}
		if note.GetId() != report.Storage.NoteID || !slices.Equal(note.GetFields(), fields) {
			return fmt.Errorf("reopened note lost its identity or Unicode fields")
		}
		report.Storage.Fields = note.GetFields()
		checked, err := mediaClient.CheckMedia(ctx, &generic.Empty{})
		if err != nil {
			return fmt.Errorf("inspect shared media directory through Rust: %w", err)
		}
		if len(checked.GetMissing()) != 0 || !slices.Contains(checked.GetUnused(), "go-media.bin") || !slices.Contains(checked.GetUnused(), "rust-media.bin") {
			return fmt.Errorf("Rust did not observe both shared media files: %v", checked)
		}
		report.Cases = append(report.Cases, "fresh backend reads persisted note and Go-created media after recoverable native error")
	}
	if _, err := scheduler.NewSchedulerServiceClient(connection).SchedTimingToday(ctx, &generic.Empty{}); err != nil {
		return fmt.Errorf("compute scheduling day with browser timezone: %w", err)
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		return fmt.Errorf("close shared collection before shutdown: %w", err)
	}
	opened = false
	if err := client.Stop(ctx); err != nil {
		return fmt.Errorf("stop closed collection backend: %w", err)
	}
	if _, err := items.GetNote(ctx, &notes.NoteId{Nid: report.Storage.NoteID}); err == nil {
		return fmt.Errorf("retained generated client accepted a call after shutdown")
	}
	return nil
}
