package log

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogInfo(t *testing.T) {
	writeResult, cleanup := setupLoggerTesting()
	defer cleanup()
	Info("This is an information")
	assert.Contains(t, writeResult.value, "INFO: This is an information")
}

func TestLogError(t *testing.T) {
	writeResult, cleanup := setupLoggerTesting()
	defer cleanup()
	err := errors.New("This is an error")
	Error(err)
	assert.Contains(t, writeResult.value, "ERROR: This is an error")
	errorMsgWithStack := fmt.Sprintf("%+v", err)
	assert.Contains(t, writeResult.value, errorMsgWithStack)
}
