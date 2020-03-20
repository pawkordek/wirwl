package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatConfirmingWorksCorrectly(t *testing.T) {
	called := false
	dialog := NewConfirmationDialog(test.Canvas())
	dialog.OnConfirm = func() { called = true }
	dialog.Display("")
	SimulateKeyPress(dialog, fyne.KeyY)
	assert.Equal(t, true, called)
	assert.Equal(t, true, dialog.Hidden)
	assert.Equal(t, false, dialog.Focused())
}

func TestThatCancellingWorksCorrectly(t *testing.T) {
	called := false
	dialog := NewConfirmationDialog(test.Canvas())
	dialog.OnCancel = func() { called = true }
	dialog.Display("")
	SimulateKeyPress(dialog, fyne.KeyN)
	assert.Equal(t, true, called)
	assert.Equal(t, true, dialog.Hidden)
	assert.Equal(t, false, dialog.Focused())
}

func TestThatDisplayedMessageHasYesOrNoMessageAppended(t *testing.T) {
	dialog := NewConfirmationDialog(test.Canvas())
	dialog.Display("Some message.")
	assert.Equal(t, "Some message. (y)es or (n)o?", dialog.Msg())
}

func TestThatItIsNotPossibleToExitDialogWithOtherButtonThanYOrN(t *testing.T) {
	dialog := NewConfirmationDialog(test.Canvas())
	dialog.Display("message")
	SimulateKeyPress(dialog, fyne.KeyE)
	SimulateKeyPress(dialog, fyne.KeyI)
	SimulateKeyPress(dialog, fyne.KeyA)
	SimulateKeyPress(dialog, fyne.Key1)
	assert.Equal(t, true, dialog.Visible())
	assert.Equal(t, true, dialog.Focused())
}
