package main

import (
	"flag"
	"fyne.io/fyne/app"
	wirwl "wirwl/internal"
)

func main() {
	flags := readCommandLineFlags()
	wirwlApp := wirwl.NewApp(app.New())
	wirwlApp.LoadAndDisplay(flags["configDirPath"], flags["appDataDirPath"])
}

func readCommandLineFlags() map[string]string {
	configDirPath := flag.String("c", "", "A path to a directory containing the application config file. If not provided it will default to [HOME]/.local/share/wirwl/")
	appDataDirPath := flag.String("ad", "", "A path to a directory containing the application data. If not provided it will default to [HOME]/.local/share/wirwl/")
	flag.Parse()
	return map[string]string{
		"configDirPath":  *configDirPath,
		"appDataDirPath": *appDataDirPath,
	}
}
