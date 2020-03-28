package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatNewPopUpIsHidden(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	assert.Equal(t, true, popUp.Hidden)
}

func TestThatTypeGetsSet(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	assert.Equal(t, popUp.title.Text, "")
	popUp.setType(InfoPopUp)
	assert.Equal(t, popUp.title.Text, "INFO")
	popUp.setType(SuccessPopUp)
	assert.Equal(t, popUp.title.Text, "SUCCESS")
	popUp.setType(WarningPopUp)
	assert.Equal(t, popUp.title.Text, "WARNING")
	popUp.setType(ErrorPopUp)
	assert.Equal(t, popUp.title.Text, "ERROR")
	popUp.setType("non existing type")
	assert.Equal(t, popUp.title.Text, "")
}

func TestThatDisplayShowsPopUpWithSpecifiedData(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	popUp.Display(SuccessPopUp, "some message")
	assert.True(t, popUp.Visible())
	assert.Equal(t, "SUCCESS", popUp.title.Text)
	assert.Equal(t, "some message", popUp.msg.Text)
}

func TestThatCallingDisplayViewsAndFocusesPopUp(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	popUp.Display(InfoPopUp, "testing")
	assert.Equal(t, true, popUp.Visible())
	assert.Equal(t, true, popUp.Focused())
}

func TestThatPressingAnyKeyHidesThePopUp(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	popUp.Display(InfoPopUp, "testing")
	SimulateKeyPress(popUp, fyne.KeyQ)
	assert.Equal(t, true, popUp.Hidden)
	assert.Equal(t, false, popUp.Focused())
	popUp.Display(InfoPopUp, "testing")
	SimulateKeyPress(popUp, fyne.Key1)
	assert.Equal(t, true, popUp.Hidden)
	assert.Equal(t, false, popUp.Focused())
}
