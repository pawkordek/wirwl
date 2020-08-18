package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatFunctionGetsCalledOnEnterPressed(t *testing.T) {
	functionExecuted := false
	input := NewInput()
	input.SetOnConfirm(func() { functionExecuted = true })
	SimulateKeyPress(input, fyne.KeyEnter)
	assert.Equal(t, true, functionExecuted)
	functionExecuted = false
	SimulateKeyPress(input, fyne.KeyReturn)
	assert.Equal(t, true, functionExecuted)
}

func TestThatFunctionGetsCalledOnAnyKeyPress(t *testing.T) {
	functionExecuted := false
	input := NewInput()
	input.SetOnTypedKey(func(key *fyne.KeyEvent) {
		functionExecuted = true
	})
	SimulateKeyPress(input, fyne.KeyEnter)
	assert.Equal(t, true, functionExecuted)
	functionExecuted = false
	SimulateKeyPress(input, fyne.Key1)
	assert.Equal(t, true, functionExecuted)
}

func TestThatTypingWorks(t *testing.T) {
	input := NewInput()
	input.FocusLost()
	input.FocusGained()
	input.Type("some value")
	assert.Equal(t, "some value", input.Text)
}

func TestThatFunctionsAreNotNil(t *testing.T) {
	input := NewInput()
	assert.NotNil(t, input.OnConfirm)
	assert.NotNil(t, input.OnTypedKey)
}

func TestMarkingAndUnmarking(t *testing.T) {
	input := NewInput()
	//Input needs to be placed into a test window, otherwise renderer doesn't work properly and marking sets background color again
	test.NewApp().NewWindow("").SetContent(input)
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
	input.Mark()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.FocusColor())
	input.Unmark()
	assert.Equal(t, input.bgRenderer.BackgroundColor(), theme.BackgroundColor())
}
