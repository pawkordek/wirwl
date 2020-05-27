package wirwl

import (
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"wirwl/internal/data"
)

var defaultAppDataDirPath string
var defaultConfigDirPath string

const logFileName = "wirwl.log"

type Config struct {
	AppDataDirPath string
	ConfigDirPath  string
}

func init() {
	homeDirPath, err := getCurrentUserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultAppDataDirPath = homeDirPath + "/.local/share/wirwl/"
	defaultConfigDirPath = getDefaultConfigDirPath()
}

func getDefaultConfigDirPath() string {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(userConfigDirPath, "wirwl")
}

func LoadConfigFromDir(configDirPath string) Config {
	if configDirPath == "" {
		configDirPath = defaultConfigDirPath
	}
	configFilePath := configDirPath + "wirwl.cfg"
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return getDefaultConfigWithConfigPathIn(configDirPath)
	}
	return readConfigFromFileIn(configFilePath)
}

func getDefaultConfigWithConfigPathIn(configPath string) Config {
	return Config{
		AppDataDirPath: defaultAppDataDirPath,
		ConfigDirPath:  configPath,
	}
}

func readConfigFromFileIn(path string) Config {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{}
	_, err = toml.Decode(string(fileData), &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func setupLoggingIn(path string) {
	if path == "" {
		path = defaultAppDataDirPath
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

func loadDataProviderIn(path string) data.Provider {
	if path == "" {
		path = defaultAppDataDirPath
	}
	data.CreateDirIfNotExist(path)
	return data.NewBoltProvider(path + "data.db")
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
