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

const (
	selectNextTabAction     Action = "SELECT_NEXT_TAB"
	selectPreviousTabAction Action = "SELECT_PREVIOUS_TAB"
	saveChangesAction       Action = "SAVE_CHANGES"
)

type Config struct {
	AppDataDirPath string
	ConfigDirPath  string
	Keymap         map[string]Action
}

func NewConfig(configDirPath string) Config {
	config := Config{ConfigDirPath: configDirPath, Keymap: map[string]Action{}}
	return config
}

func (config *Config) load() error {
	if _, err := os.Stat(config.ConfigFilePath()); os.IsNotExist(err) {
		return errors.New("Failed to find config file in path " + config.ConfigFilePath())
	} else {
		return config.readConfigFromConfigFile()
	}
}

func (config *Config) readConfigFromConfigFile() error {
	fileData, err := ioutil.ReadFile(config.ConfigFilePath())
	if err != nil {
		return errors.Wrap(err, "Failed to read the config file in path "+config.ConfigFilePath())
	}
	_, err = toml.Decode(string(fileData), &config)
	if err != nil {
		return errors.Wrap(err, "Failed to decode the config from the file in "+config.ConfigFilePath()+". File data: "+string(fileData))
	}
	return nil
}

func (config *Config) loadDefaults() error {
	defaultConfigDirPath, err := getDefaultConfigDirPath()
	if err != nil {
		return errors.Wrap(err, "An error occurred when trying to load default config directory path in config")
	}
	config.ConfigDirPath = defaultConfigDirPath
	defaultAppDataDirPath, err := getDefaultAppDataDirPath()
	if err != nil {
		return errors.Wrap(err, "An error occurred when trying to load default app directory path in config")
	}
	config.AppDataDirPath = defaultAppDataDirPath
	config.loadDefaultKeymap()
	return nil
}

func getDefaultConfigDirPath() (string, error) {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Wrap(err, "Failed to setup default user config dir path")
	}
	return filepath.Join(userConfigDirPath, appName), nil
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

func getCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get the current user")
	}
	return currentUser.HomeDir, nil
}

func (config *Config) loadDefaultKeymap() {
	config.Keymap["L"] = selectNextTabAction
	config.Keymap["H"] = selectPreviousTabAction
	config.Keymap["S,Y"] = saveChangesAction
}

func (config *Config) save() error {
	err := data.CreateDirIfNotExist(config.ConfigDirPath)
	if err != nil {
		return errors.Wrap(err, "Failed to save the config because config directory in "+config.ConfigDirPath+" could not be created")
	}
	configFile, err := os.OpenFile(config.ConfigFilePath(), os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return errors.Wrap(err, "Failed to save the config file because config file in "+config.ConfigFilePath()+" could not be opened")
	}
	err = toml.NewEncoder(configFile).Encode(config)
	if err != nil {
		return errors.Wrap(err, "Failed to save the config file because encoding failed. Config data: "+config.String())
	}
	err = configFile.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close the config file in "+config.ConfigFilePath()+" after saving")
	}
	return nil
}

func (config *Config) ConfigFilePath() string {
	return filepath.Join(config.ConfigDirPath, appName+".cfg")
}

func (config Config) String() string {
	return fmt.Sprintf("%#v", config)
}
