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
	outputter *Outputter
	config    *Config
}

func NewParse(config *Config, outputter *Outputter) *Parser {
	return &Parser{config: config, outputter: outputter}
}

func (self *Parser) ParsePathsAndOutput(paths []string, config *ParseConfig) error {
	result, err := self.ParsePaths(paths, config)
	if err != nil {
		return err
	}
	return self.outputter.Output(result, config.OutputType)
}

func (self *Parser) ParsePaths(paths []string, config *ParseConfig) (*ParsePathsResult, error) {
	result := &ParsePathsResult{Dir: config.RootPath}
	for _, path := range paths {
		readmes, err := self.ParsePath(path, config)
		if err != nil {
			return nil, fmt.Errorf("could not process path %s: %w", path, err)
		}
		result.Readmes = append(result.Readmes, readmes...)
	}
	return result, nil
}

func (self *Parser) ParsePath(path string, config *ParseConfig) ([]*Readme, error) {
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
		readmes, err = self.ParseDirectory(absPath, config)
		if err != nil {
			return nil, err
		}
	} else {
		readme, ok, err := self.ParseFile(absPath, config)
		if err != nil {
			return nil, err
		}
		if ok {
			readmes = []*Readme{readme}
		}
	}
	return readmes, nil
}

func (self *Parser) ParseDirectory(dirPath string, config *ParseConfig) ([]*Readme, error) {
	if config.UseGit {
		return self.ParseDirectoryGit(dirPath, config)
	}
	return self.ParseDirectoryWalk(dirPath, config)
}

func (self *Parser) ParseDirectoryGit(dirPath string, config *ParseConfig) ([]*Readme, error) {
	stdout, err := exec.Command("git", "ls-files", dirPath).Output()
	if err != nil {
		return nil, err
	}
	lines := regexp.MustCompile("\r?\n").Split(string(stdout), -1)
	result := []*Readme{}
	for _, filePath := range lines {
		absPath, err := filepath.Abs(filePath)
		if self.IgnoreReadme(filePath, config) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("could not get absolute path of %s: %w", filePath, err)
		}
		readme, ok, err := self.ParseFile(absPath, config)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, readme)
		}
	}
	return result, nil
}

func (self *Parser) IgnoreReadme(filePath string, config *ParseConfig) bool {
	name := filepath.Base(filePath)
	if name != config.ReadmeName {
		return true
	}
	ignore := filepath.Join(config.RootPath, config.ReadmeName)
	if ignore == filePath {
		return true
	}
	return false
}

func (self *Parser) ParseDirectoryWalk(dirPath string, config *ParseConfig) ([]*Readme, error) {
	result := []*Readme{}
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if self.IgnoreReadme(path, config) {
			return nil
		}
		readme, ok, err := self.ParseFile(path, config)
		if err != nil {
			return err
		}
		if ok {
			result = append(result, readme)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse directory %s: %w", dirPath, err)
	}
	return result, nil
}

func (self *Parser) ParseFiles(filePaths []string, config *ParseConfig) ([]*Readme, error) {
	result := []*Readme{}
	for _, path := range filePaths {
		readme, ok, err := self.ParseFile(path, config)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, readme)
		}
	}
	return result, nil
}

func (self *Parser) ParseFile(filePath string, config *ParseConfig) (*Readme, bool, error) {
	relativePath, err := filepath.Rel(config.RootPath, filePath)
	if err != nil {
		return nil, false, fmt.Errorf("could not create relative path of %s: %w", filePath, err)
	}
	dirRelative := filepath.Dir(relativePath)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false, fmt.Errorf("could not read file %s: %w", filePath, err)
	}
	segments := strings.Split(string(fileContent), "---")
	if len(segments) < 3 {
		return nil, false, nil
	}
	var data map[string]any
	err = yaml.Unmarshal([]byte(segments[1]), &data)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal yaml in file %s: %w", filePath, err)
	}
	return &Readme{
		Path:             filePath,
		PathRelative:     relativePath,
		Data:             data,
		Dir:              filepath.Dir(filePath),
		DirRelative:      dirRelative,
		DirRelativeSplit: strings.Split(dirRelative, string(os.PathSeparator)),
		DirName:          filepath.Base(dirRelative),
	}, true, nil
}
