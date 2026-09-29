// Command draft resolves a converted draft artifact's images and creates an X
// Article draft. It never publishes: publication is out of scope for this tool,
// so the only operation it performs is draft creation.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

func main() {
	artifactPath := flag.String("artifact", "", "path to the converted draft artifact JSON")
	workspace := flag.String("workspace", ".", "repository root the post package resolves against")
	cache := flag.String("cache", "", "path of the digest-to-media_id cache")
	flag.Parse()

	if err := run(*artifactPath, *workspace, *cache); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(artifactPath, workspace, cache string) error {
	if artifactPath == "" {
		return fmt.Errorf("--artifact is required")
	}
	credentials, err := draft.LoadCredentialsFromEnv()
	if err != nil {
		return err
	}
	if cache == "" {
		cache = filepath.Join(workspace, "out", "x_article_uploader", "media_cache.json")
	}
	client := xapi.New(credentials)
	pub := draft.New(client, workspace, cache)
	if err := pub.LoadCache(); err != nil {
		return err
	}
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		return err
	}
	draft, err := pub.CreateDraft(artifact)
	if err != nil {
		return err
	}
	fmt.Printf("created draft %s (not published)\n", draft.ID)
	return nil
}
