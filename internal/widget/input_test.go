package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatFunctionGetsCalledOnConfirm(t *testing.T) {
	functionExecuted := false
	input := NewInput(test.Canvas(), getInputHandlerForTesting())
	input.SetOnConfirm(func() { functionExecuted = true })
	SimulateKeyPress(input, fyne.KeyReturn)
	assert.Equal(t, true, functionExecuted)
}

func TestThatFunctionGetsCalledOnCancel(t *testing.T) {
	functionExecuted := false
	input := NewInput(test.Canvas(), getInputHandlerForTesting())
	input.SetOnExitInputModeFunction(func() { functionExecuted = true })
	SimulateKeyPress(input, fyne.KeyEscape)
	assert.True(t, functionExecuted)
}

func TestThatTypingWorks(t *testing.T) {
	input := NewInput(test.Canvas(), getInputHandlerForTesting())
	input.FocusLost()
	input.FocusGained()
	input.Type("some value")
	assert.Equal(t, "some value", input.Text)
}

func TestThatFunctionsAreNotNil(t *testing.T) {
	input := NewInput(test.Canvas(), getInputHandlerForTesting())
	assert.NotNil(t, input.OnConfirm)
}

func TestEnteringIntoInputMode(t *testing.T) {
	input := NewInput(test.Canvas(), getInputHandlerForTesting())
	input.EnterInputMode()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	assert.True(t, input.Focused())
	assert.Equal(t, input, input.canvas.Focused())
}

func TestMarkingAndUnmarking(t *testing.T) {
	input := NewInput(test.Canvas(), getInputHandlerForTesting())
	//Input needs to be placed into a test window, otherwise renderer doesn't work properly and marking sets background color again
	test.NewApp().NewWindow("").SetContent(input)
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	input.Mark()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.FocusColor())
	input.Unmark()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
}
