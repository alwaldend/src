// Package draft resolves a converted draft artifact, resolves its images, and
// creates an X Article draft.
package draft

import (
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// Credential environment variable references. The values are injected by the
// repository's Vault-backed flow and never read from source.
const (
	EnvAPIKey            = "X_API_KEY"
	EnvAPISecret         = "X_API_SECRET"
	EnvAccessToken       = "X_ACCESS_TOKEN"
	EnvAccessTokenSecret = "X_ACCESS_TOKEN_SECRET"
)

// LoadCredentials reads the four OAuth 1.0a fields from the injected
// environment. It reports a missing field by its reference and never returns or
// prints a value.
func LoadCredentials(getenv func(string) string) (xapi.Credentials, error) {
	credentials := xapi.Credentials{
		APIKey:            getenv(EnvAPIKey),
		APISecret:         getenv(EnvAPISecret),
		AccessToken:       getenv(EnvAccessToken),
		AccessTokenSecret: getenv(EnvAccessTokenSecret),
	}
	if err := credentials.Validate(); err != nil {
		return xapi.Credentials{}, fmt.Errorf("load credentials: %w", err)
	}
	return credentials, nil
}

// LoadCredentialsFromEnv reads credentials from the process environment.
func LoadCredentialsFromEnv() (xapi.Credentials, error) {
	return LoadCredentials(os.Getenv)
}
