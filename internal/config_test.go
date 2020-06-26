package wirwl

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

func TestThatDefaultAppDataDirPathIsSetToXDG_DATA_HOMEIfItIsSet(t *testing.T) {
	err := os.Setenv("XDG_DATA_HOME", "some path")
	if err != nil {
		log.Fatal(err)
	}
	config, err := NewConfig("")
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join("some path", "wirwl")
	assert.Equal(t, expectedPath, config.defaultAppDataDirPath)
	err = os.Unsetenv("XDG_DATA_HOME")
	if err != nil {
		log.Fatal(err)
	}
}

func TestThatDefaultAppDirDefaultsToLocalShareIfXDG_DATA_HOMEIsNotSet(t *testing.T) {
	config, err := NewConfig("")
	if err != nil {
		log.Fatal(err)
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join(currentUser.HomeDir, ".local", "share", "wirwl")
	assert.Equal(t, expectedPath, config.defaultAppDataDirPath)
}

func TestThatAppDataDirGetsCreatedWhenApplicationLaunches(t *testing.T) {
	_, cleanup := setupAndRunAppAsIfRunForFirstTime()
	defer cleanup()
	_, err := os.Stat(defaultTestAppDataDirPath)
	assert.Nil(t, err)
}

func TestThatLoggingFileWithItsDirGetsCreatedWhenAppLaunchesForFirstTime(t *testing.T) {
	logFilePath := defaultTestAppDataDirPath + "wirwl.log"
	_, cleanup := setupAndRunAppAsIfRunForFirstTime()
	defer cleanup()
	_, err := os.Stat(logFilePath)
	assert.Nil(t, err)
}

func TestThatLoggingFileWithItsDirGetsCreatedInAppDataDirFromConfig(t *testing.T) {
	logFilePath := testAppDataDirPath + "wirwl.log"
	_, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	_, err := os.Stat(logFilePath)
	assert.Nil(t, err)
}

func TestThatDbFileWithItsDirGetsCreatedInAppDataDirFromConfig(t *testing.T) {
	dbFilePath := testAppDataDirPath + "data.db"
	_, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	_, err := os.Stat(dbFilePath)
	assert.Nil(t, err)
}

func TestThatConfigGetsLoadedIfItExists(t *testing.T) {
	createCorrectWirwlConfigFileForLoadingInPath(testConfigDirPath)
	config, err := NewConfig(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	config.load()
	assert.Equal(t, "some db path", config.AppDataDirPath)
	assert.Equal(t, testConfigDirPath, config.ConfigDirPath)
}

func TestThatDefaultConfigPathIsUsedIfConfigIsCreatedWithEmptyPath(t *testing.T) {
	config, err := NewConfig("")
	if err != nil {
		log.Fatal(err)
	}
	config.defaultConfigDirPath = defaultTestConfigDirPath
	config.load()
	assert.Equal(t, defaultTestConfigDirPath, config.ConfigDirPath)
}

func TestThatDefaultConfigWithProvidedConfigPathGetsLoadedIfConfigFileDoesNotExist(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = data.DeleteDirWithContents(tmpDir)
		if err != nil {
			log.Fatal(err)
		}
	}()
	config, err := NewConfig(tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	config.defaultConfigDirPath = defaultTestConfigDirPath
	config.defaultAppDataDirPath = defaultTestAppDataDirPath
	config.load()
	assert.Equal(t, defaultTestAppDataDirPath, config.AppDataDirPath)
	assert.Equal(t, tmpDir, config.ConfigDirPath)
}

func TestThatCorrectConfigFileGetsWrittenToDiskAfterApplicationExits(t *testing.T) {
	_, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	createCorrectSavedWirwlConfigFileInPath(testConfigDirPath)
	assert.True(t, areFilesInPathsTheSame(testConfigDirPath+"wirwl.cfg", testConfigDirPath+"wirwl_correct.cfg"))
}

func TestThatAfterLoadingDefaultsConfigHasDefaultValues(t *testing.T) {
	config, err := NewConfig(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	config.loadDefaults()
	assert.Equal(t, config.ConfigDirPath, config.defaultConfigDirPath)
	assert.Equal(t, config.AppDataDirPath, config.defaultAppDataDirPath)
}
