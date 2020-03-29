package widget

import (
	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tabsData = map[string][]string{
	"First tab":  {"a1", "b1", "c1"},
	"Second tab": {"a2", "b2", "c2"},
	"Third tab":  {"a3", "b3", "c3"},
}

func TestThatOnDisplayCorrectTabAndContentDisplays(t *testing.T) {
	container := NewTabContainer(tabsData)
	assert.Equal(t, "First tab", container.it.CurrentTab().Text)
	label1 := GetLabelFromContent(container.it.CurrentTab().Content, "a1")
	assert.NotNil(t, label1)
	label2 := GetLabelFromContent(container.it.CurrentTab().Content, "b1")
	assert.NotNil(t, label2)
	label3 := GetLabelFromContent(container.it.CurrentTab().Content, "c1")
	assert.NotNil(t, label3)
}

func TestThatFirstLabelOnFirstTabIsBoldedAtFirst(t *testing.T) {
	container := NewTabContainer(tabsData)
	label1 := GetLabelFromContent(container.it.CurrentTab().Content, "a1")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label1.TextStyle)
	label2 := GetLabelFromContent(container.it.CurrentTab().Content, "b1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label2.TextStyle)
	label3 := GetLabelFromContent(container.it.CurrentTab().Content, "c1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label3.TextStyle)
}

func TestThatSelectingTabsWorks(t *testing.T) {
	container := NewTabContainer(tabsData)
	container.SelectNextTab()
	assert.Equal(t, "Second tab", container.it.CurrentTab().Text)
	container.SelectNextTab()
	assert.Equal(t, "Third tab", container.it.CurrentTab().Text)
	container.SelectNextTab()
	assert.Equal(t, "First tab", container.it.CurrentTab().Text)
	container.SelectPreviousTab()
	assert.Equal(t, "Third tab", container.it.CurrentTab().Text)
	container.SelectPreviousTab()
	assert.Equal(t, "Second tab", container.it.CurrentTab().Text)
	container.SelectPreviousTab()
	assert.Equal(t, "First tab", container.it.CurrentTab().Text)
}

func TestThatChangingTabSelectsFirstLabelOnTab(t *testing.T) {
	container := NewTabContainer(tabsData)
	container.SelectNextTab()
	label1 := GetLabelFromContent(container.it.CurrentTab().Content, "a2")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label1.TextStyle)
	label2 := GetLabelFromContent(container.it.CurrentTab().Content, "b2")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label2.TextStyle)
	label3 := GetLabelFromContent(container.it.CurrentTab().Content, "c2")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label3.TextStyle)
	container.SelectPreviousTab()
	label4 := GetLabelFromContent(container.it.CurrentTab().Content, "a1")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label4.TextStyle)
	label5 := GetLabelFromContent(container.it.CurrentTab().Content, "b1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label5.TextStyle)
	label6 := GetLabelFromContent(container.it.CurrentTab().Content, "c1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label6.TextStyle)
}
