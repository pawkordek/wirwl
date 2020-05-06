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
	defaultAppDataPath = "../testdata/default/"
	logFilePath := defaultAppDataPath + logFileName
	data.DeleteFile(defaultAppDataPath)
	_, cleanup := setupAppForTestingWithNoPathsProvided()
	defer cleanup()
	_, err := os.Stat(logFilePath)
	assert.Nil(t, err)
	data.DeleteFile(defaultAppDataPath)
}

func TestThatDefaultConfigGetsLoadedIfNoConfigExists(t *testing.T) {
	data.DeleteFile(defaultConfigPath)
	app, cleanup := setupAppForTestingWithDefaultTestingPaths()
	defer cleanup()
	app.config.DataDbPath = defaultConfigPath
}
