package input

import (
	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testAction Action = "TEST_ACTION"

func TestThatInputHandlerHandlesSingleKeyActions(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[string]Action)
	keymap["Q"] = testAction
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyQ)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionsInKeySequences(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[string]Action)
	keymap["Z,X"] = testAction
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyZ)
	assert.False(t, functionExecuted)
	inputHandler.Handle("", fyne.KeyX)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionForTheCorrectCaller(t *testing.T) {
	firstCallerFunctionExecuted := false
	secondCallerFunctionExecuted := false
	firstCallerFunction := func() { firstCallerFunctionExecuted = true }
	secondCallerFunction := func() { secondCallerFunctionExecuted = true }
	keymap := make(map[string]Action)
	keymap["T"] = testAction
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("firstCaller", testAction, firstCallerFunction)
	inputHandler.BindFunctionToAction("secondCaller", testAction, secondCallerFunction)
	inputHandler.Handle("firstCaller", fyne.KeyT)
	assert.True(t, firstCallerFunctionExecuted)
	assert.False(t, secondCallerFunctionExecuted)
	inputHandler.Handle("secondCaller", fyne.KeyT)
	assert.True(t, secondCallerFunctionExecuted)
}

func TestThatTryingToHandleInputForActionWithoutBindFunctionDoesNotCauseErrors(t *testing.T) {
	keymap := make(map[string]Action)
	keymap["L"] = testAction
	inputHandler := NewInputHandler(keymap)
	inputHandler.Handle("", fyne.KeyL)
}