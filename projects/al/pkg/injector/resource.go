package injector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Resource is a provider-neutral named graph node.
type Resource struct {
	// Name identifies the node in templates and dependency references.
	Name string
	// Type selects a built-in or advertised provider.
	Type string
	// Args holds opaque type-specific settings.
	Args map[string]any
	// Data holds JSON-compatible input for local templates and providers.
	Data any
	// Deps preserves the declared dependency order.
	Deps []string
}

// ResourceResult carries values across graph edges and provider RPCs.
type ResourceResult struct {
	// Name identifies the originating resource.
	Name string
	// Env contains exports for dependent resources and the child.
	Env map[string]string
	// Files contains invocation-owned dependency paths.
	Files []string
	// Data contains opaque JSON-compatible values for templates and providers.
	Data any
}

// ResourceFetcher evaluates one node after its dependencies have completed.
type ResourceFetcher interface {
	Get(context.Context, *Resource, []*ResourceResult) (*ResourceResult, error)
}

// FetchFunc adapts an external provider to the built-in graph evaluator.
type FetchFunc func(context.Context, *Resource, []*ResourceResult) (*ResourceResult, error)

// Get forwards one ordered resource evaluation to the registered provider.
func (f FetchFunc) Get(ctx context.Context, r *Resource, d []*ResourceResult) (*ResourceResult, error) {
	return f(ctx, r, d)
}

// Engine owns the registry, graph and invocation-scoped built-in resources.
type Engine struct {
	// resources stores selected definitions by resource identity.
	resources map[string]*Resource
	// fetchers maps resource types to their registered implementation.
	fetchers map[string]ResourceFetcher
	// cleaner owns built-in temporary files until engine shutdown.
	cleaner *Cleaner
	// lc owns built-in resource processes in dependency evaluation order.
	lc lifecycle.Manager
}

// New creates an engine with only generic built-in resource types.
func New() *Engine {
	cleaner := NewCleaner()
	e := &Engine{resources: map[string]*Resource{}, fetchers: map[string]ResourceFetcher{}, cleaner: cleaner}
	templater := &templater{}
	e.fetchers["env"] = newEnvFetcher(templater)
	e.fetchers["file"] = newFileFetcher(templater, cleaner)
	e.fetchers["process"] = newProcessFetcher(&e.lc)
	return e
}

// decodeArgs validates built-in settings without including their values in errors.
func decodeArgs(args map[string]any, out proto.Message) error {
	data, err := json.Marshal(args)
	if err != nil {
		return errors.New("invalid resource settings")
	}
	if err := protojson.Unmarshal(data, out); err != nil {
		return errors.New("invalid resource settings")
	}
	return nil
}

// resourceDefinition retains a node's explicit group namespace for precedence.
type resourceDefinition struct {
	// group identifies the contributing root map entry.
	group string
	// config contains the provider-neutral node declaration.
	config *al_proto.ResourceConfig
}

// groupAncestors resolves explicit lineage without selecting ancestor resources.
func groupAncestors(selected, declarations map[string]*al_proto.InjectorCall) (map[string]map[string]bool, error) {
	ancestors := map[string]map[string]bool{}
	active := map[string]bool{}
	var visit func(string) error
	visit = func(name string) error {
		if active[name] {
			return fmt.Errorf("injector extension cycle at %q", name)
		}
		if _, done := ancestors[name]; done {
			return nil
		}
		group := declarations[name]
		if group == nil {
			return fmt.Errorf("missing injector extension %q", name)
		}
		active[name] = true
		lineage := map[string]bool{}
		if parent := group.Extends; parent != "" {
			if err := visit(parent); err != nil {
				return fmt.Errorf("resolve injector group %q: %w", name, err)
			}
			maps.Copy(lineage, ancestors[parent])
			lineage[parent] = true
		}
		active[name] = false
		ancestors[name] = lineage
		return nil
	}
	for _, name := range slices.Sorted(maps.Keys(selected)) {
		if err := visit(name); err != nil {
			return nil, fmt.Errorf("resolve selected injection: %w", err)
		}
	}
	return ancestors, nil
}

// Load resolves named groups by explicit ancestry and validates the selected graph.
func (e *Engine) Load(calls, declarations map[string]*al_proto.InjectorCall) error {
	ancestors, err := groupAncestors(calls, declarations)
	if err != nil {
		return fmt.Errorf("resolve injector extensions: %w", err)
	}
	definitions := map[string][]resourceDefinition{}
	for _, group := range slices.Sorted(maps.Keys(calls)) {
		for name, node := range calls[group].Res {
			if node == nil || name == "" || node.Type == "" {
				return fmt.Errorf("resource %q requires a type", name)
			}
			definitions[name] = append(definitions[name], resourceDefinition{group: group, config: node})
		}
	}
	e.resources = map[string]*Resource{}
	for _, name := range slices.Sorted(maps.Keys(definitions)) {
		maximal := []resourceDefinition{}
		for _, candidate := range definitions[name] {
			overridden := false
			for _, other := range definitions[name] {
				if ancestors[other.group][candidate.group] {
					overridden = true
					break
				}
			}
			if !overridden {
				maximal = append(maximal, candidate)
			}
		}
		node := maximal[0].config
		for _, candidate := range maximal[1:] {
			if !proto.Equal(node, candidate.config) {
				return fmt.Errorf("conflicting resource %q from injector groups %q and %q; declare an explicit extends relationship", name, maximal[0].group, candidate.group)
			}
		}
		value, err := al.FromPbJson(node.Args).Get()
		if err != nil {
			return fmt.Errorf("decode resource %q arguments: %w", name, err)
		}
		args, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("resource %q requires an args map", name)
		}
		data, err := al.FromPbJson(node.Data).Get()
		if err != nil {
			return fmt.Errorf("decode resource %q data: %w", name, err)
		}
		e.resources[name] = &Resource{Name: name, Type: node.Type, Args: args, Data: data, Deps: node.Deps}
	}
	state := map[string]int{}
	var visit func(string) error
	visit = func(name string) error {
		r, ok := e.resources[name]
		if !ok {
			return fmt.Errorf("missing dependency %q", name)
		}
		if state[name] == 1 {
			return fmt.Errorf("dependency cycle at resource %q", name)
		}
		if state[name] == 2 {
			return nil
		}
		state[name] = 1
		for _, dep := range r.Deps {
			if err := visit(dep); err != nil {
				return fmt.Errorf("validate resource %q: %w", name, err)
			}
		}
		state[name] = 2
		return nil
	}
	for _, name := range slices.Sorted(maps.Keys(e.resources)) {
		if err := visit(name); err != nil {
			return fmt.Errorf("validate injector graph: %w", err)
		}
	}
	return nil
}

// NeedsProviders reports whether any selected type needs an external provider.
func (e *Engine) NeedsProviders() bool {
	for _, r := range e.resources {
		if _, ok := e.fetchers[r.Type]; !ok {
			return true
		}
	}
	return false
}

// Register rejects capability ambiguity, including collisions with built-ins.
func (e *Engine) Register(kind string, fetcher ResourceFetcher) error {
	if kind == "" {
		return errors.New("empty advertised resource type")
	}
	if _, ok := e.fetchers[kind]; ok {
		return fmt.Errorf("duplicate resource type %q", kind)
	}
	e.fetchers[kind] = fetcher
	return nil
}

// Evaluate fetches every selected node once and combines dependency-first exports.
func (e *Engine) Evaluate(ctx context.Context) (map[string]string, error) {
	for _, r := range e.resources {
		if _, ok := e.fetchers[r.Type]; !ok {
			return nil, fmt.Errorf("no resource provider for type %q", r.Type)
		}
	}
	cache := map[string]*ResourceResult{}
	env := map[string]string{}
	var get func(string) (*ResourceResult, error)
	get = func(name string) (*ResourceResult, error) {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("evaluate resource %q: %w", name, err)
		}
		if result, ok := cache[name]; ok {
			return result, nil
		}
		r := e.resources[name]
		deps := []*ResourceResult{}
		for _, name := range r.Deps {
			result, err := get(name)
			if err != nil {
				return nil, fmt.Errorf("resolve dependency %q for %q: %w", name, r.Name, err)
			}
			deps = append(deps, result)
		}
		result, err := e.fetchers[r.Type].Get(ctx, r, deps)
		if err != nil {
			return nil, fmt.Errorf("fetch resource %q (%s): %w", name, r.Type, err)
		}
		if result == nil {
			return nil, fmt.Errorf("provider returned no result for %q", name)
		}
		if result.Name != "" && result.Name != name {
			return nil, fmt.Errorf("provider returned a different resource identity for %q", name)
		}
		result.Name = name
		cache[name] = result
		maps.Copy(env, result.Env)
		return result, nil
	}
	for _, name := range slices.Sorted(maps.Keys(e.resources)) {
		if _, err := get(name); err != nil {
			return nil, fmt.Errorf("evaluate injector graph: %w", err)
		}
	}
	return env, nil
}

// Stop releases processes before files, leaving external providers to the caller.
func (e *Engine) Stop(ctx context.Context) error {
	return errors.Join(e.lc.Stop(ctx), e.cleaner.Stop(ctx))
}

// ToProto converts a dependency result to the provider-neutral wire format.
func ToProto(r *ResourceResult) (*al_proto.ResourceResult, error) {
	data, err := wireData(r.Data)
	if err != nil {
		return nil, fmt.Errorf("encode resource %q: %w", r.Name, err)
	}
	return &al_proto.ResourceResult{Name: r.Name, Env: r.Env, Files: r.Files, Data: data}, nil
}

// FromProto decodes a provider result without exposing its contents in diagnostics.
func FromProto(r *al_proto.ResourceResult) (*ResourceResult, error) {
	if r == nil {
		return nil, errors.New("provider returned no resource result")
	}
	data, err := al.FromPbJson(r.Data).Get()
	if err != nil {
		return nil, errors.New("invalid provider resource data")
	}
	return &ResourceResult{Name: r.Name, Env: r.Env, Files: r.Files, Data: data}, nil
}

// wireData accepts JSON-compatible provider data without rounding integer values.
func wireData(value any) (*al_proto.Json, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, errors.New("invalid provider resource data")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, errors.New("invalid provider resource data")
	}
	var normalize func(any) (any, error)
	normalize = func(value any) (any, error) {
		switch v := value.(type) {
		case json.Number:
			if !strings.ContainsAny(string(v), ".eE") {
				number, err := v.Int64()
				if err != nil {
					return nil, errors.New("provider integer exceeds int64")
				}
				return number, nil
			}
			number, err := v.Float64()
			if err != nil {
				return nil, errors.New("invalid provider number")
			}
			return number, nil
		case map[string]any:
			for key, entry := range v {
				normalized, err := normalize(entry)
				if err != nil {
					return nil, err
				}
				v[key] = normalized
			}
			return v, nil
		case []any:
			for i, entry := range v {
				normalized, err := normalize(entry)
				if err != nil {
					return nil, err
				}
				v[i] = normalized
			}
			return v, nil
		default:
			return value, nil
		}
	}
	normalized, err := normalize(decoded)
	if err != nil {
		return nil, err
	}
	return al.ToPbJson(normalized).Get()
}
