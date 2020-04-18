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

/* Contains data generated in TestMain which is as follows:
comics:
	some comic1
	some comic2
music:
	some music1
	some music2
videos:
	some video1
	some video2
Should be used for testing operations that don't change it's contents
*/
const exampleDbPath = "../test/exampleDb.db"

/* Contains no data. Should be used for testing situations when application has been run for the first time and operations
won't change it's contents
*/
const emptyDbPath = "../test/emptyDb.db"

/* Contains no data. Should be used for testing saving operations. Should be deleted before running a test. */
const saveTestDbPath = "../test/saveTestDb.db"

/* Should be made as a copy of exampleDb. Should be used for testing deleting operations. Should be deleted after running a test*/
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
	assert.Equal(t, 3, len(app.entriesTypesTabs.Items()))
	firstTab := app.entriesTypesTabs.Items()[0]
	assert.Equal(t, firstTab.Text, "comics")
	assert.Equal(t, 0, widget.GetLabelPositionInContent(firstTab.Content, "some comic1"))
	assert.Equal(t, 1, widget.GetLabelPositionInContent(firstTab.Content, "some comic2"))
	secondTab := app.entriesTypesTabs.Items()[1]
	assert.Equal(t, secondTab.Text, "music")
	assert.Equal(t, 0, widget.GetLabelPositionInContent(secondTab.Content, "some music1"))
	assert.Equal(t, 1, widget.GetLabelPositionInContent(secondTab.Content, "some music2"))
	thirdTab := app.entriesTypesTabs.Items()[2]
	assert.Equal(t, thirdTab.Text, "videos")
	assert.Equal(t, 0, widget.GetLabelPositionInContent(thirdTab.Content, "some video1"))
	assert.Equal(t, 1, widget.GetLabelPositionInContent(thirdTab.Content, "some video2"))
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
	currentTab := app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
	app.SimulateKeyPress(fyne.KeyL)
	currentTab = app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some music1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some music2").TextStyle)
	app.SimulateKeyPress(fyne.KeyH)
	currentTab = app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
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
	assert.Equal(t, 1, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "No entries", app.entriesTypesTabs.Items()[0].Text)
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
	app.SimulateKeyPress(fyne.KeyI)
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
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type("new entry type")
	app.SimulateKeyPress(fyne.KeyEnter)
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "comics", app.entriesTypesTabs.Items()[0].Text)
	assert.Equal(t, "music", app.entriesTypesTabs.Items()[1].Text)
	assert.Equal(t, "new entry type", app.entriesTypesTabs.Items()[2].Text)
	assert.Equal(t, "videos", app.entriesTypesTabs.Items()[3].Text)
	assert.Equal(t, true, app.addEntryTypeDialog.Hidden)
	assert.Equal(t, true, !app.addEntryTypeDialog.Focused())
}

func TestThatItIsNotPossibleToAddTheSameEntryTypeTwice(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	app.SimulateKeyPress(fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	app.SimulateKeyPress(fyne.KeyEnter)
	assert.Equal(t, true, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "Entry type with name 'type' already exists.", app.msgDialog.Msg())
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
}

func TestThatPressingAnyKeyClosesMessagePopUp(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.msgDialog.Show()
	app.SimulateKeyPress(fyne.KeyT)
	assert.Equal(t, true, app.msgDialog.Hidden)
	app.msgDialog.Show()
	app.SimulateKeyPress(fyne.KeyReturn)
	assert.Equal(t, true, app.msgDialog.Hidden)
}

func TestThatAddingNewEntryTypeDoesNotChangeCurrentlyOpenedTab(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyL)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	widget.SimulateKeyPress(app.addEntryTypeDialog, fyne.KeyEnter)
	assert.Equal(t, "music", app.getCurrentTabText())
}

func TestThatAfterAddingNewEntryOpenedTabStillHasTheSameElementHighlighted(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	app.SimulateKeyPress(fyne.KeyEnter)
	currentTab := app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
}

func TestThatSavingChangesWorks(t *testing.T) {
	data.DeleteFile(saveTestDbPath)
	app := NewApp(saveTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type("type")
	app.SimulateKeyPress(fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyS)
	app = NewApp(saveTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, 1, len(app.entriesTypesTabs.Items()))
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
	app.SimulateKeyPress(fyne.KeyY)
	assert.Equal(t, 2, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "music", app.entriesTypesTabs.Items()[0].Text)
	assert.Equal(t, "videos", app.entriesTypesTabs.Items()[1].Text)
	data.DeleteFile(deletionTestDbPath)
}

func TestThatWhenTryingToDeleteLastEntryTypeItIsPreventedAndWarningDialogIsDisplayed(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	app.SimulateKeyPress(fyne.KeyY)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	app.SimulateKeyPress(fyne.KeyY)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	assert.Equal(t, 1, len(app.entriesTypesTabs.Items()))
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
	app.SimulateKeyPress(fyne.KeyI)
	app.editEntryTypeDialog.Type("2")
	app.SimulateKeyPress(fyne.KeyEnter)
	assert.Equal(t, "2comics", app.entriesTypesTabs.CurrentTab().Text)
}

func TestThatEditingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyE)
	app.SimulateKeyPress(fyne.KeyI)
	app.editEntryTypeDialog.Type("2")
	app.SimulateKeyPress(fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyS)
	app2 := NewApp(deletionTestDbPath)
	app2.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, "2comics", app2.entriesTypesTabs.CurrentTab().Text)
	data.DeleteFile(deletionTestDbPath)
}

func TestThatDeletingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	data.CopyFile(exampleDbPath, deletionTestDbPath)
	app := NewApp(deletionTestDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
	app.SimulateKeyPress(fyne.KeyY)
	app.SimulateKeyPress(fyne.KeyS)
	app2 := NewApp(deletionTestDbPath)
	app2.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, "music", app2.entriesTypesTabs.CurrentTab().Text)
	data.DeleteFile(deletionTestDbPath)
}
