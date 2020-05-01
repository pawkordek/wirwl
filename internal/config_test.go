package wirwl

import (
	fyneTest "fyne.io/fyne/test"
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
	NewApp(testDbCopyPath)
	_, err := os.Stat(loggingFilePath)
	assert.Nil(t, err)
}

func TestThatDefaultConfigGetsLoadedIfNoConfigExists(t *testing.T) {
	data.DeleteFile(defaultConfigPath)
	app := NewApp(testDbCopyPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app, cleanup := setupAppForTesting()
	defer cleanup()
	app.config.DataDbPath = defaultConfigPath
}
