package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatFormDisplays(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "form dialog", "first", "second")
	dialog.Display()
	assert.True(t, dialog.Visible())
	assert.Equal(t, "form dialog", dialog.Title())
	assert.NotNil(t, dialog.inputs["first"])
	assert.NotNil(t, dialog.inputs["second"])
}

func TestThatFirstInputIsTheCurrentWhenDialogDisplays(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first", "second")
	dialog.Display()
	assert.Equal(t, dialog.currentInput(), dialog.inputs["first"])
}

func TestThatPressingJAndKSwitchesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first", "second", "third")
	dialog.Display()
	assert.Equal(t, dialog.currentInput(), dialog.inputs["first"])
	SimulateKeyPress(dialog, fyne.KeyJ)
	assert.Equal(t, dialog.currentInput(), dialog.inputs["second"])
	SimulateKeyPress(dialog, fyne.KeyJ)
	assert.Equal(t, dialog.currentInput(), dialog.inputs["third"])
	SimulateKeyPress(dialog, fyne.KeyK)
	assert.Equal(t, dialog.currentInput(), dialog.inputs["second"])
	SimulateKeyPress(dialog, fyne.KeyK)
	assert.Equal(t, dialog.currentInput(), dialog.inputs["first"])
}

func TestThatPressingIFocusesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first", "second", "third")
	dialog.Display()
	assert.False(t, dialog.currentInput().Focused())
	SimulateKeyPress(dialog, fyne.KeyI)
	assert.True(t, dialog.currentInput().Focused())
}

func TestThatPressingEscapeUnfocusesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first", "second", "third")
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyI)
	SimulateKeyPress(dialog.currentInput(), fyne.KeyEscape)
	assert.False(t, dialog.currentInput().Focused())
}

func TestThatPressingEnterCallsFunctionAndHidesAndUnfocusesDialog(t *testing.T) {
	functionCalled := false
	dialog := NewFormDialog(test.Canvas(), "", "first", "second", "third")
	dialog.OnEnterPressed = func() { functionCalled = true }
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyReturn)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
	assert.True(t, functionCalled)
	functionCalled = false
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyEnter)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
	assert.True(t, functionCalled)
	functionCalled = false
	dialog.Display()
	SimulateKeyPress(dialog.currentInput(), fyne.KeyReturn)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
	assert.True(t, functionCalled)
	functionCalled = false
	dialog.Display()
	SimulateKeyPress(dialog.currentInput(), fyne.KeyEnter)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
}

func TestThatFirstInputIsSelectedOnDialogReopening(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first", "second", "third")
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyJ)
	SimulateKeyPress(dialog, fyne.KeyEnter)
	dialog.Display()
	assert.Equal(t, dialog.inputs["first"], dialog.currentInput())
}

func TestThatPressingEscapeWhenNotInEditionModeClosesDialog(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first")
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyEscape)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
}

func TestThatSettingAndGettingItemValueWorksCorrectly(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first")
	setValue := "test value"
	dialog.SetItemValue("first", setValue)
	receivedValue, err := dialog.GetItemValue("first")
	assert.Nil(t, err)
	assert.Equal(t, setValue, receivedValue)
}

func TestThatSettingItemValueOnNonExistingItemDoesNotPanic(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first")
	dialog.SetItemValue("non existing", "value")
}

func TestThatGettingItemValueOnNonExistingItemReturnsError(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first")
	value, err := dialog.GetItemValue("non existing")
	assert.Equal(t, "", value)
	assert.NotNil(t, err)
}

func TestThatCleaningItemValuesWorks(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), "", "first", "second")
	dialog.SetItemValue("first", "val1")
	dialog.SetItemValue("second", "val2")
	dialog.CleanItemValues()
	assert.Empty(t, dialog.inputs["first"].Text)
	assert.Empty(t, dialog.inputs["second"].Text)
}