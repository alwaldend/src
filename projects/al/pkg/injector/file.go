package injector

import (
	"context"
	"errors"
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

// fileFetcher owns private temporary files created from resource settings.
type fileFetcher struct {
	// templater renders only explicit dependency data.
	templater *templater
	// cleaner owns temporary paths until provider or engine shutdown.
	cleaner *Cleaner
}

// newFileFetcher binds rendering and private-file cleanup.
func newFileFetcher(templater *templater, cleaner *Cleaner) *fileFetcher {
	return &fileFetcher{templater, cleaner}
}

var _ ResourceFetcher = (*fileFetcher)(nil)

// Get evaluates one resource using the supplied dependency results.
func (self *fileFetcher) Get(ctx context.Context, r *Resource, d []*ResourceResult) (result *ResourceResult, resultErr error) {
	value := ""
	f := &al_proto.InjectorFile{}
	if err := decodeArgs(r.Args, f); err != nil {
		return nil, fmt.Errorf("decode file settings: %w", err)
	}
	if f.Value != "" {
		value = f.Value
	} else if f.FromFile != "" {
		valueBytes, err := os.ReadFile(f.FromFile)
		if err != nil {
			return nil, fmt.Errorf("could not read from_file %s: %w", f.FromFile, err)
		}
		value = string(valueBytes)
	}
	if value == "" {
		return nil, fmt.Errorf("missing file contents")
	}
	content, err := self.templater.Template(ctx, value, d, r.Data).Get()
	if err != nil {
		return nil, fmt.Errorf("could not template: %w", err)
	}
	tmp, err := os.CreateTemp("", fmt.Sprintf("al_file_%s_*.txt", r.Name))
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file: %w", err)
	}
	successful := false
	defer func() {
		if !successful {
			if err := tmp.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
				resultErr = errors.Join(resultErr, err)
			}
			resultErr = errors.Join(resultErr, os.RemoveAll(tmp.Name()))
		}
	}()
	res := &ResourceResult{
		Name:  r.Name,
		Files: []string{tmp.Name()},
	}
	if _, err = tmp.WriteString(content); err != nil {
		return nil, fmt.Errorf("could not write to the temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("could not close temporary file: %w", err)
	}
	if err := self.cleaner.Add(tmp.Name()); err != nil {
		return nil, fmt.Errorf("register temporary file: %w", err)
	}
	successful = true
	return res, nil
}
