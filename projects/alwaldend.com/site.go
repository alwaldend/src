package site

import (
	"embed"
)

var (
	//go:embed site_archive.tar
	Site embed.FS
)
