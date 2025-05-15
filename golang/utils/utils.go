package utils

import (
	"bytes"
	"io"
)

func StreamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func StreamToString(stream io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
