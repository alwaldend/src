package site

import (
	_ "embed"
)

var (
	//go:embed site_archive.tar
	Site []byte
)
