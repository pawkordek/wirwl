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

const logFileName = "wirwl.log"

type Config struct {
	defaultAppDataDirPath string
	defaultConfigDirPath  string
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
	homeDirPath, err := getCurrentUserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(homeDirPath, ".local", "share", "wirwl")
}

func getDefaultConfigDirPath() string {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(userConfigDirPath, "wirwl")
}

func (config *Config) Load() {
	if config.ConfigDirPath == "" {
		config.ConfigDirPath = config.defaultConfigDirPath
	}
	configFilePath := filepath.Join(config.ConfigDirPath, "wirwl.cfg")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		config.AppDataDirPath = config.defaultAppDataDirPath
	}
	config.readConfigFromConfigFilePath(configFilePath)
}

func getDefaultConfigWithConfigPathIn(configPath string) Config {
	return Config{
		AppDataDirPath: getDefaultAppDataDirPath(),
		ConfigDirPath:  configPath,
	}
}

func (config *Config) readConfigFromConfigFilePath(path string) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = toml.Decode(string(fileData), &config)
	if err != nil {
		log.Fatal(err)
	}
}

func (config *Config) setupLoggingIn(path string) {
	if path == "" {
		path = config.defaultAppDataDirPath
	}
	data.CreateDirIfNotExist(path)
	logFile, err := os.OpenFile(path+logFileName, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)
}

func (config *Config) loadDataProviderIn(path string) data.Provider {
	if path == "" {
		path = config.defaultAppDataDirPath
	}
	data.CreateDirIfNotExist(path)
	return data.NewBoltProvider(filepath.Join(path, "data.db"))
}

func getCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

func (config *Config) save() {
	data.CreateDirIfNotExist(config.ConfigDirPath)
	configFilePath := filepath.Join(config.ConfigDirPath, "wirwl.cfg")
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
