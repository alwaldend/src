package main

import (
	"flag"
	"os"
	"os/exec"
	"text/template"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var dnscontrolFlag = flag.String("dnscontrol", "", "Dnscontrol binary")
var configFlag = flag.String("config", "", "Config file")
var credsTemplateFlag = flag.String("creds_template", "", "Creds template file")

func main() {
	flag.Parse()
	tmplText := al.Must(os.ReadFile(*credsTemplateFlag))
	tmpl := al.Must(template.New("creds").Funcs(template.FuncMap{"env": os.Getenv}).Parse(string(tmplText)))
	creds := al.Must(os.CreateTemp("", "creds.json"))
	defer os.Remove(creds.Name())
	al.Check(tmpl.Execute(creds, nil))
	cmd := exec.Command(*dnscontrolFlag, append(flag.Args(), "--config", *configFlag, "--creds", creds.Name())...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	al.Check(al.SetRunfilesEnv(cmd))
	al.Check(cmd.Run())
}
