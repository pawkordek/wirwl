package wirwl

import (
	"fyne.io/fyne"
	fyneTest "fyne.io/fyne/test"
	"log"
	"os/user"
	"wirwl/internal/data"
)

/* It's a generated db file which contains:
comics:
	some comic1
	some comic2
music:
	some music1
	some music2
videos:
	some video1
	some video2
Should be copied to perform any tests that require data to exist.
*/
const testDbPath = "../testdata/test_db.db"

/* Path which should contain a copy of test db file. It should be used for testing any tests that require data to exists
 */
const testDbCopyPath = "../testdata/test_db_copy.db"

/* Doesn't exist on disk when tests are run and will be created by wirwl if it has been loaded and run.
Should be used for testing situations when application has been run for the first time.
*/
const emptyDbPath = "../testdata/empty_db.db"

func createTestDb() {
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

func deleteTestDb() {
	data.DeleteFile(testDbPath)
}

func setupAppForTesting() (*App, func()) {
	data.CopyFile(testDbPath, testDbCopyPath)
	app := NewApp(testDbCopyPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	return app, deleteTestDbCopy
}

func setupFirstRunAppForTesting() (*App, func()) {
	app := NewApp(emptyDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	return app, deleteEmptyTestDb
}

func deleteTestDbCopy() {
	data.DeleteFile(testDbCopyPath)
}

func deleteEmptyTestDb() {
	data.DeleteFile(emptyDbPath)
}

func getLoggingDirForTesting() string {
	currentUser, err := user.Current()
	if err != nil {
		return "/tmp/wirwl/"
	}
	return currentUser.HomeDir + "/.local/share/wirwl/"
}

func (app *App) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	focusedElement := app.mainWindow.Canvas().Focused()
	if focusedElement != nil {
		focusedElement.TypedKey(event)
	} else {
		onTypedKey := app.mainWindow.Canvas().OnTypedKey()
		onTypedKey(event)
	}
}

func (app *App) SimulateSwitchingToNextEntryType() {
	app.SimulateKeyPress(fyne.KeyL)
}

func (app *App) SimulateSwitchingToPreviousEntryType() {
	app.SimulateKeyPress(fyne.KeyH)
}

func (app *App) SimulateOpeningDialogForAddingEntryType() {
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
}

func (app *App) SimulateAddingNewEntryTypeWithName(text string) {
	app.SimulateOpeningDialogForAddingEntryType()
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type(text)
	app.SimulateKeyPress(fyne.KeyEnter)
}

func (app *App) SimulateSavingChanges() {
	app.SimulateKeyPress(fyne.KeyS)
	app.SimulateKeyPress(fyne.KeyY)
}

func (app *App) SimulateAttemptAtDeletionOfCurrentEntryType() {
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
}

func (app *App) SimulateDeletionOfCurrentEntryType() {
	app.SimulateAttemptAtDeletionOfCurrentEntryType()
	app.SimulateKeyPress(fyne.KeyY)
}

func (app *App) SimulateEditionOfCurrentEntryTypeTo(text string) {
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyE)
	app.SimulateKeyPress(fyne.KeyI)
	app.editEntryTypeDialog.Type(text)
	app.SimulateKeyPress(fyne.KeyEnter)
}
