package main

import (
	"encoding/json"
	"errors"
	"io"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var vaultFlag = flag.String("vault", "tools/vault/vault.script.sh", "Vault binary")
var vaultSecretFlag = flag.String("vault_secret", "secrets/data/cloudflare.com/dns_token", "Vault secret path")
var vaultApproleFlag = flag.String("vault_approle", "src_infra_dns", "Approle name")
var dnscontrolFlag = flag.String("dnscontrol", "tools/dnscontrol/dnscontrol.script.sh", "Dnscontrol binary")
var configFlag = flag.String("config", "infra/dns/dnsconfig.js", "Config file")
var credsTemplateFlag = flag.String("creds_template", "infra/dns/creds.json.tpl", "Creds template file")

func check(errs... error) {
    err := errors.Join(errs...)
    if err != nil {
        panic(err)
    }
}

func main() {
	flag.Parse()
    runfilesEnv, err := runfiles.Env()
    check(err)

    vaultApproleCmd := exec.Command(*vaultFlag, "write", "--force", "--format", "json", fmt.Sprintf("auth/approle/role/%s/secret-id", *vaultApproleFlag))
    vaultApproleCmd.Stderr = os.Stderr
    vaultApproleCmd.Env = append(os.Environ(), runfilesEnv...)
    vaultApproleString, err := vaultApproleCmd.Output()
    check(err)
	vaultApproleData := struct { Data struct { SecretId string `json:"secret_id"` } `json:"data"` }{}
	check(json.Unmarshal(vaultApproleString, &vaultApproleData))

    vaultApproleLoginCmd := exec.Command(*vaultFlag, "write", "--format", "json", "auth/approle/login", fmt.Sprintf("role_id=%s", *vaultApproleFlag), "secret_id=-")
    vaultApproleLoginCmd.Stderr = os.Stderr
	stdin, _ := vaultApproleLoginCmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, vaultApproleData.Data.SecretId)
	}()
    vaultApproleLoginCmd.Env = append(os.Environ(), runfilesEnv...)
    vaultApproleLoginString, err := vaultApproleLoginCmd.Output()
    check(err)
	vaultApproleLoginData := struct { Data struct { Token string `json:"token"` } `json:"data"` }{}
	check(json.Unmarshal(vaultApproleLoginString, &vaultApproleLoginData))

    vaultSecretCmd := exec.Command(*vaultFlag, "kv", "get", "--format", "json", *vaultSecretFlag)
    vaultSecretCmd.Stderr = os.Stderr
    vaultSecretCmd.Env = append(os.Environ(),runfilesEnv...)
	vaultSecretCmd.Env = append(vaultSecretCmd.Env, fmt.Sprintf("VAULT_TOKEN=%s", vaultApproleLoginData.Data.Token))
    vaultSecretString, err := vaultSecretCmd.Output()
    check(err)
	vaultSecretData := make(map[string]any)
	check(json.Unmarshal(vaultSecretString, &vaultSecretData))

    tmplText, err := os.ReadFile(*credsTemplateFlag)
    check(err)
	tmpl, err := template.New("creds").Parse(string(tmplText))
    check(err)
	creds, err := os.CreateTemp("", "creds.json")
    check(err)
    defer os.Remove(creds.Name())
    check(tmpl.Execute(creds, map[string]any{"Secret": vaultSecretData}))

    dnscontrolCmd := exec.Command(*dnscontrolFlag, append(flag.Args(), "--config", *configFlag, "--creds", creds.Name())...)
    dnscontrolCmd.Env = append(os.Environ(), runfilesEnv...)
    dnscontrolCmd.Stdout = os.Stdout
    dnscontrolCmd.Stderr = os.Stderr
    check(dnscontrolCmd.Run())
}
