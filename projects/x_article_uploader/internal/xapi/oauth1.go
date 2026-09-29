// Package xapi talks to the X API with OAuth 1.0a user context.
//
// OAuth 1.0a is used because its user-context token does not expire on a fixed
// schedule, so publication does not depend on rotating a single-use refresh
// token. Signing is HMAC-SHA1 per RFC 5849, implemented with the standard
// library so no signing dependency is added.
package xapi

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Credentials are the four OAuth 1.0a fields. Every field is required.
type Credentials struct {
	APIKey            string
	APISecret         string
	AccessToken       string
	AccessTokenSecret string
}

// Validate reports the first missing field's reference so the draft command can
// name it without ever printing a value.
func (c Credentials) Validate() error {
	missing := []string{}
	if c.APIKey == "" {
		missing = append(missing, "X_API_KEY")
	}
	if c.APISecret == "" {
		missing = append(missing, "X_API_SECRET")
	}
	if c.AccessToken == "" {
		missing = append(missing, "X_ACCESS_TOKEN")
	}
	if c.AccessTokenSecret == "" {
		missing = append(missing, "X_ACCESS_TOKEN_SECRET")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing credential reference(s): %s", strings.Join(missing, ", "))
	}
	return nil
}

// Sign returns the Authorization header value for one request. bodyParams are
// the request's form parameters; a request with a multipart body passes none,
// since RFC 5849 includes only query and form-urlencoded parameters.
func (c Credentials) Sign(method, rawURL string, queryParams, bodyParams url.Values, now time.Time) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse request URL: %w", err)
	}
	nonce, err := nonce()
	if err != nil {
		return "", err
	}
	oauth := url.Values{}
	oauth.Set("oauth_consumer_key", c.APIKey)
	oauth.Set("oauth_nonce", nonce)
	oauth.Set("oauth_signature_method", "HMAC-SHA1")
	oauth.Set("oauth_timestamp", strconv.FormatInt(now.Unix(), 10))
	oauth.Set("oauth_token", c.AccessToken)
	oauth.Set("oauth_version", "1.0")

	params := url.Values{}
	for key, values := range queryParams {
		for _, value := range values {
			params.Add(key, value)
		}
	}
	for key, values := range bodyParams {
		for _, value := range values {
			params.Add(key, value)
		}
	}
	for key, values := range oauth {
		for _, value := range values {
			params.Add(key, value)
		}
	}

	signature, err := signature(c, method, parsed, params)
	if err != nil {
		return "", err
	}
	oauth.Set("oauth_signature", signature)
	return "OAuth " + encodeHeader(oauth), nil
}

// signature computes the RFC 5849 HMAC-SHA1 signature for the request.
func signature(c Credentials, method string, parsed *url.URL, params url.Values) (string, error) {
	baseURL := *parsed
	baseURL.RawQuery = ""
	baseURL.Fragment = ""

	base := strings.Join([]string{
		strings.ToUpper(method),
		percentEncode(baseURL.String()),
		percentEncode(normalize(params)),
	}, "&")

	key := percentEncode(c.APISecret) + "&" + percentEncode(c.AccessTokenSecret)
	mac := hmac.New(sha1.New, []byte(key))
	if _, err := mac.Write([]byte(base)); err != nil {
		return "", fmt.Errorf("sign request: %w", err)
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// normalize sorts and percent-encodes the request parameters per RFC 5849.
func normalize(params url.Values) string {
	pairs := make([]string, 0, len(params))
	for key, values := range params {
		for _, value := range values {
			pairs = append(pairs, percentEncode(key)+"="+percentEncode(value))
		}
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "&")
}

// encodeHeader renders the oauth parameters as a comma-separated header value.
func encodeHeader(params url.Values) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, percentEncode(key)+`="`+percentEncode(params.Get(key))+`"`)
	}
	return strings.Join(parts, ", ")
}

// percentEncode applies RFC 5849's percent-encoding, which is stricter than
// url.QueryEscape: it escapes every byte outside the unreserved set.
func percentEncode(value string) string {
	var b strings.Builder
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if isUnreserved(ch) {
			b.WriteByte(ch)
			continue
		}
		b.WriteString("%")
		b.WriteString(strings.ToUpper(hex.EncodeToString([]byte{ch})))
	}
	return b.String()
}

func isUnreserved(ch byte) bool {
	switch {
	case ch >= 'a' && ch <= 'z':
		return true
	case ch >= 'A' && ch <= 'Z':
		return true
	case ch >= '0' && ch <= '9':
		return true
	}
	return ch == '-' || ch == '.' || ch == '_' || ch == '~'
}

// nonce returns a random request nonce.
func nonce() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
