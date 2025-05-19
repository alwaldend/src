package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type Parse struct {
	state      *State
	readmeName string
	useGit     bool
}

func NewParse(state *State, readmeName string, useGit bool) *Parse {
	return &Parse{state: state, readmeName: readmeName, useGit: useGit}
}

func (self *Parse) ParsePathsAndOutput(paths []string, rootPath string) error {
	result, err := self.ParsePaths(paths, rootPath)
	if err != nil {
		return err
	}
	resultJson, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("could not serialise result: %w", err)
	}
	_, err = fmt.Fprintf(self.state.Stdout, "%s\n", resultJson)
	if err != nil {
		return fmt.Errorf("could not print: %w", err)
	}
	return nil
}

func (self *Parse) ParsePaths(paths []string, rootPath string) (*ParsePathsResult, error) {
	result := &ParsePathsResult{Dir: rootPath}
	for _, path := range paths {
		readmes, err := self.ParsePath(path, rootPath)
		if err != nil {
			return nil, fmt.Errorf("could not process path %s: %w", path, err)
		}
		result.Readmes = append(result.Readmes, readmes...)
	}
	return result, nil
}

func (self *Parse) ParsePath(path string, rootPath string) ([]*Readme, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path of %s: %w", path, err)
	}
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("could not stat file %s: %w", absPath, err)
	}
	var readmes []*Readme
	if stat.IsDir() {
		readmes, err = self.ParseDirectory(absPath, rootPath)
		if err != nil {
			return nil, err
		}
	} else {
		readme, err := self.ParseFile(absPath, rootPath)
		if err != nil {
			return nil, err
		}
		readmes = []*Readme{readme}
	}
	return readmes, nil
}

func (self *Parse) ParseDirectory(dirPath string, rootPath string) ([]*Readme, error) {
	if self.useGit {
		return self.ParseDirectoryGit(dirPath, rootPath)
	}
	return self.ParseDirectoryWalk(dirPath, rootPath)
}

func (self *Parse) ParseDirectoryGit(dirPath string, rootPath string) ([]*Readme, error) {
	stdout, err := exec.Command("git", "ls-files", dirPath).Output()
	if err != nil {
		return nil, err
	}
	lines := regexp.MustCompile("\r?\n").Split(string(stdout), -1)
	result := []*Readme{}
	ignore := filepath.Join(dirPath, self.readmeName)
	for _, filePath := range lines {
		absPath, err := filepath.Abs(filePath)
		name := filepath.Base(filePath)
		if name != self.readmeName || filePath == ignore {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("could not get absolute path of %s: %w", filePath, err)
		}
		readme, err := self.ParseFile(absPath, rootPath)
		if err != nil {
			return nil, err
		}
		result = append(result, readme)
	}
	return result, nil
}

func (self *Parse) ParseDirectoryWalk(dirPath string, rootPath string) ([]*Readme, error) {
	result := []*Readme{}
	ignore := filepath.Join(dirPath, self.readmeName)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		name := filepath.Base(path)
		if name != self.readmeName || path == ignore {
			return nil
		}
		readme, err := self.ParseFile(path, rootPath)
		if err != nil {
			return err
		}
		if readme != nil {
			result = append(result, readme)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse directory %s: %w", dirPath, err)
	}
	return result, nil
}

func (self *Parse) ParseFiles(filePaths []string, rootPath string) ([]*Readme, error) {
	result := []*Readme{}
	for _, path := range filePaths {
		readme, err := self.ParseFile(path, rootPath)
		if err != nil {
			return nil, err
		}
		if readme != nil {
			result = append(result, readme)
		}
	}
	return result, nil
}

func (self *Parse) ParseFile(filePath string, rootPath string) (*Readme, error) {
	relativePath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("could not create relative path of %s: %w", filePath, err)
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filePath, err)
	}
	segments := strings.Split(string(fileContent), "---")
	if len(segments) < 3 {
		return nil, nil
	}
	var data map[string]any
	err = yaml.Unmarshal([]byte(segments[1]), &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal yaml in file %s: %w", filePath, err)
	}
	return &Readme{Path: filePath, RelativePath: relativePath, Data: data}, nil
}
