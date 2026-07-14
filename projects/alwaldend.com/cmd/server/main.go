package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	site "git.alwaldend.com/alwaldend/src/projects/alwaldend.com"
)

func run() error {
	siteFs, err := fs.Sub(site.Site, "site.destination")
	if err != nil {
		return fmt.Errorf("could not create fs: %w", err)
	}
	fileServer := http.FileServer(http.FS(siteFs))
	http.Handle("/", fileServer)
	log.Println("Listening on 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		return fmt.Errorf("could not listen and serve: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
}
