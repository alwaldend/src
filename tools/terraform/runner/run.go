package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var terraformFlag = flag.String("terraform", "", "Terraform binary")
var chDirFlag = flag.String("chdir", ".", "--chdir flag for terraform")
var directFlag = flag.Bool("direct", false, "If set, just run the command")

func main() {
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
		cmdInit := al.Must(al.Command(al.CommandArgs{Name: *terraformFlag, Args: slices.Concat(commonArgs, []string{"init"}, backendArgs)}))
		al.Check(cmdInit.Run())
	}
	if len(args) == 0 || args[0] != "init" {
		backendArgs = nil
	}
	cmd := al.Must(al.Command(al.CommandArgs{Name: *terraformFlag, Args: slices.Concat(commonArgs, args, backendArgs)}))
	al.Check(cmd.Run())
}
