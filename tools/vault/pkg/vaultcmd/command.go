package vaultcmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/go/pkg/flagcompat"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Execute runs an operator command with signal cancellation and its original flag syntax.
func Execute(command *cobra.Command, required ...string) error {
	command.SilenceErrors = true
	command.SilenceUsage = true
	for _, name := range required {
		if err := command.MarkFlagRequired(name); err != nil {
			return fmt.Errorf("configure required flag %q: %w", name, err)
		}
	}
	flagcompat.SetArgs(command, os.Args[1:])
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("%s: %w", command.Name(), err)
	}
	return nil
}

// OutputLimit bounds individual command responses and cached certificate records.
const OutputLimit = 1 << 20

// limitedBuffer rejects subprocess output that exceeds the response size limit.
type limitedBuffer struct {
	// data holds the private output of one subprocess.
	data bytes.Buffer
}

// Write appends a chunk when the entire response still fits in memory.
func (b *limitedBuffer) Write(p []byte) (int, error) {
	if len(p) > OutputLimit-b.data.Len() {
		return 0, errors.New("command response exceeds 1 MiB")
	}
	n, err := b.data.Write(p)
	if err != nil {
		return n, fmt.Errorf("buffer response: %w", err)
	}
	return n, nil
}

// Output runs a command with its configured runfiles environment and bounds its private stdout.
func Output(ctx context.Context, name string, input io.Reader, args ...string) ([]byte, error) {
	var output limitedBuffer
	if err := al.RunCommand(al.CommandArgs{Ctx: ctx, Name: name, Args: args, Stdin: input, Stdout: &output}); err != nil {
		return nil, fmt.Errorf("run %q: %w", name, err)
	}
	return output.data.Bytes(), nil
}

// ReadKey decrypts an encrypted unseal share, prompting on the controlling terminal when omitted.
func ReadKey(ctx context.Context, encrypted string) ([]byte, error) {
	if encrypted == "" {
		value, err := readShare(ctx)
		if err != nil {
			return nil, fmt.Errorf("prompt for unseal share: %w", err)
		}
		encrypted = strings.TrimSpace(value)
	}
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decode encrypted share: %w", err)
	}
	key, err := Output(ctx, "gpg", bytes.NewReader(decoded), "--decrypt")
	if err != nil {
		return nil, fmt.Errorf("decrypt share: %w", err)
	}
	if len(key) == 0 {
		return nil, errors.New("empty unseal share")
	}
	return key, nil
}

// readShare accepts a hidden line of terminal input and restores the original
// mode before returning, including after an interrupt with no pending input.
func readShare(ctx context.Context) (value string, result error) {
	// A separately opened terminal descriptor participates in Go's file poller.
	// Closing it on cancellation interrupts Read without closing the caller's stdin.
	input, err := os.Open(os.Stdin.Name())
	if err != nil {
		return "", fmt.Errorf("open terminal input: %w", err)
	}
	// Check that this terminal supports interruptible reads before changing modes.
	if err := input.SetReadDeadline(time.Time{}); err != nil {
		return "", fmt.Errorf("configure interruptible input: %w", errors.Join(err, input.Close()))
	}
	closed := make(chan error, 1)
	stop := context.AfterFunc(ctx, func() { closed <- input.Close() })
	defer func() {
		var closeErr error
		if stop() {
			closeErr = input.Close()
		} else {
			closeErr = <-closed
		}
		if closeErr != nil {
			result = errors.Join(result, fmt.Errorf("close terminal input: %w", closeErr))
		}
	}()
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("disable terminal echo: %w", err)
	}
	defer func() {
		if err := term.Restore(fd, state); err != nil {
			result = errors.Join(result, fmt.Errorf("restore terminal mode: %w", err))
		}
	}()
	terminal := term.NewTerminal(struct {
		// Reader receives input from the privately owned descriptor.
		io.Reader
		// Writer displays the prompt on standard error.
		io.Writer
	}{input, os.Stderr}, "")
	value, err = terminal.ReadPassword("Enter the encrypted unseal key: ")
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", fmt.Errorf("await unseal share: %w", ctxErr)
	}
	if err != nil {
		return "", fmt.Errorf("read terminal input: %w", err)
	}
	return value, nil
}
