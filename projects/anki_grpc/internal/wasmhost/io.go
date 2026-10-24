package wasmhost

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// hostIO supplies native byte streams and readiness to one Rust WASI command.
// The command goroutine owns handles and progress; pumps report through bounded channels.
type hostIO struct {
	// ctx bounds host readiness waits and command execution.
	ctx context.Context
	// stdin supplies the worker protocol's standard input.
	stdin io.ReadCloser
	// stdout receives the worker protocol's standard output.
	stdout io.WriteCloser
	// stderr receives transport diagnostics owned by the host.
	stderr io.Writer
	// cacheError is reported only when serving requires runtime resources.
	cacheError error
	// failure retains an opening or listener failure for the command boundary.
	failure error
	// closeError retains resource cleanup failures independently of execution errors.
	closeError error
	// notice coalesces readiness notifications without polling or unbounded queues.
	notice chan struct{}
	// handles owns all streams/listeners opened by this command.
	handles map[uint32]io.Closer
	// next assigns distinct handles for this command's lifetime.
	next uint32
}

// newHostIO creates inert host capabilities; help opens no protocol streams.
func newHostIO(ctx context.Context, stdin io.ReadCloser, stdout io.WriteCloser, stderr io.Writer, cacheError error) *hostIO {
	return &hostIO{ctx: ctx, stdin: stdin, stdout: stdout, stderr: stderr, cacheError: cacheError, notice: make(chan struct{}, 1), handles: make(map[uint32]io.Closer)}
}

// signal records readiness even when the Rust executor has not begun waiting.
func (host *hostIO) signal() {
	select {
	case host.notice <- struct{}{}:
	default:
	}
}

// add retains a newly opened resource under its command-local handle.
func (host *hostIO) add(resource io.Closer) int32 {
	host.next++
	host.handles[host.next] = resource
	return int32(host.next)
}

// fail keeps the underlying host error without exposing unchecked guest memory.
func (host *hostIO) fail(err error) int32 { host.failure = err; return -1 }

// instantiate binds platform streams, leaving CLI/protobuf/server behavior in Rust.
func (host *hostIO) instantiate(runtime wazero.Runtime) error {
	_, err := runtime.NewHostModuleBuilder("anki_io").
		NewFunctionBuilder().WithFunc(host.open).Export("open").
		NewFunctionBuilder().WithFunc(host.accept).Export("accept").
		NewFunctionBuilder().WithFunc(host.read).Export("read").
		NewFunctionBuilder().WithFunc(host.write).Export("write").
		NewFunctionBuilder().WithFunc(host.flush).Export("flush").
		NewFunctionBuilder().WithFunc(host.close).Export("close").
		NewFunctionBuilder().WithFunc(host.signal).Export("notify").
		NewFunctionBuilder().WithFunc(func(ctx context.Context) {
		select {
		case <-host.notice:
		case <-ctx.Done():
		}
	}).Export("wait").Instantiate(host.ctx)
	if err != nil {
		return fmt.Errorf("instantiate Rust stream imports: %w", err)
	}
	return nil
}

// open creates only the transport already validated by the Rust CLI.
func (host *hostIO) open(ctx context.Context, module api.Module, kind, pointer, length uint32) int32 {
	if host.cacheError != nil {
		return host.fail(host.cacheError)
	}
	if kind == 0 {
		reader, writer := host.stdin, host.stdout
		if input, ok := reader.(*os.File); ok {
			var err error
			reader, err = pollableFile(input)
			if err != nil {
				return host.fail(fmt.Errorf("prepare protocol stdin: %w", err))
			}
		}
		if output, ok := writer.(*os.File); ok {
			var err error
			writer, err = pollableFile(output)
			if err != nil {
				return host.fail(errors.Join(fmt.Errorf("prepare protocol stdout: %w", err), reader.Close()))
			}
		}
		return host.add(newHostStream(reader, writer, host.signal))
	}
	data, ok := module.Memory().Read(pointer, length)
	if !ok {
		return host.fail(errors.New("Rust listener address is outside guest memory"))
	}
	var network string
	switch kind {
	case 1:
		network = "tcp"
	case 2:
		network = "unix"
	default:
		return host.fail(fmt.Errorf("unsupported Rust transport %d", kind))
	}
	address := string(data)
	listener, err := (&net.ListenConfig{}).Listen(ctx, network, address)
	if err != nil {
		return host.fail(fmt.Errorf("listen on Anki %s endpoint %q: %w", network, address, err))
	}
	if unix, ok := listener.(*net.UnixListener); ok {
		unix.SetUnlinkOnClose(false)
	}
	resource := newHostListener(listener, host.signal)
	if _, err := fmt.Fprintf(host.stderr, "Listening on %s://%s\n", network, listener.Addr()); err != nil {
		return host.fail(errors.Join(fmt.Errorf("report Rust listener: %w", err), resource.Close()))
	}
	return host.add(resource)
}

// accept transfers one ready connection to Rust without blocking its executor.
func (host *hostIO) accept(handle uint32) int32 {
	listener, ok := host.handles[handle].(*hostListener)
	if !ok {
		return host.fail(fmt.Errorf("invalid Rust listener handle %d", handle))
	}
	select {
	case result := <-listener.accepted:
		if result.err != nil {
			return host.fail(fmt.Errorf("accept Rust connection: %w", result.err))
		}
		return host.add(newHostStream(result.connection, result.connection, host.signal))
	default:
		return -2
	}
}

// read copies ready bytes into a checked guest buffer, returning -2 until ready.
func (host *hostIO) read(_ context.Context, module api.Module, handle, pointer, length uint32) int32 {
	stream, ok := host.handles[handle].(*hostStream)
	if !ok {
		return -1
	}
	data, ok := module.Memory().Read(pointer, length)
	if !ok {
		return host.fail(errors.New("Rust read buffer is outside guest memory"))
	}
	if len(data) == 0 {
		return 0
	}
	if len(stream.buffer) == 0 && stream.readError == nil {
		select {
		case result := <-stream.reads:
			stream.buffer, stream.readError = result.data, result.err
		default:
			return -2
		}
	}
	if len(stream.buffer) != 0 {
		count := copy(data, stream.buffer)
		stream.buffer = stream.buffer[count:]
		return int32(count)
	}
	if errors.Is(stream.readError, io.EOF) {
		return 0
	}
	return -1
}

// write retains bounded caller-owned bytes until the host completes their write.
func (host *hostIO) write(_ context.Context, module api.Module, handle, pointer, length uint32) int32 {
	stream, ok := host.handles[handle].(*hostStream)
	if !ok {
		return -1
	}
	data, ok := module.Memory().Read(pointer, length)
	if !ok {
		return host.fail(errors.New("Rust write buffer is outside guest memory"))
	}
	if status := stream.writeStatus(); status != 0 {
		return status
	}
	stream.writes <- append([]byte(nil), data...)
	stream.writing = true
	return int32(len(data))
}

// flush distinguishes queued writes from bytes actually delivered to the peer.
func (host *hostIO) flush(handle uint32) int32 {
	stream, ok := host.handles[handle].(*hostStream)
	if !ok {
		return -1
	}
	return stream.writeStatus()
}

// close removes ownership before releasing a resource, so repeated guest closes
// and command cleanup never close an already released handle.
func (host *hostIO) close(handle uint32) {
	if resource := host.handles[handle]; resource != nil {
		delete(host.handles, handle)
		if err := resource.Close(); err != nil {
			host.closeError = errors.Join(host.closeError, fmt.Errorf("close Rust stream %d: %w", handle, err))
		}
		host.signal()
	}
}

// closeAll releases every retained socket and interrupts blocked host stream pumps.
func (host *hostIO) closeAll() error {
	for handle := range host.handles {
		host.close(handle)
	}
	if host.closeError != nil {
		return fmt.Errorf("release Rust host I/O: %w", host.closeError)
	}
	return nil
}

// streamRead transfers one bounded read or terminal error to the command thread.
type streamRead struct {
	// data owns the bytes until Rust consumes them.
	data []byte
	// err reports EOF or the underlying read failure.
	err error
}

// hostStream pumps a duplex stream without blocking Rust's single-thread executor.
type hostStream struct {
	// reader supplies incoming protocol bytes.
	reader io.ReadCloser
	// writer receives outgoing protocol bytes.
	writer io.WriteCloser
	// reads bounds unread host data to one queued chunk.
	reads chan streamRead
	// writes bounds outgoing data to one queued chunk.
	writes chan []byte
	// written transfers the one outstanding write's result to the command thread.
	written chan error
	// buffer retains the unconsumed portion of a read on the command thread.
	buffer []byte
	// readError retains a terminal read result on the command thread.
	readError error
	// writing admits no further data until the command thread consumes completion.
	writing bool
	// writeError retains the terminal write failure on the command thread.
	writeError error
	// pumps tracks both stream goroutines until closure unblocks their I/O.
	pumps sync.WaitGroup
	// done interrupts pumps when their Rust stream closes.
	done chan struct{}
}

// writeStatus consumes completion before another write can enter the bounded queue.
// Only the command thread accesses writing and writeError.
func (stream *hostStream) writeStatus() int32 {
	if stream.writing {
		select {
		case stream.writeError = <-stream.written:
			stream.writing = false
		default:
			return -2
		}
	}
	if stream.writeError != nil {
		return -1
	}
	return 0
}

// newHostStream starts bounded pumps and reports every readiness transition.
func newHostStream(reader io.ReadCloser, writer io.WriteCloser, signal func()) *hostStream {
	stream := &hostStream{reader: reader, writer: writer, reads: make(chan streamRead, 1), writes: make(chan []byte, 1), written: make(chan error, 1), done: make(chan struct{})}
	stream.pumps.Add(2)
	go func() {
		defer stream.pumps.Done()
		for {
			buffer := make([]byte, 64*1024)
			count, err := reader.Read(buffer)
			if count != 0 || err != nil {
				select {
				case stream.reads <- streamRead{data: buffer[:count], err: err}:
					signal()
				case <-stream.done:
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	go func() {
		defer stream.pumps.Done()
		for {
			select {
			case data := <-stream.writes:
				var err error
				for len(data) != 0 && err == nil {
					var count int
					count, err = writer.Write(data)
					data = data[count:]
					if count == 0 && err == nil {
						err = io.ErrNoProgress
					}
				}
				stream.written <- err
				signal()
				if err != nil {
					return
				}
			case <-stream.done:
				return
			}
		}
	}()
	return stream
}

// Close interrupts both directions and joins their pumps after handle removal.
func (stream *hostStream) Close() error {
	close(stream.done)
	var closeError error
	for _, resource := range []io.Closer{stream.reader, stream.writer} {
		if err := resource.Close(); err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, os.ErrClosed) {
			closeError = errors.Join(closeError, err)
		}
	}
	stream.pumps.Wait()
	if closeError != nil {
		return fmt.Errorf("close host protocol stream: %w", closeError)
	}
	return nil
}

// listenerAccept transfers one connection or terminal accept error to Rust.
type listenerAccept struct {
	// connection owns an accepted peer until Rust consumes it.
	connection net.Conn
	// err retains a terminal listener failure.
	err error
}

// hostListener bounds accepted connections awaiting admission to Rust.
type hostListener struct {
	// listener owns the native socket.
	listener net.Listener
	// accepted queues at most one accepted peer.
	accepted chan listenerAccept
	// stopped acknowledges that the accept pump can no longer enqueue peers.
	stopped chan struct{}
	// done interrupts a pump waiting to hand off a connection.
	done chan struct{}
}

// newHostListener accepts asynchronously without polling or unbounded buffering.
func newHostListener(listener net.Listener, signal func()) *hostListener {
	result := &hostListener{listener: listener, accepted: make(chan listenerAccept, 1), done: make(chan struct{}), stopped: make(chan struct{})}
	go func() {
		defer close(result.stopped)
		for {
			connection, err := listener.Accept()
			select {
			case result.accepted <- listenerAccept{connection: connection, err: err}:
				signal()
			case <-result.done:
				if connection != nil {
					_ = connection.Close()
				}
				return
			}
			if err != nil {
				return
			}
		}
	}()
	return result
}

// Close releases the listener and queued peer; the resource owner calls it once.
func (listener *hostListener) Close() error {
	close(listener.done)
	err := listener.listener.Close()
	<-listener.stopped
	select {
	case result := <-listener.accepted:
		if result.connection != nil {
			_ = result.connection.Close()
		}
	default:
	}
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("close host listener: %w", err)
	}
	return nil
}
