package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var packageIDPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._/-][a-z0-9]+)*$`)

func validatePackageID(id string) error {
	if len(id) > 200 || !packageIDPattern.MatchString(id) ||
		strings.Contains(id, "..") {
		return fmt.Errorf("malformed package id %q", id)
	}
	return nil
}

func validState(state PackageState) bool {
	switch state {
	case PackageLoading, PackageReady, PackageDegraded, PackageFailed,
		PackageTimeout, PackageDraining, PackageDisabled:
		return true
	default:
		return false
	}
}

func jsonMarshal(value any) ([]byte, error) {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode control state: %w", err)
	}
	return append(content, '\n'), nil
}

func strictDecode(content []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("control state contains multiple JSON values")
	}
	return nil
}

func atomicWriteFile(path string, content []byte) error {
	directory := filepath.Dir(path)
	temporary := filepath.Join(
		directory,
		fmt.Sprintf(".%s.%d.tmp", filepath.Base(path), time.Now().UnixNano()),
	)
	if err := os.WriteFile(temporary, content, 0o644); err != nil {
		return err
	}
	if err := os.Rename(temporary, path); err != nil {
		os.Remove(temporary)
		return err
	}
	return nil
}
