//go:build unix

package wasmhost

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// pollableFile gives protocol pumps a file whose Close interrupts inherited pipe I/O.
// The duplicate shares nonblocking flags with the original; Run owns the worker
// process and reserves these standard streams exclusively for its Rust command.
func pollableFile(file *os.File) (*os.File, error) {
	syscall.ForkLock.RLock()
	descriptor, err := syscall.Dup(int(file.Fd()))
	if err == nil {
		syscall.CloseOnExec(descriptor)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("duplicate protocol file %q: %w", file.Name(), err)
	}
	if err := syscall.SetNonblock(descriptor, true); err != nil {
		return nil, errors.Join(fmt.Errorf("make protocol file %q nonblocking: %w", file.Name(), err), syscall.Close(descriptor))
	}
	// NewFile must observe O_NONBLOCK to register inherited pipes with Go's poller.
	return os.NewFile(uintptr(descriptor), file.Name()), nil
}
