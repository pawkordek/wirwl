package wirwl

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"wirwl/internal/data"
)

func TestThatLoggingFileWithItsDirGetsCreatedIfAppDataDirIsProvided(t *testing.T) {
	logFilePath := testAppDataDirPath + "wirwl.log"
	data.DeleteFile(logFilePath)
	_, cleanup := setupAppForTestingWithDefaultTestingPaths()
	defer cleanup()
	_, err := os.Stat(logFilePath)
	assert.Nil(t, err)
}

func TestThatLoggingFileWithItsDirGetsCreatedInDefaultPathIfAppDataDirIsNotProvided(t *testing.T) {
	defaultAppDataPath = defaultTestAppDataDirPath
	logFilePath := defaultAppDataPath + logFileName
	data.DeleteFile(defaultAppDataPath)
	_, cleanup := setupAppForTestingWithNoPathsProvided()
	defer cleanup()
	_, err := os.Stat(logFilePath)
	assert.Nil(t, err)
}

func TestThatDbFileWithItsDirGetsCreatedIfAppDataDirIsProvided(t *testing.T) {
	dbFilePath := testAppDataDirPath + "data.db"
	data.DeleteFile(dbFilePath)
	_, cleanup := setupAppForTestingWithDefaultTestingPaths()
	defer cleanup()
	_, err := os.Stat(dbFilePath)
	assert.Nil(t, err)
}

func TestThatDbFileWithItsDirGetsCreatedInDefaultPathIfAppDataDirIsNotProvided(t *testing.T) {
	defaultAppDataPath = defaultTestAppDataDirPath
	dbFilePath := defaultAppDataPath + "data.db"
	data.DeleteFile(dbFilePath)
	_, cleanup := setupAppForTestingWithNoPathsProvided()
	defer cleanup()
	_, err := os.Stat(dbFilePath)
	assert.Nil(t, err)
}

func TestThatConfigGetsLoadedIfItExists(t *testing.T) {
	createCorrectWirwlConfigFileForLoadingInPath(testConfigDirPath)
	app, cleanup := setupAppForTestingWithDefaultTestingPaths()
	defer cleanup()
	assert.Equal(t, app.config.AppDataDirPath, "some db path")
	assert.Equal(t, app.config.ConfigDirPath, testConfigDirPath)
}

func TestThatDefaultConfigWithProvidedConfigPathGetsLoadedIfConfigFileDoesNotExist(t *testing.T) {
	defaultConfigPath = defaultTestConfigDirPath
	data.DeleteFile(defaultConfigPath)
	app, cleanup := setupAppForTestingWithDefaultTestingPaths()
	defer cleanup()
	assert.Equal(t, app.config.AppDataDirPath, defaultAppDataPath)
	assert.Equal(t, app.config.ConfigDirPath, defaultConfigPath)
}

func TestThatCorrectConfigFileGetsWrittenToDiskAfterApplicationExits(t *testing.T) {
	data.DeleteFile(testConfigDirPath)
	_, cleanup := setupAppForTestingWithDefaultTestingPaths()
	defer cleanup()
	createCorrectSavedWirwlConfigFileInPath(testConfigDirPath)
	assert.True(t, areFilesInPathsTheSame(testConfigDirPath+"wirwl.cfg", testConfigDirPath+"wirwl_correct.cfg"))
}
