package web

import (
	"strings"
	"sync"
)

// logEntry is a bounded operational line with a monotonic delivery cursor.
type logEntry struct {
	// ID identifies the line across browser polls.
	ID uint64 `json:"id"`
	// Text is a plain-text logger record, never HTML.
	Text string `json:"text"`
}

// Log is a concurrent, bounded writer for the command's shared text logger.
// Its zero value is ready for use; callers retain stderr as a separate sink.
type Log struct {
	// mu protects the cursor and retained entries.
	mu sync.Mutex
	// sequence increases even when earlier entries are discarded.
	sequence uint64
	// entries retains the most recent thousand records.
	entries []logEntry
}

// Write retains one logger write, limiting each browser record to four KiB.
func (l *Log) Write(content []byte) (int, error) {
	text := strings.TrimSuffix(string(content), "\n")
	if len(text) > 4096 {
		text = text[:4096] + " [truncated]"
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sequence++
	entry := logEntry{ID: l.sequence, Text: text}
	if len(l.entries) == 1000 {
		copy(l.entries, l.entries[1:])
		l.entries[len(l.entries)-1] = entry
	} else {
		l.entries = append(l.entries, entry)
	}
	return len(content), nil
}

// since returns retained records after the supplied cursor and the current cursor.
func (l *Log) since(after uint64) ([]logEntry, uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	entries := []logEntry{}
	for _, entry := range l.entries {
		if entry.ID > after {
			entries = append(entries, entry)
		}
	}
	return entries, l.sequence
}
