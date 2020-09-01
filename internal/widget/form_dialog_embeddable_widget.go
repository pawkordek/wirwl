package widget

import "fyne.io/fyne"

type FormDialogEmbeddableWidget interface {
	fyne.CanvasObject
	fyne.Focusable
	EnterInputMode()
	SetOnConfirm(func())
	SetOnExitInputModeFunction(func())
	Highlight()
	Unhighlight()
	SetText(value string)
	GetText() string
}
