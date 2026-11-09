package archive

import (
	"context"
	"encoding/json"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/archivepb"
	"google.golang.org/protobuf/proto"
)

// parseAccount reads only the standard X account keys, without protobuf JSON aliases.
func parseAccount(row json.RawMessage) (*archivepb.Account, error) {
	input := sourceRecord(row, "account")
	value := &archivepb.Account{Id: input.identifier("accountId", true), Username: input.optionalString("username"), DisplayName: input.optionalString("accountDisplayName")}
	if *input.failure != nil {
		return nil, fmt.Errorf("import account: %w", *input.failure)
	}
	return value, nil
}

// parseFollowing constructs a canonical relationship from the standard X source record.
func parseFollowing(row json.RawMessage) (*archivepb.Following, error) {
	input := sourceRecord(row, "following")
	value := &archivepb.Following{AccountId: input.identifier("accountId", true), UserLink: input.optionalString("userLink")}
	if *input.failure != nil {
		return nil, fmt.Errorf("import following: %w", *input.failure)
	}
	return value, nil
}

// parseLike constructs a canonical liked-post record from standard X source keys.
func parseLike(row json.RawMessage) (*archivepb.Like, error) {
	input := sourceRecord(row, "like")
	value := &archivepb.Like{PostId: input.identifier("tweetId", true), Text: input.optionalString("fullText"), Url: input.optionalString("expandedUrl")}
	if *input.failure != nil {
		return nil, fmt.Errorf("import like: %w", *input.failure)
	}
	return value, nil
}

// typedRows validates whole datasets before exposing canonical protobuf records.
func typedRows[T proto.Message](ctx context.Context, e *Export, dataset string, parse func(json.RawMessage) (T, error)) (result []T, resultErr error) {
	defer func() {
		if resultErr != nil && e.failure == nil {
			e.failure = resultErr
		}
	}()
	rows, err := e.Read(ctx, dataset)
	if err != nil {
		return nil, fmt.Errorf("read typed %s dataset: %w", dataset, err)
	}
	result = make([]T, 0, len(rows))
	for _, row := range rows {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("decode typed %s dataset: %w", dataset, err)
		}
		value, err := parse(row)
		if err != nil {
			return nil, fmt.Errorf("import typed %s record: %w", dataset, err)
		}
		result = append(result, value)
	}
	return result, nil
}
