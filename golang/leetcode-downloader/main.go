package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	proto "git.alwaldend.com/proto/leetcode-downloader"

	"git.alwaldend.com/src/golang/leetcode-downloader/model"
)

func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	logger := slog.New(slog.NewTextHandler(stderr, &slog.HandlerOptions{}))
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	config, err := parseArgs(args, getenv)
	if err != nil {
		return err
	}
	repo := NewRepo(config, ctx, logger)
	if err != nil {
		return err
	}
	downloader, err := NewDownloader(config, ctx, logger)
	if err != nil {
		return err
	}
	generator, err := NewGenerator(config, ctx, logger)
	if err != nil {
		return err
	}
	switch config.Action {
	case proto.CliAction_DOWNLOAD:
		submissions, err := downloader.GetSubmissions(config.Offset, config.Limit)
		if err != nil {
			return err
		}
		text, err := json.Marshal(submissions)
		if err != nil {
			return err
		}
		text = append(text, '\n')
		stdout.Write(text)
	case proto.CliAction_UPDATE:
		var submissions []*proto.Submission
		if len(config.ActionArgs) > 0 {
			response := &proto.SubmissonsResponse{}
			if err = json.Unmarshal([]byte(config.ActionArgs[0]), response); err != nil {
				return fmt.Errorf("could not parse submissions from stdin: %w", err)
			}
			submissions = response.SubmissionsDump
		} else {
			submissions, err = downloader.GetSubmissions(config.Offset, config.Limit)
			if err != nil {
				return fmt.Errorf("could not download submissions: %w", err)
			}
		}
		return repo.AddSubmissions(submissions)
	case proto.CliAction_GENERATE:
		submissionsText, _ := os.ReadFile(config.SubmissionsFile)
		var submissions *proto.SubmissonsResponse
		err := json.Unmarshal(submissionsText, &submissions)
		if err != nil {
			return err
		}
		err = generator.Generate(submissions.SubmissionsDump)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported action: %s", config.Action)
	}
	return nil
}

func parseArgs(args []string, getenv func(string) string) (*proto.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configDefault, err := model.DefaultConfig()
	if err != nil {
		return nil, err
	}
	config := configDefault
	cookie := getenv("LEETCODE_COOKIE")
	if cookie == "" {
		cookie = config.Cookie
	}
	const sessionReplace = "${LEETCODE_SESSION}"
	flagset := flag.NewFlagSet("flags", flag.ContinueOnError)
	flagset.StringVar(&config.RootDir, "root-dir", wd, "")
	flagset.StringVar(&config.SubmissionsFile, "submissions-file", "", "")
	flagset.StringVar(&config.BaseUrl, "base_url", configDefault.BaseUrl, "")
	flagset.StringVar(&config.Cookie, "cookie", cookie, "")
	flagset.Uint64Var(&config.Offset, "offset", configDefault.Offset, "")
	flagset.Uint64Var(&config.Limit, "limit", configDefault.Limit, "")
	err = flagset.Parse(args)
	if err != nil {
		return nil, err
	}
	action := flagset.Args()
	if len(action) == 0 {
		return nil, fmt.Errorf("action is not specified")
	}
	actionName := action[0]
	val, ok := proto.CliAction_value[strings.ToUpper(actionName)]
	if !ok {
		return nil, fmt.Errorf("invalid action %s", actionName)
	}
	config.Action = proto.CliAction(val)
	if len(action) > 1 {
		config.ActionArgs = action[1:]
	}
	return config, nil
}
