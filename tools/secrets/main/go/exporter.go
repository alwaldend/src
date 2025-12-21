package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"git.alwaldend.com/alwaldend/src/tools/secrets/main/proto/contracts"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Exporter struct {
	backend        *Backend
	stdout, stderr io.Writer
}

func NewExporter(backend *Backend, stdout, stderr io.Writer) *Exporter {
	return &Exporter{
		backend: backend,
		stdout:  stdout,
		stderr:  stderr,
	}
}

type ExporterOpts struct {
	Secrets []string
	Output  []string
}

func (self *Exporter) Export(opts *ExporterOpts) error {
	for _, path := range opts.Secrets {
		secret, err := self.readSecret(path)
		if err != nil {
			return fmt.Errorf("could not read secret %s: %w", path, err)
		}
		for _, output := range opts.Output {
			err = self.exportSecret(secret, output)
			if err != nil {
				return fmt.Errorf("could not export secret to %s: %w", output, err)
			}
		}
	}
	return nil
}

func (self *Exporter) readSecret(path string) (*contracts.Secret, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", path, err)
	}
	secret := &contracts.Secret{}
	err = protojson.Unmarshal(content, secret)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal secret: %w", err)
	}
	return secret, nil
}

func (self *Exporter) exportSecret(secret *contracts.Secret, output string) error {
	res := []string{"#!/usr/bin/env sh"}
	data := &contracts.BackendData{}
	for _, backend := range secret.Backends {
		curData, err := self.backend.Pull(backend)
		if err != nil {
			return fmt.Errorf("could not pull from backend %v: %w", backend, err)
		}
		proto.Merge(data, curData)
	}
	for key, value := range secret.Env {
		tmpl, err := template.New("").Funcs(
			template.FuncMap{
				"secret": func(path string) (string, error) {
					val, ok := data.Secrets[path]
					if !ok {
						return "", fmt.Errorf("missing secret: %s", path)
					}
					return val, nil
				},
			},
		).Parse(value)
		if err != nil {
			return fmt.Errorf("could not create template for %s: %w", key, err)
		}
		var templateBuffer bytes.Buffer
		err = tmpl.Execute(&templateBuffer, data)
		if err != nil {
			return fmt.Errorf("could not execute template '%s': %w", value, err)
		}
		res = append(res, fmt.Sprintf("export '%s=%s'", key, templateBuffer.String()))
	}
	err := os.WriteFile(output, []byte(strings.Join(res, "\n")), 0o600)
	if err != nil {
		return fmt.Errorf("could not write to file %s: %w", output, err)
	}
	return nil
}
