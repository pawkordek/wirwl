package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatConfirmingWorksCorrectly(t *testing.T) {
	called := false
	dialog := NewConfirmationDialog(test.Canvas())
	dialog.OnConfirm = func() { called = true }
	dialog.Display("")
	SimulateKeyPressOnTestCanvas(fyne.KeyY)
	assert.Equal(t, true, called)
	assert.Equal(t, true, dialog.Hidden)
	assert.Equal(t, false, dialog.Focused())
}

func TestThatCancellingWorksCorrectly(t *testing.T) {
	called := false
	dialog := NewConfirmationDialog(test.Canvas())
	dialog.OnCancel = func() { called = true }
	dialog.Display("")
	SimulateKeyPressOnTestCanvas(fyne.KeyN)
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
	SimulateKeyPressOnTestCanvas(fyne.KeyE)
	SimulateKeyPressOnTestCanvas(fyne.KeyI)
	SimulateKeyPressOnTestCanvas(fyne.KeyA)
	SimulateKeyPressOnTestCanvas(fyne.Key1)
	assert.Equal(t, true, dialog.Visible())
	assert.Equal(t, true, dialog.Focused())
}
