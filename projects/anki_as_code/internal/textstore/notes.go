package textstore

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
)

// GenerateIDs gives existing note files new note, card, and synchronization
// identities. Filenames and all other declarations stay unchanged.
func GenerateIDs(ctx context.Context, paths []string) (int, error) {
	count, err := GenerateIDsWithConfig(ctx, paths, "")
	if err != nil {
		return 0, fmt.Errorf("generate text note identities: %w", err)
	}
	return count, nil
}

// GenerateIDsWithConfig reserves identities from an explicit root config when
// selected notes are stored outside its directory tree.
func GenerateIDsWithConfig(ctx context.Context, paths []string, config string) (int, error) {
	files, err := noteFiles(ctx, paths)
	if err != nil {
		return 0, fmt.Errorf("select note files: %w", err)
	}
	if len(files) == 0 {
		return 0, fmt.Errorf("no note TOML files selected")
	}
	selected := make(map[string]model.Note, len(files))
	reserved := make(map[string]model.Note, len(files))
	roots := map[string]bool{}
	if config != "" {
		path, err := collectionConfig(config)
		if err != nil {
			return 0, fmt.Errorf("read explicit identity root: %w", err)
		}
		roots[path] = true
	}
	for _, path := range files {
		if err := ctx.Err(); err != nil {
			return 0, fmt.Errorf("generate identities: %w", err)
		}
		note, err := identityNote(path)
		if err != nil {
			return 0, fmt.Errorf("load selected note: %w", err)
		}
		selected[path], reserved[path] = note, note
		root, err := collectionRoot(path)
		if err != nil {
			return 0, fmt.Errorf("find collection for %q: %w", path, err)
		}
		if root != "" {
			roots[filepath.Join(root, CollectionMarker)] = true
		}
	}
	// Copied notes may still share old IDs, so reserve neighboring identities
	// without requiring the complete desired state to be unique yet.
	for root := range roots {
		neighbors, err := referencedNotes(root)
		if err != nil {
			return 0, fmt.Errorf("scan collection identities: %w", err)
		}
		for _, path := range neighbors {
			if _, exists := reserved[path]; exists {
				continue
			}
			note, err := identityNote(path)
			if err != nil {
				return 0, fmt.Errorf("read neighboring identity: %w", err)
			}
			reserved[path] = note
		}
	}
	next := time.Now().UnixMilli()
	guids := map[string]bool{}
	for _, note := range reserved {
		guids[note.Guid] = true
		ids := []int64{note.NoteId}
		for _, card := range note.Cards {
			ids = append(ids, card.Id)
		}
		for _, id := range ids {
			if id == math.MaxInt64 {
				return 0, fmt.Errorf("note identities exhaust the integer range")
			}
			if id >= next {
				next = id + 1
			}
		}
	}
	var required int64
	for _, note := range selected {
		required += 1 + int64(len(note.Cards))
	}
	if next > math.MaxInt64-required {
		return 0, fmt.Errorf("insufficient integer range for new identities")
	}
	// Finish parsing and generating the entire batch before changing any file.
	for _, path := range files {
		if err := ctx.Err(); err != nil {
			return 0, fmt.Errorf("generate identities: %w", err)
		}
		note := selected[path]
		note.NoteId = next
		next++
		guid, err := newGUID(guids)
		if err != nil {
			return 0, fmt.Errorf("generate identity for %q: %w", path, err)
		}
		note.Guid = guid
		for index := range note.Cards {
			note.Cards[index].Id = next
			next++
		}
		selected[path] = note
	}
	for _, path := range files {
		if err := ctx.Err(); err != nil {
			return 0, fmt.Errorf("generate identities: %w", err)
		}
		if err := WriteNote(path, selected[path]); err != nil {
			return 0, fmt.Errorf("record identities in %q: %w", path, err)
		}
	}
	return len(files), nil
}

// identityNote reads a note declaration eligible for identity regeneration.
func identityNote(path string) (model.Note, error) {
	note, err := ReadNote(path)
	if err != nil {
		return note, fmt.Errorf("parse note %q: %w", path, err)
	}
	if note.NoteTypeId <= 0 || note.Fields == nil {
		return note, fmt.Errorf("%q must declare a note type and named fields", path)
	}
	return note, nil
}

// newGUID generates an unused synchronization GUID and reserves it for the batch.
func newGUID(reserved map[string]bool) (string, error) {
	for {
		guid := rand.Text()
		if !reserved[guid] {
			reserved[guid] = true
			return guid, nil
		}
	}
}

// noteFiles collects a deduplicated batch of note files while observing cancellation.
func noteFiles(ctx context.Context, paths []string) ([]string, error) {
	seen := map[string]bool{}
	for _, input := range paths {
		path, err := ExpandPath(input)
		if err != nil {
			return nil, fmt.Errorf("resolve path %q: %w", input, err)
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, fmt.Errorf("resolve note path %q: %w", path, err)
		}
		if resolved != path {
			return nil, fmt.Errorf("symlink in note path %q", path)
		}
		info, err := os.Lstat(path)
		if err != nil {
			return nil, fmt.Errorf("inspect path %q: %w", path, err)
		}
		if !info.IsDir() && !strings.HasSuffix(path, noteSuffix) {
			return nil, fmt.Errorf("note path %q must end in .note.anki.toml or be a directory", path)
		}
		if err := filepath.WalkDir(path, func(current string, entry fs.DirEntry, walkErr error) error {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("scan note identities: %w", err)
			}
			if walkErr != nil {
				return fmt.Errorf("walk note path %q: %w", current, walkErr)
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlink in note paths %q", current)
			}
			if entry.IsDir() || !strings.HasSuffix(current, noteSuffix) {
				return nil
			}
			if !entry.Type().IsRegular() {
				return fmt.Errorf("note path %q is not a regular file", current)
			}
			seen[current] = true
			return nil
		}); err != nil {
			return nil, fmt.Errorf("collect notes from %q: %w", path, err)
		}
	}
	files := make([]string, 0, len(seen))
	for path := range seen {
		files = append(files, path)
	}
	sort.Strings(files)
	return files, nil
}

// collectionRoot locates an enclosing root collection config for identity reservation.
func collectionRoot(path string) (string, error) {
	for directory := filepath.Dir(path); ; directory = filepath.Dir(directory) {
		manifest := filepath.Join(directory, CollectionMarker)
		info, err := os.Lstat(manifest)
		if err == nil {
			if !info.Mode().IsRegular() {
				return "", fmt.Errorf("collection manifest %q is not a regular file", manifest)
			}
			return directory, nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("inspect collection manifest %q: %w", manifest, err)
		}
		if filepath.Dir(directory) == directory {
			return "", nil
		}
	}
}
