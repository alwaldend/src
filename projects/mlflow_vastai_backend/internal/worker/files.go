package worker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// closeResource preserves the resource identity of a close failure.
func closeResource(resource io.Closer, name string) error {
	if err := resource.Close(); err != nil {
		return fmt.Errorf("close %s: %w", name, err)
	}
	return nil
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
