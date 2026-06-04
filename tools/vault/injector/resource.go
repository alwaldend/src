package main

import (
	"context"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

type ResourceType string

type ResourceResult struct {
	Name  string
	Env   map[string]string
	Files []string
	Data  any
}

type Resource interface {
	Get(ctx context.Context, name string, deps []*ResourceResult) fp.Either[*ResourceResult]
	Deps(ctx context.Context, name string) fp.Either[[]string]
}
