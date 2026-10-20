package native

/*
#include <stddef.h>
#include <stdint.h>
typedef struct AnkiBackend AnkiBackend;
typedef struct { uint8_t *data; size_t len; } AnkiBuffer;
AnkiBackend *anki_backend_new(AnkiBuffer *error);
int anki_backend_command(const AnkiBackend *backend, const uint8_t *path, size_t path_len,
    const uint8_t *input, size_t input_len, AnkiBuffer *output);
void anki_backend_free_buffer(AnkiBuffer buffer);
void anki_backend_free(AnkiBackend *backend);
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

// ErrUnknownCommand indicates a route absent from Anki's generated dispatch table.
var ErrUnknownCommand = errors.New("unknown Anki command")

// NativeError preserves Anki's serialized BackendError for the caller's decoder.
type NativeError struct {
	// Details contains the unchanged upstream protobuf error bytes.
	Details []byte
}

// Error identifies an unsuccessful native command without interpreting its payload.
func (failure *NativeError) Error() string { return "Anki command failed" }

// Server owns a native handle; callers must finish all Command calls before Close.
// The value must not be copied after initialization.
type Server struct {
	// handle points to Rust-owned synchronized backend state.
	handle *C.AnkiBackend
}

// NewServer initializes Anki with upstream defaults and transfers handle ownership to Go.
func NewServer() (*Server, error) {
	var failure C.AnkiBuffer
	handle := C.anki_backend_new(&failure)
	if handle == nil {
		message, err := copyBuffer(failure)
		if err != nil {
			return nil, fmt.Errorf("read Anki initialization error: %w", err)
		}
		return nil, fmt.Errorf("initialize Anki backend: %s", message)
	}
	return &Server{handle: handle}, nil
}

// Command borrows request bytes only for its synchronous native call and copies the owned response.
func (backend *Server) Command(method string, input []byte) ([]byte, error) {
	path := []byte(method)
	var output C.AnkiBuffer
	code := C.anki_backend_command(backend.handle,
		(*C.uint8_t)(unsafe.Pointer(unsafe.SliceData(path))), C.size_t(len(path)),
		(*C.uint8_t)(unsafe.Pointer(unsafe.SliceData(input))), C.size_t(len(input)), &output)
	switch code {
	case 0, 1:
		data, err := copyBuffer(output)
		if err != nil {
			return nil, fmt.Errorf("read Anki command response: %w", err)
		}
		if code == 1 {
			return nil, &NativeError{Details: data}
		}
		return data, nil
	case 2:
		return nil, fmt.Errorf("dispatch Anki %s: %w", method, ErrUnknownCommand)
	default:
		return nil, fmt.Errorf("dispatch Anki %s: invalid native arguments", method)
	}
}

// Close releases native state after the caller has drained every Command.
func (backend *Server) Close() {
	C.anki_backend_free(backend.handle)
	backend.handle = nil
}

// copyBuffer copies native results and always returns their allocations to Rust.
func copyBuffer(buffer C.AnkiBuffer) ([]byte, error) {
	defer C.anki_backend_free_buffer(buffer)
	if buffer.len > math.MaxInt32 {
		return nil, fmt.Errorf("native response exceeds the supported message size")
	}
	return C.GoBytes(unsafe.Pointer(buffer.data), C.int(buffer.len)), nil
}
