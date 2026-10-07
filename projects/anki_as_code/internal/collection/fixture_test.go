package collection

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/repository"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

// fixture creates a schema-18 collection archive with representative notes and cards.
func fixture(t *testing.T, root string, regularCards ...bool) string {
	t.Helper()
	dbpath := filepath.Join(root, "fixture.sqlite")
	db, err := openDB(dbpath)
	if err != nil {
		t.Fatal(err)
	}
	ddl := []string{
		"CREATE TABLE col(id integer primary key, ver integer, mod integer, scm integer)",
		"INSERT INTO col VALUES(1,18,0,0)",
		"CREATE TABLE notetypes(id integer primary key,name text,mtime_secs integer,usn integer,config blob)",
		"INSERT INTO notetypes VALUES(10,'Chinese',0,0,x'a00607')",
		"CREATE TABLE fields(ntid integer,ord integer,name text,config blob)",
		"INSERT INTO fields VALUES(10,0,'hanzi',x''),(10,1,'article-title',x''),(10,2,'english',x'')",
		"CREATE TABLE templates(ntid integer,ord integer,name text,mtime_secs integer,usn integer,config blob)",
		"INSERT INTO templates VALUES(10,0,'Forward',0,0,x'a00607'),(10,1,'Reverse',0,0,x'')",
		"CREATE TABLE decks(id integer primary key,name text COLLATE unicase,mtime_secs integer,usn integer,common blob,kind blob)",
		"CREATE UNIQUE INDEX idx_decks_name ON decks(name)",
		"INSERT INTO decks VALUES(1,'grammar',0,0,x'',x'0a020801'),(2,'filtered',0,0,x'',x'1200')",
		"CREATE TABLE notes(id integer primary key,guid text,mid integer,mod integer,usn integer,tags text,flds text,sfld text,csum integer,flags integer,data text)",
		"CREATE TABLE cards(id integer primary key,nid integer,did integer,ord integer,mod integer,usn integer,type integer,queue integer,due integer,ivl integer,factor integer,reps integer,lapses integer,left integer,odue integer,odid integer,flags integer,data text)",
		"INSERT INTO cards VALUES(200,100,2,0,1,0,2,2,99,10,2500,4,1,0,100,1,0,'opaque'),(201,100,1,1,1,0,0,0,4,0,2500,0,0,0,0,0,0,'')",
		"CREATE TABLE revlog(id integer,cid integer,ease integer)",
		"CREATE TABLE graves(oid integer,type integer,usn integer,PRIMARY KEY(oid,type))",
		"INSERT INTO revlog VALUES(500,200,3)",
		"CREATE TABLE tags(tag text primary key COLLATE unicase,usn integer,collapsed integer,config blob)",
		"INSERT INTO tags VALUES('A1',0,0,NULL)",
		"CREATE TABLE config(key text primary key,usn integer,mtime_secs integer,val blob)",
		"INSERT INTO config VALUES('setting',0,0,'true')",
		"INSERT INTO config VALUES('nextPos',0,0,'100')",
	}
	for _, q := range ddl {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec("INSERT INTO notes VALUES(100,'guid',10,0,0,' A1 ',?,'你好',0,0,'')", "<b>你好</b>\x1fGreeting\x1fhello\n"); err != nil {
		t.Fatal(err)
	}
	if len(regularCards) > 0 && regularCards[0] {
		if _, err := db.Exec("UPDATE cards SET did=1,odid=0,odue=0 WHERE id=200"); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dbpath)
	if err != nil {
		t.Fatal(err)
	}
	z, err := zstd.NewWriter(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	base := filepath.Join(root, "base.colpkg")
	f, err := os.Create(base)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for _, entry := range []struct {
		// name identifies the ZIP entry written into the fixture archive.
		name string
		// data contains the exact ZIP entry payload.
		data []byte
	}{{"meta", []byte{8, 3}}, {"collection.anki21b", z.EncodeAll(b, nil)}, {"collection.anki2", []byte("dummy")}, {"media", []byte("opaque media map")}, {"0", []byte("opaque image")}} {
		s, err := w.Create(entry.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.Write(entry.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return base
}

// readEntry reads a ZIP entry and fails the test on resource or I/O errors.
func readEntry(t *testing.T, f *zip.File) []byte {
	t.Helper()
	r, err := f.Open()
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	return b
}

// assertRows compares selected database rows between two fixture collections.
func assertRows(t *testing.T, a, b, root, query string, equal bool) {
	t.Helper()
	read := func(path string) string {
		p, err := openPackage(context.Background(), path, root)
		if err != nil {
			t.Fatal(err)
		}
		defer p.Close()
		r, err := p.db.Query(query)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		cols, err := r.Columns()
		if err != nil {
			t.Fatal(err)
		}
		result := ""
		for r.Next() {
			values := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range values {
				ptrs[i] = &values[i]
			}
			if err := r.Scan(ptrs...); err != nil {
				t.Fatal(err)
			}
			result += fmt.Sprint(values)
		}
		if err := r.Err(); err != nil {
			t.Fatal(err)
		}
		return result
	}
	if (read(a) == read(b)) != equal {
		t.Fatalf("rows differ: %s", query)
	}
}

// copyFixtureNote models copying an exported note and then assigning identities.
func copyFixtureNote(t *testing.T, dir, deck, filename string) (model.Note, error) {
	t.Helper()
	note, err := textstore.ReadNote(textstore.NotePath(dir, "grammar", 100))
	if err != nil {
		return note, fmt.Errorf("read fixture note: %w", err)
	}
	for field := range note.Fields {
		note.Fields[field] = ""
	}
	path := filepath.Join(textstore.DeckPath(dir, deck), "notes", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return note, fmt.Errorf("create fixture note directory: %w", err)
	}
	if err := textstore.WriteNote(path, note); err != nil {
		return note, fmt.Errorf("copy fixture note: %w", err)
	}
	if _, err := textstore.GenerateIDs(context.Background(), []string{path}); err != nil {
		return note, fmt.Errorf("assign fixture identities: %w", err)
	}
	note, err = textstore.ReadNote(path)
	if err != nil {
		return note, fmt.Errorf("read generated fixture: %w", err)
	}
	return note, nil
}

// moveFixtureNote models explicit deck declarations and a note file move.
func moveFixtureNote(t *testing.T, dir, destination string) {
	t.Helper()
	parts := strings.Split(destination, "::")
	for index := range parts {
		name := strings.Join(parts[:index+1], "::")
		if _, err := os.Stat(filepath.Join(textstore.DeckPath(dir, name), "deck.anki.toml")); os.IsNotExist(err) {
			if err := textstore.WriteDeck(dir, model.Deck{Deck: &collectionpb.Deck{Title: name}, Kind: "CgIIAQ=="}); err != nil {
				t.Fatal(err)
			}
		} else if err != nil {
			t.Fatal(err)
		}
	}
	original := textstore.NotePath(dir, "grammar", 100)
	note, err := textstore.ReadNote(original)
	if err != nil {
		t.Fatal(err)
	}
	note.HomeDeck = destination
	for index := range note.Cards {
		note.Cards[index].Deck = proto.String(destination)
	}
	if err := textstore.WriteNote(textstore.NotePath(dir, destination, note.NoteId), note); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(original); err != nil {
		t.Fatal(err)
	}
}

// archiveWithDatabase rebuilds a fixture archive with the supplied database.
func archiveWithDatabase(t *testing.T, base, database, output string) {
	t.Helper()
	original, err := zip.OpenReader(base)
	if err != nil {
		t.Fatal(err)
	}
	defer original.Close()
	file, err := os.Create(output)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	for _, entry := range original.File {
		if entry.Name != "collection.anki21b" {
			if err := writer.Copy(entry); err != nil {
				t.Fatal(err)
			}
			continue
		}
		target, err := writer.CreateHeader(&zip.FileHeader{Name: entry.Name, Method: zip.Store})
		if err != nil {
			t.Fatal(err)
		}
		encoder, err := zstd.NewWriter(target, zstd.WithEncoderConcurrency(1))
		if err != nil {
			t.Fatal(err)
		}
		source, err := os.Open(database)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.Copy(encoder, source); err != nil {
			t.Fatal(err)
		}
		if err := source.Close(); err != nil {
			t.Fatal(err)
		}
		if err := encoder.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

// fixtureArchive owns archive and SQLite resources used by workflow fixtures.
type fixtureArchive struct {
	// source owns the private archive extraction and its cleanup.
	source *archive.Source
	// db opens the extracted fixture database for behavioral assertions.
	db *sql.DB
	// database identifies the extracted fixture database file.
	database string
}

// Close releases the fixture database and archive resources.
func (fixture *fixtureArchive) Close() error {
	var closeErr error
	if err := fixture.db.Close(); err != nil {
		closeErr = fmt.Errorf("close fixture database: %w", err)
	}
	return errors.Join(closeErr, fixture.source.Close())
}

// openDB opens SQLite after explicitly configuring the Anki name collation.
func openDB(path string) (*sql.DB, error) {
	store, err := repository.Open(path)
	if err != nil {
		return nil, fmt.Errorf("initialize fixture database: %w", err)
	}
	if err := store.Close(); err != nil {
		return nil, fmt.Errorf("close fixture initializer: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open fixture SQL handle: %w", err)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

// openPackage opens a fixture archive and its private SQLite database.
func openPackage(ctx context.Context, path, scratch string) (*fixtureArchive, error) {
	source, err := archive.OpenPackage(ctx, path, scratch)
	if err != nil {
		return nil, fmt.Errorf("open fixture archive: %w", err)
	}
	db, err := openDB(source.DatabasePath())
	if err != nil {
		return nil, errors.Join(fmt.Errorf("open fixture database: %w", err), source.Close())
	}
	return &fixtureArchive{source: source, db: db, database: source.DatabasePath()}, nil
}
