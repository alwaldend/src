package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/provider"
)

// main starts the provider process and reports startup failure to Terraform.
func main() {
	if err := providerserver.Serve(context.Background(), provider.New("0.0.0-dev"), providerserver.ServeOpts{Address: "registry.terraform.io/alwaldend/vastai"}); err != nil {
		fmt.Fprintf(os.Stderr, "serve Vast.ai provider: %v\n", err)
		os.Exit(1)
	}
}
