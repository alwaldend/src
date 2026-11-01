package versioning

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const DevelopmentVersion = "0.0.0-dev"

var (
	releaseNamePattern     = regexp.MustCompile(`^([1-9][0-9]{3})\.([1-9]|[1-4][0-9]|5[0-3])$`)
	releaseTagPattern      = regexp.MustCompile(`^v([1-9][0-9]{3})\.([1-9]|[1-4][0-9]|5[0-3])\.(0|[1-9][0-9]*)$`)
	nightlyTagPattern      = regexp.MustCompile(`^v([1-9][0-9]{3})\.([1-9]|[1-4][0-9]|5[0-3])\.0-nightly\.([1-9][0-9]{7})$`)
	ownedVersionTagPattern = regexp.MustCompile(`^v[0-9]{4}\.`)
)

type Calendar struct {
	Year int
	Week int
}

func CalendarForDate(value time.Time) Calendar {
	year, week := value.UTC().ISOWeek()
	return Calendar{Year: year, Week: week}
}

func ParseDate(value string) (time.Time, error) {
	result, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date %q as YYYY-MM-DD: %w", value, err)
	}
	if err := validateVersionDate(result); err != nil {
		return time.Time{}, err
	}
	return result.UTC(), nil
}

func validateVersionDate(value time.Time) error {
	value = value.UTC()
	calendar := CalendarForDate(value)
	if value.Year() < 1000 || value.Year() > 9999 ||
		calendar.Year < 1000 || calendar.Year > 9999 {
		return fmt.Errorf(
			"date %s cannot form a four-digit calendar version",
			value.Format("2006-01-02"),
		)
	}
	return nil
}

func (c Calendar) String() string {
	return fmt.Sprintf("%d.%d", c.Year, c.Week)
}

func ParseCalendar(value string) (Calendar, error) {
	match := releaseNamePattern.FindStringSubmatch(value)
	if match == nil {
		return Calendar{}, fmt.Errorf("invalid release %q; want YYYY.W with an unpadded ISO week", value)
	}
	return parseCalendarParts("release", match[1], match[2])
}

type Version struct {
	Calendar Calendar
	Patch    int
	Nightly  string
}

func (v Version) String() string {
	base := fmt.Sprintf("%s.%d", v.Calendar.String(), v.Patch)
	if v.Nightly != "" {
		return base + "-nightly." + v.Nightly
	}
	return base
}

func (v Version) Tag() string {
	return "v" + v.String()
}

func ReleaseVersion(calendar Calendar, patch int) Version {
	return Version{Calendar: calendar, Patch: patch}
}

func NightlyVersion(date time.Time) Version {
	date = date.UTC()
	return Version{
		Calendar: CalendarForDate(date),
		Patch:    0,
		Nightly:  date.Format("20060102"),
	}
}

func ParseVersionTag(value string) (Version, error) {
	if match := releaseTagPattern.FindStringSubmatch(value); match != nil {
		calendar, err := parseCalendarParts("release", match[1], match[2])
		if err != nil {
			return Version{}, fmt.Errorf("invalid version tag %q: %w", value, err)
		}
		patch, err := parseDecimal("release patch", match[3])
		if err != nil {
			return Version{}, fmt.Errorf("invalid version tag %q: %w", value, err)
		}
		return Version{
			Calendar: calendar,
			Patch:    patch,
		}, nil
	}
	match := nightlyTagPattern.FindStringSubmatch(value)
	if match == nil {
		return Version{}, fmt.Errorf("invalid version tag %q", value)
	}
	date, err := time.Parse("20060102", match[3])
	if err != nil {
		return Version{}, fmt.Errorf("invalid nightly date in tag %q: %w", value, err)
	}
	calendar, err := parseCalendarParts("nightly", match[1], match[2])
	if err != nil {
		return Version{}, fmt.Errorf("invalid version tag %q: %w", value, err)
	}
	if CalendarForDate(date) != calendar {
		return Version{}, fmt.Errorf("nightly tag %q does not match its ISO week", value)
	}
	return Version{Calendar: calendar, Patch: 0, Nightly: match[3]}, nil
}

func parseCalendarParts(name string, yearValue string, weekValue string) (Calendar, error) {
	year, err := parseDecimal(name+" year", yearValue)
	if err != nil {
		return Calendar{}, err
	}
	week, err := parseDecimal(name+" week", weekValue)
	if err != nil {
		return Calendar{}, err
	}
	_, lastWeek := time.Date(year, time.December, 28, 0, 0, 0, 0, time.UTC).ISOWeek()
	if week > lastWeek {
		return Calendar{}, fmt.Errorf("ISO year %d has no week %d", year, week)
	}
	return Calendar{Year: year, Week: week}, nil
}

func parseDecimal(name string, value string) (int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s %q: %w", name, value, err)
	}
	return result, nil
}
