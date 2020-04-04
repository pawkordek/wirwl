package widget

import (
	"errors"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

/*
Implementation of a dialog that contains a title and a form below it.
Amount of form items depends on the argument passed when creating the dialog.
A user can switch between form items using J and K keys.
Currently selected item will be highlighted by theme's focused color.
Pressing I key starts the edition of the currently selected form item.
Pressing escape key when editing exits edition.
Pressing escape key when not editing closes the dialog.
Pressing enter key anytime closes the dialog and calls the function specified for this action.

Values of every form item can be set, retrieved and cleaned using proper functions.
*/
type FormDialog struct {
	FocusableDialog
	items           []string
	currentInputNum int
	inputs          map[string]*Input
	OnEnterPressed  func()
}

func NewFormDialog(canvas fyne.Canvas, title string, items ...string) *FormDialog {
	inputs := make(map[string]*Input)
	form := widget.NewForm()
	for _, item := range items {
		input := NewInput()
		inputs[item] = input
		formItem := widget.NewFormItem(item, input)
		form.AppendItem(formItem)
	}
	dialog := &FormDialog{
		FocusableDialog: *NewFocusableDialog(canvas, form),
		items:           items,
		currentInputNum: 0,
		inputs:          inputs,
		OnEnterPressed:  func() {},
	}
	for _, input := range dialog.inputs {
		input.SetOnTypedKey(dialog.TypedKeyInInput)
	}
	dialog.title.SetText(title)
	dialog.ExtendBaseWidget(dialog)
	dialog.Hide()
	return dialog
}

func (dialog *FormDialog) TypedKeyInInput(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		dialog.setCurrentInputTo(dialog.currentInputNum)
		dialog.Canvas.Focus(dialog)
	} else if key.Name == fyne.KeyEnter || key.Name == fyne.KeyReturn {
		dialog.handleEnterKey()
	}
}

func (dialog *FormDialog) handleEnterKey() {
	dialog.Canvas.Unfocus()
	dialog.Hide()
	dialog.OnEnterPressed()
}

func (dialog *FormDialog) Display() {
	dialog.currentInputNum = 0
	dialog.currentInput().Mark()
	dialog.Canvas.Focus(dialog)
	dialog.Show()
}

func (dialog *FormDialog) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyJ {
		dialog.setCurrentInputTo(dialog.currentInputNum + 1)
	} else if key.Name == fyne.KeyK {
		dialog.setCurrentInputTo(dialog.currentInputNum - 1)
	} else if key.Name == fyne.KeyI {
		dialog.currentInput().Unmark()
		dialog.Canvas.Focus(dialog.currentInput())
	} else if key.Name == fyne.KeyEnter || key.Name == fyne.KeyReturn {
		dialog.handleEnterKey()
	} else if key.Name == fyne.KeyEscape {
		dialog.Canvas.Unfocus()
		dialog.Hide()
	}
}

func (dialog *FormDialog) setCurrentInputTo(number int) {
	if number < len(dialog.inputs) && number >= 0 {
		dialog.currentInput().Unmark()
		dialog.currentInputNum = number
		dialog.currentInput().Mark()
	}
}

func (dialog *FormDialog) currentInput() *Input {
	return dialog.inputs[dialog.items[dialog.currentInputNum]]
}

func (dialog *FormDialog) SetItemValue(itemName string, value string) {
	if dialog.inputs[itemName] != nil {
		dialog.inputs[itemName].SetText(value)
	}
}

func (dialog *FormDialog) GetItemValue(itemName string) (string, error) {
	if dialog.inputs[itemName] != nil {
		return dialog.inputs[itemName].Text, nil
	}
	return "", errors.New("There is no item with name=" + itemName)
}

func (dialog *FormDialog) CleanItemValues() {
	for _, input := range dialog.inputs {
		input.SetText("")
	}
}
