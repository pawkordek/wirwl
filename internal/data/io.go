package data

import (
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func DeleteDirWithContents(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrap(err, "An error occurred when deleting folder in path "+path+" and it's contents")
		}
	}
	return nil
}

func DeleteAllInDir(dirPath string) {
	DeleteAllInDirExceptForDirs(dirPath, "")
}

func DeleteAllInDirExceptForDirs(dirPath string, excludedDirsNames ...string) error {
	subDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return errors.Wrap(err, "An error occurred when reading directory "+dirPath)
	}
	for _, subDir := range subDirs {
		shouldBeSkipped := false
		for _, excludedDirName := range excludedDirsNames {
			if subDir.Name() == excludedDirName {
				shouldBeSkipped = true
			}
		}
		if !shouldBeSkipped {
			err = DeleteDirWithContents(filepath.Join(dirPath, subDir.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(sourcePath string, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrap(err, "An error occurred when trying to open source file in path "+sourcePath+" to copy")
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return errors.Wrap(err, "An error occurred when trying to open a copy destination file in path "+destinationPath)
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return errors.Wrap(err, "An error occurred when copying files from "+sourcePath+" to "+destinationPath)
	}
	return nil
}

func CreateDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return errors.Wrap(err, "An error occurred when creating a directory in path "+path)
		}
	}
	return nil
}
