package wirwl

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
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

func NewConfig(configDirPath string) (Config, error) {
	config := Config{ConfigDirPath: configDirPath}
	err := config.setupDefaultDirPaths()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (config *Config) setupDefaultDirPaths() error {
	defaultConfigDirPath, err := getDefaultConfigDirPath()
	if err != nil {
		return err
	}
	config.defaultConfigDirPath = defaultConfigDirPath
	defaultAppDataDirPath, err := getDefaultAppDataDirPath()
	if err != nil {
		return err
	} else {
		config.defaultAppDataDirPath = defaultAppDataDirPath
	}
	return nil
}

func getDefaultAppDataDirPath() (string, error) {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return filepath.Join(xdgDataHome, appName), nil
	} else {
		homeDirPath, err := getCurrentUserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "Failed to setup default app data dir path")
		}
		return filepath.Join(homeDirPath, ".local", "share", appName), nil
	}
}

func getDefaultConfigDirPath() (string, error) {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Wrap(err, "Failed to setup default user config dir path")
	}
	return filepath.Join(userConfigDirPath, appName), nil
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
