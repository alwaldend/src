package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

var vaultFlag = flag.String("vault", "", "Vault binary")

func main() {
	flag.Parse()
	unsealKey := os.Getenv("UNSEAL_KEY")
	if unsealKey == "" {
		os.Stderr.WriteString("Enter the encrypted unseal key:\n")
		unsealKey = strings.TrimSpace(string(al.Must(term.ReadPassword(int(syscall.Stdin)))))
	}
	unsealDecoded := al.Must(base64.StdEncoding.DecodeString(unsealKey))
	unsealDecodedBuf := bytes.NewBuffer(unsealDecoded)
	var unsealDecryptedBuffer bytes.Buffer
	al.Check(al.RunCommand(al.CommandArgs{
		Name: "gpg",
		Args: []string{
			"--decrypt",
		},
		Stdin:  unsealDecodedBuf,
		Stdout: &unsealDecryptedBuffer,
	}))
	unsealDecrypted := unsealDecryptedBuffer.Bytes()
	var err error
	for _, address := range flag.Args() {
		os.Setenv("VAULT_ADDR", address)
		err = errors.Join(
			err,
			al.RunCommand(
				al.CommandArgs{
					Name: *vaultFlag,
					Args: []string{
						"write",
						"-address",
						address,
						"sys/unseal",
						"key=-",
					},
					Stdin: bytes.NewBuffer(unsealDecrypted),
				},
			),
		)
	}
	al.Check(err)
}
