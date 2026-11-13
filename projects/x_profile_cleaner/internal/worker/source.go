package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
)

// archiveObject is a read-only, seekable object with a provider-owned lifetime.
type archiveObject struct {
	// reader supplies ZIP random access without passing data through the backend.
	reader io.ReaderAt
	// size is the provider-verified length used by bounded ZIP decoding.
	size int64
	// close releases the file, cache or download owned by the provider.
	close func() error
}

// archiveSource interprets object keys in one operator-configured namespace.
// Future object-store providers can supply the same seekable object contract.
type archiveSource interface {
	validate(key string) error
	open(ctx context.Context, key string) (*archiveObject, error)
	close() error
}

// localSource confines object access to one directory capability.
type localSource struct {
	// root enforces path confinement even through symlinks and rename races.
	root *os.Root
}

// validate applies local path rules rather than imposing them on all providers.
func (s *localSource) validate(key string) error {
	if !validArchiveKey(key) || strings.ContainsAny(key, "\\:") || strings.HasPrefix(key, "/") || path.Clean(key) != key || key == "." || key == ".." || strings.HasPrefix(key, "../") {
		return fmt.Errorf("invalid local archive key")
	}
	return nil
}

// open returns a regular read-only file without blocking on named pipes.
func (s *localSource) open(ctx context.Context, key string) (*archiveObject, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("before source access: %w", err)
	}
	if err := s.validate(key); err != nil {
		return nil, fmt.Errorf("validate local key: %w", err)
	}
	file, err := s.root.OpenFile(key, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("open source object: %w", err)
	}
	info, err := file.Stat()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("inspect source object: %w", err), file.Close())
	}
	if !info.Mode().IsRegular() {
		return nil, errors.Join(fmt.Errorf("source object is not a regular file"), file.Close())
	}
	return &archiveObject{reader: file, size: info.Size(), close: file.Close}, nil
}

// close releases the local namespace capability after jobs finish.
func (s *localSource) close() error {
	if err := s.root.Close(); err != nil {
		return fmt.Errorf("close local source: %w", err)
	}
	return nil
}

// sources owns immutable, operator-configured source providers.
type sources struct {
	// providers resolves logical source names without exposing storage roots to jobs.
	providers map[string]archiveSource
}

// openSources resolves configured directory capabilities before accepting jobs.
func openSources(configured map[string]string) (*sources, error) {
	if len(configured) == 0 || len(configured) > 16 {
		return nil, fmt.Errorf("configure between one and sixteen archive sources")
	}
	result := &sources{providers: map[string]archiveSource{}}
	for name, directory := range configured {
		if !validSourceName(name) || !filepath.IsAbs(directory) {
			return nil, errors.Join(fmt.Errorf("archive sources require a bounded name and absolute directory"), result.close())
		}
		root, err := os.OpenRoot(directory)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("open archive source %q: %w", name, err), result.close())
		}
		result.providers[name] = &localSource{root: root}
	}
	return result, nil
}

// validate applies provider-specific key rules without reading an archive.
func (s *sources) validate(name, key string) error {
	provider := s.providers[name]
	if provider == nil {
		return fmt.Errorf("unknown archive source or invalid archive key")
	}
	if err := provider.validate(key); err != nil {
		return fmt.Errorf("validate source key: %w", err)
	}
	return nil
}

// load opens a validated export and transfers its source lifetime to the caller.
func (s *sources) load(ctx context.Context, name, key, account string) (*archive.Export, func() error, error) {
	if err := s.validate(name, key); err != nil {
		return nil, nil, fmt.Errorf("select archive source: %w", err)
	}
	object, err := s.providers[name].open(ctx, key)
	if err != nil {
		return nil, nil, fmt.Errorf("open source archive: %w", err)
	}
	data, err := archive.Open(ctx, object.reader, object.size)
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("decode source archive: %w", err), object.close())
	}
	if data.AccountID != account {
		return nil, nil, errors.Join(fmt.Errorf("archive account does not match evaluation account"), object.close())
	}
	return data, object.close, nil
}

// close releases all directory capabilities after active jobs have stopped.
func (s *sources) close() error {
	var result error
	for name, provider := range s.providers {
		if err := provider.close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close archive source %q: %w", name, err))
		}
	}
	return result
}

// validSourceName accepts simple operator-owned aliases, independent of paths.
func validSourceName(name string) bool {
	if len(name) == 0 || len(name) > 64 {
		return false
	}
	for _, character := range name {
		if !(character >= 'a' && character <= 'z') && !(character >= 'A' && character <= 'Z') && !(character >= '0' && character <= '9') && character != '-' && character != '_' {
			return false
		}
	}
	return true
}

// validArchiveKey bounds an opaque provider-specific object key at the RPC boundary.
func validArchiveKey(key string) bool {
	return len(key) > 0 && len(key) <= 1024 && utf8.ValidString(key) && strings.IndexFunc(key, func(character rune) bool { return character < ' ' || character == 127 }) < 0
}
