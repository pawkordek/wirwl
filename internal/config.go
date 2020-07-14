package wirwl

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
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

func (config *Config) load() error {
	if config.ConfigDirPath == "" {
		config.ConfigDirPath = config.defaultConfigDirPath
	}
	config.configFilePath = filepath.Join(config.ConfigDirPath, appName+".cfg")
	if _, err := os.Stat(config.configFilePath); os.IsNotExist(err) {
		config.AppDataDirPath = config.defaultAppDataDirPath
		return nil
	} else {
		return config.readConfigFromConfigFile()
	}
}

func (config *Config) readConfigFromConfigFile() error {
	fileData, err := ioutil.ReadFile(config.configFilePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the config file in path "+config.configFilePath)
	}
	_, err = toml.Decode(string(fileData), &config)
	if err != nil {
		return errors.Wrap(err, "Failed to decode the config from the file in "+config.configFilePath+". File data: "+string(fileData))
	}
	return nil
}

func (config *Config) loadDefaults() {
	config.ConfigDirPath = config.defaultConfigDirPath
	config.AppDataDirPath = config.defaultAppDataDirPath
}

func (config *Config) loadDataProvider() data.Provider {
	return data.NewBoltProvider(filepath.Join(config.AppDataDirPath, "data.db"))
}

func getCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get the current user")
	}
	return currentUser.HomeDir, nil
}

func (config *Config) save() error {
	err := data.CreateDirIfNotExist(config.ConfigDirPath)
	if err != nil {
		return errors.Wrap(err, "Failed to save the config because config directory in "+config.ConfigDirPath+" could not be created")
	}
	configFile, err := os.OpenFile(config.configFilePath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return errors.Wrap(err, "Failed to save the config file because config file in "+config.configFilePath+" could not be opened")
	}
	err = toml.NewEncoder(configFile).Encode(config)
	if err != nil {
		return errors.Wrap(err, "Failed to save the config file because encoding failed. Config data: "+config.String())
	}
	err = configFile.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close the config file in "+config.configFilePath+" after saving")
	}
	return nil
}

func (config Config) String() string {
	return fmt.Sprintf("%#v", config)
}
