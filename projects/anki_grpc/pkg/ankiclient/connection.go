package ankiclient

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// pipeConn adapts process pipes to the connection required by gRPC.
type pipeConn struct {
	// reader receives the worker's stdout.
	reader *os.File
	// writer sends requests to the worker's stdin.
	writer *os.File
	// inputOnce closes stdin once, independently of the response pipe.
	inputOnce sync.Once
	// inputError preserves the result of closing stdin.
	inputError error
	// closeOnce prevents duplicate transport closes by Go and gRPC.
	closeOnce sync.Once
	// closeError records shutdown failures for every caller.
	closeError error
}

// newPipeConn creates caller-owned transport endpoints and the child's pipe ends.
func newPipeConn() (*pipeConn, *os.File, *os.File, error) {
	input, writer, err := os.Pipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create Anki input pipe: %w", err)
	}
	reader, output, err := os.Pipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create Anki output pipe: %w", errors.Join(err, input.Close(), writer.Close()))
	}
	return &pipeConn{reader: reader, writer: writer}, input, output, nil
}

// Read receives HTTP/2 bytes from the worker.
func (connection *pipeConn) Read(data []byte) (int, error) {
	n, err := connection.reader.Read(data)
	if err != nil {
		return n, fmt.Errorf("read Anki worker: %w", err)
	}
	return n, nil
}

// Write sends HTTP/2 bytes to the worker.
func (connection *pipeConn) Write(data []byte) (int, error) {
	n, err := connection.writer.Write(data)
	if err != nil {
		return n, fmt.Errorf("write Anki worker: %w", err)
	}
	return n, nil
}

// Close releases both endpoints exactly once.
func (connection *pipeConn) Close() error {
	connection.closeOnce.Do(func() {
		if err := errors.Join(connection.closeInput(), connection.reader.Close()); err != nil {
			connection.closeError = fmt.Errorf("close Anki pipes: %w", err)
		}
	})
	return connection.closeError
}

// LocalAddr identifies the caller's pipe transport.
func (connection *pipeConn) LocalAddr() net.Addr { return &net.UnixAddr{Name: "caller", Net: "stdio"} }

// RemoteAddr identifies the worker's pipe transport.
func (connection *pipeConn) RemoteAddr() net.Addr {
	return &net.UnixAddr{Name: "anki-worker", Net: "stdio"}
}

// SetDeadline applies the requested deadline to both pipe endpoints.
func (connection *pipeConn) SetDeadline(deadline time.Time) error {
	if err := errors.Join(connection.SetReadDeadline(deadline), connection.SetWriteDeadline(deadline)); err != nil {
		return fmt.Errorf("set Anki pipe deadlines: %w", err)
	}
	return nil
}

// SetReadDeadline bounds pipe reads.
func (connection *pipeConn) SetReadDeadline(deadline time.Time) error {
	if err := connection.reader.SetReadDeadline(deadline); err != nil {
		return fmt.Errorf("set Anki read deadline: %w", err)
	}
	return nil
}

// SetWriteDeadline bounds pipe writes.
func (connection *pipeConn) SetWriteDeadline(deadline time.Time) error {
	if err := connection.writer.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("set Anki write deadline: %w", err)
	}
	return nil
}

// closeInput signals EOF while keeping stdout open for the worker's final frames.
func (connection *pipeConn) closeInput() error {
	connection.inputOnce.Do(func() {
		if err := connection.writer.Close(); err != nil {
			connection.inputError = fmt.Errorf("close Anki input: %w", err)
		}
	})
	return connection.inputError
}
