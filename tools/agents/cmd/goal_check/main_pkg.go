package main

import (
	"fmt"
	"os"

	goalcheck "git.alwaldend.com/alwaldend/src/tools/agents/cmd/goal_check"
)

func main() {
	if err := goalcheck.Run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "goal_check:", err)
		os.Exit(1)
	}
}
