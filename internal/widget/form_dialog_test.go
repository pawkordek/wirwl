package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"wirwl/internal/log"
)

func TestThatFormDisplays(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "form dialog", getTwoInputFieldsForFormDialogTesting()...)
	dialog.Display()
	assert.True(t, dialog.Visible())
	assert.Equal(t, "form dialog", dialog.Title())
	assert.NotNil(t, dialog.embeddedWidgets["first"])
	assert.NotNil(t, dialog.embeddedWidgets["second"])
}

func TestThatFirstInputIsTheCurrentWhenDialogDisplays(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getTwoInputFieldsForFormDialogTesting()...)
	dialog.Display()
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["first"])
	log.Info(dialog.form.Items[0].Text)
	log.Info(dialog.form.Items[1].Text)
}

func TestThatPressingJAndKSwitchesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting()...)
	dialog.Display()
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["first"])
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["second"])
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["third"])
	SimulateKeyPressOnTestCanvas(fyne.KeyK)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["second"])
	SimulateKeyPressOnTestCanvas(fyne.KeyK)
	assert.Equal(t, dialog.currentWidget(), dialog.embeddedWidgets["first"])
}

func TestThatPressingIFocusesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting()...)
	dialog.Display()
	assert.NotEqual(t, dialog.currentWidget(), test.Canvas().Focused())
	SimulateKeyPressOnTestCanvas(fyne.KeyI)
	assert.Equal(t, dialog.currentWidget(), test.Canvas().Focused())
}

func TestThatPressingEscapeUnfocusesCurrentInput(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting()...)
	dialog.Display()
	SimulateKeyPressOnTestCanvas(fyne.KeyI)
	SimulateKeyPressOnTestCanvas(fyne.KeyEscape)
	assert.NotEqual(t, dialog.currentWidget(), test.Canvas().Focused())
}

func TestThatPressingEnterCallsFunctionAndHidesAndUnfocusesDialog(t *testing.T) {
	functionCalled := false
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting()...)
	dialog.OnEnterPressed = func() { functionCalled = true }
	dialog.Display()
	SimulateKeyPressOnTestCanvas(fyne.KeyReturn)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
	assert.True(t, functionCalled)
}

func TestThatFirstInputIsSelectedOnDialogReopening(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getThreeInputFieldsForFormDialogTesting()...)
	dialog.Display()
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	SimulateKeyPressOnTestCanvas(fyne.KeyEnter)
	dialog.Display()
	assert.Equal(t, dialog.embeddedWidgets["first"], dialog.currentWidget())
}

func TestThatPressingEscapeWhenNotInEditionModeClosesDialog(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getOneInputFieldForDialogTesting()...)
	dialog.Display()
	SimulateKeyPressOnTestCanvas(fyne.KeyEscape)
	assert.True(t, dialog.Hidden)
	assert.False(t, dialog.Focused())
}

func TestThatSettingAndGettingItemValueWorksCorrectly(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getOneInputFieldForDialogTesting()...)
	setValue := "test value"
	dialog.SetItemValue("first", setValue)
	receivedValue := dialog.ItemValue("first")
	assert.Equal(t, setValue, receivedValue)
}

func TestThatSettingItemValueOnNonExistingItemDoesNotPanic(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getOneInputFieldForDialogTesting()...)
	dialog.SetItemValue("non existing", "value")
}

func TestThatCleaningItemValuesWorks(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getTwoInputFieldsForFormDialogTesting()...)
	dialog.SetItemValue("first", "val1")
	dialog.SetItemValue("second", "val2")
	dialog.CleanItemValues()
	assert.Empty(t, dialog.embeddedWidgets["first"].GetText())
	assert.Empty(t, dialog.embeddedWidgets["second"].GetText())
}

func TestThatFormDialogHidesBeforeItCallsOnEnterPressed(t *testing.T) {
	dialog := NewFormDialog(test.Canvas(), getInputHandlerForTesting(), "", getTwoInputFieldsForFormDialogTesting()...)
	dialog.Display()
	dialog.OnEnterPressed = func() {
		assert.Nil(t, dialog.Canvas.Overlays())
	}
	SimulateKeyPressOnTestCanvas(fyne.KeyEnter)
}

func TestThatFormDialogItemFactoryCreatesCorrectInputField(t *testing.T) {
	createdFormItem := NewFormDialogFormItemFactory(test.Canvas(), getInputHandlerForTesting()).
		FormItemWithInputField("This is input field")
	assert.Equal(t, "This is input field", createdFormItem.Text)
	assert.NotPanics(t, func() {
		_ = createdFormItem.Widget.(*InputField)
	})
}

func TestThatFormDialogItemFactoryCreatesCorrectSelect(t *testing.T) {
	createdFormItem := NewFormDialogFormItemFactory(test.Canvas(), getInputHandlerForTesting()).
		FormItemWithSelect("This is select", "1", "2")
	assert.Equal(t, "This is select", createdFormItem.Text)
	assert.NotPanics(t, func() {
		_ = createdFormItem.Widget.(*Select)
	})
}
