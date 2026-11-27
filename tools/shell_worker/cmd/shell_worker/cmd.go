package main

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/worker_protocol"
)

// execute selects a one-shot action or a stream of Bazel worker requests.
func execute(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	root := &cobra.Command{
		Use: "bazel-genrule-worker", Short: "Bazel worker that runs shell commands",
		SilenceErrors: true, SilenceUsage: true,
	}
	root.SetArgs(args)
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	var persist bool
	var flagfile string
	run := &cobra.Command{
		Use: "run", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !persist && flagfile != "" {
				response := perform(cmd.Context(), &worker_protocol.WorkRequest{Arguments: []string{"--flagfile=" + flagfile}})
				if response.ExitCode != 0 {
					return fmt.Errorf("shell action failed: %s", response.Output)
				}
				return nil
			}
			if err := serve(cmd.Context(), stdin, stdout, persist); err != nil {
				return fmt.Errorf("serve shell requests: %w", err)
			}
			return nil
		},
	}
	run.Flags().BoolVarP(&persist, "persistent_worker", "p", false, "Serve requests until EOF")
	run.Flags().StringVar(&flagfile, "flagfile", "", "Read this one-shot action from a file")
	root.AddCommand(run)
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("run shell worker: %w", err)
	}
	return nil
}
