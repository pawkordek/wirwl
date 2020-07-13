package wirwl

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

func TestThatDefaultConfigGetsReturnedIfConfigFileHasUnreadableDataAndConfiguratorHasProperLoadingError(t *testing.T) {
	err := data.CreateDirIfNotExist(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(testConfigDirPath+"wirwl.cfg", []byte("config with nonsense data"), 0666)
	if err != nil {
		log.Fatal(err)
	}
	configurator := NewAppConfigurator(testConfigDirPath)
	config, err := configurator.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, config.defaultConfigDirPath, config.ConfigDirPath)
	assert.Equal(t, config.defaultAppDataDirPath, config.AppDataDirPath)
	assert.Equal(t, configurator.loadingErrors["config"], "An error occurred when loading the config file in "+testConfigDirPath+"wirwl.cfg. Default config has been loaded instead.")
}

func TestThatProperDataProviderIsLoaded(t *testing.T) {
	expectedDataProvider := data.NewBoltProvider(testDbCopyPath)
	configurator := NewAppConfigurator(testConfigDirPath)
	loadedDataProvider := configurator.LoadDataProvider(testDbCopyPath)
	assert.Equal(t, expectedDataProvider, loadedDataProvider)
}

func TestThatNeededDirectoriesWereCreatedAfterPathsSetup(t *testing.T) {
	config, err := NewConfig(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	config.AppDataDirPath = testAppDataDirPath
	configurator := NewAppConfigurator(testConfigDirPath)
	err = configurator.SetupNeededPaths(config)
	if err != nil {
		log.Fatal(err)
	}
	assert.DirExists(t, testAppDataDirPath)
}

func TestThatProperErrorIsReturnedIfAnAttemptIsMadeToCreateDirectoriesWithConfigContainingWrongPaths(t *testing.T) {
	nonsensePath := "/nonsense path"
	config, err := NewConfig(nonsensePath)
	if err != nil {
		log.Fatal(err)
	}
	config.AppDataDirPath = nonsensePath
	configurator := NewAppConfigurator(nonsensePath)
	err = configurator.SetupNeededPaths(config)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to create application directory in /nonsense path. Application must exit")
}

func TestThatLogIsWrittenToFileAfterLoggerIsSetup(t *testing.T) {
	err := data.CreateDirIfNotExist(testAppDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
	configurator := NewAppConfigurator(testConfigDirPath)
	cleanup := configurator.SetupLoggerIn(testAppDataDirPath)
	defer cleanup()
	expectedTextInLogFile := "Just some text in log file"
	log.Info(expectedTextInLogFile)
	logFileContents, err := ioutil.ReadFile(testAppDataDirPath + "wirwl.log")
	if err != nil {
		log.Fatal(err)
	}
	assert.Contains(t, string(logFileContents), expectedTextInLogFile)
}
