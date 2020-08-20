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
	keymap := make(map[Action]string)
	keymap[testAction] = "Q"
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyQ)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionsInKeySequences(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]string)
	keymap[testAction] = "Z,X"
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
	keymap := make(map[Action]string)
	keymap[testAction] = "T"
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
	keymap := make(map[Action]string)
	keymap[testAction] = "L"
	inputHandler := NewInputHandler(keymap)
	inputHandler.Handle("", fyne.KeyL)
}

func TestThatKeySequenceDoesNotWorkIfCertainTimePasses(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]string)
	keymap[testAction] = "U,F"
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.Handle("", fyne.KeyF)
	assert.False(t, functionExecuted)
}

func TestThatKeySequenceWorksAfterPressingKeyPausingAndPressingNextKey(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]string)
	keymap[testAction] = "U,F"
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.Handle("", fyne.KeyF)
	inputHandler.Handle("", fyne.KeyU)
	inputHandler.Handle("", fyne.KeyF)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerExecutesSingleKeyActionIfKeyCombinationWasPressedAndTimePassed(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]string)
	keymap[testAction] = "T"
	keymap[emptyAction] = "Q,Y"
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
	keymap := make(map[Action]string)
	keymap[testAction] = "Q"
	keymap[testAction2] = "J,Q"
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, singleKeyFunc)
	inputHandler.BindFunctionToAction("", testAction2, combinationFunc)
	inputHandler.Handle("", fyne.KeyJ)
	inputHandler.Handle("", fyne.KeyQ)
	assert.True(t, combinationActionExecuted)
	assert.False(t, singleKeyActionExecuted)
}

func TestThatInputHandlerAllowsToUseTheSameKeyForTwoVariousActionsIfTheyAreForDifferentCaller(t *testing.T) {
	testActionExecuted := false
	testAction2Executed := false
	testActionFunc := func() { testActionExecuted = true }
	testAction2Func := func() { testAction2Executed = true }
	keymap := make(map[Action]string)
	keymap[testAction] = "C"
	keymap[testAction2] = "C"
	inputHandler := NewInputHandler(keymap)
	inputHandler.BindFunctionToAction("first caller", testAction, testActionFunc)
	inputHandler.BindFunctionToAction("second caller", testAction2, testAction2Func)
	inputHandler.Handle("first caller", fyne.KeyC)
	assert.False(t, testAction2Executed)
	assert.True(t, testActionExecuted)
	inputHandler.Handle("second caller", fyne.KeyC)
	assert.True(t, testAction2Executed)
}
