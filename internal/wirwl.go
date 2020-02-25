package wirwl

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type App struct {
	window          fyne.Window
	currentTabLabel *widget.Label
}

func NewApp() *App {
	return &App{}
}

func (app *App) LoadAndDisplayGUI(fyneApp fyne.App) {
	fyneApp.Settings().SetTheme(theme.LightTheme())
	app.window = fyneApp.NewWindow("wirwl")
	app.currentTabLabel = widget.NewLabel("No data loaded")
	app.window.SetContent(app.currentTabLabel)
	app.window.ShowAndRun()
}
