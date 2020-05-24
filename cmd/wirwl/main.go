package main

import (
	"flag"
	"fyne.io/fyne/app"
	wirwl "wirwl/internal"
)

func main() {
	flags := readCommandLineFlags()
	config := wirwl.LoadConfigFromDir(flags["configDirPath"])
	wirwlApp := wirwl.NewApp(app.New(), config)
	wirwlApp.LoadAndDisplay()
}

func readCommandLineFlags() map[string]string {
	configDirPath := flag.String("c", "", "A path to a directory containing the application config file. If not provided it will default to [HOME]/.config/wirwl/")
	flag.Parse()
	return map[string]string{
		"configDirPath": *configDirPath,
	}
}
