package worker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
)

// errLocked indicates that a live process holds an exclusive worker lock.
var errLocked = errors.New("already locked")

// createStateDirectory verifies that the durable state is private to its owner.
func createStateDirectory(fs afero.Fs, path string) error {
	if err := fs.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("create worker state directory: %w", err)
	}
	info, err := fs.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect worker state directory: %w", err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("worker state directory %q must be private to its owner", path)
	}
	return nil
}

// closeResource preserves the resource identity of a close failure.
func closeResource(resource io.Closer, name string) error {
	if err := resource.Close(); err != nil {
		return fmt.Errorf("close %s: %w", name, err)
	}
	return nil
}

// lockFile uses a kernel lock whose lifetime is the open file description.
func lockFile(fs afero.Fs, path string) (afero.File, error) {
	file, err := fs.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock %s: %w", path, err)
	}
	descriptor, ok := file.(interface{ Fd() uintptr })
	if !ok {
		return nil, errors.Join(errors.New("worker locks require a filesystem with file descriptors"), closeResource(file, path))
	}
	if err := unix.Flock(int(descriptor.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		if errors.Is(err, unix.EWOULDBLOCK) {
			err = errLocked
		}
		return nil, errors.Join(fmt.Errorf("lock %s: %w", path, err), closeResource(file, path))
	}
	return file, nil
}

// socketPath locates the rental-local gRPC endpoint in the worker state directory.
func socketPath(directory string) string {
	return filepath.Join(directory, "worker.sock")
}

// workerDirectory resolves an explicit override or the user's state directory.
func workerDirectory(directory string) (string, error) {
	if directory == "" {
		base := os.Getenv("XDG_STATE_HOME")
		if !filepath.IsAbs(base) {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("locate worker home directory: %w", err)
			}
			base = filepath.Join(home, ".local", "state")
		}
		directory = filepath.Join(base, "mlflow-vastai")
	}
	path, err := filepath.Abs(directory)
	if err != nil {
		return "", fmt.Errorf("resolve worker state directory: %w", err)
	}
	return path, nil
}
