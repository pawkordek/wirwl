package wirwl

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	fyneWidget "fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"
	"wirwl/internal/data"
	"wirwl/internal/input"
	"wirwl/internal/log"
	"wirwl/internal/widget"
)

type App struct {
	fyneApp                  fyne.App
	mainWindow               fyne.Window
	config                   Config
	loadingErrors            map[string]string
	addEntryTypeDialog       *widget.FormDialog
	msgDialog                *widget.MsgDialog
	confirmationDialog       *widget.ConfirmationDialog
	entriesTypesTabs         *widget.TabContainer
	recentlyPressedKeysLabel *fyneWidget.Label
	entriesContainer         *data.EntriesContainer
	editEntryTypeDialog      *widget.FormDialog
	inputHandler             input.Handler
	entriesTables            map[data.EntryType]*widget.Table
}

const configLoadError = "CONFIG_LOAD_ERROR"
const entriesLoadError = "ENTRIES_LOAD_ERROR"

func NewApp(fyneApp fyne.App, config Config, dataProvider data.Provider, loadingErrors map[string]string) *App {
	return &App{
		fyneApp:          fyneApp,
		config:           config,
		entriesContainer: data.NewEntriesContainer(dataProvider),
		loadingErrors:    loadingErrors,
		entriesTables:    map[data.EntryType]*widget.Table{}}
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
	app.inputHandler = input.NewHandler(app.config.Keymap)
	app.inputHandler.SetOnKeyPressedCallbackFunction(func(keyCombination input.KeyCombination) {
		app.recentlyPressedKeysLabel.SetText("Recently pressed keys: " + keyCombination.String())
	})
	app.inputHandler.BindFunctionToAction(appName, input.SelectNextTabAction, func() { app.entriesTypesTabs.SelectNextTab() })
	app.inputHandler.BindFunctionToAction(appName, input.SelectPreviousTabAction, func() { app.entriesTypesTabs.SelectPreviousTab() })
	app.inputHandler.BindFunctionToAction(appName, input.SaveChangesAction, func() { app.trySavingChangesToDb() })
	app.inputHandler.BindFunctionToAction(appName, input.AddEntryTypeAction, func() { app.displayDialogForAddingNewEntryType() })
	app.inputHandler.BindFunctionToAction(appName, input.EditCurrentEntryTypeAction, func() { app.editCurrentEntryType() })
	app.inputHandler.BindFunctionToAction(appName, input.RemoveEntryTypeAction, func() { app.tryDeletingCurrentEntryType() })
	app.inputHandler.BindFunctionToAction(appName, input.EnterInputModeAction, func() { app.getCurrentEntryTypeTable().EnterInputMode() })
}

func (app *App) loadEntries() {
	err := app.entriesContainer.LoadData()
	if err != nil {
		msg := "Failed to load entries. Application will now exit as it cannot continue."
		err = errors.Wrap(err, msg)
		log.Error(err)
		app.loadingErrors[entriesLoadError] = msg
	}
	if app.entriesContainer.AmountOfTypes() == 0 {
		noEntriesType := data.EntryType{
			Name:                  "No entries",
			CompletionElementName: "",
			ImageQuery:            "",
		}
		err = app.entriesContainer.AddEntryType(noEntriesType)
		if err != nil {
			log.Error(err)
		}
	}
	app.entriesContainer.SubscribeToChanges(app.reloadGUI)
}

func (app *App) loadEntriesTypesTabs() {
	app.entriesTypesTabs = widget.NewTabContainer(
		app.createTabsWithEntriesTableForEachEntryType(),
		func(element *fyne.CanvasObject) {},
		func(element *fyne.CanvasObject) {})
}

func (app *App) createTabsWithEntriesTableForEachEntryType() map[string][]fyne.CanvasObject {
	entriesGroupedByType := app.entriesContainer.EntriesGroupedByType()
	tabsData := make(map[string][]fyne.CanvasObject, len(entriesGroupedByType))
	for entryType, entries := range entriesGroupedByType {
		app.createEntriesTable(entryType, entries)
		tabElements := make([]fyne.CanvasObject, 0, 1)
		tabElements = append(tabElements, app.entriesTables[entryType])
		tabsData[entryType.Name] = tabElements
	}
	return tabsData
}

func (app *App) getEntriesNamesGroupedByType() map[string][]string {
	entriesGroupedByType := app.entriesContainer.EntriesGroupedByType()
	if len(entriesGroupedByType) != 0 {
		return getEntriesGroupedByTypeAsStrings(entriesGroupedByType)
	} else {
		return getNoEntriesTab()
	}
}

func getEntriesGroupedByTypeAsStrings(entriesGroupedByType map[data.EntryType][]data.Entry) map[string][]string {
	entriesGroupedByTypeAsStrings := make(map[string][]string)
	for entryType, entries := range entriesGroupedByType {
		names := getEntriesNamesFrom(entries)
		entriesGroupedByTypeAsStrings[entryType.Name] = names
	}
	return entriesGroupedByTypeAsStrings
}

func getEntriesNamesFrom(entries []data.Entry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Title)
	}
	return names
}

func getNoEntriesTab() map[string][]string {
	return map[string][]string{
		"No entries": {""},
	}
}

func (app *App) prepareMainWindowContent() {
	app.recentlyPressedKeysLabel = fyneWidget.NewLabel("Recently pressed keys: ")
	content := container.NewBorder(app.entriesTypesTabs, app.recentlyPressedKeysLabel, nil, nil)
	app.mainWindow.SetContent(content)
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
	app.createAddEntryTypeDialog()
	app.editEntryTypeDialog = widget.NewFormDialog(app.mainWindow.Canvas(), app.inputHandler, "Editing entry type: "+app.getCurrentTabText(), app.createEntryTypeRelatedDialogElements()...)
	app.editEntryTypeDialog.OnEnterPressed = app.applyChangesToCurrentEntryType
}

func (app *App) reloadGUI() {
	app.loadEntriesTypesTabs()
	app.prepareMainWindowContent()
}

func (app *App) createEntryTypeRelatedDialogElements() []*widget.FormDialogFormItem {
	formItemFactory := widget.NewFormDialogFormItemFactory(app.mainWindow.Canvas(), app.inputHandler)
	entryTypeRelatedDialogElements := []*widget.FormDialogFormItem{}
	entryTypeRelatedDialogElements = append(entryTypeRelatedDialogElements, formItemFactory.FormItemWithInputField("Name"))
	entryTypeRelatedDialogElements = append(entryTypeRelatedDialogElements, formItemFactory.FormItemWithInputField("Image query"))
	return entryTypeRelatedDialogElements
}

func (app *App) deleteCurrentEntryType() {
	nameOfTypeToDelete := app.getCurrentTabText()
	err := app.entriesContainer.DeleteEntryType(nameOfTypeToDelete)
	if err != nil {
		err = errors.Wrap(err, "There was an error when deleting an entry type. This is most likely a programming error")
		log.Error(err)
	}
}

func (app *App) getCurrentEntryType() data.EntryType {
	currentEntryTypeName := app.getCurrentTabText()
	currentEntryType, err := app.entriesContainer.EntryTypeWithName(currentEntryTypeName)
	if err != nil {
		err = errors.Wrap(err, "An error occurred when trying to get current entry type. This is most likely a programming error")
		log.Error(err)
	}
	return currentEntryType
}

func (app *App) getCurrentTabText() string {
	currentTab := app.entriesTypesTabs.CurrentTab()
	if currentTab != nil {
		return currentTab.Text
	}
	return ""
}

func (app *App) getCurrentEntryTypeTable() *widget.Table {
	return app.entriesTables[app.getCurrentEntryType()]
}

func (app *App) onKeyPressed(event *fyne.KeyEvent) {
	app.inputHandler.HandleInNormalMode(appName, event.Name)
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

func (app *App) tryDeletingCurrentEntryType() {
	if len(app.entriesTypesTabs.Items()) > 1 {
		app.confirmationDialog.Display("Are you sure you want to delete entry type '" + app.entriesTypesTabs.CurrentTab().Text + "'?")
	} else {
		app.msgDialog.Display(widget.WarningPopUp, "You cannot remove the only remaining entry type!")
	}
}

func (app *App) trySavingChangesToDb() {
	err := app.entriesContainer.SaveData()
	if err != nil {
		app.msgDialog.Display(widget.ErrorPopUp, err.Error())
	} else {
		app.msgDialog.Display(widget.SuccessPopUp, "Changes saved.")
	}
}

func (app *App) shutdown() {
	if _, configLoadingErrorExists := app.loadingErrors[configLoadError]; !configLoadingErrorExists {
		err := app.config.save()
		if err != nil {
			log.Error(err)
		}
	}
}
