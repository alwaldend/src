package experiments

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// collect downloads one completed job independently of submission and observation.
func (c *controller) collect(ctx context.Context, index int) {
	e := *c.journal.Experiments[index]
	c.workers++
	c.tasks.Add(1)
	go func() {
		defer c.tasks.Done()
		var directory string
		var result struct {
			// Files names the original media downloaded by the native CLI.
			Files []struct {
				// Path identifies one downloaded original.
				Path string `json:"path"`
			} `json:"files"`
		}
		var warning bool
		var err error
	retrieval:
		for retry := 0; retry < 3; retry++ {
			if retry > 0 {
				timer := time.NewTimer(time.Duration(retry) * time.Second)
				select {
				case <-ctx.Done():
					timer.Stop()
					err = fmt.Errorf("artifact retry interrupted: %w", ctx.Err())
					break retrieval
				case <-timer.C:
				}
			}
			// Every attempt retains separate receipts and any earlier partial originals.
			attempt := fmt.Sprintf("%d", time.Now().UnixNano())
			directory = filepath.Join(c.directory, e.Directory, "media-"+attempt)
			result.Files = nil
			warning, err = callCLI(ctx, c.journal.Config, filepath.Join(c.directory, e.Directory, "download-"+attempt), &result, "download", e.PromptID, "-o", directory)
			if err == nil || ctx.Err() != nil {
				break
			}
		}
		files := []media{}
		if err == nil && len(result.Files) == 0 {
			err = fmt.Errorf("native download returned no media")
		}
		seen := map[string]bool{}
		for _, item := range result.Files {
			if err != nil {
				break
			}
			var file media
			file, err = validateMedia(c.directory, directory, item.Path)
			if seen[file.Path] {
				err = fmt.Errorf("download receipt repeated a media path")
			}
			seen[file.Path] = true
			files = append(files, file)
		}
		select {
		case c.events <- completion{index: index, collected: true, media: files, warning: warning, err: err}:
		case <-ctx.Done():
		}
	}()
}

// validateMedia confines receipts to regular original files beneath the experiment output.
func validateMedia(run, directory, path string) (media, error) {
	if !filepath.IsAbs(path) {
		return media{}, fmt.Errorf("download path must be absolute")
	}
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		return media{}, fmt.Errorf("resolve downloaded media: %w", err)
	}
	realDirectory, err := filepath.EvalSymlinks(directory)
	if err != nil {
		return media{}, fmt.Errorf("resolve media directory: %w", err)
	}
	// The directory itself may not redirect outside the private experiment tree.
	if filepath.Clean(realDirectory) != filepath.Clean(directory) {
		return media{}, fmt.Errorf("media directory contains a symbolic link")
	}
	relative, err := filepath.Rel(realDirectory, real)
	if err != nil {
		return media{}, fmt.Errorf("locate downloaded media: %w", err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) || real != filepath.Clean(path) {
		return media{}, fmt.Errorf("downloaded media escapes its directory or contains a symbolic link")
	}
	f, err := os.Open(real)
	if err != nil {
		return media{}, fmt.Errorf("open downloaded media: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return media{}, fmt.Errorf("inspect downloaded media: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		f.Close()
		return media{}, fmt.Errorf("downloaded media is not a nonempty regular file")
	}
	hash := sha256.New()
	_, copyErr := io.Copy(hash, f)
	closeErr := f.Close()
	if copyErr != nil {
		return media{}, fmt.Errorf("hash downloaded media: %w", copyErr)
	}
	if closeErr != nil {
		return media{}, fmt.Errorf("close downloaded media: %w", closeErr)
	}
	runRelative, err := filepath.Rel(run, real)
	if err != nil {
		return media{}, fmt.Errorf("record downloaded media path: %w", err)
	}
	return media{Path: filepath.ToSlash(runRelative), SHA256: hex.EncodeToString(hash.Sum(nil)), Bytes: info.Size()}, nil
}

// report coalesces immutable projections without blocking native queue replenishment.
func (c *controller) report() error {
	data, err := json.Marshal(c.journal)
	if err != nil {
		return fmt.Errorf("snapshot report data: %w", err)
	}
	select {
	case c.reports <- data:
	default:
		select {
		case <-c.reports:
		default:
		}
		c.reports <- data
	}
	return nil
}

// renderReports serializes coalesced projections separately from the controller.
func (c *controller) renderReports() {
	defer close(c.reportDone)
	for data := range c.reports {
		var snapshot journal
		err := decodeJSON(bytes.NewReader(data), &snapshot)
		if err == nil {
			err = renderReport(c.directory, snapshot)
		}
		if err != nil {
			select {
			case c.reportErrors <- fmt.Errorf("render report snapshot: %w", err):
			default:
			}
		}
	}
}

// renderReport writes original media and collapsed data from an immutable journal projection.
func renderReport(directory string, snapshot journal) error {
	var page strings.Builder
	page.WriteString("<!doctype html><meta charset=utf-8><meta http-equiv=Content-Security-Policy content=\"default-src 'none'; img-src 'self' file:; media-src 'self' file:; style-src 'unsafe-inline'\"><title>ComfyUI experiments</title><style>body{background:#171717;color:#eee;font:16px system-ui;max-width:1200px;margin:32px auto;padding:20px}img,video{max-width:100%;max-height:85vh}pre{white-space:pre-wrap;overflow-wrap:anywhere}article{margin-bottom:48px}a{color:#9cf}</style><h1>ComfyUI experiments</h1><p>Execution receipts do not establish visual quality. Data and originals are private until reviewed for publication.</p>")
	for _, e := range snapshot.Experiments {
		if e.State != "completed" {
			continue
		}
		page.WriteString(fmt.Sprintf("<article><h2>Experiment %d — %s</h2>", e.Sequence, html.EscapeString(e.Candidate.Description)))
		for _, file := range e.Media {
			source := (&url.URL{Path: file.Path}).String()
			switch strings.ToLower(filepath.Ext(file.Path)) {
			case ".png", ".jpg", ".jpeg", ".webp", ".gif":
				page.WriteString("<img loading=lazy src=\"" + html.EscapeString(source) + "\" alt=\"Original experiment output\">")
			case ".mp4", ".webm":
				page.WriteString("<video controls src=\"" + html.EscapeString(source) + "\"></video>")
			case ".wav", ".mp3", ".ogg", ".flac":
				page.WriteString("<audio controls src=\"" + html.EscapeString(source) + "\"></audio>")
			default:
				page.WriteString("<a href=\"" + html.EscapeString(source) + "\">Original media</a>")
			}
		}
		page.WriteString("<details><summary>Prompts, settings and provenance</summary>")
		for _, key := range []string{"positive_prompt", "negative_prompt"} {
			if value, exists := e.Candidate.Notes[key]; exists {
				page.WriteString("<p>" + key + "</p><pre>" + html.EscapeString(fmt.Sprint(value)) + "</pre>")
			}
		}
		data, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			return fmt.Errorf("encode experiment report: %w", err)
		}
		page.WriteString("<pre>" + html.EscapeString(string(data)) + "</pre><p><a href=\"" + e.Directory + "/workflow.json\">Exact workflow snapshot</a></p></details></article>")
	}
	if err := os.WriteFile(filepath.Join(directory, "gallery.html"), []byte(page.String()), 0o600); err != nil {
		return fmt.Errorf("write local gallery: %w", err)
	}
	return nil
}
