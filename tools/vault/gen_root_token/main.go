package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"golang.org/x/term"
)

var vaultFlag = flag.String("vault", "", "Vault binary")
var pgpKeyFlag = flag.String("pgp_key", "", "Public pgp key (in base64)")
var nonceFlag = flag.String("nonce", "", "Init nonce")

func main() {
	flag.Parse()
	al.Must(os.Stderr.WriteString("Enter the encrypted unseal key: "))
	unsealEncrypted := strings.TrimSpace(string(al.Must(term.ReadPassword(int(syscall.Stdin)))))
	unsealDecoded := al.Must(base64.StdEncoding.DecodeString(unsealEncrypted))
	unsealDecodedBuf := bytes.NewBuffer(unsealDecoded)
	var unsealDecrypted bytes.Buffer
	al.Check(al.RunCommand(al.CommandArgs{
		Name:   "gpg",
		Args:   []string{"--decrypt"},
		Stdin:  unsealDecodedBuf,
		Stdout: &unsealDecrypted,
	}))
	nonce := *nonceFlag
	if nonce == "" {
		var initOutput bytes.Buffer
		al.Check(al.RunCommand(al.CommandArgs{
			Name: *vaultFlag,
			Args: []string{
				"operator",
				"generate-root",
				"-init",
				fmt.Sprintf("-pgp-key=%s", *pgpKeyFlag),
				"-format",
				"json",
			},
			Stdout: &initOutput,
		}))
		initVal := struct {
			Nonce string `json:"nonce"`
		}{}
		al.Check(json.Unmarshal(initOutput.Bytes(), &initVal))
		nonce = initVal.Nonce
	}
	var genOutput bytes.Buffer
	al.Check(al.RunCommand(al.CommandArgs{
		Name: *vaultFlag,
		Args: []string{
			"operator",
			"generate-root",
			fmt.Sprintf("-nonce=%s", nonce),
			"-format",
			"json",
			"-",
		},
		Stdin:  &unsealDecrypted,
		Stdout: &genOutput,
	}))
	genVal := struct {
		EncodedToken string `json:"encoded_token"`
	}{}
	al.Check(json.Unmarshal(genOutput.Bytes(), &genVal))
	tokenEnc := al.Must(base64.StdEncoding.DecodeString(genVal.EncodedToken))
	tokenEncBuffer := bytes.NewBuffer(tokenEnc)
	var token bytes.Buffer
	al.Check(al.RunCommand(al.CommandArgs{
		Name:   "gpg",
		Args:   []string{"--decrypt"},
		Stdin:  tokenEncBuffer,
		Stdout: &token,
	}))
	os.WriteFile(filepath.Join(al.Must(os.UserHomeDir()), ".vault-token"), token.Bytes(), 0600)
}
