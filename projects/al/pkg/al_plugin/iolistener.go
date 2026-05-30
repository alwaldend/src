package al_plugin

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Addr represents a network end point address.
//
// The two methods [Addr.Network] and [Addr.String] conventionally return strings that can be passed as the arguments to Dial, but the exact form and meaning of the strings is up to the implementation.
// https://pkg.go.dev/net#Addr
type IOAddr struct {
	reader io.Reader
	writer io.Writer
}

func NewIOAddr(reader io.Reader, writer io.Writer) *IOAddr {
	return &IOAddr{
		reader: reader,
		writer: writer,
	}
}

var _ net.Addr = (*IOAddr)(nil)

// Name of the network (for example, "tcp", "udp")
func (self *IOAddr) Network() string {
	return "io"
}

// String form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func (self *IOAddr) String() string {
	return fmt.Sprintf("%p.%p", self.reader, self.writer)
}

// Conn is a generic stream-oriented network connection.
//
// Multiple goroutines may invoke methods on a Conn simultaneously.
// https://pkg.go.dev/net#Conn
type IOConn struct {
	reader io.Reader
	writer io.Writer
	addr   *IOAddr
	closed bool
}

var _ net.Conn = (*IOConn)(nil)

func NewIOConn(reader io.Reader, writer io.Writer, addr *IOAddr) (*IOConn, error) {
	return &IOConn{
		reader: reader,
		writer: writer,
		addr:   addr,
		closed: false,
	}, nil
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (self *IOConn) Read(b []byte) (n int, err error) {
	if self.closed {
		return 0, fmt.Errorf("closed")
	}
	return self.reader.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (self *IOConn) Write(b []byte) (n int, err error) {
	if self.closed {
		return 0, fmt.Errorf("closed")
	}
	return self.writer.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
// Close may or may not block until any buffered data is sent;
// for TCP connections see [*TCPConn.SetLinger].
func (self *IOConn) Close() error {
	self.closed = true
	return nil
}

// LocalAddr returns the local network address, if known.
func (self *IOConn) LocalAddr() net.Addr {
	return self.addr
}

// RemoteAddr returns the remote network address, if known.
func (self IOConn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (self *IOConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (self *IOConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (self *IOConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// A Listener is a generic network listener for stream-oriented protocols.
//
// Multiple goroutines may invoke methods on a Listener simultaneously.
// https://pkg.go.dev/net#Listener
type IOListener struct {
	reader io.Reader
	writer io.Writer
	addr   *IOAddr
	conn   *IOConn
	closed chan any
	lock   sync.Mutex
}

var _ net.Listener = (*IOListener)(nil)

func NewIOListener(reader io.Reader, writer io.Writer) (*IOListener, error) {
	return &IOListener{
		reader: reader,
		writer: writer,
		addr:   &IOAddr{reader: reader, writer: writer},
		conn:   nil,
		closed: make(chan any, 1),
	}, nil
}

// Accept waits for and returns the next connection to the listener.
func (self *IOListener) Accept() (net.Conn, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.conn == nil {
		conn, err := NewIOConn(self.reader, self.writer, self.addr)
		if err != nil {
			return nil, fmt.Errorf("could not create an io connection: %w", err)
		}
		self.conn = conn
		return conn, nil
	}
	<-self.closed
	return nil, fmt.Errorf("closed")
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (self *IOListener) Close() error {
	err := self.conn.Close()
	close(self.closed)
	return err
}

// Addr returns the listener's network address.
func (self *IOListener) Addr() net.Addr {
	return self.addr
}
