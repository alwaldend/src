//go:build cgo && linux && amd64

package native

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

// connection adapts generated unary calls to Rust state owned by its Client.
// The client starts and stops this connection; native crashes share application lifetime.
type connection struct {
	// mutex serializes lifecycle transitions and native-call registration.
	mutex sync.Mutex
	// native owns Rust state while the backend is ready or draining.
	native *Server
	// started records the one permitted startup attempt.
	started bool
	// stopped prevents native-call registration after shutdown starts.
	stopped bool
	// active tracks startup and native calls beyond cancellation of their waiting callers.
	active sync.WaitGroup
	// done closes after all native work and final backend destruction.
	done chan struct{}
}

// start checks cancellation before and after synchronous native initialization.
// Initialization cannot be interrupted; ctx does not own the native lifetime.
func (backend *connection) start(ctx context.Context) error {
	backend.mutex.Lock()
	if backend.started || backend.stopped {
		backend.mutex.Unlock()
		return fmt.Errorf("start direct Anki backend: already started or stopped")
	}
	backend.started = true
	backend.active.Add(1)
	backend.mutex.Unlock()
	defer backend.active.Done()
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start direct Anki backend: %w", err)
	}
	native, err := NewServer()
	if err != nil {
		return fmt.Errorf("start direct Anki backend: %w", err)
	}
	if err := ctx.Err(); err != nil {
		native.Close()
		return fmt.Errorf("finish direct Anki startup: %w", err)
	}
	backend.mutex.Lock()
	if backend.stopped {
		backend.mutex.Unlock()
		native.Close()
		return fmt.Errorf("finish direct Anki startup: client stopped during initialization")
	}
	backend.native = native
	backend.mutex.Unlock()
	return nil
}

// stop rejects new calls and releases native state after active calls finish.
// A cancelled ctx bounds waiting; cleanup continues, and another Stop can await it.
// Close an open collection through its generated client before stopping the backend.
func (backend *connection) stop(ctx context.Context) error {
	backend.mutex.Lock()
	if !backend.stopped {
		backend.stopped = true
		if backend.done == nil {
			backend.done = make(chan struct{})
		}
		native := backend.native
		go func() {
			backend.active.Wait()
			if native != nil {
				native.Close()
			}
			close(backend.done)
		}()
	}
	done := backend.done
	backend.mutex.Unlock()
	select {
	case <-done:
		return nil
	default:
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for direct Anki shutdown: %w", ctx.Err())
	}
}

// Invoke marshals a unary request, executes native dispatch and reuses the shared error decoder.
// Supported call options are static-method markers, headers, trailers and message-size limits.
// Cancellation stops waiting without interrupting native work or mutating the reply afterward.
func (backend *connection) Invoke(ctx context.Context, method string, input, output any, options ...grpc.CallOption) error {
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
	backend.mutex.Lock()
	if backend.native == nil || backend.stopped {
		backend.mutex.Unlock()
		return fmt.Errorf("call direct Anki %s: backend is not ready", method)
	}
	native := backend.native
	backend.active.Add(1)
	backend.mutex.Unlock()
	completed := make(chan directResult, 1)
	go func() {
		defer backend.active.Done()
		response, err := native.Command(method, data)
		result := directResult{data: response, err: err}
		var failure *NativeError
		switch {
		case errors.As(err, &failure):
			result.trailers = metadata.Pairs(ankiclient.ErrorTrailer, string(failure.Details))
			result.err = status.Error(codes.Unknown, "Anki command failed")
		case errors.Is(err, ErrUnknownCommand):
			result.err = status.Errorf(codes.Unimplemented, "unknown Anki command: %s", method)
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

// NewStream rejects streaming because the upstream Anki command protocol is unary.
func (backend *connection) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "direct Anki supports only unary commands")
}

// directResult holds a completed native response without writing caller-owned reply state.
type directResult struct {
	// data contains the copied response protobuf.
	data []byte
	// trailers contains serialized native failure metadata.
	trailers metadata.MD
	// err distinguishes native errors from unknown routes or bridge failures.
	err error
}

// directCallOptions records supported options without emulating network transport behavior.
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

// directOptions validates options before any native command is dispatched.
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
