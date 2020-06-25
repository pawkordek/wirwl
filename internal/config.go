package wirwl

import (
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"wirwl/internal/data"
)

const appName = "wirwl"
const logFileName = appName + ".log"

type Config struct {
	defaultAppDataDirPath string
	defaultConfigDirPath  string
	configFilePath        string
	AppDataDirPath        string
	ConfigDirPath         string
}

func NewConfig(configDirPath string) Config {
	config := Config{ConfigDirPath: configDirPath}
	config.setupDefaultDirPaths()
	return config
}

func (config *Config) setupDefaultDirPaths() {
	config.defaultConfigDirPath = getDefaultConfigDirPath()
	config.defaultAppDataDirPath = getDefaultAppDataDirPath()
}

func getDefaultAppDataDirPath() string {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return filepath.Join(xdgDataHome, appName)
	} else {
		homeDirPath, err := getCurrentUserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		return filepath.Join(homeDirPath, ".local", "share", appName)
	}
}

func getDefaultConfigDirPath() string {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(userConfigDirPath, appName)
}

func (config *Config) load() {
	if config.ConfigDirPath == "" {
		config.ConfigDirPath = config.defaultConfigDirPath
	}
	config.configFilePath = filepath.Join(config.ConfigDirPath, appName+".cfg")
	if _, err := os.Stat(config.configFilePath); os.IsNotExist(err) {
		config.AppDataDirPath = config.defaultAppDataDirPath
	} else {
		config.readConfigFromConfigFile()
	}
}

func (config *Config) readConfigFromConfigFile() {
	fileData, err := ioutil.ReadFile(config.configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = toml.Decode(string(fileData), &config)
	if err != nil {
		log.Fatal(err)
	}
}

func (config *Config) setupLogger() {
	logFile, err := os.OpenFile(filepath.Join(config.AppDataDirPath, logFileName), os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)
}

func (config *Config) loadDataProvider() data.Provider {
	return data.NewBoltProvider(filepath.Join(config.AppDataDirPath, "data.db"))
}

func getCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

func (config *Config) save() {
	err := data.CreateDirIfNotExist(config.ConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	configFile, err := os.OpenFile(config.configFilePath, os.O_CREATE|os.O_WRONLY, 0700)
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
