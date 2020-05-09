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
	data.DeleteFile(defaultAppDataPath)
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
	data.DeleteFile(defaultAppDataPath)
}

func TestThatDefaultConfigGetsLoadedIfNoConfigExists(t *testing.T) {
	defaultConfigPath = defaultTestConfigDirPath
	data.DeleteFile(defaultConfigPath)
	app, cleanup := setupAppForTestingWithNoPathsProvided()
	defer cleanup()
	app.config.DataDbPath = defaultConfigPath
}
