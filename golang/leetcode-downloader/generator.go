package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	proto "git.alwaldend.com/src/golang/leetcode-downloader/proto"
)

type Generator struct {
	config *proto.Config
	ctx    context.Context
	logger *slog.Logger
}

func NewGenerator(config *proto.Config, ctx context.Context, logger *slog.Logger) (*Generator, error) {
	generator := &Generator{
		config: config,
		ctx:    ctx,
		logger: logger,
	}
	return generator, nil
}

func (g *Generator) Generate(submissions []*proto.Submission) error {
	count := 0
	for _, submission := range submissions {
		if submission.StatusDisplay != "Accepted" {
			continue
		}
		if err := g.write(submission); err != nil {
			return err
		}
		count += 1
	}
	g.logger.Info(fmt.Sprintf("wrote %d submissions to %s", count, g.config.RootDir))
	return nil
}

func (self *Generator) write(submission *proto.Submission) error {
	name, err := self.submissionName(submission)
	if err != nil {
		return err
	}
	extension, err := self.submissionExtension(submission)
	if err != nil {
		return err
	}
	codeContent := []byte(submission.Code)
	dataContent, err := json.MarshalIndent(submission, "", "    ")
	if err != nil {
		return err
	}
	dir, err := self.submissionDir(submission)
	if err != nil {
		return err
	}
	dir = filepath.Join(self.config.RootDir, dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	codePath := filepath.Join(dir, fmt.Sprintf("%s.%s", name, extension))
	dataPath := filepath.Join(dir, fmt.Sprintf("%s.%s", name, "json"))
	err = os.WriteFile(codePath, codeContent, 0644)
	if err != nil {
		return fmt.Errorf("could not write '%s': %w", codePath, err)
	}
	err = os.WriteFile(dataPath, dataContent, 0644)
	if err != nil {
		return fmt.Errorf("could not write '%s': %w", codePath, err)
	}
	return nil
}

func (g *Generator) submissionDir(submission *proto.Submission) (string, error) {
	config, err := g.submissionConfig(submission)
	if err != nil {
		return "", nil
	}
	return config.Dir, nil
}

func (g *Generator) submissionConfig(submission *proto.Submission) (*proto.SubmissionConfig, error) {
	for _, config := range g.config.Submissions {
		if slices.Contains(config.Types, submission.Lang) {
			return config, nil
		}
	}
	return nil, fmt.Errorf("could not find config for submission %v", submission)
}

func (g *Generator) submissionExtension(submission *proto.Submission) (string, error) {
	config, err := g.submissionConfig(submission)
	if err != nil {
		return "", err
	}
	return config.Extension, nil
}

func (g *Generator) submissionName(submission *proto.Submission) (string, error) {
	name := fmt.Sprintf("%05d-%05d-%s", submission.QuestionId, submission.Id, submission.TitleSlug)
	return name, nil
}
