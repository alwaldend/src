package fsstore

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

func readMarkdownFile(path string, maximum int64) ([]byte, error) {
	content, err := readRegularFile(path, maximum)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(content) {
		return nil, errors.New("Markdown content must be valid UTF-8")
	}
	if bytes.IndexByte(content, 0) >= 0 {
		return nil, errors.New("Markdown content must not contain NUL bytes")
	}
	return content, nil
}
