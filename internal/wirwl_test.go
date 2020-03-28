package wirwl

import (
	"fyne.io/fyne"
	fyneTest "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/widget"
)

const exampleDbPath = "../test/exampleDb.db"
const emptyDbPath = "../test/emptyDb.db"
const saveTestDbPath = "../test/saveTestDb.db"
const deletionTestDbPath = "../test/deletionTestDb.db"

func TestMain(m *testing.M) {
	dataProvider := data.NewBoltProvider(exampleDbPath)
	entriesTypes := data.GetEntriesTypes()
	err := dataProvider.SaveEntriesTypesToDb(entriesTypes)
	if err != nil {
		log.Fatal(err)
	}
	videos := data.GetExampleVideoEntries()
	err = dataProvider.SaveEntriesToDb("videos", videos)
	if err != nil {
		log.Fatal(err)
	}
	comics := data.GetExampleComicEntries()
	err = dataProvider.SaveEntriesToDb("comics", comics)
	if err != nil {
		log.Fatal(err)
	}
	music := data.GetExampleMusicEntries()
	err = dataProvider.SaveEntriesToDb("music", music)
	if err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	data.DeleteFile(exampleDbPath)
	os.Exit(exitCode)
}

func TestThatEntriesTabsWithContentDisplayInCorrectOrder(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, 3, len(app.entriesTabContainer.Items))
	assert.Equal(t, app.entriesTabContainer.Items[0].Text, "comics")
	assert.Equal(t, "some comic1", app.entriesLabels["comics"][0].Text)
	assert.Equal(t, "some comic2", app.entriesLabels["comics"][1].Text)
	assert.Equal(t, app.entriesTabContainer.Items[1].Text, "music")
	assert.Equal(t, "some video1", app.entriesLabels["videos"][0].Text)
	assert.Equal(t, "some video2", app.entriesLabels["videos"][1].Text)
	assert.Equal(t, app.entriesTabContainer.Items[2].Text, "videos")
	assert.Equal(t, "some music1", app.entriesLabels["music"][0].Text)
	assert.Equal(t, "some music2", app.entriesLabels["music"][1].Text)
}

func TestSwitchingTabs(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, "comics", app.getCurrentTabText())
	app.SimulateKeyPress(fyne.KeyL)
	assert.Equal(t, "music", app.getCurrentTabText())
	app.SimulateKeyPress(fyne.KeyL)
	assert.Equal(t, "videos", app.getCurrentTabText())
	app.SimulateKeyPress(fyne.KeyL)
	assert.Equal(t, "comics", app.getCurrentTabText())
	app.SimulateKeyPress(fyne.KeyH)
	assert.Equal(t, "videos", app.getCurrentTabText())
	app.SimulateKeyPress(fyne.KeyH)
	assert.Equal(t, "music", app.getCurrentTabText())
	app.SimulateKeyPress(fyne.KeyH)
	assert.Equal(t, "comics", app.getCurrentTabText())
}

func TestEntryHighlightingWhenSwitchingTabs(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, fyne.TextStyle{Bold: true}, app.entriesLabels["comics"][0].TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, app.entriesLabels["comics"][1].TextStyle)
	app.SimulateKeyPress(fyne.KeyL)
	assert.Equal(t, fyne.TextStyle{Bold: true}, app.entriesLabels["music"][0].TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, app.entriesLabels["music"][1].TextStyle)
	app.SimulateKeyPress(fyne.KeyH)
	assert.Equal(t, fyne.TextStyle{Bold: true}, app.entriesLabels["comics"][0].TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, app.entriesLabels["comics"][1].TextStyle)
}

func TestThatApplicationDoesNotCrashWhenTryingToSwitchToATabThatDoesNotExist(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyL)
	app.SimulateKeyPress(fyne.KeyL)
	app.SimulateKeyPress(fyne.KeyH)
	app.SimulateKeyPress(fyne.KeyH)
}

func TestThatIfThereAreNoEntriesCorrectMessageDisplays(t *testing.T) {
	app := NewApp(emptyDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, 1, len(app.entriesTabContainer.Items))
	assert.Equal(t, "No entries", app.entriesTabContainer.Items[0].Text)
}

func TestWhetherDialogForAddingEntryTypesOpens(t *testing.T) {
	app := NewApp(emptyDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, true, app.addEntryTypeDialog.Hidden)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	assert.Equal(t, true, app.addEntryTypeDialog.Visible())
	assert.Equal(t, true, app.addEntryTypeDialog.Focused())
}

func TestWhetherReopeningDialogForAddingEntriesTypesDoesNotPersistPreviouslyInputText(t *testing.T) {
	app := NewApp(emptyDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("some type")
	app.SimulateKeyPress(fyne.KeyReturn)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	name, _ := app.addEntryTypeDialog.GetItemValue("Name")
	imageQuery, _ := app.addEntryTypeDialog.GetItemValue("Image query")
	assert.Empty(t, name)
	assert.Empty(t, imageQuery)
}

func TestAddingOfNewEntryType(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("new entry type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	assert.Equal(t, 4, len(app.entriesTabContainer.Items))
	assert.Equal(t, "comics", app.entriesTabContainer.Items[0].Text)
	assert.Equal(t, "music", app.entriesTabContainer.Items[1].Text)
	assert.Equal(t, "new entry type", app.entriesTabContainer.Items[2].Text)
	assert.Equal(t, "videos", app.entriesTabContainer.Items[3].Text)
	assert.Equal(t, true, app.addEntryTypeDialog.Hidden)
	assert.Equal(t, true, !app.addEntryTypeDialog.Focused())
}

func TestThatItIsNotPossibleToAddTheSameEntryTypeTwice(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	assert.Equal(t, true, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "Entry type with name 'type' already exists.", app.msgDialog.Msg())
	assert.Equal(t, 4, len(app.entriesTabContainer.Items))
}

func TestThatPressingAnyKeyClosesMessagePopUp(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.msgDialog.Show()
	widget.SimulateKeyPress(app.msgDialog, fyne.KeyT)
	assert.Equal(t, true, app.msgDialog.Hidden)
	app.msgDialog.Show()
	widget.SimulateKeyPress(app.msgDialog, fyne.KeyReturn)
	assert.Equal(t, true, app.msgDialog.Hidden)
}

func TestThatAddingNewEntryTypeDoesNotChangeCurrentlyOpenedTab(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyL)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	assert.Equal(t, "music", app.getCurrentTabText())
}

func TestThatAfterAddingNewEntryOpenedTabStillHasTheSameElementHighlighted(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	assert.Equal(t, fyne.TextStyle{Bold: true}, app.entriesLabels["comics"][0].TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, app.entriesLabels["comics"][1].TextStyle)
}

func TestThatSavingChangesWorks(t *testing.T) {
	data.DeleteFile(saveTestDbPath)
	app := NewApp(saveTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyS)
	app = NewApp(saveTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, 1, len(app.entriesTabContainer.Items))
	assert.Equal(t, "type", app.getCurrentTabText())
}

func TestThatAfterSavingSuccessfullySuccessDialogDisplays(t *testing.T) {
	data.DeleteFile(saveTestDbPath)
	app := NewApp(saveTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyS)
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "SUCCESS", app.msgDialog.Title())
	assert.Equal(t, "Changes saved.", app.msgDialog.Msg())
}

func TestThatAfterSavingUnsuccessfullyErrorDialogDisplays(t *testing.T) {
	data.DeleteFile(saveTestDbPath)
	app := NewApp(saveTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.dataProvider = data.NewAlwaysFailingProvider()
	app.SimulateKeyPress(fyne.KeyS)
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, data.AlwaysFailingProviderError.Error(), app.msgDialog.Msg())
}

func TestThatDeletingEntriesTypesWorks(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	widget.SimulateKeyPress(app.confirmationDialog, fyne.KeyY)
	assert.Equal(t, 2, len(app.entriesTabContainer.Items))
	assert.Equal(t, "music", app.entriesTabContainer.Items[0].Text)
	assert.Equal(t, "videos", app.entriesTabContainer.Items[1].Text)
	data.DeleteFile(deletionTestDbPath)
}

func TestThatWhenTryingToDeleteLastEntryTypeItIsPreventedAndWarningDialogIsDisplayed(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	widget.SimulateKeyPress(app.confirmationDialog, fyne.KeyY)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	widget.SimulateKeyPress(app.confirmationDialog, fyne.KeyY)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	assert.Equal(t, 1, len(app.entriesTabContainer.Items))
	assert.Equal(t, true, app.msgDialog.Visible())
	assert.Equal(t, "WARNING", app.msgDialog.Title())
	assert.Equal(t, "You cannot remove the only remaining entry type!", app.msgDialog.Msg())
	data.DeleteFile(deletionTestDbPath)
}

func TestThatEditingEntryTypeWorks(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyE)
	widget.SimulateKeyPress(app.editEntryTypeDialog, fyne.KeyI)
	app.editEntryTypeDialog.Type("2")
	widget.SimulateKeyPress(app.editEntryTypeDialog, fyne.KeyEnter)
	assert.Equal(t, "2comics", app.entriesTabContainer.CurrentTab().Text)
}

func TestThatEditingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyE)
	widget.SimulateKeyPress(app.editEntryTypeDialog, fyne.KeyI)
	app.editEntryTypeDialog.Type("2")
	widget.SimulateKeyPress(app.editEntryTypeDialog, fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyS)
	app2 := NewApp(deletionTestDbPath)
	app2.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, "2comics", app2.entriesTabContainer.CurrentTab().Text)
	data.DeleteFile(deletionTestDbPath)
}

func TestThatDeletingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	widget.SimulateKeyPress(app.confirmationDialog, fyne.KeyY)
	app.SimulateKeyPress(fyne.KeyS)
	app2 := NewApp(deletionTestDbPath)
	app2.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, "music", app2.entriesTabContainer.CurrentTab().Text)
	data.DeleteFile(deletionTestDbPath)
}

func (app *App) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	onTypedKey := app.mainWindow.Canvas().OnTypedKey()
	onTypedKey(event)
}
