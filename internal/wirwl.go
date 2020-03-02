package wirwl

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"sort"
	"wirwl/internal/data"
)

type App struct {
	window              fyne.Window
	entriesTabContainer *widget.TabContainer
	currentTab          string
	currentEntryNr      int
	entries             map[string][]data.Entry
	entriesLabels       map[string][]widget.Label
	dataProvider        *data.DataProvider
}

func NewApp(entriesPath string) *App {
	return &App{dataProvider: data.NewDataProvider(entriesPath)}
}

func (app *App) LoadAndDisplay(fyneApp fyne.App) {
	fyneApp.Settings().SetTheme(theme.LightTheme())
	app.window = fyneApp.NewWindow("wirwl")
	app.loadEntries()
	app.loadEntriesTabContainer()
	app.window.SetContent(widget.NewVBox(app.entriesTabContainer))
	app.window.ShowAndRun()
}

func (app *App) loadEntriesTabContainer() {
	tabs := app.loadEntriesTypesTabsWithTheirContent()
	if len(tabs) != 0 {
		app.entriesTabContainer = widget.NewTabContainer(tabs...)
	}
}

func (app *App) loadEntries() {
	app.entries = make(map[string][]data.Entry)
	entriesTypes, err := app.dataProvider.LoadEntriesTypesFromDb()
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(entriesTypes)
	for _, entryType := range entriesTypes {
		entries, err := app.dataProvider.LoadEntriesFromDb(entryType + "s")
		if err != nil {
			log.Fatal(err)
		}
		app.entries[entryType] = entries
	}
}

func (app *App) loadEntriesTypesTabsWithTheirContent() []*widget.TabItem {
	var tabs []*widget.TabItem
	app.entriesLabels = make(map[string][]widget.Label, len(app.entries))
	for entryType, entriesOfCertainType := range app.entries {
		labels := app.getEntriesNamesAsLabels(entriesOfCertainType)
		app.entriesLabels[entryType+"s"] = labels
		labelsAsCanvasObjects := app.getLabelsAsCanvasObjects(labels)
		tab := widget.NewTabItem(entryType+"s", widget.NewVBox(labelsAsCanvasObjects...))
		tabs = append(tabs, tab)
	}
	return tabs
}

func (app *App) getLabelsAsCanvasObjects(labels []widget.Label) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(labels))
	for i, _ := range labels {
		objects[i] = &labels[i]
	}
	return objects
}

func (app *App) getEntriesNamesAsLabels(entries []data.Entry) []widget.Label {
	var labels []widget.Label
	for _, entry := range entries {
		label := widget.NewLabel(entry.Title)
		labels = append(labels, *label)
	}
	return labels
}
