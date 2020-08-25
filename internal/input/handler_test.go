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

func TestThatInputHandlerDoesNotPanickWhenCallbackFunctionIsNotSet(t *testing.T) {
	keymap := make(map[Action]KeyCombination)
	handler := NewHandler(keymap)
	handler.Handle("", fyne.KeyU)
}

func TestThatInputHandlerHandlesSingleKeyActions(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyQ)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyQ)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionsInKeySequences(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyZ, fyne.KeyX)
	inputHandler := NewHandler(keymap)
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
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyT)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("firstCaller", testAction, firstCallerFunction)
	inputHandler.BindFunctionToAction("secondCaller", testAction, secondCallerFunction)
	inputHandler.Handle("firstCaller", fyne.KeyT)
	assert.True(t, firstCallerFunctionExecuted)
	assert.False(t, secondCallerFunctionExecuted)
	inputHandler.Handle("secondCaller", fyne.KeyT)
	assert.True(t, secondCallerFunctionExecuted)
}

func TestThatTryingToHandleInputForActionWithoutBindFunctionDoesNotCauseErrors(t *testing.T) {
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyL)
	inputHandler := NewHandler(keymap)
	inputHandler.Handle("", fyne.KeyL)
}

func TestThatKeySequenceDoesNotWorkIfCertainTimePasses(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyU, fyne.KeyF)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.Handle("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.Handle("", fyne.KeyF)
	assert.False(t, functionExecuted)
}

func TestThatKeySequenceWorksAfterPressingKeyPausingAndPressingNextKey(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyU, fyne.KeyF)
	inputHandler := NewHandler(keymap)
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
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyT)
	keymap[emptyAction] = TwoKeyCombination(fyne.KeyQ, fyne.KeyY)
	inputHandler := NewHandler(keymap)
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
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyQ)
	keymap[testAction2] = TwoKeyCombination(fyne.KeyJ, fyne.KeyQ)
	inputHandler := NewHandler(keymap)
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
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyC)
	keymap[testAction2] = SingleKeyCombination(fyne.KeyC)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("first caller", testAction, testActionFunc)
	inputHandler.BindFunctionToAction("second caller", testAction2, testAction2Func)
	inputHandler.Handle("first caller", fyne.KeyC)
	assert.False(t, testAction2Executed)
	assert.True(t, testActionExecuted)
	inputHandler.Handle("second caller", fyne.KeyC)
	assert.True(t, testAction2Executed)
}

func TestThatOnKeyPressedCallbackFunctionIsCalledWithProperArgumentsOnKeyPress(t *testing.T) {
	pressedKeyCombination := KeyCombination{}
	keymap := make(map[Action]KeyCombination)
	handler := NewHandler(keymap)
	handler.SetOnKeyPressedCallbackFunction(func(keyCombination KeyCombination) { pressedKeyCombination = keyCombination })
	handler.Handle("", fyne.KeyR)
	assert.Equal(t, fyne.KeyR, pressedKeyCombination.firstKey)
	assert.Equal(t, fyne.KeyName(""), pressedKeyCombination.secondKey)
	handler.Handle("", fyne.KeyY)
	assert.Equal(t, fyne.KeyR, pressedKeyCombination.firstKey)
	assert.Equal(t, fyne.KeyY, pressedKeyCombination.secondKey)
}
