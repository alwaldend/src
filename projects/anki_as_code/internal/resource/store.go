package resource

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/deck_config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/ankitects/anki/proto/anki/search"
	"google.golang.org/grpc"
)

// Store shares native clients and resource reads on a borrowed prepared connection.
// Consumers use generated client methods directly and own mutation policy and lifetime.
type Store struct {
	// Decks exposes native deck operations for resource consumers.
	Decks decks.DecksServiceClient
	// NoteTypes exposes Anki's note-type structure and presentation operations.
	NoteTypes notetypes.NotetypesServiceClient
	// Notes exposes native content operations and card generation.
	Notes notes.NotesServiceClient
	// Cards exposes generated-card reads and native updates.
	Cards cards.CardsServiceClient
	// Search enumerates native identities and supports caller queries.
	Search search.SearchServiceClient
	// Presets reads existing scheduling references; consumers own managed scope.
	Presets deck_config.DeckConfigServiceClient
}

// New constructs reusable generated clients without native work or lifecycle ownership.
func New(connection grpc.ClientConnInterface) *Store {
	return &Store{
		Decks:     decks.NewDecksServiceClient(connection),
		NoteTypes: notetypes.NewNotetypesServiceClient(connection),
		Notes:     notes.NewNotesServiceClient(connection),
		Cards:     cards.NewCardsServiceClient(connection),
		Search:    search.NewSearchServiceClient(connection),
		Presets:   deck_config.NewDeckConfigServiceClient(connection),
	}
}

// Snapshot reads fresh native values without writing resources or collection state.
// Results belong to the caller; errors return no partial snapshot. The caller
// supplies exclusive offline access because sequential reads are not transactional.
func (store *Store) Snapshot(ctx context.Context) (*Snapshot, error) {
	result := &Snapshot{Presets: map[int64]string{}}
	types, err := store.NoteTypes.GetNotetypeNames(ctx, &generic.Empty{})
	if err != nil {
		return nil, fmt.Errorf("list note types: %w", err)
	}
	for _, entry := range types.Entries {
		value, err := store.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: entry.Id})
		if err != nil {
			return nil, fmt.Errorf("get note type %d: %w", entry.Id, err)
		}
		result.NoteTypes = append(result.NoteTypes, value)
	}
	slices.SortFunc(result.NoteTypes, func(a, b *notetypes.Notetype) int { return cmp.Compare(a.Id, b.Id) })
	listed, err := store.Decks.GetDeckNames(ctx, &decks.GetDeckNamesRequest{IncludeFiltered: true})
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	for _, entry := range listed.Entries {
		value, err := store.Decks.GetDeck(ctx, &decks.DeckId{Did: entry.Id})
		if err != nil {
			return nil, fmt.Errorf("get deck %d: %w", entry.Id, err)
		}
		result.Decks = append(result.Decks, value)
	}
	slices.SortFunc(result.Decks, func(a, b *decks.Deck) int { return strings.Compare(a.Name, b.Name) })
	for _, deck := range result.Decks {
		if deck.GetNormal() == nil {
			continue
		}
		values, err := store.Presets.GetDeckConfigsForUpdate(ctx, &decks.DeckId{Did: deck.Id})
		if err != nil {
			return nil, fmt.Errorf("read scheduling preset references for deck %d: %w", deck.Id, err)
		}
		for _, value := range values.AllConfig {
			result.Presets[value.Config.Id] = value.Config.Name
		}
		break
	}
	noteIDs, err := store.Search.SearchNotes(ctx, &search.SearchRequest{})
	if err != nil {
		return nil, fmt.Errorf("search notes: %w", err)
	}
	slices.Sort(noteIDs.Ids)
	for _, id := range noteIDs.Ids {
		value, err := store.Notes.GetNote(ctx, &notes.NoteId{Nid: id})
		if err != nil {
			return nil, fmt.Errorf("get note %d: %w", id, err)
		}
		result.Notes = append(result.Notes, value)
	}
	cardIDs, err := store.Search.SearchCards(ctx, &search.SearchRequest{})
	if err != nil {
		return nil, fmt.Errorf("search cards: %w", err)
	}
	slices.Sort(cardIDs.Ids)
	for _, id := range cardIDs.Ids {
		value, err := store.Cards.GetCard(ctx, &cards.CardId{Cid: id})
		if err != nil {
			return nil, fmt.Errorf("get card %d: %w", id, err)
		}
		result.Cards = append(result.Cards, value)
	}
	return result, nil
}
