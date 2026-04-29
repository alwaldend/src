package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var terraformFlag = flag.String("terraform", "tools/terraform/terraform.script.sh", "Terraform binary")
var vaultFlag = flag.String("vault", "tools/vault/vault.script.sh", "Vault binary")
var alConfigFlag = flag.String("al_config", "infra/yandex_cloud/org1/al_/al.json", "Al config")
var chDirFlag = flag.String("chdir", ".", "--chdir flag for terraform")
var isVaultFlag = flag.Bool("run_vault", false, "If set, run vault")
var isTfFlag = flag.Bool("run_tf", false, "If set, run just tf")

func main() {
	ctx := context.Background()
	flag.Parse()
	config := al.Must(al.LoadConfigs(ctx, *alConfigFlag))
	client := al.Must(al.VaultAuthDefault(ctx, config))
    vaultEnv := al.Must(al.VaultDefaultEnv(config, client))
	if (*isVaultFlag) {
		cmdVault := al.Must(al.Command(*vaultFlag, flag.Args()...))
		cmdVault.Env = append(cmdVault.Env, vaultEnv...)
		al.Check(cmdVault.Run())
	} else {
		secrets := al.Must(al.VaultFetchSecrets(ctx, config))
		serviceAccountFile := al.Must(os.CreateTemp("", "service_account.json"))
		defer os.Remove(serviceAccountFile.Name())
		al.Must(serviceAccountFile.WriteString(secrets["tf_auth"]["service_account_key"].(string)))
		ycEnv := []string{
			fmt.Sprintf("YC_CLOUD_ID=%s", secrets["tf_auth"]["cloud_id"]),
			fmt.Sprintf("YC_FOLDER_ID=%s", secrets["tf_auth"]["folder_id"]),
			fmt.Sprintf("YC_SERVICE_ACCOUNT_KEY_FILE=%s", serviceAccountFile.Name()),
		}
		if (!*isTfFlag) {
			cmdInit := al.Must(al.Command(*terraformFlag, append([]string{fmt.Sprintf("-chdir=%s", *chDirFlag)}, "init")...))
			cmdInit.Env = append(cmdInit.Env, vaultEnv...)
			cmdInit.Env = append(cmdInit.Env, ycEnv...)
			al.Check(cmdInit.Run())
		}
		cmd := al.Must(al.Command(*terraformFlag, append([]string{fmt.Sprintf("-chdir=%s", *chDirFlag)}, flag.Args()...)...))
		cmd.Env = append(cmd.Env, vaultEnv...)
		cmd.Env = append(cmd.Env, ycEnv...)
		al.Check(cmd.Run())
	}
}
