package ankiascode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/BurntSushi/toml"
	"github.com/ankitects/anki/proto/anki/cards"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/deck_config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/scheduler"
	"github.com/ankitects/anki/proto/anki/search"
	"github.com/ankitects/anki/proto/anki/stats"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// fixtureClient creates and inspects synthetic collections through Anki's API.
type fixtureClient struct {
	// opened records whether cleanup must close a collection; native package export closes it itself.
	opened bool
	// worker owns the fixture server and its connection.
	worker ankiclient.Backend
	// collections controls the open collection lifetime.
	collections ankicol.BackendCollectionServiceClient
	// decks exposes native deck operations.
	decks decks.DecksServiceClient
	// types exposes native note-type operations.
	types notetypes.NotetypesServiceClient
	// notes exposes native note operations.
	notes notes.NotesServiceClient
	// cards exposes generated cards and their supported scheduling state.
	cards cards.CardsServiceClient
	// search locates native note and card identities.
	search search.SearchServiceClient
	// scheduler creates real review events for preservation workflows.
	scheduler scheduler.SchedulerServiceClient
	// stats reads review history through supported methods.
	stats stats.StatsServiceClient
	// presets creates synthetic presets and inspects their preserved native values.
	presets deck_config.DeckConfigServiceClient
}

// close closes the collection before stopping its server, even after cancellation.
func (client *fixtureClient) close() error {
	if client.worker == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var collectionErr error
	if client.opened {
		_, collectionErr = client.collections.CloseCollection(ctx, &ankicol.CloseCollectionRequest{})
		client.opened = false
	}
	workerErr := client.worker.Stop(ctx)
	client.worker = nil
	if collectionErr != nil {
		collectionErr = fmt.Errorf("close fixture collection: %w", collectionErr)
	}
	if workerErr != nil {
		workerErr = fmt.Errorf("stop fixture server: %w", workerErr)
	}
	return errors.Join(collectionErr, workerErr)
}

// fixtureCLIPath locates the explicitly declared application executable for native fixtures.
func fixtureCLIPath() string {
	return filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
}

// fixtureApplication configures application workflows with the linked CLI's server command.
func fixtureApplication() *App {
	return New(Options{BackendBinary: fixtureCLIPath(), BackendArguments: []string{"server", "serve"}})
}

// connectFixture opens an existing or new private collection with generated clients.
func connectFixture(t *testing.T, database string) *fixtureClient {
	t.Helper()
	worker := ankiclient.New(ankiclient.BackendSTDOptions{Binary: fixtureCLIPath(), Arguments: []string{"server", "serve"}})
	if err := worker.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	connection, err := worker.Connection()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if stopErr := worker.Stop(ctx); stopErr != nil {
			t.Error(stopErr)
		}
		t.Fatal(err)
	}
	client := &fixtureClient{
		worker:      worker,
		collections: ankicol.NewBackendCollectionServiceClient(connection),
		decks:       decks.NewDecksServiceClient(connection),
		types:       notetypes.NewNotetypesServiceClient(connection),
		notes:       notes.NewNotesServiceClient(connection),
		cards:       cards.NewCardsServiceClient(connection),
		search:      search.NewSearchServiceClient(connection),
		scheduler:   scheduler.NewSchedulerServiceClient(connection),
		stats:       stats.NewStatsServiceClient(connection),
		presets:     deck_config.NewDeckConfigServiceClient(connection),
	}
	t.Cleanup(func() {
		if err := client.close(); err != nil {
			t.Error(err)
		}
	})
	if _, err := client.collections.OpenCollection(context.Background(), &ankicol.OpenCollectionRequest{
		CollectionPath: database, MediaFolderPath: filepath.Join(filepath.Dir(database), "collection.media"), MediaDbPath: filepath.Join(filepath.Dir(database), "collection.media.db2"),
	}); err != nil {
		t.Fatal(err)
	}
	client.opened = true
	return client
}

// inspectFixture restores through Anki so observing a package cannot change it.
func inspectFixture(t *testing.T, input string) *fixtureClient {
	t.Helper()
	database := filepath.Join(t.TempDir(), "collection.anki2")
	client := connectFixture(t, database)
	if _, err := client.collections.CloseCollection(context.Background(), &ankicol.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	client.opened = false
	conn, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(conn).ImportCollectionPackage(context.Background(), &import_export.ImportCollectionPackageRequest{
		ColPath: database, BackupPath: input, MediaFolder: filepath.Join(filepath.Dir(database), "collection.media"), MediaDb: filepath.Join(filepath.Dir(database), "collection.media.db2"),
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.collections.OpenCollection(context.Background(), &ankicol.OpenCollectionRequest{
		CollectionPath: database, MediaFolderPath: filepath.Join(filepath.Dir(database), "collection.media"), MediaDbPath: filepath.Join(filepath.Dir(database), "collection.media.db2"),
	}); err != nil {
		t.Fatal(err)
	}
	client.opened = true
	return client
}

// exportFixture writes an open fixture through Anki and records its consuming lifetime.
func (client *fixtureClient) exportFixture(t *testing.T, output string) {
	t.Helper()
	conn, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	client.opened = false
	if _, err := import_export.NewBackendImportExportServiceClient(conn).ExportCollectionPackage(context.Background(), &import_export.ExportCollectionPackageRequest{
		OutPath: output, IncludeMedia: true,
	}); err != nil {
		t.Fatal(err)
	}
}

// fixtureData identifies the API-created resources used by the CLI workflows.
type fixtureData struct {
	// base is the original public synthetic archive.
	base string
	// database is its closed standalone collection.
	database string
	// noteID is allocated by Anki for the representative note.
	noteID int64
	// typeID is allocated by Anki for the Chinese note type.
	typeID int64
	// cardIDs identifies its two generated cards.
	cardIDs []int64
}

// newFixture creates a reviewed note, templates, and an archive through Anki's API.
func newFixture(t *testing.T) fixtureData {
	t.Helper()
	root := t.TempDir()
	fixture := fixtureData{base: filepath.Join(root, "base.colpkg"), database: filepath.Join(root, "collection.anki2")}
	client := connectFixture(t, fixture.database)
	ctx := context.Background()
	deck, err := client.decks.GetDeck(ctx, &decks.DeckId{Did: 1})
	if err != nil {
		t.Fatal(err)
	}
	deck.Name = "grammar"
	if _, err := client.decks.UpdateDeck(ctx, deck); err != nil {
		t.Fatal(err)
	}
	definition := &notetypes.Notetype{
		Name: "Chinese", Config: &notetypes.Notetype_Config{Css: ".card { font-family: sans-serif; }"},
		Fields: []*notetypes.Notetype_Field{
			{Name: "hanzi", Config: &notetypes.Notetype_Field_Config{}},
			{Name: "article-title", Config: &notetypes.Notetype_Field_Config{}},
			{Name: "english", Config: &notetypes.Notetype_Field_Config{}},
		},
		Templates: []*notetypes.Notetype_Template{
			{Name: "Forward", Config: &notetypes.Notetype_Template_Config{QFormat: "{{hanzi}}", AFormat: "{{english}}"}},
			{Name: "Reverse", Config: &notetypes.Notetype_Template_Config{QFormat: "{{english}}", AFormat: "{{hanzi}}"}},
		},
	}
	created, err := client.types.AddNotetype(ctx, definition)
	if err != nil {
		t.Fatal(err)
	}
	fixture.typeID = created.Id
	added, err := client.notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{
		Guid: "fixture-note", NotetypeId: fixture.typeID, Fields: []string{"<b>你好</b>", "Greeting", "hello\n"}, Tags: []string{"A1"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	fixture.noteID = added.NoteId
	found, err := client.search.SearchCards(ctx, &search.SearchRequest{Search: fmt.Sprintf("nid:%d", fixture.noteID)})
	if err != nil || len(found.GetIds()) != 2 {
		t.Fatalf("native card generation: %v", err)
	}
	fixture.cardIDs = found.Ids
	states, err := client.scheduler.GetSchedulingStates(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.scheduler.AnswerCard(ctx, &scheduler.CardAnswer{
		CardId: fixture.cardIDs[0], CurrentState: states.Current, NewState: states.Good,
		Rating: scheduler.CardAnswer_GOOD, AnsweredAtMillis: time.Now().UnixMilli(), MillisecondsTaken: 250,
	}); err != nil {
		t.Fatal(err)
	}
	history, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil || len(history.GetEntries()) == 0 {
		t.Fatalf("native review creation: %v", err)
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	packageFixture(t, fixture.database, fixture.base)
	return fixture
}

// packageFixture creates a native package from a closed collection and public media.
func packageFixture(t *testing.T, database, output string) {
	t.Helper()
	client := connectFixture(t, database)
	conn, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := media.NewMediaServiceClient(conn).AddMediaFile(context.Background(), &media.AddMediaFileRequest{
		DesiredName: "example.png", Data: []byte("fixture media"),
	}); err != nil {
		t.Fatal(err)
	}
	client.exportFixture(t, output)
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
}

// assertFixtureMedia verifies native media paths and bytes after package reimport.
func assertFixtureMedia(t *testing.T, client *fixtureClient) {
	t.Helper()
	conn, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	path, err := media.NewMediaServiceClient(conn).GetAbsoluteMediaPath(context.Background(), &generic.String{Val: "example.png"})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path.Val)
	if err != nil || string(data) != "fixture media" {
		t.Fatalf("native media changed: %q %v", data, err)
	}
}

// readFixtureTOML parses fixture declarations through their protobuf schema.
func readFixtureTOML(t *testing.T, path string, message proto.Message) {
	t.Helper()
	var values map[string]any
	if _, err := toml.DecodeFile(path, &values); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(values)
	if err != nil {
		t.Fatal(err)
	}
	if err := protojson.Unmarshal(data, message); err != nil {
		t.Fatal(err)
	}
}

// writeFixtureTOML records editable declarations with ordinary TOML syntax.
func writeFixtureTOML(t *testing.T, path string, message proto.Message) {
	t.Helper()
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	var values map[string]any
	if err := json.Unmarshal(data, &values); err != nil {
		t.Fatal(err)
	}
	var content strings.Builder
	if err := toml.NewEncoder(&content).Encode(values); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content.String()), 0o600); err != nil {
		t.Fatal(err)
	}
}

// fixtureNotePath locates the representative note in the default export layout.
func fixtureNotePath(root string) string {
	return filepath.Join(root, "decks", "grammar", "notes", "你好.note.anki.toml")
}

// fixturePlan parses a native reconciliation plan and fails on workflow errors.
func fixturePlan(t *testing.T, input, text string) *ankiascodepb.CollectionChangePlan {
	t.Helper()
	plan, err := planForFixture(context.Background(), input, text, "")
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

// assertConverged verifies the original declarations describe the resulting collection.
func assertConverged(t *testing.T, input, text string) {
	t.Helper()
	if plan := fixturePlan(t, input, text); len(plan.Changes) != 0 {
		t.Fatalf("not converged: %v", plan.Changes)
	}
}

// fixtureCardSchedule selects supported study values independently of mutable deck metadata.
func fixtureCardSchedule(card *cards.Card) *cards.Card {
	value := proto.Clone(card).(*cards.Card)
	value.MtimeSecs, value.Usn, value.DeckId, value.TemplateIdx = 0, 0, 0, 0
	return value
}

// runFixtureCLI executes the built command with bounded lifetime and separate streams.
func runFixtureCLI(t *testing.T, arguments ...string) []byte {
	t.Helper()
	binary := fixtureCLIPath()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, binary, arguments...)
	var stdout, stderr strings.Builder
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("CLI %s: %v: %s", arguments[0], err, stderr.String())
	}
	operation := arguments
	if operation[0] == "col" || operation[0] == "collection" {
		operation = operation[1:]
	}
	if operation[0] != "plan" && !(len(operation) > 1 && operation[0] == "state" && operation[1] == "dump") && stdout.Len() != 0 {
		t.Fatalf("status polluted stdout for %s", arguments[0])
	}
	return []byte(stdout.String())
}

// saveFixtureArtifact emits a repeatable public collection or text result for inspection.
func saveFixtureArtifact(t *testing.T, input, name string) {
	t.Helper()
	root := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if root == "" {
		return
	}
	data, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, name), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// exportForFixture invokes the application workflow with the fixture namespace.
func exportForFixture(ctx context.Context, input, output, key string) error {
	application := fixtureApplication()
	if err := application.ColExport(ctx, ankiascodeapi.CollectionPath(input), output, "", key); err != nil {
		return fmt.Errorf("export fixture: %w", err)
	}
	return nil
}

// buildForFixture invokes the application workflow with the fixture namespace.
func buildForFixture(ctx context.Context, base, dir, output, key string) error {
	application := fixtureApplication()
	if err := application.ColBuild(ctx, base, dir, output, key); err != nil {
		return fmt.Errorf("build fixture: %w", err)
	}
	return nil
}

// applyForFixture invokes the application workflow with the fixture namespace.
func applyForFixture(ctx context.Context, input, dir, key string) error {
	application := fixtureApplication()
	if err := application.ColApply(ctx, ankiascodeapi.CollectionPath(input), dir, key); err != nil {
		return fmt.Errorf("apply fixture: %w", err)
	}
	return nil
}

// planForFixture invokes the application workflow with the fixture namespace.
func planForFixture(ctx context.Context, input, dir, key string) (*ankiascodepb.CollectionChangePlan, error) {
	application := fixtureApplication()
	value, err := application.ColPlan(ctx, ankiascodeapi.CollectionPath(input), dir, key)
	if err != nil {
		return nil, fmt.Errorf("plan fixture: %w", err)
	}
	return value, nil
}
