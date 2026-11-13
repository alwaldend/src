package worker

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// wireTime requires a present, representable protobuf timestamp.
func wireTime(value *timestamppb.Timestamp) (time.Time, error) {
	if value == nil {
		return time.Time{}, fmt.Errorf("missing timestamp")
	}
	if err := value.CheckValid(); err != nil {
		return time.Time{}, fmt.Errorf("validate timestamp: %w", err)
	}
	return value.AsTime(), nil
}
