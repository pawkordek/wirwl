package widget

import (
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatTypeGetsSet(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	assert.Equal(t, popUp.title.Text, "")
	popUp.SetType(InfoPopUp)
	assert.Equal(t, popUp.title.Text, "INFO")
	popUp.SetType(SuccessPopUp)
	assert.Equal(t, popUp.title.Text, "SUCCESS")
	popUp.SetType(WarningPopUp)
	assert.Equal(t, popUp.title.Text, "WARNING")
	popUp.SetType(ErrorPopUp)
	assert.Equal(t, popUp.title.Text, "ERROR")
	popUp.SetType("non existing type")
	assert.Equal(t, popUp.title.Text, "")
}

func TestThatDisplayShowsPopUpWithSpecifiedData(t *testing.T) {
	popUp := NewMsgPopUp(test.Canvas())
	popUp.Display(SuccessPopUp, "some message")
	assert.True(t, popUp.Visible())
	assert.Equal(t, "SUCCESS", popUp.title.Text)
	assert.Equal(t, "some message", popUp.msg.Text)
}
