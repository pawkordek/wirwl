package data

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func DeleteFile(path string) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DeleteAllInDir(dirPath string) {
	subDirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, subDir := range subDirs {
		err = os.RemoveAll(filepath.Join(dirPath, subDir.Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CopyFile(sourcePath string, destinationPath string) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		log.Fatal(err)
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		log.Fatal(err)
	}
}
