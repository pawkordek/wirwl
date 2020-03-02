package wirwl

import (
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
