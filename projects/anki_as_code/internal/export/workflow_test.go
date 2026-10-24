package export

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/scheduler"
	"github.com/ankitects/anki/proto/anki/search"
	"github.com/ankitects/anki/proto/anki/stats"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
)

// TestExportWorkflow exercises the independent exporter against a native collection.
func TestExportWorkflow(t *testing.T) {
	ctx := context.Background()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: []string{"serve"}})
	var collections collection.BackendCollectionServiceClient
	opened := false
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if opened {
			if _, err := collections.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
				t.Error(err)
			}
		}
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	collections = collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath:  filepath.Join(directory, "collection.anki2"),
		MediaFolderPath: filepath.Join(directory, "collection.media"),
		MediaDbPath:     filepath.Join(directory, "collection.media.db2"),
	}); err != nil {
		t.Fatal(err)
	}
	opened = true
	types := notetypes.NewNotetypesServiceClient(connection)
	noteClient := notes.NewNotesServiceClient(connection)
	deckClient := decks.NewDecksServiceClient(connection)
	store := state.New(connection, "export-test")
	otherStore := state.New(connection, "unrelated-test")
	exporter := New(resource.New(connection), store, filesystem.FS{})
	definition := &notetypes.Notetype{
		Name: "Readable", Config: &notetypes.Notetype_Config{Css: ".card { color: red; }", SortFieldIdx: 1},
		Fields: []*notetypes.Notetype_Field{
			{Name: "Front", Config: &notetypes.Notetype_Field_Config{}},
			{Name: "Back", Config: &notetypes.Notetype_Field_Config{}},
		},
		Templates: []*notetypes.Notetype_Template{{Name: "Card", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"}}},
	}
	created, err := types.AddNotetype(ctx, definition)
	if err != nil {
		t.Fatal(err)
	}
	deck, err := deckClient.GetDeck(ctx, &decks.DeckId{Did: 1})
	if err != nil {
		t.Fatal(err)
	}
	deck.Name = "Language::Lesson"
	if _, err := deckClient.UpdateDeck(ctx, deck); err != nil {
		t.Fatal(err)
	}
	fields := []string{"A \"quote\"\nand another line\n", "<b>Hello &amp; 世界</b>"}
	added, err := noteClient.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: created.Id, Fields: fields, Tags: []string{"lesson"}}})
	if err != nil {
		t.Fatal(err)
	}
	cardIDs, err := search.NewSearchServiceClient(connection).SearchCards(ctx, &search.SearchRequest{Search: fmt.Sprintf("nid:%d", added.NoteId)})
	if err != nil || len(cardIDs.GetIds()) != 1 {
		t.Fatalf("generated cards: %v", err)
	}
	schedulerClient := scheduler.NewSchedulerServiceClient(connection)
	scheduling, err := schedulerClient.GetSchedulingStates(ctx, &cards.CardId{Cid: cardIDs.Ids[0]})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := schedulerClient.AnswerCard(ctx, &scheduler.CardAnswer{CardId: cardIDs.Ids[0], CurrentState: scheduling.Current, NewState: scheduling.Good, Rating: scheduler.CardAnswer_GOOD, AnsweredAtMillis: time.Now().UnixMilli(), MillisecondsTaken: 100}); err != nil {
		t.Fatal(err)
	}
	statsClient := stats.NewStatsServiceClient(connection)
	history, err := statsClient.GetReviewLogs(ctx, &cards.CardId{Cid: cardIDs.Ids[0]})
	if err != nil || len(history.GetEntries()) == 0 {
		t.Fatalf("review history: %v", err)
	}
	other := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Notes: map[string]*ankiascodepb.NoteState{"notes/unrelated.note.anki.toml": {Id: added.NoteId}}}}
	if _, err := otherStore.Replace(ctx, other); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "text")
	if err := exporter.Export(ctx, output, ""); err != nil {
		t.Fatal(err)
	}
	configs, err := config.New(output, filesystem.FS{}).Load()
	if err != nil {
		t.Fatal(err)
	}
	address := "decks/Language/decks/Lesson/notes/Hello-世界.note.anki.toml"
	for _, absoluteOutput := range []bool{false, true} {
		t.Run("rooted-output-"+fmt.Sprint(absoluteOutput), func(t *testing.T) {
			storage := filesystem.New(afero.NewBasePathFs(afero.NewMemMapFs(), "/scope"))
			root, stateOutput := "configs", "/state.toml"
			if absoluteOutput {
				root, stateOutput = "/configs", "state.toml"
			}
			if err := New(resource.New(connection), store, storage).Export(ctx, root, stateOutput); err != nil {
				t.Fatal(err)
			}
			actual, err := config.New(root, storage).Load()
			if err != nil || !proto.Equal(actual, configs) {
				t.Fatalf("independent rooted export changed declarations: %v %v", actual, err)
			}
			data, err := storage.ReadFile(stateOutput)
			if err != nil {
				t.Fatal(err)
			}
			bindings := &ankiascodepb.State{}
			if err := codec.Decode(data, codec.TOML, bindings); err != nil || bindings.GetResources().GetNotes()[address].GetId() != added.NoteId {
				t.Fatalf("independent rooted export lost note bindings: %v %v", bindings, err)
			}
			if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
				if err := os.WriteFile(filepath.Join(artifact, "independent-rooted-state-"+fmt.Sprint(absoluteOutput)+".toml"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	for _, name := range []string{"state.anki.toml", "state.anki.json"} {
		if _, err := os.Stat(filepath.Join(output, name)); !os.IsNotExist(err) {
			t.Fatalf("export generated %s: %v", name, err)
		}
	}
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.CopyFS(filepath.Join(artifact, "initial"), os.DirFS(output)); err != nil {
			t.Fatal(err)
		}
	}
	note := &ankiascodepb.NoteConfig{}
	readExportConfig(t, filepath.Join(output, address), note)
	if note.Fields["Front"] != fields[0] || note.Fields["Back"] != fields[1] || note.NoteType != "note_types/Readable.notetype.anki.toml" {
		t.Fatalf("note content: %v", note)
	}
	settings := &ankiascodepb.NoteTypeConfig{}
	readExportConfig(t, filepath.Join(output, note.NoteType), settings)
	if settings.SortField != "Back" || settings.Config.Css != definition.Config.Css || settings.Config.SortFieldIdx != 0 || settings.Fields[0].Config.Id != nil || settings.Templates[0].Config.Id != nil || settings.Templates[0].Config.QFormat != "{{Front}}" {
		t.Fatalf("editable type: %v", settings)
	}
	stored, err := store.Dump(ctx)
	if err != nil || len(stored.GetResources().GetNotes()) != 0 {
		t.Fatalf("export installed state: %v", err)
	}
	if _, err := store.Import(ctx, configs, map[string]string{address: fmt.Sprint(added.NoteId), note.NoteType: fmt.Sprint(created.Id)}); err != nil {
		t.Fatal(err)
	}
	bindings, err := store.Dump(ctx)
	if err != nil || len(bindings.GetResources().GetNoteTypes()[note.NoteType].GetFields()) != 2 || len(bindings.GetResources().GetNoteTypes()[note.NoteType].GetTemplates()) != 1 {
		t.Fatalf("explicit state import omitted native member bindings: %v", err)
	}
	nativeType, err := types.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: created.Id})
	if err != nil {
		t.Fatal(err)
	}
	nativeType.Name = "Renamed"
	if _, err := types.UpdateNotetype(ctx, nativeType); err != nil {
		t.Fatal(err)
	}
	definition.Id, definition.Name = 0, "Readable"
	discoveredType, err := types.AddNotetype(ctx, definition)
	if err != nil {
		t.Fatal(err)
	}
	discoveredNote, err := noteClient.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: discoveredType.Id, Fields: fields}})
	if err != nil {
		t.Fatal(err)
	}
	filtered, err := deckClient.AddDeck(ctx, &decks.Deck{Name: "Filtered", Kind: &decks.Deck_Filtered_{Filtered: &decks.Deck_Filtered{Reschedule: true}}})
	if err != nil {
		t.Fatal(err)
	}
	filteredNote, err := noteClient.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: created.Id, Fields: []string{"Filtered", "Resident"}}})
	if err != nil {
		t.Fatal(err)
	}
	residentIDs, err := search.NewSearchServiceClient(connection).SearchCards(ctx, &search.SearchRequest{Search: fmt.Sprintf("nid:%d", filteredNote.NoteId)})
	if err != nil || len(residentIDs.GetIds()) != 1 {
		t.Fatalf("filtered fixture: %v", err)
	}
	cardClient := cards.NewCardsServiceClient(connection)
	resident, err := cardClient.GetCard(ctx, &cards.CardId{Cid: residentIDs.Ids[0]})
	if err != nil {
		t.Fatal(err)
	}
	resident.OriginalDeckId, resident.DeckId, resident.OriginalDue = 1, filtered.Id, resident.Due
	if _, err := cardClient.UpdateCards(ctx, &cards.UpdateCardsRequest{Cards: []*cards.Card{resident}, SkipUndoEntry: true}); err != nil {
		t.Fatal(err)
	}
	t.Run("partial-output", func(t *testing.T) {
		blocked := filepath.Join(t.TempDir(), "text")
		obstruction := filepath.Join(blocked, "note_types", "Readable.notetype.anki.toml")
		if err := os.MkdirAll(obstruction, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(obstruction, "keep"), []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := exporter.Export(ctx, blocked, ""); err == nil {
			t.Fatal("config-file directory obstruction was ignored")
		}
		if _, err := os.Stat(filepath.Join(blocked, "collection.anki.toml")); err != nil {
			t.Fatalf("successful config writes were rolled back: %v", err)
		}
		if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
			if err := os.CopyFS(filepath.Join(artifact, "partial"), os.DirFS(blocked)); err != nil {
				t.Fatal(err)
			}
		}
	})
	private := filepath.Join(output, "private")
	if err := os.Mkdir(private, 0o700); err != nil {
		t.Fatal(err)
	}
	privateFile := filepath.Join(private, "keep.txt")
	if err := os.WriteFile(privateFile, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(output, "empty")
	if err := os.Mkdir(empty, 0o700); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(output, "obsolete.note.anki.toml")
	if err := os.WriteFile(stale, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"state.anki.toml", "state.anki.json"} {
		if err := os.WriteFile(filepath.Join(output, name), []byte("caller state"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := exporter.Export(ctx, output, ""); err != nil {
		t.Fatal(err)
	}
	for name, title := range map[string]string{"Readable": "Renamed", "Readable-2": "Readable"} {
		definition := &ankiascodepb.NoteTypeConfig{}
		readExportConfig(t, filepath.Join(output, "note_types", name+".notetype.anki.toml"), definition)
		if definition.Name != title {
			t.Fatal("tracked/discovered filenames collided")
		}
	}
	duplicate := &ankiascodepb.NoteConfig{}
	readExportConfig(t, filepath.Join(output, "decks/Language/decks/Lesson/notes/Hello-世界-2.note.anki.toml"), duplicate)
	if duplicate.NoteType != "note_types/Readable-2.notetype.anki.toml" {
		t.Fatal("discovered note references wrong type")
	}
	if _, err := config.New(output, filesystem.FS{}).Load(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"state.anki.toml", "state.anki.json"} {
		if data, err := os.ReadFile(filepath.Join(output, name)); err != nil || string(data) != "caller state" {
			t.Fatalf("export changed caller-owned state file %q: %v", name, err)
		}
	}
	retained, err := store.Dump(ctx)
	if err != nil || !proto.Equal(bindings, retained) {
		t.Fatalf("export changed stored state: %v", err)
	}
	if data, err := os.ReadFile(stale); err != nil || string(data) != "stale" {
		t.Fatal("export removed or changed an unselected config file")
	}
	for name, mode := range map[string]os.FileMode{private: 0o700, privateFile: 0o600, empty: 0o700} {
		info, err := os.Stat(name)
		if err != nil || info.Mode().Perm() != mode {
			t.Fatalf("unrelated permissions changed: %s: %v", name, err)
		}
	}
	content, err := os.ReadFile(privateFile)
	if err != nil || string(content) != "keep" {
		t.Fatal("unrelated bytes changed")
	}
	filteredConfig := &ankiascodepb.DeckConfig{}
	readExportConfig(t, filepath.Join(output, "decks/Filtered/deck.anki.toml"), filteredConfig)
	if filteredConfig.Filtered == nil || !filteredConfig.Filtered.Reschedule {
		t.Fatal("filtered config missing")
	}
	residentConfig := &ankiascodepb.NoteConfig{}
	readExportConfig(t, filepath.Join(output, "notes/Resident.note.anki.toml"), residentConfig)
	if len(residentConfig.Cards) != 0 {
		t.Fatal("resident filtered card exported as managed")
	}
	before := exportTree(t, output)
	t.Run("cancelled", func(t *testing.T) {
		cancelled, cancel := context.WithCancel(ctx)
		cancel()
		if err := exporter.Export(cancelled, output, ""); err == nil {
			t.Fatal("cancelled export succeeded")
		}
		if after := exportTree(t, output); !maps.Equal(before, after) {
			t.Fatal("cancelled export changed output")
		}
	})
	configs, err = config.New(output, filesystem.FS{}).Load()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Import(ctx, configs, map[string]string{
		"note_types/Readable-2.notetype.anki.toml":                    fmt.Sprint(discoveredType.Id),
		"decks/Language/decks/Lesson/notes/Hello-世界-2.note.anki.toml": fmt.Sprint(discoveredNote.NoteId),
	}); err != nil {
		t.Fatal(err)
	}
	if err := exporter.Export(ctx, output, ""); err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(before, exportTree(t, output)) {
		t.Fatal("installed state did not stabilize re-export")
	}
	unchanged, err := noteClient.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil || !slices.Equal(unchanged.Fields, fields) {
		t.Fatalf("native fields changed: %v", err)
	}
	remainingHistory, err := statsClient.GetReviewLogs(ctx, &cards.CardId{Cid: cardIDs.Ids[0]})
	if err != nil || !proto.Equal(history, remainingHistory) {
		t.Fatalf("review history changed: %v", err)
	}
	remainingOther, err := otherStore.Dump(ctx)
	if err != nil || !proto.Equal(other, remainingOther) {
		t.Fatalf("other namespace changed: %v", err)
	}
	if _, err := types.GetNotetypeNames(ctx, &generic.Empty{}); err != nil {
		t.Fatalf("export closed borrowed collection: %v", err)
	}
	t.Run("native-filename-text", func(t *testing.T) {
		cases := map[string]string{
			"Hello-world":           "<b>Hel</b>lo<br>world",
			"Sound":                 "[sound:ignored.mp3]<i>Sound</i>",
			"note":                  "<img src=\"image.png\">[sound:ignored.mp3]",
			"note-CON":              "CON",
			"école-日本語":             "<b>e\u0301</b>cole 日本語",
			strings.Repeat("界", 40): strings.Repeat("界", 80),
		}
		for name, field := range cases {
			added, err := noteClient.AddNote(ctx, &notes.AddNoteRequest{
				DeckId: 1,
				Note:   &notes.Note{NotetypeId: created.Id, Fields: []string{"Filename fixture", field}},
			})
			if err != nil {
				t.Fatal(err)
			}
			native, err := noteClient.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
			if err != nil {
				t.Fatal(err)
			}
			cases[name] = native.Fields[1]
		}
		nativeOutput := filepath.Join(t.TempDir(), "text")
		if err := exporter.Export(ctx, nativeOutput, ""); err != nil {
			t.Fatal(err)
		}
		for name, field := range cases {
			value := &ankiascodepb.NoteConfig{}
			readExportConfig(t, filepath.Join(nativeOutput, "decks/Language/decks/Lesson/notes", name+".note.anki.toml"), value)
			if value.Fields["Back"] != field {
				t.Fatalf("filename conversion changed sort field %q", name)
			}
		}
		if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
			if err := os.CopyFS(filepath.Join(artifact, "native-names"), os.DirFS(nativeOutput)); err != nil {
				t.Fatal(err)
			}
		}
	})
	if _, err := media.NewMediaServiceClient(connection).AddMediaFile(ctx, &media.AddMediaFileRequest{DesiredName: "example.txt", Data: []byte("public native fixture")}); err != nil {
		t.Fatal(err)
	}
	artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if artifact != "" {
		if err := os.CopyFS(filepath.Join(artifact, "text"), os.DirFS(output)); err != nil {
			t.Fatal(err)
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: filepath.Join(artifact, "export.colpkg"), IncludeMedia: true}); err != nil {
			t.Fatal(err)
		}
		opened = false
	}
}

// readExportConfig reads one independently loadable protobuf-backed TOML file.
func readExportConfig(t *testing.T, filename string, value proto.Message) {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if err := codec.Decode(data, codec.TOML, value); err != nil {
		t.Fatal(err)
	}
}

// exportTree captures all exported bytes for read-only and stable re-export checks.
func exportTree(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, failure error) error {
		if failure != nil {
			return fmt.Errorf("inspect exported tree: %w", failure)
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read exported file: %w", err)
		}
		result[strings.TrimPrefix(name, root)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
