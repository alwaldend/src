package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"os/exec"
	"text/template"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var vaultFlag = flag.String("vault", "tools/vault/vault.script.sh", "Vault binary")
var vaultSecretFlag = flag.String("vault_secret", "sre/data/infra/dns", "Vault secret path")
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
    tmplText, err := os.ReadFile(*credsTemplateFlag)
    check(err)
	tmpl, err := template.New("creds").Parse(string(tmplText))
    check(err)
	creds, err := os.CreateTemp("", "creds.json")
    check(err)
    defer os.Remove(creds.Name())
    vaultCmd := exec.Command(*vaultFlag, "kv", "get", "--format", "json", *vaultSecretFlag)
    vaultCmd.Stderr = os.Stderr
    vaultCmd.Env = append(os.Environ(), runfilesEnv...)
    secretString, err := vaultCmd.Output()
    check(err)
    secret := make(map[string]any)
    check(json.Unmarshal(secretString, &secret))
    check(tmpl.Execute(creds, map[string]any{"Secret": secret}))
    dnscontrolCmd := exec.Command(*dnscontrolFlag, append(flag.Args(), "--config", *configFlag, "--creds", creds.Name())...)
    dnscontrolCmd.Env = append(os.Environ(), runfilesEnv...)
    dnscontrolCmd.Stdout = os.Stdout
    dnscontrolCmd.Stderr = os.Stderr
    check(dnscontrolCmd.Run())
}
