package widget

import (
	"fyne.io/fyne"
	"unicode"
	"wirwl/internal/input"
)

type NumericInputField struct {
	InputField
}

func NewNumericInputField(canvas fyne.Canvas, inputHandler input.Handler) *NumericInputField {
	//The order of things is a workaround for a possible bug in fyne where numeric input field doesn't refresh (draw correctly)
	//unless it's base widget is the one that is set (it can't be done by extend as no matter what the order of operations
	//is, input field as a parent will always extend it overriding whatever numeric input field sets)
	//Seems like a bug as other widgets that are extended multiple times don't have this problem
	//TODO: Check
	inputField := &NumericInputField{}
	inputField.ExtendBaseWidget(inputField)
	baseWidget := inputField.BaseWidget
	inputField.InputField = *NewInputField(canvas, inputHandler)
	inputField.BaseWidget = baseWidget
	inputField.SetRuneFilteringFunction(func(r rune) bool {
		return unicode.IsDigit(r)
	})
	return inputField
}
