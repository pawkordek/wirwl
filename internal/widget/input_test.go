package widget

import (
	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatFunctionGetsCalledOnEnterPressed(t *testing.T) {
	functionExecuted := false
	input := NewInput()
	input.SetOnEnterPressed(func() { functionExecuted = true })
	SimulateKeyPress(input, fyne.KeyEnter)
	assert.Equal(t, true, functionExecuted)
	functionExecuted = false
	SimulateKeyPress(input, fyne.KeyReturn)
	assert.Equal(t, true, functionExecuted)
}

func TestThatTypingWorks(t *testing.T) {
	input := NewInput()
	input.FocusLost()
	input.FocusGained()
	input.Type("some value")
	assert.Equal(t, "some value", input.Text)
}
