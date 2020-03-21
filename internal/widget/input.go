package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type Input struct {
	widget.Entry
	OnEnterPressed   func()
	firstRuneIgnored bool
}

func NewInput() *Input {
	input := &Input{
		Entry:            widget.Entry{},
		OnEnterPressed:   func() {},
		firstRuneIgnored: false,
	}
	input.ExtendBaseWidget(input)
	return input
}

func (input *Input) SetOnEnterPressed(function func()) {
	input.OnEnterPressed = function
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
