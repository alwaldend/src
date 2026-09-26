package main

import (
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/tools/molecule/internal"
)

func main() {
	if err := internal.Main(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "molecule runner:", err)
		os.Exit(1)
	}
}
