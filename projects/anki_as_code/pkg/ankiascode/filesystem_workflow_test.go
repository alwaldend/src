package ankiascode

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
)

// TestBasePathConfigurationBoundary keeps config roots and declared files inside injected storage.
func TestBasePathConfigurationBoundary(t *testing.T) {
	for _, boundary := range []string{"parent-root", "declared-parent", "root-symlink", "declared-symlink"} {
		t.Run(boundary, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), "base")
			outside := root + "-sibling"
			for _, directory := range []string{root, outside} {
				if err := os.Mkdir(directory, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			original := []byte("# outside storage must remain unchanged\nversion = \"1\"\n")
			marker := filepath.Join(outside, "collection.anki.toml")
			if err := os.WriteFile(marker, original, 0o600); err != nil {
				t.Fatal(err)
			}
			storage := afero.NewBasePathFs(afero.NewOsFs(), root)
			application := New(Options{FS: storage})
			output := "config"
			declared := "notes/new.note.anki.toml"
			switch boundary {
			case "parent-root":
				output = "../base-sibling"
			case "declared-parent":
				declared = "../../base-sibling/collection.anki.toml"
			case "root-symlink":
				output = "escape"
				if err := os.Symlink(outside, filepath.Join(root, output)); err != nil {
					t.Fatal(err)
				}
			case "declared-symlink":
				if err := storage.Mkdir(output, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, filepath.Join(root, output, "escape")); err != nil {
					t.Fatal(err)
				}
				declared = "escape/collection.anki.toml"
			}
			configs := &ankiascodepb.CollectionConfigSet{
				Config: &ankiascodepb.CollectionConfig{Version: 1},
				Notes:  []*ankiascodepb.NoteConfigFile{{Path: declared, Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "boundary"}}}},
			}
			if err := application.ConfigExport(t.Context(), output, configs); !errors.Is(err, os.ErrPermission) {
				t.Errorf("config export did not reject an external path: %v", err)
			}
			if data, err := os.ReadFile(marker); err != nil || string(data) != string(original) {
				t.Errorf("config export changed outside storage: %q %v", data, err)
			}
			if entries, err := os.ReadDir(outside); err != nil || len(entries) != 1 {
				t.Errorf("config export created outside entries: %v %v", entries, err)
			}
			if _, err := application.ConfigDump(t.Context(), "../base-sibling", "json"); !errors.Is(err, os.ErrPermission) {
				t.Errorf("config dump did not reject an external root: %v", err)
			}
		})
	}
}

// TestBasePathConfigurationNormalization accepts parent components that stay within storage.
func TestBasePathConfigurationNormalization(t *testing.T) {
	storage := afero.NewBasePathFs(afero.NewOsFs(), t.TempDir())
	application := New(Options{FS: storage})
	configs := &ankiascodepb.CollectionConfigSet{
		Config: &ankiascodepb.CollectionConfig{Version: 1},
		Notes:  []*ankiascodepb.NoteConfigFile{{Path: "notes/unused/../hello.note.anki.toml", Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "normalized"}}}},
	}
	if err := application.ConfigExport(t.Context(), "unused/../config", configs); err != nil {
		t.Fatal(err)
	}
	data, err := application.ConfigDump(t.Context(), "unused/../config", "json")
	if err != nil {
		t.Fatal(err)
	}
	actual := &ankiascodepb.CollectionConfigSet{}
	if err := codec.Decode(data, codec.JSON, actual); err != nil {
		t.Fatal(err)
	}
	expected := proto.Clone(configs).(*ankiascodepb.CollectionConfigSet)
	expected.Notes[0].Path = "notes/hello.note.anki.toml"
	if !proto.Equal(actual, expected) {
		t.Fatalf("normalization changed the configuration: %v", actual)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "basepath-normalized-configuration.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestBasePathExportBoundary rejects backing-store symlinks outside the injected namespace.
func TestBasePathExportBoundary(t *testing.T) {
	for _, targetKind := range []string{"absolute", "relative"} {
		t.Run(targetKind, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), "base")
			outside := root + "-sibling"
			for _, directory := range []string{root, outside} {
				if err := os.Mkdir(directory, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			storage := afero.NewBasePathFs(afero.NewOsFs(), root)
			if err := afero.WriteFile(storage, "source.anki2", []byte("unchanged"), 0o600); err != nil {
				t.Fatal(err)
			}
			target := outside
			if targetKind == "relative" {
				target = "../base-sibling"
			}
			if err := os.Symlink(target, filepath.Join(root, "escape")); err != nil {
				t.Fatal(err)
			}
			application := New(Options{FS: storage})
			err := application.ColExport(t.Context(), ankiascodeapi.CollectionOptions{Path: "/source.anki2"}, "/escape/config", "")
			if !errors.Is(err, os.ErrPermission) {
				t.Fatalf("external symlink was not rejected at the filesystem boundary: %v", err)
			}
			if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
				t.Fatalf("export wrote beyond the injected base: %v %v", entries, err)
			}
		})
	}
}

// TestInjectedExportSourceConflict rejects state output aliasing a source before backend startup.
func TestInjectedExportSourceConflict(t *testing.T) {
	for _, source := range []string{"/collection.anki2", "collection.anki2"} {
		t.Run(source, func(t *testing.T) {
			storage := afero.NewMemMapFs()
			original := []byte("collection source must remain unchanged")
			if err := afero.WriteFile(storage, source, original, 0o600); err != nil {
				t.Fatal(err)
			}
			application := New(Options{FS: storage})
			output := filepath.Join(filepath.Dir(source), "configuration")
			stateOutput := filepath.Join(filepath.Dir(source), "unused", "..", filepath.Base(source))
			err := application.ColExport(t.Context(), ankiascodeapi.CollectionOptions{Path: source}, output, stateOutput)
			if err == nil || !strings.Contains(err.Error(), "state output must be a regular file distinct from the source collection") {
				t.Fatalf("source conflict reached backend startup: %v", err)
			}
			data, err := afero.ReadFile(storage, source)
			if err != nil || string(data) != string(original) {
				t.Fatalf("source conflict changed collection: %q %v", data, err)
			}
			if _, err := storage.Stat(output); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("source conflict wrote config output: %v", err)
			}
		})
	}
}

// TestInjectedOverlayConfiguration merges base and overlay files while retaining base contents.
func TestInjectedOverlayConfiguration(t *testing.T) {
	base, layer := afero.NewMemMapFs(), afero.NewMemMapFs()
	root := "configuration"
	baseConfigs := &ankiascodepb.CollectionConfigSet{
		Config: &ankiascodepb.CollectionConfig{Version: 1},
		Notes: []*ankiascodepb.NoteConfigFile{
			{Path: "notes/base.note.anki.toml", Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "base only"}}},
			{Path: "notes/shared.note.anki.toml", Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "original"}}},
		},
	}
	if err := New(Options{FS: base}).ConfigExport(t.Context(), root, baseConfigs); err != nil {
		t.Fatal(err)
	}
	originals := map[string][]byte{}
	for _, path := range []string{"collection.anki.toml", "notes/base.note.anki.toml", "notes/shared.note.anki.toml"} {
		data, err := afero.ReadFile(base, filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		originals[path] = data
	}
	layerConfigs := &ankiascodepb.CollectionConfigSet{
		Config: &ankiascodepb.CollectionConfig{Version: 1},
		Notes: []*ankiascodepb.NoteConfigFile{
			{Path: "notes/layer.note.anki.toml", Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "layer only"}}},
			{Path: "notes/shared.note.anki.toml", Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "updated"}}},
		},
	}
	application := New(Options{FS: afero.NewCopyOnWriteFs(base, layer)})
	if err := application.ConfigExport(t.Context(), root, layerConfigs); err != nil {
		t.Fatal(err)
	}
	data, err := application.ConfigDump(t.Context(), root, "json")
	if err != nil {
		t.Fatal(err)
	}
	actual := &ankiascodepb.CollectionConfigSet{}
	if err := codec.Decode(data, codec.JSON, actual); err != nil {
		t.Fatal(err)
	}
	expected := &ankiascodepb.CollectionConfigSet{
		Config: baseConfigs.Config,
		Notes:  []*ankiascodepb.NoteConfigFile{baseConfigs.Notes[0], layerConfigs.Notes[0], layerConfigs.Notes[1]},
	}
	if !proto.Equal(actual, expected) {
		t.Errorf("overlay dump lost base files or layer precedence: %v", actual)
	}
	for path, original := range originals {
		if current, err := afero.ReadFile(base, filepath.Join(root, path)); err != nil || string(current) != string(original) {
			t.Errorf("overlay changed base file %s: %q %v", path, current, err)
		}
	}
	if _, err := base.Stat(filepath.Join(root, "notes/layer.note.anki.toml")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("overlay created a base file: %v", err)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "overlay-configuration.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestInjectedConfigurationWorkflow exercises complete file-only operations without host storage.
func TestInjectedConfigurationWorkflow(t *testing.T) {
	for _, namespace := range []string{"absolute", "relative", "literal-home"} {
		t.Run(namespace, func(t *testing.T) {
			host := filepath.Join(t.TempDir(), "configuration")
			root := host
			switch namespace {
			case "relative":
				root = "configuration"
			case "literal-home":
				t.Setenv("HOME", "")
				root = "~/configuration"
			}
			storage := afero.NewMemMapFs()
			application := New(Options{FS: storage})
			configs := &ankiascodepb.CollectionConfigSet{
				Config:    &ankiascodepb.CollectionConfig{Version: 1},
				NoteTypes: []*ankiascodepb.NoteTypeConfigFile{{Path: "note_types/Basic.notetype.anki.toml", Config: &ankiascodepb.NoteTypeConfig{Name: "Basic"}}},
				Decks: []*ankiascodepb.DeckConfigSet{
					{Path: "decks/Lesson/deck.anki.toml", Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}}},
					{Path: "decks/Lesson/decks/Child/deck.anki.toml", Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}}, Notes: []*ankiascodepb.NoteConfigFile{{Path: "decks/Lesson/decks/Child/notes/你好.note.anki.toml", Config: &ankiascodepb.NoteConfig{NoteType: "note_types/Basic.notetype.anki.toml", Cards: []*ankiascodepb.CardConfig{{}}, Fields: map[string]string{"Front": "你好\nsecond line"}}}}},
				},
				Notes: []*ankiascodepb.NoteConfigFile{{Path: "notes/root.note.anki.toml", Config: &ankiascodepb.NoteConfig{Fields: map[string]string{"Front": "root"}}}},
			}
			if err := application.ConfigExport(t.Context(), root, configs); err != nil {
				t.Fatal(err)
			}
			expected := proto.Clone(configs).(*ankiascodepb.CollectionConfigSet)
			expected.Decks[0].Title = "Lesson"
			expected.Decks[1].Title = "Lesson::Child"
			expected.Decks[1].Notes[0].Config.Cards[0].Deck = proto.String(expected.Decks[1].Path)
			for _, input := range []string{root, filepath.Join(root, "collection.anki.toml")} {
				for _, format := range []string{"json", "toml"} {
					data, err := application.ConfigDump(t.Context(), input, format)
					if err != nil {
						t.Fatal(err)
					}
					actual := &ankiascodepb.CollectionConfigSet{}
					if err := codec.Decode(data, codec.Format(format), actual); err != nil || !proto.Equal(actual, expected) {
						t.Fatalf("%s dump differs: %v\n%v", format, err, actual)
					}
					if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
						if err := os.WriteFile(filepath.Join(artifacts, "memory-configuration-"+namespace+"."+format), data, 0o600); err != nil {
							t.Fatal(err)
						}
					}
				}
			}
			if configs.Decks[0].Title != "" || configs.Decks[1].Notes[0].Config.Cards[0].Deck != nil {
				t.Fatal("export mutated caller messages")
			}
			other := New(Options{FS: afero.NewMemMapFs()})
			if _, err := other.ConfigDump(t.Context(), root, "json"); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("another application accessed this storage: %v", err)
			}
			readonly := New(Options{FS: afero.NewReadOnlyFs(storage)})
			if _, err := readonly.ConfigDump(t.Context(), root, "json"); err != nil {
				t.Fatalf("read-only dump: %v", err)
			}
			if err := readonly.ConfigExport(t.Context(), root, configs); !errors.Is(err, os.ErrPermission) {
				t.Fatalf("read-only export lost failure: %v", err)
			}
			marker := filepath.Join(root, "collection.anki.toml")
			if err := afero.WriteFile(storage, marker, []byte("version = ["), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := application.ConfigDump(t.Context(), root, "json"); err == nil {
				t.Fatal("dump ignored a malformed fresh read")
			}
			if err := storage.Remove(marker); err != nil {
				t.Fatal(err)
			}
			if _, err := application.ConfigDump(t.Context(), root, "json"); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("missing input lost cause: %v", err)
			}
			if _, err := os.Stat(host); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("configuration escaped to disk: %v", err)
			}
		})
	}
}
