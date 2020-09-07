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
	inputField := &NumericInputField{
		InputField: *NewInputField(canvas, inputHandler),
	}
	inputField.SetRuneFilteringFunction(func(r rune) bool {
		return unicode.IsDigit(r)
	})
	return inputField
}
