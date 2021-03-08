package main

import (
	"flag"
	"fyne.io/fyne/v2/app"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	wirwl "wirwl/internal"
	"wirwl/internal/log"
)

func main() {
	flags := readCommandLineFlags()
	configurator := wirwl.NewAppConfigurator(flags["configDirPath"])
	config, err := configurator.LoadConfig()
	if err == nil {
		err = configurator.SetupNeededPaths(config)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		cleanup := configurator.SetupLoggerIn(config.AppDataDirPath)
		defer cleanup()
		dataProvider := configurator.LoadDataProvider(filepath.Join(config.AppDataDirPath, "data.db"))
		wirwlApp := wirwl.NewApp(app.New(), config, dataProvider, configurator.LoadingErrors())
		err = wirwlApp.LoadAndDisplay()
		if err != nil {
			err = errors.Wrap(err, "An error occurred when loading the application preventing it from continuing")
			log.Error(err)
			os.Exit(1)
		}
	} else {
		err = errors.Wrap(err, "A fatal error occurred. Application cannot continue")
		log.Error(err)
		os.Exit(1)
	}
}

func readCommandLineFlags() map[string]string {
	configDirPath := flag.String("c", "",
		"A path to a directory containing the application's config file. "+
			"If not provided it will default to: \n"+
			"On Unix: $XDG_CONFIG_HOME/wirwl/ and if this env variable is not set, then $HOME/.config/wirwl/\n"+
			"On Windows %AppData%/wirwl/")
	flag.Parse()
	return map[string]string{
		"configDirPath": *configDirPath,
	}
}
