package main

import "fmt"

// runServer reports that this executable was built without a local engine.
func runServer([]string) error {
	return fmt.Errorf("run Anki server: no local engines are compiled")
}

// subprocessSchemes permits explicitly supplied external worker executables.
var subprocessSchemes = []string{"subprocess"}
