//go:build js && wasm

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"syscall/js"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	backend "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/browser"
	"github.com/spf13/afero"
)

// browser owns one application, its shared virtual files and scoped backend lifetimes.
type browser struct {
	// application coordinates existing configuration and collection workflows.
	application *ankiascode.App
	// storage supplies workspace transfers in the application's path namespace.
	storage filesystem.FS
	// clients records cleanup completion before persistence can inspect the workspace.
	clients []<-chan struct{}
}

// main registers an asynchronous worker callback and keeps the Go runtime alive.
func main() {
	storage := afero.NewMemMapFs()
	if err := storage.MkdirAll("/tmp", 0o755); err != nil {
		panic(fmt.Errorf("create browser temporary directory: %w", err))
	}
	application := &browser{}
	application.useStorage(storage)
	callback := js.FuncOf(application.invoke)
	defer callback.Release()
	js.Global().Set("ankiascodeInvoke", callback)
	js.Global().Call("ankiascodeReady")
	select {}
}

// useStorage constructs the existing App with one shared Afero filesystem and lazy backend factory.
func (application *browser) useStorage(storage afero.Fs) {
	application.storage = filesystem.New(storage)
	application.application = ankiascode.New(ankiascode.Options{
		FS: storage,
		NewClient: ankigrpc.ToClientFactory(func(address string) (*backend.Client, error) {
			if address != "" {
				return nil, fmt.Errorf("browser collection operations require the local server")
			}
			client := backend.NewClient(storage, js.Global().Get("ankiascodeCreateServer"))
			application.clients = append(application.clients, client.Done())
			return client, nil
		}),
	})
}

// invoke returns immediately so asynchronous backend startup can use the worker event loop.
func (application *browser) invoke(_ js.Value, arguments []js.Value) any {
	executor := js.FuncOf(func(_ js.Value, callbacks []js.Value) any {
		resolve := callbacks[0]
		go func() {
			if len(arguments) != 2 || arguments[0].Type() != js.TypeString || arguments[1].Type() != js.TypeObject || arguments[1].IsNull() {
				resolve.Invoke(map[string]any{"error": "expected a method name and argument object"})
				return
			}
			value, err := application.dispatch(arguments[0].String(), arguments[1])
			if err != nil {
				resolve.Invoke(map[string]any{"error": err.Error()})
				return
			}
			resolve.Invoke(map[string]any{"value": value})
		}()
		return nil
	})
	defer executor.Release()
	return js.Global().Get("Promise").New(executor)
}

// dispatch adapts requests without duplicating the App's workflow behavior.
func (application *browser) dispatch(method string, arguments js.Value) (any, error) {
	if err := application.requireClosedClients(); err != nil {
		return nil, fmt.Errorf("wait for previous collection cleanup: %w", err)
	}
	switch method {
	case "configDump":
		input, err := pathArgument(arguments, "text")
		if err != nil {
			return nil, fmt.Errorf("read dump input: %w", err)
		}
		format, err := stringArgument(arguments, "format")
		if err != nil {
			return nil, fmt.Errorf("read dump format: %w", err)
		}
		data, err := application.application.ConfigDump(context.Background(), input, format)
		if err != nil {
			return nil, fmt.Errorf("browser configuration dump: %w", err)
		}
		return string(data), nil
	case "configExport":
		output, err := pathArgument(arguments, "output")
		if err != nil {
			return nil, fmt.Errorf("read export output: %w", err)
		}
		encoded, err := stringArgument(arguments, "configs")
		if err != nil {
			return nil, fmt.Errorf("read export configuration: %w", err)
		}
		configs := &ankiascodepb.CollectionConfigSet{}
		if err := codec.Decode([]byte(encoded), codec.JSON, configs); err != nil {
			return nil, fmt.Errorf("decode browser configuration: %w", err)
		}
		if err := application.application.ConfigExport(context.Background(), output, configs); err != nil {
			return nil, fmt.Errorf("browser configuration export: %w", err)
		}
		return nil, nil
	case "writeFile", "readFile", "listFiles", "remove":
		name, err := pathArgument(arguments, "path")
		if err != nil {
			return nil, fmt.Errorf("read workspace path: %w", err)
		}
		value, err := application.fileOperation(method, name, arguments)
		if err != nil {
			return nil, fmt.Errorf("workspace %s %q: %w", method, name, err)
		}
		return value, nil
	case "rename":
		source, err := pathArgument(arguments, "source")
		if err != nil {
			return nil, fmt.Errorf("read rename source: %w", err)
		}
		destination, err := pathArgument(arguments, "destination")
		if err != nil {
			return nil, fmt.Errorf("read rename destination: %w", err)
		}
		if err := application.renameFile(source, destination); err != nil {
			return nil, fmt.Errorf("rename workspace file %q to %q: %w", source, destination, err)
		}
		return nil, nil
	case "collectionExport", "collectionPlan", "collectionApply", "collectionBuild", "collectionStateUpdate":
		value, err := application.collectionOperation(method, arguments)
		if err != nil {
			return nil, fmt.Errorf("browser %s: %w", method, err)
		}
		return value, nil
	case "snapshot":
		value, err := application.snapshot()
		if err != nil {
			return nil, fmt.Errorf("snapshot browser workspace: %w", err)
		}
		return value, nil
	case "restoreSnapshot":
		if err := application.restoreSnapshot(arguments.Get("entries")); err != nil {
			return nil, fmt.Errorf("restore browser workspace: %w", err)
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported browser method %q", method)
	}
}

// collectionOperation adapts browser paths to the existing managed collection workflows.
func (application *browser) collectionOperation(method string, arguments js.Value) (any, error) {
	input := "collection"
	if method == "collectionBuild" {
		input = "base"
	}
	collection, err := packageArgument(arguments, input)
	if err != nil {
		return nil, fmt.Errorf("read collection source: %w", err)
	}
	stateKey, err := stringArgument(arguments, "stateKey")
	if err != nil {
		return nil, fmt.Errorf("read state namespace: %w", err)
	}
	options := ankiascodeapi.CollectionOptions{Path: collection, StateKey: stateKey}
	ctx := context.Background()
	switch method {
	case "collectionExport":
		output, err := pathArgument(arguments, "output")
		if err != nil {
			return nil, fmt.Errorf("read configuration destination: %w", err)
		}
		stateOutput, err := stringArgument(arguments, "stateOutput")
		if err != nil {
			return nil, fmt.Errorf("read bindings destination: %w", err)
		}
		if stateOutput != "" {
			stateOutput, err = pathArgument(arguments, "stateOutput")
			if err != nil {
				return nil, fmt.Errorf("validate bindings destination: %w", err)
			}
		}
		if err := application.application.ColExport(ctx, options, output, stateOutput); err != nil {
			return nil, fmt.Errorf("export collection configuration: %w", err)
		}
	case "collectionStateUpdate":
		name, err := pathArgument(arguments, "state")
		if err != nil {
			return nil, fmt.Errorf("read bindings path: %w", err)
		}
		data, err := application.storage.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("read collection bindings %q: %w", name, err)
		}
		state := &ankiascodepb.State{}
		if err := codec.Decode(data, codec.TOML, state); err != nil {
			return nil, fmt.Errorf("decode collection bindings: %w", err)
		}
		if err := application.application.ColStateUpdate(ctx, options, state); err != nil {
			return nil, fmt.Errorf("update collection bindings: %w", err)
		}
	default:
		text, err := pathArgument(arguments, "text")
		if err != nil {
			return nil, fmt.Errorf("read configuration source: %w", err)
		}
		switch method {
		case "collectionPlan":
			plan, err := application.application.ColPlan(ctx, options, text)
			if err != nil {
				return nil, fmt.Errorf("plan collection changes: %w", err)
			}
			data, err := codec.Encode(plan, codec.JSON)
			if err != nil {
				return nil, fmt.Errorf("encode collection plan: %w", err)
			}
			return string(data), nil
		case "collectionApply":
			if err := application.application.ColApply(ctx, options, text); err != nil {
				return nil, fmt.Errorf("apply collection changes: %w", err)
			}
		case "collectionBuild":
			output, err := packageArgument(arguments, "output")
			if err != nil {
				return nil, fmt.Errorf("read package destination: %w", err)
			}
			if err := application.application.ColBuild(ctx, options, text, output); err != nil {
				return nil, fmt.Errorf("build collection package: %w", err)
			}
		}
	}
	return nil, nil
}

// requireClosedClients prevents persistence while any scoped backend can still access storage.
func (application *browser) requireClosedClients() error {
	remaining := application.clients[:0]
	for _, done := range application.clients {
		select {
		case <-done:
		default:
			remaining = append(remaining, done)
		}
	}
	application.clients = remaining
	if len(remaining) != 0 {
		return fmt.Errorf("collection cleanup is still running; retry after it finishes")
	}
	return nil
}

// snapshot copies completed workspace files and directories, excluding private runtime storage.
func (application *browser) snapshot() ([]any, error) {
	entries := []any{}
	err := afero.Walk(application.storage.Afero(), "/", func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("inspect snapshot path %q: %w", name, err)
		}
		if name == "/tmp" {
			return filepath.SkipDir
		}
		entry := map[string]any{"path": name, "directory": info.IsDir()}
		if !info.IsDir() {
			data, err := application.storage.ReadFile(name)
			if err != nil {
				return fmt.Errorf("read snapshot path %q: %w", name, err)
			}
			value := js.Global().Get("Uint8Array").New(len(data))
			js.CopyBytesToJS(value, data)
			entry["data"] = value
		}
		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk browser workspace: %w", err)
	}
	return entries, nil
}

// restoreSnapshot validates and fills a fresh filesystem before replacing the application.
func (application *browser) restoreSnapshot(entries js.Value) error {
	if !js.Global().Get("Array").Call("isArray", entries).Bool() {
		return fmt.Errorf("snapshot entries must be an array")
	}
	storage := filesystem.New(afero.NewMemMapFs())
	seen := map[string]bool{}
	for index := 0; index < entries.Length(); index++ {
		entry := entries.Index(index)
		if entry.Type() != js.TypeObject || entry.IsNull() {
			return fmt.Errorf("snapshot entry %d must be an object", index)
		}
		name, err := pathArgument(entry, "path")
		if err != nil {
			return fmt.Errorf("read snapshot entry %d: %w", index, err)
		}
		if name != entry.Get("path").String() || seen[name] {
			return fmt.Errorf("snapshot path %q must be canonical and unique", name)
		}
		seen[name] = true
		directory := entry.Get("directory")
		if directory.Type() != js.TypeBoolean {
			return fmt.Errorf("snapshot path %q must specify its directory status", name)
		}
		info, err := storage.Stat(name)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect snapshot path %q: %w", name, err)
		}
		if err == nil && info.IsDir() != directory.Bool() {
			return fmt.Errorf("snapshot path %q has conflicting file and directory entries", name)
		}
		if directory.Bool() {
			if err := storage.MkdirAll(name, 0o755); err != nil {
				return fmt.Errorf("restore directory %q: %w", name, err)
			}
			continue
		}
		value := entry.Get("data")
		if value.Type() != js.TypeObject || value.IsNull() || !value.InstanceOf(js.Global().Get("Uint8Array")) {
			return fmt.Errorf("snapshot file %q must contain Uint8Array bytes", name)
		}
		data := make([]byte, value.Length())
		js.CopyBytesToGo(data, value)
		if err := storage.MkdirAll(path.Dir(name), 0o755); err != nil {
			return fmt.Errorf("restore file parent %q: %w", name, err)
		}
		if err := storage.WriteFile(name, data, 0o644); err != nil {
			return fmt.Errorf("restore file %q: %w", name, err)
		}
	}
	if err := storage.MkdirAll("/tmp", 0o755); err != nil {
		return fmt.Errorf("restore runtime temporary directory: %w", err)
	}
	application.useStorage(storage.Afero())
	return nil
}

// fileOperation transfers files through the same storage used by the App.
func (application *browser) fileOperation(method, name string, arguments js.Value) (any, error) {
	switch method {
	case "writeFile":
		value := arguments.Get("data")
		if value.Type() != js.TypeObject || value.IsNull() || !value.InstanceOf(js.Global().Get("Uint8Array")) {
			return nil, fmt.Errorf("file data must be a Uint8Array")
		}
		info, err := application.storage.Stat(name)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("inspect file destination: %w", err)
		}
		if err == nil && !info.Mode().IsRegular() {
			return nil, fmt.Errorf("file destination must be a regular file")
		}
		data := make([]byte, value.Length())
		js.CopyBytesToGo(data, value)
		if err := application.storage.MkdirAll(path.Dir(name), 0o755); err != nil {
			return nil, fmt.Errorf("create file parent: %w", err)
		}
		if err := application.storage.WriteFile(name, data, 0o644); err != nil {
			return nil, fmt.Errorf("write file bytes: %w", err)
		}
		return nil, nil
	case "readFile":
		data, err := application.storage.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("read file bytes: %w", err)
		}
		value := js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(value, data)
		return value, nil
	case "listFiles":
		names, err := application.listFiles(name)
		if err != nil {
			return nil, fmt.Errorf("list workspace files: %w", err)
		}
		slices.Sort(names)
		values := make([]any, len(names))
		for index, name := range names {
			values[index] = name
		}
		return values, nil
	case "remove":
		if err := application.storage.RemoveAll(name); err != nil {
			return nil, fmt.Errorf("remove workspace path: %w", err)
		}
		if err := application.storage.MkdirAll("/tmp", 0o755); err != nil {
			return nil, fmt.Errorf("restore runtime temporary directory: %w", err)
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported file operation %q", method)
	}
}

// renameFile publishes a regular file without allowing directory replacement or implicit parent creation.
func (application *browser) renameFile(source, destination string) error {
	info, err := application.storage.Stat(source)
	if err != nil {
		return fmt.Errorf("inspect source: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("source must be a regular file")
	}
	parent, err := application.storage.Stat(path.Dir(destination))
	if err != nil {
		return fmt.Errorf("inspect destination parent: %w", err)
	}
	if !parent.IsDir() {
		return fmt.Errorf("destination parent must be a directory")
	}
	target, err := application.storage.Stat(destination)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect destination: %w", err)
	}
	if err == nil && !target.Mode().IsRegular() {
		return fmt.Errorf("destination must be a regular file")
	}
	if source == destination {
		return nil
	}
	if err := application.storage.Rename(source, destination); err != nil {
		return fmt.Errorf("publish file: %w", err)
	}
	return nil
}

// listFiles gathers regular files recursively, reporting missing or unreadable paths.
func (application *browser) listFiles(name string) ([]string, error) {
	info, err := application.storage.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("inspect path %q: %w", name, err)
	}
	if !info.IsDir() {
		return []string{name}, nil
	}
	entries, err := application.storage.ReadDir(name)
	if err != nil {
		return nil, fmt.Errorf("list directory %q: %w", name, err)
	}
	var names []string
	for _, entry := range entries {
		children, err := application.listFiles(path.Join(name, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("list child %q: %w", entry.Name(), err)
		}
		names = append(names, children...)
	}
	return names, nil
}

// stringArgument validates a required string before using syscall/js conversions.
func stringArgument(arguments js.Value, name string) (string, error) {
	value := arguments.Get(name)
	if value.Type() != js.TypeString {
		return "", fmt.Errorf("%s must be a string", name)
	}
	return value.String(), nil
}

// pathArgument requires an absolute virtual path without host path expansion.
func pathArgument(arguments js.Value, name string) (string, error) {
	value, err := stringArgument(arguments, name)
	if err != nil {
		return "", fmt.Errorf("read path: %w", err)
	}
	if !path.IsAbs(value) || strings.ContainsRune(value, '\x00') {
		return "", fmt.Errorf("%s must be an absolute virtual path without NUL bytes", name)
	}
	value = path.Clean(value)
	if value == "/tmp" || strings.HasPrefix(value, "/tmp/") {
		return "", fmt.Errorf("/tmp is reserved for the collection runtime")
	}
	return value, nil
}

// packageArgument limits managed browser collections to portable completed packages.
func packageArgument(arguments js.Value, name string) (string, error) {
	value, err := pathArgument(arguments, name)
	if err != nil {
		return "", fmt.Errorf("read package path: %w", err)
	}
	if path.Ext(value) != ".colpkg" {
		return "", fmt.Errorf("%s must identify a .colpkg collection package", name)
	}
	return value, nil
}
