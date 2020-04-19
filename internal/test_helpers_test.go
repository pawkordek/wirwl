package wirwl

import (
	"fyne.io/fyne"
	"os/user"
)

func getLoggingDirForTesting() string {
	currentUser, err := user.Current()
	if err != nil {
		return "/tmp/wirwl/"
	}
	return currentUser.HomeDir + "/.local/share/wirwl/"
}

func (app *App) SimulateKeyPress(key fyne.KeyName) {
	event := &fyne.KeyEvent{Name: key}
	focusedElement := app.mainWindow.Canvas().Focused()
	if focusedElement != nil {
		focusedElement.TypedKey(event)
	} else {
		onTypedKey := app.mainWindow.Canvas().OnTypedKey()
		onTypedKey(event)
	}
}

func (app *App) SimulateSwitchingToNextEntryType() {
	app.SimulateKeyPress(fyne.KeyL)
}

func (app *App) SimulateSwitchingToPreviousEntryType() {
	app.SimulateKeyPress(fyne.KeyH)
}

func (app *App) SimulateOpeningDialogForAddingEntryType() {
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyI)
}

func (app *App) SimulateAddingNewEntryTypeWithName(text string) {
	app.SimulateOpeningDialogForAddingEntryType()
	app.SimulateKeyPress(fyne.KeyI)
	app.addEntryTypeDialog.Type(text)
	app.SimulateKeyPress(fyne.KeyEnter)
}

func (app *App) SimulateSavingChanges() {
	app.SimulateKeyPress(fyne.KeyS)
}

func (app *App) SimulateAttemptAtDeletionOfCurrentEntryType() {
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyD)
}

func (app *App) SimulateDeletionOfCurrentEntryType() {
	app.SimulateAttemptAtDeletionOfCurrentEntryType()
	app.SimulateKeyPress(fyne.KeyY)
}

func (app *App) SimulateEditionOfCurrentEntryTypeTo(text string) {
	app.SimulateKeyPress(fyne.KeyT)
	app.SimulateKeyPress(fyne.KeyE)
	app.SimulateKeyPress(fyne.KeyI)
	app.editEntryTypeDialog.Type(text)
	app.SimulateKeyPress(fyne.KeyEnter)
}
