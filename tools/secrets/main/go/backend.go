package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"git.alwaldend.com/alwaldend/src/tools/secrets/main/proto/contracts"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Backend struct {
	stdout, stderr io.Writer
}

func NewBackend(
	stdout, stderr io.Writer,
) *Backend {
	return &Backend{
		stdout: stdout,
		stderr: stderr,
	}
}

func (self *Backend) ReadFile(paths ...string) (*contracts.Backend, error) {
	res := &contracts.Backend{}
	for _, path := range paths {
		path = os.ExpandEnv(path)
		backend := &contracts.Backend{}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not read file %s: %w", paths, err)
		}
		err = protojson.Unmarshal(content, backend)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal backend %s: %w", paths, err)
		}
		proto.Merge(res, backend)
	}
	return res, nil
}

func (self *Backend) WriteFile(backend *contracts.Backend, outputs ...string) error {
	res, err := protojson.MarshalOptions{Indent: "    "}.Marshal(backend)
	if err != nil {
		return fmt.Errorf("could not marshal backend: %w", err)
	}
	for _, output := range outputs {
		output = os.ExpandEnv(output)
		err = os.WriteFile(output, res, os.FileMode(0o600))
		if err != nil {
			return fmt.Errorf("could not write backend to %s: %w", output, err)
		}
	}
	return nil
}

func (self *Backend) Pull(backend *contracts.Backend) (*contracts.BackendData, error) {
	systemdCreds := backend.SystemdCreds
	if systemdCreds == nil {
		return nil, fmt.Errorf("missing backend configuration")
	}
	path := os.ExpandEnv(systemdCreds.EncryptedFile)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not stat encrypted systemd-creds file %s: %w", path, err)
	}
	if stat.Size() == 0 {
		return &contracts.BackendData{}, nil
	}
	tempOutput, err := os.CreateTemp("", "*.backend_data.json")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary output file: %w", err)
	}
	defer os.Remove(tempOutput.Name())
	err = self.decryptSystemdCreds(
		systemdCreds.Name, path,
		tempOutput.Name(), systemdCreds.System,
	)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt systemd creds: %w", err)
	}
	res, err := self.readDataFile(tempOutput.Name())
	if err != nil {
		return nil, fmt.Errorf("could not parse decrypted data file: %w", err)
	}
	return res, nil
}

func (self *Backend) Push(backend *contracts.Backend, data *contracts.BackendData) error {
	systemdCreds := backend.SystemdCreds
	if systemdCreds == nil {
		return fmt.Errorf("missing backend configuration")
	}
	tempOutput, err := os.CreateTemp("", "*.backend_data.json")
	if err != nil {
		return fmt.Errorf("could not create temporary output file: %w", err)
	}
	defer os.Remove(tempOutput.Name())
	err = self.writeDataFile(data, tempOutput.Name())
	if err != nil {
		return fmt.Errorf("could not write data file: %w", err)
	}
	err = self.encryptSystemdCreds(
		systemdCreds.Name, tempOutput.Name(),
		os.ExpandEnv(systemdCreds.EncryptedFile), systemdCreds.System,
	)
	if err != nil {
		return fmt.Errorf("could not encrypt the backend data: %w", err)
	}
	return nil
}

func (self *Backend) Edit(backend *contracts.Backend) error {
	data, err := self.Pull(backend)
	if err != nil {
		return fmt.Errorf("could not edit backend: %w", err)
	}
	newData, err := self.editData(data)
	if err != nil {
		return fmt.Errorf("could not edit backend data: %w", err)
	}
	err = self.Push(backend, newData)
	if err != nil {
		return fmt.Errorf("could not push edited backend data: %w", err)
	}
	return nil
}

func (self *Backend) editData(data *contracts.BackendData) (*contracts.BackendData, error) {
	dataFile, err := os.CreateTemp("", "*.backend_data.json")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(dataFile.Name())
	err = self.writeDataFile(data, dataFile.Name())
	if err != nil {
		return nil, fmt.Errorf("could not write data to temporary file: %w", err)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmdEdit := exec.Command(editor, dataFile.Name())
	cmdEdit.Stderr = self.stderr
	cmdEdit.Stdout = self.stdout
	err = cmdEdit.Run()
	if err != nil {
		return nil, fmt.Errorf("could not edit temp data file: %w", err)
	}
	updatedData, err := self.readDataFile(dataFile.Name())
	if err != nil {
		return nil, fmt.Errorf("could not parse updated data file: %w", err)
	}
	return updatedData, nil
}

func (self *Backend) writeDataFile(data *contracts.BackendData, path string) error {
	dataMarshaled, err := protojson.MarshalOptions{Indent: "    "}.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal backend data: %w", err)
	}
	err = os.WriteFile(path, dataMarshaled, 0o600)
	if err != nil {
		return fmt.Errorf("could not write temp data file: %w", err)
	}
	return nil
}

func (self *Backend) readDataFile(path string) (*contracts.BackendData, error) {
	res := &contracts.BackendData{}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read backend data file %s: %w", path, err)
	}
	err = protojson.Unmarshal(content, res)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal backend data file %s: %w", path, err)
	}
	return res, nil
}

func (self *Backend) decryptSystemdCreds(name, input, output string, system bool) error {
	args := []string{"decrypt"}
	if !system {
		args = append(args, "--user")
	}
	args = append(args, input, output)
	return self.runSystemdCredsCmd(args)
}

func (self *Backend) encryptSystemdCreds(name, input, output string, system bool) error {
	args := []string{"encrypt"}
	if !system {
		args = append(args, "--user")
	}
	args = append(args, input, output)
	return self.runSystemdCredsCmd(args)
}

func (self *Backend) runSystemdCredsCmd(args []string) error {
	cmd := exec.Command("systemd-creds", args...)
	cmd.Stderr = self.stderr
	cmd.Stdout = self.stdout
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run cmd %s: %w", args, err)
	}
	return nil
}
