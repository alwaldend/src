package workload

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"github.com/spf13/afero"
)

// workloadEnvironment supplies the same tracking context to fetch and execution.
func workloadEnvironment(command *exec.Cmd, id string, specification *jobpb.Workload) []string {
	return append(command.Environ(),
		"MLFLOW_TRACKING_URI="+specification.GetTrackingUri(),
		"MLFLOW_EXPERIMENT_ID="+specification.GetExperimentId(),
		"MLFLOW_RUN_ID="+id,
	)
}

// prepareProject retrieves the MLflow artifact and expands it into a new directory.
func prepareProject(ctx context.Context, filesystem afero.Fs, directory string, id string, specification *jobpb.Workload) (project string, err error) {
	download := filepath.Join(directory, "download")
	if err := filesystem.Mkdir(download, 0o700); err != nil {
		return "", fmt.Errorf("create artifact download directory: %w", err)
	}
	defer func() {
		if cleanupErr := filesystem.RemoveAll(download); cleanupErr != nil {
			err = errors.Join(err, fmt.Errorf("remove artifact download directory: %w", cleanupErr))
		}
	}()
	artifact := specification.GetProjectArtifactPath()
	command := exec.Command("mlflow", "artifacts", "download", "--run-id", id, "--artifact-path", artifact, "--dst-path", download)
	var output strings.Builder
	command.Env = workloadEnvironment(command, id, specification)
	command.Stdout = &output
	command.Stderr = os.Stderr
	if _, err := runCommand(ctx, command); err != nil {
		return "", fmt.Errorf("download MLflow project for %s: %w", id, err)
	}
	// MLflow prints the downloaded path, which depends on the artifact store.
	archive := strings.TrimSpace(output.String())
	archive = archive[strings.LastIndexByte(archive, '\n')+1:]
	relative, err := filepath.Rel(download, archive)
	if err != nil {
		return "", fmt.Errorf("resolve downloaded MLflow artifact: %w", err)
	}
	if relative == "." || !filepath.IsLocal(relative) {
		return "", errors.New("MLflow artifact path is outside its download directory")
	}
	project = filepath.Join(directory, "project")
	if err := filesystem.Mkdir(project, 0o700); err != nil {
		return "", fmt.Errorf("create project directory: %w", err)
	}
	if err := extractProject(ctx, filesystem, archive, project); err != nil {
		return "", fmt.Errorf("expand MLflow project for %s: %w", id, err)
	}
	return project, nil
}

// interruptibleReader bounds cancellation latency while reading the archive.
type interruptibleReader struct {
	// ctx is cancelled when the listener stops the workload.
	ctx context.Context
	// reader supplies the archive bytes.
	reader io.Reader
}

// Read checks cancellation before requesting another block of archive data.
func (reader interruptibleReader) Read(buffer []byte) (int, error) {
	if err := reader.ctx.Err(); err != nil {
		return 0, fmt.Errorf("read project archive: %w", err)
	}
	n, err := reader.reader.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, fmt.Errorf("read compressed project archive: %w", err)
	}
	return n, err
}

// extractProject creates only directories and ordinary files beneath a fresh root.
func extractProject(ctx context.Context, filesystem afero.Fs, archive, root string) (err error) {
	file, err := filesystem.Open(archive)
	if err != nil {
		return fmt.Errorf("open downloaded project archive: %w", err)
	}
	defer func() { err = errors.Join(err, closeResource(file, "project archive")) }()
	compressed, err := gzip.NewReader(interruptibleReader{ctx: ctx, reader: file})
	if err != nil {
		return fmt.Errorf("open gzip project archive: %w", err)
	}
	defer func() { err = errors.Join(err, closeResource(compressed, "gzip archive")) }()
	reader := tar.NewReader(compressed)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			// Consume the gzip footer so a corrupt checksum is reported.
			if _, err := io.Copy(io.Discard, compressed); err != nil {
				return fmt.Errorf("verify gzip project archive: %w", err)
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("read project tar entry: %w", err)
		}
		if err := extractEntry(filesystem, root, header, reader); err != nil {
			return fmt.Errorf("extract project tar entry: %w", err)
		}
	}
}

// extractEntry validates the archive member before creating a filesystem entry.
func extractEntry(filesystem afero.Fs, root string, header *tar.Header, reader io.Reader) error {
	name := strings.TrimSuffix(header.Name, "/")
	if name == "." || !fs.ValidPath(name) || strings.ContainsRune(name, 0) {
		return fmt.Errorf("invalid archive member %q", name)
	}
	destination := filepath.Join(root, filepath.FromSlash(name))
	switch header.Typeflag {
	case tar.TypeDir:
		if err := filesystem.MkdirAll(destination, 0o700); err != nil {
			return fmt.Errorf("create project subdirectory %q: %w", name, err)
		}
	case tar.TypeReg:
		if err := filesystem.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return fmt.Errorf("create parent for archive member %q: %w", name, err)
		}
		mode := fs.FileMode(header.Mode)&0o100 | 0o600
		file, err := filesystem.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
		if err != nil {
			return fmt.Errorf("create archive member %q: %w", name, err)
		}
		_, copyErr := io.Copy(file, reader)
		if copyErr != nil {
			copyErr = fmt.Errorf("copy archive member %q: %w", name, copyErr)
		}
		return errors.Join(copyErr, closeResource(file, name))
	default:
		return fmt.Errorf("unsupported archive member %q of type %d", name, header.Typeflag)
	}
	return nil
}
