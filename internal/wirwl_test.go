package wirwl

import (
	"fyne.io/fyne"
	fyneTest "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"wirwl/internal/data"
)

const exampleDbPath = "../test/exampleDb.db"
const emptyDbPath = "../test/emptyDb.db"

func TestMain(m *testing.M) {
	dataProvider := data.NewDataProvider(exampleDbPath)
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

func TestThatEntriesTabsWithContentDisplay(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	for _, tab := range app.entriesTabContainer.Items {
		if tab.Text == "comics" {
			assert.Equal(t, "some comic1", app.entriesLabels["comics"][0].Text)
			assert.Equal(t, "some comic2", app.entriesLabels["comics"][1].Text)
		} else if tab.Text == "videos" {
			assert.Equal(t, "some video1", app.entriesLabels["videos"][0].Text)
			assert.Equal(t, "some video2", app.entriesLabels["videos"][1].Text)
		} else if tab.Text == "music" {
			assert.Equal(t, "some music1", app.entriesLabels["music"][0].Text)
			assert.Equal(t, "some music2", app.entriesLabels["music"][1].Text)
		} else {
			assert.Fail(t, "There is an unexpected tab called "+tab.Text+" displayed!")
		}
	}
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

func TestWhetherEntryTypeInputOpens(t *testing.T) {
	app := NewApp(emptyDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, true, app.addEntryTypePopUp.Hidden)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	assert.Equal(t, true, app.addEntryTypePopUp.Visible())
	assert.Equal(t, true, app.typeInput.Visible())
	assert.Equal(t, true, app.typeInput.Focused())
}

func TestAddingOfNewEntryType(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.typeInput.Type("new entry type")
	app.typeInput.SimulateKeyPress(fyne.KeyEnter)
	assert.Equal(t, 4, len(app.entriesTabContainer.Items))
	assert.Equal(t, "comics", app.entriesTabContainer.Items[0].Text)
	assert.Equal(t, "music", app.entriesTabContainer.Items[1].Text)
	assert.Equal(t, "new entry type", app.entriesTabContainer.Items[2].Text)
	assert.Equal(t, "videos", app.entriesTabContainer.Items[3].Text)
	assert.Equal(t, true, app.addEntryTypePopUp.Hidden)
	assert.Equal(t, true, !app.typeInput.Focused())
}

func TestThatItIsNotPossibleToAddTheSameEntryTypeTwice(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.typeInput.Type("type")
	app.typeInput.SimulateKeyPress(fyne.KeyEnter)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.typeInput.Type("type")
	app.typeInput.SimulateKeyPress(fyne.KeyEnter)
	assert.Equal(t, true, app.errorPopUp.Visible())
	assert.Equal(t, "Entry type with name 'type' already exists.", app.errorMsg.Text)
	assert.Equal(t, 4, len(app.entriesTabContainer.Items))
}

func TestThatPressingAnyKeyClosesErrorPopUp(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.errorPopUp.Show()
	app.SimulateKeyPress(fyne.KeyT)
	assert.Equal(t, true, app.errorPopUp.Hidden)
	app.errorPopUp.Show()
	app.SimulateKeyPress(fyne.KeyReturn)
	assert.Equal(t, true, app.errorPopUp.Hidden)
}

func TestThatAddingNewEntryTypeDoesNotChangeCurrentlyOpenedTab(t *testing.T) {
	app := NewApp(exampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	app.SimulateKeyPress(fyne.KeyL)
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
	app.typeInput.Type("type")
	app.typeInput.SimulateKeyPress(fyne.KeyEnter)
	assert.Equal(t, "music", app.getCurrentTabText())
}

func (app *App) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	onTypedKey := app.mainWindow.Canvas().OnTypedKey()
	onTypedKey(event)
}
