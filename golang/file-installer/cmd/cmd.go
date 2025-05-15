package cmd

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	cobra "github.com/spf13/cobra"

	"git.alwaldend.com/src/golang/utils"
)

type configType struct {
	install_files  []string
	default_target string
	fileMode       fs.FileMode
}

var config = &configType{fileMode: fs.FileMode(int(0777))}

var rootCmd = &cobra.Command{
	Use:   "file-installer",
	Short: "File installer",
	Long:  "File installer",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var installCmd = &cobra.Command{
	Use:   "install [files-to-install...]",
	Short: "Install files",
	Args:  cobra.ArbitraryArgs,
	Long:  "Install files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return installArgs(args, config.default_target, config.fileMode)
	},
}

func installArgs(files []string, default_target string, directoryFileMode os.FileMode) error {
	var waitGroup sync.WaitGroup
	type resultType struct {
		index int
		path  string
		err   error
	}
	results := make(chan *resultType, len(files))
	for i, file := range files {
		waitGroup.Add(1)
		go func(file string) {
			defer waitGroup.Done()
			err := installArg(file, default_target, directoryFileMode)
			results <- &resultType{i, file, err}
		}(file)
	}
	waitGroup.Wait()

	count := 0
	for range len(files) {
		result := <-results
		if result.err != nil {
			count += 1
			log.Printf(
				"could not install target %d (%s): %v\n",
				result.index, result.path, result.err,
			)
		}
	}
	if count > 0 {
		return fmt.Errorf("could not install %d target(s)", count)
	}
	return nil
}

func installArg(filePath string, default_target string, directoryFileMode os.FileMode) error {
	target := default_target
	fileObj, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if fileObj.IsDir() {
		return installDirectory(filePath, target, directoryFileMode)
	}
	extension := filepath.Ext(filePath)
	if extension == ".tar" {
		return installTar(filePath, target)
	} else {
		return fmt.Errorf("file %s has invalid extension: %s", filePath, extension)
	}
}

func installTar(sourcePath string, targetDirectory string) error {
	log.Printf("installing tar %s to %s", sourcePath, targetDirectory)
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func() {
		err := sourceFile.Close()
		if err != nil {
			log.Printf("could not close %v: %s\n", sourceFile, err)
		}
	}()
	tarReader := tar.NewReader(sourceFile)
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		err = installTarHeader(header, tarReader, targetDirectory)
		if err != nil {
			return fmt.Errorf("could not install tar header %s: %w", header.Name, err)
		}
	}
}

func installTarHeader(header *tar.Header, archive io.Reader, targetDirectory string) error {
	switch header.Typeflag {
	case tar.TypeDir:
		return nil
	case tar.TypeReg:
		targetPath := filepath.Join(targetDirectory, header.Name)
		targetBytes, err := utils.StreamToByte(archive)
		if err != nil {
			return err
		}
		err = os.WriteFile(targetPath, targetBytes, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
	}
	return nil
}

func installDirectory(sourceDirectory string, targetDirectory string, directoryFileMode os.FileMode) error {
	log.Printf("installing directory %s to %s", sourceDirectory, targetDirectory)
	return filepath.WalkDir(sourceDirectory, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(sourceDirectory, path)
		if err != nil {
			return err
		}
		return installFile(path, filepath.Join(targetDirectory, relPath), directoryFileMode)
	})
}

func installFile(sourcePath string, targetPath string, directoryFileMode os.FileMode) error {
	os.MkdirAll(filepath.Dir(targetPath), directoryFileMode)
	log.Printf("%s -> %s\n", sourcePath, targetPath)
	os.Symlink(sourcePath, targetPath)
	return nil
}

func init() {
	installCmd.PersistentFlags().StringVarP(
		&config.default_target,
		"target",
		"t",
		"",
		"default target for all files",
	)
	rootCmd.AddCommand(installCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
