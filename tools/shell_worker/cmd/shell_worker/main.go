package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// main reports worker startup and protocol failures on stderr, leaving stdout for Bazel.
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
