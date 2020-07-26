package wirwl

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

//AppConfigurator is responsible for:
//- preparing all of the dependencies that will be passed to an App instance through a constructor
//- preparing the elements of the application that are needed for it to even start, i.e. directories

type AppConfigurator struct {
	configDirPath string
	logFile       *os.File
	loadingErrors map[string]string
}

func NewAppConfigurator(configDirPath string) AppConfigurator {
	return AppConfigurator{configDirPath: configDirPath, loadingErrors: make(map[string]string)}
}

func (configurator *AppConfigurator) LoadConfig() (Config, error) {
	config := NewConfig(configurator.configDirPath)
	err := config.load()
	if err != nil {
		configurator.loadingErrors[configLoadError] = "An error occurred when loading the config file in " + config.ConfigDirPath + "wirwl.cfg. Default config has been loaded instead."
		log.Error(err)
		err = config.loadDefaults()
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

func (configurator *AppConfigurator) LoadDataProvider(dbPath string) data.Provider {
	return data.NewBoltProvider(dbPath)
}

func (configurator *AppConfigurator) SetupNeededPaths(config Config) error {
	err := data.CreateDirIfNotExist(config.AppDataDirPath)
	if err != nil {
		return errors.Wrap(err, "Failed to create application directory in "+config.AppDataDirPath+". Application must exit")
	}
	return nil
}

func (configurator *AppConfigurator) SetupLoggerIn(loggingPath string) func() {
	logFilePath := filepath.Join(loggingPath, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		err = errors.Wrap(err, "Failed to open the logfile in path "+logFilePath)
		log.Error(err)
	} else {
		writer := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(writer)
	}
	configurator.logFile = logFile
	return configurator.CloseLogFile
}

func (configurator *AppConfigurator) CloseLogFile() {
	log.SetOutput(os.Stdout)
	err := configurator.logFile.Close()
	if err != nil {
		log.Error(err)
	}
}

func (configurator *AppConfigurator) LoadingErrors() map[string]string {
	return configurator.loadingErrors
}
