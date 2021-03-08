package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	fyneWidget "fyne.io/fyne/widget"
	"wirwl/internal/input"
)

/*
A pop up menu allowing a user to choose one of the choices presented vertically using up and down actions.
*/
type PopUpMenu struct {
	fyneWidget.PopUp
	focused                  bool
	inputHandler             input.Handler
	choices                  []*fyneWidget.Label
	currentChoiceNum         int
	OnChoiceSelectedCallback func(string)
}

func NewPopUpMenu(canvas fyne.Canvas, handler input.Handler, choicesNames ...string) *PopUpMenu {
	choices := generateChoicesFromNames(choicesNames...)
	content := container.NewVBox(choicesAsCanvasObjects(choices)...)
	menu := &PopUpMenu{
		PopUp: fyneWidget.PopUp{
			Content: content,
			Canvas:  canvas,
		},
		focused:                  false,
		inputHandler:             handler,
		choices:                  choices,
		currentChoiceNum:         0,
		OnChoiceSelectedCallback: func(s string) {},
	}
	menu.ExtendBaseWidget(menu)
	menu.inputHandler.BindFunctionToAction(menu, input.MoveDownAction, func() { menu.selectNextChoice() })
	menu.inputHandler.BindFunctionToAction(menu, input.MoveUpAction, func() { menu.selectPreviousChoice() })
	menu.inputHandler.BindFunctionToAction(menu, input.ConfirmAction, func() { menu.onChoiceSelected() })
	menu.currentChoice().TextStyle = fyne.TextStyle{Bold: true}
	return menu
}

func generateChoicesFromNames(choicesNames ...string) []*fyneWidget.Label {
	choices := []*fyneWidget.Label{}
	for _, itemName := range choicesNames {
		choices = append(choices, fyneWidget.NewLabel(itemName))
	}
	return choices
}

func choicesAsCanvasObjects(choices []*fyneWidget.Label) []fyne.CanvasObject {
	objects := []fyne.CanvasObject{}
	for _, choice := range choices {
		objects = append(objects, choice)
	}
	return objects
}

func (menu *PopUpMenu) currentChoice() *fyneWidget.Label {
	return menu.choices[menu.currentChoiceNum]
}

func (menu *PopUpMenu) selectNextChoice() {
	menu.selectChoice(menu.currentChoiceNum + 1)
}

func (menu *PopUpMenu) selectPreviousChoice() {
	menu.selectChoice(menu.currentChoiceNum - 1)
}

func (menu *PopUpMenu) selectChoice(num int) {
	if num >= 0 && num < len(menu.choices) {
		menu.unselectCurrentChoice()
		menu.currentChoiceNum = num
		menu.selectCurrentChoice()
	}
}

func (menu *PopUpMenu) unselectCurrentChoice() {
	menu.currentChoice().TextStyle = fyne.TextStyle{Bold: false}
	menu.currentChoice().Refresh()
}

func (menu *PopUpMenu) selectCurrentChoice() {
	menu.currentChoice().TextStyle = fyne.TextStyle{Bold: true}
	menu.currentChoice().Refresh()
}

func (menu *PopUpMenu) FocusGained() {
	menu.focused = true
}

func (menu *PopUpMenu) FocusLost() {
	menu.focused = false
}

func (menu *PopUpMenu) Focused() bool {
	return menu.focused
}

func (menu *PopUpMenu) TypedRune(r rune) {
	//Do nothing as text input is not needed in a menu
}

func (menu *PopUpMenu) TypedKey(key *fyne.KeyEvent) {
	menu.inputHandler.HandleInNormalMode(menu, key.Name)
}

func (menu *PopUpMenu) Show() {
	menu.PopUp.Show()
	menu.Canvas.Focus(menu)
}

func (menu *PopUpMenu) ShowAtPosition(position fyne.Position) {
	menu.PopUp.ShowAtPosition(position)
	menu.Canvas.Focus(menu)
}

func (menu *PopUpMenu) onChoiceSelected() {
	menu.Canvas.Unfocus()
	menu.Hide()
	menu.OnChoiceSelectedCallback(menu.currentChoice().Text)
}
