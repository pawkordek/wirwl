package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/pkg/errors"
	"wirwl/internal/input"
	"wirwl/internal/log"
)

/*
Implementation of a dialog that contains a title and a form below it.
Amount of embedded widgets depends on the argument passed when creating the dialog.
All keys described here assume default keymap configuration.
A user can switch between embedded widgets using J and K keys.
Currently selected item will be highlighted by theme's focused color.
Pressing I key starts the edition of the currently selected form item.
Pressing escape key when editing exits edition.
Pressing escape key when not editing closes the dialog.
Pressing enter key anytime closes the dialog and calls the function specified for this action.

Values of every form item can be set, retrieved and cleaned using proper functions.
*/
type FormDialog struct {
	FocusableDialog
	currentInputNum int
	OnEnterPressed  func()
	embeddedWidgets map[string]FormDialogEmbeddableWidget
	form            widget.Form
	inputHandler    input.Handler
}

type FormDialogElement struct {
	LabelText  string
	WidgetType Type
}

//Passed in elements will be displayed in the order they are in the array.
func NewFormDialog(canvas fyne.Canvas, inputHandler input.Handler, title string, elements []FormDialogElement) *FormDialog {
	embeddedWidgets := make(map[string]FormDialogEmbeddableWidget, len(elements))
	form := widget.NewForm()
	for _, element := range elements {
		embeddedWidget := widgetFromType(element.WidgetType, canvas, inputHandler)
		form.Append(element.LabelText, embeddedWidget)
		embeddedWidgets[element.LabelText] = embeddedWidget
	}
	dialog := &FormDialog{
		FocusableDialog: *NewFocusableDialog(canvas, form),
		currentInputNum: 0,
		OnEnterPressed:  func() {},
		embeddedWidgets: embeddedWidgets,
		form:            *form,
	}
	for _, embeddedWidget := range dialog.embeddedWidgets {
		embeddedWidget.SetOnConfirm(dialog.handleEnterKey)
		embeddedWidget.SetOnExitInputModeFunction(func() {
			dialog.setCurrentInputTo(dialog.currentInputNum)
			dialog.Canvas.Focus(dialog)
		})
	}
	dialog.title.SetText(title)
	dialog.ExtendBaseWidget(dialog)
	dialog.Hide()
	dialog.inputHandler = inputHandler
	dialog.setupInputHandler()
	return dialog
}

func widgetFromType(widgetType Type, canvas fyne.Canvas, inputHandler input.Handler) FormDialogEmbeddableWidget {
	switch widgetType {
	case InputFieldType:
		return NewInputField(canvas, inputHandler)
	}
	err := errors.New("A widget type called '" + string(widgetType) + "'that cannot be instantiated has been passed to Form Dialog. This is a coding error and should be fixed!")
	log.Panic(err)
	return nil
}

func (dialog *FormDialog) setupInputHandler() {
	dialog.inputHandler.BindFunctionToAction(dialog, input.MoveDownAction, func() {
		dialog.setCurrentInputTo(dialog.currentInputNum + 1)
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.MoveUpAction, func() {
		dialog.setCurrentInputTo(dialog.currentInputNum - 1)
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.EnterInputModeAction, func() {
		dialog.currentWidget().EnterInputMode()
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.ConfirmAction, func() {
		dialog.handleEnterKey()
	})
	dialog.inputHandler.BindFunctionToAction(dialog, input.CancelAction, func() {
		dialog.Canvas.Unfocus()
		dialog.Hide()
	})
}

func (dialog *FormDialog) handleEnterKey() {
	dialog.Canvas.Unfocus()
	dialog.Hide()
	dialog.OnEnterPressed()
}

func (dialog *FormDialog) Display() {
	dialog.currentInputNum = 0
	dialog.currentWidget().Highlight()
	dialog.Canvas.Focus(dialog)
	dialog.Show()
}

func (dialog *FormDialog) TypedKey(key *fyne.KeyEvent) {
	dialog.inputHandler.HandleInNormalMode(dialog, key.Name)
}

func (dialog *FormDialog) setCurrentInputTo(number int) {
	if number < len(dialog.embeddedWidgets) && number >= 0 {
		dialog.currentWidget().Unhighlight()
		dialog.currentInputNum = number
		dialog.currentWidget().Highlight()
	}
}

func (dialog *FormDialog) currentWidget() FormDialogEmbeddableWidget {
	return dialog.form.Items[dialog.currentInputNum].Widget.(FormDialogEmbeddableWidget)
}

func (dialog *FormDialog) SetItemValue(itemName string, value string) {
	if dialog.embeddedWidgets[itemName] != nil {
		dialog.embeddedWidgets[itemName].SetText(value)
	}
}

func (dialog *FormDialog) ItemValue(itemName string) string {
	if dialog.embeddedWidgets[itemName] != nil {
		return dialog.embeddedWidgets[itemName].GetText()
	}
	return ""
}

func (dialog *FormDialog) CleanItemValues() {
	for _, inputField := range dialog.embeddedWidgets {
		inputField.SetText("")
	}
}
