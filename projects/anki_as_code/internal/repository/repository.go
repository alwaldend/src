package repository

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/planning"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/validation"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
	"modernc.org/sqlite"
)

// reconcileAppearance updates editable CSS and templates while preserving other configuration.
func (repository *Repository) reconcileAppearance(run executeStatement, appearances []*collectionpb.NoteTypeAppearance, s snapshot, now int64) error {
	for _, appearance := range appearances {
		current := s.Types[appearance.Id]
		if current.Appearance.Css != appearance.Css {
			config := &notetypes.Notetype_Config{}
			if err := proto.Unmarshal(s.noteTypes[appearance.Id].config, config); err != nil {
				return fmt.Errorf("read note type %d configuration: %w", appearance.Id, err)
			}
			config.Css = appearance.Css
			data, err := ankiformat.MarshalConfig(config)
			if err != nil {
				return fmt.Errorf("encode note type %d CSS: %w", appearance.Id, err)
			}
			if err := run("UPDATE notetypes SET config=?,mtime_secs=?,usn=-1 WHERE id=?", data, now, appearance.Id); err != nil {
				return fmt.Errorf("update note type CSS: %w", err)
			}
		}
		for i, template := range appearance.Templates {
			old := current.Appearance.Templates[i]
			if proto.Equal(old, template) {
				continue
			}
			config := &notetypes.Notetype_Template_Config{}
			if err := proto.Unmarshal(s.noteTypes[appearance.Id].templates[int(template.Ordinal)], config); err != nil {
				return fmt.Errorf("read note type %d template %d: %w", appearance.Id, template.Ordinal, err)
			}
			config.QFormat = template.Front
			config.AFormat = template.Back
			config.QFormatBrowser = template.BrowserFront
			config.AFormatBrowser = template.BrowserBack
			config.BrowserFontName = template.BrowserFont
			config.BrowserFontSize = template.BrowserFontSize
			data, err := ankiformat.MarshalConfig(config)
			if err != nil {
				return fmt.Errorf("encode note type %d template %d: %w", appearance.Id, template.Ordinal, err)
			}
			if err := run("UPDATE templates SET config=?,mtime_secs=?,usn=-1 WHERE ntid=? AND ord=?", data, now, appearance.Id, template.Ordinal); err != nil {
				return fmt.Errorf("update card template: %w", err)
			}
			// Anki uploads complete note types selected by the parent's sync marker.
			if err := run("UPDATE notetypes SET mtime_secs=?,usn=-1 WHERE id=?", now, appearance.Id); err != nil {
				return fmt.Errorf("mark template's note type changed: %w", err)
			}
		}
	}
	return nil
}

// Repository owns collection queries and transactional reconciliation.
type Repository struct {
	// db owns the private SQLite connection used for snapshots and reconciliation.
	db *sql.DB
}

// configureSQLite registers the Anki Unicode collation once and retains any setup error.
var configureSQLite = sync.OnceValue(func() error {
	if err := sqlite.RegisterCollationUtf8("unicase", func(a, b string) int {
		return strings.Compare(ankiformat.Folded(a), ankiformat.Folded(b))
	}); err != nil {
		return fmt.Errorf("register Anki Unicode collation: %w", err)
	}
	return nil
})

// openDB opens SQLite after explicitly configuring the Anki name collation.
func openDB(path string) (*sql.DB, error) {
	if err := configureSQLite(); err != nil {
		return nil, fmt.Errorf("configure SQLite: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open SQLite %q: %w", path, err)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

// validateSchema rejects databases whose schema version is unsupported.
func (repository *Repository) validateSchema(ctx context.Context) error {
	var version int
	if err := repository.db.QueryRowContext(ctx, "SELECT ver FROM col").Scan(&version); err != nil {
		return fmt.Errorf("read collection schema: %w", err)
	}
	if version != 18 {
		return fmt.Errorf("unsupported collection schema %d; expected 18", version)
	}
	return nil
}

// query runs a contextual SQLite query and closes its rows after scanning.
func (repository *Repository) query(ctx context.Context, q string, scan func(*sql.Rows) error) (err error) {
	r, err := repository.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("query collection: %w", err)
	}
	defer func() {
		if e := r.Close(); e != nil {
			err = errors.Join(err, fmt.Errorf("close query: %w", e))
		}
	}()
	for r.Next() {
		if err := scan(r); err != nil {
			return fmt.Errorf("scan collection: %w", err)
		}
	}
	if err := r.Err(); err != nil {
		return fmt.Errorf("iterate collection: %w", err)
	}
	return nil
}

// Snapshot loads the current shared collection state from SQLite.
func (repository *Repository) Snapshot(ctx context.Context) (model.CollectionState, error) {
	s, err := repository.loadSnapshot(ctx)
	if err != nil {
		return model.CollectionState{}, fmt.Errorf("load collection snapshot: %w", err)
	}
	return s.CollectionState, nil
}

// loadSnapshot loads shared collection state and private preservation metadata.
func (repository *Repository) loadSnapshot(ctx context.Context) (snapshot, error) {
	if err := repository.validateSchema(ctx); err != nil {
		return snapshot{}, fmt.Errorf("validate collection schema: %w", err)
	}
	s := snapshot{
		CollectionState: model.CollectionState{
			Notes:       map[int64]model.Note{},
			Types:       map[int64]*model.NoteType{},
			CardDeckIDs: map[int64]int64{},
		},
		noteTypes:    map[int64]noteTypeStorage{},
		newPositions: map[int64]int64{},
	}
	if err := repository.loadNoteTypes(ctx, &s); err != nil {
		return s, fmt.Errorf("loadNoteTypes: %w", err)
	}
	if err := repository.loadFields(ctx, &s); err != nil {
		return s, fmt.Errorf("loadFields: %w", err)
	}
	if err := repository.loadTemplates(ctx, &s); err != nil {
		return s, fmt.Errorf("loadTemplates: %w", err)
	}
	for _, nt := range s.Types {
		s.Appearance = append(s.Appearance, nt.Appearance)
	}
	sort.Slice(s.Appearance, func(i, j int) bool { return s.Appearance[i].Id < s.Appearance[j].Id })

	if err := repository.loadDecks(ctx, &s); err != nil {
		return s, fmt.Errorf("loadDecks: %w", err)
	}

	if err := repository.loadNotes(ctx, &s); err != nil {
		return s, fmt.Errorf("loadNotes: %w", err)
	}
	if err := repository.loadCards(ctx, &s); err != nil {
		return s, fmt.Errorf("loadCards: %w", err)
	}
	if err := repository.loadNextPosition(ctx, &s); err != nil {
		return s, fmt.Errorf("load next position: %w", err)
	}
	return s, nil
}

// loadNoteTypes loads upstream note-type identities and original configuration bytes.
func (repository *Repository) loadNoteTypes(ctx context.Context, s *snapshot) error {
	if err := repository.query(ctx, "SELECT id,name,config FROM notetypes ORDER BY id", func(r *sql.Rows) error {
		var id int64
		var name string
		var config []byte
		if err := r.Scan(&id, &name, &config); err != nil {
			return fmt.Errorf("read note type: %w", err)
		}
		definition := &notetypes.Notetype_Config{}
		if err := proto.Unmarshal(config, definition); err != nil {
			return fmt.Errorf("parse note type %d: %w", id, err)
		}
		s.Types[id] = &model.NoteType{
			Notetype:   &notetypes.Notetype{Id: id, Name: name, Config: definition},
			Ordinals:   map[int]bool{},
			Appearance: &collectionpb.NoteTypeAppearance{Id: id, Name: name, Css: definition.Css},
		}
		s.noteTypes[id] = noteTypeStorage{config: config, templates: map[int][]byte{}}
		return nil
	}); err != nil {
		return fmt.Errorf("load note types: %w", err)
	}
	return nil
}

// loadFields loads note-type field definitions in their stored order.
func (repository *Repository) loadFields(ctx context.Context, s *snapshot) error {
	if err := repository.query(ctx, "SELECT ntid,ord,name FROM fields ORDER BY ntid,ord", func(r *sql.Rows) error {
		var id int64
		var ord int
		var name string
		if err := r.Scan(&id, &ord, &name); err != nil {
			return fmt.Errorf("read field: %w", err)
		}
		nt := s.Types[id]
		if nt == nil || ord != len(nt.Fields) {
			return fmt.Errorf("invalid field ordering for type %d", id)
		}
		nt.Fields = append(nt.Fields, &notetypes.Notetype_Field{Name: name})
		return nil
	}); err != nil {
		return fmt.Errorf("load field layout: %w", err)
	}
	return nil
}

// loadTemplates loads template definitions, appearance, and original configuration bytes.
func (repository *Repository) loadTemplates(ctx context.Context, s *snapshot) error {
	if err := repository.query(ctx, "SELECT ntid,ord,name,config FROM templates ORDER BY ntid,ord", func(r *sql.Rows) error {
		var id int64
		var ord int
		var name string
		var config []byte
		if err := r.Scan(&id, &ord, &name, &config); err != nil {
			return fmt.Errorf("read template: %w", err)
		}
		if s.Types[id] == nil {
			return fmt.Errorf("template references unknown type %d", id)
		}
		s.Types[id].Ordinals[ord] = true
		definition := &notetypes.Notetype_Template_Config{}
		if err := proto.Unmarshal(config, definition); err != nil {
			return fmt.Errorf("parse template configuration: %w", err)
		}
		nt := s.Types[id]
		s.noteTypes[id].templates[ord] = config
		nt.Appearance.Templates = append(nt.Appearance.Templates, &collectionpb.TemplateAppearance{
			Ordinal:         int64(ord),
			Name:            name,
			Front:           definition.QFormat,
			Back:            definition.AFormat,
			BrowserFront:    definition.QFormatBrowser,
			BrowserBack:     definition.AFormatBrowser,
			BrowserFont:     definition.BrowserFontName,
			BrowserFontSize: definition.BrowserFontSize,
		})
		return nil
	}); err != nil {
		return fmt.Errorf("load templates: %w", err)
	}
	return nil
}

// loadDecks loads deck identities and decodes their upstream configuration.
func (repository *Repository) loadDecks(ctx context.Context, s *snapshot) error {
	if err := repository.query(ctx, "SELECT id,name,kind FROM decks ORDER BY name", func(r *sql.Rows) error {
		d := model.Deck{Deck: &collectionpb.Deck{}}
		var kind []byte
		if err := r.Scan(&d.Id, &d.Title, &kind); err != nil {
			return fmt.Errorf("read deck: %w", err)
		}
		d.Title = strings.ReplaceAll(d.Title, "\x1f", "::")
		d.Kind = base64.StdEncoding.EncodeToString(kind)
		s.Decks = append(s.Decks, d)
		return nil
	}); err != nil {
		return fmt.Errorf("load decks: %w", err)
	}
	return nil
}

// loadNotes loads note identities, tags, and ordered field content.
func (repository *Repository) loadNotes(ctx context.Context, s *snapshot) error {
	if err := repository.query(ctx, "SELECT id,guid,mid,tags,flds FROM notes ORDER BY id", func(r *sql.Rows) error {
		n := model.Note{Note: &collectionpb.Note{}}
		var tags, fields string
		if err := r.Scan(&n.NoteId, &n.Guid, &n.NoteTypeId, &tags, &fields); err != nil {
			return fmt.Errorf("read note: %w", err)
		}
		nt := s.Types[n.NoteTypeId]
		if nt == nil {
			return fmt.Errorf("note %d references unknown type", n.NoteId)
		}
		n.NoteType = nt.Name
		n.Tags = strings.Fields(tags)
		n.Bodies = strings.Split(fields, "\x1f")
		if len(n.Bodies) != len(nt.Fields) {
			return fmt.Errorf("note %d field count mismatch", n.NoteId)
		}
		for _, field := range nt.Fields {
			n.OrderedFields = append(n.OrderedFields, field.Name)
		}
		s.Notes[n.NoteId] = n
		return nil
	}); err != nil {
		return fmt.Errorf("load notes: %w", err)
	}
	return nil
}

// loadCards loads managed cards and private references protecting filtered cards.
func (repository *Repository) loadCards(ctx context.Context, s *snapshot) error {
	s.OccupiedDecks = map[int64]bool{}
	deckNames := make(map[int64]string, len(s.Decks))
	normalDecks := make(map[int64]bool, len(s.Decks))
	for _, deck := range s.Decks {
		deckNames[deck.Id] = deck.Title
		kind, err := ankiformat.DecodeDeckKind(deck.Kind)
		if err != nil {
			return fmt.Errorf("read deck %d kind: %w", deck.Id, err)
		}
		normalDecks[deck.Id] = kind.GetNormal() != nil
	}
	if err := repository.query(ctx, "SELECT id,nid,ord,did,type,due,odid FROM cards ORDER BY id", func(r *sql.Rows) error {
		c := &collectionpb.Card{}
		var nid, did, position, homeDeck int64
		var cardType int
		if err := r.Scan(&c.Id, &nid, &c.Ordinal, &did, &cardType, &position, &homeDeck); err != nil {
			return fmt.Errorf("read card: %w", err)
		}
		s.OccupiedDecks[did] = true
		n, ok := s.Notes[nid]
		if !ok {
			return fmt.Errorf("card %d references missing note", c.Id)
		}
		name, ok := deckNames[did]
		if !ok {
			return fmt.Errorf("card %d references missing deck", c.Id)
		}
		if homeDeck != 0 || !normalDecks[did] {
			s.filteredCards = append(s.filteredCards, filteredCard{id: c.Id, noteID: nid, ordinal: int(c.Ordinal), deckID: did, homeDeckID: homeDeck})
			return nil
		}
		if cardType == 0 {
			if old, ok := s.newPositions[nid]; !ok || position < old {
				s.newPositions[nid] = position
			}
			if position >= s.nextPosition {
				s.nextPosition = position + 1
			}
		}
		c.Deck = proto.String(name)
		s.CardDeckIDs[c.Id] = did
		n.Cards = append(n.Cards, c)
		s.Notes[nid] = n
		return nil
	}); err != nil {
		return fmt.Errorf("load cards: %w", err)
	}
	return nil
}

// loadNextPosition loads the collection new-card position cursor.
func (repository *Repository) loadNextPosition(ctx context.Context, s *snapshot) error {
	if err := repository.query(ctx, "SELECT val FROM config WHERE key='nextPos'", func(r *sql.Rows) error {
		var value []byte
		if err := r.Scan(&value); err != nil {
			return fmt.Errorf("read next card position: %w", err)
		}
		var cursor uint32
		if err := json.Unmarshal(value, &cursor); err != nil {
			return fmt.Errorf("decode next card position: %w", err)
		}
		if int64(cursor) > s.nextPosition {
			s.nextPosition = int64(cursor)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("load next card position: %w", err)
	}
	return nil
}

// ValidateDesired validates managed resources and protects unmanaged filtered
// cards, returning the current managed state for planning.
func (repository *Repository) ValidateDesired(ctx context.Context, desired *model.CollectionDesiredState) (model.CollectionState, error) {
	current, err := repository.loadSnapshot(ctx)
	if err != nil {
		return model.CollectionState{}, fmt.Errorf("read current collection: %w", err)
	}
	if err := validateDesired(desired, current); err != nil {
		return model.CollectionState{}, fmt.Errorf("validate desired collection: %w", err)
	}
	return current.CollectionState, nil
}

// validateDesired checks desired resources and rejects conflicts with unmanaged cards.
func validateDesired(desired *model.CollectionDesiredState, current snapshot) error {
	if err := validation.ValidateDesired(desired, current.CollectionState); err != nil {
		return fmt.Errorf("validate managed resources: %w", err)
	}
	notes := make(map[int64]model.Note, len(desired.Notes))
	cards := make(map[int64]bool)
	decks := make(map[int64]bool, len(desired.Decks))
	for _, note := range desired.Notes {
		notes[note.NoteId] = note
		for _, card := range note.Cards {
			cards[card.Id] = true
		}
	}
	for _, deck := range desired.Decks {
		decks[deck.Id] = true
	}
	for _, card := range current.filteredCards {
		note, exists := notes[card.noteID]
		if !exists {
			return fmt.Errorf("cannot remove note %d while it has unmanaged filtered cards; empty the filtered deck in Anki first", card.noteID)
		}
		if cards[card.id] {
			return fmt.Errorf("card %d is in a filtered deck and cannot be managed", card.id)
		}
		for _, declared := range note.Cards {
			if declared.Ordinal == int64(card.ordinal) {
				return fmt.Errorf("note %d already has an unmanaged filtered card at ordinal %d", card.noteID, card.ordinal)
			}
		}
		if !decks[card.deckID] || card.homeDeckID != 0 && !decks[card.homeDeckID] {
			return fmt.Errorf("cannot remove a deck referenced by unmanaged filtered card %d; empty the filtered deck in Anki first", card.id)
		}
	}
	return nil
}

// Reconcile transactionally synchronizes managed resources while retaining unmanaged state.
func (repository *Repository) Reconcile(ctx context.Context, d model.CollectionDesiredState) (err error) {
	s, err := repository.loadSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("read current collection: %w", err)
	}
	if err := validateDesired(&d, s); err != nil {
		return fmt.Errorf("validate desired collection: %w", err)
	}
	plan := planning.CalculatePlan(d, s.CollectionState)
	if len(plan.Changes) == 0 {
		return nil
	}
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reconciliation: %w", err)
	}
	defer func() {
		if e := tx.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback reconciliation: %w", e))
		}
	}()
	run := func(q string, args ...any) error {
		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("reconcile collection row: %w", err)
		}
		return nil
	}
	positions := cardPositions{Next: s.nextPosition}
	now := time.Now().Unix()
	decks, deckIDs, err := reconcileDecks(run, d.Decks, s.Decks, now)
	if err != nil {
		return fmt.Errorf("reconcile decks: %w", err)
	}
	desiredNotes, desiredCards, oldCards, err := reconcileNotes(run, d.Notes, s, decks, &positions, now)
	if err != nil {
		return fmt.Errorf("reconcile notes and cards: %w", err)
	}
	if err := repository.removeMissing(ctx, tx, run, s.CollectionState, desiredNotes, desiredCards, oldCards, deckIDs); err != nil {
		return fmt.Errorf("remove absent resources: %w", err)
	}
	if positions.Advanced {
		if err := run("INSERT INTO config(key,usn,mtime_secs,val) VALUES('nextPos',-1,?,?) ON CONFLICT(key) DO UPDATE SET usn=-1,mtime_secs=excluded.mtime_secs,val=excluded.val", now, []byte(fmt.Sprint(positions.Next))); err != nil {
			return fmt.Errorf("advance new card position: %w", err)
		}
	}
	if err := repository.reconcileAppearance(run, d.Appearance, s, now); err != nil {
		return fmt.Errorf("reconcile card appearance: %w", err)
	}
	if err := run("UPDATE col SET mod=?", time.Now().UnixMilli()); err != nil {
		return fmt.Errorf("mark collection changed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit reconciliation: %w", err)
	}
	var integrity string
	if err := repository.db.QueryRowContext(ctx, "PRAGMA integrity_check").Scan(&integrity); err != nil {
		return fmt.Errorf("verify database integrity: %w", err)
	}
	if integrity != "ok" {
		return fmt.Errorf("database integrity failed: %s", integrity)
	}
	after, err := repository.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("verify reconciled state: %w", err)
	}
	if remaining := planning.CalculatePlan(d, after); len(remaining.Changes) > 0 {
		return fmt.Errorf("reconciliation did not converge (%d remaining changes)", len(remaining.Changes))
	}
	return nil
}

// cardPositions tracks the next available new-card position during reconciliation.
type cardPositions struct {
	// Next is the next unused new-card position.
	Next int64
	// Advanced records whether reconciliation allocated a new-card position.
	Advanced bool
}

// executeStatement executes reconciliation SQL through the active transaction.
type executeStatement func(query string, args ...any) error

// reconcileDecks updates deck declarations and returns allocated identities for reconciliation.
func reconcileDecks(run executeStatement, desiredDecks, currentDecks []model.Deck, now int64) (map[string]int64, map[int64]bool, error) {
	decks := map[string]int64{}
	deckIDs := map[int64]bool{}
	existingDecks := map[int64]model.Deck{}
	for _, x := range currentDecks {
		existingDecks[x.Id] = x
	}
	if err := releaseDeckNames(run, desiredDecks, currentDecks); err != nil {
		return nil, nil, fmt.Errorf("release changing deck names: %w", err)
	}
	for _, deck := range desiredDecks {
		decks[ankiformat.Folded(deck.Title)] = deck.Id
		deckIDs[deck.Id] = true
		kind, err := base64.StdEncoding.DecodeString(deck.Kind)
		if err != nil {
			return nil, nil, fmt.Errorf("decode desired deck: %w", err)
		}
		if old, ok := existingDecks[deck.Id]; !ok {
			if err := run("INSERT INTO decks(id,name,mtime_secs,usn,common,kind) VALUES(?,?,?,-1,?,?)", deck.Id, strings.ReplaceAll(deck.Title, "::", "\x1f"), now, []byte{}, kind); err != nil {
				return nil, nil, fmt.Errorf("add deck %q: %w", deck.Title, err)
			}
			if err := run("DELETE FROM graves WHERE oid=? AND type=?", deck.Id, 2); err != nil {
				return nil, nil, fmt.Errorf("clear restored deck tombstone: %w", err)
			}
		} else if old.Title != deck.Title || old.Kind != deck.Kind {
			if err := run("UPDATE decks SET name=?,kind=?,mtime_secs=?,usn=-1 WHERE id=?", strings.ReplaceAll(deck.Title, "::", "\x1f"), kind, now, deck.Id); err != nil {
				return nil, nil, fmt.Errorf("update deck %q: %w", deck.Title, err)
			}
		}
	}
	return decks, deckIDs, nil
}

// Release names inside the private transaction before assigning final titles.
// Both renames and deletions may otherwise block a valid desired name.
func releaseDeckNames(run executeStatement, desiredDecks, currentDecks []model.Deck) error {
	desiredNames := map[int64]string{}
	reserved := map[string]bool{}
	for _, deck := range desiredDecks {
		desiredNames[deck.Id] = deck.Title
		reserved[ankiformat.Folded(deck.Title)] = true
	}
	for _, deck := range currentDecks {
		reserved[ankiformat.Folded(deck.Title)] = true
	}
	for _, deck := range currentDecks {
		if name, retained := desiredNames[deck.Id]; retained && name == deck.Title {
			continue
		}
		var temporary string
		for attempt := 0; ; attempt++ {
			temporary = fmt.Sprintf("_anki_as_code_pending_%d_%d", deck.Id, attempt)
			if !reserved[ankiformat.Folded(temporary)] {
				break
			}
		}
		reserved[ankiformat.Folded(temporary)] = true
		if err := run("UPDATE decks SET name=? WHERE id=?", temporary, deck.Id); err != nil {
			return fmt.Errorf("temporarily rename deck %d: %w", deck.Id, err)
		}
	}
	return nil
}

// reconcileNotes updates notes and card declarations, assigning positions to new cards.
func reconcileNotes(run executeStatement, notes []model.Note, s snapshot, decks map[string]int64, positions *cardPositions, now int64) (map[int64]bool, map[int64]bool, map[int64]*collectionpb.Card, error) {
	desiredNotes := map[int64]bool{}
	desiredCards := map[int64]bool{}
	oldCards := map[int64]*collectionpb.Card{}
	for _, n := range s.Notes {
		for _, c := range n.Cards {
			oldCards[c.Id] = c
		}
	}
	for _, n := range notes {
		desiredNotes[n.NoteId] = true
		nt := s.Types[n.NoteTypeId]
		fields := strings.Join(n.Bodies, "\x1f")
		sortField := ankiformat.PlainText(n.Bodies[nt.Notetype.Config.SortFieldIdx])
		checksum := sha1.Sum([]byte(ankiformat.PlainText(n.Bodies[0])))
		csum := binary.BigEndian.Uint32(checksum[:4])
		tags := ""
		if len(n.Tags) > 0 {
			tags = " " + strings.Join(n.Tags, " ") + " "
		}
		old, exists := s.Notes[n.NoteId]
		if !exists {
			if err := run("INSERT INTO notes(id,guid,mid,mod,usn,tags,flds,sfld,csum,flags,data) VALUES(?,?,?,?,-1,?,?,?,?,0,'')", n.NoteId, n.Guid, n.NoteTypeId, now, tags, fields, sortField, csum); err != nil {
				return nil, nil, nil, fmt.Errorf("add note %d: %w", n.NoteId, err)
			}
			if err := run("DELETE FROM graves WHERE oid=? AND type=?", n.NoteId, 1); err != nil {
				return nil, nil, nil, fmt.Errorf("clear restored note tombstone: %w", err)
			}
		} else if !slices.Equal(old.Bodies, n.Bodies) || strings.Join(old.Tags, " ") != strings.Join(n.Tags, " ") {
			if err := run("UPDATE notes SET tags=?,flds=?,sfld=?,csum=?,mod=?,usn=-1 WHERE id=?", tags, fields, sortField, csum, now, n.NoteId); err != nil {
				return nil, nil, nil, fmt.Errorf("update note %d: %w", n.NoteId, err)
			}
		}
		for _, tag := range n.Tags {
			if err := run("INSERT OR IGNORE INTO tags(tag,usn,collapsed) VALUES(?,-1,0)", tag); err != nil {
				return nil, nil, nil, fmt.Errorf("register tag: %w", err)
			}
		}
		position, hasPosition := s.newPositions[n.NoteId]
		for _, c := range n.Cards {
			desiredCards[c.Id] = true
			did := decks[ankiformat.Folded(c.GetDeck())]
			old, exists := oldCards[c.Id]
			if !exists {
				if !hasPosition {
					if positions.Next > 2147483647 {
						return nil, nil, nil, fmt.Errorf("new card positions exceed Anki range")
					}
					position = positions.Next
					positions.Next++
					positions.Advanced = true
					hasPosition = true
				}
				if err := run("INSERT INTO cards(id,nid,did,ord,mod,usn,type,queue,due,ivl,factor,reps,lapses,left,odue,odid,flags,data) VALUES(?,?,?,?,?,-1,0,0,?,0,0,0,0,0,0,0,0,'')", c.Id, n.NoteId, did, c.Ordinal, now, position); err != nil {
					return nil, nil, nil, fmt.Errorf("add card %d: %w", c.Id, err)
				}
				if err := run("DELETE FROM graves WHERE oid=? AND type=?", c.Id, 0); err != nil {
					return nil, nil, nil, fmt.Errorf("clear restored card tombstone: %w", err)
				}
			} else if planning.CardNeedsUpdate(old, c, s.CardDeckIDs[c.Id], did) {
				if err := run("UPDATE cards SET did=?,ord=?,mod=?,usn=-1 WHERE id=?", did, c.Ordinal, now, c.Id); err != nil {
					return nil, nil, nil, fmt.Errorf("update card %d: %w", c.Id, err)
				}
			}
		}
	}
	return desiredNotes, desiredCards, oldCards, nil
}

// removeMissing removes absent managed resources and records synchronization tombstones.
func (repository *Repository) removeMissing(ctx context.Context, tx *sql.Tx, run executeStatement, s model.CollectionState, desiredNotes, desiredCards map[int64]bool, oldCards map[int64]*collectionpb.Card, deckIDs map[int64]bool) error {
	for id := range oldCards {
		if !desiredCards[id] {
			if err := run("INSERT OR REPLACE INTO graves(oid,type,usn) VALUES(?,0,-1)", id); err != nil {
				return fmt.Errorf("record card deletion: %w", err)
			}
			if err := run("DELETE FROM cards WHERE id=?", id); err != nil {
				return fmt.Errorf("delete card %d: %w", id, err)
			}
		}
	}
	for id := range s.Notes {
		if !desiredNotes[id] {
			if err := run("INSERT OR REPLACE INTO graves(oid,type,usn) VALUES(?,1,-1)", id); err != nil {
				return fmt.Errorf("record note deletion: %w", err)
			}
			if err := run("DELETE FROM notes WHERE id=?", id); err != nil {
				return fmt.Errorf("delete note %d: %w", id, err)
			}
		}
	}
	for _, deck := range s.Decks {
		if !deckIDs[deck.Id] {
			var cards int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM cards WHERE did=? OR odid=?", deck.Id, deck.Id).Scan(&cards); err != nil {
				return fmt.Errorf("check removed deck: %w", err)
			}
			if cards > 0 {
				return fmt.Errorf("cannot remove deck %q while filtered cards remain in it", deck.Title)
			}
			if err := run("INSERT OR REPLACE INTO graves(oid,type,usn) VALUES(?,2,-1)", deck.Id); err != nil {
				return fmt.Errorf("record deck deletion: %w", err)
			}
			if err := run("DELETE FROM decks WHERE id=?", deck.Id); err != nil {
				return fmt.Errorf("delete deck: %w", err)
			}
		}
	}
	return nil
}

// Open configures SQLite and opens an offline collection. Snapshot and Reconcile
// validate its supported schema before querying collection entities.
func Open(path string) (*Repository, error) {
	db, err := openDB(path)
	if err != nil {
		return nil, fmt.Errorf("open collection repository: %w", err)
	}
	return &Repository{db: db}, nil
}

// Close releases the repository database connection.
func (repository *Repository) Close() error {
	if repository.db == nil {
		return nil
	}
	err := repository.db.Close()
	repository.db = nil
	if err != nil {
		return fmt.Errorf("close collection repository: %w", err)
	}
	return nil
}
