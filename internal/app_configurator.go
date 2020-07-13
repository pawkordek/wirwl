package wirwl

import (
	"github.com/pkg/errors"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

//AppConfigurator is responsible for:
//- preparing all of the dependencies that will be passed to an App instance through a constructor
//- preparing the elements of the application that are needed for it to even start, i.e. directories

type AppConfigurator struct {
	configDirPath string
	loadingErrors map[string]string
}

func NewAppConfigurator(configDirPath string) AppConfigurator {
	return AppConfigurator{configDirPath: configDirPath, loadingErrors: make(map[string]string)}
}

func (configurator *AppConfigurator) LoadConfig() (Config, error) {
	config, err := NewConfig(configurator.configDirPath)
	if err != nil {
		return config, err
	}
	err = config.load()
	if err != nil {
		configurator.loadingErrors["config"] = "An error occurred when loading the config file in " + config.ConfigDirPath + "wirwl.cfg. Default config has been loaded instead."
		log.Error(err)
		config.loadDefaults()
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
