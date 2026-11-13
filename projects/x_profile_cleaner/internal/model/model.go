package model

import (
	"fmt"
	"time"
)

// TimeWindow bounds discovery by activity time; zero bounds are independently open.
type TimeWindow struct {
	// Start includes activities at this instant or later when nonzero.
	Start time.Time
	// End excludes activities at this instant or later when nonzero.
	End time.Time
}

// Validate keeps archive windows ordered and expressed at whole seconds.
func (w TimeWindow) Validate() error {
	for _, bound := range []time.Time{w.Start, w.End} {
		if !bound.IsZero() && bound.Nanosecond() != 0 {
			return fmt.Errorf("time bounds must be whole seconds")
		}
	}
	if !w.Start.IsZero() && !w.End.IsZero() && !w.Start.Before(w.End) {
		return fmt.Errorf("time window start must precede its end")
	}
	return nil
}

// Contains applies inclusive-start/exclusive-end semantics for adjacent chunks.
func (w TimeWindow) Contains(stamp time.Time) bool {
	return (w.Start.IsZero() || !stamp.Before(w.Start)) && (w.End.IsZero() || stamp.Before(w.End))
}

// ItemKind distinguishes author-owned posts from repost activities.
type ItemKind string

const (
	// ItemPost includes original posts, quotes, and replies.
	ItemPost ItemKind = "post"
	// ItemRepost identifies an author's repost activity, not its source post.
	ItemRepost ItemKind = "repost"
)

// ReferenceKind identifies the relationship to a referenced post.
type ReferenceKind string

const (
	// ReferenceQuote identifies the post quoted by an activity.
	ReferenceQuote ReferenceKind = "quoted"
	// ReferenceReply identifies the post answered by a reply.
	ReferenceReply ReferenceKind = "replied_to"
	// ReferenceRepost identifies the source of a repost activity.
	ReferenceRepost ReferenceKind = "retweeted"
)

// Reference identifies an original post referenced by an author-owned activity.
type Reference struct {
	// ID is the referenced post's decimal identifier.
	ID string
	// Kind is quoted, replied_to, or retweeted.
	Kind ReferenceKind
}

// Item is an author-owned post or repost activity, never foreign source content.
type Item struct {
	// ID identifies this post or repost activity.
	ID string
	// EditHistoryIDs identifies versions of this same activity, oldest first.
	EditHistoryIDs []string
	// AuthorID identifies the account which performed the activity.
	AuthorID string
	// Kind distinguishes post from repost.
	Kind ItemKind
	// SourceID identifies the original post for a repost removal.
	SourceID string
	// CreatedAt timestamps this activity, not its reposted source.
	CreatedAt time.Time
	// Text contains the complete available text.
	Text string
	// Language is the returned language code, or nil when absent.
	Language *string
	// Attachments describes returned media and poll identifiers without expansion.
	Attachments *Attachments
	// Entities describes the exposed text; nil means metadata was not returned.
	Entities *TextEntities
	// Metrics contains only counters actually returned by X.
	Metrics map[string]int64
	// References contains the activity's typed references.
	References []Reference
}

// Attachments contains optional attachment identifiers returned on an own activity.
type Attachments struct {
	// MediaKeys identifies attached media without loading the media objects.
	MediaKeys []string
	// PollIDs identifies attached polls without loading their contents.
	PollIDs []string
	// MediaSourceIDs identifies the sources of reused media.
	MediaSourceIDs []string
}

// TextEntity contains the supported fields of a hashtag, mention, or URL entity.
type TextEntity struct {
	// Start is the entity's starting offset in the returned text.
	Start int64
	// End is the entity's exclusive ending offset.
	End int64
	// Tag is present on hashtags and cashtags.
	Tag *string
	// ID is the mentioned user identifier when returned.
	ID *string
	// Username is the mentioned user's handle when returned.
	Username *string
	// URL is the shortened link when returned.
	URL *string
	// ExpandedURL is the expanded link when returned.
	ExpandedURL *string
	// DisplayURL is the link's display form when returned.
	DisplayURL *string
}

// TextEntities groups optional entity arrays for the exposed activity text.
type TextEntities struct {
	// Hashtags contains returned hashtag spans.
	Hashtags []TextEntity
	// Cashtags contains returned cashtag spans.
	Cashtags []TextEntity
	// Mentions contains returned user mention spans.
	Mentions []TextEntity
	// URLs contains returned URL spans.
	URLs []TextEntity
}

// ValidID accepts canonical nonzero decimal X identifiers without numeric coercion.
func ValidID(id string) bool {
	if len(id) == 0 || len(id) > 19 || id[0] == '0' {
		return false
	}
	for _, c := range id {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
