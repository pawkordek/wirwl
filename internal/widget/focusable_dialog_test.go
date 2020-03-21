package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatNewFocusableDialogIsHidden(t *testing.T) {
	dialog := NewFocusableDialog(test.Canvas())
	assert.Equal(t, true, dialog.Hidden)
}

func TestThatDisplayShowsFocusedFocusableDialogWithSpecifiedData(t *testing.T) {
	label := widget.NewLabel("some title")
	dialog := NewFocusableDialog(test.Canvas(), label)
	dialog.Display("some title")
	assert.True(t, dialog.Visible())
	assert.True(t, dialog.Focused())
	assert.Equal(t, "some title", dialog.Title())
	assert.True(t, ContainsWidget(dialog.Content, label))
}

func TestThatPressingAnyKeyHidesFocusableDialog(t *testing.T) {
	dialog := NewFocusableDialog(test.Canvas())
	dialog.Display("")
	SimulateKeyPress(dialog, fyne.KeyQ)
	assert.Equal(t, true, dialog.Hidden)
	assert.Equal(t, false, dialog.Focused())
	dialog = NewFocusableDialog(test.Canvas())
	dialog.Display("")
	SimulateKeyPress(dialog, fyne.Key1)
	assert.Equal(t, true, dialog.Hidden)
	assert.Equal(t, false, dialog.Focused())
}
