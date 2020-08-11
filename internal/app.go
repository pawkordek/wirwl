package wirwl

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	fyneWidget "fyne.io/fyne/widget"
	"github.com/pkg/errors"
	"wirwl/internal/data"
	"wirwl/internal/log"
	"wirwl/internal/widget"
)

type App struct {
	fyneApp             fyne.App
	mainWindow          fyne.Window
	config              Config
	loadingErrors       map[string]string
	addEntryTypeDialog  *widget.FormDialog
	msgDialog           *widget.MsgDialog
	confirmationDialog  *widget.ConfirmationDialog
	entriesTypesTabs    *widget.TabContainer
	entries             map[data.EntryType][]data.Entry
	dataProvider        data.Provider
	lastKeyPress        fyne.KeyName
	editEntryTypeDialog *widget.FormDialog
	inputHandler        InputHandler
}

const configLoadError = "CONFIG_LOAD_ERROR"
const entriesLoadError = "ENTRIES_LOAD_ERROR"

func NewApp(fyneApp fyne.App, config Config, dataProvider data.Provider, loadingErrors map[string]string) *App {
	return &App{fyneApp: fyneApp, config: config, dataProvider: dataProvider, loadingErrors: loadingErrors}
}

func (app *App) LoadAndDisplay() error {
	app.prepare()
	app.displayLoadingErrors()
	app.mainWindow.ShowAndRun()
	app.shutdown()
	return nil
}

func (app *App) prepare() {
	app.setupBasicSettings()
	app.setupInputHandler()
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

func (app *App) setupInputHandler() {
	app.inputHandler = NewInputHandler(app.config.Keymap)
	app.inputHandler.bindFunctionToAction(selectNextTabAction, func() { app.entriesTypesTabs.SelectNextTab() })
	app.inputHandler.bindFunctionToAction(selectPreviousTabAction, func() { app.entriesTypesTabs.SelectPreviousTab() })
}

func (app *App) loadEntries() {
	entries, err := app.dataProvider.LoadEntries()
	if err != nil {
		msg := "Failed to load entries. Application will now exit as it cannot continue."
		err = errors.Wrap(err, msg)
		log.Error(err)
		app.loadingErrors[entriesLoadError] = msg
	}
	app.entries = entries
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
		entriesNames[entryType.Name] = names
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

func (app *App) displayLoadingErrors() {
	if len(app.loadingErrors) != 0 {
		if app.loadingErrors[entriesLoadError] != "" {
			app.msgDialog.SetOneTimeOnHideCallback(func() {
				app.fyneApp.Quit()
			})
		}
		errorsList := ""
		for _, err := range app.loadingErrors {
			errorsList += fmt.Sprintln("- " + err)
		}
		app.msgDialog.Display(widget.ErrorPopUp, errorsList)
	}
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
	currentEntryType := app.getCurrentEntryType()
	delete(app.entries, currentEntryType)
	app.loadEntriesTypesTabs()
	app.prepareMainWindowContent()
}

func (app *App) getCurrentEntryType() data.EntryType {
	for entryType := range app.entries {
		if entryType.Name == app.getCurrentTabText() {
			return entryType
		}
	}
	return data.EntryType{}
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
	currentEntryType := app.getCurrentEntryType()
	oldType := currentEntryType
	currentEntryType.Name = app.editEntryTypeDialog.ItemValue("Name")
	currentEntryType.ImageQuery = app.editEntryTypeDialog.ItemValue("Image query")
	app.entries[currentEntryType] = app.entries[oldType]
	delete(app.entries, oldType)
	currentTabIndex := app.entriesTypesTabs.CurrentTabIndex()
	app.loadEntriesTypesTabs()
	app.prepareMainWindowContent()
	app.entriesTypesTabs.SelectTabIndex(currentTabIndex)
}

func (app *App) getCurrentTabText() string {
	return app.entriesTypesTabs.CurrentTab().Text
}

func (app *App) onKeyPressed(event *fyne.KeyEvent) {
	app.inputHandler.handle(event.Name)
	if app.lastKeyPress == fyne.KeyT {
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
	currentEntryType := app.getCurrentEntryType()
	app.editEntryTypeDialog.SetItemValue("Name", currentEntryType.Name)
	app.editEntryTypeDialog.SetItemValue("Image query", currentEntryType.ImageQuery)
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
	if app.noEntryTypeWithNameExists(newEntryTypeName) {
		newEntryType := data.EntryType{
			Name:       newEntryTypeName,
			ImageQuery: app.addEntryTypeDialog.ItemValue("Image query"),
		}
		app.entries[newEntryType] = []data.Entry{}
		app.loadEntriesTypesTabs()
		app.prepareMainWindowContent()
		return nil
	} else {
		return errors.New("Entry type with name '" + newEntryTypeName + "' already exists.")
	}
}

func (app *App) noEntryTypeWithNameExists(name string) bool {
	for entryType := range app.entries {
		if entryType.Name == name {
			return false
		}
	}
	return true
}

func (app *App) saveChangesToDb() error {
	err := app.dataProvider.SaveEntries(app.entries)
	return err
}

func (app *App) shutdown() {
	if _, configLoadingErrorExists := app.loadingErrors[configLoadError]; !configLoadingErrorExists {
		err := app.config.save()
		if err != nil {
			log.Error(err)
		}
	}
}
