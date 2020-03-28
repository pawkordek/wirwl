package widget

import "fyne.io/fyne"

/*
Implementation of a typical yes/no dialog except that confirmation and cancellation are done by using 'y' and 'n' keys
respectively, therefore no buttons are shown.
Dialog disappears only when any of those two buttons is pressed.
Message passed into has an instruction message '(y)es or (n)o' appended at the end so caller doesn't have to add it
every time.
*/
type ConfirmationDialog struct {
	*MsgDialog
	OnConfirm func()
	OnCancel  func()
	focused   bool
}

func NewConfirmationDialog(canvas fyne.Canvas) *ConfirmationDialog {
	dialog := ConfirmationDialog{}
	dialog.MsgDialog = NewMsgPopUp(canvas)
	dialog.OnConfirm = func() {}
	dialog.OnCancel = func() {}
	dialog.ExtendBaseWidget(dialog)
	return &dialog
}

func (dialog *ConfirmationDialog) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyY {
		dialog.OnConfirm()
		dialog.MsgDialog.TypedKey(key)
	} else if key.Name == fyne.KeyN {
		dialog.OnCancel()
		dialog.MsgDialog.TypedKey(key)
	}
}

func (dialog *ConfirmationDialog) Display(msg string) {
	msg += " (y)es or (n)o?"
	dialog.MsgDialog.Display(InfoPopUp, msg)
	dialog.Canvas.Focus(dialog)
}
