package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type MsgDialog struct {
	*FocusableDialog
	msg *widget.Label
}

const (
	InfoPopUp    = "info"
	SuccessPopUp = "success"
	WarningPopUp = "warning"
	ErrorPopUp   = "error"
)

func NewMsgPopUp(canvas fyne.Canvas) *MsgDialog {
	msg := widget.NewLabel("")
	msg.Alignment = fyne.TextAlignCenter
	popUp := &MsgDialog{
		FocusableDialog: newFocusableDialog(canvas, msg),
		msg:             msg,
	}
	popUp.ExtendBaseWidget(popUp)
	popUp.Hide()
	return popUp
}

func (popUp *MsgDialog) Title() string {
	return popUp.title.Text
}

func (popUp *MsgDialog) Msg() string {
	return popUp.msg.Text
}

func (popUp *MsgDialog) setType(t string) {
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

func (popUp *MsgDialog) Display(popUpType string, msg string) {
	popUp.setType(popUpType)
	popUp.msg.SetText(msg)
	popUp.Canvas.Focus(popUp)
	popUp.Show()
}
