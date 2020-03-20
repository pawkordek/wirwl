package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type MsgPopUp struct {
	*widget.PopUp
	title   *widget.Label
	msg     *widget.Label
	focused bool
}

const (
	InfoPopUp    = "info"
	SuccessPopUp = "success"
	WarningPopUp = "warning"
	ErrorPopUp   = "error"
)

func NewMsgPopUp(canvas fyne.Canvas) *MsgPopUp {
	title := widget.NewLabel("")
	title.Alignment = fyne.TextAlignCenter
	msg := widget.NewLabel("")
	msg.Alignment = fyne.TextAlignCenter
	content := widget.NewVBox(title, msg)
	popUp := &MsgPopUp{widget.NewModalPopUp(content, canvas), title, msg, false}
	popUp.ExtendBaseWidget(popUp)
	return popUp
}

func (popUp *MsgPopUp) Title() string {
	return popUp.title.Text
}

func (popUp *MsgPopUp) Msg() string {
	return popUp.msg.Text
}

func (popUp *MsgPopUp) setType(t string) {
	switch t {
	case InfoPopUp:
		popUp.title.SetText("INFO")
	case SuccessPopUp:
		popUp.title.SetText("SUCCESS")
	case WarningPopUp:
		popUp.title.SetText("WARNING")
	case ErrorPopUp:
		popUp.title.SetText("ERROR")
	default:
		popUp.title.SetText("")
	}
}

func (popUp *MsgPopUp) Display(popUpType string, msg string) {
	popUp.setType(popUpType)
	popUp.msg.Text = msg
	popUp.Show()
}

func (popUp *MsgPopUp) TypedKey(key *fyne.KeyEvent) {
	popUp.Hide()
}

func (popUp *MsgPopUp) FocusGained() {
	popUp.focused = true
}

func (popUp *MsgPopUp) FocusLost() {
	popUp.focused = false
}

func (popUp *MsgPopUp) Focused() bool {
	return popUp.focused
}

func (popUp *MsgPopUp) TypedRune(r rune) {
	//Do nothing as inputting text handling is not needed, only key presses
}
