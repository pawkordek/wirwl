package widget

import "fyne.io/fyne"

func (input *Input) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	input.TypedKey(event)
}

/* When opening input in the running application, the last typed in character is still handled which normally means it
gets typed into the input. There is code that prevents this but as this situation doesn't happen when running the
test code any string typed into the input needs an additional character at the beginning as if the bug happened.
*/
func (input *Input) Type(chars string) {
	fixedChars := " " + chars
	for _, char := range fixedChars {
		input.TypedRune(char)
	}
}

func (dialog *ConfirmationDialog) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	dialog.TypedKey(event)
}
