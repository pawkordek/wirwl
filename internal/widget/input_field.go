package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"wirwl/internal/input"
)

type InputField struct {
	widget.Entry
	canvas           fyne.Canvas
	bgRenderer       *backgroundRenderer
	inputHandler     input.Handler
	OnConfirm        func()
	OnExitInputMode  func()
	firstRuneIgnored bool
}

type backgroundRenderer struct {
	fyne.WidgetRenderer
	color color.Color
}

func (renderer *backgroundRenderer) BackgroundColor() color.Color {
	return renderer.color
}

func (renderer *backgroundRenderer) SetColor(color color.Color) {
	renderer.color = color
}

func (inputField *InputField) CreateRenderer() fyne.WidgetRenderer {
	renderer := inputField.Entry.CreateRenderer()
	bgRenderer := &backgroundRenderer{renderer, theme.BackgroundColor()}
	inputField.bgRenderer = bgRenderer
	return bgRenderer
}

func NewInputField(canvas fyne.Canvas, handler input.Handler) *InputField {
	newInput := &InputField{
		Entry:            widget.Entry{},
		canvas:           canvas,
		bgRenderer:       &backgroundRenderer{},
		inputHandler:     handler,
		OnConfirm:        func() {},
		OnExitInputMode:  func() {},
		firstRuneIgnored: false,
	}
	newInput.ExtendBaseWidget(newInput)
	newInput.inputHandler.BindFunctionToAction(newInput, input.ConfirmAction, func() { newInput.OnConfirm() })
	newInput.inputHandler.BindFunctionToAction(newInput, input.ExitInputModeAction, func() { newInput.OnExitInputMode() })
	return newInput
}

func (inputField *InputField) SetOnConfirm(function func()) {
	inputField.OnConfirm = function
}

func (inputField *InputField) SetOnExitInputModeFunction(function func()) {
	inputField.OnExitInputMode = function
}

func (inputField *InputField) TypedKey(key *fyne.KeyEvent) {
	/*If TypedRune function ignores the first key, TypedKey has to do so as well, otherwise inputField will think
	there is one more character and crash the application on backspace press.
	This cannot be tested automatically as this doesn't happen there.
	*/
	if inputField.firstRuneIgnored {
		inputField.Entry.TypedKey(key)
	}
	inputField.inputHandler.HandleInInputMode(inputField, key.Name)
}

func (inputField *InputField) TypedRune(r rune) {
	//Prevents the last key pressed to display the hidden inputField from being typed into it
	if inputField.firstRuneIgnored {
		inputField.Entry.TypedRune(r)
	} else {
		inputField.firstRuneIgnored = true
	}
}

func (inputField *InputField) FocusGained() {
	inputField.firstRuneIgnored = false
	inputField.Entry.FocusGained()
}

func (inputField *InputField) setBgColor(color color.Color) {
	inputField.bgRenderer.SetColor(color)
	inputField.Refresh()
}

func (inputField *InputField) Mark() {
	inputField.setBgColor(theme.FocusColor())
}

func (inputField *InputField) Unmark() {
	inputField.setBgColor(theme.BackgroundColor())
}

func (inputField *InputField) EnterInputMode() {
	inputField.Unmark()
	inputField.canvas.Focus(inputField)
}
