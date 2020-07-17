package wirwl

import (
	"bytes"
	"fyne.io/fyne"
	fyneTest "fyne.io/fyne/test"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

/*File containing various functions allowing to setup/clean up the test environment
All path variables below should be treated as constants. They cannot be made const as they need to have slashes adapted
for cross platform compatibility.
*/

//Folder that should be used for storing any temporary data when testing and for storing directories used as paths in
//passed into application's Config file.
//It's contents are cleared after every test, except for folder 'persistent' and it's contents.
//The folder itself is created at the beginning before first test is run and removed after all tests are run
var testDataDirPath = filepath.FromSlash("../testdata/")

//Should be used for storing data that must persist between tests as it's the only folder in testdata directory
//which is not removed after each test's execution
var persistentTestDataDirPath = filepath.FromSlash(testDataDirPath + "persistent/")

//Used as application's data directory path when testing
var testAppDataDirPath = filepath.FromSlash(testDataDirPath + "app_data/")

//Used as application's config directory path when testing
var testConfigDirPath = filepath.FromSlash(testDataDirPath + "config/")

//Used as application's default data directory path when testing
var defaultTestAppDataDirPath = filepath.FromSlash(testDataDirPath + "default_app_data/")

//Used as application's default config directory path when testing
var defaultTestConfigDirPath = filepath.FromSlash(testDataDirPath + "default_config/")

/* It's a path to a database file which is generated every time tests are run but before any test executes.
If shown in the application, the data would look as follows:
comics:
	some comic1
	some comic2
music:
	some music1
	some music2
videos:
	some video1
	some video2
The file should be copied to perform any tests that require an existing data.
*/
var testDbPath = filepath.FromSlash(persistentTestDataDirPath + "data.db")

//Path used by tests to store a copy of test database file so that they don't affect the original file
var testDbCopyPath = filepath.FromSlash(testAppDataDirPath + "data.db")

func testSetup() {
	/*Cleanup is run in the case that a test crashed in the previous run and couldn't run it's cleanup functions, leaving
	potentially unwanted files and directories*/
	testCleanup()
	err := data.CreateDirIfNotExist(testDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
	createTestDb()
}

func createTestDb() {
	err := data.CreateDirIfNotExist(persistentTestDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
	dataProvider := data.NewBoltProvider(testDbPath)
	saveTestEntriesTypes(dataProvider)
	saveTestComics(dataProvider)
	saveTestMusic(dataProvider)
	saveTestVideos(dataProvider)
}

func saveTestEntriesTypes(provider data.Provider) {
	entriesTypes := data.GetEntriesTypes()
	err := provider.SaveEntriesTypesToDb(entriesTypes)
	if err != nil {
		log.Fatal(err)
	}
}

func saveTestVideos(provider data.Provider) {
	videos := data.GetExampleVideoEntries()
	err := provider.SaveEntriesToDb("videos", videos)
	if err != nil {
		log.Fatal(err)
	}
}

func saveTestComics(provider data.Provider) {
	comics := data.GetExampleComicEntries()
	err := provider.SaveEntriesToDb("comics", comics)
	if err != nil {
		log.Fatal(err)
	}
}

func saveTestMusic(provider data.Provider) {
	music := data.GetExampleMusicEntries()
	err := provider.SaveEntriesToDb("music", music)
	if err != nil {
		log.Fatal(err)
	}
}

func testCleanup() {
	err := data.DeleteDirWithContents(testDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
}

func removeAllNonPersistentFilesInTestDataDir() {
	err := data.DeleteAllInDirExceptForDirs(testDataDirPath, "persistent")
	if err != nil {
		log.Fatal(err)
	}
}

func setupAndRunAppForTestingWithFailingToLoadConfig() (*App, func()) {
	err := data.CreateDirIfNotExist(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	testConfigFile := filepath.Join(testConfigDirPath + "wirwl.cfg")
	nonsenseContents := []byte("qkrhqwroqwprhqr")
	//Having a config file with non-parsable contents will always cause an error when it gets loaded
	err = ioutil.WriteFile(testConfigFile, nonsenseContents, 0644)
	if err != nil {
		log.Fatal(err)
	}
	config, err := NewConfig(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	app := NewApp(fyneTest.NewApp(), config, config.loadDataProvider())
	app.LoadAndDisplay()
	return app, removeAllNonPersistentFilesInTestDataDir
}

func setupAndRunAppForTestingWithTestConfig() (*App, func()) {
	err := data.CreateDirIfNotExist(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	config, err := getTestConfigWithConfigPathIn(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	config.AppDataDirPath = testAppDataDirPath
	config.saveConfigIn(config.ConfigDirPath + "wirwl.cfg")
	return setupAndRunAppForTestingWithExistingTestData(config)
}

func setupAndRunAppForTestingWithExistingTestData(config Config) (*App, func()) {
	data.CreateDirIfNotExist(testAppDataDirPath)
	err := data.CopyFile(testDbPath, testDbCopyPath)
	if err != nil {
		log.Fatal(err)
	}
	app := NewApp(fyneTest.NewApp(), config, config.loadDataProvider())
	app.LoadAndDisplay()
	return app, removeAllNonPersistentFilesInTestDataDir
}

func setupAndRunAppAsIfRunForFirstTime() (*App, func()) {
	config, err := getTestConfigWithConfigPathIn(testConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}
	config.defaultAppDataDirPath = defaultTestAppDataDirPath
	app := NewApp(fyneTest.NewApp(), config, config.loadDataProvider())
	app.LoadAndDisplay()
	return app, removeAllNonPersistentFilesInTestDataDir
}

func getTestConfigWithConfigPathIn(path string) (Config, error) {
	return NewConfig(path)
}

func areFilesInPathsTheSame(filePath1 string, filePath2 string) bool {
	file1, err := os.Open(filePath1)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	file2, err := os.Open(filePath2)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	return areFilesTheSame(file1, file2)
}

func areFilesTheSame(file1 *os.File, file2 *os.File) bool {
	const chunkSize = 4000
	for {
		bytesOfFile1 := make([]byte, chunkSize)
		_, err1 := file1.Read(bytesOfFile1)
		bytesOfFile2 := make([]byte, chunkSize)
		_, err2 := file2.Read(bytesOfFile2)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Error(err1)
				log.Fatal(err2)
			}
		}
		if !bytes.Equal(bytesOfFile1, bytesOfFile2) {
			return false
		}
	}
	return true
}

func createCorrectSavedWirwlConfigFileInPath(path string) {
	err := data.CreateDirIfNotExist(path)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		AppDataDirPath: testAppDataDirPath,
		ConfigDirPath:  testConfigDirPath,
	}
	config.saveConfigIn(path + "wirwl_correct.cfg")
}

func createCorrectWirwlConfigFileForLoadingInPath(path string) {
	err := data.CreateDirIfNotExist(path)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		AppDataDirPath: "some db path",
		ConfigDirPath:  testConfigDirPath,
	}
	config.saveConfigIn(path + "wirwl.cfg")
}

func (config *Config) saveConfigIn(configFilePath string) {
	configFile, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	err = toml.NewEncoder(configFile).Encode(config)
	if err != nil {
		log.Fatal(err)
	}
	err = configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) simulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	focusedElement := app.mainWindow.Canvas().Focused()
	if focusedElement != nil {
		focusedElement.TypedKey(event)
	} else {
		onTypedKey := app.mainWindow.Canvas().OnTypedKey()
		onTypedKey(event)
	}
}

func (app *App) simulateSwitchingToNextEntryType() {
	app.simulateKeyPress(fyne.KeyL)
}

func (app *App) simulateSwitchingToPreviousEntryType() {
	app.simulateKeyPress(fyne.KeyH)
}

func (app *App) simulateOpeningDialogForAddingEntryType() {
	app.simulateKeyPress(fyne.KeyT)
	app.simulateKeyPress(fyne.KeyI)
}

func (app *App) simulateAddingNewEntryTypeWithName(text string) {
	app.simulateOpeningDialogForAddingEntryType()
	app.simulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type(text)
	app.simulateKeyPress(fyne.KeyEnter)
}

func (app *App) simulateSavingChanges() {
	app.simulateKeyPress(fyne.KeyS)
	app.simulateKeyPress(fyne.KeyY)
}

func (app *App) simulateAttemptAtDeletionOfCurrentEntryType() {
	app.simulateKeyPress(fyne.KeyT)
	app.simulateKeyPress(fyne.KeyD)
}

func (app *App) simulateDeletionOfCurrentEntryType() {
	app.simulateAttemptAtDeletionOfCurrentEntryType()
	app.simulateKeyPress(fyne.KeyY)
}

func (app *App) simulateEditionOfCurrentEntryTypeTo(text string) {
	app.simulateKeyPress(fyne.KeyT)
	app.simulateKeyPress(fyne.KeyE)
	app.simulateKeyPress(fyne.KeyI)
	app.editEntryTypeDialog.Type(text)
	app.simulateKeyPress(fyne.KeyEnter)
}
