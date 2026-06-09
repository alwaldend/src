package fp

import (
	"context"
	"fmt"
	"time"
)

func Run(ctx context.Context, start func(context.Context) error, stop func(context.Context) error, stopTimeout time.Duration) error {
	if err := start(ctx); err != nil {
		return fmt.Errorf("could not start: %w", err)
	}
	<-ctx.Done()
	stopCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	defer cancel()
	if err := stop(stopCtx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}
