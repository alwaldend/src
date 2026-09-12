package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

var (
	zonePattern = regexp.MustCompile(`^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+\.?$`)
	namePattern = regexp.MustCompile(`^(\*\.)?[_a-zA-Z0-9-]+(\.[_a-zA-Z0-9-]+)*\.?$`)
)

// domainName resolves only the name equivalence needed for file ownership.
// Provider values, TTLs, priorities, and destination expansion remain owned by
// the Terraform module.
func domainName(name, zone string) (string, error) {
	if name == "@" {
		return zone, nil
	}
	if !namePattern.MatchString(name) {
		return "", fmt.Errorf("invalid domain name %q", name)
	}
	canonical := strings.ToLower(strings.TrimSuffix(name, "."))
	if canonical == zone || strings.HasSuffix(canonical, "."+zone) {
		return canonical, nil
	}
	if strings.HasSuffix(name, ".") {
		return "", fmt.Errorf("absolute domain name %q is outside zone %s", name, zone)
	}
	return canonical + "." + zone, nil
}

func parseRecords(source string, data []byte, zone string) ([]Record, error) {
	if err := CheckSource(source, data); err != nil {
		return nil, err
	}
	var document struct {
		Records map[string]map[string]json.RawMessage `json:"records"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}
	if document.Records == nil {
		return nil, fmt.Errorf("%s: expected a records object", source)
	}
	keys := make([]string, 0, len(document.Records))
	for key := range document.Records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var records []Record
	var problems []error
	for _, key := range keys {
		members := document.Records[key]
		var views []string
		if err := json.Unmarshal(members["dsp"], &views); err != nil || len(views) == 0 {
			problems = append(problems, fmt.Errorf("%s record %q: expected a nonempty dsp array", source, key))
			continue
		}
		validViews := true
		for _, view := range views {
			if view != "global" && view != "dc1" && view != "all" {
				problems = append(problems, fmt.Errorf("%s record %q: unsupported DNS view %q", source, key, view))
				validViews = false
			}
		}
		if !validViews {
			continue
		}
		views = sortedUnique(views)
		types := make([]string, 0, len(members))
		for kind := range members {
			if kind != "dsp" {
				types = append(types, kind)
			}
		}
		sort.Strings(types)
		if len(types) == 0 {
			problems = append(problems, fmt.Errorf("%s record %q: no DNS type members", source, key))
		}
		for _, kind := range types {
			var member struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(members[kind], &member); err != nil {
				problems = append(problems, fmt.Errorf("%s record %q type %s: %w", source, key, kind, err))
				continue
			}
			name, err := domainName(member.Name, zone)
			if err != nil {
				problems = append(problems, fmt.Errorf("%s record %q type %s: %w", source, key, kind, err))
				continue
			}
			records = append(records, Record{Source: source, Key: key, Name: name, Type: kind, Views: views})
		}
	}
	return records, errors.Join(problems...)
}

// CheckRecords requires one dnsconfig.json file per exact canonical domain.
// Different values, types, and views may coexist in the same file. Parent and
// child domain names are distinct ownership units.
func CheckRecords(records []Record) error {
	names := map[string]Record{}
	var problems []error
	for _, record := range records {
		previous, exists := names[record.Name]
		if exists && previous.Source != record.Source {
			problems = append(problems, fmt.Errorf("domain %s is managed by multiple dnsconfig.json files: %s record %q and %s record %q", record.Name, previous.Source, previous.Key, record.Source, record.Key))
			continue
		}
		if !exists {
			names[record.Name] = record
		}
	}
	return errors.Join(problems...)
}

// WriteMarkdown prints a stable table grouped by source file and domain.
// Empty source files appear with empty domain, type, and view columns.
func WriteMarkdown(output io.Writer, report Report) error {
	type row struct {
		types []string
		views []string
	}
	files := map[string]map[string]*row{}
	for _, source := range report.Sources {
		files[source] = map[string]*row{}
	}
	for _, record := range report.Records {
		if files[record.Source] == nil {
			files[record.Source] = map[string]*row{}
		}
		item := files[record.Source][record.Name]
		if item == nil {
			item = &row{}
			files[record.Source][record.Name] = item
		}
		item.types = append(item.types, record.Type)
		item.views = append(item.views, record.Views...)
	}
	sources := make([]string, 0, len(files))
	for source := range files {
		sources = append(sources, source)
	}
	sort.Strings(sources)
	var table strings.Builder
	table.WriteString("| Source | Domain | Types | Views |\n| --- | --- | --- | --- |\n")
	for _, source := range sources {
		names := make([]string, 0, len(files[source]))
		for name := range files[source] {
			names = append(names, name)
		}
		sort.Strings(names)
		if len(names) == 0 {
			fmt.Fprintf(&table, "| %s | - | - | - |\n", tableCell(source))
		}
		for _, name := range names {
			item := files[source][name]
			fmt.Fprintf(&table, "| %s | %s | %s | %s |\n", tableCell(source), tableCell(name), tableCell(strings.Join(sortedUnique(item.types), ", ")), tableCell(strings.Join(sortedUnique(item.views), ", ")))
		}
	}
	_, err := io.WriteString(output, table.String())
	return err
}

func sortedUnique(values []string) []string {
	seen := map[string]bool{}
	for _, value := range values {
		seen[value] = true
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func tableCell(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "|", "&#124;", "\n", " ", "\r", " ").Replace(value)
}
