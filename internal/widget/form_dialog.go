package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"wirwl/internal/input"
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
	*FocusableDialog
	currentInputNum int
	OnEnterPressed  func()
	embeddedWidgets map[string]FormDialogEmbeddableWidget
	form            widget.Form
	inputHandler    input.Handler
}

type FormDialogFormItem struct {
	widget.FormItem
}

func newFormDialogFormItem(labelText string, embeddableWidget FormDialogEmbeddableWidget) *FormDialogFormItem {
	return &FormDialogFormItem{
		FormItem: *widget.NewFormItem(labelText, embeddableWidget),
	}
}

type FormDialogFormItemFactory struct {
	canvas       fyne.Canvas
	inputHandler input.Handler
}

func NewFormDialogFormItemFactory(canvas fyne.Canvas, inputHandler input.Handler) *FormDialogFormItemFactory {
	return &FormDialogFormItemFactory{
		canvas:       canvas,
		inputHandler: inputHandler,
	}
}

func (factory *FormDialogFormItemFactory) FormItemWithInputField(labelText string) *FormDialogFormItem {
	return newFormDialogFormItem(labelText, NewInputField(factory.canvas, factory.inputHandler))
}

func (factory *FormDialogFormItemFactory) FormItemWithSelect(labelText string, selectChoices ...string) *FormDialogFormItem {
	return newFormDialogFormItem(labelText, NewSelect(factory.canvas, factory.inputHandler, selectChoices...))
}

//Form items will be displayed in the order they are passed in
func NewFormDialog(canvas fyne.Canvas, inputHandler input.Handler, title string, formItems ...*FormDialogFormItem) *FormDialog {
	embeddedWidgets := make(map[string]FormDialogEmbeddableWidget, len(formItems))
	form := widget.NewForm()
	for _, formItem := range formItems {
		form.AppendItem(&formItem.FormItem)
		embeddedWidgets[formItem.Text] = formItem.Widget.(FormDialogEmbeddableWidget)
	}
	dialog := &FormDialog{
		FocusableDialog: newFocusableDialog(canvas, form),
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
	dialog.Show()
	dialog.Canvas.Focus(dialog)
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
