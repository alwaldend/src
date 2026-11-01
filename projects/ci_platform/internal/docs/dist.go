package docs

import (
	"embed"
	"net/http"
)

var (
	// returns swagger.json or swagger.yaml
	SwaggerFilesystem http.FileSystem
	// embeds files into the variable on initialization
	//go:embed swagger.*
	dist embed.FS
)

func init() {
	SwaggerFilesystem = http.FS(dist)
}
