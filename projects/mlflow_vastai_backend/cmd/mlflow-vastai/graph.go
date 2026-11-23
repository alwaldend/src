package main

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
)

// composeStateMachineGraph combines the configured machines and their call edges.
func composeStateMachineGraph(cli, application, worker machine.Definition) (string, error) {
	const applicationSection, workerSection = "application", "worker"
	sections := []machine.GraphSection{
		{ID: "cli", Title: "CLI", DOT: cli.Graph()},
		{ID: applicationSection, Title: "Application", DOT: application.Graph()},
		{ID: workerSection, Title: "Worker", DOT: worker.Graph()},
	}
	handoffs := []machine.Handoff{
		{From: stateDispatching, To: applicationSection, Label: "runs application command"},
		{From: app.StateDispatchWorker, To: workerSection, Label: "runs worker command"},
	}
	dot, err := machine.ComposeGraph(sections, handoffs)
	if err != nil {
		return "", fmt.Errorf("compose CLI graph: %w", err)
	}
	return dot, nil
}
