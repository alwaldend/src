//go:build !unix

package wasmhost

import "os"

// pollableFile uses the platform's os.File implementation outside Unix.
func pollableFile(file *os.File) (*os.File, error) {
	return file, nil
}
