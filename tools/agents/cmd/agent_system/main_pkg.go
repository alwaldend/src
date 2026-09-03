package main

import (
	"fmt"
	"os"

	agentsystem "git.alwaldend.com/alwaldend/src/tools/agents/cmd/agent_system"
)

func main() {
	if err := agentsystem.Run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "agent_system:", err)
		os.Exit(1)
	}
}
