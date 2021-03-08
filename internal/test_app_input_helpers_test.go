package wirwl

import "fyne.io/fyne/v2"

func (app *App) simulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	focusedElement := app.mainWindow.Canvas().Focused()
	if focusedElement != nil {
		focusedElement.TypedKey(event)
	} else {
		onTypedKey := app.mainWindow.Canvas().OnTypedKey()
		onTypedKey(event)
	}
}

func (app *App) simulateSwitchingToNextEntryType() {
	app.simulateKeyPress(fyne.KeyL)
}

func (app *App) simulateSwitchingToPreviousEntryType() {
	app.simulateKeyPress(fyne.KeyH)
}

func (app *App) simulateOpeningDialogForAddingEntryType() {
	app.simulateKeyPress(fyne.KeyT)
	app.simulateKeyPress(fyne.KeyI)
}

func (app *App) simulateAddingNewEntryTypeWithName(text string) {
	app.simulateOpeningDialogForAddingEntryType()
	app.simulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type(text)
	app.simulateKeyPress(fyne.KeyReturn)
}

func (app *App) simulateSavingChanges() {
	app.simulateKeyPress(fyne.KeyS)
	app.simulateKeyPress(fyne.KeyY)
}

func (app *App) simulateAttemptAtDeletionOfCurrentEntryType() {
	app.simulateKeyPress(fyne.KeyT)
	app.simulateKeyPress(fyne.KeyD)
}

func (app *App) simulateDeletionOfCurrentEntryType() {
	app.simulateAttemptAtDeletionOfCurrentEntryType()
	app.simulateKeyPress(fyne.KeyY)
}

func (app *App) simulateEditionOfCurrentEntryTypeTo(text string) {
	app.simulateKeyPress(fyne.KeyT)
	app.simulateKeyPress(fyne.KeyE)
	app.simulateKeyPress(fyne.KeyI)
	app.editEntryTypeDialog.Type(text)
	app.simulateKeyPress(fyne.KeyReturn)
}
