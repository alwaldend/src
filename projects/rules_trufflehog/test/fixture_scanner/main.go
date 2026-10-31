package main

import (
	"fmt"
	"os"
	"slices"
)

// main checks rule-to-command wiring without contacting credential services.
func main() {
	if !slices.Contains(os.Args, "--no-verification") {
		os.Exit(1)
	}
	fmt.Println(`{"DetectorName":"URI","Verified":false,"Raw":"rule-fixture","SourceMetadata":{"Data":{"Git":{"file":"fixture.txt","commit":"0000000000000000000000000000000000000000","line":1}}}}`)
	os.Exit(183)
}
