package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateArgsAcceptsSavedPlan(t *testing.T) {
	chdir := t.TempDir()
	planPath := filepath.Join(chdir, "reviewed.tfplan")
	if err := os.WriteFile(planPath, []byte("plan contents validated by Terraform"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name string
		path string
	}{
		{name: "absolute", path: planPath},
		{name: "relative to chdir", path: "reviewed.tfplan"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := validateArgs([]string{"apply", test.path}, chdir, true); err != nil {
				t.Fatalf("saved plan rejected: %v", err)
			}
		})
	}
}

func TestValidateArgsRejectsUnsavedApply(t *testing.T) {
	chdir := t.TempDir()
	planPath := filepath.Join(chdir, "reviewed.tfplan")
	if err := os.WriteFile(planPath, []byte("plan fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name string
		args []string
	}{
		{name: "empty command"},
		{name: "bare apply", args: []string{"apply"}},
		{name: "empty path", args: []string{"apply", ""}},
		{name: "plan command", args: []string{"plan", planPath}},
		{name: "destroy command", args: []string{"destroy", planPath}},
		{name: "approval flag", args: []string{"apply", "-auto-approve"}},
		{name: "target flag", args: []string{"apply", "-target=module.dns"}},
		{name: "extra flag", args: []string{"apply", planPath, "-auto-approve"}},
		{name: "extra file", args: []string{"apply", planPath, planPath}},
		{name: "missing file", args: []string{"apply", "missing.tfplan"}},
		{name: "directory", args: []string{"apply", chdir}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := validateArgs(test.args, chdir, true); err == nil {
				t.Fatal("command accepted without exactly one saved plan file")
			}
		})
	}
}

func TestValidateArgsPreservesStandardCommands(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"apply"},
		{"plan", "-target=module.dns"},
		{"apply", "-auto-approve"},
		{"apply", "missing.tfplan"},
	} {
		if err := validateArgs(args, "missing-directory", false); err != nil {
			t.Fatalf("command %q changed with guard disabled: %v", args, err)
		}
	}
}

func TestValidateEnvironmentRejectsInjectedArguments(t *testing.T) {
	for _, name := range []string{
		"TF_CLI_ARGS",
		"TF_CLI_ARGS_apply",
		"TF_CLI_ARGS_init",
		"TF_CLI_ARGS_future_operation",
	} {
		t.Run(name, func(t *testing.T) {
			// This flag can consume the supplied plan path and leave bare apply.
			err := validateEnvironment([]string{name + "=-backup"}, true)
			if err == nil {
				t.Fatal("injected arguments accepted")
			}
			want := "--require-saved-plan rejects nonempty " + name
			if err.Error() != want {
				t.Fatal("diagnostic must identify only the variable name, not its value")
			}
		})
	}
}

func TestValidateEnvironmentAllowsEmptyAndUnrelatedVariables(t *testing.T) {
	for _, environ := range [][]string{
		nil,
		{"TF_CLI_ARGS=", "TF_CLI_ARGS_apply=", "TF_CLI_ARGS_future_operation="},
		{"TF_CLI_ARGSOTHER=-backup", "TF_CLI_ARG=-backup", "OTHER_TF_CLI_ARGS=-backup"},
	} {
		if err := validateEnvironment(environ, true); err != nil {
			t.Fatalf("allowed environment rejected: %v", err)
		}
	}
	if err := validateEnvironment([]string{"TF_CLI_ARGS_apply= "}, true); err == nil {
		t.Fatal("nonempty whitespace value accepted")
	}
}

func TestValidateEnvironmentPreservesStandardCommands(t *testing.T) {
	environ := []string{"TF_CLI_ARGS=-input=false", "TF_CLI_ARGS_apply=-backup"}
	if err := validateEnvironment(environ, false); err != nil {
		t.Fatalf("environment changed with guard disabled: %v", err)
	}
}
