package bazelcmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

// maxMCPOutputBytes keeps escaped structured and text copies below the SDK frame limit.
const maxMCPOutputBytes = mcp.DefaultMaxLineLength / 32

// mcpOptions fixes the workspace, executable and resource limits for one session.
type mcpOptions struct {
	// workspace selects the enclosing Bazel workspace.
	workspace string
	// executable selects Bazel or Bazelisk without shell interpretation.
	executable string
	// startup contains operator-selected Bazel startup arguments.
	startup []string
	// timeout bounds each operation, including time waiting for another request.
	timeout time.Duration
	// maxOutput bounds each captured output stream in bytes.
	maxOutput int
}

// mcpQueryInput describes a graph expression and Bazel command options.
type mcpQueryInput struct {
	// Expression is passed as one literal query argument.
	Expression string `json:"expression" jsonschema:"Bazel query expression, passed as a single argument"`
	// Output selects a textual Bazel output format.
	Output string `json:"output,omitempty" jsonschema:"Text output format; defaults to label for query/cquery and jsonproto for aquery"`
	// Options contains command flags, each in --name or --name=value form.
	Options []string `json:"options,omitempty" jsonschema:"Bazel command flags; use --name=value for options with values"`
}

// mcpTargetsInput selects explicit absolute target patterns for build or test.
type mcpTargetsInput struct {
	// Targets contains the absolute labels or patterns to execute.
	Targets []string `json:"targets" jsonschema:"Nonempty absolute Bazel labels or patterns, such as //package:target or //package/..."`
	// Options contains Bazel flags before the target separator.
	Options []string `json:"options,omitempty" jsonschema:"Bazel command flags; use --name=value for options with values"`
}

// mcpInfoInput selects optional Bazel information keys.
type mcpInfoInput struct {
	// Keys selects information values; omitted keys request the full info report.
	Keys []string `json:"keys,omitempty" jsonschema:"Optional Bazel info keys, such as workspace or bazel-bin"`
	// Options contains Bazel flags before the key separator.
	Options []string `json:"options,omitempty" jsonschema:"Bazel command flags; use --name=value for options with values"`
}

// mcpReport records bounded process output without implying complete graph coverage.
type mcpReport struct {
	// Workspace identifies the selected workspace root.
	Workspace string `json:"workspace"`
	// Command identifies the Bazel operation.
	Command string `json:"command"`
	// ExitCode is the child status, or -1 when execution did not complete normally.
	ExitCode int `json:"exit_code"`
	// Stdout contains at most the configured output byte limit.
	Stdout string `json:"stdout"`
	// Stderr contains bounded child diagnostics separate from protocol output.
	Stderr string `json:"stderr"`
	// StdoutTruncated marks an incomplete stdout stream.
	StdoutTruncated bool `json:"stdout_truncated"`
	// StderrTruncated marks an incomplete stderr stream.
	StderrTruncated bool `json:"stderr_truncated"`
	// TimedOut identifies a server-selected deadline expiration.
	TimedOut bool `json:"timed_out"`
	// Cancelled identifies request or server cancellation.
	Cancelled bool `json:"cancelled"`
	// Error describes failure to start, finish or capture a complete result.
	Error string `json:"error,omitempty"`
}

// limitedOutput retains a prefix while draining all child writes.
type limitedOutput struct {
	// data holds the retained byte prefix.
	data []byte
	// limit is the maximum retained size.
	limit int
	// truncated indicates that at least one byte was discarded.
	truncated bool
}

// Write consumes every byte while keeping memory bounded by limit.
func (output *limitedOutput) Write(data []byte) (int, error) {
	retained := min(len(data), output.limit-len(output.data))
	output.data = append(output.data, data[:retained]...)
	output.truncated = output.truncated || retained < len(data)
	return len(data), nil
}

// mcpWriter borrows a command output stream without owning its lifetime.
type mcpWriter struct {
	// Writer receives protocol frames without transferring ownership.
	io.Writer
}

// Close releases the transport without closing the caller-owned output stream.
func (mcpWriter) Close() error { return nil }

// mcpReader cancels active operations when the stdio connection loses its input.
type mcpReader struct {
	// ReadCloser supplies protocol frames and is closed during server cleanup.
	io.ReadCloser
	// cancel ends the session without waiting for active requests to finish first.
	cancel context.CancelCauseFunc
}

// Read propagates input termination before the SDK waits for pending handlers.
func (reader mcpReader) Read(data []byte) (int, error) {
	n, err := reader.ReadCloser.Read(data)
	if err != nil {
		reader.cancel(err)
		return n, fmt.Errorf("read MCP input: %w", err)
	}
	return n, nil
}

// newMCPCommand exposes AL's Bazel MCP server without starting Bazel during construction.
func newMCPCommand() *cobra.Command {
	command := &cobra.Command{Use: "mcp", Short: "Serve AL tools through the Model Context Protocol"}
	command.AddCommand(newMCPServeCommand())
	return command
}

// newMCPServeCommand constructs the opt-in stdio interface without starting Bazel.
func newMCPServeCommand() *cobra.Command {
	options := mcpOptions{}
	command := &cobra.Command{
		Use: "serve", Short: "Serve workspace-bound Bazel tools over MCP stdio", Args: cobra.NoArgs,
		Long: "Serve query, cquery, aquery, build, test, and info tools over MCP stdio. Uses normal Bazel configuration. Requires a trusted local client and workspace; rules, tests, and options can execute code and contact configured services.",
		RunE: func(command *cobra.Command, _ []string) error {
			if options.timeout <= 0 || options.maxOutput < 1 || options.maxOutput > maxMCPOutputBytes {
				return fmt.Errorf("timeout must be positive and max-output-bytes must be between 1 and %d", maxMCPOutputBytes)
			}
			workspace, err := testWorkspace(options.workspace)
			if err != nil {
				return fmt.Errorf("select MCP workspace: %w", err)
			}
			options.workspace = workspace
			executable, err := exec.LookPath(options.executable)
			if err != nil {
				return fmt.Errorf("find MCP Bazel executable: %w", err)
			}
			options.executable, err = filepath.Abs(executable)
			if err != nil {
				return fmt.Errorf("resolve MCP Bazel executable: %w", err)
			}
			ctx, cancel := context.WithCancelCause(command.Context())
			defer cancel(nil)
			server := newBazelMCPServer(ctx, options)
			reader, ok := command.InOrStdin().(io.ReadCloser)
			if !ok {
				reader = io.NopCloser(command.InOrStdin())
			}
			transport := &mcp.IOTransport{Reader: mcpReader{reader, cancel}, Writer: mcpWriter{command.OutOrStdout()}, MaxLineLength: 1024 * 1024}
			if err := server.Run(ctx, transport); err != nil {
				if errors.Is(context.Cause(ctx), io.EOF) && command.Context().Err() == nil {
					return nil
				}
				return fmt.Errorf("serve Bazel MCP: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&options.workspace, "workspace", "", "Workspace directory (default: current enclosing workspace)")
	command.Flags().StringVar(&options.executable, "bazel-path", "bazel", "Bazel or Bazelisk executable")
	command.Flags().StringArrayVar(&options.startup, "startup-option", nil, "Operator-selected Bazel startup argument (repeatable)")
	command.Flags().DurationVar(&options.timeout, "timeout", 10*time.Minute, "Maximum duration per tool call, including queue time")
	command.Flags().IntVar(&options.maxOutput, "max-output-bytes", 64*1024, fmt.Sprintf("Maximum bytes retained per stdout/stderr stream (1 to %d)", maxMCPOutputBytes))
	return command
}

// newBazelMCPServer registers the fixed catalog and serializes access to one workspace.
func newBazelMCPServer(serverContext context.Context, options mcpOptions) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "al-bazel", Version: "1"}, &mcp.ServerOptions{
		Instructions: "Use query tools to inspect the graph, then build or test explicit targets. Output may be truncated: inspect truncation and exit status before drawing conclusions. Tests run directly; Git-affected selection remains available through al bazel test. Bazel configuration, rules, tests and options are trusted executable inputs. Cancellation stops the client, without rollback guarantees.",
	})
	gate := make(chan struct{}, 1)
	for _, operation := range []string{"query", "cquery", "aquery"} {
		mcp.AddTool(server, &mcp.Tool{Name: "bazel_" + operation, Description: "Inspect the Bazel " + operation + " graph using one expression. Returns bounded textual output and explicit completion status."},
			func(ctx context.Context, _ *mcp.CallToolRequest, input mcpQueryInput) (*mcp.CallToolResult, mcpReport, error) {
				if strings.TrimSpace(input.Expression) == "" || strings.HasPrefix(input.Expression, "-") || strings.ContainsRune(input.Expression, 0) {
					return nil, mcpReport{}, fmt.Errorf("expression must be nonempty and cannot begin with '-' or contain NUL")
				}
				if input.Output == "" {
					input.Output = "label"
					if operation == "aquery" {
						input.Output = "jsonproto"
					}
				}
				if !validMCPOutput(operation, input.Output) {
					return nil, mcpReport{}, fmt.Errorf("unsupported text output format %q for %s", input.Output, operation)
				}
				if err := validateMCPOptions(input.Options); err != nil {
					return nil, mcpReport{}, fmt.Errorf("validate query options: %w", err)
				}
				args := append([]string{operation}, input.Options...)
				args = append(args, "--output="+input.Output, "--", input.Expression)
				return runMCPBazel(ctx, serverContext, options, gate, args)
			})
	}
	for _, operation := range []string{"build", "test"} {
		mcp.AddTool(server, &mcp.Tool{Name: "bazel_" + operation, Description: "Run bazel " + operation + " for explicit target patterns. Executes trusted workspace code using normal Bazel configuration."},
			func(ctx context.Context, _ *mcp.CallToolRequest, input mcpTargetsInput) (*mcp.CallToolResult, mcpReport, error) {
				if len(input.Targets) == 0 {
					return nil, mcpReport{}, fmt.Errorf("at least one target is required")
				}
				if _, err := patternSet(input.Targets); err != nil {
					return nil, mcpReport{}, fmt.Errorf("validate targets: %w", err)
				}
				if err := validateMCPOptions(input.Options); err != nil {
					return nil, mcpReport{}, fmt.Errorf("validate target options: %w", err)
				}
				args := append([]string{operation}, input.Options...)
				args = append(args, "--")
				args = append(args, input.Targets...)
				return runMCPBazel(ctx, serverContext, options, gate, args)
			})
	}
	mcp.AddTool(server, &mcp.Tool{Name: "bazel_info", Description: "Read Bazel workspace and output information. Omit keys for the full report."},
		func(ctx context.Context, _ *mcp.CallToolRequest, input mcpInfoInput) (*mcp.CallToolResult, mcpReport, error) {
			if err := validateMCPOptions(input.Options); err != nil {
				return nil, mcpReport{}, fmt.Errorf("validate info options: %w", err)
			}
			for _, key := range input.Keys {
				if key == "" || strings.HasPrefix(key, "-") || strings.ContainsAny(key, "\x00\r\n") {
					return nil, mcpReport{}, fmt.Errorf("info keys must be nonempty and cannot begin with '-' or contain control characters")
				}
			}
			args := append([]string{"info"}, input.Options...)
			args = append(args, "--")
			args = append(args, input.Keys...)
			return runMCPBazel(ctx, serverContext, options, gate, args)
		})
	return server
}

// validateMCPOptions prevents options from injecting positional commands or separators.
func validateMCPOptions(options []string) error {
	for _, option := range options {
		if !strings.HasPrefix(option, "--") || len(option) <= 2 || strings.ContainsAny(option, "\x00\r\n") {
			return fmt.Errorf("options must use --name or --name=value syntax")
		}
	}
	return nil
}

// validMCPOutput admits textual formats and excludes binary protocol output.
func validMCPOutput(operation, format string) bool {
	formats := map[string]string{
		"query":  " label label_kind build minrank maxrank package location graph xml streamed_jsonproto ",
		"cquery": " label label_kind files transitions jsonproto textproto ",
		"aquery": " text textproto jsonproto commands summary ",
	}
	return format != "" && !strings.ContainsAny(format, " \t\r\n") && strings.Contains(formats[operation], " "+format+" ")
}

// runMCPBazel executes one bounded client and preserves errors as structured tool results.
func runMCPBazel(requestContext, serverContext context.Context, options mcpOptions, gate chan struct{}, args []string) (*mcp.CallToolResult, mcpReport, error) {
	ctx, cancel := context.WithTimeout(requestContext, options.timeout)
	defer cancel()
	stop := context.AfterFunc(serverContext, cancel)
	defer stop()
	report := mcpReport{Workspace: options.workspace, Command: args[0], ExitCode: -1}
	select {
	case gate <- struct{}{}:
		defer func() { <-gate }()
	case <-ctx.Done():
	}
	stdout, stderr := &limitedOutput{limit: options.maxOutput}, &limitedOutput{limit: options.maxOutput}
	child := exec.CommandContext(ctx, options.executable, append(append([]string(nil), options.startup...), args...)...)
	child.Dir = options.workspace
	child.Stdout, child.Stderr = stdout, stderr
	child.WaitDelay = time.Second
	err := child.Run()
	report.Stdout, report.Stderr = string(stdout.data), string(stderr.data)
	report.StdoutTruncated, report.StderrTruncated = stdout.truncated, stderr.truncated
	if child.ProcessState != nil {
		report.ExitCode = child.ProcessState.ExitCode()
	}
	report.TimedOut = errors.Is(ctx.Err(), context.DeadlineExceeded)
	report.Cancelled = errors.Is(ctx.Err(), context.Canceled)
	if err != nil {
		report.Error = fmt.Errorf("run Bazel %s: %w", args[0], err).Error()
	} else if stdout.truncated || stderr.truncated {
		report.Error = "Bazel output exceeded the configured limit; result is incomplete"
	}
	return &mcp.CallToolResult{IsError: report.Error != "" || report.TimedOut || report.Cancelled}, report, nil
}
