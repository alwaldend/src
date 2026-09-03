// Package main implements a persistent Bazel worker for building Hugo sites.
package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorkRequest mirrors the subset of Bazel's worker_protocol.WorkRequest
// needed by this worker.
type WorkRequest struct {
	Arguments  []string
	RequestID  int64
	Cancel     bool
	SandboxDir string
}

// WorkResponse mirrors the subset of Bazel's worker_protocol.WorkResponse
// needed by this worker.
type WorkResponse struct {
	ExitCode  int32
	Output    string
	RequestID int64
}

// WorkerFlagfile is the JSON content passed from the Hugo rule.
type WorkerFlagfile struct {
	SiteArchive   string            `json:"site_archive"`
	Arguments     []string          `json:"arguments"`
	Hugo          string            `json:"hugo"`
	EnvFile       string            `json:"env_file"`
	Postcss       string            `json:"postcss"`
	PostcssBindir string            `json:"postcss_bindir"`
	ToolDirs      []string          `json:"tool_dirs"`
	Env           map[string]string `json:"env"`
	OutputDir     string            `json:"output_dir"`
	Shell         string            `json:"shell"`
}

func readVarint(r *bufio.Reader) (uint64, error) {
	var result uint64
	var shift uint
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
}

func writeVarint(w io.Writer, v uint64) {
	var buf [10]byte
	n := 0
	for v >= 0x80 {
		buf[n] = byte(v) | 0x80
		v >>= 7
		n++
	}
	buf[n] = byte(v)
	w.Write(buf[:n+1])
}

func decodeRequest(data []byte) (*WorkRequest, error) {
	req := &WorkRequest{}
	i := 0
	for i < len(data) {
		tag, err := readVarintFrom(data, &i)
		if err != nil {
			return nil, err
		}
		fieldNum := tag >> 3
		wireType := tag & 0x07
		switch {
		case fieldNum == 1 && wireType == 2: // arguments
			length, err := readVarintFrom(data, &i)
			if err != nil {
				return nil, err
			}
			val := string(data[i : i+int(length)])
			i += int(length)
			req.Arguments = append(req.Arguments, val)
		case fieldNum == 3 && wireType == 0: // request_id
			val, err := readVarintFrom(data, &i)
			if err != nil {
				return nil, err
			}
			req.RequestID = int64(val)
		case fieldNum == 4 && wireType == 0: // cancel
			val, err := readVarintFrom(data, &i)
			if err != nil {
				return nil, err
			}
			req.Cancel = val != 0
		case fieldNum == 6 && wireType == 2: // sandbox_dir
			length, err := readVarintFrom(data, &i)
			if err != nil {
				return nil, err
			}
			req.SandboxDir = string(data[i : i+int(length)])
			i += int(length)
		case wireType == 0: // unknown varint field
			_, err := readVarintFrom(data, &i)
			if err != nil {
				return nil, err
			}
		case wireType == 1: // unknown 64-bit field
			i += 8
		case wireType == 2: // unknown length-delimited field
			length, err := readVarintFrom(data, &i)
			if err != nil {
				return nil, err
			}
			i += int(length)
		case wireType == 5: // unknown 32-bit field
			i += 4
		default:
			return nil, fmt.Errorf("unexpected field %d wire type %d", fieldNum, wireType)
		}
	}
	return req, nil
}

func readVarintFrom(data []byte, pos *int) (uint64, error) {
	var result uint64
	var shift uint
	for {
		if *pos >= len(data) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		b := data[*pos]
		*pos++
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
}

func encodeResponse(resp *WorkResponse) []byte {
	var buf []byte
	// field 1: exit_code (varint)
	buf = appendVarint(buf, (1<<3)|0)
	buf = appendVarint(buf, uint64(resp.ExitCode))
	// field 2: output (string)
	buf = appendVarint(buf, (2<<3)|2)
	buf = appendVarint(buf, uint64(len(resp.Output)))
	buf = append(buf, resp.Output...)
	// field 3: request_id (varint)
	buf = appendVarint(buf, (3<<3)|0)
	buf = appendVarint(buf, uint64(resp.RequestID))
	return buf
}

func appendVarint(buf []byte, v uint64) []byte {
	for v >= 0x80 {
		buf = append(buf, byte(v)|0x80)
		v >>= 7
	}
	return append(buf, byte(v))
}

func run(r io.Reader, w io.Writer) {
	reader := bufio.NewReader(r)
	for {
		length, err := readVarint(reader)
		if err != nil {
			return
		}
		data := make([]byte, length)
		if _, err := io.ReadFull(reader, data); err != nil {
			return
		}
		req, err := decodeRequest(data)
		if err != nil {
			return
		}
		if req.Cancel {
			continue
		}
		resp := handleRequest(req)
		respBytes := encodeResponse(resp)
		writeVarint(w, uint64(len(respBytes)))
		w.Write(respBytes)
	}
}

func handleRequest(req *WorkRequest) *WorkResponse {
	resp := &WorkResponse{RequestID: req.RequestID}
	flagfilePath := ""
	for _, arg := range req.Arguments {
		if strings.HasPrefix(arg, "--flagfile=") {
			flagfilePath = strings.TrimPrefix(arg, "--flagfile=")
		}
	}
	if flagfilePath == "" {
		resp.ExitCode = 1
		resp.Output = "no --flagfile argument provided"
		return resp
	}
	flagfileBytes, err := os.ReadFile(flagfilePath)
	if err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not read flagfile: %v", err)
		return resp
	}
	var ff WorkerFlagfile
	if err := json.Unmarshal(flagfileBytes, &ff); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not parse flagfile: %v", err)
		return resp
	}
	workDir, err := os.Getwd()
	if err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not get working directory: %v", err)
		return resp
	}
	execDir := workDir
	resolve := func(path string) string {
		resolvedPath := filepath.Join(execDir, path)
		if _, err := os.Stat(resolvedPath); err == nil {
			return resolvedPath
		}
		return path
	}
	siteDir := filepath.Join(execDir, fmt.Sprintf("hugo-site-%d", req.RequestID))
	os.RemoveAll(siteDir)
	if err := os.MkdirAll(siteDir, 0o755); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not create site dir: %v", err)
		return resp
	}
	defer os.RemoveAll(siteDir)

	if err := os.MkdirAll(filepath.Join(siteDir, "node_modules/.bin"), 0o755); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not create node_modules: %v", err)
		return resp
	}

	postcssPath := resolve(ff.Postcss)
	postcssBindir := resolve(ff.PostcssBindir)
	postcssShim := fmt.Sprintf(
		"#!/usr/bin/env sh\nset -eu\nexport BAZEL_BINDIR='%s'\nexport NODE_PATH='%s/node_modules'\nexec '%s' \"$@\"\n",
		postcssBindir,
		postcssBindir,
		postcssPath,
	)
	if err := os.WriteFile(filepath.Join(siteDir, "node_modules/.bin/postcss"), []byte(postcssShim), 0o755); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not write postcss shim: %v", err)
		return resp
	}
	archivePath := resolve(ff.SiteArchive)
	if err := extractArchive(archivePath, siteDir); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not extract site archive: %v", err)
		return resp
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
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not create home dir: %v", err)
		return resp
	}
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("could not create temp dir: %v", err)
		return resp
	}
	env = append(env,
		"HOME="+homeDir,
		"TMPDIR="+tmpDir,
		"HUGO_CACHE_DIR="+filepath.Join(siteDir, "_hugo_cache"),
		"HUGO_MODULE_PROXY=off",
		"GOPROXY=off",
		"BAZEL_BINDIR="+postcssBindir,
		"NODE_PATH="+filepath.Join(execDir, ff.PostcssBindir, "node_modules"),
	)
	pathEntries := []string{}
	for _, dir := range ff.ToolDirs {
		pathEntries = append(pathEntries, resolve(dir))
	}
	env = append(env, "PATH="+strings.Join(pathEntries, ":")+":/usr/bin:/bin")
	for key, value := range ff.Env {
		env = append(env, key+"="+value)
	}

	hugoArgs := append([]string{}, ff.Arguments...)
	hugoCmd := exec.Command(resolve(ff.Hugo), hugoCmdArgs(hugoArgs, &ff, resolve(ff.OutputDir), siteDir)...)
	hugoCmd.Dir = siteDir
	hugoCmd.Env = env
	hugoCmd.Stdout = os.Stderr
	hugoCmd.Stderr = os.Stderr
	if err := hugoCmd.Run(); err != nil {
		resp.ExitCode = 1
		resp.Output = fmt.Sprintf("hugo build failed: %v", err)
		return resp
	}
	return resp
}

func extractArchive(archivePath string, destination string) error {
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	var reader io.Reader = archive
	if isGzip(archive) {
		gzipReader, err := gzip.NewReader(archive)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		reader = gzipReader
	} else if _, err := archive.Seek(0, io.SeekStart); err != nil {
		return err
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		path, err := safeArchivePath(destination, header.Name)
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, path); err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}
		}
	}
}

func isGzip(file *os.File) bool {
	magic := make([]byte, 2)
	if _, err := file.ReadAt(magic, 0); err != nil {
		return false
	}
	return magic[0] == 0x1f && magic[1] == 0x8b
}

func safeArchivePath(destination string, name string) (string, error) {
	path := filepath.Join(destination, name)
	relative, err := filepath.Rel(destination, path)
	if err != nil || relative == ".." || filepath.IsAbs(relative) {
		return "", fmt.Errorf("unsafe archive path %q", name)
	}
	return path, nil
}

func hugoCmdArgs(args []string, ff *WorkerFlagfile, outputDir string, siteDir string) []string {
	return append(args,
		"--destination", outputDir,
		"--source", siteDir,
	)
}

func main() {
	run(os.Stdin, os.Stdout)
}
