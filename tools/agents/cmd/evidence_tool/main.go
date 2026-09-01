// Command evidence_tool emits immutable validation evidence and assertions
// from the command line. It performs no network or stateful operation.
package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	if err := execute(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
