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
	NewApp(testDbPath)
	_, err := os.Stat(loggingFilePath)
	assert.Nil(t, err)
}
