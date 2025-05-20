package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type Parser struct {
	state     *State
	outputter *Outputter
}

func NewParse(state *State, outputter *Outputter) *Parser {
	return &Parser{state: state, outputter: outputter}
}

func (self *Parser) ParsePathsAndOutput(paths []string, rootPath string) error {
	result, err := self.ParsePaths(paths, rootPath)
	if err != nil {
		return err
	}
	switch self.state.Config.Parse.OutputType {
	case ParseOutputTypeJson:
		return self.outputter.OutputJson(result)
	case ParseOutputTypeMarkdown:
		return self.outputter.OutputMarkdown(result)
	default:
		return fmt.Errorf("output type %s is not supported", self.state.Config.Parse.OutputType)
	}
}

func (self *Parser) ParsePaths(paths []string, rootPath string) (*ParsePathsResult, error) {
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

func (self *Parser) ParsePath(path string, rootPath string) ([]*Readme, error) {
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
		if readme != nil {
			readmes = []*Readme{readme}
		}
	}
	return readmes, nil
}

func (self *Parser) ParseDirectory(dirPath string, rootPath string) ([]*Readme, error) {
	if self.state.Config.UseGit {
		return self.ParseDirectoryGit(dirPath, rootPath)
	}
	return self.ParseDirectoryWalk(dirPath, rootPath)
}

func (self *Parser) ParseDirectoryGit(dirPath string, rootPath string) ([]*Readme, error) {
	stdout, err := exec.Command("git", "ls-files", dirPath).Output()
	if err != nil {
		return nil, err
	}
	lines := regexp.MustCompile("\r?\n").Split(string(stdout), -1)
	result := []*Readme{}
	ignore := filepath.Join(dirPath, self.state.Config.ReadmeName)
	for _, filePath := range lines {
		absPath, err := filepath.Abs(filePath)
		name := filepath.Base(filePath)
		if name != self.state.Config.ReadmeName || filePath == ignore {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("could not get absolute path of %s: %w", filePath, err)
		}
		readme, err := self.ParseFile(absPath, rootPath)
		if err != nil {
			return nil, err
		}
		if readme != nil {
			result = append(result, readme)
		}
	}
	return result, nil
}

func (self *Parser) ParseDirectoryWalk(dirPath string, rootPath string) ([]*Readme, error) {
	result := []*Readme{}
	ignore := filepath.Join(dirPath, self.state.Config.ReadmeName)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		name := filepath.Base(path)
		if name != self.state.Config.ReadmeName || path == ignore {
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

func (self *Parser) ParseFiles(filePaths []string, rootPath string) ([]*Readme, error) {
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

func (self *Parser) ParseFile(filePath string, rootPath string) (*Readme, error) {
	relativePath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("could not create relative path of %s: %w", filePath, err)
	}
	dirRelative := filepath.Dir(relativePath)
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
	return &Readme{
		Path:             filePath,
		PathRelative:     relativePath,
		Data:             data,
		Dir:              filepath.Dir(filePath),
		DirRelative:      dirRelative,
		DirRelativeSplit: strings.Split(dirRelative, string(os.PathSeparator)),
		DirName:          filepath.Base(dirRelative),
	}, nil
}
