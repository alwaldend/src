package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	proto "git.alwaldend.com/src/golang/leetcode-downloader/proto"
)

const SUBMISSIONS_REQUEST_LIMIT = 20

type Downloader struct {
	config *proto.Config
	client *http.Client
	ctx    context.Context
	logger *slog.Logger
}

func NewDownloader(config *proto.Config, ctx context.Context, logger *slog.Logger) (*Downloader, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	downloader := &Downloader{
		config: config,
		client: &http.Client{Jar: jar},
		ctx:    ctx,
		logger: logger,
	}
	return downloader, nil
}

func (d *Downloader) GetSubmissions(offset uint64, limit uint64) ([]*proto.Submission, error) {
	res := []*proto.Submission{}
	lastKey := ""
	for limit > 0 {
		count := min(limit, SUBMISSIONS_REQUEST_LIMIT)
		d.logger.Info(fmt.Sprintf("fetching %d submissions with offset %d", count, offset))
		response, err := d.getSubmissions(offset, count, lastKey)
		if err != nil {
			return nil, fmt.Errorf("could not get submissions %d-%d: %w", offset, limit, err)
		}
		res = append(res, response.SubmissionsDump...)
		if !response.HasNext {
			break
		}
		lastKey = response.LastKey
		offset += count
		limit -= count
		time.Sleep(2)
	}
	d.logger.Info(fmt.Sprintf("fetched %d submissions with offset %d", offset, limit))
	return res, nil
}

func (d *Downloader) getSubmissions(offset uint64, limit uint64, lastKey string) (*proto.SubmissonsResponse, error) {
	reqUrl, err := url.Parse(strings.Join([]string{d.config.BaseUrl, "api", "submissions"}, "/"))
	if err != nil {
		return nil, fmt.Errorf("request '%s' failed: %w", reqUrl, err)
	}
	query := reqUrl.Query()
	query.Set("offset", fmt.Sprint(offset))
	query.Set("limit", fmt.Sprint(limit))
	query.Set("lastkey", lastKey)
	reqUrl.RawQuery = query.Encode()
	request, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request '%s' failed: %w", reqUrl, err)
	}
	for key, value := range d.config.Headers {
		request.Header.Set(key, value)
	}
	request.Header.Set("Cookie", d.config.Cookie)
	response, err := d.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request '%s' failed: %w", reqUrl, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("request '%s' failed with '%d': %w", reqUrl, response.StatusCode, err)
	}
	bodyJson := &proto.SubmissonsResponse{}
	err = json.Unmarshal(body, bodyJson)
	if err != nil {
		return nil, fmt.Errorf("request '%s' failed with '%d' and response '%s': %w", reqUrl, response.StatusCode, string(body), err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request '%s' failed with '%d' and response '%s'", reqUrl, response.StatusCode, string(body))
	}
	return bodyJson, nil
}
