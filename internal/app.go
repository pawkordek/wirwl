package wirwl

import (
	"errors"
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	fyneWidget "fyne.io/fyne/widget"
	"log"
	"sort"
	"wirwl/internal/data"
	"wirwl/internal/widget"
)

type App struct {
	fyneApp             fyne.App
	mainWindow          fyne.Window
	addEntryTypeDialog  *widget.FormDialog
	msgDialog           *widget.MsgDialog
	confirmationDialog  *widget.ConfirmationDialog
	entriesTypesTabs    *widget.TabContainer
	currentEntryNr      int
	entries             map[string][]data.Entry
	entriesTypes        map[string]data.EntryType
	entriesLabels       map[string][]fyneWidget.Label
	dataProvider        data.Provider
	lastKeyPress        fyne.KeyName
	editEntryTypeDialog *widget.FormDialog
}

func NewApp(entriesPath string) *App {
	provider := data.NewBoltProvider(entriesPath)
	return &App{dataProvider: provider}
}

func (app *App) LoadAndDisplay(fyneApp fyne.App) {
	app.fyneApp = fyneApp
	app.fyneApp.Settings().SetTheme(theme.LightTheme())
	app.prepare()
	app.mainWindow.ShowAndRun()
}

func (app *App) prepare() {
	app.mainWindow = app.fyneApp.NewWindow("wirwl")
	app.loadEntriesTypes()
	app.loadEntries()
	app.loadEntriesTypesTabs()
	app.resetSelectedEntry()
	app.prepareDialogs()
	app.prepareMainWindowContent()
	app.mainWindow.Canvas().SetOnTypedKey(app.onKeyPressed)
}

func (app *App) onEnterPressedInAddEntryTypeDialog() {
	currentTabText := app.getCurrentTabText()
	err := app.addNewEntryType()
	if err != nil {
		app.msgDialog.Display(widget.ErrorPopUp, err.Error())
	} else {
		for _, tab := range app.entriesTypesTabs.Items() {
			if tab.Text == currentTabText {
				app.entriesTypesTabs.SelectTab(tab)
				app.updateCurrentlySelectedEntry()
				break
			}
		}
	}
}

func (app *App) prepareMainWindowContent() {
	app.mainWindow.SetContent(fyneWidget.NewVBox(app.entriesTypesTabs))
}

func (app *App) prepareDialogs() {
	app.msgDialog = widget.NewMsgPopUp(app.mainWindow.Canvas())
	app.confirmationDialog = widget.NewConfirmationDialog(app.mainWindow.Canvas())
	app.confirmationDialog.OnConfirm = app.deleteCurrentEntryType
	app.addEntryTypeDialog = widget.NewFormDialog(app.mainWindow.Canvas(), "Add new entry type", "Name", "Image query")
	app.addEntryTypeDialog.OnEnterPressed = app.onEnterPressedInAddEntryTypeDialog
	app.editEntryTypeDialog = widget.NewFormDialog(app.mainWindow.Canvas(), "Editing entry type: "+app.getCurrentTabText(), "Name", "Image query")
	app.editEntryTypeDialog.OnEnterPressed = app.applyChangesToCurrentEntryType
}

func (app *App) deleteCurrentEntryType() {
	currentTab := app.entriesTypesTabs.CurrentTab()
	delete(app.entries, currentTab.Text)
	delete(app.entriesTypes, currentTab.Text)
	delete(app.entriesLabels, currentTab.Text)
	app.entriesTypesTabs.Remove(currentTab)
}

func (app *App) applyChangesToCurrentEntryType() {
	currentEntryType := app.entriesTypes[app.getCurrentTabText()]
	oldTypeName := currentEntryType.Name
	currentEntryType.Name, _ = app.editEntryTypeDialog.GetItemValue("Name")
	currentEntryType.ImageQuery, _ = app.editEntryTypeDialog.GetItemValue("Image query")
	newTypeName := currentEntryType.Name
	app.entries[newTypeName] = app.entries[oldTypeName]
	app.entriesTypes[newTypeName] = currentEntryType
	delete(app.entriesTypes, oldTypeName)
	delete(app.entries, oldTypeName)
	currentTabIndex := app.entriesTypesTabs.CurrentTabIndex()
	app.loadEntriesTypesTabs()
	app.prepareMainWindowContent()
	app.entriesTypesTabs.SelectTabIndex(currentTabIndex)
	app.updateCurrentlySelectedEntry()
}

func (app *App) loadEntriesTypesTabs() {
	entriesGroupedByType := app.getEntriesNamesGroupedByType()
	app.entriesTypesTabs = widget.NewLabelsInTabContainer(entriesGroupedByType)
}

func (app *App) loadEntriesTypes() {
	app.entriesTypes = make(map[string]data.EntryType)
	entriesTypes, err := app.dataProvider.LoadEntriesTypesFromDb()
	if err != nil {
		log.Fatal(err)
	}
	for _, entryType := range entriesTypes {
		app.entriesTypes[entryType.Name] = entryType
	}
}

func (app *App) loadEntries() {
	app.entries = make(map[string][]data.Entry)
	for typeName, _ := range app.entriesTypes {
		entries, err := app.dataProvider.LoadEntriesFromDb(typeName)
		if err != nil {
			log.Fatal(err)
		}
		app.entries[typeName] = entries
	}
}

func (app *App) getEntriesNamesGroupedByType() map[string][]string {
	if len(app.entries) != 0 {
		return app.getEntriesNamesGroupedByTypeMap()
	} else {
		return app.getNoEntriesTab()
	}
}

func (app *App) getEntriesNamesGroupedByTypeMap() map[string][]string {
	entriesNames := make(map[string][]string)
	for entryType, entries := range app.entries {
		names := app.getEntriesNames(entries)
		entriesNames[entryType] = names
	}
	return entriesNames
}

func (app *App) getEntriesNames(entries []data.Entry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Title)
	}
	return names
}

func (app *App) getNoEntriesTab() map[string][]string {
	return map[string][]string{
		"No entries": {""},
	}
}

func (app *App) getAlphabeticallySortedEntriesTypesNames() []string {
	sortedTypesNames := make([]string, 0, len(app.entries))
	for typeName, _ := range app.entries {
		sortedTypesNames = append(sortedTypesNames, typeName)
	}
	sort.Strings(sortedTypesNames)
	return sortedTypesNames
}

func (app *App) loadEntryTabWithItsContent(typeName string) *fyneWidget.TabItem {
	labels := app.getEntriesNamesAsLabels(app.entries[typeName])
	app.entriesLabels[typeName] = labels
	labelsAsCanvasObjects := app.getLabelsAsCanvasObjects(labels)
	return fyneWidget.NewTabItem(typeName, fyneWidget.NewVBox(labelsAsCanvasObjects...))
}

func (app *App) getEntriesNamesAsLabels(entries []data.Entry) []fyneWidget.Label {
	var labels []fyneWidget.Label
	for _, entry := range entries {
		label := fyneWidget.NewLabel(entry.Title)
		labels = append(labels, *label)
	}
	return labels
}

func (app *App) getLabelsAsCanvasObjects(labels []fyneWidget.Label) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(labels))
	for i, _ := range labels {
		objects[i] = &labels[i]
	}
	return objects
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
	return app.entriesTypesTabs.CurrentTab().Text
}

func (app *App) onKeyPressed(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyL {
		app.selectNextTab()
	} else if event.Name == fyne.KeyH {
		app.selectPreviousTab()
	} else if app.lastKeyPress == fyne.KeyT {
		app.handleTabRelatedKeyPress(event)
	} else if event.Name == fyne.KeyS {
		err := app.saveChangesToDb()
		if err != nil {
			app.msgDialog.Display(widget.ErrorPopUp, err.Error())
		} else {
			app.msgDialog.Display(widget.SuccessPopUp, "Changes saved.")
		}
	}

	app.lastKeyPress = event.Name
}

func (app *App) handleTabRelatedKeyPress(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyI {
		app.displayDialogForAddingNewEntryType()
	} else if event.Name == fyne.KeyD {
		if len(app.entriesTypesTabs.Items()) > 1 {
			app.confirmationDialog.Display("Are you sure you want to delete entry type '" + app.entriesTypesTabs.CurrentTab().Text + "'?")
		} else {
			app.msgDialog.Display(widget.WarningPopUp, "You cannot remove the only remaining entry type!")
		}
	} else if event.Name == fyne.KeyE {
		app.editCurrentEntryType()
	}
}

func (app *App) displayDialogForAddingNewEntryType() {
	app.addEntryTypeDialog.CleanItemValues()
	app.addEntryTypeDialog.Display()
}

func (app *App) editCurrentEntryType() {
	app.editEntryTypeDialog.SetItemValue("Name", app.getCurrentTabText())
	app.editEntryTypeDialog.SetItemValue("Image query", app.entriesTypes[app.getCurrentTabText()].ImageQuery)
	app.editEntryTypeDialog.Display()
}

func (app *App) selectNextTab() {
	if app.entriesTypesTabs.CurrentTabIndex() < len(app.entries)-1 {
		app.changeTab(1)
	} else {
		app.selectTab(0)
	}
}

func (app *App) selectPreviousTab() {
	if app.entriesTypesTabs.CurrentTabIndex() > 0 {
		app.changeTab(-1)
	} else {
		app.selectTab(len(app.entriesTypesTabs.Items()) - 1)
	}
}

func (app *App) selectTab(tabNum int) {
	app.entriesTypesTabs.SelectTabIndex(tabNum)
	app.updateCurrentlySelectedEntry()
}

func (app *App) changeTab(byHowManyTabs int) {
	currentIndex := app.entriesTypesTabs.CurrentTabIndex()
	newIndex := currentIndex + byHowManyTabs
	app.selectTab(newIndex)
}

func (app *App) addNewEntryType() error {
	newEntryTypeName, _ := app.addEntryTypeDialog.GetItemValue("Name")
	if _, exists := app.entriesTypes[newEntryTypeName]; !exists {
		app.entries[newEntryTypeName] = nil
		imageQuery, _ := app.addEntryTypeDialog.GetItemValue("Image query")
		app.entriesTypes[newEntryTypeName] = data.EntryType{
			Name:       newEntryTypeName,
			ImageQuery: imageQuery,
		}
		app.loadEntriesTypesTabs()
		app.prepareMainWindowContent()
		return nil
	} else {
		return errors.New("Entry type with name '" + newEntryTypeName + "' already exists.")
	}
}

func (app *App) saveChangesToDb() error {
	var types []data.EntryType
	for _, entryType := range app.entriesTypes {
		types = append(types, entryType)
	}
	err := app.dataProvider.SaveEntriesTypesToDb(types)
	for entryType, entry := range app.entries {
		err = app.dataProvider.SaveEntriesToDb(entryType, entry)
	}
	return err
}
