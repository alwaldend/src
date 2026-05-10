package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var vaultFlag = flag.String("vault", "", "Vault binary")
var mountFlag = flag.String("path", "pki/ica_clients", "Mount path")
var roleFlag = flag.String("role", "pki_ica_clients_users", "Role name")
var outputDirFlag = flag.String("output_dir", "", "Output directory")
var usernameFlag = flag.String("user", "", "Username")
var hostFlag = flag.String("host", "", "Hostname")
var noCreateFlag = flag.Bool("no_create", false, "If set, do not create a new cert")
var ttlFlag = flag.String("ttl", "28927206", "TTL")

func main() {
	flag.Parse()
	if (*usernameFlag == "" && *hostFlag == "") || *mountFlag == "" || *roleFlag == "" || *outputDirFlag == "" || *ttlFlag == "" {
		os.Stderr.WriteString("missing required flags\n")
	}
	target := *usernameFlag
	commonName := fmt.Sprintf("%s.users.alwaldend.com", *usernameFlag)
	if *hostFlag != "" {
		target = *hostFlag
		if *usernameFlag != "" {
			commonName = fmt.Sprintf("%s.%s", *hostFlag, commonName)
		} else {
			commonName = fmt.Sprintf("%s.infra.alwaldend.com", *hostFlag)
		}
	}
	jsonPath := filepath.Join(*outputDirFlag, fmt.Sprintf("%s.json", target))
	pubPath := filepath.Join(*outputDirFlag, fmt.Sprintf("%s.pub", target))
	chainPath := filepath.Join(*outputDirFlag, fmt.Sprintf("%s_chain.pub", target))
	keyPath := filepath.Join(*outputDirFlag, fmt.Sprintf("%s.key", target))
	combinedPath := filepath.Join(*outputDirFlag, fmt.Sprintf("%s.pem", target))
	pkcs12Path := filepath.Join(*outputDirFlag, fmt.Sprintf("%s.pfx", target))
	al.Check(os.MkdirAll(*outputDirFlag, 0700))
	if !*noCreateFlag {
		output := al.Must(os.OpenFile(jsonPath, os.O_WRONLY|os.O_CREATE, 0600))
		defer output.Close()
		al.Check(al.RunCommand(al.CommandArgs{
			Name: *vaultFlag,
			Args: append(
				[]string{
					"write",
					"-format=json",
					fmt.Sprintf("%s/issue/%s", *mountFlag, *roleFlag),
					fmt.Sprintf("common_name=%s", commonName),
					fmt.Sprintf("ttl=%s", *ttlFlag),
				},
				flag.Args()...,
			),
			Stdout: output,
		}))
	}
	data := struct {
		Data struct {
			Certificate string   `json:"certificate"`
			CaChain     []string `json:"ca_chain"`
			PrivateKey  string   `json:"private_key"`
		} `json:"data"`
	}{}
	al.Check(json.Unmarshal(al.Must(os.ReadFile(jsonPath)), &data))
	combined := []string{data.Data.Certificate}
	combined = append(combined, data.Data.CaChain...)
	combined = append(combined, data.Data.PrivateKey)
	combined = append(combined, "")
	al.Check(os.WriteFile(pubPath, []byte(data.Data.Certificate), 0600))
	al.Check(os.WriteFile(keyPath, []byte(data.Data.PrivateKey), 0600))
	al.Check(os.WriteFile(chainPath, []byte(strings.Join(data.Data.CaChain, "\n")), 0600))
	al.Check(os.WriteFile(combinedPath, []byte(strings.Join(combined, "\n")), 0600))
	al.Check(al.RunCommand(al.CommandArgs{
		Name: "openssl",
		Args: []string{
			"pkcs12",
			"-export",
			"-certfile",
			chainPath,
			"-out",
			pkcs12Path,
			"-in",
			pubPath,
			"-inkey",
			keyPath,
		},
	}))
}
