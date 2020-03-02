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

func TestMain(m *testing.M) {
	dataProvider := data.NewDataProvider(data.ExampleDbPath)
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
	os.Exit(m.Run())
}

func TestThatEntriesTabsWithContentDisplay(t *testing.T) {
	app := NewApp(data.ExampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	for _, tab := range app.entriesTabContainer.Items {
		if tab.Text == "comics" {
			assert.Equal(t, "some comic1", app.entriesLabels["comics"][0].Text)
			assert.Equal(t, "some comic2", app.entriesLabels["comics"][1].Text)
		} else if tab.Text == "videos" {
			assert.Equal(t, "some video1", app.entriesLabels["videos"][0].Text)
			assert.Equal(t, "some video2", app.entriesLabels["videos"][1].Text)
		} else {
			assert.Fail(t, "There is an unexpected tab called "+tab.Text+" displayed!")
		}
	}
}

func TestSwitchingTabs(t *testing.T) {
	app := NewApp(data.ExampleDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, fyne.TextStyle{Bold: true}, app.entriesLabels["comics"][0].TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, app.entriesLabels["comics"][1].TextStyle)
	app.SimulateKeyPress(fyne.KeyL)
	assert.Equal(t, "videos", app.currentTab)
	assert.Equal(t, fyne.TextStyle{Bold: true}, app.entriesLabels["videos"][0].TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, app.entriesLabels["videos"][1].TextStyle)
	app.SimulateKeyPress(fyne.KeyH)
	assert.Equal(t, "comics", app.currentTab)
}

func TestThatIfThereAreNoEntriesCorrectMessageDisplays(t *testing.T) {
	app := NewApp(data.EmptyDbPath)
	app.LoadAndDisplay(fyneTest.NewApp())
	assert.Equal(t, 1, len(app.entriesTabContainer.Items))
	assert.Equal(t, "No entries", app.entriesTabContainer.Items[0].Text)
}

func (app *App) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	onTypedKey := app.window.Canvas().OnTypedKey()
	onTypedKey(event)
}
