package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

var (
	terraformFlag        = flag.String("terraform", "", "Terraform binary")
	chDirFlag            = flag.String("chdir", ".", "--chdir flag for terraform")
	directFlag           = flag.Bool("direct", false, "If set, just run the command")
	requireSavedPlanFlag = flag.Bool("require-saved-plan", false, "Require apply with exactly one existing saved plan file")
	logger               = log.New(os.Stderr, "com.alwaldend.src.tools.terraform.runner ", log.Flags())
)

func validateArgs(args []string, chdir string, requireSavedPlan bool) error {
	if !requireSavedPlan {
		return nil
	}
	if len(args) != 2 || args[0] != "apply" || args[1] == "" || strings.HasPrefix(args[1], "-") {
		return errors.New("--require-saved-plan requires apply <saved-plan-file> with no other arguments")
	}
	planPath := args[1]
	if !filepath.IsAbs(planPath) {
		planPath = filepath.Join(chdir, planPath)
	}
	info, err := os.Stat(planPath)
	if err != nil {
		return fmt.Errorf("could not inspect saved plan file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return errors.New("saved plan path must refer to a regular file")
	}
	return nil
}

func validateEnvironment(environ []string, requireSavedPlan bool) error {
	if !requireSavedPlan {
		return nil
	}
	for _, env := range environ {
		name, value, _ := strings.Cut(env, "=")
		if value != "" && (name == "TF_CLI_ARGS" || strings.HasPrefix(name, "TF_CLI_ARGS_")) {
			return fmt.Errorf("--require-saved-plan rejects nonempty %s", name)
		}
	}
	return nil
}

func run() int {
	flag.Parse()
	args := flag.Args()
	if err := validateArgs(args, *chDirFlag, *requireSavedPlanFlag); err != nil {
		logger.Printf("Invalid command: %s\n", err)
		return 2
	}
	environ := os.Environ()
	if err := validateEnvironment(environ, *requireSavedPlanFlag); err != nil {
		logger.Printf("Invalid environment: %s\n", err)
		return 2
	}
	commonArgs := []string{fmt.Sprintf("-chdir=%s", *chDirFlag)}
	backendArgs := []string{}
	for _, env := range environ {
		split := strings.SplitN(env, "=", 2)
		if len(split) != 2 {
			continue
		}
		if strings.HasPrefix(split[0], "AL_TF_BACKEND_CONFIG") {
			backendArgs = append(backendArgs, "--backend-config", split[1])
		}
	}
	if !*directFlag {
		cmdInit, err := al.Command(al.CommandArgs{
			Name: *terraformFlag,
			Args: slices.Concat(commonArgs, []string{"init"}, backendArgs),
		})
		if err != nil {
			logger.Printf("Could not create the init command: %s\n", err)
			return 2
		}
		cmdInit.Stdout = os.Stderr
		if err := cmdInit.Run(); err != nil {
			logger.Printf("Could not execute the init command: %s\n", err)
			return 3
		}
	}
	if len(args) == 0 || args[0] != "init" {
		backendArgs = nil
	}
	cmd, err := al.Command(al.CommandArgs{
		Name: *terraformFlag,
		Args: slices.Concat(commonArgs, args, backendArgs),
	})
	if err != nil {
		logger.Printf("Could not create the main command: %s\n", err)
		return 4
	}
	if err := cmd.Run(); err != nil {
		logger.Printf("Could not run the main command: %s\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
