package widget

import (
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatItIsOnlyPossibleToInputNumbers(t *testing.T) {
	inputField := NewNumericInputField(test.Canvas(), getInputHandlerForTesting())
	inputField.canvas.Focus(inputField)
	TypeIntoFocusable(inputField, "nw1eut2hqr3kjlkj4lb`;,5;[!67\\rqr8zzzz9yyyy0")
	assert.Equal(t, "1234567890", inputField.Text)
}
