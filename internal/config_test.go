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

func TestThatAppDataDirPathDefaultsToXDG_DATA_HOMEIfItIsSet(t *testing.T) {
	err := os.Setenv("XDG_DATA_HOME", "some path")
	if err != nil {
		log.Fatal(err)
	}
	config := NewConfig("")
	err = config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join("some path", "wirwl")
	assert.Equal(t, expectedPath, config.AppDataDirPath)
	err = os.Unsetenv("XDG_DATA_HOME")
	if err != nil {
		log.Fatal(err)
	}
}

func TestThatAppDirDefaultsToLocalShareIfXDG_DATA_HOMEIsNotSet(t *testing.T) {
	config := NewConfig("")
	err := config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join(currentUser.HomeDir, ".local", "share", "wirwl")
	assert.Equal(t, expectedPath, config.AppDataDirPath)
}

func TestThatConfigDirDefaultsToUserConfigDir(t *testing.T) {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	config := NewConfig("")
	err = config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join(userConfigDirPath, "wirwl")
	assert.Equal(t, expectedPath, config.ConfigDirPath)
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
	defer func() {
		err = data.DeleteDirWithContents(tmpDir)
		if err != nil {
			log.Fatal(err)
		}
	}()
	config := NewConfig(tmpDir)
	config.defaultConfigDirPath = defaultTestConfigDirPath
	config.defaultAppDataDirPath = defaultTestAppDataDirPath
	config.load()
	assert.Equal(t, defaultTestAppDataDirPath, config.AppDataDirPath)
	assert.Equal(t, tmpDir, config.ConfigDirPath)
}
