package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"testing"
	"wirwl/internal/input"
)

func TestThatFunctionGetsCalledOnConfirm(t *testing.T) {
	functionExecuted := false
	inputField := NewInputField(test.Canvas(), getInputHandlerForTesting())
	inputField.SetOnConfirm(func() { functionExecuted = true })
	inputField.canvas.Focus(inputField)
	SimulateKeyPress(inputField, fyne.KeyReturn)
	assert.Equal(t, true, functionExecuted)
}

func TestThatFunctionGetsCalledOnCancel(t *testing.T) {
	functionExecuted := false
	inputField := NewInputField(test.Canvas(), getInputHandlerForTesting())
	inputField.SetOnExitInputModeFunction(func() { functionExecuted = true })
	inputField.canvas.Focus(inputField)
	SimulateKeyPress(inputField, fyne.KeyEscape)
	assert.True(t, functionExecuted)
}

func TestThatTypingWorks(t *testing.T) {
	inputField := NewInputField(test.Canvas(), getInputHandlerForTesting())
	inputField.FocusLost()
	inputField.FocusGained()
	inputField.Type("some value")
	assert.Equal(t, "some value", inputField.Text)
}

func TestThatFunctionsAreNotNil(t *testing.T) {
	inputField := NewInputField(test.Canvas(), getInputHandlerForTesting())
	assert.NotNil(t, inputField.OnConfirm)
}

func TestEnteringIntoInputMode(t *testing.T) {
	inputField := NewInputField(test.Canvas(), getInputHandlerForTesting())
	inputField.EnterInputMode()
	assert.Equal(t, inputField.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	assert.True(t, inputField.Focused())
	assert.Equal(t, inputField, inputField.canvas.Focused())
}

func TestHighlightingAndUnhiglighting(t *testing.T) {
	inputField := NewInputField(test.Canvas(), getInputHandlerForTesting())
	//InputField needs to be placed into a test window, otherwise renderer doesn't work properly and marking sets background color again
	test.NewApp().NewWindow("").SetContent(inputField)
	assert.Equal(t, inputField.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	inputField.Highlight()
	assert.Equal(t, inputField.bgRenderer.BackgroundColor(), theme.FocusColor())
	inputField.Unhighlight()
	assert.Equal(t, inputField.bgRenderer.BackgroundColor(), theme.BackgroundColor())
}

func TestThatWhenExitingInputModeWithTwoKeyCombinationNeitherKeyOfCombinationGetsLeftInFieldsText(t *testing.T) {
	keymap := make(map[input.Action]input.KeyCombination)
	keymap[input.ExitInputModeAction] = input.TwoKeyCombination(fyne.KeyJ, fyne.KeyJ)
	inputHandler := input.NewHandler(keymap)
	inputField := NewInputField(test.Canvas(), inputHandler)
	inputField.canvas.Focus(inputField)
	inputField.Type("abc")
	SimulateKeyPress(inputField, fyne.KeyJ)
	SimulateKeyPress(inputField, fyne.KeyJ)
	assert.Equal(t, "abc", inputField.Text)
}
