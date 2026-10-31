package bazelcmd

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

	"github.com/spf13/cobra"
)

type doctorReport struct {
	Schema   string         `json:"schema"`
	Runner   doctorRunner   `json:"runner"`
	Source   doctorSource   `json:"source"`
	Bazel    doctorBazel    `json:"bazel"`
	Platform doctorPlatform `json:"platform"`
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

// doctorOptions selects workspace-relative diagnostic inputs and declared rc files.
type doctorOptions struct {
	// workspaceRoot is the base directory for relative input paths.
	workspaceRoot string
	// scratchPath selects the task scratch location to classify.
	scratchPath string
	// bazeliskConfig selects the pin file, defaulting to .bazeliskrc.
	bazeliskConfig string
	// sourceBinary selects the built executable to compare with the running AL.
	sourceBinary string
	// rcFiles declares the rc inventory without inspecting files or their imports.
	rcFiles []string
}

func fileDigest(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

// bazeliskPins reads the required version and archive digest from the selected configuration.
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
		return "", "", fmt.Errorf("Bazelisk configuration must pin version and archive digest")
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

// doctorPath resolves an explicitly selected diagnostic input against the workspace.
func doctorPath(root, value string) string {
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Join(root, value)
}

// buildDoctorReport reads selected pins and executable digests without discovering rc imports.
func buildDoctorReport(options doctorOptions, executable, bazelPath string) (doctorReport, error) {
	root, err := filepath.Abs(options.workspaceRoot)
	if err != nil {
		return doctorReport{}, fmt.Errorf("resolve workspace root: %w", err)
	}
	runnerDigest, err := fileDigest(executable)
	if err != nil {
		return doctorReport{}, fmt.Errorf("digest runner: %w", err)
	}
	if options.bazeliskConfig == "" {
		options.bazeliskConfig = ".bazeliskrc"
	}
	pinFile := doctorPath(root, options.bazeliskConfig)
	version, archiveDigest, err := bazeliskPins(pinFile)
	if err != nil {
		return doctorReport{}, fmt.Errorf("read Bazelisk pins %q: %w", pinFile, err)
	}
	if options.sourceBinary == "" {
		options.sourceBinary = filepath.Join("bazel-bin", "external", "rules_al+", "cmd", "al", "al_", "al")
	}
	sourceBinary := doctorPath(root, options.sourceBinary)
	source := doctorSource{Target: "@rules_al//:al"}
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
	if len(options.rcFiles) == 0 {
		options.rcFiles = []string{".bazelrc"}
	}
	rcFiles := make([]string, 0, len(options.rcFiles))
	for _, configured := range options.rcFiles {
		resolved := doctorPath(root, configured)
		relative, err := filepath.Rel(root, resolved)
		if err != nil {
			rcFiles = append(rcFiles, filepath.ToSlash(resolved))
		} else {
			rcFiles = append(rcFiles, filepath.ToSlash(relative))
		}
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
		RCFiles:  rcFiles,
		Scratch:  taskScratch(root, options.scratchPath),
	}, nil
}

// runDoctor parses diagnostic paths and writes one local runner and workspace report.
func runDoctor(args []string) error {
	options := doctorOptions{workspaceRoot: "."}
	for index := 0; index < len(args); index++ {
		flag := args[index]
		if flag == "--help" {
			fmt.Println("Usage: al doctor [--workspace-root PATH] [--task-scratch PATH] [--bazelisk-config PATH] [--source-binary PATH] [--rc-file PATH ...]")
			return nil
		}
		if index+1 >= len(args) {
			return fmt.Errorf("%s requires a value", flag)
		}
		value := args[index+1]
		index++
		switch flag {
		case "--workspace-root":
			options.workspaceRoot = value
		case "--task-scratch":
			options.scratchPath = value
		case "--bazelisk-config", "--source-binary", "--rc-file":
			if value == "" {
				return fmt.Errorf("%s requires a nonempty path", flag)
			}
			switch flag {
			case "--bazelisk-config":
				options.bazeliskConfig = value
			case "--source-binary":
				options.sourceBinary = value
			case "--rc-file":
				options.rcFiles = append(options.rcFiles, value)
			}
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
	report, err := buildDoctorReport(options, executable, bazelPath)
	if err != nil {
		return fmt.Errorf("build doctor report: %w", err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("encode doctor report: %w", err)
	}
	return nil
}

// NewDoctorCommand exposes bounded runner and workspace diagnostics.
func NewDoctorCommand() *cobra.Command {
	return &cobra.Command{Use: "doctor", Short: "Inspect AL and Bazel workspace configuration", DisableFlagParsing: true, RunE: func(cmd *cobra.Command, args []string) error {
		if err := runDoctor(args); err != nil {
			return fmt.Errorf("inspect runner: %w", err)
		}
		return nil
	}}
}
