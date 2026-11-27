package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// maxFrameBytes bounds the memory used to receive a Bazel request.
const maxFrameBytes = 64 << 20

// maxFlagfileBytes limits the configuration decoded for a single action.
const maxFlagfileBytes = 4 << 20

// workRequest contains the Bazel fields needed by the serial Hugo worker.
type workRequest struct {
	// arguments carries the arguments from the Hugo build action.
	arguments []string
	// id correlates the response with the request.
	id int32
	// cancel identifies a cancellation for an action already served.
	cancel bool
}

// workerFlagfile is the JSON configuration produced by the Hugo rule.
type workerFlagfile struct {
	// SiteArchive is the site source archive in the execution root.
	SiteArchive string `json:"site_archive"`
	// Arguments contains additional Hugo command arguments.
	Arguments []string `json:"arguments"`
	// Hugo identifies the resolved build tool.
	Hugo string `json:"hugo"`
	// Postcss identifies the PostCSS executable.
	Postcss string `json:"postcss"`
	// PostcssBindir contains the PostCSS Node modules.
	PostcssBindir string `json:"postcss_bindir"`
	// ToolDirs lists additional executable search directories.
	ToolDirs []string `json:"tool_dirs"`
	// Env contains action-specific environment settings.
	Env map[string]string `json:"env"`
	// OutputDir is the declared output tree in the execution root.
	OutputDir string `json:"output_dir"`
}

// decodeRequest reads the supported protobuf fields and skips unknown scalar fields.
func decodeRequest(data []byte) (workRequest, error) {
	var req workRequest
	reader := bytes.NewReader(data)
	for reader.Len() > 0 {
		tag, err := binary.ReadUvarint(reader)
		if err != nil {
			return req, fmt.Errorf("read worker field tag: %w", err)
		}
		field := tag >> 3
		if field == 0 {
			return req, fmt.Errorf("invalid worker field number 0")
		}
		var length uint64
		switch tag & 7 {
		case 0:
			n, err := binary.ReadUvarint(reader)
			if err != nil {
				return req, fmt.Errorf("read worker field %d: %w", field, err)
			}
			switch field {
			case 3:
				req.id = int32(n)
			case 4:
				req.cancel = n != 0
			}
			continue
		case 1:
			length = 8
		case 2:
			length, err = binary.ReadUvarint(reader)
			if err != nil {
				return req, fmt.Errorf("read worker field %d length: %w", field, err)
			}
		case 5:
			length = 4
		default:
			return req, fmt.Errorf("unsupported worker wire type %d", tag&7)
		}
		if length > uint64(reader.Len()) {
			return req, fmt.Errorf("worker field %d extends past the request", field)
		}
		if field == 1 && tag&7 == 2 {
			// The rule passes a single flagfile. Leave room for alternate callers.
			if len(req.arguments) >= 1024 {
				return req, fmt.Errorf("worker request exceeds 1024 arguments")
			}
			offset := len(data) - reader.Len()
			req.arguments = append(req.arguments, string(data[offset:offset+int(length)]))
		}
		if _, err := reader.Seek(int64(length), io.SeekCurrent); err != nil {
			return req, fmt.Errorf("skip worker field %d: %w", field, err)
		}
	}
	return req, nil
}

// writeResponse emits a framed protobuf, returning command diagnostics to Bazel.
func writeResponse(w io.Writer, id int32, buildErr error) error {
	var response []byte
	if buildErr != nil {
		message := strings.ToValidUTF8(buildErr.Error(), "\uFFFD")
		response = append(response, 8, 1, 18)
		response = binary.AppendUvarint(response, uint64(len(message)))
		response = append(response, message...)
	}
	response = append(response, 24)
	response = binary.AppendUvarint(response, uint64(id))
	frame := binary.AppendUvarint(nil, uint64(len(response)))
	frame = append(frame, response...)
	if n, err := w.Write(frame); err != nil {
		return fmt.Errorf("write worker response: %w", err)
	} else if n != len(frame) {
		return fmt.Errorf("write worker response: %w", io.ErrShortWrite)
	}
	return nil
}

// run serves length-delimited protobuf requests until the parent closes the pipe.
func run(ctx context.Context, r io.Reader, w io.Writer) error {
	reader := bufio.NewReader(r)
	for {
		length, err := binary.ReadUvarint(reader)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read worker request length: %w", err)
		}
		if length > maxFrameBytes {
			return fmt.Errorf("worker frame length %d exceeds %d bytes", length, maxFrameBytes)
		}
		data := make([]byte, length)
		if _, err := io.ReadFull(reader, data); err != nil {
			return fmt.Errorf("read worker frame: %w", err)
		}
		req, err := decodeRequest(data)
		if err != nil {
			return fmt.Errorf("decode worker frame: %w", err)
		}
		if req.cancel {
			continue
		}
		if err := writeResponse(w, req.id, handleRequest(ctx, req.arguments)); err != nil {
			return fmt.Errorf("respond to request %d: %w", req.id, err)
		}
	}
}

// handleRequest prepares an isolated site and executes Hugo for the requested flagfile.
func handleRequest(ctx context.Context, arguments []string) (result error) {
	var flagfilePath string
	for _, arg := range arguments {
		if path, ok := strings.CutPrefix(arg, "--flagfile="); ok {
			flagfilePath = path
		}
	}
	if flagfilePath == "" {
		return fmt.Errorf("no --flagfile argument provided")
	}
	f, err := os.Open(flagfilePath)
	if err != nil {
		return fmt.Errorf("open flagfile %q: %w", flagfilePath, err)
	}
	data, readErr := io.ReadAll(io.LimitReader(f, maxFlagfileBytes+1))
	if err := errors.Join(readErr, f.Close()); err != nil {
		return fmt.Errorf("read flagfile %q: %w", flagfilePath, err)
	}
	if len(data) > maxFlagfileBytes {
		return fmt.Errorf("flagfile %q exceeds %d bytes", flagfilePath, maxFlagfileBytes)
	}
	var ff workerFlagfile
	if err := json.Unmarshal(data, &ff); err != nil {
		return fmt.Errorf("decode flagfile %q: %w", flagfilePath, err)
	}
	execDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("locate execution root: %w", err)
	}
	resolve := func(path string) string {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(execDir, path)
	}
	siteDir, err := os.MkdirTemp(execDir, "hugo-site-")
	if err != nil {
		return fmt.Errorf("create site directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(siteDir); err != nil {
			result = errors.Join(result, fmt.Errorf("remove site directory: %w", err))
		}
	}()

	if err := os.MkdirAll(filepath.Join(siteDir, "node_modules/.bin"), 0o755); err != nil {
		return fmt.Errorf("could not create node_modules: %w", err)
	}

	postcssPath := resolve(ff.Postcss)
	postcssBindir := resolve(ff.PostcssBindir)
	// Hugo filters the environment of its asset processors. Restore the
	// runfile locations in the shim so they reach the PostCSS binary.
	postcssShim := fmt.Sprintf(
		"#!/usr/bin/env sh\nset -eu\nexport BAZEL_BINDIR=%s\nexport NODE_PATH=%s\nexec %s \"$@\"\n",
		shellLiteral(postcssBindir),
		shellLiteral(filepath.Join(postcssBindir, "node_modules")),
		shellLiteral(postcssPath),
	)
	if err := os.WriteFile(filepath.Join(siteDir, "node_modules/.bin/postcss"), []byte(postcssShim), 0o755); err != nil {
		return fmt.Errorf("could not write postcss shim: %w", err)
	}
	archivePath := resolve(ff.SiteArchive)
	if err := extractArchive(archivePath, siteDir); err != nil {
		return fmt.Errorf("could not extract site archive: %w", err)
	}

	env := os.Environ()
	filteredEnv := make([]string, 0, len(env)+8)
	for _, value := range env {
		if !strings.HasPrefix(value, "HOME=") && !strings.HasPrefix(value, "TMPDIR=") {
			filteredEnv = append(filteredEnv, value)
		}
	}
	env = filteredEnv
	homeDir := filepath.Join(siteDir, "_home")
	tmpDir := filepath.Join(siteDir, "_tmp")
	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		return fmt.Errorf("could not create home dir: %w", err)
	}
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return fmt.Errorf("could not create temp dir: %w", err)
	}
	env = append(env,
		"HOME="+homeDir,
		"TMPDIR="+tmpDir,
		"HUGO_CACHE_DIR="+filepath.Join(siteDir, "_hugo_cache"),
		"HUGO_MODULE_PROXY=off",
		"GOPROXY=off",
		"BAZEL_BINDIR="+postcssBindir,
		"NODE_PATH="+filepath.Join(postcssBindir, "node_modules"),
	)
	pathEntries := make([]string, 0, len(ff.ToolDirs))
	for _, dir := range ff.ToolDirs {
		pathEntries = append(pathEntries, resolve(dir))
	}
	env = append(env, "PATH="+strings.Join(pathEntries, ":")+":/usr/bin:/bin")
	for key, value := range ff.Env {
		env = append(env, key+"="+value)
	}

	hugoArgs := append(ff.Arguments, "--destination", resolve(ff.OutputDir), "--source", siteDir)
	hugoCmd := exec.CommandContext(ctx, resolve(ff.Hugo), hugoArgs...)
	hugoCmd.Dir = siteDir
	hugoCmd.Env = env
	hugoCmd.Stdout = os.Stderr
	hugoCmd.Stderr = os.Stderr
	if err := hugoCmd.Run(); err != nil {
		return fmt.Errorf("hugo build failed: %w", err)
	}
	return nil
}

// shellLiteral quotes a literal path for the POSIX shell used by the PostCSS shim.
func shellLiteral(path string) string {
	return "'" + strings.ReplaceAll(path, "'", `'\''`) + "'"
}

// extractArchive expands regular files, directories and symlinks within the new site root.
func extractArchive(archivePath, destination string) (result error) {
	archive, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open site archive %q: %w", archivePath, err)
	}
	defer func() {
		if err := archive.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close site archive: %w", err))
		}
	}()
	root, err := os.OpenRoot(destination)
	if err != nil {
		return fmt.Errorf("open extraction directory %q: %w", destination, err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close extraction directory: %w", err))
		}
	}()
	buffer := bufio.NewReader(archive)
	var reader io.Reader = buffer
	magic, err := buffer.Peek(2)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("inspect site archive: %w", err)
	}
	if bytes.Equal(magic, []byte{0x1f, 0x8b}) {
		compressed, err := gzip.NewReader(buffer)
		if err != nil {
			return fmt.Errorf("open compressed site archive: %w", err)
		}
		defer func() {
			if err := compressed.Close(); err != nil {
				result = errors.Join(result, fmt.Errorf("close compressed archive: %w", err))
			}
		}()
		reader = compressed
	}
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read site archive header: %w", err)
		}
		path := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if err := root.MkdirAll(path, header.FileInfo().Mode().Perm()); err != nil {
				return fmt.Errorf("create archive directory %q: %w", path, err)
			}
		case tar.TypeReg, tar.TypeSymlink:
			if err := root.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return fmt.Errorf("create parent of %q: %w", path, err)
			}
			if header.Typeflag == tar.TypeSymlink {
				if err := root.Symlink(header.Linkname, path); err != nil && !errors.Is(err, os.ErrExist) {
					return fmt.Errorf("create archive link %q: %w", path, err)
				}
				continue
			}
			file, err := root.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode().Perm())
			if err != nil {
				return fmt.Errorf("create archive file %q: %w", path, err)
			}
			_, copyErr := io.Copy(file, tr)
			if err := errors.Join(copyErr, file.Close()); err != nil {
				return fmt.Errorf("extract archive file %q: %w", path, err)
			}
		}
	}
}

// execute coordinates one-shot actions and the persistent worker protocol.
func execute(ctx context.Context, args []string, input io.Reader, output, diagnostics io.Writer) error {
	var flagfile string
	var persist bool
	action := func(cmd *cobra.Command, args []string) error {
		if flagfile != "" && !persist {
			if err := handleRequest(cmd.Context(), []string{"--flagfile=" + flagfile}); err != nil {
				return fmt.Errorf("build Hugo site: %w", err)
			}
		} else if err := run(cmd.Context(), input, output); err != nil {
			return fmt.Errorf("serve Hugo worker: %w", err)
		}
		return nil
	}
	root := &cobra.Command{Use: "hugo_worker", Short: "Build Hugo sites for Bazel", Args: cobra.NoArgs, RunE: action, SilenceErrors: true, SilenceUsage: true}
	root.SetArgs(args)
	root.SetIn(input)
	root.SetOut(diagnostics)
	root.SetErr(diagnostics)
	root.PersistentFlags().StringVar(&flagfile, "flagfile", "", "Read this one-shot action from a file")
	root.PersistentFlags().BoolVar(&persist, "persistent_worker", false, "Serve requests until EOF")
	root.AddCommand(&cobra.Command{Use: "run", Args: cobra.NoArgs, RunE: action})
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("run Hugo worker: %w", err)
	}
	return nil
}

// main reserves stdout for Bazel and reports command errors to stderr.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// An inherited stdin can be a blocking file descriptor. A pipe lets shutdown
	// interrupt the protocol reader even while Bazel keeps its end open.
	input, writer := io.Pipe()
	go func() {
		_, err := io.Copy(writer, os.Stdin)
		_ = writer.CloseWithError(err)
	}()
	closeInput := context.AfterFunc(ctx, func() { _ = input.CloseWithError(ctx.Err()) })
	err := execute(ctx, os.Args[1:], input, os.Stdout, os.Stderr)
	closeInput()
	_ = input.Close()
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
