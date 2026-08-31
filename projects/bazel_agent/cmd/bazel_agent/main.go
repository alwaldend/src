package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

type (
	lookPathFunc       func(string) (string, error)
	replaceProcessFunc func(string, []string, []string) error
)

type doctorReport struct {
	Schema   string         `json:"schema"`
	Runner   doctorRunner   `json:"runner"`
	Source   doctorSource   `json:"source"`
	Bazel    doctorBazel    `json:"bazel"`
	Platform doctorPlatform `json:"platform"`
	Profile  []string       `json:"profile"`
	RCFiles  []string       `json:"rcFiles"`
	Scratch  doctorScratch  `json:"scratch"`
}

type doctorRunner struct {
	Executable string `json:"executable"`
	Digest     string `json:"digest"`
}

type doctorSource struct {
	Target       string `json:"target"`
	Binary       string `json:"binary,omitempty"`
	Digest       string `json:"digest,omitempty"`
	StaleInstall *bool  `json:"staleInstall,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

type doctorBazel struct {
	Executable    string `json:"executable"`
	Version       string `json:"version"`
	ArchiveSHA256 string `json:"archiveSha256"`
}

type doctorPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type doctorScratch struct {
	Path       string `json:"path,omitempty"`
	Namespaced bool   `json:"namespaced"`
	Reason     string `json:"reason,omitempty"`
}

func bazelArguments(args []string) []string {
	result := make([]string, 0, len(args)+2)
	result = append(result, "--batch")
	if len(args) == 0 {
		return result
	}
	result = append(result, args[0], "--config=agent")
	return append(result, args[1:]...)
}

func run(
	args []string,
	environment []string,
	lookPath lookPathFunc,
	replaceProcess replaceProcessFunc,
) error {
	bazelPath, err := lookPath("bazel")
	if err != nil {
		return fmt.Errorf("find bazel in PATH: %w", err)
	}
	processArgs := append([]string{bazelPath}, bazelArguments(args)...)
	if err := replaceProcess(bazelPath, processArgs, environment); err != nil {
		return fmt.Errorf("execute bazel: %w", err)
	}
	return nil
}

func fileDigest(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func bazeliskPins(path string) (string, string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	values := map[string]string{}
	for _, line := range strings.Split(string(content), "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if ok {
			values[key] = value
		}
	}
	version := values["USE_BAZEL_VERSION"]
	digest := values["BAZELISK_VERIFY_SHA256"]
	if version == "" || digest == "" {
		return "", "", fmt.Errorf(".bazeliskrc must pin version and archive digest")
	}
	return version, digest, nil
}

func taskScratch(workspaceRoot, candidate string) doctorScratch {
	if candidate == "" {
		return doctorScratch{Reason: "unavailable"}
	}
	absolute := candidate
	if !filepath.IsAbs(absolute) {
		absolute = filepath.Join(workspaceRoot, absolute)
	}
	absolute = filepath.Clean(absolute)
	relative, err := filepath.Rel(filepath.Join(workspaceRoot, "out"), absolute)
	if err != nil || relative == "." || relative == "" || relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return doctorScratch{Path: absolute, Reason: "conflict"}
	}
	parts := strings.Split(filepath.ToSlash(relative), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return doctorScratch{Path: absolute, Reason: "partial"}
	}
	return doctorScratch{Path: absolute, Namespaced: true}
}

func buildDoctorReport(workspaceRoot, scratchPath, executable, bazelPath string) (doctorReport, error) {
	root, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return doctorReport{}, fmt.Errorf("resolve workspace root: %w", err)
	}
	runnerDigest, err := fileDigest(executable)
	if err != nil {
		return doctorReport{}, fmt.Errorf("digest runner: %w", err)
	}
	version, archiveDigest, err := bazeliskPins(filepath.Join(root, ".bazeliskrc"))
	if err != nil {
		return doctorReport{}, fmt.Errorf("read Bazelisk pins: %w", err)
	}
	sourceBinary := filepath.Join(root, "bazel-bin", "projects", "bazel_agent", "cmd", "bazel_agent", "bazel_agent_", "bazel_agent")
	source := doctorSource{Target: "//projects/bazel_agent/cmd/bazel_agent:bazel_agent"}
	if sourceDigest, digestErr := fileDigest(sourceBinary); digestErr == nil {
		stale := sourceDigest != runnerDigest
		source.Binary = sourceBinary
		source.Digest = sourceDigest
		source.StaleInstall = &stale
	} else if os.IsNotExist(digestErr) {
		source.Reason = "unavailable"
	} else {
		return doctorReport{}, fmt.Errorf("digest built runner: %w", digestErr)
	}
	return doctorReport{
		Schema: "agents.alwaldend.com/bazel-agent-doctor/v1alpha1",
		Runner: doctorRunner{Executable: executable, Digest: runnerDigest},
		Source: source,
		Bazel: doctorBazel{
			Executable:    bazelPath,
			Version:       version,
			ArchiveSHA256: archiveDigest,
		},
		Platform: doctorPlatform{OS: runtime.GOOS, Arch: runtime.GOARCH},
		Profile:  []string{"--batch", "--config=agent"},
		RCFiles: []string{
			".bazelrc",
			"tools/bazelrc/root.bazelrc",
			"tools/bazelrc/preset.bazelrc",
			"tools/bazelrc/project.bazelrc",
			"user.bazelrc (optional, last)",
		},
		Scratch: taskScratch(root, scratchPath),
	}, nil
}

func runDoctor(args []string) error {
	workspaceRoot := "."
	scratchPath := ""
	for index := 0; index < len(args); index++ {
		flag := args[index]
		if flag == "--help" {
			fmt.Println("Usage: bazel_agent doctor [--workspace-root PATH] [--task-scratch PATH]")
			return nil
		}
		if index+1 >= len(args) {
			return fmt.Errorf("%s requires a value", flag)
		}
		value := args[index+1]
		index++
		switch flag {
		case "--workspace-root":
			workspaceRoot = value
		case "--task-scratch":
			scratchPath = value
		default:
			return fmt.Errorf("unknown doctor option %q", flag)
		}
	}
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve runner executable: %w", err)
	}
	bazelPath, err := exec.LookPath("bazel")
	if err != nil {
		return fmt.Errorf("find bazel in PATH: %w", err)
	}
	report, err := buildDoctorReport(workspaceRoot, scratchPath, executable, bazelPath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func main() {
	args := os.Args[1:]
	var err error
	if len(args) > 0 && args[0] == "doctor" {
		err = runDoctor(args[1:])
	} else {
		err = run(args, os.Environ(), exec.LookPath, syscall.Exec)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "bazel_agent: %v\n", err)
		os.Exit(1)
	}
}
