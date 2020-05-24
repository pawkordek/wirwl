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
	_, err := os.Stat(testAppDataDirPath)
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
	config := LoadConfigFromDir(testConfigDirPath)
	assert.Equal(t, "some db path", config.AppDataDirPath)
	assert.Equal(t, testConfigDirPath, config.ConfigDirPath)
}

func TestThatDefaultConfigWithProvidedConfigPathGetsLoadedIfConfigFileDoesNotExist(t *testing.T) {
	defaultConfigDirPath = defaultTestConfigDirPath
	defaultAppDataDirPath = defaultTestAppDataDirPath
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer data.DeleteFile(tmpDir)
	config := LoadConfigFromDir(tmpDir)
	assert.Equal(t, defaultTestAppDataDirPath, config.AppDataDirPath)
	assert.Equal(t, tmpDir, config.ConfigDirPath)
}

func TestThatCorrectConfigFileGetsWrittenToDiskAfterApplicationExits(t *testing.T) {
	_, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	createCorrectSavedWirwlConfigFileInPath(testConfigDirPath)
	assert.True(t, areFilesInPathsTheSame(testConfigDirPath+"wirwl.cfg", testConfigDirPath+"wirwl_correct.cfg"))
}
