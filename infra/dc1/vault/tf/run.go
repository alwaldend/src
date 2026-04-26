package main

import (
	"context"
	"flag"
	"fmt"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var terraformFlag = flag.String("terraform", "tools/terraform/terraform.script.sh", "Terraform binary")
var alConfigFlag = flag.String("al_config", "infra/dc1/vault/al_/al.json", "Al config")
var chDirFlag = flag.String("chdir", ".", "--chdir flag for terraform")

func main() {
	ctx := context.Background()
	flag.Parse()
	config := al.Must(al.LoadConfigs(ctx, *alConfigFlag))
	client := al.Must(al.VaultAuthDefault(ctx, config))
    vaultEnv := al.Must(al.VaultDefaultEnv(config, client))
	cmdInit := al.Must(al.Command(*terraformFlag, append([]string{fmt.Sprintf("-chdir=%s", *chDirFlag)}, "init")...))
	cmdInit.Env = append(cmdInit.Env, vaultEnv...)
	al.Check(cmdInit.Run())
	cmd := al.Must(al.Command(*terraformFlag, append([]string{fmt.Sprintf("-chdir=%s", *chDirFlag)}, flag.Args()...)...))
	cmd.Env = append(cmd.Env, vaultEnv...)
	al.Check(cmd.Run())
}
