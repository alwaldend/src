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

// Validate rejects bounds the timeline API cannot represent without truncation.
func (w TimeWindow) Validate() error {
	minimum := time.Date(2010, time.November, 6, 0, 0, 0, 0, time.UTC)
	for _, bound := range []time.Time{w.Start, w.End} {
		if !bound.IsZero() && (bound.Before(minimum) || bound.Nanosecond() != 0) {
			return fmt.Errorf("time bounds must be whole seconds on or after 2010-11-06")
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

// ActionKind identifies an executable cleanup operation in a saved plan.
type ActionKind string

const (
	// ActionDeletePost deletes an owned post, including a quote or reply.
	ActionDeletePost ActionKind = "delete_post"
	// ActionUnrepost removes the author's relationship to a source post.
	ActionUnrepost ActionKind = "unrepost"
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
	MediaKeys []string `json:"media_keys"`
	// PollIDs identifies attached polls without loading their contents.
	PollIDs []string `json:"poll_ids"`
	// MediaSourceIDs identifies the sources of reused media.
	MediaSourceIDs []string `json:"media_source_tweet_id"`
}

// TextEntity contains the supported fields of a hashtag, mention, or URL entity.
type TextEntity struct {
	// Start is the entity's starting offset in the returned text.
	Start int64 `json:"start"`
	// End is the entity's exclusive ending offset.
	End int64 `json:"end"`
	// Tag is present on hashtags and cashtags.
	Tag *string `json:"tag"`
	// ID is the mentioned user identifier when returned.
	ID *string `json:"id"`
	// Username is the mentioned user's handle when returned.
	Username *string `json:"username"`
	// URL is the shortened link when returned.
	URL *string `json:"url"`
	// ExpandedURL is the expanded link when returned.
	ExpandedURL *string `json:"expanded_url"`
	// DisplayURL is the link's display form when returned.
	DisplayURL *string `json:"display_url"`
}

// TextEntities groups optional entity arrays for the exposed activity text.
type TextEntities struct {
	// Hashtags contains returned hashtag spans.
	Hashtags []TextEntity `json:"hashtags"`
	// Cashtags contains returned cashtag spans.
	Cashtags []TextEntity `json:"cashtags"`
	// Mentions contains returned user mention spans.
	Mentions []TextEntity `json:"mentions"`
	// URLs contains returned URL spans.
	URLs []TextEntity `json:"urls"`
}

// Reason identifies a matching declaration and one preservation witness seed.
type Reason struct {
	// RuleID uniquely identifies the declaration within this plan.
	RuleID int `json:"rule_id"`
	// Name is the optional policy-supplied label.
	Name string `json:"name,omitempty"`
	// Operation is the direct or predicate selection function.
	Operation string `json:"operation"`
	// SeedID identifies one canonical preservation seed; deletions omit it.
	SeedID string `json:"seed_id,omitempty"`
}

// Decision explains a loaded activity without granting execution authority.
type Decision struct {
	// ItemID is the own activity identifier, including repost activity IDs.
	ItemID string `json:"item_id"`
	// SourceID identifies a repost's original content for informational display.
	SourceID string `json:"source_id,omitempty"`
	// Kind distinguishes a post from a repost activity.
	Kind ItemKind `json:"kind"`
	// CreatedAt is the activity timestamp used by the policy.
	CreatedAt time.Time `json:"created_at"`
	// Text is the snapshot text available when planning.
	Text string `json:"text"`
	// Outcome is unchanged, preserved, delete, or force_delete.
	Outcome string `json:"outcome"`
	// Reasons contains at most one matching reason per declaration.
	Reasons []Reason `json:"reasons"`
}

// Action describes one explicit X mutation.
type Action struct {
	// Kind is delete_post or unrepost.
	Kind ActionKind `json:"kind"`
	// TargetID is a post ID for deletion or a source post ID for unreposting.
	TargetID string `json:"target_id"`
	// Forced reports selection overriding preservation.
	Forced bool `json:"forced"`
}

// Plan is the portable, account-bound list executed by apply.
type Plan struct {
	// Version is the supported serialized schema version.
	Version int `json:"version"`
	// AccountID binds execution to the account which was planned.
	AccountID string `json:"account_id"`
	// EvaluatedAt fixes the time used by every age predicate.
	EvaluatedAt time.Time `json:"evaluated_at"`
	// Actions contains only final selected mutations, in stable order.
	Actions []Action `json:"actions"`
	// Decisions is informational provenance; execution reads only Actions.
	Decisions []Decision `json:"decisions,omitempty"`
}

// Validate rejects every malformed action before any mutation is attempted.
func (p Plan) Validate() error {
	if p.Version != 1 || !ValidID(p.AccountID) || p.EvaluatedAt.IsZero() || p.Actions == nil {
		return fmt.Errorf("invalid plan version, account, evaluation time, or action list")
	}
	seen := make(map[string]bool, len(p.Actions))
	for _, action := range p.Actions {
		if action.Kind != ActionDeletePost && action.Kind != ActionUnrepost {
			return fmt.Errorf("invalid action kind %q", action.Kind)
		}
		if !ValidID(action.TargetID) {
			return fmt.Errorf("invalid target identifier for %s", action.Kind)
		}
		key := string(action.Kind) + ":" + action.TargetID
		if seen[key] {
			return fmt.Errorf("duplicate action %s", key)
		}
		seen[key] = true
	}
	return nil
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
