package directclient

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ErrUnknownCommand identifies routes absent from Anki's generated dispatch table.
var ErrUnknownCommand = errors.New("unknown Anki command")

// ErrTerminated marks an unusable backend whose active calls must drain before cleanup.
var ErrTerminated = errors.New("Anki backend terminated")

// NativeError retains Anki's encoded BackendError for the shared call decoder.
type NativeError struct {
	// Details contains the unchanged upstream protobuf error bytes.
	Details []byte
}

// Error identifies an unsuccessful native command without interpreting its payload.
func (failure *NativeError) Error() string { return "Anki command failed" }

// Backend owns command state until Close, which runs only after all calls finish.
type Backend interface {
	// Command borrows the request during execution and returns owned response bytes.
	Command(method string, request []byte) ([]byte, error)
	// Close destroys command state and releases its platform resources.
	Close() error
}

// Client owns a direct backend and implements the generated clients' connection.
// Its zero value is ready for one Start attempt; do not copy it after use.
type Client struct {
	// mutex serializes lifecycle transitions and active-work registration.
	mutex sync.Mutex
	// backend is available after startup until cleanup takes ownership of it.
	backend Backend
	// started records the one permitted startup attempt.
	started bool
	// stopped rejects registration after terminal shutdown begins.
	stopped bool
	// active tracks backend creation and commands beyond caller cancellation.
	active sync.WaitGroup
	// done closes after shutdown has drained work and destroyed backend state.
	done chan struct{}
	// exitError retains terminal command and cleanup failures for repeated Stop calls.
	exitError error
}

// Start initializes one backend; cancellation bounds waiting, not backend work.
// Call Stop after every attempt so Done can observe any late startup cleanup.
func (client *Client) Start(ctx context.Context, factory func() (Backend, error)) error {
	client.mutex.Lock()
	if client.started || client.stopped {
		client.mutex.Unlock()
		return fmt.Errorf("start direct Anki backend: already started or stopped")
	}
	client.started = true
	client.active.Add(1)
	client.mutex.Unlock()
	completed := make(chan error, 1)
	go func() {
		defer client.active.Done()
		completed <- client.initialize(ctx, factory)
	}()
	select {
	case err := <-completed:
		if err != nil {
			return fmt.Errorf("start direct Anki backend: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for direct Anki startup: %w", ctx.Err())
	}
}

// initialize publishes a new backend only while startup and the client remain live.
func (client *Client) initialize(ctx context.Context, factory func() (Backend, error)) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("check direct Anki startup: %w", err)
	}
	if factory == nil {
		return fmt.Errorf("direct Anki backend factory is required")
	}
	backend, err := factory()
	if err != nil {
		if backend != nil {
			err = errors.Join(err, client.closeBackend(backend))
		}
		return fmt.Errorf("create direct Anki backend: %w", err)
	}
	if backend == nil {
		return fmt.Errorf("direct Anki backend factory returned no backend")
	}
	failure := ctx.Err()
	client.mutex.Lock()
	if client.stopped {
		failure = errors.Join(failure, fmt.Errorf("client stopped during initialization"))
	}
	if failure == nil {
		client.backend = backend
	}
	client.mutex.Unlock()
	if failure != nil {
		return fmt.Errorf("finish direct Anki startup: %w", errors.Join(failure, client.closeBackend(backend)))
	}
	return nil
}

// Connection lends the client's connection after startup and before shutdown.
func (client *Client) Connection() (grpc.ClientConnInterface, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.backend == nil || client.stopped {
		return nil, fmt.Errorf("get direct Anki connection: client is not ready")
	}
	return client, nil
}

// Stop rejects new calls and waits for final cleanup within the supplied deadline.
// Cleanup continues after a deadline, and repeated calls observe the same result.
func (client *Client) Stop(ctx context.Context) error {
	client.mutex.Lock()
	client.stopLocked()
	done := client.done
	client.mutex.Unlock()
	select {
	case <-done:
	default:
		select {
		case <-done:
		case <-ctx.Done():
			return fmt.Errorf("wait for direct Anki shutdown: %w", ctx.Err())
		}
	}
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.exitError != nil {
		return fmt.Errorf("stop direct Anki backend: %w", client.exitError)
	}
	return nil
}

// stopLocked begins terminal cleanup once while mutex is held.
func (client *Client) stopLocked() {
	if client.stopped {
		return
	}
	client.stopped = true
	if client.done == nil {
		client.done = make(chan struct{})
	}
	go func() {
		client.active.Wait()
		client.mutex.Lock()
		backend := client.backend
		client.backend = nil
		client.mutex.Unlock()
		if backend != nil {
			_ = client.closeBackend(backend)
		}
		close(client.done)
	}()
}

// closeBackend records destruction failures before active work or shutdown completes.
func (client *Client) closeBackend(backend Backend) error {
	if err := backend.Close(); err != nil {
		failure := fmt.Errorf("release direct Anki backend: %w", err)
		client.mutex.Lock()
		client.exitError = errors.Join(client.exitError, failure)
		client.mutex.Unlock()
		return failure
	}
	return nil
}

// Done returns the stable channel observing final cleanup after Stop or termination.
func (client *Client) Done() <-chan struct{} {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.done == nil {
		client.done = make(chan struct{})
	}
	return client.done
}

// Invoke marshals unary calls and preserves cancellation, options and native errors.
// Cancelled callers never have their reply mutated by later backend completion.
func (client *Client) Invoke(ctx context.Context, method string, input, output any, options ...grpc.CallOption) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("call direct Anki %s: %w", method, err)
	}
	settings, err := directOptions(options)
	if err != nil {
		return fmt.Errorf("configure direct Anki %s: %w", method, err)
	}
	request, ok := input.(proto.Message)
	if !ok {
		return fmt.Errorf("call direct Anki %s: request must be a protobuf message", method)
	}
	reply, ok := output.(proto.Message)
	if !ok || !reply.ProtoReflect().IsValid() {
		return fmt.Errorf("call direct Anki %s: reply must be a writable protobuf message", method)
	}
	data, err := proto.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode direct Anki %s: %w", method, err)
	}
	if len(data) > settings.maxSend {
		return status.Errorf(codes.ResourceExhausted, "direct Anki request exceeds %d bytes", settings.maxSend)
	}
	client.mutex.Lock()
	if client.backend == nil || client.stopped {
		client.mutex.Unlock()
		return fmt.Errorf("call direct Anki %s: backend is not ready", method)
	}
	backend := client.backend
	client.active.Add(1)
	client.mutex.Unlock()
	completed := make(chan directResult, 1)
	go func() {
		defer client.active.Done()
		response, err := backend.Command(method, data)
		result := directResult{data: response, err: err}
		var failure *NativeError
		switch {
		case errors.As(err, &failure):
			result.trailers = metadata.Pairs(ankiclient.ErrorTrailer, string(failure.Details))
			result.err = status.Error(codes.Unknown, "Anki command failed")
		case errors.Is(err, ErrUnknownCommand):
			result.err = status.Errorf(codes.Unimplemented, "unknown Anki command: %s", method)
		case errors.Is(err, ErrTerminated):
			client.mutex.Lock()
			client.exitError = errors.Join(client.exitError, err)
			client.stopLocked()
			client.mutex.Unlock()
		}
		completed <- result
	}()
	select {
	case <-ctx.Done():
		return fmt.Errorf("call direct Anki %s: %w", method, ctx.Err())
	case result := <-completed:
		for _, destination := range settings.headers {
			*destination = metadata.MD{}
		}
		for _, destination := range settings.trailers {
			*destination = result.trailers.Copy()
		}
		if err := ankiclient.DecodeCallError(ctx, method, result.err, result.trailers); err != nil {
			return fmt.Errorf("execute direct Anki: %w", err)
		}
		if len(result.data) > settings.maxReceive {
			return status.Errorf(codes.ResourceExhausted, "direct Anki response exceeds %d bytes", settings.maxReceive)
		}
		if err := proto.Unmarshal(result.data, reply); err != nil {
			return fmt.Errorf("decode direct Anki %s: %w", method, err)
		}
		return nil
	}
}

// NewStream rejects streaming because the upstream command protocol is unary.
func (client *Client) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "direct Anki supports only unary commands")
}

// directResult holds backend completion without touching caller-owned reply state.
type directResult struct {
	// data contains copied response protobuf bytes.
	data []byte
	// trailers contains serialized native failure metadata.
	trailers metadata.MD
	// err distinguishes native errors from unknown routes or bridge failures.
	err error
}

// directCallOptions records supported options without emulating network behavior.
type directCallOptions struct {
	// headers retains every requested response-header destination.
	headers []*metadata.MD
	// trailers retains every requested response-trailer destination.
	trailers []*metadata.MD
	// maxSend bounds encoded request bytes.
	maxSend int
	// maxReceive bounds encoded response bytes.
	maxReceive int
}

// directOptions validates options before any backend command is dispatched.
func directOptions(options []grpc.CallOption) (directCallOptions, error) {
	settings := directCallOptions{maxSend: math.MaxInt32, maxReceive: math.MaxInt32}
	for _, option := range options {
		switch option := option.(type) {
		case grpc.StaticMethodCallOption, grpc.EmptyCallOption:
		case grpc.HeaderCallOption:
			if option.HeaderAddr == nil {
				return settings, fmt.Errorf("direct Anki header destination is nil")
			}
			settings.headers = append(settings.headers, option.HeaderAddr)
		case grpc.TrailerCallOption:
			if option.TrailerAddr == nil {
				return settings, fmt.Errorf("direct Anki trailer destination is nil")
			}
			settings.trailers = append(settings.trailers, option.TrailerAddr)
		case grpc.MaxSendMsgSizeCallOption:
			settings.maxSend = option.MaxSendMsgSize
		case grpc.MaxRecvMsgSizeCallOption:
			settings.maxReceive = option.MaxRecvMsgSize
		default:
			return settings, status.Errorf(codes.Unimplemented, "unsupported direct Anki call option %T", option)
		}
	}
	return settings, nil
}
