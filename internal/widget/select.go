package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	fyneWidget "fyne.io/fyne/widget"
	"image/color"
	"wirwl/internal/input"
)

type Select struct {
	fyneWidget.Select
	backgroundRenderer *selectBackgroundRenderer
	canvas             fyne.Canvas
	inputHandler       input.Handler
	menu               *PopUpMenu
	focused            bool
	onExitInputMode    func()
}

type selectBackgroundRenderer struct {
	fyne.WidgetRenderer
	color color.Color
}

func (renderer *selectBackgroundRenderer) BackgroundColor() color.Color {
	return renderer.color
}

func (renderer *selectBackgroundRenderer) SetColor(color color.Color) {
	renderer.color = color
}

func (selectWidget *Select) CreateRenderer() fyne.WidgetRenderer {
	renderer := selectWidget.Select.CreateRenderer()
	bgRenderer := &selectBackgroundRenderer{renderer, theme.BackgroundColor()}
	selectWidget.backgroundRenderer = bgRenderer
	return bgRenderer
}

func NewSelect(canvas fyne.Canvas, handler input.Handler, choices ...string) *Select {
	menu := NewPopUpMenu(canvas, handler, choices...)
	selectWidget := &Select{
		backgroundRenderer: &selectBackgroundRenderer{},
		canvas:             canvas,
		inputHandler:       handler,
		menu:               menu,
		focused:            false,
		onExitInputMode:    func() {},
	}
	selectWidget.Options = choices
	menu.OnChoiceSelectedCallback = func(s string) {
		selectWidget.SetSelected(s)
		selectWidget.onExitInputMode()
	}
	selectWidget.ExtendBaseWidget(selectWidget)
	selectWidget.inputHandler.BindFunctionToAction(selectWidget, input.ExitInputModeAction, func() {
		selectWidget.canvas.Unfocus()
		selectWidget.onExitInputMode()
	})
	return selectWidget
}

func (selectWidget *Select) EnterInputMode() {
	selectWidget.menu.ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(selectWidget))
}

func (selectWidget *Select) FocusGained() {
	selectWidget.focused = true
}

func (selectWidget *Select) FocusLost() {
	selectWidget.focused = false
}

func (selectWidget *Select) Focused() bool {
	return selectWidget.focused
}

func (selectWidget *Select) TypedRune(r rune) {
	//Not handled as select doesn't do anything with typed in text
}

func (selectWidget *Select) TypedKey(event *fyne.KeyEvent) {
	selectWidget.inputHandler.HandleInNormalMode(selectWidget, event.Name)
}

func (selectWidget *Select) SetOnConfirm(f func()) {
	//Not handled as OnConfirm callbacks will be removed from forms elements in the future
}

func (selectWidget *Select) SetOnExitInputModeFunction(function func()) {
	selectWidget.onExitInputMode = function
}

func (selectWidget *Select) Highlight() {
	selectWidget.backgroundRenderer.SetColor(theme.FocusColor())
	selectWidget.Refresh()
}

func (selectWidget *Select) Unhighlight() {
	selectWidget.backgroundRenderer.SetColor(theme.BackgroundColor())
	selectWidget.Refresh()
}

func (selectWidget *Select) SetText(value string) {
	selectWidget.Selected = value
}

func (selectWidget *Select) GetText() string {
	return selectWidget.Selected
}
