package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

/*
Implementation of a dialog displaying a the center of the screen that can be focused and responds to key events.
It contains a title at the top and can contain any content below it.
When any key gets pressed it hides itself.
*/
type FocusableDialog struct {
	*widget.PopUp
	title   *widget.Label
	focused bool
}

func NewFocusableDialog(canvas fyne.Canvas, content ...fyne.CanvasObject) *FocusableDialog {
	title := widget.NewLabel("")
	title.Alignment = fyne.TextAlignCenter
	content = append([]fyne.CanvasObject{title}, content...)
	popupContent := widget.NewVBox(content...)
	dialog := &FocusableDialog{
		PopUp:   widget.NewModalPopUp(popupContent, canvas),
		title:   title,
		focused: false,
	}
	dialog.ExtendBaseWidget(dialog)
	dialog.Hide()
	return dialog
}

func (dialog *FocusableDialog) Title() string {
	return dialog.title.Text
}

func (dialog *FocusableDialog) Display(title string) {
	dialog.title.SetText(title)
	dialog.Canvas.Focus(dialog)
	dialog.Show()
}

func (dialog *FocusableDialog) TypedKey(key *fyne.KeyEvent) {
	dialog.Hide()
	dialog.Canvas.Unfocus()
}

func (dialog *FocusableDialog) FocusGained() {
	dialog.focused = true
}

func (dialog *FocusableDialog) FocusLost() {
	dialog.focused = false
}

func (dialog *FocusableDialog) Focused() bool {
	return dialog.focused
}

func (dialog *FocusableDialog) TypedRune(r rune) {
	//Do nothing as inputting text handling is not needed, only key presses
}
