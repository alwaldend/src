package planning

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"google.golang.org/protobuf/proto"
)

// CardNeedsUpdate reports whether a managed card changes ordinal or destination deck.
func CardNeedsUpdate(current, desired *collectionpb.Card, currentDeckID, desiredDeckID int64) bool {
	return current.Ordinal != desired.Ordinal || currentDeckID != desiredDeckID
}

// CalculatePlan compares desired resources with the snapshot and returns their change plan.
func CalculatePlan(d model.CollectionDesiredState, s model.CollectionState) *collectionpb.CollectionChangePlan {
	p := &collectionpb.CollectionChangePlan{Changes: []*collectionpb.CollectionChange{}}
	add := func(resource, id, action string) {
		p.Changes = append(p.Changes, &collectionpb.CollectionChange{Resource: resource, Id: id, Action: action})
	}
	notes := map[int64]model.Note{}
	cards := map[int64]*collectionpb.Card{}
	oldCards := map[int64]*collectionpb.Card{}
	desiredDeckIDs := make(map[string]int64, len(d.Decks))
	for _, deck := range d.Decks {
		desiredDeckIDs[ankiformat.Folded(deck.Title)] = deck.Id
	}
	for _, n := range s.Notes {
		for _, c := range n.Cards {
			oldCards[c.Id] = c
		}
	}
	for _, n := range d.Notes {
		notes[n.NoteId] = n
		old, ok := s.Notes[n.NoteId]
		id := fmt.Sprint(n.NoteId)
		if !ok {
			p.NotesAdded++
			add("note", id, "add")
		} else if !slices.Equal(old.Bodies, n.Bodies) || strings.Join(old.Tags, " ") != strings.Join(n.Tags, " ") {
			p.NotesChanged++
			add("note", id, "update")
		}
		for _, c := range n.Cards {
			cards[c.Id] = c
			old, ok := oldCards[c.Id]
			if !ok {
				p.CardsAdded++
				add("card", fmt.Sprint(c.Id), "add")
			} else if CardNeedsUpdate(old, c, s.CardDeckIDs[c.Id], desiredDeckIDs[ankiformat.Folded(c.GetDeck())]) {
				p.CardsChanged++
				add("card", fmt.Sprint(c.Id), "update")
			}
		}
	}
	for id := range s.Notes {
		if _, ok := notes[id]; !ok {
			p.NotesDeleted++
			add("note", fmt.Sprint(id), "delete")
		}
	}
	for id := range oldCards {
		if _, ok := cards[id]; !ok {
			p.CardsDeleted++
			add("card", fmt.Sprint(id), "delete")
		}
	}
	decks := map[int64]model.Deck{}
	oldDecks := map[int64]model.Deck{}
	for _, x := range s.Decks {
		oldDecks[x.Id] = x
	}
	for _, x := range d.Decks {
		decks[x.Id] = x
		old, ok := oldDecks[x.Id]
		if !ok {
			p.DecksAdded++
			add("deck", x.Title, "add")
		} else if old.Title != x.Title || old.Kind != x.Kind {
			p.DecksChanged++
			add("deck", x.Title, "update")
		}
	}
	for id, x := range oldDecks {
		if _, ok := decks[id]; !ok {
			p.DecksDeleted++
			add("deck", x.Title, "delete")
		}
	}
	for _, appearance := range d.Appearance {
		if !proto.Equal(appearance, s.Types[appearance.Id].Appearance) {
			p.AppearanceChanged++
			add("appearance", fmt.Sprint(appearance.Id), "update")
		}
	}

	sort.Slice(p.Changes, func(i, j int) bool {
		a, b := p.Changes[i], p.Changes[j]
		return a.Resource+"/"+a.Id < b.Resource+"/"+b.Id
	})
	return p
}
