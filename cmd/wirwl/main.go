package main

import (
	"fyne.io/fyne/app"
	wirwl "wirwl/internal"
)

func main() {
	wirwlApp := wirwl.NewApp("exampleDb.db")
	wirwlApp.LoadAndDisplay(app.New())
}
