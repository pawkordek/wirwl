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
	input := NewInputField(test.Canvas(), getInputHandlerForTesting())
	input.SetOnConfirm(func() { functionExecuted = true })
	input.canvas.Focus(input)
	SimulateKeyPress(input, fyne.KeyReturn)
	assert.Equal(t, true, functionExecuted)
}

func TestThatFunctionGetsCalledOnCancel(t *testing.T) {
	functionExecuted := false
	input := NewInputField(test.Canvas(), getInputHandlerForTesting())
	input.SetOnExitInputModeFunction(func() { functionExecuted = true })
	input.canvas.Focus(input)
	SimulateKeyPress(input, fyne.KeyEscape)
	assert.True(t, functionExecuted)
}

func TestThatTypingWorks(t *testing.T) {
	input := NewInputField(test.Canvas(), getInputHandlerForTesting())
	input.FocusLost()
	input.FocusGained()
	input.Type("some value")
	assert.Equal(t, "some value", input.Text)
}

func TestThatFunctionsAreNotNil(t *testing.T) {
	input := NewInputField(test.Canvas(), getInputHandlerForTesting())
	assert.NotNil(t, input.OnConfirm)
}

func TestEnteringIntoInputMode(t *testing.T) {
	input := NewInputField(test.Canvas(), getInputHandlerForTesting())
	input.EnterInputMode()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	assert.True(t, input.Focused())
	assert.Equal(t, input, input.canvas.Focused())
}

func TestHighlightingAndUnhiglighting(t *testing.T) {
	input := NewInputField(test.Canvas(), getInputHandlerForTesting())
	//InputField needs to be placed into a test window, otherwise renderer doesn't work properly and marking sets background color again
	test.NewApp().NewWindow("").SetContent(input)
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	input.Highlight()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.FocusColor())
	input.Unhighlight()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
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
