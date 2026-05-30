package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"os"
	"strings"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"golang.org/x/term"
)

var vaultFlag = flag.String("vault", "", "Vault binary")

func main() {
	flag.Parse()
	al.Must(os.Stderr.WriteString("Enter the encrypted unseal key: "))
	unsealEncrypted := strings.TrimSpace(string(al.Must(term.ReadPassword(int(syscall.Stdin)))))
	unsealDecoded := al.Must(base64.StdEncoding.DecodeString(unsealEncrypted))
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
