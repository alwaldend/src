package main

import (
	"context"
	"fmt"
	"io"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
)

const (
	// stateReady is the initial state before parsing the command line.
	stateReady machine.State = "cli-ready"
	// stateParsing validates arguments and selects a route.
	stateParsing machine.State = "cli-parsing"
	// stateDispatching hands a typed request to the application.
	stateDispatching machine.State = "cli-dispatch"
	// stateExporting writes the full graph of the composed machines.
	stateExporting machine.State = "cli-graph"
)

const (
	// start begins the CLI invocation.
	start machine.Event = "start"
	// dispatch selects the application child.
	dispatch machine.Event = "dispatch"
	// showGraph requests the composed graph.
	showGraph machine.Event = "state-machine graph"
)

// cliRun holds the state and inputs of one CLI invocation.
type cliRun struct {
	// state is the CLI's current state.
	state machine.State
	// arguments are the raw command-line arguments for Cobra.
	arguments []string
	// command holds the application request selected by Cobra.
	command app.Command
	// output receives the state machine graph.
	output io.Writer
}

// cliMachine owns command parsing and starts the application child.
type cliMachine struct {
	// definition holds the root CLI machine's states.
	definition *machine.Machine[*cliRun]
	// triggerApplication starts the application child for the selected command.
	triggerApplication func(context.Context, app.Command) error
	// exportGraph generates the composed DOT for the configured CLI machine.
	exportGraph func(machine.Definition) (string, error)
}

// newMachine configures the CLI state machine around its application child.
func newMachine(
	triggerApplication func(context.Context, app.Command) error,
	exportGraph func(machine.Definition) (string, error),
) *cliMachine {
	cli := &cliMachine{
		definition:         machine.New[*cliRun](stateReady),
		triggerApplication: triggerApplication,
		exportGraph:        exportGraph,
	}
	definition := cli.definition
	definition.Configure(stateReady).Permit(start, stateParsing)
	definition.Configure(stateParsing).OnEntry(cli.parseCommand).
		Permit(dispatch, stateDispatching).Permit(showGraph, stateExporting)
	definition.Configure(stateDispatching).OnEntry(cli.runApplication)
	definition.Configure(stateExporting).OnEntry(cli.writeGraph)
	return cli
}

// run fires the entry event for a CLI invocation.
func (cli *cliMachine) run(ctx context.Context, arguments []string) error {
	current := &cliRun{state: stateReady, arguments: arguments}
	if err := cli.definition.Fire(ctx, &current.state, start, current); err != nil {
		return fmt.Errorf("run CLI: %w", err)
	}
	return nil
}

// parseCommand connects the selected Cobra route to the next CLI transition.
func (cli *cliMachine) parseCommand(ctx context.Context, _ ...any) error {
	current, err := cli.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read arguments for parsing: %w", err)
	}
	runService := func(ctx context.Context, command app.Command) error {
		current.command = command
		if err := cli.definition.Advance(ctx, dispatch); err != nil {
			return fmt.Errorf("enter dispatch state: %w", err)
		}
		return nil
	}
	exportGraph := func(ctx context.Context, output io.Writer) error {
		current.output = output
		if err := cli.definition.Advance(ctx, showGraph); err != nil {
			return fmt.Errorf("enter graph export state: %w", err)
		}
		return nil
	}
	if err := executeCommand(ctx, current.arguments, runService, exportGraph); err != nil {
		return fmt.Errorf("execute CLI command: %w", err)
	}
	return nil
}

// runApplication runs the selected application command.
func (cli *cliMachine) runApplication(ctx context.Context, _ ...any) error {
	current, err := cli.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read the selected application command: %w", err)
	}
	if err := cli.triggerApplication(ctx, current.command); err != nil {
		return fmt.Errorf("run application command: %w", err)
	}
	return nil
}

// writeGraph emits the CLI and all application children in one DOT graph.
func (cli *cliMachine) writeGraph(ctx context.Context, _ ...any) error {
	current, err := cli.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read the graph output stream: %w", err)
	}
	graph, err := cli.exportGraph(cli.definition)
	if err != nil {
		return fmt.Errorf("compose state machine graph: %w", err)
	}
	if _, err := fmt.Fprint(current.output, graph); err != nil {
		return fmt.Errorf("write state machine graph: %w", err)
	}
	return nil
}
