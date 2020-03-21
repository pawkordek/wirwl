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
	assert.NotNil(t, input.OnEnterPressed)
	assert.NotNil(t, input.OnTypedKey)
}
