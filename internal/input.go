package wirwl

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type Input struct {
	widget.Entry
	OnEnterPressed   func()
	firstRuneIgnored bool
}

func newInput() *Input {
	input := &Input{}
	input.ExtendBaseWidget(input)
	return input
}

func (input *Input) SetOnEnterPressed(function func()) {
	input.OnEnterPressed = function
}

func (input *Input) TypedKey(key *fyne.KeyEvent) {
	input.Entry.TypedKey(key)
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
