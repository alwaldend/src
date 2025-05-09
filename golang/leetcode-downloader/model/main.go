package model

import (
	"embed"
	"encoding/json"
	"fmt"
	"text/template"

	proto "git.alwaldend.com/proto/leetcode-downloader"
)

//go:embed config.json
var embedFS embed.FS

func DefaultConfig() (*proto.Config, error) {
	content, err := embedFS.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("could not open embedded config: %w", err)
	}
	config := &proto.Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, fmt.Errorf("could not parse embedded config: %w", err)
	}
	return config, nil
}

func SubmissionTemplate() (*template.Template, error) {
	templ, err := template.ParseFS(embedFS, "submission.txt.tmpl")
	if err != nil {
		return nil, fmt.Errorf("could not parse embedded template: %w", err)
	}
	return templ, nil
}
