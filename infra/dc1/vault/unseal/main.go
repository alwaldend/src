package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"os"
	"strings"
	"syscall"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
	"golang.org/x/term"
)

var vaultFlag = flag.String("vault", "", "Vault binary")

func main() {
	flag.Parse()
	al.Must(os.Stderr.WriteString("Enter the encrypted unseal key: "))
	unsealEncrypted := strings.TrimSpace(string(al.Must(term.ReadPassword(int(syscall.Stdin)))))
	unsealDecoded := al.Must(base64.StdEncoding.DecodeString(unsealEncrypted))
	unsealDecodedBuf := bytes.NewBuffer(unsealDecoded)
	var unsealDecrypted bytes.Buffer
    al.Check(al.RunCommand(al.CommandArgs{
        Name: "gpg",
        Args: []string{
            "--decrypt",
        },
		Stdin: unsealDecodedBuf,
		Stdout: &unsealDecrypted,
    }))
    al.Check(al.RunCommand(al.CommandArgs{
        Name: *vaultFlag,
        Args: []string{
            "write",
			"sys/unseal",
			"key=-",
        },
		Stdin: &unsealDecrypted,
    }))
}
