package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
	"wirwl/internal/input"
)

/* When opening input in the running application, the last typed in character is still handled which normally means it
gets typed into the input. There is code that prevents this but as this situation doesn't happen when running the
test code any string typed into the input needs an additional character at the beginning as if the bug happened.
*/
func TypeIntoFocusable(focusable fyne.Focusable, chars string) {
	fixedChars := " " + chars
	for _, char := range fixedChars {
		focusable.TypedRune(char)
	}
}

func (inputField *InputField) Type(chars string) {
	TypeIntoFocusable(inputField, chars)
}

func (dialog *FormDialog) Type(chars string) {
	TypeIntoFocusable(dialog.currentWidget(), chars)
}

//Please make sure that tested focusable is focused if used on a standalone focusable, otherwise this won't work
func SimulateKeyPress(focusable fyne.Focusable, key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	if focusable.Focused() {
		focusable.TypedKey(event)
	}
	if focusable.Focused() {
		//This is a very naive way of retrieving a character to type in and will work incorrectly for e.g. Escape key
		character := []rune(event.Name)[0]
		focusable.TypedRune(character)
	}
}

func ContainsWidget(content fyne.CanvasObject, searchedWidget interface{}) bool {
	for _, existingWidget := range content.(*widget.Box).Children {
		if existingWidget == searchedWidget {
			return true
		}
	}
	return false
}

func ContainsLabelWithSameText(content fyne.CanvasObject, searchedText string) bool {
	for _, existingWidget := range content.(*widget.Box).Children {
		if existingWidget.(*widget.Label).Text == searchedText {
			return true
		}
	}
	return false
}

func GetLabelFromContent(content fyne.CanvasObject, labelText string) *widget.Label {
	for _, existingWidget := range content.(*widget.Box).Children {
		if existingWidget.(*widget.Label).Text == labelText {
			return existingWidget.(*widget.Label)
		}
	}
	return nil
}

func GetLabelPositionInContent(content fyne.CanvasObject, labelText string) int {
	for position, existingWidget := range content.(*widget.Box).Children {
		if existingWidget.(*widget.Label).Text == labelText {
			return position
		}
	}
	return -1
}

func getInputHandlerForTesting() input.Handler {
	keymap := make(map[input.Action]input.KeyCombination)
	//Default keys are the same as if they were set by default config
	keymap[input.MoveDownAction] = input.SingleKeyCombination(fyne.KeyJ)
	keymap[input.MoveUpAction] = input.SingleKeyCombination(fyne.KeyK)
	keymap[input.EnterInputModeAction] = input.SingleKeyCombination(fyne.KeyI)
	keymap[input.ExitInputModeAction] = input.SingleKeyCombination(fyne.KeyEscape)
	keymap[input.ConfirmAction] = input.SingleKeyCombination(fyne.KeyReturn)
	keymap[input.CancelAction] = input.SingleKeyCombination(fyne.KeyEscape)
	return input.NewHandler(keymap)
}

func getOneInputFieldForDialogTesting() []*FormDialogFormItem {
	elements := []*FormDialogFormItem{}
	elements = append(elements, newFormDialogFormItem("first", NewInputField(test.Canvas(), getInputHandlerForTesting())))
	return elements
}

func getTwoInputFieldsForFormDialogTesting() []*FormDialogFormItem {
	elements := []*FormDialogFormItem{}
	elements = append(elements, newFormDialogFormItem("first", NewInputField(test.Canvas(), getInputHandlerForTesting())))
	elements = append(elements, newFormDialogFormItem("second", NewInputField(test.Canvas(), getInputHandlerForTesting())))
	return elements
}

func getThreeInputFieldsForFormDialogTesting() []*FormDialogFormItem {
	elements := []*FormDialogFormItem{}
	elements = append(elements, newFormDialogFormItem("first", NewInputField(test.Canvas(), getInputHandlerForTesting())))
	elements = append(elements, newFormDialogFormItem("second", NewInputField(test.Canvas(), getInputHandlerForTesting())))
	elements = append(elements, newFormDialogFormItem("third", NewInputField(test.Canvas(), getInputHandlerForTesting())))
	return elements
}
