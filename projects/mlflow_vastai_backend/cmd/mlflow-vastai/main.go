package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/app"
	"github.com/qmuntal/stateless"
	"github.com/spf13/cobra"
)

// serviceApp provides the service actions used by the command machine.
type serviceApp interface {
	ServeJobs(context.Context) error
	ServeWorker(context.Context) error
}

// commandActions holds the dependencies of the CLI state entry actions.
type commandActions struct {
	// application supplies the server and worker entry points.
	application serviceApp
	// command parses arguments and provides the output stream.
	command *cobra.Command
	// dot captures the command graph in its initial state.
	dot string
}

// mode identifies a command state or trigger.
type mode string

const (
	// ready is the initial state before argument parsing.
	ready mode = "ready"
	// start begins command processing.
	start mode = "start"
	// parsing handles command-line arguments and help.
	parsing mode = "parsing"
	// serveJobs selects the job server.
	serveJobs mode = "server"
	// serveWorker selects the rental-local worker.
	serveWorker mode = "worker"
	// graph selects DOT export for the configured command machine.
	graph mode = "state-machine"
)

// main reports command failures and sets the process exit status.
func main() {
	if err := newMachine(app.App{}, os.Args[1:]).Fire(start); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newMachine processes arguments before dispatching the selected service.
func newMachine(application serviceApp, args []string) *stateless.StateMachine {
	machine := stateless.NewStateMachineWithMode(ready, stateless.FiringQueued)
	command := &cobra.Command{
		Use:           "mlflow-vastai",
		Short:         "Run MLflow Vast.ai services",
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	actions := commandActions{application: application, command: command}
	command.AddCommand(
		modeCommand(machine, serveJobs, "serve", "Run the job server", actions.serveJobs),
		modeCommand(machine, serveWorker, "serve", "Run the rental-local worker", actions.serveWorker),
		modeCommand(machine, graph, "graph", "Print the command state machine as DOT", actions.writeGraph),
	)
	command.SetArgs(args)
	machine.Configure(ready).Permit(start, parsing)
	machine.Configure(parsing).OnEntry(actions.parseArguments)
	// Capture the graph while the machine is still in its initial state.
	// Strip the compiler's method value suffix from entry action labels.
	actions.dot = strings.ReplaceAll(machine.ToGraph(), `-fm\l`, `\l`)
	return machine
}

// parseArguments runs Cobra to select a service or graph operation.
func (actions *commandActions) parseArguments(ctx context.Context, _ ...any) error {
	if err := actions.command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("process command: %w", err)
	}
	return nil
}

// serveJobs invokes the job server entry point.
func (actions *commandActions) serveJobs(ctx context.Context, _ ...any) error {
	if err := actions.application.ServeJobs(ctx); err != nil {
		return fmt.Errorf("run mlflow-vastai server serve: %w", err)
	}
	return nil
}

// serveWorker invokes the rental-local worker entry point.
func (actions *commandActions) serveWorker(ctx context.Context, _ ...any) error {
	if err := actions.application.ServeWorker(ctx); err != nil {
		return fmt.Errorf("run mlflow-vastai worker serve: %w", err)
	}
	return nil
}

// writeGraph emits the configured command machine as DOT.
func (actions *commandActions) writeGraph(_ context.Context, _ ...any) error {
	if _, err := fmt.Fprint(actions.command.OutOrStdout(), actions.dot); err != nil {
		return fmt.Errorf("write DOT graph: %w", err)
	}
	return nil
}

// modeCommand registers a transition and its corresponding Cobra route.
func modeCommand(
	machine *stateless.StateMachine,
	target mode,
	verb string,
	description string,
	action stateless.ActionFunc,
) *cobra.Command {
	machine.Configure(parsing).Permit(target, target)
	machine.Configure(target).OnEntry(action)
	group := &cobra.Command{Use: string(target), Short: description}
	group.AddCommand(&cobra.Command{
		Use:   verb,
		Short: description,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := machine.FireCtx(command.Context(), target); err != nil {
				return fmt.Errorf("select %s: %w", command.CommandPath(), err)
			}
			return nil
		},
	})
	return group
}
