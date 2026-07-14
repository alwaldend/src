package site

import (
	"embed"
)

var (
	//go:embed site.destination
	Site embed.FS
)
