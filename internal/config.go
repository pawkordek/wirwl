package wirwl

import (
	"log"
	"os"
	"os/user"
)

func setupLogging() {
	loggingDir := getLoggingDir()
	createLoggingDirIfNotExist(loggingDir)
	logFile, err := os.OpenFile(loggingDir+"wirwl.log", os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
}

func getCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

func getLoggingDir() string {
	userHomeDir, err := getCurrentUserHomeDir()
	if err != nil {
		return "/tmp/wirwl/"
	}
	loggingDir := userHomeDir + "/.local/share/wirwl/"
	return loggingDir
}

func createLoggingDirIfNotExist(loggingDir string) {
	if _, err := os.Stat(loggingDir); os.IsNotExist(err) {
		err := os.Mkdir(loggingDir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}
