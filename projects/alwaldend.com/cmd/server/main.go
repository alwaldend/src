package main

import (
	"fmt"
	"log"
	"io/fs"
	"net/http"
	"os"

	site "git.alwaldend.com/alwaldend/src/projects/alwaldend.com"
)

func run() error {
	siteFs, err := New(site.Site)
	if err != nil {
		return fmt.Errorf("could not create fs: %w", err)
	}
    siteFsSub, err := fs.Sub(siteFs, "site.destination")
    if err != nil {
        return fmt.Errorf("could not get site fs subtree: %w", err)
    }
	fileServer := http.FileServer(http.FS(siteFsSub))
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
