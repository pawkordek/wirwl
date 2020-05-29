package wirwl

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"wirwl/internal/data"
)

func TestThatAppDataDirGetsCreatedWhenApplicationLaunches(t *testing.T) {
	_, cleanup := setupAndRunAppAsIfRunForFirstTime()
	defer cleanup()
	_, err := os.Stat(defaultTestAppDataDirPath)
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
	config := NewConfig(testConfigDirPath)
	config.load()
	assert.Equal(t, "some db path", config.AppDataDirPath)
	assert.Equal(t, testConfigDirPath, config.ConfigDirPath)
}

func TestThatDefaultConfigPathIsUsedIfConfigIsCreatedWithEmptyPath(t *testing.T) {
	config := NewConfig("")
	config.defaultConfigDirPath = defaultTestConfigDirPath
	config.load()
	assert.Equal(t, defaultTestConfigDirPath, config.ConfigDirPath)
}

func TestThatDefaultConfigWithProvidedConfigPathGetsLoadedIfConfigFileDoesNotExist(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer data.DeleteFile(tmpDir)
	config := NewConfig(tmpDir)
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
