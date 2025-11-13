package main

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/src/tools/release/contracts"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Generator struct{}

type GeneratorParams struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

// Generate a manifest file
func (self *Generator) Generate(
	release *contracts.Release,
	items []string,
	output string,
	marshalOptions *protojson.MarshalOptions,
) error {
	release = proto.CloneOf(release)
	for _, item := range items {
		itemSplit := strings.SplitN(item, ":", 2)
		if len(itemSplit) != 2 {
			return fmt.Errorf("invalid item string %s: %s", itemSplit, item)
		}
		itemPath := itemSplit[1]
		switch itemType := itemSplit[0]; itemType {
		case "file":
			if err := self.addFileItem(release, itemPath); err != nil {
				return fmt.Errorf("could not add file %s to release: %w", item, err)
			}
		case "manifest":
			if err := self.addManifest(release, itemPath); err != nil {
				return fmt.Errorf("could not add manifest %s to release: %w", item, err)
			}
		default:
			return fmt.Errorf("invalid item type %s: %s", itemType, item)
		}
	}
	data, err := marshalOptions.Marshal(release)
	if err != nil {
		return fmt.Errorf("could not marshal release %v: %w", release, err)
	}
	if err := os.WriteFile(output, data, 0o444); err != nil {
		return fmt.Errorf("could not write data to file %s: %w", output, err)
	}
	return nil
}

// Add manifest information in-place
func (self *Generator) addManifest(release *contracts.Release, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}
	fileManifest := &contracts.Release{}
	if err := protojson.Unmarshal(content, fileManifest); err != nil {
		return fmt.Errorf("could not parse file %s: %w", path, err)
	}
	proto.Merge(release, fileManifest)
	return nil
}

// Add file item information in-place
func (self *Generator) addFileItem(release *contracts.Release, path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("could not stat file %s: %w", path, err)
	}
	releaseFile := &contracts.ReleaseFile{
		Name: filepath.Base(path),
		Size: stat.Size(),
	}
	for _, hashType := range []crypto.Hash{crypto.SHA256, crypto.MD5} {
		hashObj := hashType.New()
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("could not open file %s: %w", path, err)
		}
		defer file.Close()
		if _, err := io.Copy(hashObj, file); err != nil {
			return fmt.Errorf("could not copy file %s to hash: %w", path, err)
		}
		releaseFile.Hashes = append(releaseFile.Hashes, &contracts.ReleaseHash{
			Algo:    hashType.String(),
			Content: fmt.Sprintf("%x", hashObj.Sum(nil)),
		})
	}
	release.Items = append(release.Items, &contracts.ReleaseItem{File: releaseFile})
	return nil
}
