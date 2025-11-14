package main

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	items []string,
	manifests []string,
	output string,
	marshalOptions *protojson.MarshalOptions,
) error {
	release := &contracts.Release{}
	for _, item := range items {
		if err := self.addItem(release, item); err != nil {
			return fmt.Errorf("could not add item %s to release: %w", item, err)
		}
	}
	for _, manifest := range manifests {
		if err := self.addManifest(release, manifest); err != nil {
			return fmt.Errorf("could not add manifest %s to release: %w", manifest, err)
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

// Add item information in-place
func (self *Generator) addItem(release *contracts.Release, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read item manifest %s: %w", path, err)
	}
	item := &contracts.ReleaseItem{}
	if err := protojson.Unmarshal(content, item); err != nil {
		return fmt.Errorf("could not unmarshal item manifest %s: %w", path, err)
	}
	if item.File == nil || item.File.LocalPath == "" {
		return fmt.Errorf("item manifest %s is missing file.local_path", path)
	}
	stat, err := os.Stat(item.File.LocalPath)
	if err != nil {
		return fmt.Errorf("could not stat local file %s: %w", path, err)
	}
	item.File.Name = filepath.Base(item.File.LocalPath)
	item.File.Size = stat.Size()
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
		item.File.Hashes = append(item.File.Hashes, &contracts.ReleaseHash{
			Algo:    hashType.String(),
			Content: fmt.Sprintf("%x", hashObj.Sum(nil)),
		})
	}
	release.Items = append(release.Items, item)
	return nil
}
