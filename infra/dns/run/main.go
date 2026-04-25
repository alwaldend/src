package main

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"text/template"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var dnscontrolFlag = flag.String("dnscontrol", "tools/dnscontrol/dnscontrol.script.sh", "Dnscontrol binary")
var configFlag = flag.String("config", "infra/dns/dnsconfig.js", "Config file")
var alConfigFlag = flag.String("al_config", "infra/dns/al_/al.json", "Al config")
var credsTemplateFlag = flag.String("creds_template", "infra/dns/creds.json.tpl", "Creds template file")


func main() {
	ctx := context.Background()
	flag.Parse()
	config, err := al.LoadConfigs(ctx, *alConfigFlag)
	al.Check(err)
	secrets, err := al.FetchSecrets(ctx, config)
	al.Check(err)
    tmplText, err := os.ReadFile(*credsTemplateFlag)
    al.Check(err)
	tmpl, err := template.New("creds").Parse(string(tmplText))
    al.Check(err)
	creds, err := os.CreateTemp("", "creds.json")
    al.Check(err)
    defer os.Remove(creds.Name())
    err = tmpl.Execute(creds, map[string]any{"Secrets": secrets})
	al.Check(err)
	cmd := exec.Command(*dnscontrolFlag, append(flag.Args(), "--config", *configFlag, "--creds", creds.Name())...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	al.Check(al.SetRunfilesInfo(cmd))
	al.Check(cmd.Run())
}
