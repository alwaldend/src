package wasmhost

import (
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/timezone"
)

// offsetForUTC returns host-local seconds east of UTC at the requested instant.
func offsetForUTC(seconds int64) int32 {
	return timezone.OffsetForUTC(time.Local, seconds)
}

// offsetsForLocal packs the earliest/latest offsets for a local civil timestamp.
// Seconds encode its calendar fields as UTC, without applying a timezone yet.
// A MinInt32 offset denotes a nonexistent local time; equal offsets are unique.
func offsetsForLocal(seconds int64) uint64 {
	return timezone.OffsetsForLocal(time.Local, seconds)
}
