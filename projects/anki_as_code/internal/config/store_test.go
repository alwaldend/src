package config_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"github.com/BurntSushi/toml"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// configFixture identifies a file-only collection with path-addressed resources.
type configFixture struct {
	// root contains every source file for preservation checks.
	root string
	// config identifies the root collection marker.
	config string
	// note identifies one deck-owned note file.
	note string
}

// newConfigFixture writes current-schema source files without native identities.
func newConfigFixture(t *testing.T) configFixture {
	t.Helper()
	root := t.TempDir()
	f := configFixture{root: root, config: filepath.Join(root, "collection.anki.toml"), note: filepath.Join(root, "decks", "Lesson", "notes", "greeting.note.anki.toml")}
	files := map[string]proto.Message{
		f.config: &collectionpb.CollectionConfig{Version: 1, Moved: []*collectionpb.ResourceMove{{From: "decks/Old", To: "decks/Lesson"}}},
		filepath.Join(root, "decks", "Lesson", "deck.anki.toml"):      &collectionpb.DeckConfig{Preset: "Default", Normal: &decks.Deck_Normal{NewLimit: proto.Uint32(0), DesiredRetention: proto.Float32(0.91)}},
		filepath.Join(root, "note_types", "Basic.notetype.anki.toml"): &collectionpb.NoteTypeConfig{Name: "Basic", Config: &notetypes.Notetype_Config{Css: ".card { color: blue; }"}, Fields: []*notetypes.Notetype_Field{{Name: "Front", Config: &notetypes.Notetype_Field_Config{FontSize: 24}}}, Templates: []*collectionpb.TemplateConfig{{Name: "Forward", Config: &notetypes.Notetype_Template_Config{QFormat: "{{UnknownField}}", AFormat: "{{Front}}"}}}},
		f.note: &collectionpb.NoteConfig{NoteType: "note_types/Basic.notetype.anki.toml", Tags: []string{" spaced "}, Cards: []*collectionpb.CardConfig{{}}, Fields: map[string]string{"Front": "<div class=\"example\">你好</div>\nSecond line"}},
		filepath.Join(root, "notes", "root.note.anki.toml"): &collectionpb.NoteConfig{NoteType: "note_types/Basic.notetype.anki.toml", Cards: []*collectionpb.CardConfig{{}}, Fields: map[string]string{"Front": ""}},
	}
	for path, message := range files {
		writeConfigMessage(t, path, message)
	}
	writeConfigFixture(t, filepath.Join(root, "decks", "Lesson", "notes", "README.txt"), "unrelated file\n")
	return f
}

// writeConfigMessage generates valid fixture TOML from the schema through independent library codecs.
// ProtoJSON's 64-bit integer strings remain exact; production parsing accepts that mapping.
func writeConfigMessage(t *testing.T, path string, message proto.Message) {
	t.Helper()
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	var values map[string]any
	if err := json.Unmarshal(data, &values); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := toml.NewEncoder(&output).Encode(values); err != nil {
		t.Fatal(err)
	}
	writeConfigFixture(t, path, output.String())
}

// writeConfigFixture writes a complete synthetic source file at its chosen location.
func writeConfigFixture(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

// configSources captures all source bytes so read-only operations can be verified independently.
func configSources(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk fixture: %w", err)
		}
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read fixture: %w", err)
		}
		files[path] = string(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return files
}

// assertConfigSources verifies dump successes and failures cannot alter source configuration.
func assertConfigSources(t *testing.T, before, after map[string]string) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatal("source file set changed")
	}
	for path, data := range before {
		if after[path] != data {
			t.Fatalf("source %s changed", path)
		}
	}
}

// TestConfigPathWorkflow verifies source paths, readable references, moves, presence, and lossless dump formats.
func TestConfigPathWorkflow(t *testing.T) {
	fixture := newConfigFixture(t)
	before := configSources(t, fixture.root)
	for _, format := range []string{"json", "toml"} {
		data, err := config.New(fixture.config).Dump(format)
		if err != nil {
			t.Fatalf("dump %s: %v", format, err)
		}
		actual := &collectionpb.CollectionConfigSet{}
		if format == "json" {
			if err := protojson.Unmarshal(data, actual); err != nil {
				t.Fatal(err)
			}
		} else {
			var values map[string]any
			if _, err := toml.Decode(string(data), &values); err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(values)
			if err != nil {
				t.Fatal(err)
			}
			if err := protojson.Unmarshal(encoded, actual); err != nil {
				t.Fatal(err)
			}
		}
		if actual.Config.Version != 1 || len(actual.Config.Moved) != 1 || len(actual.Decks) != 1 || len(actual.Notes) != 1 || len(actual.NoteTypes) != 1 {
			t.Fatalf("grouped values lost: %v", actual)
		}
		deck := actual.Decks[0]
		if deck.Path != "decks/Lesson/deck.anki.toml" || deck.Notes[0].Path != "decks/Lesson/notes/greeting.note.anki.toml" || deck.Notes[0].Config.Cards[0].GetDeck() != deck.Path {
			t.Fatalf("path ownership lost: %v", deck)
		}
		if actual.Notes[0].Config.Cards[0].Deck != nil || actual.Notes[0].Config.Fields["Front"] != "" || deck.Deck.Normal.NewLimit == nil || *deck.Deck.Normal.NewLimit != 0 {
			t.Fatal("optional presence changed")
		}
		if actual.NoteTypes[0].Path != "note_types/Basic.notetype.anki.toml" || deck.Notes[0].Config.Fields["Front"] != "<div class=\"example\">你好</div>\nSecond line" {
			t.Fatal("content changed")
		}
		assertConfigSources(t, before, configSources(t, fixture.root))
		if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
			writeConfigFixture(t, filepath.Join(artifacts, "path-configuration."+format), string(data))
		}
	}
}

// TestConfigExportWorkflow verifies supplied path objects are written directly without normalization or validation.
func TestConfigExportWorkflow(t *testing.T) {
	fixture := newConfigFixture(t)
	configs, err := config.New(fixture.root).Load()
	if err != nil {
		t.Fatal(err)
	}
	original := proto.Clone(configs)
	output := filepath.Join(t.TempDir(), "exported")
	other := filepath.Join(output, "README.txt")
	writeConfigFixture(t, other, "unrelated\n")
	store := config.New(output)
	if err := store.Export(configs); err != nil {
		t.Fatal(err)
	}
	actual, err := store.Load()
	if err != nil || !proto.Equal(configs, actual) || !proto.Equal(configs, original) {
		t.Fatalf("object export changed values: %v", err)
	}
	data, err := os.ReadFile(other)
	if err != nil || string(data) != "unrelated\n" {
		t.Fatal("unrelated file changed")
	}
	unchecked := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
	unchecked.Config.Version = 9007199254740993
	unchecked.Notes[0].Config.NoteType = ""
	uncheckedOutput := filepath.Join(t.TempDir(), "unchecked")
	if err := config.New(uncheckedOutput).Export(unchecked); err != nil {
		t.Fatalf("writer validated input: %v", err)
	}
	var root map[string]any
	if _, err := toml.DecodeFile(filepath.Join(uncheckedOutput, "collection.anki.toml"), &root); err != nil {
		t.Fatal(err)
	}
	if root["version"] != "9007199254740993" {
		t.Fatal("integer precision lost")
	}
	blocked := filepath.Join(t.TempDir(), "blocked")
	writeConfigFixture(t, blocked, "file\n")
	if err := config.New(blocked).Export(configs); err == nil {
		t.Fatal("filesystem error ignored")
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := config.New(filepath.Join(artifacts, "exported-configuration")).Export(configs); err != nil {
			t.Fatal(err)
		}
	}
}

// TestConfigFailureWorkflow verifies malformed/unsupported sources fail without partial output or edits.
func TestConfigFailureWorkflow(t *testing.T) {
	cases := map[string]string{
		"syntax":    "version=[",
		"unknown":   "version=1\nunexpected=true\n",
		"type":      "version=true\n",
		"old-key":   "format_version=1\n",
		"version":   "version=2\n",
		"root-path": "version=1\ndecks_path='elsewhere'\n",
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			f := newConfigFixture(t)
			writeConfigFixture(t, f.config, input)
			before := configSources(t, f.root)
			data, err := config.New(f.root).Dump("json")
			if err == nil || len(data) != 0 {
				t.Fatalf("invalid source accepted: %v", err)
			}
			assertConfigSources(t, before, configSources(t, f.root))
		})
	}
	for _, field := range []string{"guid='old'", "note_id=1", "note_type_id=1", "[fields]\nFront=2026-01-01"} {
		t.Run(field, func(t *testing.T) {
			f := newConfigFixture(t)
			writeConfigFixture(t, f.note, field+"\n")
			before := configSources(t, f.root)
			data, err := config.New(f.root).Dump("json")
			if err == nil || len(data) != 0 {
				t.Fatal("removed identity or datetime accepted")
			}
			assertConfigSources(t, before, configSources(t, f.root))
		})
	}
	f := newConfigFixture(t)
	if _, err := config.New(f.root).Dump("yaml"); err == nil {
		t.Fatal("unsupported output accepted")
	}
}

// TestConfigFilesystemWorkflow verifies hierarchy boundaries, empty roots, markers, and symlink rejection.
func TestConfigFilesystemWorkflow(t *testing.T) {
	for _, scenario := range []string{"missing-marker", "dangling-notes", "symlink-note", "non-directory"} {
		t.Run(scenario, func(t *testing.T) {
			f := newConfigFixture(t)
			switch scenario {
			case "missing-marker":
				if err := os.Remove(filepath.Join(f.root, "decks", "Lesson", "deck.anki.toml")); err != nil {
					t.Fatal(err)
				}
			case "dangling-notes":
				if err := os.Symlink("missing", filepath.Join(f.root, "notes-link")); err != nil {
					t.Fatal(err)
				}
				if err := os.RemoveAll(filepath.Join(f.root, "notes")); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(filepath.Join(f.root, "notes-link"), filepath.Join(f.root, "notes")); err != nil {
					t.Fatal(err)
				}
			case "symlink-note":
				if err := os.Symlink(f.note, filepath.Join(f.root, "notes", "link.note.anki.toml")); err != nil {
					t.Fatal(err)
				}
			case "non-directory":
				if err := os.RemoveAll(filepath.Join(f.root, "notes")); err != nil {
					t.Fatal(err)
				}
				writeConfigFixture(t, filepath.Join(f.root, "notes"), "file")
			}
			if _, err := config.New(f.root).Load(); err == nil {
				t.Fatalf("unsafe layout accepted: %s", scenario)
			}
		})
	}
	f := newConfigFixture(t)
	child := filepath.Join(f.root, "decks", "Lesson", "decks", "Child")
	writeConfigFixture(t, filepath.Join(child, "deck.anki.toml"), "[normal]\n")
	writeConfigFixture(t, filepath.Join(child, "notes", "child.note.anki.toml"), "[fields]\nFront=''\n")
	writeConfigFixture(t, filepath.Join(child, "notes", "nested", "ignored.note.anki.toml"), "invalid=[")
	configs, err := config.New(f.root).Load()
	if err != nil || len(configs.Decks) != 2 || configs.Decks[1].Title != "Lesson::Child" || len(configs.Decks[1].Notes) != 1 {
		t.Fatalf("nested ownership lost: %v", err)
	}
	empty := t.TempDir()
	writeConfigFixture(t, filepath.Join(empty, "collection.anki.toml"), "version=1\n")
	if cfg, err := config.New(empty).Load(); err != nil || len(cfg.Decks)+len(cfg.Notes)+len(cfg.NoteTypes) != 0 {
		t.Fatalf("empty layout rejected: %v", err)
	}
}
