package wirwl

import (
	"log"
	"os"
	"os/user"
)

const defaultAppDataPath = "/.local/share/wirwl/"
const defaultConfigPath = defaultAppDataPath + "wirwl.cfg"
const defaultDataDbPath = defaultAppDataPath + "data.db"

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

func setupLogging() {
	loggingDir := getLoggingDir()
	createDirIfNotExist(loggingDir)
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

func createDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}
