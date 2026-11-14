package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/comfyui_tools/internal/experiments"
	"github.com/spf13/cobra"
)

// command constructs the literal-argv experiment controller interface.
func command() *cobra.Command {
	root := &cobra.Command{Use: "experiments", Short: "Run buffered experiments on an existing ComfyUI endpoint", SilenceErrors: true, SilenceUsage: true}
	var directory, config string
	var cancelOwned bool
	run := &cobra.Command{Use: "run", Short: "Start or resume a private experiment stream", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if err := experiments.Run(cmd.Context(), config, directory); err != nil {
			return fmt.Errorf("run experiments: %w", err)
		}
		return nil
	}}
	run.Flags().StringVar(&config, "config", "", "Absolute JSON configuration path")
	status := &cobra.Command{Use: "status", Short: "Read the private durable journal as JSON", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		data, err := experiments.Status(directory)
		if err != nil {
			return fmt.Errorf("read experiment status: %w", err)
		}
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(data); err != nil {
			return fmt.Errorf("write experiment status: %w", err)
		}
		return nil
	}}
	stop := &cobra.Command{Use: "stop", Short: "Request drain; optionally cancel only recorded jobs", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if err := experiments.Stop(directory, cancelOwned); err != nil {
			return fmt.Errorf("stop experiments: %w", err)
		}
		return nil
	}}
	stop.Flags().BoolVar(&cancelOwned, "cancel-owned", false, "Cancel only IDs owned by this run when its controller consumes the request")
	for _, cmd := range []*cobra.Command{run, status, stop} {
		cmd.Flags().StringVar(&directory, "directory", "", "Absolute private run directory")
		if err := cmd.MarkFlagRequired("directory"); err != nil {
			panic(fmt.Errorf("declare required directory flag: %w", err))
		}
	}
	if err := run.MarkFlagRequired("config"); err != nil {
		panic(fmt.Errorf("declare required configuration flag: %w", err))
	}
	root.AddCommand(run, status, stop)
	return root
}

// main handles process signals without cancelling or changing the remote instance.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
