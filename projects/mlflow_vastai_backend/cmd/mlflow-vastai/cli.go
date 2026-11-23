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
	// operation runs the route selected by Cobra.
	operation func(context.Context) error
}

// cli owns command parsing and starts the application child.
type cli struct {
	// definition holds the root CLI machine's states.
	definition *machine.Machine[*cliRun]
	// application runs the selected command.
	application *app.App
}

// newCLI configures the CLI state machine around its application child.
func newCLI(application *app.App) *cli {
	cli := &cli{
		definition:  machine.New[*cliRun](stateReady),
		application: application,
	}
	definition := cli.definition
	definition.Configure(stateReady).Permit(start, stateParsing)
	definition.Configure(stateParsing).OnEntry(cli.parseCommand).
		Permit(dispatch, stateDispatching).Permit(showGraph, stateExporting)
	definition.Configure(stateDispatching).OnEntry(cli.runOperation)
	definition.Configure(stateExporting).OnEntry(cli.runOperation)
	return cli
}

// run fires the entry event for a CLI invocation.
func (cli *cli) run(ctx context.Context, arguments []string) error {
	current := &cliRun{state: stateReady, arguments: arguments}
	ctx = cli.definition.WithInvocation(ctx, &current.state, current)
	if err := cli.definition.FireCtx(ctx, start); err != nil {
		return fmt.Errorf("run CLI: %w", err)
	}
	return nil
}

// parseCommand connects the selected Cobra route to the next CLI transition.
func (cli *cli) parseCommand(ctx context.Context, _ ...any) error {
	current, err := cli.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read arguments for parsing: %w", err)
	}
	if err := executeCommand(ctx, current.arguments, cli); err != nil {
		return fmt.Errorf("execute CLI command: %w", err)
	}
	return nil
}

// selectRoute enters the state that runs the selected operation.
func (cli *cli) selectRoute(ctx context.Context, event machine.Event, operation func(context.Context) error) error {
	current, err := cli.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read CLI invocation for route selection: %w", err)
	}
	current.operation = operation
	if err := cli.definition.FireCtx(ctx, event); err != nil {
		return fmt.Errorf("select CLI route %s: %w", event, err)
	}
	return nil
}

// runOperation calls the method selected by Cobra.
func (cli *cli) runOperation(ctx context.Context, _ ...any) error {
	current, err := cli.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read the selected CLI operation: %w", err)
	}
	if err := current.operation(ctx); err != nil {
		return fmt.Errorf("run CLI operation: %w", err)
	}
	return nil
}

// graph composes the CLI graph and the application's child graphs.
func (cli *cli) graph() (string, error) {
	section := machine.GraphSection{
		ID: "cli", Title: "CLI", DOT: cli.definition.ToGraph(),
		Calls: []machine.GraphCall{
			{From: stateDispatching, Label: "runs application command", Target: cli.application.Graph()},
		},
	}
	dot, err := machine.ComposeGraph(section)
	if err != nil {
		return "", fmt.Errorf("compose CLI graph: %w", err)
	}
	return dot, nil
}

// writeGraph emits the CLI and all application children in one DOT graph.
func (cli *cli) writeGraph(output io.Writer) error {
	graph, err := cli.graph()
	if err != nil {
		return fmt.Errorf("compose state machine graph: %w", err)
	}
	if _, err := fmt.Fprint(output, graph); err != nil {
		return fmt.Errorf("write state machine graph: %w", err)
	}
	return nil
}
