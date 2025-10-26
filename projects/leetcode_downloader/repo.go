package main

import (
	"context"
	"fmt"
	"log/slog"

	proto "git.alwaldend.com/src/contracts/leetcode_downloader"
)

type Repo struct {
	config *proto.Config
	ctx    context.Context
	logger *slog.Logger
}

func NewRepo(config *proto.Config, ctx context.Context, logger *slog.Logger) *Repo {
	return &Repo{config: config, ctx: ctx, logger: logger}
}

func (r *Repo) Submissions() ([]*proto.Submission, error) {
	submissions, err := r.readStorageFile()
	if err != nil {
		return nil, err
	}
	return submissions.Submissions, nil
}

func (r *Repo) AddSubmissions(submissions []*proto.Submission) error {
	curSubmissions, err := r.readStorageFile()
	if err != nil {
		return err
	}
	r.addSubmissions(curSubmissions, submissions)
	return r.writeStorageFile(curSubmissions)
}

func (r *Repo) addSubmissions(curSubmissions *proto.SubmissonsStorage, newSubmissions []*proto.Submission) {
	ids := map[uint64]struct{}{}
	for _, submission := range curSubmissions.Submissions {
		ids[submission.Id] = struct{}{}
	}
	countNew := 0
	for _, submission := range newSubmissions {
		if _, ok := ids[submission.Id]; !ok {
			countNew += 1
			curSubmissions.Submissions = append(curSubmissions.Submissions, submission)
		}
	}
	r.logger.Info(fmt.Sprintf("added %d submissions", countNew))
}

func (r *Repo) readStorageFile() (*proto.SubmissonsStorage, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Repo) writeStorageFile(storage *proto.SubmissonsStorage) error {
	return fmt.Errorf("not implemented")
}
