package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"wirwl/internal/log"
)

func TestThatFormDisplays(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "form dialog", getTwoInputFieldsForFormDialogTesting())
	dialog.Display()
	assert.True(t, dialog.Visible())
	assert.Equal(t, "form dialog", dialog.Title())
	assert.NotNil(t, dialog.embeddedWidgets["first"])
	assert.NotNil(t, dialog.embeddedWidgets["second"])
}

func TestThatFirstInputIsTheCurrentWhenDialogDisplays(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getTwoInputFieldsForFormDialogTesting())
	dialog.Display()
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["first"])
	log.Info(dialog.form.Items[0].Text)
	log.Info(dialog.form.Items[1].Text)
}

func TestThatPressingJAndKSwitchesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting())
	dialog.Display()
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["first"])
	SimulateKeyPress(dialog, fyne.KeyJ)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["second"])
	SimulateKeyPress(dialog, fyne.KeyJ)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["third"])
	SimulateKeyPress(dialog, fyne.KeyK)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["second"])
	SimulateKeyPress(dialog, fyne.KeyK)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["first"])
}

func TestThatPressingIFocusesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting())
	dialog.Display()
	assert.False(t, dialog.currentWidget().Focused())
	SimulateKeyPress(dialog, fyne.KeyI)
	assert.True(t, dialog.currentWidget().Focused())
}

func TestThatPressingEscapeUnfocusesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting())
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyI)
	SimulateKeyPress(dialog.currentWidget(), fyne.KeyEscape)
	assert.False(t, dialog.currentWidget().Focused())
}

func TestThatPressingEnterCallsFunctionAndHidesAndUnfocusesDialog(t *testing.T) {
	functionCalled := false
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting())
	dialog.OnEnterPressed = func() { functionCalled = true }
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyReturn)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
	assert.True(t, functionCalled)
}

func TestThatFirstInputIsSelectedOnDialogReopening(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting())
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyJ)
	SimulateKeyPress(dialog, fyne.KeyEnter)
	dialog.Display()
	assert.Equal(t, dialog.embeddedWidgets["first"], dialog.currentWidget())
}

func TestThatPressingEscapeWhenNotInEditionModeClosesDialog(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getOneInputFieldForDialogTesting())
	dialog.Display()
	SimulateKeyPress(dialog, fyne.KeyEscape)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
}

func TestThatSettingAndGettingItemValueWorksCorrectly(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getOneInputFieldForDialogTesting())
	setValue := "test value"
	dialog.SetItemValue("first", setValue)
	receivedValue := dialog.ItemValue("first")
	assert.Equal(t, setValue, receivedValue)
}

func TestThatSettingItemValueOnNonExistingItemDoesNotPanic(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getOneInputFieldForDialogTesting())
	dialog.SetItemValue("non existing", "value")
}

func TestThatCleaningItemValuesWorks(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getTwoInputFieldsForFormDialogTesting())
	dialog.SetItemValue("first", "val1")
	dialog.SetItemValue("second", "val2")
	dialog.CleanItemValues()
	assert.Empty(t, dialog.embeddedWidgets["first"].GetText())
	assert.Empty(t, dialog.embeddedWidgets["second"].GetText())
}

func TestThatFormDialogHidesBeforeItCallsOnEnterPressed(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getTwoInputFieldsForFormDialogTesting())
	dialog.Display()
	dialog.OnEnterPressed = func() {
		assert.Nil(t, dialog.Canvas.Overlay())
	}
	SimulateKeyPress(dialog, fyne.KeyEnter)
}
