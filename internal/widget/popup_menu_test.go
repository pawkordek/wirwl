package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatPopUpMenuIsFocusedAndVisibleWhenItGetsShown(t *testing.T) {
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "")
	menu.Show()
	assert.True(t, menu.focused)
	assert.True(t, menu.Visible())
	assert.Equal(t, menu, test.Canvas().Focused())
	menu = NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "")
	menu.ShowAtPosition(fyne.Position{})
	assert.True(t, menu.focused)
	assert.True(t, menu.Visible())
	assert.Equal(t, menu, test.Canvas().Focused())
}

func TestThatPopUpMenuHidesAndUnfocusesWhenChoiceGetsSelected(t *testing.T) {
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "")
	menu.Show()
	SimulateKeyPressOnTestCanvas(fyne.KeyReturn)
	assert.False(t, menu.focused)
	assert.False(t, menu.Visible())
	assert.NotEqual(t, menu, test.Canvas().Focused())
}

func TestThatMenuChoicesCanBeSelectedAndIndicateThatTheyAre(t *testing.T) {
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "1", "2", "3")
	menu.Show()
	assert.True(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[0])
	SimulateKeyPressOnTestCanvas(fyne.KeyK)
	assert.True(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[0])
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.True(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[1])
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.True(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[2])
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.True(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[2])
	SimulateKeyPressOnTestCanvas(fyne.KeyK)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.True(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[1])
}

func TestThatWhenChoiceGetsSelectedMenuReturnsProperValue(t *testing.T) {
	returnedValue := ""
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "1", "2", "3")
	menu.OnChoiceSelectedCallback = func(selectedChoice string) {
		returnedValue = selectedChoice
	}
	menu.Show()
	SimulateKeyPressOnTestCanvas(fyne.KeyJ)
	SimulateKeyPressOnTestCanvas(fyne.KeyReturn)
	assert.Equal(t, "2", returnedValue)
}
