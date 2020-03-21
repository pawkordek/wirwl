package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

type Input struct {
	widget.Entry
	bgRenderer       *backgroundRenderer
	OnEnterPressed   func()
	OnTypedKey       func(key *fyne.KeyEvent)
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

func NewInput() *Input {
	input := &Input{
		Entry:            widget.Entry{},
		OnEnterPressed:   func() {},
		OnTypedKey:       func(key *fyne.KeyEvent) {},
		firstRuneIgnored: false,
	}
	input.ExtendBaseWidget(input)
	return input
}

func (input *Input) SetOnEnterPressed(function func()) {
	input.OnEnterPressed = function
}

func (input *Input) SetOnTypedKey(function func(key *fyne.KeyEvent)) {
	input.OnTypedKey = function
}

func (input *Input) TypedKey(key *fyne.KeyEvent) {
	/*If TypedRune function ignores the first key, TypedKey has to do so as well, otherwise input will think
	there is one more character and crash the application on backspace press.
	This cannot be tested automatically as this doesn't happen there.
	*/
	if input.firstRuneIgnored {
		input.Entry.TypedKey(key)
	}
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		input.OnEnterPressed()
	}
	input.OnTypedKey(key)
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
