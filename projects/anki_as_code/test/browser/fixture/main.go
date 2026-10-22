package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/native"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/scheduler"
	"github.com/ankitects/anki/proto/anki/search"
	"github.com/ankitects/anki/proto/anki/stats"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// snapshot records independently observable identities, review state, and media.
type snapshot struct {
	// Note retains the original identity and editable contents.
	Note *notes.Note `json:"note"`
	// Card retains every scheduling field, excluding mutable metadata.
	Card *cards.Card `json:"card"`
	// History is the canonical native review-history response.
	History json.RawMessage `json:"history"`
	// MediaSHA256 identifies all binary media bytes.
	MediaSHA256 string `json:"mediaSHA256"`
}

// main exposes the fixture operations to the JavaScript browser test.
func main() {
	command := &cobra.Command{Use: "fixture", SilenceUsage: true, SilenceErrors: true}
	command.AddCommand(&cobra.Command{
		Use: "create <package> <baseline>", Args: cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			if err := create(command.Context(), args[0], args[1]); err != nil {
				return fmt.Errorf("create native browser fixture: %w", err)
			}
			return nil
		},
	}, &cobra.Command{
		Use: "verify <package> <baseline> <evidence>", Args: cobra.ExactArgs(3),
		RunE: func(command *cobra.Command, args []string) error {
			if err := verify(command.Context(), args[0], args[1], args[2]); err != nil {
				return fmt.Errorf("verify browser package with native Anki: %w", err)
			}
			return nil
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if err := command.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// withCollection borrows an application-owned collection using the native backend.
func withCollection(ctx context.Context, source string, operation func(grpc.ClientConnInterface) error) (result error) {
	application := ankiascode.New(ankiascode.Options{
		NewClient: ankigrpc.ToClientFactory(func(string) (*native.Client, error) {
			return native.New(), nil
		}),
	})
	collection, err := application.ColOpen(ctx, ankiascodeapi.ColOpenOptions{Path: source})
	if err != nil {
		return fmt.Errorf("prepare fixture collection %q: %w", source, err)
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := collection.Stop(cleanup); err != nil {
			result = errors.Join(result, fmt.Errorf("close fixture collection: %w", err))
		}
	}()
	connection, err := collection.Connection()
	if err != nil {
		return fmt.Errorf("borrow fixture connection: %w", err)
	}
	if err := operation(connection); err != nil {
		return fmt.Errorf("run fixture operation: %w", err)
	}
	return nil
}

// create produces a native package with an actual review and binary media.
func create(ctx context.Context, archive, baseline string) (result error) {
	directory, err := os.MkdirTemp(os.Getenv("TEST_TMPDIR"), "native-browser-fixture-")
	if err != nil {
		return fmt.Errorf("create native fixture directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(directory); err != nil {
			result = errors.Join(result, fmt.Errorf("remove native fixture directory: %w", err))
		}
	}()
	database := filepath.Join(directory, "collection.anki2")
	// The application opens existing sources; Anki initializes this empty SQLite file.
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		return fmt.Errorf("create empty fixture collection: %w", err)
	}
	if err := withCollection(ctx, database, func(connection grpc.ClientConnInterface) error {
		definition, err := notetypes.NewNotetypesServiceClient(connection).AddNotetype(ctx, &notetypes.Notetype{
			Name: "Browser fixture", Config: &notetypes.Notetype_Config{},
			Fields:    []*notetypes.Notetype_Field{{Name: "Front", Config: &notetypes.Notetype_Field_Config{}}, {Name: "Back", Config: &notetypes.Notetype_Field_Config{}}},
			Templates: []*notetypes.Notetype_Template{{Name: "Forward", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"}}},
		})
		if err != nil {
			return fmt.Errorf("create fixture note type: %w", err)
		}
		if _, err := notes.NewNotesServiceClient(connection).AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{
			Guid: "browser-fixture", NotetypeId: definition.Id, Fields: []string{"browser fixture original", `<img src="browser-media.bin">`}, Tags: []string{"browser"},
		}}); err != nil {
			return fmt.Errorf("create fixture note: %w", err)
		}
		found, err := search.NewSearchServiceClient(connection).SearchCards(ctx, &search.SearchRequest{})
		if err != nil {
			return fmt.Errorf("find fixture card: %w", err)
		}
		if len(found.Ids) != 1 {
			return fmt.Errorf("expected one generated fixture card, found %d", len(found.Ids))
		}
		reviews := scheduler.NewSchedulerServiceClient(connection)
		states, err := reviews.GetSchedulingStates(ctx, &cards.CardId{Cid: found.Ids[0]})
		if err != nil {
			return fmt.Errorf("read fixture scheduling states: %w", err)
		}
		if _, err := reviews.AnswerCard(ctx, &scheduler.CardAnswer{CardId: found.Ids[0], CurrentState: states.Current, NewState: states.Good, Rating: scheduler.CardAnswer_GOOD, AnsweredAtMillis: time.Now().UnixMilli(), MillisecondsTaken: 250}); err != nil {
			return fmt.Errorf("review fixture card: %w", err)
		}
		binary := make([]byte, 8193)
		for index := range binary {
			binary[index] = byte(index % 251)
		}
		if _, err := media.NewMediaServiceClient(connection).AddMediaFile(ctx, &media.AddMediaFileRequest{DesiredName: "browser-media.bin", Data: binary}); err != nil {
			return fmt.Errorf("write binary fixture media: %w", err)
		}
		observed, err := inspect(ctx, connection)
		if err != nil {
			return fmt.Errorf("capture native baseline: %w", err)
		}
		if err := writeJSON(baseline, observed); err != nil {
			return fmt.Errorf("write fixture baseline: %w", err)
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: archive, IncludeMedia: true}); err != nil {
			return fmt.Errorf("export native fixture package: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("create fixture package %q: %w", archive, err)
	}
	return nil
}

// inspect reads collection state exclusively through generated native APIs.
func inspect(ctx context.Context, connection grpc.ClientConnInterface) (*snapshot, error) {
	listed, err := search.NewSearchServiceClient(connection).SearchCards(ctx, &search.SearchRequest{})
	if err != nil {
		return nil, fmt.Errorf("find inspected card: %w", err)
	}
	if len(listed.Ids) != 1 {
		return nil, fmt.Errorf("expected one preserved card, found %d", len(listed.Ids))
	}
	card, err := cards.NewCardsServiceClient(connection).GetCard(ctx, &cards.CardId{Cid: listed.Ids[0]})
	if err != nil {
		return nil, fmt.Errorf("read preserved card: %w", err)
	}
	card.MtimeSecs, card.Usn = 0, 0
	note, err := notes.NewNotesServiceClient(connection).GetNote(ctx, &notes.NoteId{Nid: card.NoteId})
	if err != nil {
		return nil, fmt.Errorf("read preserved note: %w", err)
	}
	history, err := stats.NewStatsServiceClient(connection).GetReviewLogs(ctx, &cards.CardId{Cid: card.Id})
	if err != nil {
		return nil, fmt.Errorf("read preserved review history: %w", err)
	}
	if len(history.Entries) != 1 {
		return nil, fmt.Errorf("expected one preserved review, found %d", len(history.Entries))
	}
	encoded, err := protojson.Marshal(history)
	if err != nil {
		return nil, fmt.Errorf("encode review history: %w", err)
	}
	var decoded any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("normalize review history: %w", err)
	}
	canonical, err := json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("canonicalize review history: %w", err)
	}
	location, err := media.NewMediaServiceClient(connection).GetAbsoluteMediaPath(ctx, &generic.String{Val: "browser-media.bin"})
	if err != nil {
		return nil, fmt.Errorf("resolve preserved media: %w", err)
	}
	contents, err := os.ReadFile(location.Val)
	if err != nil {
		return nil, fmt.Errorf("read preserved media: %w", err)
	}
	digest := sha256.Sum256(contents)
	return &snapshot{Note: note, Card: card, History: canonical, MediaSHA256: hex.EncodeToString(digest[:])}, nil
}

// verify requires the browser edit while preserving native identity and history.
func verify(ctx context.Context, archive, baseline, evidence string) error {
	contents, err := os.ReadFile(baseline)
	if err != nil {
		return fmt.Errorf("read native baseline: %w", err)
	}
	var expected snapshot
	if err := json.Unmarshal(contents, &expected); err != nil {
		return fmt.Errorf("decode native baseline: %w", err)
	}
	// MarshalIndent expands RawMessage whitespace in the readable baseline.
	var expectedHistory bytes.Buffer
	if err := json.Compact(&expectedHistory, expected.History); err != nil {
		return fmt.Errorf("normalize baseline review history: %w", err)
	}
	if err := withCollection(ctx, archive, func(connection grpc.ClientConnInterface) error {
		actual, err := inspect(ctx, connection)
		if err != nil {
			return fmt.Errorf("inspect browser-produced package: %w", err)
		}
		if err := writeJSON(evidence, actual); err != nil {
			return fmt.Errorf("write native verification evidence: %w", err)
		}
		if actual.Note.Id != expected.Note.Id || actual.Note.Guid != expected.Note.Guid || actual.Note.NotetypeId != expected.Note.NotetypeId {
			return fmt.Errorf("browser round trip changed native note identity")
		}
		expected.Note.Fields[0] = "browser edited 你好"
		if !reflect.DeepEqual(actual.Note.Fields, expected.Note.Fields) || !reflect.DeepEqual(actual.Note.Tags, expected.Note.Tags) {
			return fmt.Errorf("browser note contents differ: %q", actual.Note.Fields)
		}
		if !proto.Equal(actual.Card, expected.Card) {
			return fmt.Errorf("browser round trip changed card identity or scheduling: before=%v after=%v", expected.Card, actual.Card)
		}
		if !bytes.Equal(actual.History, expectedHistory.Bytes()) {
			return fmt.Errorf("browser round trip changed review history: before=%s after=%s", expectedHistory.Bytes(), actual.History)
		}
		if actual.MediaSHA256 != expected.MediaSHA256 {
			return fmt.Errorf("browser round trip changed binary media: before=%s after=%s", expected.MediaSHA256, actual.MediaSHA256)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("verify browser archive %q: %w", archive, err)
	}
	return nil
}

// writeJSON persists a readable verification artifact in the test output tree.
func writeJSON(filename string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode JSON artifact: %w", err)
	}
	if err := os.WriteFile(filename, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write JSON artifact %q: %w", filename, err)
	}
	return nil
}
