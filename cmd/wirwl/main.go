package main

import (
	"flag"
	"fyne.io/fyne/app"
	"github.com/pkg/errors"
	"os"
	wirwl "wirwl/internal"
	"wirwl/internal/log"
)

func main() {
	flags := readCommandLineFlags()
	config, err := wirwl.NewConfig(flags["configDirPath"])
	if err == nil {
		wirwlApp := wirwl.NewApp(app.New(), config)
		wirwlApp.LoadAndDisplay()
	} else {
		err = errors.Wrap(err, "A fatal error occurred. Application cannot continue")
		log.Error(err)
		os.Exit(1)
	}
}

func readCommandLineFlags() map[string]string {
	configDirPath := flag.String("c", "", "A path to a directory containing the application config file. If not provided it will default to [HOME]/.config/wirwl/")
	flag.Parse()
	return map[string]string{
		"configDirPath": *configDirPath,
	}
}
