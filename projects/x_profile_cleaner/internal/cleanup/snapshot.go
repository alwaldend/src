package cleanup

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/archivepb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
)

// snapshot is the cleanup graph's projection of canonical protobuf records.
// It belongs to the cleanup application, not the archive importer or evaluator API.
type snapshot struct {
	// records contains the owned activities admitted to cleanup graph operations.
	records []model.Item
}

// load returns a bounded cleanup view without source expansion or network reads.
func (s *snapshot) load(ctx context.Context, kind model.ItemKind, window model.TimeWindow) ([]model.Item, error) {
	if err := window.Validate(); err != nil {
		return nil, fmt.Errorf("validate cleanup window: %w", err)
	}
	result := []model.Item{}
	for _, item := range s.records {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("filter cleanup inventory: %w", err)
		}
		if item.Kind != kind || !window.Contains(item.CreatedAt) {
			continue
		}
		if kind == model.ItemRepost && item.SourceID == "" {
			return nil, fmt.Errorf("repost %s has no explicit source post ID", item.ID)
		}
		result = append(result, item)
	}
	return result, nil
}

// cleanupSnapshot derives only the data needed by cleanup graph algorithms.
func cleanupSnapshot(ctx context.Context, values []*archivepb.Tweet) (*snapshot, error) {
	result := &snapshot{records: make([]model.Item, 0, len(values))}
	for _, value := range values {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("project cleanup graph: %w", err)
		}
		item := model.Item{
			ID: value.Id, AuthorID: value.AuthorId, SourceID: value.GetSourceId(),
			EditHistoryIDs: value.EditHistoryIds, CreatedAt: value.CreatedAt.AsTime(),
			Text: value.Text, Language: value.Language, Metrics: map[string]int64{},
			Attachments: &model.Attachments{
				MediaKeys: value.Attachments.GetMediaKeys(), PollIDs: value.Attachments.GetPollIds(), MediaSourceIDs: value.Attachments.GetMediaSourceIds(),
			},
			Entities: &model.TextEntities{
				Hashtags: graphEntities(value.Entities.GetHashtags()), Cashtags: graphEntities(value.Entities.GetCashtags()),
				Mentions: graphEntities(value.Entities.GetMentions()), URLs: graphEntities(value.Entities.GetUrls()),
			},
		}
		switch value.Kind {
		case archivepb.ActivityKind_POST:
			item.Kind = model.ItemPost
		case archivepb.ActivityKind_REPOST:
			item.Kind = model.ItemRepost
		default:
			return nil, fmt.Errorf("unsupported cleanup activity kind")
		}
		counter := func(name string, value *int64) {
			if value != nil {
				item.Metrics[name] = *value
			}
		}
		if metrics := value.Metrics; metrics != nil {
			counter("like_count", metrics.LikeCount)
			counter("repost_count", metrics.RepostCount)
			counter("reply_count", metrics.ReplyCount)
			counter("quote_count", metrics.QuoteCount)
			counter("bookmark_count", metrics.BookmarkCount)
			counter("impression_count", metrics.ImpressionCount)
		}
		for _, reference := range value.References {
			kind := model.ReferenceKind("")
			switch reference.Kind {
			case archivepb.ReferenceKind_REPLY:
				kind = model.ReferenceReply
			case archivepb.ReferenceKind_QUOTE:
				kind = model.ReferenceQuote
			case archivepb.ReferenceKind_REPOST_SOURCE:
				kind = model.ReferenceRepost
			default:
				return nil, fmt.Errorf("unsupported cleanup reference kind")
			}
			item.References = append(item.References, model.Reference{ID: reference.Id, Kind: kind})
		}
		result.records = append(result.records, item)
	}
	return result, nil
}

// graphEntities projects canonical spans for cleanup callbacks without rereading JSON.
func graphEntities(values []*archivepb.TextEntity) []model.TextEntity {
	result := make([]model.TextEntity, 0, len(values))
	for _, value := range values {
		result = append(result, model.TextEntity{Start: value.Start, End: value.End, Tag: value.Tag, ID: value.Id, Username: value.Username, URL: value.Url, ExpandedURL: value.ExpandedUrl, DisplayURL: value.DisplayUrl})
	}
	return result
}
