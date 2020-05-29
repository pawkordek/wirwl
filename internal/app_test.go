package wirwl

import (
	"fyne.io/fyne"
	fyneTest "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/widget"
)

func TestMain(m *testing.M) {
	testSetup()
	exitCode := m.Run()
	testCleanup()
	os.Exit(exitCode)
}

func TestThatEntriesTabsWithContentDisplayInCorrectOrder(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	assert.Equal(t, 3, len(app.entriesTypesTabs.Items()))
	firstTab := app.entriesTypesTabs.Items()[0]
	assert.Equal(t, firstTab.Text, "comics")
	assert.Equal(t, 0, widget.GetLabelPositionInContent(firstTab.Content, "some comic1"))
	assert.Equal(t, 1, widget.GetLabelPositionInContent(firstTab.Content, "some comic2"))
	secondTab := app.entriesTypesTabs.Items()[1]
	assert.Equal(t, secondTab.Text, "music")
	assert.Equal(t, 0, widget.GetLabelPositionInContent(secondTab.Content, "some music1"))
	assert.Equal(t, 1, widget.GetLabelPositionInContent(secondTab.Content, "some music2"))
	thirdTab := app.entriesTypesTabs.Items()[2]
	assert.Equal(t, thirdTab.Text, "videos")
	assert.Equal(t, 0, widget.GetLabelPositionInContent(thirdTab.Content, "some video1"))
	assert.Equal(t, 1, widget.GetLabelPositionInContent(thirdTab.Content, "some video2"))
}

func TestSwitchingTabs(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	assert.Equal(t, "comics", app.getCurrentTabText())
	app.simulateSwitchingToNextEntryType()
	assert.Equal(t, "music", app.getCurrentTabText())
	app.simulateSwitchingToNextEntryType()
	assert.Equal(t, "videos", app.getCurrentTabText())
	app.simulateSwitchingToNextEntryType()
	assert.Equal(t, "comics", app.getCurrentTabText())
	app.simulateSwitchingToPreviousEntryType()
	assert.Equal(t, "videos", app.getCurrentTabText())
	app.simulateSwitchingToPreviousEntryType()
	assert.Equal(t, "music", app.getCurrentTabText())
	app.simulateSwitchingToPreviousEntryType()
	assert.Equal(t, "comics", app.getCurrentTabText())
}

func TestEntryHighlightingWhenSwitchingTabs(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	currentTab := app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
	app.simulateSwitchingToNextEntryType()
	currentTab = app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some music1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some music2").TextStyle)
	app.simulateSwitchingToPreviousEntryType()
	currentTab = app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
}

func TestThatApplicationDoesNotCrashWhenTryingToSwitchToATabThatDoesNotExist(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateSwitchingToNextEntryType()
	app.simulateSwitchingToNextEntryType()
	app.simulateSwitchingToPreviousEntryType()
	app.simulateSwitchingToPreviousEntryType()
}

func TestThatIfThereAreNoEntriesCorrectMessageDisplays(t *testing.T) {
	app, cleanup := setupAndRunAppAsIfRunForFirstTime()
	defer cleanup()
	assert.Equal(t, 1, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "No entries", app.entriesTypesTabs.Items()[0].Text)
}

func TestWhetherDialogForAddingEntryTypesOpens(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	assert.Equal(t, true, app.addEntryTypeDialog.Hidden)
	app.simulateOpeningDialogForAddingEntryType()
	assert.Equal(t, true, app.addEntryTypeDialog.Visible())
	assert.Equal(t, true, app.addEntryTypeDialog.Focused())
}

func TestWhetherReopeningDialogForAddingEntriesTypesDoesNotPersistPreviouslyInputText(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("some type")
	app.simulateOpeningDialogForAddingEntryType()
	name := app.addEntryTypeDialog.ItemValue("Name")
	imageQuery := app.addEntryTypeDialog.ItemValue("Image query")
	assert.Empty(t, name)
	assert.Empty(t, imageQuery)
}

func TestAddingOfNewEntryType(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("new entry type")
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "comics", app.entriesTypesTabs.Items()[0].Text)
	assert.Equal(t, "music", app.entriesTypesTabs.Items()[1].Text)
	assert.Equal(t, "new entry type", app.entriesTypesTabs.Items()[2].Text)
	assert.Equal(t, "videos", app.entriesTypesTabs.Items()[3].Text)
	assert.Equal(t, true, app.addEntryTypeDialog.Hidden)
	assert.Equal(t, true, !app.addEntryTypeDialog.Focused())
}

func TestThatItIsNotPossibleToAddTheSameEntryTypeTwice(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("type")
	app.simulateAddingNewEntryTypeWithName("type")
	assert.Equal(t, true, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "Entry type with name 'type' already exists.", app.msgDialog.Msg())
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
}

func TestThatPressingAnyKeyClosesMessagePopUp(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.msgDialog.Display("", "")
	app.simulateKeyPress(fyne.KeyT)
	assert.Equal(t, true, app.msgDialog.Hidden)
	app.msgDialog.Display("", "")
	app.simulateKeyPress(fyne.KeyReturn)
	assert.Equal(t, true, app.msgDialog.Hidden)
}

func TestThatAddingNewEntryTypeDoesNotChangeCurrentlyOpenedTab(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateSwitchingToNextEntryType()
	app.simulateAddingNewEntryTypeWithName("type")
	assert.Equal(t, "music", app.getCurrentTabText())
}

func TestThatAfterAddingNewEntryOpenedTabStillHasTheSameElementHighlighted(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("type")
	currentTab := app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
}

func TestThatAfterTryingToAddExistingEntryTypeAndClosingWarningMessageAboutItDialogForAddingEntriesTypesIsStillOpen(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("comics")
	assert.True(t, app.addEntryTypeDialog.Hidden)
	assert.False(t, app.addEntryTypeDialog.Focused())
	app.simulateKeyPress(fyne.KeyEscape)
	assert.True(t, app.addEntryTypeDialog.Visible())
	assert.True(t, app.addEntryTypeDialog.Focused())
}

func TestThatAfterTryingToAddEntryTypeWithEmptyNameWarningMessageDisplaysAndEntryTypeIsNotAdded(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("")
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "You cannot add entry type with empty name", app.msgDialog.Msg())
	assert.Equal(t, 3, len(app.entriesTypesTabs.Items()))
}

func TestThatSavingChangesWorks(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("type")
	app.simulateSavingChanges()
	config := getTestConfigWithConfigPathIn("/tmp/")
	config.AppDataDirPath = testAppDataDirPath
	app = NewApp(fyneTest.NewApp(), config)
	app.LoadAndDisplay()
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "comics", app.getCurrentTabText())
}

func TestThatAfterSavingSuccessfullySuccessDialogDisplays(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateSavingChanges()
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "SUCCESS", app.msgDialog.Title())
	assert.Equal(t, "Changes saved.", app.msgDialog.Msg())
}

func TestThatAfterSavingUnsuccessfullyErrorDialogDisplays(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.dataProvider = data.NewAlwaysFailingProvider()
	app.simulateSavingChanges()
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, data.AlwaysFailingProviderError.Error(), app.msgDialog.Msg())
}

func TestThatDeletingEntriesTypesWorks(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateDeletionOfCurrentEntryType()
	assert.Equal(t, 2, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "music", app.entriesTypesTabs.Items()[0].Text)
	assert.Equal(t, "videos", app.entriesTypesTabs.Items()[1].Text)
}

func TestThatWhenTryingToDeleteLastEntryTypeItIsPreventedAndWarningDialogIsDisplayed(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateDeletionOfCurrentEntryType()
	app.simulateDeletionOfCurrentEntryType()
	app.simulateAttemptAtDeletionOfCurrentEntryType()
	assert.Equal(t, 1, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, true, app.msgDialog.Visible())
	assert.Equal(t, "WARNING", app.msgDialog.Title())
	assert.Equal(t, "You cannot remove the only remaining entry type!", app.msgDialog.Msg())
}

func TestThatEditingEntryTypeWorks(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateEditionOfCurrentEntryTypeTo("2")
	assert.Equal(t, "2comics", app.entriesTypesTabs.CurrentTab().Text)
}

func TestThatEditingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.LoadAndDisplay()
	app.simulateEditionOfCurrentEntryTypeTo("2")
	app.simulateSavingChanges()
	config := getTestConfigWithConfigPathIn(testAppDataDirPath)
	config.AppDataDirPath = testAppDataDirPath
	app2 := NewApp(fyneTest.NewApp(), config)
	app2.LoadAndDisplay()
	assert.Equal(t, "2comics", app2.entriesTypesTabs.CurrentTab().Text)
}

func TestThatDeletingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	app, cleanup := setupAndRunAppForTestingWithTestConfig()
	defer cleanup()
	app.simulateDeletionOfCurrentEntryType()
	app.simulateSavingChanges()
	config := app.config
	app2 := NewApp(fyneTest.NewApp(), config)
	app2.LoadAndDisplay()
	assert.Equal(t, "music", app2.entriesTypesTabs.CurrentTab().Text)
}
