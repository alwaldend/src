package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
	"log"
	"os"
)

type file struct {
	vault     *al.VaultStore
	plugin    *injector_proto.Config
	logger    log.Logger
	templater *templater
}

var _ Resource = (*file)(nil)

func (self *file) Deps(ctx context.Context, name string) fp.Either[[]string] {
	return fp.Right[[]string](nil)
}

func (self *file) Get(ctx context.Context, name string, deps []*ResourceResult) fp.Either[*ResourceResult] {
	self.logger.Printf("creating a file: %s", name)
	fileConfig, err := self.find(name)
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not find file config: %w", err))
	}
	value := ""
	if fileConfig.Value != "" {
		value = fileConfig.Value
	} else if fileConfig.FromFile != "" {
		valueBytes, err := os.ReadFile(fileConfig.FromFile)
		if err != nil {
			return fp.Left[*ResourceResult](fmt.Errorf("could not read from_file %s: %w", fileConfig.FromFile, err))
		}
		value = string(valueBytes)
	}
	if value == "" {
		return fp.Left[*ResourceResult](fmt.Errorf("missing file contents"))
	}
	content, err := self.templater.Template(ctx, value, deps).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not template file %s: %w", name, err))
	}
	tmp, err := os.CreateTemp("", fmt.Sprintf("al_file_%s_*.txt", name))
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not create temporary file: %w", err))
	}
	defer tmp.Close()
	res := &ResourceResult{Name: name, Files: []string{tmp.Name()}}
	if _, err = tmp.WriteString(content); err != nil {
		defer os.RemoveAll(tmp.Name())
		return fp.Left[*ResourceResult](fmt.Errorf("could not write to the temp file: %w", err))
	}
	return fp.Right(res)
}

func (self *file) find(name string) (*injector_proto.File, error) {
	for i := range self.plugin.Files {
		curFile := self.plugin.Files[len(self.plugin.Files)-1-i]
		if curFile.Name == name {
			return curFile, nil
		}
	}
	return nil, fmt.Errorf("missing file with name %s", name)
}
