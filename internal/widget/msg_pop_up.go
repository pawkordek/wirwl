package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type MsgPopUp struct {
	*widget.PopUp
	title *widget.Label
	msg   *widget.Label
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
	popUp := &MsgPopUp{widget.NewModalPopUp(content, canvas), title, msg}
	popUp.ExtendBaseWidget(popUp)
	return popUp
}

func (popUp *MsgPopUp) Title() string {
	return popUp.title.Text
}

func (popUp *MsgPopUp) Msg() string {
	return popUp.msg.Text
}

func (popUp *MsgPopUp) SetType(t string) {
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
	popUp.SetType(popUpType)
	popUp.msg.Text = msg
	popUp.Show()
}
