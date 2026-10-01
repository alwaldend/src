package xapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Only response headers useful for rate-limit diagnosis belong in APIError.
// Retain every value, but never copy arbitrary headers such as Set-Cookie.
var rateLimitHeaderNames = [...]string{
	"X-Rate-Limit-Limit",
	"X-Rate-Limit-Remaining",
	"X-Rate-Limit-Reset",
	"X-User-Limit-24hour-Limit",
	"X-User-Limit-24hour-Remaining",
	"X-User-Limit-24hour-Reset",
	"Retry-After",
}

const utcTimeLayout = "2006-01-02 15:04:05 UTC"

func copyRateLimitHeaders(source http.Header) http.Header {
	headers := make(http.Header)
	for _, name := range rateLimitHeaderNames {
		if values := source.Values(name); len(values) > 0 {
			headers[http.CanonicalHeaderKey(name)] = append([]string(nil), values...)
		}
	}
	return headers
}

func (e *APIError) rateLimitDetails() string {
	if len(e.Headers) == 0 {
		return ""
	}
	var details strings.Builder
	for _, name := range rateLimitHeaderNames {
		for _, value := range e.Headers.Values(name) {
			visible := value
			if strings.ContainsAny(value, "\r\n\t") {
				visible = strconv.Quote(value)
			}
			fmt.Fprintf(&details, "\n%s: %s", name, visible)
			if strings.HasSuffix(name, "-Reset") {
				if reset, ok := unixReset(value); ok {
					fmt.Fprintf(&details, " (%s)", reset.Format(utcTimeLayout))
				}
			} else if name == "Retry-After" {
				if retry, ok := e.retryAfter(); ok {
					fmt.Fprintf(&details, " (%s; server retry advice)", retry.Format(utcTimeLayout))
				}
			}
		}
	}
	if boundary, ok := e.retryNotBefore(); ok {
		fmt.Fprintf(&details, "\nRate-limit retry not before: %s (service availability is not guaranteed)", boundary.Format(utcTimeLayout))
	} else {
		details.WriteString("\nCurrent retry time is unknown from this response.")
	}
	return details.String()
}

// retryNotBefore describes only a lower bound from known exhausted windows.
// Positive or missing remaining counts never establish an exhausted window;
// malformed counts and untimed exhausted windows prevent a combined boundary.
func (e *APIError) retryNotBefore() (time.Time, bool) {
	if e.ReceivedAt.IsZero() {
		return time.Time{}, false
	}
	var boundary time.Time
	for _, prefix := range [...]string{"X-Rate-Limit-", "X-User-Limit-24hour-"} {
		remainingValues := e.Headers.Values(prefix + "Remaining")
		if len(remainingValues) == 0 {
			continue
		}
		if len(remainingValues) != 1 {
			return time.Time{}, false
		}
		remaining, ok := nonnegativeDecimal(remainingValues[0])
		if !ok {
			return time.Time{}, false
		}
		if remaining != 0 {
			continue
		}
		resets := e.Headers.Values(prefix + "Reset")
		if len(resets) != 1 {
			return time.Time{}, false
		}
		reset, ok := unixReset(resets[0])
		if !ok || !reset.After(e.ReceivedAt) {
			return time.Time{}, false
		}
		if reset.After(boundary) {
			boundary = reset
		}
	}
	if boundary.IsZero() {
		return time.Time{}, false
	}
	if retry, ok := e.retryAfter(); ok && retry.After(boundary) {
		boundary = retry
	}
	return boundary, true
}

// retryAfter accepts either decimal seconds from response receipt or the HTTP
// date formats. An HTTP date contains a comma, so it must not be split as a list.
func (e *APIError) retryAfter() (time.Time, bool) {
	values := e.Headers.Values("Retry-After")
	if len(values) != 1 {
		return time.Time{}, false
	}
	value := strings.TrimSpace(values[0])
	if seconds, ok := nonnegativeDecimal(value); ok {
		const maxSeconds = int64(time.Duration(1<<63-1) / time.Second)
		if e.ReceivedAt.IsZero() || seconds > maxSeconds {
			return time.Time{}, false
		}
		retry := e.ReceivedAt.Add(time.Duration(seconds) * time.Second).UTC()
		// Round relative advice upward so a display with second precision
		// cannot recommend a time before the server's requested delay ends.
		if retry.Nanosecond() != 0 {
			retry = retry.Truncate(time.Second).Add(time.Second)
		}
		return retry, true
	}
	if retry, err := http.ParseTime(value); err == nil {
		return retry.UTC(), true
	}
	return time.Time{}, false
}

func unixReset(value string) (time.Time, bool) {
	seconds, ok := nonnegativeDecimal(value)
	// Keep diagnostic timestamps in the conventional four-digit year range.
	if !ok || seconds > 253402300799 {
		return time.Time{}, false
	}
	return time.Unix(seconds, 0).UTC(), true
}

func nonnegativeDecimal(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return 0, false
		}
	}
	number, err := strconv.ParseInt(value, 10, 64)
	return number, err == nil
}
