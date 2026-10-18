package ankiascode

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"github.com/ankitects/anki/proto/anki/decks"
	"google.golang.org/protobuf/proto"
)

// TestExportStateOnboardingWorkflow exercises whole-collection binding through the native CLI.
// Export must preserve source data, use selected paths, and converge after state installation.
func TestExportStateOnboardingWorkflow(t *testing.T) {
	for _, inputFormat := range []string{"package", "database"} {
		t.Run(inputFormat, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if inputFormat == "database" {
				input = fixture.database
			}
			before := stateFixtureDump(t, input)
			root := t.TempDir()
			text, stateFile := filepath.Join(root, "configs"), filepath.Join(root, "state.anki.toml")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text, "--state-output", stateFile)
			if !proto.Equal(before, stateFixtureDump(t, input)) {
				t.Fatal("state export changed stored collection state")
			}
			data, err := os.ReadFile(stateFile)
			if err != nil {
				t.Fatal(err)
			}
			generated := &ankiascodepb.State{}
			if err := codec.Decode(data, codec.TOML, generated); err != nil {
				t.Fatal(err)
			}
			notePath := "decks/grammar/notes/你好.note.anki.toml"
			if generated.GetResources().GetNotes()[notePath].GetId() != fixture.noteID ||
				generated.GetResources().GetDecks()["decks/grammar/deck.anki.toml"].GetId() != 1 ||
				generated.GetResources().GetNoteTypes()["note_types/Chinese.notetype.anki.toml"].GetId() != fixture.typeID ||
				len(generated.GetResources().GetPresets()) == 0 {
				t.Fatal("generated bindings do not identify exported native objects", generated)
			}
			definition := generated.Resources.NoteTypes["note_types/Chinese.notetype.anki.toml"]
			if len(definition.Fields) == 0 || len(definition.Templates) == 0 {
				t.Fatal("missing native member bindings")
			}
			if len(stateFixtureDump(t, input).GetResources().GetNotes()) != 0 {
				t.Fatal("export installed state")
			}
			runFixtureCLI(t, "col", "state", "update", "--collection", input, stateFile)
			if !proto.Equal(generated, stateFixtureDump(t, input)) {
				t.Fatal("generated state did not round-trip through update")
			}
			assertConverged(t, input, text)
			plan := runFixtureCLI(t, "col", "plan", "--collection", input, "--text", text)
			planFile := filepath.Join(root, "plan.json")
			if err := os.WriteFile(planFile, plan, 0o600); err != nil {
				t.Fatal(err)
			}
			// A second export discovers new bindings while preserving a custom namespace's names.
			key := "export-state/custom"
			custom := proto.Clone(generated).(*ankiascodepb.State)
			customPath := "decks/grammar/notes/custom.note.anki.toml"
			custom.Resources.Notes[customPath] = custom.Resources.Notes[notePath]
			delete(custom.Resources.Notes, notePath)
			for _, preset := range custom.Resources.Presets {
				custom.Resources.Presets = map[string]*ankiascodepb.PresetState{"study": preset}
				break
			}
			writeFixtureTOML(t, stateFile, custom)
			runFixtureCLI(t, "col", "state", "replace", "--collection", input, "--state-key", key, stateFile)
			selected := filepath.Join(root, "selected")
			selectedState := filepath.Join(selected, "state.anki.toml")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", selected, "--state-output", selectedState, "--state-key", key)
			selectedData, err := os.ReadFile(selectedState)
			if err != nil {
				t.Fatal(err)
			}
			selectedBindings := &ankiascodepb.State{}
			if err := codec.Decode(selectedData, codec.TOML, selectedBindings); err != nil || !proto.Equal(custom, selectedBindings) {
				t.Fatalf("selected bindings differ from generated paths: %v\n%s", err, selectedData)
			}
			if _, err := config.New(selected).Load(); err != nil {
				t.Fatal("state output broke config loading", err)
			}
			saveFixtureArtifact(t, selectedState, "onboarding-state-"+inputFormat+".toml")
			saveFixtureArtifact(t, filepath.Join(selected, filepath.FromSlash(customPath)), "onboarding-note-"+inputFormat+".toml")
			saveFixtureArtifact(t, planFile, "onboarding-plan-"+inputFormat+".json")
		})
	}
}

// TestExportStateDestinationWorkflow rejects destructive destinations and retains prior state on failure.
func TestExportStateDestinationWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	text := filepath.Join(root, "configs")
	application := fixtureApplication()
	for _, destination := range []string{fixture.base, filepath.Join(text, "collection.anki.toml"), filepath.Join(text, "collection.anki.toml", "state.toml"), filepath.Join(text, "decks"), filepath.Join(text, "note_types", "Chinese.notetype.anki.toml")} {
		if err := application.ColExport(t.Context(), ankiascodeapi.CollectionPath(fixture.base), text, destination, ""); err == nil {
			t.Fatalf("accepted unsafe state destination %q", destination)
		}
		if _, err := os.Stat(text); !os.IsNotExist(err) {
			t.Fatal("rejected destination wrote config output", err)
		}
	}
	for _, linkKind := range []string{"symlink", "hardlink"} {
		alias := filepath.Join(root, linkKind)
		link := os.Link
		if linkKind == "symlink" {
			link = os.Symlink
		}
		if err := link(fixture.base, alias); err != nil {
			t.Fatal(err)
		}
		if err := application.ColExport(t.Context(), ankiascodeapi.CollectionPath(fixture.base), text, alias, ""); err == nil {
			t.Fatal("accepted source alias", linkKind)
		}
	}
	for _, scenario := range []string{"dangling file", "dangling directory"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			text := filepath.Join(root, "configs")
			if err := os.Mkdir(text, 0o755); err != nil {
				t.Fatal(err)
			}
			stateFile := filepath.Join(root, "state.anki.toml")
			linkPath, target := filepath.Join(text, "collection.anki.toml"), stateFile
			if scenario == "dangling directory" {
				target = filepath.Join(root, "external-types")
				linkPath = filepath.Join(text, "note_types")
				stateFile = filepath.Join(target, "Chinese.notetype.anki.toml")
			}
			if err := os.Symlink(target, linkPath); err != nil {
				t.Fatal(err)
			}
			if err := application.ColExport(t.Context(), ankiascodeapi.CollectionPath(fixture.base), text, stateFile, ""); err == nil {
				t.Fatal("accepted dangling config alias to state output")
			}
			if _, err := os.Stat(stateFile); !os.IsNotExist(err) {
				t.Fatal("rejected dangling alias wrote state destination", err)
			}
			if actual, err := os.Readlink(linkPath); err != nil || actual != target {
				t.Fatal("rejected export changed symlink", err)
			}
			if entries, err := os.ReadDir(text); err != nil || len(entries) != 1 {
				t.Fatal("rejected export wrote configs", err)
			}
		})
	}
	stateFile := filepath.Join(root, "state.anki.toml")
	if err := os.WriteFile(stateFile, []byte("previous state"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(text, []byte("blocks config directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := application.ColExport(t.Context(), ankiascodeapi.CollectionPath(fixture.base), text, stateFile, ""); err == nil {
		t.Fatal("accepted blocked config output")
	}
	if data, err := os.ReadFile(stateFile); err != nil || string(data) != "previous state" {
		t.Fatal("failed config export changed state output", err)
	}
}

// TestExportConfigOnlyWorkflow verifies config-only output, read-only export and independent state dumps.
func TestExportConfigOnlyWorkflow(t *testing.T) {
	for _, inputFormat := range []string{"package", "database"} {
		t.Run(inputFormat, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if inputFormat == "database" {
				input = fixture.database
			}
			before, err := fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			text := filepath.Join(t.TempDir(), "text")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			after, err := fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("export modified source collection content")
			}
			initial := stateFixtureDump(t, input)
			if len(initial.GetResources().GetDecks())+len(initial.GetResources().GetNotes())+len(initial.GetResources().GetNoteTypes()) != 0 {
				t.Fatal("export initialized source state")
			}
			for _, format := range []string{"toml", "json"} {
				file := filepath.Join(text, "state.anki."+format)
				if _, err := os.Stat(file); !os.IsNotExist(err) {
					t.Fatalf("export generated %s state: %v", format, err)
				}
				if err := os.WriteFile(file, []byte("caller state"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			for _, format := range []string{"toml", "json"} {
				if data, err := os.ReadFile(filepath.Join(text, "state.anki."+format)); err != nil || string(data) != "caller state" {
					t.Fatalf("export changed caller-owned %s state: %v", format, err)
				}
			}
			runFixtureCLI(t, "col", "state", "import", "--collection", input, "--text", text,
				"decks/grammar/deck.anki.toml=1",
				fmt.Sprintf("note_types/Chinese.notetype.anki.toml=%d", fixture.typeID),
				fmt.Sprintf("decks/grammar/notes/你好.note.anki.toml=%d", fixture.noteID))
			stored := stateFixtureDump(t, input)
			if stored.GetResources().GetNotes()["decks/grammar/notes/你好.note.anki.toml"].GetId() != fixture.noteID {
				t.Fatal("separate state setup did not bind the native note")
			}
			for _, format := range []string{"toml", "json"} {
				data := runFixtureCLI(t, "col", "state", "dump", "--format", format, "--collection", input)
				value := &ankiascodepb.State{}
				if err := codec.Decode(data, codec.Format(format), value); err != nil || !proto.Equal(value, stored) {
					t.Fatalf("separate %s state dump: %v", format, err)
				}
				file := filepath.Join(t.TempDir(), "state.anki."+format)
				if err := os.WriteFile(file, data, 0o600); err != nil {
					t.Fatal(err)
				}
				saveFixtureArtifact(t, file, "dumped-state-"+inputFormat+"."+format)
			}
			key := "export-fixture/custom"
			custom := proto.Clone(stored).(*ankiascodepb.State)
			path := "decks/grammar/notes/custom.note.anki.toml"
			custom.Resources.Notes[path] = custom.Resources.Notes["decks/grammar/notes/你好.note.anki.toml"]
			delete(custom.Resources.Notes, "decks/grammar/notes/你好.note.anki.toml")
			stateFile := filepath.Join(t.TempDir(), "custom-state.anki.toml")
			writeFixtureTOML(t, stateFile, custom)
			runFixtureCLI(t, "col", "state", "replace", "--collection", input, "--state-key", key, stateFile)
			beforeExport := fixtureSources(t, text)
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			assertFixtureSources(t, text, beforeExport)
			selected := filepath.Join(t.TempDir(), "selected")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", selected, "--state-key", key)
			if _, err := os.Stat(filepath.Join(selected, filepath.FromSlash(path))); err != nil {
				t.Fatal("export ignored the selected state namespace", err)
			}
			if !proto.Equal(custom, stateFixtureDumpKey(t, input, key)) || !proto.Equal(stored, stateFixtureDump(t, input)) {
				t.Fatal("export changed stored bindings")
			}
			configs, err := config.New(selected).Load()
			if err != nil || len(configs.GetDecks()) == 0 || len(configs.GetNoteTypes()) == 0 {
				t.Fatalf("exported hierarchy is not loadable: %v", err)
			}
			saveFixtureArtifact(t, filepath.Join(selected, filepath.FromSlash(path)), "exported-note-"+inputFormat+".toml")
		})
	}
}

// TestExportWritesWorkflow verifies in-place output and caller-owned files/directories.
func TestExportWritesWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	t.Setenv("HOME", root)
	text := filepath.Join(root, "text")
	reference := filepath.Join(root, "permissions")
	if err := os.Mkdir(reference, 0o755); err != nil {
		t.Fatal(err)
	}
	defaultMode, err := os.Stat(reference)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), fixture.base, "~/text", ""); err != nil {
		t.Fatal(err)
	}
	created, err := os.Stat(text)
	if err != nil || created.Mode().Perm() != defaultMode.Mode().Perm() {
		t.Fatalf("new export directory permissions differ from normal creation: %v: %v", created, err)
	}
	if err := os.Chmod(text, 0o775); err != nil {
		t.Fatal(err)
	}
	unrelated := filepath.Join(text, "README.txt")
	if err := os.WriteFile(unrelated, []byte("keep me"), 0o600); err != nil {
		t.Fatal(err)
	}
	private := filepath.Join(text, "private")
	if err := os.Mkdir(private, 0o700); err != nil {
		t.Fatal(err)
	}
	privateFile := filepath.Join(private, "personal.txt")
	if err := os.WriteFile(privateFile, []byte("unrelated private content"), 0o600); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(text, "stale.note.anki.toml")
	if err := os.WriteFile(stale, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	staleDeck := filepath.Join(text, "decks", "Removed")
	marker := filepath.Join(staleDeck, "decks", "Child", "deck.anki.toml")
	writeFixtureTOML(t, marker, &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
	writeFixtureTOML(t, filepath.Join(filepath.Dir(marker), "notes", "old.note.anki.toml"), &ankiascodepb.NoteConfig{})
	staleType := filepath.Join(text, "note_types", "Removed.notetype.anki.toml")
	writeFixtureTOML(t, staleType, &ankiascodepb.NoteTypeConfig{Name: "Removed"})
	if err := exportForFixture(context.Background(), fixture.base, "~/text", ""); err != nil {
		t.Fatal(err)
	}
	refreshed, err := os.Stat(text)
	if err != nil || refreshed.Mode().Perm() != 0o775 {
		t.Fatalf("export refresh changed directory permissions: %v: %v", refreshed, err)
	}
	if data, err := os.ReadFile(stale); err != nil || string(data) != "stale" {
		t.Fatal("export changed an unselected declaration")
	}
	for _, retained := range []string{marker, staleType} {
		if _, err := os.Stat(retained); err != nil {
			t.Fatalf("export removed an unselected resource at %q: %v", retained, err)
		}
	}
	data, err := os.ReadFile(unrelated)
	if err != nil || string(data) != "keep me" {
		t.Fatal("unrelated file changed")
	}
	for path, mode := range map[string]os.FileMode{private: 0o700, privateFile: 0o600, unrelated: 0o600} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != mode {
			t.Errorf("export changed private permissions at %q: %v: %v", path, info, err)
		}
	}
	privateContent, err := os.ReadFile(privateFile)
	if err != nil || string(privateContent) != "unrelated private content" {
		t.Fatal("unrelated private content changed")
	}
	configPath := filepath.Join(text, "collection.anki.toml")
	before, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), filepath.Join(root, "missing.colpkg"), text, ""); err == nil {
		t.Fatal("invalid input accepted")
	}
	after, err := os.ReadFile(configPath)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed export changed existing declarations")
	}
	fresh := filepath.Join(root, "fresh")
	if err := exportForFixture(context.Background(), fixture.base, fresh, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := config.New(fresh).Load(); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(fixture.database)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), fixture.database, filepath.Dir(fixture.database), ""); err == nil {
		t.Fatal("export accepted a destination containing its source")
	}
	if after, err := os.ReadFile(fixture.database); err != nil || string(after) != string(original) {
		t.Fatal("rejected export changed its source database")
	}
	for _, scenario := range []string{"source symlink", "destination symlink"} {
		t.Run(scenario, func(t *testing.T) {
			sourceDirectory := filepath.Dir(fixture.database)
			link := filepath.Join(t.TempDir(), "alias")
			if err := os.Symlink(sourceDirectory, link); err != nil {
				t.Fatal(err)
			}
			input, output := fixture.database, sourceDirectory
			if scenario == "source symlink" {
				input = filepath.Join(link, filepath.Base(fixture.database))
			} else {
				output = link
			}
			before := fixtureSources(t, sourceDirectory)
			if err := exportForFixture(context.Background(), input, output, ""); err == nil {
				t.Fatal("export accepted a destination physically containing its source")
			}
			assertFixtureSources(t, sourceDirectory, before)
		})
	}

	saveFixtureArtifact(t, configPath, "exported.anki.toml")
}
