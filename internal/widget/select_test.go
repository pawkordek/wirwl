package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatSelectHighlightingUnhighlightingWorks(t *testing.T) {
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "")
	//Select needs to be placed into a test window, otherwise renderer doesn't work properly
	test.NewApp().NewWindow("").SetContent(selectWidget)
	assert.Equal(t, theme.BackgroundColor(), selectWidget.backgroundRenderer.BackgroundColor())
	selectWidget.Highlight()
	assert.Equal(t, theme.FocusColor(), selectWidget.backgroundRenderer.BackgroundColor())
	selectWidget.Unhighlight()
	assert.Equal(t, theme.BackgroundColor(), selectWidget.backgroundRenderer.BackgroundColor())
}

func TestThatWhenEnteringInputModeMenuGetsShownOverSelectWithChoices(t *testing.T) {
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "1", "2")
	selectWidget.EnterInputMode()
	assert.True(t, selectWidget.menu.Visible())
	assert.True(t, selectWidget.menu.Focused())
	assert.Equal(t, selectWidget.menu, test.Canvas().Focused())
	assert.Equal(t, 2, len(selectWidget.menu.choices))
}

func TestThatWhenExitingInputModeSelectedChoiceInMenuGetsSetAsSelectValue(t *testing.T) {
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "1", "2")
	selectWidget.EnterInputMode()
	SimulateKeyPress(selectWidget.menu, fyne.KeyJ)
	SimulateKeyPress(selectWidget.menu, fyne.KeyReturn)
	assert.Equal(t, "2", selectWidget.Selected)
}

func TestThatWhenExitingInputModeMenuGetsHidden(t *testing.T) {
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "1", "2")
	selectWidget.EnterInputMode()
	SimulateKeyPress(selectWidget.menu, fyne.KeyReturn)
	assert.False(t, selectWidget.menu.Visible())
	assert.False(t, selectWidget.menu.Focused())
	assert.NotEqual(t, selectWidget.menu, test.Canvas().Focused())
}

func TestThatWhenChoiceWasMadeInMenuSelectUnfocusesAndCallsOnExitInpuModeFunction(t *testing.T) {
	functionCalled := false
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "1", "2")
	selectWidget.EnterInputMode()
	selectWidget.SetOnExitInputModeFunction(func() { functionCalled = true })
	SimulateKeyPress(selectWidget.menu, fyne.KeyReturn)
	assert.False(t, selectWidget.Focused())
	assert.NotEqual(t, selectWidget, test.Canvas().Focused())
	assert.True(t, functionCalled)
}

func TestThatSetTextFunctionSetsSelectedValue(t *testing.T) {
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "1", "2")
	selectWidget.SetText("2")
	assert.Equal(t, "2", selectWidget.Selected)
}

func TestThatReturnedTextIsEqualToSelectedValue(t *testing.T) {
	selectWidget := NewSelect(test.Canvas(), getInputHandlerForTesting(), "1", "2")
	selectWidget.Selected = "2"
	assert.Equal(t, "2", selectWidget.GetText())
}
