package wirwl

import (
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"os"
	"os/user"
	"wirwl/internal/data"
)

var defaultAppDataPath string
var defaultConfigPath string
var defaultConfigFilePath = defaultConfigPath + "wirwl.cfg"
var defaultDataDbPath = defaultAppDataPath + "data.db"

const logFileName = "wirwl.log"

var defaultConfig = Config{
	DataDbPath:    defaultDataDbPath,
	ConfigDirPath: defaultConfigPath,
}

type Config struct {
	DataDbPath    string
	ConfigDirPath string
}

func init() {
	homeDirPath, err := getCurrentUserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultAppDataPath = homeDirPath + "/.local/share/wirwl/"
	defaultConfigPath = homeDirPath + "/.config/wirwl/"
}

func loadConfigFromDir(configDirPath string) Config {
	if configDirPath == "" {
		configDirPath = defaultConfigPath
	}
	createDirIfNotExist(configDirPath)
	if _, err := os.Stat(defaultConfigFilePath); os.IsNotExist(err) {
		return defaultConfig
	}
	return Config{DataDbPath: "", ConfigDirPath: configDirPath}
}

func setupLoggingIn(path string) {
	if path == "" {
		path = defaultAppDataPath
	}
	createDirIfNotExist(path)
	logFile, err := os.OpenFile(path+logFileName, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)
}

func loadDataProviderIn(path string) data.Provider {
	if path == "" {
		path = defaultAppDataPath
	}
	createDirIfNotExist(path)
	return data.NewBoltProvider(path + "data.db")
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

func (config *Config) save() {
	configFilePath := config.ConfigDirPath + "wirwl.cfg"
	configFile, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal("Failed to save the config file due to an error", err)
	}
	err = toml.NewEncoder(configFile).Encode(config)
	if err != nil {
		log.Fatal("Failed to save the config file due to an error", err)
	}
	err = configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}
