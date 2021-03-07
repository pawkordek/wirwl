package wirwl

import (
	"errors"
	"fyne.io/fyne"
	fyneWidget "fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/input"
	"wirwl/internal/log"
	"wirwl/internal/widget"
)

func TestThatLoadingErrorsMsgDialogDoesNotDisplayIfThereAreNoErrors(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	assert.True(t, app.msgDialog.Hidden)
}

func TestThatErrorsPassedInOnAppCreationDisplayAfterItRuns(t *testing.T) {
	configurator := NewTestAppConfigurator()
	loadingErrors := make(map[string]string)
	loadingErrors["some error"] = "Some error occurred"
	loadingErrors["some other error"] = "Some other error occurred"
	app, cleanup := configurator.prepareConfiguratorForTestingWithExistingData().
		setLoadingErrors(loadingErrors).
		createTestApplication().
		getRunningTestApplication()
	defer cleanup()
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Contains(t, app.msgDialog.Msg(), "Some error occurred")
	assert.Contains(t, app.msgDialog.Msg(), "Some other error occurred")
}

func TestThatErrorDisplaysWhenEntriesFailToLoad(t *testing.T) {
	configurator := NewTestAppConfigurator()
	dataProvider := data.NewAbstractProvider()
	dataProvider.LoadEntriesFunc = func() (map[data.EntryType][]data.Entry, error) {
		return nil, errors.New("An error occured when entries failed to load")
	}
	app, cleanup := configurator.prepareConfiguratorForTestingWithExistingData().
		setDataProvider(dataProvider).
		createTestApplication().
		getRunningTestApplication()
	defer cleanup()
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Contains(t, app.msgDialog.Msg(), "Failed to load entries. Application will now exit as it cannot continue.")
}

func TestThatDbFileWithItsDirGetsCreatedInAppDataDirFromConfig(t *testing.T) {
	dbFilePath := testAppDataDirPath + "data.db"
	configurator := NewTestAppConfigurator()
	_, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	_, err := os.Stat(dbFilePath)
	assert.Nil(t, err)
}

func TestThatCorrectConfigFileGetsWrittenToDiskAfterApplicationExits(t *testing.T) {
	configurator := NewTestAppConfigurator()
	savedConfig := Config{
		AppDataDirPath: testAppDataDirPath,
		ConfigDirPath:  testConfigDirPath,
		Keymap:         map[input.Action]input.KeyCombination{},
	}
	savedConfig.loadDefaultKeymap()
	_, cleanup := configurator.
		prepareConfiguratorForTestingWithExistingData().setConfig(savedConfig).
		createTestApplication().
		getRunningTestApplication()
	defer cleanup()
	loadedConfig := NewConfig(testConfigDirPath)
	err := loadedConfig.load()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, savedConfig, loadedConfig)
}

func TestThatEntriesTabsWithContentDisplayInCorrectOrder(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
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
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
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
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
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
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateSwitchingToNextEntryType()
	app.simulateSwitchingToNextEntryType()
	app.simulateSwitchingToPreviousEntryType()
	app.simulateSwitchingToPreviousEntryType()
}

func TestThatIfThereAreNoEntriesCorrectMessageDisplays(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatWillRunForFirstTime().getRunningTestApplication()
	defer cleanup()
	assert.Equal(t, 1, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "No entries", app.entriesTypesTabs.Items()[0].Text)
}

func TestWhetherDialogForAddingEntryTypesOpens(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	assert.Equal(t, true, app.addEntryTypeDialog.Hidden)
	app.simulateOpeningDialogForAddingEntryType()
	assert.Equal(t, true, app.addEntryTypeDialog.Visible())
	assert.Equal(t, true, app.addEntryTypeDialog.Focused())
}

func TestWhetherReopeningDialogForAddingEntriesTypesDoesNotPersistPreviouslyInputText(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("some type")
	app.simulateOpeningDialogForAddingEntryType()
	name := app.addEntryTypeDialog.ItemValue("Name")
	imageQuery := app.addEntryTypeDialog.ItemValue("Image query")
	assert.Empty(t, name)
	assert.Empty(t, imageQuery)
}

func TestAddingOfNewEntryType(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
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
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("type")
	app.simulateAddingNewEntryTypeWithName("type")
	assert.Equal(t, true, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "Entry type with name 'type' already exists", app.msgDialog.Msg())
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
}

func TestThatPressingAnyKeyClosesMessagePopUp(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.msgDialog.Display("", "")
	app.simulateKeyPress(fyne.KeyT)
	assert.Equal(t, true, app.msgDialog.Hidden)
	app.msgDialog.Display("", "")
	app.simulateKeyPress(fyne.KeyReturn)
	assert.Equal(t, true, app.msgDialog.Hidden)
}

func TestThatAddingNewEntryTypeDoesNotChangeCurrentlyOpenedTab(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateSwitchingToNextEntryType()
	app.simulateAddingNewEntryTypeWithName("type")
	assert.Equal(t, "music", app.getCurrentTabText())
}

func TestThatAfterAddingNewEntryOpenedTabStillHasTheSameElementHighlighted(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("type")
	currentTab := app.entriesTypesTabs.CurrentTab()
	assert.Equal(t, fyne.TextStyle{Bold: true}, widget.GetLabelFromContent(currentTab.Content, "some comic1").TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: false}, widget.GetLabelFromContent(currentTab.Content, "some comic2").TextStyle)
}

func TestThatAfterTryingToAddExistingEntryTypeAndClosingWarningMessageAboutItDialogForAddingEntriesTypesIsStillOpen(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("comics")
	assert.True(t, app.addEntryTypeDialog.Hidden)
	assert.False(t, app.addEntryTypeDialog.Focused())
	app.simulateKeyPress(fyne.KeyEscape)
	assert.True(t, app.addEntryTypeDialog.Visible())
	assert.True(t, app.addEntryTypeDialog.Focused())
}

func TestThatAfterTryingToAddEntryTypeWithEmptyNameWarningMessageDisplaysAndEntryTypeIsNotAdded(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("")
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "Cannot add entry type with an empty name", app.msgDialog.Msg())
	assert.Equal(t, 3, len(app.entriesTypesTabs.Items()))
}

func TestThatSavingChangesWorks(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateAddingNewEntryTypeWithName("type")
	app.simulateSavingChanges()
	app, cleanup = configurator.getRunningTestApplication()
	defer cleanup()
	assert.Equal(t, 4, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "comics", app.getCurrentTabText())
}

func TestThatAfterSavingSuccessfullySuccessDialogDisplays(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateSavingChanges()
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "SUCCESS", app.msgDialog.Title())
	assert.Equal(t, "Changes saved.", app.msgDialog.Msg())
}

func TestThatAfterSavingUnsuccessfullyErrorDialogDisplays(t *testing.T) {
	providerFailingOnSave := data.NewAbstractProvider()
	providerFailingOnSave.SaveEntriesFunc = func(m map[data.EntryType][]data.Entry) error {
		return errors.New("Testing that saving failed")
	}
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.prepareConfiguratorForTestingWithExistingData().
		setDataProvider(providerFailingOnSave).
		createTestApplication().getRunningTestApplication()
	defer cleanup()
	app.simulateSavingChanges()
	assert.True(t, app.msgDialog.Visible())
	assert.Equal(t, "ERROR", app.msgDialog.Title())
	assert.Equal(t, "Testing that saving failed", app.msgDialog.Msg())
}

func TestThatDeletingEntriesTypesWorks(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateDeletionOfCurrentEntryType()
	assert.Equal(t, 2, len(app.entriesTypesTabs.Items()))
	assert.Equal(t, "music", app.entriesTypesTabs.Items()[0].Text)
	assert.Equal(t, "videos", app.entriesTypesTabs.Items()[1].Text)
}

func TestThatWhenTryingToDeleteLastEntryTypeItIsPreventedAndWarningDialogIsDisplayed(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
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
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateEditionOfCurrentEntryTypeTo("2")
	assert.Equal(t, "2comics", app.entriesTypesTabs.CurrentTab().Text)
}

func TestThatEditingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateEditionOfCurrentEntryTypeTo("2")
	app.simulateSavingChanges()
	app, cleanup = configurator.getRunningTestApplication()
	defer cleanup()
	assert.Equal(t, "2comics", app.entriesTypesTabs.CurrentTab().Text)
}

func TestThatDeletingEntryTypePersistsAfterReopeningTheApplication(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	app.simulateDeletionOfCurrentEntryType()
	app.simulateSavingChanges()
	app, cleanup = configurator.getRunningTestApplication()
	defer cleanup()
	assert.Equal(t, "music", app.entriesTypesTabs.CurrentTab().Text)
}

func TestThatConfigIsNotSavedIfItFailedToLoad(t *testing.T) {
	configurator := NewTestAppConfigurator()
	loadingErrors := make(map[string]string)
	loadingErrors[configLoadError] = "Config failed to load"
	_, cleanup := configurator.prepareConfiguratorForTestingWithExistingData().
		setLoadingErrors(loadingErrors).
		createFailingToLoadConfigFile().
		createTestApplication().
		getRunningTestApplication()
	defer cleanup()
	configFileContents, err := ioutil.ReadFile(filepath.Join(testConfigDirPath + "wirwl.cfg"))
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, failingToLoadConfigFileContents, string(configFileContents))
}

func TestThatRecentlyPressedKeysLabelDisplaysTheCorrectInformation(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	//The first single key has to not have any action in default bindings, otherwise the action is gonna get executed
	//instead of the key making a combination with the second key which will fail the test
	assert.Equal(t, "Recently pressed keys: ", app.recentlyPressedKeysLabel.Text)
	app.simulateKeyPress(fyne.KeyQ)
	assert.Equal(t, "Recently pressed keys: Q", app.recentlyPressedKeysLabel.Text)
	app.simulateKeyPress(fyne.KeyG)
	assert.Equal(t, "Recently pressed keys: Q,G", app.recentlyPressedKeysLabel.Text)
	app.simulateKeyPress(fyne.KeyEscape)
	assert.Equal(t, "Recently pressed keys: Escape", app.recentlyPressedKeysLabel.Text)
}

func TestThatEntriesTablesHaveCorrectAmountOfHeaderColumns(t *testing.T) {
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	for entryType, entriesTable := range app.entriesTables {
		amountOfHeaderColumns := len(entriesTable.HeaderColumns())
		assert.Equal(t, 14, amountOfHeaderColumns, "The table for entry type "+entryType.Name+" has incorrect amount of header columns")
	}
}

func TestThatEntriesTablesHaveCorrectHeaderLabels(t *testing.T) {
	headerColumnsNames := []string{
		"Num",
		"Image",
		"Status",
		"Title",
		"Elements completed",
		"Total amount",
		"Score",
		"Start date",
		"Finish date",
		"Link",
		"Description",
		"Comment",
		"Tags",
		"Image query",
	}
	configurator := NewTestAppConfigurator()
	app, cleanup := configurator.createTestApplicationThatUsesExistingData().getRunningTestApplication()
	defer cleanup()
	for entryType, entriesTable := range app.entriesTables {
		for _, column := range entriesTable.HeaderColumns() {
			label := column.(*fyneWidget.Label)
			assert.Contains(t, headerColumnsNames, label.Text, "The table for entry type "+entryType.Name+" doesn't have a header label with text "+label.Text)
		}
	}
}
