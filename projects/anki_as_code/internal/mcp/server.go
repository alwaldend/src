package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rpc"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/protobuf/proto"
)

// Server hosts MCP tools using the injected application API.
type Server struct {
	// application owns workflow behavior independently of protocol hosting.
	application ankiascodeapi.App
}

// New constructs an MCP host without starting protocol or native resources.
func New(application ankiascodeapi.App) *Server { return &Server{application: application} }

// server borrows one application and adapts tool envelopes to its operations.
type server struct {
	// context cancels application work when the hosting command is interrupted.
	context context.Context
	// application retains domain behavior and collection source reservations.
	application ankiascodeapi.App
	// protocol owns MCP framing, discovery, concurrency and request cancellation.
	protocol *mcp.Server
	// names prevents accidental replacement of an already registered tool.
	names map[string]bool
	// sessions lends explicit collection handles without changing path-based calls.
	sessions *collectionSessions
}

// workflowTool declares the envelope and protobuf contracts for one App workflow.
type workflowTool struct {
	// name is the stable command-derived MCP tool name.
	name string
	// description explains the operation's externally visible effect.
	description string
	// fields selects the envelope properties meaningful for this operation.
	fields []string
	// required identifies mandatory, nonempty envelope values.
	required []string
	// input supplies an optional protobuf-valued envelope property.
	input proto.Message
	// property names the optional protobuf input property.
	property string
	// output supplies the structured response contract, or nil for an empty object.
	output proto.Message
	// operation executes the declared workflow with its resolved collection source.
	operation func(context.Context, ankiascodeapi.CollectionOptions, arguments) ([]byte, error)
}

// arguments is the serialization contract shared by the explicitly scoped tool envelopes.
type arguments struct {
	// Handle selects an explicitly opened collection instead of path/address controls.
	Handle string `json:"handle"`
	// Collection identifies a managed collection database or archive.
	Collection string `json:"collection"`
	// Address selects a supported Anki client strategy or external server.
	Address string `json:"address"`
	// StateKey selects the collection binding namespace for workflows.
	StateKey string `json:"state_key"`
	// Text identifies the TOML configuration directory.
	Text string `json:"text"`
	// Output identifies an exported directory or built archive.
	Output string `json:"output"`
	// StateOutput optionally receives exported snapshot bindings as TOML.
	StateOutput string `json:"state_output"`
	// Base identifies the existing collection used for an independent archive build.
	Base string `json:"base"`
	// WriteBack optionally overrides native API archive publication, which defaults to true.
	WriteBack *bool `json:"write_back"`
	// Request retains native ProtoJSON without converting numbers through float64.
	Request json.RawMessage `json:"request"`
	// Configs retains the structured CollectionConfigSet for configuration export.
	Configs json.RawMessage `json:"configs"`
	// State retains the structured State object for update or replacement.
	State json.RawMessage `json:"state"`
	// Bindings maps configuration addresses to existing native identifier strings.
	Bindings map[string]string `json:"bindings"`
}

// Serve hosts one MCP connection on SDK-managed stdio until disconnect or cancellation.
// It registers the complete catalog before native work and waits for request cleanup on exit.
func (host *Server) Serve(ctx context.Context) (err error) {
	application := host.application
	adapter := &server{
		context: ctx, application: application, names: map[string]bool{},
		protocol: mcp.NewServer(&mcp.Implementation{Name: "anki_as_code", Version: "1"}, &mcp.ServerOptions{
			PageSize:     32,
			Capabilities: &mcp.ServerCapabilities{},
			Instructions: "Use col_open to reuse a collection handle, then col_close to publish archives and release its worker. Path-based calls open collections per call. Keep other Anki writers closed. Errors and cancellation do not imply rollback. Archive-backed reads can publish an unchanged archive. Mutating calls are not automatically retried.",
		}),
	}
	adapter.sessions = &collectionSessions{context: ctx, application: application, entries: map[string]*collectionHandle{}}
	defer func() {
		if closeErr := adapter.sessions.closeAll(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("finish MCP sessions: %w", closeErr))
		}
	}()
	if err := adapter.registerNative(); err != nil {
		return fmt.Errorf("register native MCP tools: %w", err)
	}
	if err := adapter.registerWorkflows(); err != nil {
		return fmt.Errorf("register workflow MCP tools: %w", err)
	}
	if err := adapter.sessions.register(adapter); err != nil {
		return fmt.Errorf("register collection session tools: %w", err)
	}
	if err := adapter.protocol.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("run stdio MCP server: %w", err)
	}
	return nil
}

// registerNative exposes declared RPCs once, omitting implicit forwarding aliases.
func (adapter *server) registerNative() error {
	services, err := rpc.Services()
	if err != nil {
		return fmt.Errorf("read compiled Anki catalog: %w", err)
	}
	for _, service := range services {
		for _, method := range service.Methods {
			if method.Descriptor.Parent().FullName() != service.FullName {
				continue
			}
			name := "api_" + strings.ReplaceAll(service.Name+"_"+method.Name, "-", "_")
			input := collectionToolSchema([]string{"collection", "address", "write_back", "request"}, []string{"request"}, true)
			if err := protobufProperty(input, "request", method.Descriptor.Input()); err != nil {
				return fmt.Errorf("describe %s input: %w", method.Path, err)
			}
			output, err := messageSchema(method.Descriptor.Output())
			if err != nil {
				return fmt.Errorf("describe %s output: %w", method.Path, err)
			}
			operation := func(ctx context.Context, value arguments) ([]byte, error) {
				if value.Collection == "" && value.Address == "" && value.Handle == "" {
					return nil, fmt.Errorf("native API requires collection, address or handle")
				}
				writeBack := true
				if value.WriteBack != nil {
					writeBack = *value.WriteBack
				}
				options, release, err := adapter.collectionOptions(ctx, value)
				if err != nil {
					return nil, fmt.Errorf("select native API source: %w", err)
				}
				defer release()
				response, err := adapter.application.APICall(ctx, ankiascodeapi.APIRequest{
					Collection: options, WriteBack: writeBack, Method: method.Path, Input: value.Request,
					InputFormat: "json", OutputFormat: "json",
				})
				if err != nil {
					return nil, fmt.Errorf("call native RPC %s: %w", method.Path, err)
				}
				return response, nil
			}
			if err := adapter.add(name, "Invoke "+method.Path+" using ProtoJSON. Collection lifetime and backend capability restrictions apply; errors do not imply rollback.", input, output, operation); err != nil {
				return fmt.Errorf("register %s: %w", method.Path, err)
			}
		}
	}
	return nil
}

// registerWorkflows exposes the application's ten configuration and collection operations.
func (adapter *server) registerWorkflows() error {
	for _, workflow := range []workflowTool{
		{"config_dump", "Read parsed TOML configurations without opening Anki.", []string{"text"}, []string{"text"}, nil, "", &ankiascodepb.CollectionConfigSet{}, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			data, err := adapter.application.ConfigDump(ctx, value.Text, "json")
			if err != nil {
				return nil, fmt.Errorf("dump config: %w", err)
			}
			return data, nil
		}},
		{"config_export", "Write a structured config set as a TOML tree; errors can leave partial output.", []string{"output", "configs"}, []string{"output", "configs"}, &ankiascodepb.CollectionConfigSet{}, "configs", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			configs := &ankiascodepb.CollectionConfigSet{}
			if err := codec.Decode(value.Configs, codec.JSON, configs); err != nil {
				return nil, fmt.Errorf("decode config set: %w", err)
			}
			if err := adapter.application.ConfigExport(ctx, value.Output, configs); err != nil {
				return nil, fmt.Errorf("export config: %w", err)
			}
			return []byte("{}"), nil
		}},
		{"col_export", "Export editable TOML and optional snapshot bindings. Archive cleanup can repackage the source.", []string{"collection", "address", "state_key", "output", "state_output"}, []string{"collection", "output"}, nil, "", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			if err := adapter.application.ColExport(ctx, options, value.Output, value.StateOutput); err != nil {
				return nil, fmt.Errorf("export collection: %w", err)
			}
			return []byte("{}"), nil
		}},
		{"col_plan", "Compare desired TOML with native state. Archive cleanup can repackage the source.", []string{"collection", "address", "state_key", "text"}, []string{"collection", "text"}, nil, "", &ankiascodepb.CollectionChangePlan{}, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			output, err := adapter.application.ColPlan(ctx, options, value.Text)
			if err != nil {
				return nil, fmt.Errorf("plan collection: %w", err)
			}
			data, err := codec.Encode(output, codec.JSON)
			if err != nil {
				return nil, fmt.Errorf("encode plan collection: %w", err)
			}
			return data, nil
		}},
		{"col_apply", "Reconcile a collection in place. Committed writes can persist after an error.", []string{"collection", "address", "state_key", "text"}, []string{"collection", "text"}, nil, "", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			if err := adapter.application.ColApply(ctx, options, value.Text); err != nil {
				return nil, fmt.Errorf("apply collection: %w", err)
			}
			return []byte("{}"), nil
		}},
		{"col_build", "Build an independent reconciled archive while preserving the base and TOML inputs.", []string{"base", "address", "state_key", "text", "output"}, []string{"base", "text", "output"}, nil, "", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			options.Path = value.Base
			if err := adapter.application.ColBuild(ctx, options, value.Text, value.Output); err != nil {
				return nil, fmt.Errorf("build collection: %w", err)
			}
			return []byte("{}"), nil
		}},
		{"col_state_dump", "Read stored bindings. Archive cleanup can repackage the source.", []string{"collection", "address", "state_key"}, []string{"collection"}, nil, "", &ankiascodepb.State{}, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			output, err := adapter.application.ColStateDump(ctx, options)
			if err != nil {
				return nil, fmt.Errorf("dump state: %w", err)
			}
			data, err := codec.Encode(output, codec.JSON)
			if err != nil {
				return nil, fmt.Errorf("encode dump state: %w", err)
			}
			return data, nil
		}},
		{"col_state_import", "Bind existing native identities to declared configuration paths.", []string{"collection", "address", "state_key", "text", "bindings"}, []string{"collection", "text", "bindings"}, nil, "", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			if err := adapter.application.ColStateImport(ctx, options, value.Text, value.Bindings); err != nil {
				return nil, fmt.Errorf("import state: %w", err)
			}
			return []byte("{}"), nil
		}},
		{"col_state_update", "Merge structured bindings into the selected state namespace.", []string{"collection", "address", "state_key", "state"}, []string{"collection", "state"}, &ankiascodepb.State{}, "state", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			state := &ankiascodepb.State{}
			if err := codec.Decode(value.State, codec.JSON, state); err != nil {
				return nil, fmt.Errorf("decode state: %w", err)
			}
			if err := adapter.application.ColStateUpdate(ctx, options, state); err != nil {
				return nil, fmt.Errorf("update state: %w", err)
			}
			return []byte("{}"), nil
		}},
		{"col_state_replace", "Replace all bindings in the selected state namespace.", []string{"collection", "address", "state_key", "state"}, []string{"collection", "state"}, &ankiascodepb.State{}, "state", nil, func(ctx context.Context, options ankiascodeapi.CollectionOptions, value arguments) ([]byte, error) {
			state := &ankiascodepb.State{}
			if err := codec.Decode(value.State, codec.JSON, state); err != nil {
				return nil, fmt.Errorf("decode state: %w", err)
			}
			if err := adapter.application.ColStateReplace(ctx, options, state); err != nil {
				return nil, fmt.Errorf("replace state: %w", err)
			}
			return []byte("{}"), nil
		}},
	} {
		input := toolSchema(workflow.fields, workflow.required)
		if _, managed := input["properties"].(map[string]any)["collection"]; managed {
			input = collectionToolSchema(workflow.fields, workflow.required, false)
		}
		if workflow.input != nil {
			if err := protobufProperty(input, workflow.property, workflow.input.ProtoReflect().Descriptor()); err != nil {
				return fmt.Errorf("describe %s input: %w", workflow.name, err)
			}
		}
		output := toolSchema(nil, nil)
		if workflow.output != nil {
			var err error
			output, err = messageSchema(workflow.output.ProtoReflect().Descriptor())
			if err != nil {
				return fmt.Errorf("describe %s output: %w", workflow.name, err)
			}
		}
		if err := adapter.add(workflow.name, workflow.description, input, output, func(ctx context.Context, value arguments) ([]byte, error) {
			options, release, err := adapter.collectionOptions(ctx, value)
			if err != nil {
				return nil, fmt.Errorf("select workflow source: %w", err)
			}
			defer release()
			data, err := workflow.operation(ctx, options, value)
			if err != nil {
				return nil, fmt.Errorf("call workflow %s: %w", workflow.name, err)
			}
			return data, nil
		}); err != nil {
			return fmt.Errorf("register %s: %w", workflow.name, err)
		}
	}
	return nil
}

// collectionOptions resolves a path or lends a handle until the returned release is called.
func (adapter *server) collectionOptions(ctx context.Context, value arguments) (ankiascodeapi.CollectionOptions, func(), error) {
	options := ankiascodeapi.CollectionOptions{Path: value.Collection, Address: value.Address, StateKey: value.StateKey}
	if value.Handle == "" {
		return options, func() {}, nil
	}
	connection, release, err := adapter.sessions.borrow(ctx, value.Handle)
	if err != nil {
		return options, nil, fmt.Errorf("select collection handle: %w", err)
	}
	return ankiascodeapi.CollectionOptions{Connection: connection, StateKey: value.StateKey}, release, nil
}

// add registers an envelope validator and preserves operation errors as tool errors.
func (adapter *server) add(name, description string, input, output map[string]any, operation func(context.Context, arguments) ([]byte, error)) error {
	if adapter.names[name] {
		return fmt.Errorf("duplicate MCP tool name %q", name)
	}
	adapter.names[name] = true
	adapter.protocol.AddTool(&mcp.Tool{Name: name, Description: description, InputSchema: input, OutputSchema: output}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, cancel := context.WithCancel(ctx)
		stop := context.AfterFunc(adapter.context, cancel)
		defer stop()
		defer cancel()
		value, err := decodeArguments(request.Params.Arguments, input)
		var data []byte
		if err == nil {
			data, err = operation(ctx, value)
		}
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%s: %v", name, err)}}}, nil
		}
		if !json.Valid(data) {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: name + ": operation returned invalid JSON"}}}, nil
		}
		return &mcp.CallToolResult{StructuredContent: json.RawMessage(data), Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil
	})
	return nil
}

// decodeArguments validates the envelope before App validates its protobuf payloads.
func decodeArguments(data json.RawMessage, schema map[string]any) (arguments, error) {
	var value arguments
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return value, fmt.Errorf("decode tool envelope: %w", err)
	}
	if fields == nil {
		return value, fmt.Errorf("tool arguments must be an object")
	}
	properties := schema["properties"].(map[string]any)
	for name, raw := range fields {
		if _, exists := properties[name]; !exists {
			return value, fmt.Errorf("unknown tool argument %q", name)
		}
		raw = bytes.TrimSpace(raw)
		if bytes.Equal(raw, []byte("null")) {
			return value, fmt.Errorf("tool argument %q cannot be null", name)
		}
		if name == "request" || name == "configs" || name == "state" || name == "bindings" {
			if len(raw) == 0 || raw[0] != '{' {
				return value, fmt.Errorf("tool argument %q must be an object", name)
			}
		}
	}
	for _, name := range schema["required"].([]string) {
		raw, exists := fields[name]
		if !exists || bytes.Equal(raw, []byte(`""`)) {
			return value, fmt.Errorf("tool argument %q is required and cannot be empty", name)
		}
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("decode typed tool arguments: %w", err)
	}
	if _, hasHandle := fields["handle"]; hasHandle {
		if value.Handle == "" {
			return value, fmt.Errorf("handle cannot be empty")
		}
		for _, name := range []string{"collection", "address", "write_back"} {
			if _, exists := fields[name]; exists {
				return value, fmt.Errorf("handle cannot accompany %q; select these controls when opening the collection", name)
			}
		}
	}
	return value, nil
}

// collectionToolSchema adds mutually exclusive handle and path/address source forms.
func collectionToolSchema(fields, required []string, addressOnly bool) map[string]any {
	var mandatory []string
	for _, name := range required {
		if name != "collection" {
			mandatory = append(mandatory, name)
		}
	}
	input := toolSchema(append(fields, "handle"), mandatory)
	paths := []any{map[string]any{"required": []string{"collection"}}}
	if addressOnly {
		paths = append(paths, map[string]any{"required": []string{"address"}})
	}
	input["oneOf"] = []any{
		map[string]any{"required": []string{"handle"}, "not": map[string]any{"anyOf": []any{map[string]any{"required": []string{"collection"}}, map[string]any{"required": []string{"address"}}, map[string]any{"required": []string{"write_back"}}}}},
		map[string]any{"anyOf": paths, "not": map[string]any{"required": []string{"handle"}}},
	}
	return input
}

// toolSchema selects only the envelope properties meaningful for a given tool.
func toolSchema(fields, required []string) map[string]any {
	properties := map[string]any{}
	for _, name := range fields {
		property := map[string]any{"type": "string"}
		switch name {
		case "handle":
			property["description"] = "Opaque handle returned by col_open; publish and release it with col_close."
		case "collection", "base":
			property["description"] = "Server-side path to an existing .anki2 database or modern .colpkg archive."
		case "address":
			property["description"] = "Optional Anki client address. With a managed path it selects the backend; address-only API access uses the backend's current collection."
		case "state_key":
			property["description"] = "Binding namespace; empty or omitted uses anki_as_code/state/v1."
		case "text":
			property["description"] = "Server-side TOML configuration directory."
		case "output", "state_output":
			property["description"] = "Server-side output path. state_output receives TOML snapshot bindings."
		case "write_back":
			property = map[string]any{"type": "boolean", "default": true, "description": "Publish managed archives during cleanup, including after errors. False does not undo database writes."}
		case "bindings":
			property = map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}, "description": "Configuration paths mapped to native identifier strings."}
		case "request", "configs", "state":
			property = map[string]any{"type": "object"}
		}
		properties[name] = property
	}
	if required == nil {
		required = []string{}
	}
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}
