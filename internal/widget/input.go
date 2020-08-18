package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"wirwl/internal/input"
)

type Input struct {
	widget.Entry
	bgRenderer       *backgroundRenderer
	inputHandler     input.InputHandler
	OnConfirm        func()
	OnCancel         func()
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

func (input *Input) CreateRenderer() fyne.WidgetRenderer {
	renderer := input.Entry.CreateRenderer()
	bgRenderer := &backgroundRenderer{renderer, theme.BackgroundColor()}
	input.bgRenderer = bgRenderer
	return bgRenderer
}

func NewInput(handler input.InputHandler) *Input {
	newInput := &Input{
		Entry:            widget.Entry{},
		bgRenderer:       &backgroundRenderer{},
		inputHandler:     handler,
		OnConfirm:        func() {},
		OnCancel:         func() {},
		firstRuneIgnored: false,
	}
	newInput.ExtendBaseWidget(newInput)
	newInput.inputHandler.BindFunctionToAction(newInput, input.ConfirmAction, func() { newInput.OnConfirm() })
	newInput.inputHandler.BindFunctionToAction(newInput, input.CancelAction, func() { newInput.OnCancel() })
	return newInput
}

func (input *Input) SetOnConfirm(function func()) {
	input.OnConfirm = function
}

func (input *Input) SetOnCancel(function func()) {
	input.OnCancel = function
}

func (input *Input) TypedKey(key *fyne.KeyEvent) {
	/*If TypedRune function ignores the first key, TypedKey has to do so as well, otherwise input will think
	there is one more character and crash the application on backspace press.
	This cannot be tested automatically as this doesn't happen there.
	*/
	if input.firstRuneIgnored {
		input.Entry.TypedKey(key)
	}
	input.inputHandler.Handle(input, key.Name)
}

func (input *Input) TypedRune(r rune) {
	//Prevents the last key pressed to display the hidden input from being typed into it
	if input.firstRuneIgnored {
		input.Entry.TypedRune(r)
	} else {
		input.firstRuneIgnored = true
	}
}

func (input *Input) FocusGained() {
	input.firstRuneIgnored = false
	input.Entry.FocusGained()
}

func (input *Input) setBgColor(color color.Color) {
	input.bgRenderer.SetColor(color)
	input.Refresh()
}

func (input *Input) Mark() {
	input.setBgColor(theme.FocusColor())
}

func (input *Input) Unmark() {
	input.setBgColor(theme.BackgroundColor())
}
