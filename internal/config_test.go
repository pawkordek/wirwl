package wirwl

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"wirwl/internal/data"
)

func TestThatLoggingFileWithItsDirGetsCreated(t *testing.T) {
	logFilePath := testAppDataDirPath + "wirwl.log"
	data.DeleteFile(logFilePath)
	_, cleanup := setupAppForTesting()
	defer cleanup()
	_, err := os.Stat(logFilePath)
	assert.Nil(t, err)
}

func TestThatDefaultConfigGetsLoadedIfNoConfigExists(t *testing.T) {
	data.DeleteFile(defaultConfigPath)
	app, cleanup := setupAppForTesting()
	defer cleanup()
	app.config.DataDbPath = defaultConfigPath
}
