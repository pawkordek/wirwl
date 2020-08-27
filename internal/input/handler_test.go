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
	handler.HandleInNormalMode("", fyne.KeyU)
}

func TestThatInputHandlerHandlesSingleKeyActions(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyQ)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInNormalMode("", fyne.KeyQ)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionsInKeySequences(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyZ, fyne.KeyX)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInNormalMode("", fyne.KeyZ)
	assert.False(t, functionExecuted)
	inputHandler.HandleInNormalMode("", fyne.KeyX)
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
	inputHandler.HandleInNormalMode("firstCaller", fyne.KeyT)
	assert.True(t, firstCallerFunctionExecuted)
	assert.False(t, secondCallerFunctionExecuted)
	inputHandler.HandleInNormalMode("secondCaller", fyne.KeyT)
	assert.True(t, secondCallerFunctionExecuted)
}

func TestThatTryingToHandleInputForActionWithoutBindFunctionDoesNotCauseErrors(t *testing.T) {
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyL)
	inputHandler := NewHandler(keymap)
	inputHandler.HandleInNormalMode("", fyne.KeyL)
}

func TestThatKeySequenceDoesNotWorkIfCertainTimePasses(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyU, fyne.KeyF)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInNormalMode("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.HandleInNormalMode("", fyne.KeyF)
	assert.False(t, functionExecuted)
}

func TestThatKeySequenceWorksAfterPressingKeyPausingAndPressingNextKey(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyU, fyne.KeyF)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInNormalMode("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.HandleInNormalMode("", fyne.KeyF)
	inputHandler.HandleInNormalMode("", fyne.KeyU)
	inputHandler.HandleInNormalMode("", fyne.KeyF)
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
	inputHandler.HandleInNormalMode("", fyne.KeyQ)
	inputHandler.HandleInNormalMode("", fyne.KeyY)
	time.Sleep(2 * time.Second)
	inputHandler.HandleInNormalMode("", fyne.KeyT)
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
	inputHandler.HandleInNormalMode("", fyne.KeyJ)
	inputHandler.HandleInNormalMode("", fyne.KeyQ)
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
	inputHandler.HandleInNormalMode("first caller", fyne.KeyC)
	assert.False(t, testAction2Executed)
	assert.True(t, testActionExecuted)
	inputHandler.HandleInNormalMode("second caller", fyne.KeyC)
	assert.True(t, testAction2Executed)
}

func TestThatOnKeyPressedCallbackFunctionIsCalledWithProperArgumentsOnKeyPress(t *testing.T) {
	pressedKeyCombination := KeyCombination{}
	keymap := make(map[Action]KeyCombination)
	handler := NewHandler(keymap)
	handler.SetOnKeyPressedCallbackFunction(func(keyCombination KeyCombination) { pressedKeyCombination = keyCombination })
	handler.HandleInNormalMode("", fyne.KeyR)
	assert.Equal(t, fyne.KeyR, pressedKeyCombination.firstKey)
	assert.Equal(t, fyne.KeyName(""), pressedKeyCombination.secondKey)
	handler.HandleInNormalMode("", fyne.KeyY)
	assert.Equal(t, fyne.KeyR, pressedKeyCombination.firstKey)
	assert.Equal(t, fyne.KeyY, pressedKeyCombination.secondKey)
}

func TestThatInputHandlerHandlesSingleKeyActionsInInputMode(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyQ)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInInputMode("", fyne.KeyQ)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionsInKeySequencesInInputMode(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyZ, fyne.KeyX)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInInputMode("", fyne.KeyZ)
	assert.False(t, functionExecuted)
	inputHandler.HandleInInputMode("", fyne.KeyX)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerHandlesActionForTheCorrectCallerInInputMode(t *testing.T) {
	firstCallerFunctionExecuted := false
	secondCallerFunctionExecuted := false
	firstCallerFunction := func() { firstCallerFunctionExecuted = true }
	secondCallerFunction := func() { secondCallerFunctionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyT)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("firstCaller", testAction, firstCallerFunction)
	inputHandler.BindFunctionToAction("secondCaller", testAction, secondCallerFunction)
	inputHandler.HandleInInputMode("firstCaller", fyne.KeyT)
	assert.True(t, firstCallerFunctionExecuted)
	assert.False(t, secondCallerFunctionExecuted)
	inputHandler.HandleInInputMode("secondCaller", fyne.KeyT)
	assert.True(t, secondCallerFunctionExecuted)
}

func TestThatTryingToHandleInputForActionWithoutBindFunctionDoesNotCauseErrorsInInputMode(t *testing.T) {
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyL)
	inputHandler := NewHandler(keymap)
	inputHandler.HandleInInputMode("", fyne.KeyL)
}

func TestThatKeySequenceDoesNotWorkIfCertainTimePassesInInputMode(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyU, fyne.KeyF)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInInputMode("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.HandleInInputMode("", fyne.KeyF)
	assert.False(t, functionExecuted)
}

func TestThatKeySequenceWorksAfterPressingKeyPausingAndPressingNextKeyInInputMode(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = TwoKeyCombination(fyne.KeyU, fyne.KeyF)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.HandleInInputMode("", fyne.KeyU)
	time.Sleep(1 * time.Second)
	inputHandler.HandleInInputMode("", fyne.KeyF)
	inputHandler.HandleInInputMode("", fyne.KeyU)
	inputHandler.HandleInInputMode("", fyne.KeyF)
	assert.True(t, functionExecuted)
}

func TestThatInputHandlerExecutesSingleKeyActionIfKeyCombinationWasPressedAndTimePassedInInputMode(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyT)
	keymap[emptyAction] = TwoKeyCombination(fyne.KeyQ, fyne.KeyY)
	inputHandler := NewHandler(keymap)
	inputHandler.BindFunctionToAction("", testAction, function)
	inputHandler.BindFunctionToAction("", emptyAction, func() {})
	inputHandler.HandleInInputMode("", fyne.KeyQ)
	inputHandler.HandleInInputMode("", fyne.KeyY)
	time.Sleep(2 * time.Second)
	inputHandler.HandleInInputMode("", fyne.KeyT)
	assert.True(t, functionExecuted)
}

func TestThatCombinationActionGetsExecutedInsteadOfSingleKeyActionIfBothCouldHappenInInputMode(t *testing.T) {
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
	inputHandler.HandleInInputMode("", fyne.KeyJ)
	inputHandler.HandleInInputMode("", fyne.KeyQ)
	assert.True(t, combinationActionExecuted)
	assert.False(t, singleKeyActionExecuted)
}

func TestThatInputHandlerAllowsToUseTheSameKeyForTwoVariousActionsIfTheyAreForDifferentCallerInInputMode(t *testing.T) {
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
	inputHandler.HandleInInputMode("first caller", fyne.KeyC)
	assert.False(t, testAction2Executed)
	assert.True(t, testActionExecuted)
	inputHandler.HandleInInputMode("second caller", fyne.KeyC)
	assert.True(t, testAction2Executed)
}

func TestThatOnKeyPressedCallbackFunctionIsCalledWithProperArgumentsOnKeyPressInInputMode(t *testing.T) {
	pressedKeyCombination := KeyCombination{}
	keymap := make(map[Action]KeyCombination)
	handler := NewHandler(keymap)
	handler.SetOnKeyPressedCallbackFunction(func(keyCombination KeyCombination) { pressedKeyCombination = keyCombination })
	handler.HandleInInputMode("", fyne.KeyR)
	assert.Equal(t, fyne.KeyR, pressedKeyCombination.firstKey)
	assert.Equal(t, fyne.KeyName(""), pressedKeyCombination.secondKey)
	handler.HandleInInputMode("", fyne.KeyY)
	assert.Equal(t, fyne.KeyR, pressedKeyCombination.firstKey)
	assert.Equal(t, fyne.KeyY, pressedKeyCombination.secondKey)
}

func TestThatLastKeyOfCombinationBecomesFirstKeyAfterNextPressWhenHandlingKeyPressesInInputMode(t *testing.T) {
	keymap := make(map[Action]KeyCombination)
	handler := NewHandler(keymap)
	handler.HandleInInputMode("", fyne.KeyC)
	handler.HandleInInputMode("", fyne.KeyB)
	assert.Equal(t, TwoKeyCombination(fyne.KeyC, fyne.KeyB), handler.currentKeyCombination)
	handler.HandleInInputMode("", fyne.KeyA)
	assert.Equal(t, TwoKeyCombination(fyne.KeyB, fyne.KeyA), handler.currentKeyCombination)
}

func TestThatSingleKeyActionWorksIfThePressedKeyMatchesSecondKeyOfCurrentCombinationInInputMode(t *testing.T) {
	testActionExecuted := false
	testActionFunc := func() { testActionExecuted = true }
	keymap := make(map[Action]KeyCombination)
	keymap[testAction] = SingleKeyCombination(fyne.KeyEscape)
	handler := NewHandler(keymap)
	handler.BindFunctionToAction("", testAction, testActionFunc)
	handler.HandleInInputMode("", fyne.KeyF)
	handler.HandleInInputMode("", fyne.KeyEscape)
	assert.True(t, testActionExecuted)
}
