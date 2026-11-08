package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// main implements only the version invocation used by the isolated fixture.
func main() {
	if err := verifyInvocation(); err != nil {
		log.Fatal(err)
	}
}

// verifyInvocation checks the real AL/injector/launcher environment and files
// before retaining a public receipt at the test's undeclared output boundary.
func verifyInvocation() error {
	if os.Getenv("AL_TERRAFORM_FIXTURE") != "selected" {
		return fmt.Errorf("selected injector environment did not reach Terraform")
	}
	if len(os.Args) != 3 || !strings.HasPrefix(os.Args[1], "-chdir=") || os.Args[2] != "version" {
		return fmt.Errorf("Terraform launcher changed the expected invocation")
	}
	directory := strings.TrimPrefix(os.Args[1], "-chdir=")
	content, err := os.ReadFile(filepath.Join(directory, "main.tf"))
	if err != nil {
		return fmt.Errorf("read packaged configuration: %w", err)
	}
	if !strings.Contains(string(content), `fixture = "terraform-wrapper"`) {
		return fmt.Errorf("packaged configuration does not match the fixture")
	}
	receipt, err := json.MarshalIndent(map[string]any{
		"verified": true, "arguments": os.Args[1:], "environment": "selected", "configuration": "main.tf",
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode invocation receipt: %w", err)
	}
	receipt = append(receipt, '\n')
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "terraform-invocation.json"), receipt, 0o600); err != nil {
			return fmt.Errorf("retain invocation receipt: %w", err)
		}
	}
	if _, err := os.Stdout.Write(receipt); err != nil {
		return fmt.Errorf("write invocation receipt: %w", err)
	}
	return nil
}
