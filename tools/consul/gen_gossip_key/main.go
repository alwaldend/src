package main

import (
	"bytes"
	"flag"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

var vault = flag.String("vault", "", "Vault path")
var consul = flag.String("consul", "", "Consul path")
var secret = flag.String("secret", "", "Secret path")
var secretMount = flag.String("secret_mount", "", "Secret mount")

func main() {
	flag.Parse()
	var keyBuffer bytes.Buffer
	al.Check(al.RunCommand(al.CommandArgs{
		Name:   *consul,
		Args:   []string{"keygen"},
		Stdout: &keyBuffer,
	}))
	key := strings.Trim(keyBuffer.String(), "\n")
	al.Check(al.RunCommand(al.CommandArgs{
		Name:  *vault,
		Args:  []string{"kv", "put", "-mount", *secretMount, *secret, "key=-"},
		Stdin: strings.NewReader(key),
	}))
}
