package timezone

import (
	"math"
	"time"
	_ "time/tzdata" // Keep named host timezones available without system zoneinfo files.
)

// OffsetForUTC returns seconds east of UTC at the requested instant in location.
// Native and browser timezone imports must supply a non-nil location.
func OffsetForUTC(location *time.Location, seconds int64) int32 {
	_, offset := time.Unix(seconds, 0).In(location).Zone()
	return int32(offset)
}

// OffsetsForLocal packs the earliest/latest offsets for a local civil timestamp.
// Seconds encode its calendar fields as UTC, without applying a timezone yet.
// A MinInt32 offset denotes a nonexistent local time; equal offsets are unique.
// Native and browser timezone imports must supply a non-nil location.
func OffsetsForLocal(location *time.Location, seconds int64) uint64 {
	const day = 24 * 60 * 60
	first, last := int32(math.MinInt32), int32(math.MinInt32)
	if seconds >= math.MinInt64+day && seconds <= math.MaxInt64-day {
		limit := time.Unix(seconds+day, 0)
		for cursor := time.Unix(seconds-day, 0); !cursor.After(limit); {
			local := cursor.In(location)
			_, offset := local.Zone()
			// Chrono offsets are strictly within one day of UTC. Walk every
			// intersecting zone interval so gaps, overlaps and skipped dates
			// work even when a transition is not a one-hour DST adjustment.
			if offset > -day && offset < day {
				_, actual := time.Unix(seconds-int64(offset), 0).In(location).Zone()
				if actual == offset {
					value := int32(offset)
					if first == math.MinInt32 || value > first {
						first = value
					}
					if last == math.MinInt32 || value < last {
						last = value
					}
				}
			}
			_, next := local.ZoneBounds()
			if next.IsZero() || !next.After(cursor) || next.After(limit) {
				break
			}
			cursor = next
		}
	}
	return uint64(uint32(first))<<32 | uint64(uint32(last))
}
