package wirwl

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

func TestThatDefaultConfigGetsReturnedIfConfigFileHasUnreadableDataAndConfiguratorHasProperLoadingError(t *testing.T) {
	defer cleanupAfterTestRun()
	err := data.CreateDirIfNotExist(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(testConfigDirPath+"wirwl.cfg", []byte("config with nonsense data"), 0666)
	if err != nil {
		log.Fatal(err)
	}
	configurator := NewAppConfigurator(testConfigDirPath)
	actualConfig, err := configurator.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	expectedConfig := NewConfig(testConfigDirPath)
	err = expectedConfig.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedConfig.ConfigDirPath, actualConfig.ConfigDirPath)
	assert.Equal(t, expectedConfig.AppDataDirPath, actualConfig.AppDataDirPath)
	assert.Equal(t, configurator.loadingErrors[configLoadError], "An error occurred when loading the config file in "+testConfigDirPath+"wirwl.cfg. Default config has been loaded instead.")
}

func TestThatProperDataProviderIsLoaded(t *testing.T) {
	defer cleanupAfterTestRun()
	expectedDataProvider := data.NewBoltProvider(testDbCopyPath)
	configurator := NewAppConfigurator(testConfigDirPath)
	loadedDataProvider := configurator.LoadDataProvider(testDbCopyPath)
	assert.Equal(t, expectedDataProvider, loadedDataProvider)
}

func TestThatNeededDirectoriesWereCreatedAfterPathsSetup(t *testing.T) {
	defer cleanupAfterTestRun()
	config := NewConfig(testConfigDirPath)
	config.AppDataDirPath = testAppDataDirPath
	configurator := NewAppConfigurator(testConfigDirPath)
	err := configurator.SetupNeededPaths(config)
	if err != nil {
		log.Fatal(err)
	}
	assert.DirExists(t, testAppDataDirPath)
}

func TestThatProperErrorIsReturnedIfAnAttemptIsMadeToCreateDirectoriesWithConfigContainingWrongPaths(t *testing.T) {
	defer cleanupAfterTestRun()
	nonsensePath := "/nonsense path"
	config := NewConfig(nonsensePath)
	config.AppDataDirPath = nonsensePath
	configurator := NewAppConfigurator(nonsensePath)
	err := configurator.SetupNeededPaths(config)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to create application directory in /nonsense path. Application must exit")
}

func TestThatLogIsWrittenToFileAfterLoggerIsSetup(t *testing.T) {
	defer cleanupAfterTestRun()
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
