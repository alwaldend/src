package app

import (
	"context"
	"errors"
)

// App provides the service entry points invoked by the command.
type App struct{}

// ServeJobs is the entry point for the job server.
func (App) ServeJobs(context.Context) error {
	return errors.New("job server is not implemented")
}

// ServeWorker is the entry point for the rental-local worker.
func (App) ServeWorker(context.Context) error {
	return errors.New("worker service is not implemented")
}
