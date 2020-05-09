package wirwl

import (
	"errors"
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	fyneWidget "fyne.io/fyne/widget"
	"log"
	"wirwl/internal/data"
	"wirwl/internal/widget"
)

type App struct {
	fyneApp             fyne.App
	mainWindow          fyne.Window
	config              Config
	addEntryTypeDialog  *widget.FormDialog
	msgDialog           *widget.MsgDialog
	confirmationDialog  *widget.ConfirmationDialog
	entriesTypesTabs    *widget.TabContainer
	entries             map[string][]data.Entry
	entriesTypes        map[string]data.EntryType
	dataProvider        data.Provider
	lastKeyPress        fyne.KeyName
	editEntryTypeDialog *widget.FormDialog
}

func NewApp(fyneApp fyne.App) *App {
	return &App{fyneApp: fyneApp}
}

func (app *App) LoadAndDisplay(configDirPath string, appDataDirPath string) {
	setupLoggingIn(appDataDirPath)
	app.dataProvider = loadDataProviderIn(appDataDirPath)
	app.config = loadConfigFromDir(configDirPath)
	app.prepare()
	app.mainWindow.ShowAndRun()
	app.shutdown()
}

func (app *App) prepare() {
	app.setupBasicSettings()
	app.loadEntriesTypes()
	app.loadEntries()
	app.loadEntriesTypesTabs()
	app.prepareDialogs()
	app.prepareMainWindowContent()
	app.mainWindow.Canvas().SetOnTypedKey(app.onKeyPressed)
}

func (app *App) setupBasicSettings() {
	app.mainWindow = app.fyneApp.NewWindow("wirwl")
	app.fyneApp.Settings().SetTheme(theme.LightTheme())
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
	for typeName := range app.entriesTypes {
		entries, err := app.dataProvider.LoadEntriesFromDb(typeName)
		if err != nil {
			log.Fatal(err)
		}
		app.entries[typeName] = entries
	}
}

func (app *App) loadEntriesTypesTabs() {
	entriesGroupedByType := app.getEntriesNamesGroupedByType()
	app.entriesTypesTabs = widget.NewLabelsInTabContainer(entriesGroupedByType)
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
	app.loadEntriesTypesTabs()
	app.prepareMainWindowContent()
}

func (app *App) onEnterPressedInAddEntryTypeDialog() {
	currentTabText := app.getCurrentTabText()
	err := app.addNewEntryType()
	if err != nil {
		app.msgDialog.SetOneTimeOnHideCallback(func() {
			app.addEntryTypeDialog.Display()
		})
		app.msgDialog.Display(widget.ErrorPopUp, err.Error())
	} else {
		for _, tab := range app.entriesTypesTabs.Items() {
			if tab.Text == currentTabText {
				app.entriesTypesTabs.SelectTab(tab)
				break
			}
		}
	}
}

func (app *App) applyChangesToCurrentEntryType() {
	currentEntryType := app.entriesTypes[app.getCurrentTabText()]
	oldTypeName := currentEntryType.Name
	currentEntryType.Name = app.editEntryTypeDialog.ItemValue("Name")
	currentEntryType.ImageQuery = app.editEntryTypeDialog.ItemValue("Image query")
	newTypeName := currentEntryType.Name
	app.entries[newTypeName] = app.entries[oldTypeName]
	app.entriesTypes[newTypeName] = currentEntryType
	delete(app.entriesTypes, oldTypeName)
	delete(app.entries, oldTypeName)
	currentTabIndex := app.entriesTypesTabs.CurrentTabIndex()
	app.loadEntriesTypesTabs()
	app.prepareMainWindowContent()
	app.entriesTypesTabs.SelectTabIndex(currentTabIndex)
}

func (app *App) getCurrentTabText() string {
	return app.entriesTypesTabs.CurrentTab().Text
}

func (app *App) onKeyPressed(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyL {
		app.entriesTypesTabs.SelectNextTab()
	} else if event.Name == fyne.KeyH {
		app.entriesTypesTabs.SelectPreviousTab()
	} else if app.lastKeyPress == fyne.KeyT {
		app.handleTabRelatedKeyPress(event)
	} else if app.lastKeyPress == fyne.KeyS && event.Name == fyne.KeyY {
		app.trySavingChangesToDb()
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

func (app *App) trySavingChangesToDb() {
	err := app.saveChangesToDb()
	if err != nil {
		app.msgDialog.Display(widget.ErrorPopUp, err.Error())
	} else {
		app.msgDialog.Display(widget.SuccessPopUp, "Changes saved.")
	}
}

func (app *App) addNewEntryType() error {
	newEntryTypeName := app.addEntryTypeDialog.ItemValue("Name")
	if newEntryTypeName == "" {
		return errors.New("You cannot add entry type with empty name")
	}
	if _, exists := app.entriesTypes[newEntryTypeName]; !exists {
		app.entries[newEntryTypeName] = nil
		imageQuery := app.addEntryTypeDialog.ItemValue("Image query")
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

func (app *App) shutdown() {
	app.config.save()
}
