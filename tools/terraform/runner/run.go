package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

var terraformFlag = flag.String("terraform", "", "Terraform binary")
var chDirFlag = flag.String("chdir", ".", "--chdir flag for terraform")
var directFlag = flag.Bool("direct", false, "If set, just run the command")
var logger = log.New(os.Stderr, "com.alwaldend.src.tools.terraform.runner ", log.Flags())

func run() int {
	flag.Parse()
	commonArgs := []string{fmt.Sprintf("-chdir=%s", *chDirFlag)}
	backendArgs := []string{}
	args := flag.Args()
	for _, env := range os.Environ() {
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
