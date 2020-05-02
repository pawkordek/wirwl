package wirwl

import (
	"io"
	"log"
	"os"
	"os/user"
)

const defaultAppDataPath = "/.local/share/wirwl/"
const defaultConfigPath = defaultAppDataPath + "wirwl.cfg"
const defaultDataDbPath = defaultAppDataPath + "data.db"
const logFileName = "wirwl.log"

var defaultConfig = Config{
	DataDbPath: defaultDataDbPath,
}

type Config struct {
	DataDbPath string
}

func loadConfigFromDir(configDirPath string) Config {
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		return defaultConfig
	}
	return Config{DataDbPath: ""}
}

func setupLoggingIn(path string) {
	logFile, err := os.OpenFile(path+logFileName, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)
}

func getCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

func createDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}
