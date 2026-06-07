package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger = log.New(os.Stderr, "com.alwaldend.src.tools.vault.tf_backend ", log.Flags())

func main() {
	logger.Printf("Starting plugin")
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	plugin := NewPlugin(ctx)
	go func() {
		<-ctx.Done()
		logger.Printf("Got a cancel signal")
	}()
	if err := plugin.Run(os.Stdin, os.Stdout); err != nil {
		logger.Printf("error: %s", err)
		os.Exit(1)
	}
}
