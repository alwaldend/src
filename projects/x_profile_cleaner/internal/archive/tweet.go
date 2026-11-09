package archive

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/archivepb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// parseTweet constructs the sole typed record directly from standard X source keys.
func parseTweet(row json.RawMessage, account string) (*archivepb.Tweet, error) {
	input := sourceRecord(row, "tweet")
	value := &archivepb.Tweet{
		Id: input.identifier("id_str", true), AuthorId: account,
		Kind: archivepb.ActivityKind_POST, Text: input.requiredString("full_text"),
		Language: input.optionalString("lang"), EditHistoryIds: input.strings("edit_history_tweet_ids"),
		Metrics: &archivepb.Metrics{
			LikeCount: input.integer("favorite_count"), RepostCount: input.integer("retweet_count"),
			ReplyCount: input.integer("reply_count"), QuoteCount: input.integer("quote_count"),
			BookmarkCount: input.integer("bookmark_count"), ImpressionCount: input.integer("impression_count"),
		},
		Attachments: &archivepb.Attachments{}, Entities: &archivepb.TextEntities{},
	}
	author := input.text("user_id_str")
	if author != "" && author != account {
		input.fail(fmt.Errorf("foreign export activity"))
	}
	created, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", input.requiredString("created_at"))
	if err != nil {
		input.fail(fmt.Errorf("parse tweet timestamp: %w", err))
	}
	if created.IsZero() {
		input.fail(fmt.Errorf("invalid tweet timestamp"))
	}
	value.CreatedAt = timestamppb.New(created)
	if err := value.CreatedAt.CheckValid(); err != nil {
		input.fail(fmt.Errorf("validate tweet timestamp: %w", err))
	}
	reply := input.identifier("in_reply_to_status_id_str", false)
	quote := input.identifier("quoted_status_id_str", false)
	repost := input.identifier("retweeted_status_id_str", false)
	if reply != "" {
		value.References = append(value.References, &archivepb.Reference{Id: reply, Kind: archivepb.ReferenceKind_REPLY})
	}
	if quote != "" {
		value.References = append(value.References, &archivepb.Reference{Id: quote, Kind: archivepb.ReferenceKind_QUOTE})
	}
	original := input.object("retweeted_status", false)
	if original.fields != nil {
		target := original.identifier("id_str", true)
		if repost != "" && repost != target {
			input.fail(fmt.Errorf("ambiguous repost source"))
		}
		repost = target
	}
	if repost != "" {
		value.Kind = archivepb.ActivityKind_REPOST
		value.SourceId = &repost
		value.References = append(value.References, &archivepb.Reference{Id: repost, Kind: archivepb.ReferenceKind_REPOST_SOURCE})
	} else if strings.HasPrefix(value.Text, "RT @") {
		value.Kind = archivepb.ActivityKind_REPOST
	}
	groups := input.object("entities", false)
	value.Entities.Hashtags = parseEntities(groups, "hashtags")
	value.Entities.Cashtags = parseEntities(groups, "symbols")
	value.Entities.Mentions = parseEntities(groups, "user_mentions")
	value.Entities.Urls = parseEntities(groups, "urls")
	media := groups
	if extended := input.object("extended_entities", false); extended.fields != nil {
		media = extended
	}
	parseMedia(media, value.Attachments)
	if input.boolean("is_quote_status") && quote == "" {
		targets := map[string]bool{}
		for _, entity := range value.Entities.Urls {
			if target := statusID(entity.GetExpandedUrl()); target != "" {
				targets[target] = true
			}
		}
		if len(targets) != 1 {
			input.fail(fmt.Errorf("quoted post has no unambiguous source ID"))
		}
		for target := range targets {
			reference := &archivepb.Reference{Id: target, Kind: archivepb.ReferenceKind_QUOTE}
			if repost != "" {
				last := len(value.References) - 1
				value.References = append(value.References[:last], reference, value.References[last])
			} else {
				value.References = append(value.References, reference)
			}
		}
	}
	if *input.failure != nil {
		return nil, fmt.Errorf("import tweet: %w", *input.failure)
	}
	return value, nil
}

// parseEntities reads supported source spans directly into protobuf messages.
func parseEntities(group objectReader, name string) []*archivepb.TextEntity {
	var result []*archivepb.TextEntity
	for _, row := range group.array(name) {
		if *group.failure != nil {
			break
		}
		input := readObject(row, group.failure, true)
		offsets := input.array("indices")
		if len(offsets) != 2 {
			input.fail(fmt.Errorf("entity requires two offsets"))
			break
		}
		start, err := number(offsets[0])
		if err != nil {
			input.fail(fmt.Errorf("read entity start: %w", err))
		}
		end, err := number(offsets[1])
		if err != nil {
			input.fail(fmt.Errorf("read entity end: %w", err))
		}
		if start > end {
			input.fail(fmt.Errorf("invalid entity offset order"))
		}
		result = append(result, &archivepb.TextEntity{
			Start: start, End: end, Tag: input.optionalString("text"), Id: input.optionalString("id_str"),
			Username: input.optionalString("screen_name"), Url: input.optionalString("url"),
			ExpandedUrl: input.optionalString("expanded_url"), DisplayUrl: input.optionalString("display_url"),
		})
	}
	return result
}

// parseMedia imports X media identifiers without reading bundled binary files.
func parseMedia(group objectReader, result *archivepb.Attachments) {
	for _, row := range group.array("media") {
		if *group.failure != nil {
			break
		}
		input := readObject(row, group.failure, true)
		id := input.identifier("id_str", true)
		key := input.text("media_key")
		if len(key) > 128 {
			input.fail(fmt.Errorf("invalid media key"))
		}
		if key == "" {
			key = id
		}
		result.MediaKeys = append(result.MediaKeys, key)
		if source := input.identifier("source_status_id_str", false); source != "" {
			result.MediaSourceIds = append(result.MediaSourceIds, source)
		}
	}
}

// statusID recognizes an explicit X status link without requesting its content.
func statusID(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || (parsed.Host != "x.com" && parsed.Host != "twitter.com" && parsed.Host != "www.twitter.com") {
		return ""
	}
	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(segments) < 3 || segments[len(segments)-2] != "status" {
		return ""
	}
	id := segments[len(segments)-1]
	if !model.ValidID(id) {
		return ""
	}
	return id
}
