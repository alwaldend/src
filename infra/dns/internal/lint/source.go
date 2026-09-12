package lint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

// CheckSource rejects duplicate JSON object keys before decoding can silently
// replace a record or field. JSON pointers identify the offending declaration.
func CheckSource(source string, data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := checkValue(decoder, ""); err != nil {
		return fmt.Errorf("%s: %w", source, err)
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("multiple JSON values")
		}
		return fmt.Errorf("%s: %w", source, err)
	}
	return nil
}

func checkValue(decoder *json.Decoder, pointer string) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, nested := token.(json.Delim)
	if !nested {
		return nil
	}
	switch delimiter {
	case '{':
		keys := map[string]bool{}
		for decoder.More() {
			token, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := token.(string)
			if !ok {
				return fmt.Errorf("invalid object key at %s", pointer)
			}
			child := pointer + "/" + strings.ReplaceAll(strings.ReplaceAll(key, "~", "~0"), "/", "~1")
			if keys[key] {
				return fmt.Errorf("duplicate JSON key at %s", child)
			}
			keys[key] = true
			if err := checkValue(decoder, child); err != nil {
				return err
			}
		}
	case '[':
		for index := 0; decoder.More(); index++ {
			if err := checkValue(decoder, fmt.Sprintf("%s/%d", pointer, index)); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unexpected delimiter %s at %s", delimiter, pointer)
	}
	_, err = decoder.Token()
	return err
}

// Scan loads every canonical dnsconfig.json below the workspace at runtime.
// Invalid files do not prevent other source files from being read and checked.
func Scan(workspace fs.FS, zone string) (Report, error) {
	if !zonePattern.MatchString(zone) {
		return Report{}, fmt.Errorf("invalid DNS zone %q", zone)
	}
	zone = strings.ToLower(strings.TrimSuffix(zone, "."))
	sources, err := DiscoverSources(workspace)
	if err != nil {
		return Report{}, fmt.Errorf("discover DNS sources: %w", err)
	}
	report := Report{Sources: sources}
	var problems []error
	for _, source := range sources {
		data, err := fs.ReadFile(workspace, source)
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", source, err))
			continue
		}
		records, err := parseRecords(source, data, zone)
		report.Records = append(report.Records, records...)
		problems = append(problems, err)
	}
	problems = append(problems, CheckRecords(report.Records))
	return report, errors.Join(problems...)
}

// DiscoverSources recursively finds dnsconfig.json files across nested
// workspaces. Scratch, tool caches, VCS metadata, and symlinked directories are
// excluded so generated copies cannot become additional record owners.
func DiscoverSources(workspace fs.FS) ([]string, error) {
	var sources []string
	err := fs.WalkDir(workspace, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && name != "." {
			base := path.Base(name)
			switch base {
			case ".git", "out", "node_modules", ".terraform", ".cache", ".venv", "__pycache__":
				return fs.SkipDir
			}
			if strings.HasPrefix(base, "bazel-") {
				return fs.SkipDir
			}
		}
		if !entry.IsDir() && entry.Name() == "dnsconfig.json" {
			sources = append(sources, name)
		}
		return nil
	})
	return sources, err
}
