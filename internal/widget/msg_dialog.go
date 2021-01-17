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
	dialog := newMsgDialog(canvas)
	dialog.ExtendBaseWidget(dialog)
	return dialog
}

//Should be used for dialog creation by any widget that embed this widget so it can properly extend fyne's BaseWidget
func newMsgDialog(canvas fyne.Canvas) *MsgDialog {
	msg := widget.NewLabel("")
	msg.Alignment = fyne.TextAlignCenter
	dialog := &MsgDialog{
		FocusableDialog: newFocusableDialog(canvas, msg),
		msg:             msg,
	}
	dialog.Hide()
	return dialog
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
	popUp.Show()
	popUp.Canvas.Focus(popUp)
}
