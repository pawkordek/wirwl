package wirwl

import (
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatTypeGetsSet(t *testing.T) {
	popUp := newMsgPopUp(test.Canvas())
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

func TestThatMsgGetsSet(t *testing.T) {
	popUp := newMsgPopUp(test.Canvas())
	assert.Equal(t, "", popUp.msg.Text)
	popUp.SetMsg("Some msg")
	assert.Equal(t, "Some msg", popUp.msg.Text)
}
