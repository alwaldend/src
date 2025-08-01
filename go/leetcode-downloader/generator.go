package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"

	proto "git.alwaldend.com/src/go/leetcode-downloader/proto"
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
	storage := &proto.SubmissonsStorage{Submissions: []*proto.Submission{submission}}
	dataContent, err := toml.Marshal(storage)
	if err != nil {
		return err
	}
	dir, err := self.submissionDir(submission)
	if err != nil {
		return err
	}
	dir = filepath.Join(self.config.RootDir, dir)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}
	codePath := filepath.Join(dir, fmt.Sprintf("%s.%s", name, extension))
	dataPath := filepath.Join(dir, fmt.Sprintf("%s.%s", name, "toml"))
	if self.config.WriteCode {
		err = os.WriteFile(codePath, codeContent, 0o644)
		if err != nil {
			return fmt.Errorf("could not write '%s': %w", codePath, err)
		}
	}
	err = os.WriteFile(dataPath, dataContent, 0o644)
	if err != nil {
		return fmt.Errorf("could not write '%s': %w", dataPath, err)
	}
	return nil
}

func (self *Generator) submissionDir(submission *proto.Submission) (string, error) {
	config, err := self.submissionConfig(submission)
	if err != nil {
		return "", err
	}
	question := fmt.Sprintf("%06d", submission.QuestionId)
	res := filepath.Join(config.Dir, "questions", question, "submissions")
	return res, nil
}

func (self *Generator) submissionConfig(submission *proto.Submission) (*proto.SubmissionConfig, error) {
	for _, config := range self.config.Submissions {
		if slices.Contains(config.Types, submission.Lang) {
			return config, nil
		}
	}
	return nil, fmt.Errorf("could not find config for submission %v", submission)
}

func (self *Generator) submissionExtension(submission *proto.Submission) (string, error) {
	config, err := self.submissionConfig(submission)
	if err != nil {
		return "", err
	}
	return config.Extension, nil
}

func (self *Generator) submissionName(submission *proto.Submission) (string, error) {
	name := fmt.Sprintf("%012d", submission.Id)
	return name, nil
}
