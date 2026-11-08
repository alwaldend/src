package injector

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

// envFetcher renders environment exports from ordered dependency results.
type envFetcher struct {
	// templater renders only explicit dependency data.
	templater *templater
}

// newEnvFetcher binds dependency template rendering.
func newEnvFetcher(templater *templater) *envFetcher {
	return &envFetcher{templater}
}

var _ ResourceFetcher = (*envFetcher)(nil)

// Get evaluates one resource using the supplied dependency results.
func (self *envFetcher) Get(ctx context.Context, r *Resource, d []*ResourceResult) (*ResourceResult, error) {
	settings := &al_proto.InjectorEnv{}
	if err := decodeArgs(r.Args, settings); err != nil {
		return nil, fmt.Errorf("decode environment settings: %w", err)
	}
	content, err := self.templater.Template(ctx, settings.Value, d, r.Data).Get()
	if err != nil {
		return nil, fmt.Errorf("could not template: %w", err)
	}
	res := &ResourceResult{
		Name: r.Name,
		Env:  map[string]string{r.Name: content},
	}
	return res, nil
}
