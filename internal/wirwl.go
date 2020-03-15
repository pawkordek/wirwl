package wirwl

import (
	"errors"
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"sort"
	"wirwl/internal/data"
)

type App struct {
	fyneApp             fyne.App
	mainWindow          fyne.Window
	addEntryTypePopUp   *widget.PopUp
	errorPopUp          *widget.PopUp
	errorMsg            *widget.Label
	entriesTabContainer *widget.TabContainer
	currentEntryNr      int
	entries             map[string][]data.Entry
	entriesLabels       map[string][]widget.Label
	dataProvider        *data.DataProvider
	lastKeyPress        fyne.KeyName
	typeInput           *Input
}

func NewApp(entriesPath string) *App {
	return &App{dataProvider: data.NewDataProvider(entriesPath)}
}

func (app *App) LoadAndDisplay(fyneApp fyne.App) {
	app.fyneApp = fyneApp
	app.fyneApp.Settings().SetTheme(theme.LightTheme())
	app.prepareMainWindow()
	app.mainWindow.ShowAndRun()
}

func (app *App) prepareMainWindow() {
	app.mainWindow = app.fyneApp.NewWindow("wirwl")
	app.loadEntries()
	app.loadEntriesTabContainer()
	app.resetSelectedEntry()
	app.prepareAddEntryTypePopUp()
	app.prepareErrorPopUp()
	app.prepareMainWindowContent()
	app.mainWindow.Canvas().SetOnTypedKey(app.onKeyPressed)
}

func (app *App) prepareAddEntryTypePopUp() {
	app.typeInput = newInput()
	app.typeInput.SetOnEnterPressed(app.onTypeInputEnterPressed)
	popUpTitle := widget.NewLabel("Add new entry type")
	popUpTitle.Alignment = fyne.TextAlignCenter
	popUpContent := widget.NewVBox(popUpTitle, app.typeInput)
	app.addEntryTypePopUp = widget.NewModalPopUp(popUpContent, app.mainWindow.Canvas())
	app.addEntryTypePopUp.Hide()
}

func (app *App) onTypeInputEnterPressed() {
	currentTabText := app.getCurrentTabText()
	err := app.addNewTab(app.typeInput.Text)
	app.mainWindow.Canvas().Unfocus()
	app.addEntryTypePopUp.Hide()
	app.typeInput.Text = ""
	if err != nil {
		app.errorMsg.SetText(err.Error())
		app.errorPopUp.Show()
	} else {
		for _, tab := range app.entriesTabContainer.Items {
			if tab.Text == currentTabText {
				app.entriesTabContainer.SelectTab(tab)
				app.updateCurrentlySelectedEntry()
				break
			}
		}
	}
}

func (app *App) prepareMainWindowContent() {
	app.mainWindow.SetContent(widget.NewVBox(app.entriesTabContainer))
}

func (app *App) prepareErrorPopUp() {
	app.errorMsg = widget.NewLabel("")
	title := widget.NewLabel("ERROR!")
	title.Alignment = fyne.TextAlignCenter
	content := widget.NewVBox(title, app.errorMsg)
	app.errorPopUp = widget.NewModalPopUp(content, app.mainWindow.Canvas())
	app.errorPopUp.Hide()
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
		entries, err := app.dataProvider.LoadEntriesFromDb(entryType)
		if err != nil {
			log.Fatal(err)
		}
		app.entries[entryType] = entries
	}
}

func (app *App) loadEntriesTypesTabsWithTheirContent() []*widget.TabItem {
	var tabs []*widget.TabItem
	if len(app.entries) != 0 {
		app.entriesLabels = make(map[string][]widget.Label, len(app.entries))
		orderedEntriesKeys := app.getOrderedEntriesKeys()
		for _, entryType := range orderedEntriesKeys {
			labels := app.getEntriesNamesAsLabels(app.entries[entryType])
			app.entriesLabels[entryType] = labels
			labelsAsCanvasObjects := app.getLabelsAsCanvasObjects(labels)
			tab := widget.NewTabItem(entryType, widget.NewVBox(labelsAsCanvasObjects...))
			tabs = append(tabs, tab)
		}
		return tabs
	} else {
		tab := widget.NewTabItem("No entries", widget.NewVBox())
		return append(tabs, tab)
	}
}

func (app *App) getOrderedEntriesKeys() []string {
	orderedEntriesKeys := make([]string, 0, len(app.entries))
	for key, _ := range app.entries {
		orderedEntriesKeys = append(orderedEntriesKeys, key)
	}
	sort.Strings(orderedEntriesKeys)
	return orderedEntriesKeys
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

func (app *App) resetSelectedEntry() {
	app.currentEntryNr = 0
	app.updateCurrentlySelectedEntry()
}

func (app *App) updateCurrentlySelectedEntry() {
	for _, label := range app.entriesLabels[app.getCurrentTabText()] {
		(&label).TextStyle = fyne.TextStyle{
			Bold: false,
		}
		(&label).Refresh()
	}
	if len(app.entriesLabels[app.getCurrentTabText()]) > 0 {
		label := &app.entriesLabels[app.getCurrentTabText()][app.currentEntryNr]
		label.TextStyle = fyne.TextStyle{
			Bold: true,
		}
		label.Refresh()
	}
}

func (app *App) getCurrentTabText() string {
	return app.entriesTabContainer.CurrentTab().Text
}

func (app *App) onKeyPressed(event *fyne.KeyEvent) {
	if app.errorPopUp.Visible() {
		app.errorPopUp.Hide()
	}
	if event.Name == fyne.KeyL {
		app.selectNextTab()
	}
	if event.Name == fyne.KeyH {
		app.selectPreviousTab()
	}
	if event.Name == fyne.KeyI && app.lastKeyPress == fyne.KeyT {
		app.addEntryTypePopUp.Show()
		app.mainWindow.Canvas().Focus(app.typeInput)
	}
	app.lastKeyPress = event.Name
}

func (app *App) selectNextTab() {
	if app.entriesTabContainer.CurrentTabIndex() < len(app.entries)-1 {
		app.changeTab(1)
	} else {
		app.selectTab(0)
	}
}

func (app *App) selectPreviousTab() {
	if app.entriesTabContainer.CurrentTabIndex() > 0 {
		app.changeTab(-1)
	} else {
		app.selectTab(len(app.entriesTabContainer.Items) - 1)
	}
}

func (app *App) selectTab(tabNum int) {
	app.entriesTabContainer.SelectTabIndex(tabNum)
	app.updateCurrentlySelectedEntry()
}

func (app *App) changeTab(byHowManyTabs int) {
	currentIndex := app.entriesTabContainer.CurrentTabIndex()
	newIndex := currentIndex + byHowManyTabs
	app.selectTab(newIndex)
}

func (app *App) addNewTab(name string) error {
	if _, exists := app.entries[name]; !exists {
		app.entries[name] = nil
		app.loadEntriesTabContainer()
		app.prepareMainWindowContent()
		return nil
	} else {
		return errors.New("Entry type with name '" + name + "' already exists.")
	}
}
