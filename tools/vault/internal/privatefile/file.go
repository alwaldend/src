package privatefile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CheckDestination rejects nonregular existing files before authentication.
func CheckDestination(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("inspect credential file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return errors.New("existing credential destination must be a regular file, not a symbolic link")
	}
	return nil
}

// Write atomically replaces the destination using a same-directory 0600 file.
func Write(path string, data []byte) (result error) {
	if err := CheckDestination(path); err != nil {
		return fmt.Errorf("recheck credential destination: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".vault-credential-*")
	if err != nil {
		return fmt.Errorf("create private credential file: %w", err)
	}
	closed, installed := false, false
	defer func() {
		if !closed {
			if err := file.Close(); err != nil {
				result = errors.Join(result, fmt.Errorf("close temporary credential file: %w", err))
			}
		}
		if !installed {
			if err := os.Remove(file.Name()); err != nil {
				result = errors.Join(result, fmt.Errorf("remove temporary credential file: %w", err))
			}
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return fmt.Errorf("set credential file permissions: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write private credential file: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("synchronize private credential file: %w", err)
	}
	err = file.Close()
	closed = true
	if err != nil {
		return fmt.Errorf("close private credential file: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace credential file: %w", err)
	}
	installed = true
	return nil
}
