package wirwl

import (
	"fmt"
	"fyne.io/fyne"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"wirwl/internal/data"
	"wirwl/internal/input"
)

const appName = "wirwl"
const configFileName = appName + ".cfg"
const logFileName = appName + ".log"

type Config struct {
	AppDataDirPath string
	ConfigDirPath  string
	Keymap         map[input.Action]input.KeyCombination
}

/*As TOML can't encode/decode maps that contain something else than strings, a helper struct is needed to convert
before encoding/decoding.
*/
type encodableDecodableConfig struct {
	AppDataDirPath string
	ConfigDirPath  string
	Keymap         map[string]string
}

func NewConfig(configDirPath string) Config {
	config := Config{ConfigDirPath: configDirPath, Keymap: map[input.Action]input.KeyCombination{}}
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
	decodedConfig := encodableDecodableConfig{}
	_, err = toml.Decode(string(fileData), &decodedConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to decode the config from the file in "+config.ConfigFilePath()+". File data: \n"+string(fileData))
	}
	config.readDataFromDecodedConfig(decodedConfig)
	return nil
}

func (config *Config) readDataFromDecodedConfig(decodedConfig encodableDecodableConfig) {
	config.AppDataDirPath = decodedConfig.AppDataDirPath
	config.ConfigDirPath = decodedConfig.ConfigDirPath
	config.Keymap = convertStringKeymapToFormatUsableByConfig(decodedConfig.Keymap)
}

func convertStringKeymapToFormatUsableByConfig(stringKeymap map[string]string) map[input.Action]input.KeyCombination {
	properKeymap := make(map[input.Action]input.KeyCombination)
	for action, keyCombination := range stringKeymap {
		properKeymap[input.Action(action)] = input.KeyCombinationFromString(keyCombination)
	}
	return properKeymap
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
	config.Keymap[input.SelectNextTabAction] = input.SingleKeyCombination(fyne.KeyL)
	config.Keymap[input.SelectPreviousTabAction] = input.SingleKeyCombination(fyne.KeyH)
	config.Keymap[input.SaveChangesAction] = input.TwoKeyCombination(fyne.KeyS, fyne.KeyY)
	config.Keymap[input.AddEntryTypeAction] = input.TwoKeyCombination(fyne.KeyT, fyne.KeyI)
	config.Keymap[input.RemoveEntryTypeAction] = input.TwoKeyCombination(fyne.KeyT, fyne.KeyD)
	config.Keymap[input.EditCurrentEntryTypeAction] = input.TwoKeyCombination(fyne.KeyT, fyne.KeyE)
	config.Keymap[input.MoveDownAction] = input.SingleKeyCombination(fyne.KeyJ)
	config.Keymap[input.MoveUpAction] = input.SingleKeyCombination(fyne.KeyK)
	config.Keymap[input.EnterInputModeAction] = input.SingleKeyCombination(fyne.KeyI)
	config.Keymap[input.ExitInputModeAction] = input.SingleKeyCombination(fyne.KeyEscape)
	config.Keymap[input.ExitTableAction] = input.SingleKeyCombination(fyne.KeySpace)
	config.Keymap[input.ConfirmAction] = input.SingleKeyCombination(fyne.KeyReturn)
	config.Keymap[input.CancelAction] = input.SingleKeyCombination(fyne.KeyEscape)
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
	//Can't pass the keymap as it is as toml can't encode maps that don't use strings as key so properly modified config has to be used
	err = toml.NewEncoder(configFile).Encode(config.madeEncodable())
	if err != nil {
		return errors.Wrap(err, "Failed to save the config file because encoding failed. Config data: "+config.String())
	}
	err = configFile.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close the config file in "+config.ConfigFilePath()+" after saving")
	}
	return nil
}

func (config *Config) madeEncodable() encodableDecodableConfig {
	encodableKeymap := make(map[string]string)
	for action, key := range config.Keymap {
		encodableKeymap[string(action)] = key.String()
	}
	return encodableDecodableConfig{
		AppDataDirPath: config.AppDataDirPath,
		ConfigDirPath:  config.ConfigDirPath,
		Keymap:         encodableKeymap,
	}
}

func (config *Config) ConfigFilePath() string {
	return filepath.Join(config.ConfigDirPath, configFileName)
}

func (config Config) String() string {
	return fmt.Sprintf("%#v", config)
}
