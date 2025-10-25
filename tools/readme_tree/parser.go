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

func NewParser(config *Config, outputter *Outputter) *Parser {
	return &Parser{config: config, outputter: outputter}
}

func (self *Parser) ParseAndOutput(paths []string, config *ParseConfig) error {
	result, err := self.Parse(paths, config)
	if err != nil {
		return err
	}
	return self.outputter.Output(result, config.OutputType)
}

func (self *Parser) Parse(paths []string, config *ParseConfig) (*ParsePathsResult, error) {
	result := &ParsePathsResult{Dir: config.RootPath}
	for _, path := range paths {
		readmes, err := self.parsePath(path, config)
		if err != nil {
			return nil, fmt.Errorf("could not process path %s: %w", path, err)
		}
		result.Readmes = append(result.Readmes, readmes...)
	}
	return result, nil
}

func (self *Parser) parsePath(path string, config *ParseConfig) ([]*Readme, error) {
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
		readmes, err = self.parseDirectory(absPath, config)
		if err != nil {
			return nil, err
		}
	} else {
		readme, ok, err := self.parseFile(absPath, config)
		if err != nil {
			return nil, err
		}
		if ok {
			readmes = []*Readme{readme}
		}
	}
	return readmes, nil
}

func (self *Parser) parseDirectory(dirPath string, config *ParseConfig) ([]*Readme, error) {
	if config.UseGit {
		return self.parseDirectoryGit(dirPath, config)
	}
	return self.parseDirectoryWalk(dirPath, config)
}

func (self *Parser) parseDirectoryGit(dirPath string, config *ParseConfig) ([]*Readme, error) {
	stdout, err := exec.Command("git", "ls-files", dirPath).Output()
	if err != nil {
		return nil, err
	}
	expr, err := regexp.Compile("\r?\n")
	if err != nil {
		return nil, err
	}
	lines := expr.Split(string(stdout), -1)
	result := []*Readme{}
	for _, filePath := range lines {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return nil, fmt.Errorf("could not get absolute path of %s: %w", filePath, err)
		}
		ignore, err := self.ignorePath(filePath, config)
		if err != nil {
			return nil, err
		}
		if ignore {
			continue
		}
		readme, ok, err := self.parseFile(absPath, config)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, readme)
		}
	}
	return result, nil
}

func (self *Parser) ignorePath(file string, config *ParseConfig) (bool, error) {
	name := filepath.Base(file)
	if name != config.ReadmeName {
		return true, nil
	}
	ignore := filepath.Join(config.RootPath, config.ReadmeName)
	if ignore == file || file == config.ReadmeName {
		return true, nil
	}
	for _, exclude := range config.Exclude {
		if exclude == file {
			return true, nil
		}
		rel, err := filepath.Rel(exclude, file)
		if err != nil {
			return true, fmt.Errorf("could not create relative path of %s and %s: %w", exclude, file, err)
		}
		if !strings.Contains(rel, "..") {
			return true, nil
		}
	}
	return false, nil
}

func (self *Parser) parseDirectoryWalk(dirPath string, config *ParseConfig) ([]*Readme, error) {
	result := []*Readme{}
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ignore, err := self.ignorePath(path, config)
		if err != nil || ignore {
			return err
		}
		readme, ignore, err := self.parseFile(path, config)
		if err != nil {
			return err
		}
		if ignore {
			result = append(result, readme)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse directory %s: %w", dirPath, err)
	}
	return result, nil
}

func (self *Parser) parseFile(filePath string, config *ParseConfig) (*Readme, bool, error) {
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
