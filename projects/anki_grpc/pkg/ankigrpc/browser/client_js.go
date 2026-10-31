//go:build js && wasm

package browser

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syscall/js"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/directclient"
	"google.golang.org/grpc"
)

// Client owns a browser Rust backend and its timezone callbacks.
// Use NewClient; each instance allows one Start attempt and must not be copied.
type Client struct {
	// connection shares command adaptation and lifecycle with the native backend.
	connection directclient.Client
	// factory resolves with a fresh backend's synchronous command and close methods.
	factory js.Value
}

// NewClient constructs an inert backend using the application's JavaScript factory.
// factory receives a timezone host and resolves with {command, close}.
// The factory captures the caller-owned filesystem and root in JavaScript.
// Start must run outside a synchronous JavaScript callback's event stack.
func NewClient(factory js.Value) *Client {
	return &Client{factory: factory}
}

// Start waits for the factory to finish, checking ctx before and after initialization.
// Cancellation during initialization waits for the factory and prevents publication.
// A stuck factory delays startup failure and cleanup. Call Stop after every attempt;
// Done observes complete cleanup even when a concurrent Stop stops waiting.
func (client *Client) Start(ctx context.Context) error {
	if err := client.connection.Start(ctx, client.newBackend); err != nil {
		return fmt.Errorf("start browser Anki client: %w", err)
	}
	return nil
}

// newBackend prepares timezone callbacks before invoking the application's factory.
// A returned backend also owns cleanup when the factory fails or returns invalid state.
func (client *Client) newBackend() (directclient.Backend, error) {
	if client.factory.Type() != js.TypeFunction {
		return nil, fmt.Errorf("expected a JavaScript server factory")
	}
	timezone, err := newTimezoneHost()
	if err != nil {
		return nil, fmt.Errorf("prepare browser Anki timezone: %w", err)
	}
	backend := &browserBackend{timezone: timezone}
	server, err := invokeClientJavaScript(func() js.Value {
		return client.factory.Invoke(timezone.endpoint)
	})
	if err == nil {
		server, err = awaitClientPromise(server)
	}
	if err != nil {
		return backend, fmt.Errorf("create browser Anki backend: %w", err)
	}
	backend.server = server
	valid, err := invokeClientJavaScript(func() js.Value {
		return js.ValueOf(server.Type() == js.TypeObject && !server.IsNull() &&
			server.Get("command").Type() == js.TypeFunction && server.Get("close").Type() == js.TypeFunction)
	})
	if err != nil {
		return backend, fmt.Errorf("inspect browser Anki backend: %w", err)
	}
	if !valid.Bool() {
		return backend, fmt.Errorf("server factory did not return command and close methods")
	}
	return backend, nil
}

// Connection lends the generated clients' connection while the backend is ready.
func (client *Client) Connection() (grpc.ClientConnInterface, error) {
	connection, err := client.connection.Connection()
	if err != nil {
		return nil, fmt.Errorf("get browser Anki connection: %w", err)
	}
	return connection, nil
}

// Stop drains commands, destroys Rust state, and closes filesystem descriptors.
// Close collections through Anki first; ctx bounds only the caller's cleanup wait.
func (client *Client) Stop(ctx context.Context) error {
	if err := client.connection.Stop(ctx); err != nil {
		return fmt.Errorf("stop browser Anki client: %w", err)
	}
	return nil
}

// Done observes final cleanup even after a Stop deadline or a trapped command.
func (client *Client) Done() <-chan struct{} { return client.connection.Done() }

// browserBackend adapts the synchronous JavaScript command API to owned Go bytes.
type browserBackend struct {
	// commandMutex prevents reactor reentry while synchronous host callbacks run.
	commandMutex sync.Mutex
	// server owns the Rust instance and its command and close methods.
	server js.Value
	// timezone owns callbacks until Rust has been destroyed.
	timezone *timezoneHost
}

// Command copies protobuf bytes across the Go and Rust memories through JavaScript.
// Calls are serialized because Go callbacks may run another command goroutine.
// A host exception or invalid response terminates the backend; Anki errors do not.
func (backend *browserBackend) Command(method string, request []byte) ([]byte, error) {
	backend.commandMutex.Lock()
	defer backend.commandMutex.Unlock()
	var response []byte
	var code int
	var protocolError error
	_, err := invokeClientJavaScript(func() js.Value {
		input := js.Global().Get("Uint8Array").New(len(request))
		js.CopyBytesToJS(input, request)
		result := backend.server.Call("command", method, input)
		if result.Type() != js.TypeObject || result.IsNull() || result.Get("code").Type() != js.TypeNumber {
			protocolError = fmt.Errorf("command did not return a response code")
			return js.Undefined()
		}
		status := result.Get("code").Float()
		if status != 0 && status != 1 && status != 2 {
			protocolError = fmt.Errorf("command returned an invalid response code %g", status)
			return js.Undefined()
		}
		code = int(status)
		data := result.Get("data")
		if data.Type() != js.TypeObject || data.IsNull() || !data.InstanceOf(js.Global().Get("Uint8Array")) {
			protocolError = fmt.Errorf("command did not return Uint8Array response bytes")
			return js.Undefined()
		}
		response = make([]byte, data.Length())
		js.CopyBytesToGo(response, data)
		return js.Undefined()
	})
	if failure := errors.Join(err, protocolError); failure != nil {
		return nil, fmt.Errorf("execute browser Anki %s: %w", method, errors.Join(directclient.ErrTerminated, failure))
	}
	switch code {
	case 1:
		return nil, &directclient.NativeError{Details: response}
	case 2:
		return nil, fmt.Errorf("execute browser Anki %s: %w", method, directclient.ErrUnknownCommand)
	default:
		return response, nil
	}
}

// Close releases the Rust backend before its owned timezone callbacks.
func (backend *browserBackend) Close() error {
	_, serverError := invokeClientJavaScript(func() js.Value {
		if backend.server.Type() == js.TypeObject && !backend.server.IsNull() && backend.server.Get("close").Type() == js.TypeFunction {
			return backend.server.Call("close")
		}
		return js.Undefined()
	})
	backend.timezone.Close()
	if serverError != nil {
		return fmt.Errorf("release browser Anki backend: %w", serverError)
	}
	return nil
}

// clientPromiseResult carries one JavaScript Promise settlement into a Go goroutine.
type clientPromiseResult struct {
	// value is the resolved JavaScript value when err is nil.
	value js.Value
	// err describes a JavaScript Promise rejection.
	err error
}

// awaitClientPromise joins a JavaScript Promise without blocking its callback stack.
func awaitClientPromise(value js.Value) (js.Value, error) {
	settled := make(chan clientPromiseResult, 1)
	resolve := js.FuncOf(func(_ js.Value, arguments []js.Value) any {
		result := js.Undefined()
		if len(arguments) != 0 {
			result = arguments[0]
		}
		settled <- clientPromiseResult{value: result}
		return nil
	})
	reject := js.FuncOf(func(_ js.Value, arguments []js.Value) any {
		message := "JavaScript Promise rejected without a reason"
		if len(arguments) != 0 {
			reason := arguments[0]
			message = reason.String()
			if reason.Type() == js.TypeObject && !reason.IsNull() && reason.Get("message").Type() == js.TypeString {
				message = reason.Get("message").String()
			}
		}
		settled <- clientPromiseResult{err: errors.New(message)}
		return nil
	})
	defer resolve.Release()
	defer reject.Release()
	if _, err := invokeClientJavaScript(func() js.Value {
		return js.Global().Get("Promise").Call("resolve", value).Call("then", resolve, reject)
	}); err != nil {
		return js.Undefined(), fmt.Errorf("observe browser Anki Promise: %w", err)
	}
	result := <-settled
	if result.err != nil {
		return js.Undefined(), fmt.Errorf("await browser Anki Promise: %w", result.err)
	}
	return result.value, nil
}

// invokeClientJavaScript translates thrown JavaScript errors at the host boundary.
func invokeClientJavaScript(operation func() js.Value) (value js.Value, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if failure, ok := recovered.(js.Error); ok {
				err = fmt.Errorf("call browser Anki host: %w", failure)
				return
			}
			panic(recovered)
		}
	}()
	return operation(), nil
}
