package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"git.alwaldend.com/alwaldend/src/tools/secrets/main/proto/contracts"
	"google.golang.org/protobuf/encoding/protojson"
)

type Editor struct {
	stdout, stderr io.Writer
}

func NewEditor(
	stdout, stderr io.Writer,
) *Editor {
	return &Editor{
		stdout: stdout,
		stderr: stderr,
	}
}

type EditorOpts struct {
	Backend []string
	Output  []string
}

func (self *Editor) Edit(opts *EditorOpts) error {
	for _, backendPath := range opts.Backend {
		err := self.updateBackend(backendPath, opts.Output)
		if err != nil {
			return fmt.Errorf("could not edit backend: %w", err)
		}
	}
	return nil
}

func (self *Editor) updateBackend(backendPath string, outputPath []string) error {
	backend := &contracts.Backend{}
	backendContent, err := os.ReadFile(backendPath)
	if err != nil {
		return fmt.Errorf("could not read backend: %w", err)
	}
	err = protojson.Unmarshal(backendContent, backend)
	if err != nil {
		return fmt.Errorf("could not unmarshal backend: %w", err)
	}
	err = self.updateEncryped(backend)
	if err != nil {
		return fmt.Errorf("could not update encrypted part: %w", err)
	}
	res, err := protojson.MarshalOptions{Indent: "    "}.Marshal(backend)
	if err != nil {
		return fmt.Errorf("could not marshal: %w", err)
	}
	for _, output := range outputPath {
		err = os.WriteFile(output, res, os.FileMode(0o600))
		if err != nil {
			return fmt.Errorf("could not write updated backend: %w", err)
		}
	}
	return nil
}

func (self *Editor) updateEncryped(backend *contracts.Backend) error {
	systemdCredsFile := backend.SystemdCreds.Encrypted
	tempInput, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(tempInput.Name())
	tempOutput, err := os.CreateTemp("", "*.secret.json")
	if err != nil {
		return fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(tempOutput.Name())

	if systemdCredsFile != "" {
		err = os.WriteFile(tempInput.Name(), []byte(systemdCredsFile), 0x555)
		if err != nil {
			return fmt.Errorf("could not write systemd-creds file %s: %w", tempInput.Name(), err)
		}
		var cmdDecrypt *exec.Cmd
		if backend.SystemdCreds.System {
			cmdDecrypt = exec.Command("systemd-creds", "decrypt", "--name", backend.SystemdCreds.Name, tempInput.Name(), tempOutput.Name())
		} else {
			cmdDecrypt = exec.Command("systemd-creds", "decrypt", "--user", "--name", backend.SystemdCreds.Name, tempInput.Name(), tempOutput.Name())
		}
		cmdDecrypt.Stderr = self.stderr
		cmdDecrypt.Stdout = self.stdout
		err = cmdDecrypt.Run()
		if err != nil {
			return fmt.Errorf("could not decrypt input %s: %w", tempInput.Name(), err)
		}
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmdEdit := exec.Command(editor, tempOutput.Name())
	cmdEdit.Stderr = self.stderr
	cmdEdit.Stdout = self.stdout
	err = cmdEdit.Run()
	if err != nil {
		return fmt.Errorf("could not edit temp output %s: %w", tempOutput.Name(), err)
	}

	var cmdEncrypt *exec.Cmd
	if backend.SystemdCreds.System {
		cmdEncrypt = exec.Command("systemd-creds", "encrypt", "--name", backend.SystemdCreds.Name, tempOutput.Name(), tempInput.Name())
	} else {
		cmdEncrypt = exec.Command("systemd-creds", "encrypt", "--user", "--name", backend.SystemdCreds.Name, tempOutput.Name(), tempInput.Name())
	}
	cmdEncrypt.Stderr = self.stderr
	cmdEncrypt.Stdout = self.stdout
	err = cmdEncrypt.Run()
	if err != nil {
		return fmt.Errorf("could not encrypt output %s: %w", tempOutput.Name(), err)
	}
	newCreds, err := os.ReadFile(tempInput.Name())
	if err != nil {
		return fmt.Errorf("could not read temp input %s: %w", tempInput.Name(), err)
	}
	backend.SystemdCreds.Encrypted = string(newCreds)
	return nil
}
