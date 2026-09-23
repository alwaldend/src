// Command prepare_content derives Hugo metadata for plain repository Markdown.
package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

var (
	leadingHeading = regexp.MustCompile(`^\s*#[ \t]+([^\r\n]+)`)
	closingHashes  = regexp.MustCompile(`[ \t]+#+[ \t]*$`)
)

func frontMatter(title string, derived bool) ([]byte, error) {
	metadata := map[string]any{"title": title}
	if derived {
		metadata["al_title_from_heading"] = true
	} else {
		metadata["al_generated_section"] = true
	}
	data, err := json.Marshal(metadata)
	return append(data, '\n', '\n'), err
}

func sectionTitle(dir string) string {
	name := path.Base(dir)
	switch name {
	case "openspec":
		return "OpenSpec"
	case "specs":
		return "Specifications"
	}
	name = strings.NewReplacer("-", " ", "_", " ").Replace(name)
	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func prepare(input, output string) error {
	src, err := os.Open(input)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(output)
	if err != nil {
		return err
	}
	defer dst.Close()
	reader, writer := tar.NewReader(src), tar.NewWriter(dst)
	existing, sections := map[string]bool{}, map[string]bool{}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := path.Clean(header.Name)
		existing[name] = true
		var content io.Reader = reader
		if header.Typeflag == tar.TypeReg && strings.HasPrefix(name, "content/docs/") && path.Ext(name) == ".md" {
			data, err := io.ReadAll(reader)
			if err != nil {
				return err
			}
			if match := leadingHeading.FindSubmatch(data); match != nil {
				title := strings.TrimSpace(string(match[1]))
				title = strings.TrimSpace(closingHashes.ReplaceAllString(title, ""))
				metadata, err := frontMatter(title, true)
				if err != nil {
					return err
				}
				data = append(metadata, data...)
				for dir := path.Dir(name); strings.HasPrefix(dir, "content/docs/"); dir = path.Dir(dir) {
					sections[dir] = true
				}
			}
			header.Size = int64(len(data))
			content = bytes.NewReader(data)
		}
		if err := writer.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(writer, content); err != nil {
			return err
		}
	}
	var dirs []string
	for dir := range sections {
		if existing[dir+"/_index.md"] {
			continue
		}
		// Never turn a source leaf bundle, or one of its resources, into a section.
		leaf := false
		for parent := dir; strings.HasPrefix(parent, "content/docs/"); parent = path.Dir(parent) {
			leaf = leaf || existing[parent+"/index.md"]
		}
		if !leaf {
			dirs = append(dirs, dir)
		}
	}
	slices.Sort(dirs)
	for _, dir := range dirs {
		data, err := frontMatter(sectionTitle(dir), false)
		if err != nil {
			return err
		}
		if err := writer.WriteHeader(&tar.Header{Name: dir + "/_index.md", Mode: 0o644, Size: int64(len(data))}); err != nil {
			return err
		}
		if _, err := writer.Write(data); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return dst.Close()
}

func main() {
	input := flag.String("input", "", "Packaged Hugo source tar")
	output := flag.String("output", "", "Prepared Hugo source tar")
	flag.Parse()
	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "prepare_content requires --input and --output")
		os.Exit(2)
	}
	if err := prepare(*input, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
