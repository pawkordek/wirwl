package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"wirwl/internal/input"
)

/* When opening input in the running application, the last typed in character is still handled which normally means it
gets typed into the input. There is code that prevents this but as this situation doesn't happen when running the
test code any string typed into the input needs an additional character at the beginning as if the bug happened.
*/
func (inputField *InputField) Type(chars string) {
	fixedChars := " " + chars
	for _, char := range fixedChars {
		inputField.TypedRune(char)
	}
}

func (dialog *FormDialog) Type(chars string) {
	dialog.currentInput().Type(chars)
}

func SimulateKeyPress(focusable fyne.Focusable, key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	focusable.TypedKey(event)
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
