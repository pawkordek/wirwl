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
	canvas               fyne.Canvas
	bgRenderer           *backgroundRenderer
	inputHandler         input.Handler
	OnConfirm            func()
	OnExitInputMode      func()
	runeAllowedToBeTyped func(r rune) bool
	firstRuneIgnored     bool
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
	newInput := newInputField(canvas, handler)
	newInput.ExtendBaseWidget(newInput)
	return newInput
}

//Should be used for dialog creation by any widget that embed this widget so it can properly extend fyne's BaseWidget
func newInputField(canvas fyne.Canvas, handler input.Handler) *InputField {
	newInput := &InputField{
		Entry:                widget.Entry{},
		canvas:               canvas,
		bgRenderer:           &backgroundRenderer{},
		inputHandler:         handler,
		OnConfirm:            func() {},
		OnExitInputMode:      func() {},
		runeAllowedToBeTyped: func(r rune) bool { return true },
		firstRuneIgnored:     false,
	}
	newInput.inputHandler.BindFunctionToAction(newInput, input.ConfirmAction, func() { newInput.OnConfirm() })
	newInput.inputHandler.BindFunctionToAction(newInput, input.ExitInputModeAction, func() { newInput.ExitInputMode() })
	return newInput
}

func (inputField *InputField) SetOnConfirm(function func()) {
	inputField.OnConfirm = function
}

func (inputField *InputField) SetOnExitInputModeFunction(function func()) {
	inputField.OnExitInputMode = function
}

func (inputField *InputField) SetRuneFilteringFunction(function func(r rune) bool) {
	inputField.runeAllowedToBeTyped = function
}

func (inputField *InputField) TypedKey(key *fyne.KeyEvent) {
	/*If TypedRune function ignores the first key, TypedKey has to do so as well, otherwise inputField will think
	there is one more character and crash the application on backspace press.
	This cannot be tested automatically as this doesn't happen there.
	*/
	if inputField.firstRuneIgnored {
		inputField.Entry.TypedKey(key)
	}
	handled, handlingResult := inputField.inputHandler.HandleInInputMode(inputField, key.Name)
	if handled && handlingResult.Action == input.ExitInputModeAction && handlingResult.KeyCombination.BothKeysPressed() {
		//When two key combination for exiting input mode gets pressed, it's first character has already been typed so it has to be removed
		//The second one doesn't have to be removed as it won't be typed in due to input field unfocusing when it exits
		inputField.removeLastCharacterFromText()
	}
}

func (inputField *InputField) removeLastCharacterFromText() {
	textLength := len(inputField.Text)
	if textLength > 0 {
		inputField.Text = inputField.Text[0 : textLength-1]
		inputField.Refresh()
	}
}

func (inputField *InputField) TypedRune(r rune) {
	//Prevents the last key pressed to display the hidden inputField from being typed into it
	if inputField.firstRuneIgnored && inputField.runeAllowedToBeTyped(r) {
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

func (inputField *InputField) Highlight() {
	inputField.setBgColor(theme.FocusColor())
}

func (inputField *InputField) Unhighlight() {
	inputField.setBgColor(theme.BackgroundColor())
}

func (inputField *InputField) EnterInputMode() {
	inputField.Unhighlight()
	inputField.canvas.Focus(inputField)
}

func (inputField *InputField) ExitInputMode() {
	inputField.canvas.Unfocus()
	inputField.OnExitInputMode()
}

func (inputField *InputField) GetText() string {
	return inputField.Text
}
