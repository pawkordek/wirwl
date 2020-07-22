package wirwl

import (
	fyneTest "fyne.io/fyne/test"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

/*
The purpose of TestAppConfigurator is to provide an object that can allow to quickly create the most commonly used
testing setups while also allowing to make small changes to such, if needed.
*/

type TestAppConfigurator struct {
	config       Config
	dataProvider data.Provider
	app          *App
}

func NewTestAppConfigurator() TestAppConfigurator {
	return TestAppConfigurator{}
}

//Always has to be run as these directories will be used by the test application, no matter the other settings
func (configurator *TestAppConfigurator) createTestDirectories() *TestAppConfigurator {
	err := data.CreateDirIfNotExist(testConfigDirPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to create test config directory in "+testConfigDirPath))
	}
	err = data.CreateDirIfNotExist(testAppDataDirPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to create app data directory in "+testAppDataDirPath))
	}
	return configurator
}

func (configurator *TestAppConfigurator) createTestDb() *TestAppConfigurator {
	err := data.CopyFile(testDbPath, testDbCopyPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to create test db by copying "+testDbPath+"into "+testDbCopyPath))
	}
	return configurator
}

func (configurator *TestAppConfigurator) createTestConfig() *TestAppConfigurator {
	config := NewConfig(testConfigDirPath)
	config.AppDataDirPath = testAppDataDirPath
	configurator.config = config
	return configurator
}

//Should be only called if config has been already created
func (configurator *TestAppConfigurator) createTestConfigFile() *TestAppConfigurator {
	configurator.config.saveConfigIn(configurator.config.ConfigDirPath + "wirwl.cfg")
	return configurator
}

func (configurator *TestAppConfigurator) createFailingToLoadConfigFile() *TestAppConfigurator {
	testConfigFile := filepath.Join(testConfigDirPath + "wirwl.cfg")
	nonsenseContents := []byte("qkrhqwroqwprhqr")
	//Having a config file with non-parsable contents will always cause an error when it gets loaded
	err := ioutil.WriteFile(testConfigFile, nonsenseContents, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return configurator
}

func (configurator *TestAppConfigurator) createDefaultDataProvider() *TestAppConfigurator {
	configurator.dataProvider = data.NewBoltProvider(filepath.Join(testAppDataDirPath, "data.db"))
	return configurator
}

func (configurator *TestAppConfigurator) createTestApplication() *TestAppConfigurator {
	app := NewApp(fyneTest.NewApp(), configurator.config, configurator.dataProvider)
	configurator.app = app
	return configurator
}

func (configurator *TestAppConfigurator) prepareConfiguratorForTestingWithExistingData() *TestAppConfigurator {
	configurator.
		createTestDirectories().
		createTestDb().
		createTestConfig().
		createTestConfigFile().
		createDefaultDataProvider()
	return configurator
}

func (configurator *TestAppConfigurator) createTestApplicationThatUsesExistingData() *TestAppConfigurator {
	configurator.prepareConfiguratorForTestingWithExistingData().createTestApplication()
	return configurator
}

func (configurator *TestAppConfigurator) createTestApplicationThatWillRunForFirstTime() *TestAppConfigurator {
	configurator.
		createTestDirectories().
		createDefaultDataProvider().
		createTestApplication()
	return configurator
}

func (configurator *TestAppConfigurator) getRunningTestApplication() (*App, func()) {
	err := configurator.app.LoadAndDisplay()
	if err != nil {
		log.Fatal(err)
	}
	return configurator.app, removeAllNonPersistentFilesInTestDataDir
}
