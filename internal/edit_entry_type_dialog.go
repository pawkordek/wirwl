package wirwl

import (
	"wirwl/internal/data"
	"wirwl/internal/log"
	"wirwl/internal/widget"
)

func (app *App) createEditEntryTypeDialog() {
	app.editEntryTypeDialog = widget.NewFormDialog(app.mainWindow.Canvas(), app.inputHandler, "Editing entry type: "+app.getCurrentTabText(), app.createEntryTypeRelatedDialogElements()...)
	app.editEntryTypeDialog.OnEnterPressed = app.applyChangesToCurrentEntryType
}

func (app *App) applyChangesToCurrentEntryType() {
	currentTabIndex := app.entriesTypesTabs.CurrentTabIndex()
	nameOfEntryToUpdate := app.getCurrentTabText()
	entryToUpdateWith := app.getEntryToUpdateWith()
	err := app.entriesContainer.UpdateEntryType(nameOfEntryToUpdate, entryToUpdateWith)
	if err != nil {
		log.Error(err)
	}
	app.entriesTypesTabs.SelectTabIndex(currentTabIndex)
}

func (app *App) getEntryToUpdateWith() data.EntryType {
	return data.EntryType{
		Name:                  app.editEntryTypeDialog.ItemValue("Name"),
		CompletionElementName: "",
		ImageQuery:            app.editEntryTypeDialog.ItemValue("Image query"),
	}
}
