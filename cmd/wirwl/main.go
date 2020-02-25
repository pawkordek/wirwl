package main

import (
	"fyne.io/fyne/app"
	wirwl "wirwl/internal"
)

func main() {
	wirwlApp := wirwl.NewApp()
	wirwlApp.LoadAndDisplayGUI(app.New())
}
