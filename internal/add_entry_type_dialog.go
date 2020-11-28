package wirwl

import (
	"wirwl/internal/data"
	"wirwl/internal/widget"
)

func (app *App) createAddEntryTypeDialog() {
	app.addEntryTypeDialog = widget.NewFormDialog(app.mainWindow.Canvas(), app.inputHandler, "Add new entry type", app.createEntryTypeRelatedDialogElements()...)
	app.addEntryTypeDialog.OnEnterPressed = app.onEnterPressedInAddEntryTypeDialog
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

func (app *App) addNewEntryType() error {
	newEntryType := app.getNewEntryType()
	err := app.entriesContainer.AddEntryType(newEntryType)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) getNewEntryType() data.EntryType {
	return data.EntryType{
		Name:       app.addEntryTypeDialog.ItemValue("Name"),
		ImageQuery: app.addEntryTypeDialog.ItemValue("Image query"),
	}
}