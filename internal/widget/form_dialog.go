package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"wirwl/internal/input"
)

/*
Implementation of a dialog that contains a title and a form below it.
Amount of form items depends on the argument passed when creating the dialog.
All keys described here assume default keymap configuration.
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
	inputHandler    input.InputHandler
}

func NewFormDialog(canvas fyne.Canvas, inputHandler input.InputHandler, title string, items ...string) *FormDialog {
	inputs := make(map[string]*Input)
	form := widget.NewForm()
	for _, item := range items {
		input := NewInput(inputHandler)
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
		input.SetOnConfirm(dialog.handleEnterKey)
	}
	dialog.title.SetText(title)
	dialog.ExtendBaseWidget(dialog)
	dialog.Hide()
	dialog.inputHandler = inputHandler
	dialog.setupInputHandler()
	return dialog
}

func (dialog *FormDialog) setupInputHandler() {
	dialog.inputHandler.BindFunctionToAction(dialog, input.MoveDownAction, func() {
		dialog.setCurrentInputTo(dialog.currentInputNum + 1)
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.MoveUpAction, func() {
		dialog.setCurrentInputTo(dialog.currentInputNum - 1)
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.EnterInputModeAction, func() {
		dialog.currentInput().Unmark()
		dialog.Canvas.Focus(dialog.currentInput())
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.ConfirmAction, func() {
		dialog.handleEnterKey()
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.CancelAction, func() {
		dialog.Canvas.Unfocus()
		dialog.Hide()
	})
}

func (dialog *FormDialog) TypedKeyInInput(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		dialog.setCurrentInputTo(dialog.currentInputNum)
		dialog.Canvas.Focus(dialog)
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
	dialog.inputHandler.Handle(dialog, key.Name)
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

func (dialog *FormDialog) ItemValue(itemName string) string {
	if dialog.inputs[itemName] != nil {
		return dialog.inputs[itemName].Text
	}
	return ""
}

func (dialog *FormDialog) CleanItemValues() {
	for _, input := range dialog.inputs {
		input.SetText("")
	}
}
