package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Commands other than these initialization-free operations may load provider
// schemas or provider implementations. Defaulting to verification also covers
// new Terraform commands without maintaining a growing command allowlist.
func commandUsesProviders(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "init", "fmt", "version", "-v", "-version", "--version", "help", "-help", "--help":
		return false
	}
	return true
}

func (call invocation) verifyProviderSelections(config configuration, chdir string, environ []string, execute func([]string, io.Writer) error) error {
	// version reads lock selections without loading provider plugins or the
	// state backend. The invocation's CLI configuration disables checkpoints.
	var output bytes.Buffer
	if err := execute([]string{"version", "-json"}, &output); err != nil {
		return fmt.Errorf("inspect installed provider selections: %w", err)
	}
	var version struct {
		ProviderSelections map[string]string `json:"provider_selections"`
	}
	if err := json.Unmarshal(output.Bytes(), &version); err != nil {
		return fmt.Errorf("decode Terraform provider selections: %w", err)
	}
	if version.ProviderSelections == nil {
		return errors.New("Terraform version output did not contain provider_selections")
	}
	declared := make(map[string]providerPackage, len(config.Providers))
	for _, provider := range config.Providers {
		declared[provider.Source] = provider
	}
	selected := make([]string, 0, len(version.ProviderSelections))
	for source := range version.ProviderSelections {
		selected = append(selected, source)
	}
	sort.Strings(selected)
	for _, source := range selected {
		provider, found := declared[source]
		if !found {
			return fmt.Errorf("installed provider selection %s is not declared by this Bazel target", source)
		}
		selectedVersion := version.ProviderSelections[source]
		if selectedVersion != provider.Version {
			return fmt.Errorf("installed provider %s selects %s, but this Bazel target declares %s; initialize with the declared provider version", source, selectedVersion, provider.Version)
		}
		archive, err := call.runfiles.Rlocation(provider.Archive)
		if err != nil {
			return fmt.Errorf("locate declared provider archive %s: %w", source, err)
		}
		directory := filepath.Join(terraformDataDirectory(chdir, environ), "providers", filepath.FromSlash(source), provider.Version, provider.Platform)
		if err := verifyInstalledPackage(archive, directory); err != nil {
			return fmt.Errorf("installed provider %s %s does not match its declared Bazel archive: %w", source, provider.Version, err)
		}
	}
	return nil
}

func terraformDataDirectory(chdir string, environ []string) string {
	directory := environmentValue(environ, "TF_DATA_DIR")
	if directory == "" {
		return filepath.Join(chdir, ".terraform")
	}
	if filepath.IsAbs(directory) {
		return directory
	}
	// Terraform handles -chdir before constructing its working-directory
	// object, so relative TF_DATA_DIR paths are relative to that directory.
	return filepath.Join(chdir, directory)
}

func verifyInstalledPackage(archive, directory string) error {
	packaged, err := zip.OpenReader(archive)
	if err != nil {
		return fmt.Errorf("open provider ZIP: %w", err)
	}
	defer packaged.Close()
	// Terraform permits its package directory to be a symlink to an unpacked
	// mirror or shared cache. Verify the resolved tree rather than treating
	// that root link as a package file or modifying it.
	directory, err = filepath.EvalSymlinks(directory)
	if err != nil {
		return fmt.Errorf("locate installed package: %w", err)
	}
	files := make(map[string][sha256.Size]byte)
	for _, file := range packaged.File {
		name := strings.TrimSuffix(file.Name, "/")
		if !validRunfilePath(name) {
			return fmt.Errorf("provider ZIP contains an invalid path %q", file.Name)
		}
		if file.FileInfo().IsDir() {
			continue
		}
		if !file.Mode().IsRegular() {
			return fmt.Errorf("provider ZIP entry %q is not a regular file", file.Name)
		}
		if _, duplicate := files[name]; duplicate {
			return fmt.Errorf("provider ZIP contains duplicate file %q", name)
		}
		reader, err := file.Open()
		if err != nil {
			return fmt.Errorf("read provider ZIP entry %q: %w", name, err)
		}
		digest, hashError := hashContents(reader)
		closeError := reader.Close()
		if err := errors.Join(hashError, closeError); err != nil {
			return fmt.Errorf("hash provider ZIP entry %q: %w", name, err)
		}
		files[name] = digest
	}
	if len(files) == 0 {
		return errors.New("provider ZIP contains no regular files")
	}
	if err := filepath.WalkDir(directory, func(name string, entry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(directory, name)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		expected, found := files[relative]
		if !found {
			return fmt.Errorf("installed package contains undeclared file %q", relative)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("installed package entry %q is not a regular file", relative)
		}
		installed, err := os.Open(name)
		if err != nil {
			return err
		}
		actual, hashError := hashContents(installed)
		closeError := installed.Close()
		if err := errors.Join(hashError, closeError); err != nil {
			return err
		}
		if actual != expected {
			return fmt.Errorf("installed package entry %q differs from the declared archive", relative)
		}
		delete(files, relative)
		return nil
	}); err != nil {
		return err
	}
	if len(files) != 0 {
		missing := make([]string, 0, len(files))
		for name := range files {
			missing = append(missing, name)
		}
		sort.Strings(missing)
		return fmt.Errorf("installed package is missing archive file %q", missing[0])
	}
	return nil
}

func hashContents(reader io.Reader) ([sha256.Size]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return [sha256.Size]byte{}, err
	}
	var digest [sha256.Size]byte
	copy(digest[:], hash.Sum(nil))
	return digest, nil
}
