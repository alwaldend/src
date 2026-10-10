package collection

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
	"google.golang.org/protobuf/proto"
)

// These behavioral cases precede implementation. The fixture contains shared
// note fields, a filtered card, a normal card, review history, and opaque media.
func TestCollectionWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	textDir := filepath.Join(root, "text")
	if err := Export(ctx, base, textDir); err != nil {
		t.Fatal(err)
	}
	unchanged := filepath.Join(root, "unchanged.colpkg")
	if err := Build(ctx, base, textDir, unchanged); err != nil {
		t.Fatal(err)
	}
	assertRows(t, base, unchanged, root, "select * from notes order by id", true)
	assertRows(t, base, unchanged, root, "select * from cards order by id", true)
	assertRows(t, base, unchanged, root, "select * from revlog order by id", true)
	file := textstore.NotePath(textDir, "grammar", 100)
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	b = []byte(strings.ReplaceAll(string(b), "<b>你好</b>", "<strong>你好，世界</strong>"))
	if err := os.WriteFile(file, b, 0o600); err != nil {
		t.Fatal(err)
	}
	moveFixtureNote(t, textDir, "grammar::A1::Greeting")
	output := filepath.Join(root, "edited.colpkg")
	if err := Build(ctx, base, textDir, output); err != nil {
		t.Fatal(err)
	}
	lesson := model.Deck{Deck: &collectionpb.Deck{}}
	if err := textstore.ReadTOML(filepath.Join(textstore.DeckPath(textDir, "grammar::A1::Greeting"), "deck.anki.toml"), lesson.Deck); err != nil {
		t.Fatal(err)
	}
	if lesson.Id <= 0 {
		t.Fatal("generated deck identity was not recorded in text")
	}
	p, err := openPackage(ctx, output, root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	var fields, deck string
	if err := p.db.QueryRow("select flds from notes where id=100").Scan(&fields); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(fields, "<strong>你好，世界</strong>") {
		t.Fatalf("TOML field edit missing: %s", fields)
	}
	if err := p.db.QueryRow("select name from decks where id=(select odid from cards where id=200)").Scan(&deck); err != nil {
		t.Fatal(err)
	}
	if deck != "grammar" {
		t.Fatalf("wrong filtered home deck %q", deck)
	}
	var did int64
	if err := p.db.QueryRow("select did from cards where id=200").Scan(&did); err != nil {
		t.Fatal(err)
	}
	if did != 2 {
		t.Fatal("filtered deck membership changed")
	}
	assertRows(t, base, output, root, "select id,nid,ord,type,queue,due,ivl,factor,reps,lapses,left,odue,flags,data from cards order by id", true)
	assertRows(t, base, output, root, "select * from cards where id=200", true)
	assertRows(t, base, output, root, "select * from revlog order by id", true)
	z1, err := zip.OpenReader(base)
	if err != nil {
		t.Fatal(err)
	}
	defer z1.Close()
	z2, err := zip.OpenReader(output)
	if err != nil {
		t.Fatal(err)
	}
	defer z2.Close()
	for _, f := range z1.File {
		if f.Name == "collection.anki21b" {
			continue
		}
		a := readEntry(t, f)
		var other *zip.File
		for _, g := range z2.File {
			if g.Name == f.Name {
				other = g
			}
		}
		if other == nil || string(a) != string(readEntry(t, other)) {
			t.Fatalf("entry changed: %s", f.Name)
		}
	}
	artifactRoot := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if artifactRoot == "" {
		artifactRoot = root
	}
	if err := os.WriteFile(filepath.Join(artifactRoot, "verification.txt"), []byte("round trip, multiline TOML edit, lesson decks, scheduling, filtered membership, preserved review history, and media verified\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestRejectInvalidText verifies malformed desired declarations fail before mutation.
func TestRejectInvalidText(t *testing.T) {
	for _, mode := range []string{"duplicate", "malformed", "field-name", "cloze", "same-output", "filtered-deck"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			base := fixture(t, root)
			dir := filepath.Join(root, "text")
			if err := Export(ctx, base, dir); err != nil {
				t.Fatal(err)
			}
			file := textstore.NotePath(dir, "grammar", 100)
			b, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(root, "output.colpkg")
			switch mode {
			case "missing":
				err = os.Remove(file)
			case "duplicate":
				err = os.WriteFile(filepath.Join(filepath.Dir(file), "duplicate.note.anki.toml"), b, 0o600)
			case "malformed":
				err = os.WriteFile(file, []byte("---\nunknown: true\n---\nbroken"), 0o600)
			case "card-id":
				err = os.WriteFile(file, []byte(strings.Replace(string(b), "id = 200", "id = 999", 1)), 0o600)
			case "field-name":
				err = os.WriteFile(file, []byte(strings.Replace(string(b), `hanzi =`, `invalid =`, 1)), 0o600)
			case "cloze":
				err = os.WriteFile(file, []byte(strings.Replace(string(b), "<b>你好</b>", "{{c1::你好}}", 1)), 0o600)
			case "filtered-deck":
				note, readErr := textstore.ReadNote(file)
				if readErr != nil {
					t.Fatal(readErr)
				}
				note.Cards[0].Deck = proto.String("filtered")
				err = textstore.WriteNote(file, note)
			case "same-output":
				out = base
			}
			if err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(base)
			if err != nil {
				t.Fatal(err)
			}
			if out != base {
				if err := os.WriteFile(out, []byte("sentinel"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if err := Build(ctx, base, dir, out); err == nil {
				t.Fatal("invalid input accepted")
			}
			after, err := os.ReadFile(base)
			if err != nil {
				t.Fatal(err)
			}
			if string(before) != string(after) {
				t.Fatal("base changed")
			}
			if out != base {
				got, err := os.ReadFile(out)
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != "sentinel" {
					t.Fatal("existing output changed")
				}
			}
		})
	}
}

// TestDesiredStateReconciliation verifies additions, updates, and removals match the desired state.
func TestDesiredStateReconciliation(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root, true)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	path := textstore.NotePath(dir, "grammar", 100)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	newNote := strings.NewReplacer("note_id = 100", "note_id = 101", "guid = \"guid\"", "guid = \"new-guid\"", "id = 200", "id = 202", "id = 201", "id = 203").Replace(string(b))
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(textstore.NotePath(dir, "grammar", 101), []byte(newNote), 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := Plan(ctx, base, dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.NotesAdded != 1 || p.NotesDeleted != 1 || p.CardsAdded != 2 || p.CardsDeleted != 2 {
		t.Fatalf("wrong plan: %+v", p)
	}
	out := filepath.Join(root, "result.colpkg")
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	a, err := openPackage(ctx, out, root)
	if err != nil {
		t.Fatal(err)
	}
	dbfile := filepath.Join(root, "live.anki2")
	data, err := os.ReadFile(a.database)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dbfile, data, 0o600); err != nil {
		t.Fatal(err)
	}
	db, err := openDB(dbfile)
	if err != nil {
		t.Fatal(err)
	}
	var graves int
	if err := db.QueryRow("select count(*) from graves").Scan(&graves); err != nil {
		t.Fatal(err)
	}
	if graves != 3 {
		t.Fatalf("expected note and card deletion tombstones, got %d", graves)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	p, err = Plan(ctx, dbfile, dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.NotesAdded+p.NotesChanged+p.NotesDeleted+p.CardsAdded+p.CardsChanged+p.CardsDeleted != 0 {
		t.Fatalf("not converged: %+v", p)
	}
	if err := Apply(ctx, dbfile, dir); err != nil {
		t.Fatal(err)
	}
	// A second apply must be a no-op and keep identities/state stable.
	if err := Apply(ctx, dbfile, dir); err != nil {
		t.Fatal(err)
	}
	archiveTarget := filepath.Join(root, "apply.colpkg")
	if err := archive.CopyFile(base, archiveTarget); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(archiveTarget, 0o640); err != nil {
		t.Fatal(err)
	}
	if err := Apply(ctx, archiveTarget, dir); err != nil {
		t.Fatal(err)
	}
	archiveInfo, err := os.Stat(archiveTarget)
	if err != nil || archiveInfo.Mode().Perm() != 0o640 {
		t.Fatalf("archive permissions changed: %v, %v", archiveInfo, err)
	}
	archivePlan, err := Plan(ctx, archiveTarget, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(archivePlan.Changes) != 0 {
		t.Fatalf("archive apply did not converge: %+v", archivePlan)
	}
	// Restoring text must remove tombstones for the restored identities.
	if err := os.Remove(textstore.NotePath(dir, "grammar", 101)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Apply(ctx, dbfile, dir); err != nil {
		t.Fatal(err)
	}
	db, err = openDB(dbfile)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.QueryRow("select count(*) from graves where oid in (100,200,201)").Scan(&graves); err != nil {
		t.Fatal(err)
	}
	if graves != 0 {
		t.Fatalf("restored records retain %d deletion tombstones", graves)
	}
}

// TestDeckSettingsEdit verifies editable deck settings retain unmanaged configuration.
func TestDeckSettingsEdit(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(textstore.DeckPath(dir, "grammar"), "deck.anki.toml")
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	b = []byte(strings.Replace(string(b), "extend_new = 0", "extend_new = 25", 1))
	if err := os.WriteFile(file, b, 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := Plan(ctx, base, dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.DecksChanged != 1 {
		t.Fatalf("deck settings not planned: %+v", p)
	}
	out := filepath.Join(root, "settings.colpkg")
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	a, err := openPackage(ctx, out, root)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	var kind []byte
	if err := a.db.QueryRow("select kind from decks where id=1").Scan(&kind); err != nil {
		t.Fatal(err)
	}
	if string(kind) != string([]byte{10, 4, 8, 1, 16, 25}) {
		t.Fatalf("unexpected deck config %x", kind)
	}
	deck := model.Deck{Deck: &collectionpb.Deck{}}
	if err := textstore.ReadTOML(file, deck.Deck); err != nil {
		t.Fatal(err)
	}
	for index, description := range []string{
		`<div class="example">Description</div>`,
		`It's a description with "quotes" and \ paths`,
		"First line\n<div class=\"example\">Second line</div>\n",
		"Delimiter \"\"\"\nCRLF\r\nTrailing space \n",
	} {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			deck.Normal.Description = description
			if err := textstore.WriteTOML(file, deck.Deck); err != nil {
				t.Fatal(err)
			}
			if err := Build(ctx, base, dir, out); err != nil {
				t.Fatal(err)
			}
			exported := filepath.Join(root, fmt.Sprintf("deck-settings-%d", index))
			if err := Export(ctx, out, exported); err != nil {
				t.Fatal(err)
			}
			actual := model.Deck{Deck: &collectionpb.Deck{}}
			if err := textstore.ReadTOML(filepath.Join(textstore.DeckPath(exported, "grammar"), textstore.DeckMarker), actual.Deck); err != nil {
				t.Fatal(err)
			}
			if actual.Normal.Description != description || actual.Normal.ExtendNew != 25 {
				t.Fatalf("deck settings changed: %+v", actual.Normal)
			}
			plan, err := Plan(ctx, out, exported)
			if err != nil || len(plan.Changes) != 0 {
				t.Fatalf("deck settings did not converge: %+v, %v", plan, err)
			}
		})
	}
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.WriteFile(filepath.Join(artifact, "deck-string-quoting.txt"), []byte("Deck export and rebuild preserve quoted HTML, apostrophes, delimiters, backslashes, controls, and newline boundaries; plans converge.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestMultilineFields verifies note fields with newlines round-trip without content changes.
func TestMultilineFields(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	path := textstore.NotePath(dir, "grammar", 100)
	note, err := textstore.ReadNote(path)
	if err != nil {
		t.Fatal(err)
	}
	note.Fields["hanzi"] = "Line one \nLine two with \"quotes\" and \\ backslash"
	if err := textstore.WriteNote(path, note); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "multiline.colpkg")
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	a, err := openPackage(ctx, out, root)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	var fields string
	if err := a.db.QueryRow("select flds from notes where id=100").Scan(&fields); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(fields, "Line one \nLine two with \"quotes\" and \\ backslash") {
		t.Fatalf("multiline content changed: %q", fields)
	}
	for index, value := range []string{
		`<div class="example">你好</div>`,
		`C:\notes\cards`,
		`It's an example with "quotes"`,
		"\nLeading newline",
		"Trailing newline\n",
		"\n\nBoth boundaries\n\n",
		"Multiline\nIt's quoted 'once' and ''twice''",
		"Multiline\nEnds with an apostrophe'",
		"Multiline\nEnds with two apostrophes''",
		"Multiline\nConflicting ''' delimiter and \"\"\" quotes",
		"CRLF\r\nLine and bare CR\r",
		"Tabs\tand\ntrailing space \n",
		"Control\b\f\x7f\ncharacters",
	} {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			note.Fields["hanzi"] = value
			if err := textstore.WriteNote(path, note); err != nil {
				t.Fatal(err)
			}
			if err := Build(ctx, base, dir, out); err != nil {
				t.Fatal(err)
			}
			result, err := openPackage(ctx, out, root)
			if err != nil {
				t.Fatal(err)
			}
			defer result.Close()
			var actual string
			if err := result.db.QueryRow("select flds from notes where id=100").Scan(&actual); err != nil {
				t.Fatal(err)
			}
			if strings.Split(actual, "\x1f")[0] != value {
				t.Fatalf("field changed: got %q, want %q", actual, value)
			}
		})
	}
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.WriteFile(filepath.Join(artifact, "field-strings.txt"), []byte("Field string boundaries, HTML, backslashes, apostrophes, delimiters, controls, and exact newlines verified through archive builds.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestNewNoteWorkflow verifies newly declared notes and cards reconcile successfully.
func TestNewNoteWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	n, err := copyFixtureNote(t, dir, "grammar", "greetings.note.anki.toml")
	if err != nil {
		t.Fatal(err)
	}
	if n.NoteId <= 0 || n.Guid == "" || len(n.Cards) != 1 {
		t.Fatalf("invalid generated identity: %+v", n)
	}
	path := filepath.Join(textstore.DeckPath(dir, "grammar"), "notes", "greetings.note.anki.toml")
	saved, err := textstore.ReadNote(path)
	if err != nil {
		t.Fatal(err)
	}
	if value, present := saved.Fields["hanzi"]; !present || value != "" {
		t.Fatal("missing empty note type field")
	}
	p, err := Plan(ctx, base, dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.NotesAdded != 1 || p.CardsAdded != 1 {
		t.Fatalf("wrong create plan: %+v", p)
	}
	if err := os.Rename(path, filepath.Join(filepath.Dir(path), "renamed.note.anki.toml")); err != nil {
		t.Fatal(err)
	}
	after, err := Plan(ctx, base, dir)
	if err != nil {
		t.Fatal(err)
	}
	if after.NotesAdded != p.NotesAdded || after.CardsAdded != p.CardsAdded {
		t.Fatal("rename changed identity")
	}
	out := filepath.Join(root, "created.colpkg")
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	a, err := openPackage(ctx, out, root)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	var count int
	if err := a.db.QueryRow("select count(*) from notes where id=?", n.NoteId).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("generated note missing")
	}
}

// TestNewCardPositions verifies new-card position allocation retains existing scheduling.
func TestNewCardPositions(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root, true)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"first.note.anki.toml", "second.note.anki.toml"} {
		if _, err := copyFixtureNote(t, dir, "grammar", name); err != nil {
			t.Fatal(err)
		}
	}
	out := filepath.Join(root, "positions.colpkg")
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	archive, err := openPackage(ctx, out, root)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	var siblings, distinctPositions int
	if err := archive.db.QueryRow("select count(*), count(distinct due) from cards where nid!=100").Scan(&siblings, &distinctPositions); err != nil {
		t.Fatal(err)
	}
	if siblings != 4 || distinctPositions != 2 {
		t.Fatalf("new cards: %d cards at %d positions, want four siblings at two note positions", siblings, distinctPositions)
	}
	var mismatchedNotes int
	if err := archive.db.QueryRow("select count(*) from (select nid from cards where nid!=100 group by nid having min(due)!=max(due))").Scan(&mismatchedNotes); err != nil {
		t.Fatal(err)
	}
	var minimum, maximum, nextPosition int64
	if err := archive.db.QueryRow("select min(due),max(due) from cards where nid!=100").Scan(&minimum, &maximum); err != nil {
		t.Fatal(err)
	}
	if err := archive.db.QueryRow("select cast(val as integer) from config where key='nextPos'").Scan(&nextPosition); err != nil {
		t.Fatal(err)
	}
	if mismatchedNotes != 0 || minimum != 100 || maximum != 101 || nextPosition != 102 {
		t.Fatalf("invalid allocation: mismatched=%d, min=%d, max=%d, cursor=%d", mismatchedNotes, minimum, maximum, nextPosition)
	}
	var storageType string
	if err := archive.db.QueryRow("select typeof(val) from config where key='nextPos'").Scan(&storageType); err != nil {
		t.Fatal(err)
	}
	if storageType != "blob" {
		t.Fatalf("allocation cursor storage: %s, want Anki's blob encoding", storageType)
	}
	plan, err := Plan(ctx, out, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("allocation did not converge: %v, %+v", err, plan)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "new-card-positions.txt"), []byte("four siblings; two note positions: 100,101; nextPos=102; converged\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestEmptyDeckNoteDirectories verifies missing note directories represent empty decks.
func TestEmptyDeckNoteDirectories(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(textstore.DeckPath(dir, "filtered"), "notes")); err != nil {
		t.Fatal(err)
	}
	emptyDeck := textstore.DeckPath(dir, "empty")
	if err := os.MkdirAll(emptyDeck, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := textstore.WriteTOML(filepath.Join(emptyDeck, "deck.anki.toml"), &collectionpb.Deck{Id: 3, Title: "empty", NotesPath: "notes", Normal: &collectionpb.DeckConfig{ConfigId: 1}}); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "empty-deck.colpkg")
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(ctx, out, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("empty deck did not reconcile: %v, %+v", err, plan)
	}
	if _, err := copyFixtureNote(t, dir, "empty", "first.note.anki.toml"); err != nil {
		t.Fatal(err)
	}
	if err := Build(ctx, base, dir, out); err != nil {
		t.Fatal(err)
	}
	plan, err = Plan(ctx, out, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("first note did not reconcile: %v, %+v", err, plan)
	}
	if artifactDirectory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactDirectory != "" {
		if err := os.WriteFile(filepath.Join(artifactDirectory, "empty-deck-verification.txt"), []byte("empty notes directories accepted; first note created; collection converged\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestReviewContracts verifies reconciliation guarantees recorded during review.
func TestReviewContracts(t *testing.T) {
	ctx := context.Background()
	t.Run("media-caches", func(t *testing.T) {
		root := t.TempDir()
		base := fixture(t, root)
		dir := filepath.Join(root, "text")
		if err := Export(ctx, base, dir); err != nil {
			t.Fatal(err)
		}
		file := textstore.NotePath(dir, "grammar", 100)
		n, err := textstore.ReadNote(file)
		if err != nil {
			t.Fatal(err)
		}
		n.Fields["hanzi"] = `<audio src="audio.mp3"></audio><video src="video.mp4"></video><source src="stream.mp4"><img src="image.png"><object data="object.bin"></object>`
		if err := textstore.WriteNote(file, n); err != nil {
			t.Fatal(err)
		}
		output := filepath.Join(root, "media.colpkg")
		if err := Build(ctx, base, dir, output); err != nil {
			t.Fatal(err)
		}
		p, err := openPackage(ctx, output, root)
		if err != nil {
			t.Fatal(err)
		}
		defer p.Close()
		var sortField, fields string
		var checksum uint32
		if err := p.db.QueryRow("select sfld,csum,flds from notes where id=100").Scan(&sortField, &checksum, &fields); err != nil {
			t.Fatal(err)
		}
		expected := " audio.mp3  video.mp4  stream.mp4  image.png  object.bin "
		digest := sha1.Sum([]byte(expected))
		if sortField != expected || checksum != binary.BigEndian.Uint32(digest[:4]) || strings.Split(fields, "\x1f")[0] != n.Fields["hanzi"] {
			t.Fatalf("incorrect media caches: sort=%q, checksum=%d", sortField, checksum)
		}
	})
	t.Run("declared-deck-preset", func(t *testing.T) {
		root := t.TempDir()
		base := fixture(t, root)
		dir := filepath.Join(root, "text")
		if err := Export(ctx, base, dir); err != nil {
			t.Fatal(err)
		}
		if err := textstore.WriteDeck(dir, model.Deck{Deck: &collectionpb.Deck{Id: 3, Title: "grammar::A1"}, Kind: "CgIIAg=="}); err != nil {
			t.Fatal(err)
		}
		if err := textstore.WriteDeck(dir, model.Deck{Deck: &collectionpb.Deck{Title: "grammar::A1::Greeting"}, Kind: "CgIIAg=="}); err != nil {
			t.Fatal(err)
		}
		moveFixtureNote(t, dir, "grammar::A1::Greeting")
		output := filepath.Join(root, "parent.colpkg")
		if err := Build(ctx, base, dir, output); err != nil {
			t.Fatal(err)
		}
		exported := filepath.Join(root, "reexport")
		if err := Export(ctx, output, exported); err != nil {
			t.Fatal(err)
		}
		lesson := model.Deck{Deck: &collectionpb.Deck{}}
		if err := textstore.ReadTOML(filepath.Join(textstore.DeckPath(exported, "grammar::A1::Greeting"), "deck.anki.toml"), lesson.Deck); err != nil {
			t.Fatal(err)
		}
		if lesson.Normal == nil || lesson.Normal.ConfigId != 2 {
			t.Fatalf("declared deck did not retain its preset: %+v", lesson.Normal)
		}
	})
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" && !t.Failed() {
		if err := os.WriteFile(filepath.Join(artifacts, "review-contracts.txt"), []byte("media caches match Anki; declared deck preset retained\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestDeckTitleReconciliation verifies deck renames, title swaps, and title reuse.
func TestDeckTitleReconciliation(t *testing.T) {
	for _, scenario := range []string{"swap", "reuse-deleted", "swap-keep-title"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			fixture(t, root)
			source := filepath.Join(root, "fixture.sqlite")
			db, err := openDB(source)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec("INSERT INTO decks VALUES(3,'other',0,0,x'',x'0a020801')"); err != nil {
				t.Fatal(err)
			}
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			dir := filepath.Join(root, "text")
			if err := Export(ctx, source, dir); err != nil {
				t.Fatal(err)
			}
			if err := textstore.WriteDeck(dir, model.Deck{Deck: &collectionpb.Deck{Id: 1, Title: "other"}, Kind: "CgIIAQ=="}); err != nil {
				t.Fatal(err)
			}
			// Keep the declarations in arbitrary folders; only titles change.
			if err := os.Remove(filepath.Join(textstore.DeckPath(dir, "grammar"), textstore.DeckMarker)); err != nil {
				t.Fatal(err)
			}
			deck := model.Deck{Deck: &collectionpb.Deck{}}
			marker := filepath.Join(textstore.DeckPath(dir, "other"), textstore.DeckMarker)
			if err := textstore.ReadTOML(marker, deck.Deck); err != nil {
				t.Fatal(err)
			}
			deck.NotesPath = "../grammar/notes"
			if err := textstore.WriteTOML(marker, deck.Deck); err != nil {
				t.Fatal(err)
			}
			if scenario != "reuse-deleted" {
				if err := textstore.WriteTOML(filepath.Join(textstore.DeckPath(dir, "grammar"), textstore.DeckMarker), &collectionpb.Deck{Id: 3, Title: "grammar", NotesPath: "empty", Normal: &collectionpb.DeckConfig{ConfigId: 1}}); err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "swap-keep-title" {
				path := textstore.NotePath(dir, "grammar", 100)
				note, err := textstore.ReadNote(path)
				if err != nil {
					t.Fatal(err)
				}
				for i := range note.Cards {
					note.Cards[i].Deck = proto.String("grammar")
				}
				if err := textstore.WriteNote(path, note); err != nil {
					t.Fatal(err)
				}
			}
			plan, err := Plan(ctx, source, dir)
			if err != nil {
				t.Fatal(err)
			}
			wantCardChanges := 0
			wantHome := int64(1)
			if scenario == "swap-keep-title" {
				wantCardChanges = 1
				wantHome = 3
			}
			if plan.CardsChanged != int64(wantCardChanges) {
				t.Fatalf("planned %d card moves, want %d", plan.CardsChanged, wantCardChanges)
			}
			if err := Apply(ctx, source, dir); err != nil {
				t.Fatal(err)
			}
			db, err = openDB(source)
			if err != nil {
				t.Fatal(err)
			}
			var normalHome, filteredHome, filteredMembership, originalDue int64
			if err := db.QueryRow("SELECT did FROM cards WHERE id=201").Scan(&normalHome); err != nil {
				t.Fatal(err)
			}
			if err := db.QueryRow("SELECT odid,did,odue FROM cards WHERE id=200").Scan(&filteredHome, &filteredMembership, &originalDue); err != nil {
				t.Fatal(err)
			}
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			if normalHome != wantHome || filteredHome != 1 || filteredMembership != 2 || originalDue != 100 {
				t.Fatalf("wrong homes or filtered scheduling: normal=%d, filtered=%d, membership=%d, due=%d", normalHome, filteredHome, filteredMembership, originalDue)
			}
			plan, err = Plan(ctx, source, dir)
			if err != nil || len(plan.Changes) != 0 {
				t.Fatalf("deck titles did not converge: %+v, %v", plan, err)
			}
			if err := Apply(ctx, source, dir); err != nil {
				t.Fatal(err)
			}
			if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
				if err := os.WriteFile(filepath.Join(artifacts, "deck-titles-"+scenario+".txt"), []byte(fmt.Sprintf("Deck titles reconciled with Anki's unique name index; %d planned card moves; normal and filtered home IDs=%d; filtered membership and original due preserved; repeated apply converged.\n", wantCardChanges, wantHome)), 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestOccupiedFilteredDeckConversion verifies occupied filtered decks cannot become normal decks.
func TestOccupiedFilteredDeckConversion(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	if err := textstore.WriteDeck(dir, model.Deck{Deck: &collectionpb.Deck{Id: 2, Title: "filtered"}, Kind: "CgIIAQ=="}); err != nil {
		t.Fatal(err)
	}
	before, err := archive.FileHash(base)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(ctx, base, dir); err == nil {
		t.Fatal("plan accepted conversion of an occupied filtered deck")
	}
	if err := Apply(ctx, base, dir); err == nil {
		t.Fatal("apply accepted conversion of an occupied filtered deck")
	}
	after, err := archive.FileHash(base)
	if err != nil || after != before {
		t.Fatal("rejected conversion changed the input archive")
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "filtered-conversion.txt"), []byte("Plan and apply reject occupied filtered-deck conversion; input archive preserved exactly.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestFilteredCardProtection checks both public reconciliation entry points
// against attempts to claim a filtered card or orphan it by deleting its note.
func TestFilteredCardProtection(t *testing.T) {
	for _, scenario := range []string{"identity", "ordinal", "delete-note"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			base := fixture(t, root)
			dir := filepath.Join(root, "text")
			if err := Export(ctx, base, dir); err != nil {
				t.Fatal(err)
			}
			path := textstore.NotePath(dir, "grammar", 100)
			note, err := textstore.ReadNote(path)
			if err != nil {
				t.Fatal(err)
			}
			if len(note.Cards) != 1 || note.Cards[0].Id != 201 {
				t.Fatalf("export includes unmanaged filtered cards: %+v", note.Cards)
			}
			switch scenario {
			case "identity":
				note.Cards[0].Id = 200
			case "ordinal":
				note.Cards = append(note.Cards, &collectionpb.Card{Id: 202, Ordinal: 0, Deck: proto.String("grammar")})
			case "delete-note":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
			if scenario != "delete-note" {
				if err := textstore.WriteNote(path, note); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.ReadFile(base)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Plan(ctx, base, dir); err == nil {
				t.Fatal("plan accepted a conflict with an unmanaged filtered card")
			}
			if err := Apply(ctx, base, dir); err == nil {
				t.Fatal("apply accepted a conflict with an unmanaged filtered card")
			}
			after, err := os.ReadFile(base)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("rejected reconciliation changed the archive: %v", err)
			}
			if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
				if err := os.WriteFile(filepath.Join(artifacts, "filtered-protection-"+scenario+".txt"), []byte("Filtered card excluded from export; conflicting plan and apply rejected; original archive preserved exactly.\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestBuildCompatibleArchiveWithoutBinding verifies desired text can rebuild a compatible unbound base.
func TestBuildCompatibleArchiveWithoutBinding(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, textstore.CollectionMarker)
	contents, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(contents), "base_sha256") {
		t.Error("exported text is still bound to an archive hash")
	}
	source := filepath.Join(root, "fixture.sqlite")
	db, err := openDB(source)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("UPDATE cards SET ivl=30,reps=12 WHERE id=200"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO revlog VALUES(501,200,4)"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	updated := filepath.Join(root, "updated.colpkg")
	archiveWithDatabase(t, base, source, updated)
	// The same desired text applies to another compatible archive, including
	// old exports whose now-unused hash metadata must not constrain builds.
	for _, legacy := range []bool{false, true} {
		text := contents
		if legacy {
			text = []byte("base_sha256 = 'unused-old-hash'\n" + string(contents))
		}
		if err := os.WriteFile(config, text, 0o600); err != nil {
			t.Fatal(err)
		}
		output := filepath.Join(root, fmt.Sprintf("rebuilt-%t.colpkg", legacy))
		if err := Build(ctx, updated, dir, output); err != nil {
			t.Fatal(err)
		}
		assertRows(t, updated, output, root, "SELECT * FROM cards ORDER BY id", true)
		assertRows(t, updated, output, root, "SELECT * FROM revlog ORDER BY id", true)
		assertRows(t, updated, output, root, "SELECT * FROM config ORDER BY key", true)
		plan, err := Plan(ctx, output, dir)
		if err != nil || len(plan.Changes) != 0 {
			t.Fatalf("compatible archive did not converge: %+v, %v", plan, err)
		}
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "unbound-text.txt"), []byte("Config omits archive hash; changed compatible archive and legacy metadata accepted; current scheduling, review history, and configuration preserved; empty plans.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
