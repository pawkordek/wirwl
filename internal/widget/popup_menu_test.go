package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatPopUpMenuIsFocusedAndVisibleWhenItGetsShown(t *testing.T) {
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "")
	menu.Show()
	assert.True(t, menu.focused)
	assert.Equal(t, menu, test.Canvas().Focused())
	menu = NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "")
	menu.ShowAtPosition(fyne.Position{})
	assert.True(t, menu.focused)
	assert.Equal(t, menu, test.Canvas().Focused())
}

func TestThatPopUpMenuHidesAndUnfocusesWhenChoiceGetsSelected(t *testing.T) {
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "")
	menu.Show()
	SimulateKeyPress(menu, fyne.KeyReturn)
	assert.False(t, menu.focused)
	assert.NotEqual(t, menu, test.Canvas().Focused())
}

func TestThatMenuChoicesCanBeSelectedAndIndicateThatTheyAre(t *testing.T) {
	menu := NewPopUpMenu(test.Canvas(), getInputHandlerForTesting(), "1", "2", "3")
	menu.Show()
	assert.True(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[0])
	SimulateKeyPress(menu, fyne.KeyK)
	assert.True(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[0])
	SimulateKeyPress(menu, fyne.KeyJ)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.True(t, menu.choices[1].TextStyle.Bold)
	assert.False(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[1])
	SimulateKeyPress(menu, fyne.KeyJ)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.True(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[2])
	SimulateKeyPress(menu, fyne.KeyJ)
	assert.False(t, menu.choices[0].TextStyle.Bold)
	assert.False(t, menu.choices[1].TextStyle.Bold)
	assert.True(t, menu.choices[2].TextStyle.Bold)
	assert.Equal(t, menu.currentChoice(), menu.choices[2])
	SimulateKeyPress(menu, fyne.KeyK)
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
	SimulateKeyPress(menu, fyne.KeyJ)
	SimulateKeyPress(menu, fyne.KeyReturn)
	assert.Equal(t, "2", returnedValue)
}
