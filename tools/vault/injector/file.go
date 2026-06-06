package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
	"os"
)

type FileFetcher struct {
	vault     *al.VaultStore
	templater *Templater
	cleaner   *Cleaner
}

var _ ResourceFetcher = (*FileFetcher)(nil)

func (self *FileFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.FileFetcher"
}

func (self *FileFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) fp.Either[*ResourceResult] {
	value := ""
	f := r.GetFile()
	if r.GetFile() == nil {
		return fp.Left[*ResourceResult](fmt.Errorf("missing file config"))
	}
	if f.Value != "" {
		value = f.Value
	} else if f.FromFile != "" {
		valueBytes, err := os.ReadFile(f.FromFile)
		if err != nil {
			return fp.Left[*ResourceResult](fmt.Errorf("could not read from_file %s: %w", f.FromFile, err))
		}
		value = string(valueBytes)
	}
	if value == "" {
		return fp.Left[*ResourceResult](fmt.Errorf("missing file contents"))
	}
	content, err := self.templater.Template(ctx, value, d).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not template: %w", err))
	}
	tmp, err := os.CreateTemp("", fmt.Sprintf("al_file_%s_*.txt", r.Name))
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not create temporary file: %w", err))
	}
	defer tmp.Close()
	res := &ResourceResult{Name: r.Name, Files: []string{tmp.Name()}}
	if _, err = tmp.WriteString(content); err != nil {
		defer os.RemoveAll(tmp.Name())
		return fp.Left[*ResourceResult](fmt.Errorf("could not write to the temp file: %w", err))
	}
	self.cleaner.Add(tmp.Name())
	return fp.Right(res)
}
