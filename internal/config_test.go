package wirwl

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"wirwl/internal/data"
)

func TestThatLoggingFileWithItsDirGetsCreated(t *testing.T) {
	loggingDir := getLoggingDirForTesting()
	loggingFilePath := loggingDir + "wirwl.log"
	data.DeleteFile(loggingFilePath)
	data.DeleteFile(loggingDir)
	_, cleanup := setupAppForTesting()
	defer cleanup()
	_, err := os.Stat(loggingFilePath)
	assert.Nil(t, err)
}

func TestThatDefaultConfigGetsLoadedIfNoConfigExists(t *testing.T) {
	data.DeleteFile(defaultConfigPath)
	app, cleanup := setupAppForTesting()
	defer cleanup()
	app.config.DataDbPath = defaultConfigPath
}
