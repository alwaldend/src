package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type stringListFlag []string

func (values *stringListFlag) String() string {
	return strings.Join(*values, ",")
}

func (values *stringListFlag) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func run(arguments []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("skill_validator", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var skillPath string
	var openAIPath string
	var expectedName string
	var outputPath string
	var declaredFiles stringListFlag
	flags.StringVar(&skillPath, "skill", "", "path to SKILL.md")
	flags.StringVar(
		&openAIPath,
		"openai-yaml",
		"",
		"optional path to agents/openai.yaml",
	)
	flags.StringVar(
		&expectedName,
		"expected-name",
		"",
		"skill name expected from the Bazel package",
	)
	flags.StringVar(&outputPath, "output", "", "validation report output")
	flags.Var(
		&declaredFiles,
		"file",
		"file declared by skill_library; may be repeated",
	)
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "ERROR: unexpected positional arguments")
		return 2
	}
	requiredOptions := []struct {
		name  string
		value string
	}{
		{name: "--skill", value: skillPath},
		{name: "--expected-name", value: expectedName},
		{name: "--output", value: outputPath},
	}
	for _, option := range requiredOptions {
		if option.value == "" {
			fmt.Fprintf(stderr, "ERROR: %s is required\n", option.name)
			return 2
		}
	}
	if len(declaredFiles) == 0 {
		fmt.Fprintln(stderr, "ERROR: at least one --file is required")
		return 2
	}

	errors := validateSkill(
		skillPath,
		openAIPath,
		expectedName,
		declaredFiles,
	)
	if len(errors) != 0 {
		for _, validationError := range errors {
			fmt.Fprintf(
				stderr,
				"ERROR: skill %s: %s\n",
				expectedName,
				validationError,
			)
		}
		return 1
	}

	report := fmt.Sprintf("Skill is valid: %s\n", expectedName)
	if err := os.WriteFile(outputPath, []byte(report), 0o644); err != nil {
		fmt.Fprintf(stderr, "ERROR: could not write validation report: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, report)
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
