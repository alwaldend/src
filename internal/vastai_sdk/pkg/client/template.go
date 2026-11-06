package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

// Template contains resolved public metadata without retaining environment values.
type Template struct {
	// Image is the image reference including its separately configured tag.
	Image string
	// EnvironmentOnly is true only for well-formed literal environment assignments.
	EnvironmentOnly bool
}

// Template reads one exact template version without creating or editing anything.
func (c *Client) Template(ctx context.Context, hash string) (Template, error) {
	filters, err := json.Marshal(map[string]any{"hash_id": map[string]string{"eq": hash}})
	if err != nil {
		return Template{}, fmt.Errorf("encode template filter: %w", err)
	}
	var result struct {
		// Templates holds only the fields needed for identity and networking validation.
		Templates []struct {
			// Hash identifies the exact template version returned.
			Hash string `json:"hash_id"`
			// Image names the container repository or full reference.
			Image string `json:"image"`
			// Tag is the optional separate container tag.
			Tag string `json:"tag"`
			// Environment distinguishes an explicitly empty option string from missing data.
			Environment *string `json:"env"`
		} `json:"templates"`
	}
	request := c.api.SearchAPI.SearchTemplates(ctx).SelectFilters(string(filters)).SelectCols(`["hash_id","image","tag","env"]`)
	if err := generated(c, http.MethodGet, "/template/", &result, true, request.Execute); err != nil {
		return Template{}, fmt.Errorf("read template: %w", err)
	}
	if len(result.Templates) != 1 {
		return Template{}, errors.New("template lookup must return exactly one matching template")
	}
	remote := result.Templates[0]
	if remote.Hash != hash || strings.TrimSpace(remote.Image) == "" || remote.Environment == nil {
		return Template{}, errors.New("template response has incomplete or mismatched metadata")
	}
	image := remote.Image
	if remote.Tag != "" && !strings.Contains(image, "@") && !strings.Contains(image[strings.LastIndex(image, "/")+1:], ":") {
		image += ":" + remote.Tag
	}
	return Template{Image: image, EnvironmentOnly: environmentOptionsOnly(*remote.Environment)}, nil
}

// environmentOptionsOnly accepts only literal Docker environment assignments, failing closed on every other option.
func environmentOptionsOnly(options string) bool {
	words, valid := optionWords(options)
	if !valid {
		return false
	}
	for i := 0; i < len(words); i++ {
		word := words[i]
		var assignment string
		switch {
		case word == "-e" || word == "--env":
			i++
			if i == len(words) {
				return false
			}
			assignment = words[i]
		case strings.HasPrefix(word, "--env="):
			assignment = strings.TrimPrefix(word, "--env=")
		case strings.HasPrefix(word, "-e"):
			assignment = strings.TrimPrefix(word, "-e")
		default:
			return false
		}
		name, _, found := strings.Cut(assignment, "=")
		if !found || name == "" {
			return false
		}
		for index, character := range name {
			if character != '_' && !(character >= 'A' && character <= 'Z') && !(character >= 'a' && character <= 'z') && !(index > 0 && character >= '0' && character <= '9') {
				return false
			}
		}
	}
	return true
}

// optionWords tokenizes literal Docker options without evaluating shell expressions or reading files.
func optionWords(input string) ([]string, bool) {
	var words []string
	var word strings.Builder
	var quote rune
	active, escaped := false, false
	for _, character := range input {
		if escaped {
			word.WriteRune(character)
			escaped = false
			continue
		}
		if character == '\\' && quote != '\'' {
			escaped, active = true, true
			continue
		}
		if quote != 0 {
			if character == quote {
				quote = 0
			} else {
				word.WriteRune(character)
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote, active = character, true
		} else if unicode.IsSpace(character) {
			if active {
				words = append(words, word.String())
				word.Reset()
				active = false
			}
		} else {
			word.WriteRune(character)
			active = true
		}
	}
	if escaped || quote != 0 {
		return nil, false
	}
	if active {
		words = append(words, word.String())
	}
	return words, true
}
