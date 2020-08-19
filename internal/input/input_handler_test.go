package input

import (
	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testAction Action = "TEST_ACTION"
const testAction2 Action = "TEST_ACION2"
const emptyAction Action = "EMPTY_ACTION"

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

func TestThatKeySequenceDoesNotWorkIfCertainTimePasses(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[string]Action)
	keymap["U,F"] = testAction
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.Handle("", fyne.KeyF)
	assert.False(t, functionExecuted)
	inputHandler.Handle("", fyne.KeyU)
	inputHandler.Handle("", fyne.KeyF)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerExecutesSingleKeyActionIfKeyCombinationWasPressedAndTimePassed(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[string]Action)
	keymap["T"] = testAction
	keymap["Q,Y"] = emptyAction
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.BindFunctionToAction("", emptyAction, func() {})
	inputHandler.Handle("", fyne.KeyQ)
	inputHandler.Handle("", fyne.KeyY)
	time.Sleep(2 * time.Second)
	inputHandler.Handle("", fyne.KeyT)
	assert.True(t, functionExecuted)
}

func TestThatCombinationActionGetsExecutedInsteadOfSingleKeyActionIfBothCouldHappen(t *testing.T) {
	combinationActionExecuted := false
	singleKeyActionExecuted := false
	combinationFunc := func() { combinationActionExecuted = true }
	singleKeyFunc := func() { singleKeyActionExecuted = true }
	keymap := make(map[string]Action)
	keymap["Q"] = testAction
	keymap["J,Q"] = testAction2
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, singleKeyFunc)
	inputHandler.BindFunctionToAction("", testAction2, combinationFunc)
	inputHandler.Handle("", fyne.KeyJ)
	inputHandler.Handle("", fyne.KeyQ)
	assert.True(t, combinationActionExecuted)
	assert.False(t, singleKeyActionExecuted)
}
