// embeds a directory with the frontend into FrontendFilesystem
package web

import (
	"embed"
	"net/http"
)

var (
	FrontendFilesystem http.FileSystem
	// embeds directory into the variable
	// //go:embed dist
	dist embed.FS
)

func init() {
	FrontendFilesystem = http.FS(dist)
}
