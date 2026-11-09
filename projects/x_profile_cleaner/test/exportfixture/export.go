package exportfixture

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
)

// Encode writes independently specified records in the literal YTD export format.
func Encode(account string, items []model.Item) ([]byte, error) {
	rows := []any{}
	for _, item := range items {
		tweet := map[string]any{"id_str": item.ID, "user_id_str": item.AuthorID, "created_at": item.CreatedAt.Format("Mon Jan 02 15:04:05 -0700 2006"), "full_text": item.Text}
		if item.EditHistoryIDs != nil {
			tweet["edit_history_tweet_ids"] = item.EditHistoryIDs
		}
		if item.Language != nil {
			tweet["lang"] = *item.Language
		}
		for source, target := range map[string]string{"like_count": "favorite_count", "repost_count": "retweet_count", "reply_count": "reply_count", "quote_count": "quote_count", "bookmark_count": "bookmark_count", "impression_count": "impression_count"} {
			if value, ok := item.Metrics[source]; ok {
				tweet[target] = fmt.Sprint(value)
			}
		}
		for _, reference := range item.References {
			switch reference.Kind {
			case model.ReferenceReply:
				tweet["in_reply_to_status_id_str"] = reference.ID
			case model.ReferenceQuote:
				tweet["quoted_status_id_str"] = reference.ID
			case model.ReferenceRepost:
				tweet["retweeted_status_id_str"] = reference.ID
			}
		}
		if item.Kind == model.ItemRepost {
			if item.SourceID != "" {
				tweet["retweeted_status_id_str"] = item.SourceID
			} else if !strings.HasPrefix(item.Text, "RT @") {
				tweet["full_text"] = "RT @fixture: " + item.Text
			}
		}
		if item.Entities != nil {
			entities := map[string]any{}
			for key, values := range map[string][]model.TextEntity{"hashtags": item.Entities.Hashtags, "symbols": item.Entities.Cashtags, "user_mentions": item.Entities.Mentions, "urls": item.Entities.URLs} {
				if values == nil {
					continue
				}
				entries := []any{}
				for _, value := range values {
					entry := map[string]any{"indices": []string{fmt.Sprint(value.Start), fmt.Sprint(value.End)}}
					for name, text := range map[string]*string{"text": value.Tag, "id_str": value.ID, "screen_name": value.Username, "url": value.URL, "expanded_url": value.ExpandedURL, "display_url": value.DisplayURL} {
						if text != nil {
							entry[name] = *text
						}
					}
					entries = append(entries, entry)
				}
				entities[key] = entries
			}
			tweet["entities"] = entities
		}
		if item.Attachments != nil {
			media := []any{}
			for index, key := range item.Attachments.MediaKeys {
				id := key
				if _, suffix, ok := strings.Cut(key, "_"); ok {
					id = suffix
				}
				entry := map[string]any{"id_str": id, "media_key": key}
				if index < len(item.Attachments.MediaSourceIDs) {
					entry["source_status_id_str"] = item.Attachments.MediaSourceIDs[index]
				}
				media = append(media, entry)
			}
			tweet["extended_entities"] = map[string]any{"media": media}
		}
		rows = append(rows, map[string]any{"tweet": tweet})
	}
	accountData, err := json.Marshal([]any{map[string]any{"account": map[string]string{"accountId": account, "username": "fixture"}}})
	if err != nil {
		return nil, fmt.Errorf("encode fixture account: %w", err)
	}
	tweetData, err := json.Marshal(rows)
	if err != nil {
		return nil, fmt.Errorf("encode fixture tweets: %w", err)
	}
	var data bytes.Buffer
	writer := zip.NewWriter(&data)
	for name, content := range map[string]string{"data/account.js": "window.YTD.account.part0 = " + string(accountData) + ";", "data/tweets.js": "window.YTD.tweets.part0 = " + string(tweetData) + ";"} {
		file, err := writer.Create(name)
		if err != nil {
			return nil, fmt.Errorf("create fixture ZIP entry: %w", err)
		}
		if _, err := io.WriteString(file, content); err != nil {
			return nil, fmt.Errorf("write fixture ZIP entry: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("finish fixture ZIP: %w", err)
	}
	return data.Bytes(), nil
}
