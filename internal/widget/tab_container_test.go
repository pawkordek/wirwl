package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tabsData = map[string][]fyne.CanvasObject{
	"First tab":  {widget.NewLabel("a1"), widget.NewLabel("b1"), widget.NewLabel("c1")},
	"Second tab": {widget.NewLabel("a2"), widget.NewLabel("b2"), widget.NewLabel("c2")},
	"Third tab":  {widget.NewLabel("a3"), widget.NewLabel("b3"), widget.NewLabel("c3")},
}

func createTabContainerForTesting() *TabContainer {
	return NewTabContainer(tabsData, boldLabelOnSelection, unboldLabelOnUnselection, )
}

func boldLabelOnSelection(element *fyne.CanvasObject) {
	label := *element
	label.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
}

func unboldLabelOnUnselection(element *fyne.CanvasObject) {
	label := *element
	label.(*widget.Label).TextStyle = fyne.TextStyle{Bold: false}
}

func TestThatOnDisplayCorrectTabAndContentDisplays(t *testing.T) {
	container := createTabContainerForTesting()
	assert.Equal(t, "First tab", container.CurrentTab().Text)
	label1 := GetLabelFromContent(container.CurrentTab().Content, "a1")
	assert.NotNil(t, label1)
	label2 := GetLabelFromContent(container.CurrentTab().Content, "b1")
	assert.NotNil(t, label2)
	label3 := GetLabelFromContent(container.CurrentTab().Content, "c1")
	assert.NotNil(t, label3)
}

func TestThatChangingGraphicalPropertiesOnSelectionUnselectionWorks(t *testing.T) {
	container := createTabContainerForTesting()
	label1 := GetLabelFromContent(container.CurrentTab().Content, "a1")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label1.TextStyle)
	label2 := GetLabelFromContent(container.CurrentTab().Content, "b1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label2.TextStyle)
	label3 := GetLabelFromContent(container.CurrentTab().Content, "c1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label3.TextStyle)
}

func TestThatSelectingTabsWorks(t *testing.T) {
	container := createTabContainerForTesting()
	container.SelectNextTab()
	assert.Equal(t, "Second tab", container.CurrentTab().Text)
	container.SelectNextTab()
	assert.Equal(t, "Third tab", container.CurrentTab().Text)
	container.SelectNextTab()
	assert.Equal(t, "First tab", container.CurrentTab().Text)
	container.SelectPreviousTab()
	assert.Equal(t, "Third tab", container.CurrentTab().Text)
	container.SelectPreviousTab()
	assert.Equal(t, "Second tab", container.CurrentTab().Text)
	container.SelectPreviousTab()
	assert.Equal(t, "First tab", container.CurrentTab().Text)
}

func TestThatChangingTabSelectsFirstLabelOnTab(t *testing.T) {
	container := createTabContainerForTesting()
	container.SelectNextTab()
	label1 := GetLabelFromContent(container.CurrentTab().Content, "a2")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label1.TextStyle)
	label2 := GetLabelFromContent(container.CurrentTab().Content, "b2")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label2.TextStyle)
	label3 := GetLabelFromContent(container.CurrentTab().Content, "c2")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label3.TextStyle)
	container.SelectPreviousTab()
	label4 := GetLabelFromContent(container.CurrentTab().Content, "a1")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label4.TextStyle)
	label5 := GetLabelFromContent(container.CurrentTab().Content, "b1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label5.TextStyle)
	label6 := GetLabelFromContent(container.CurrentTab().Content, "c1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label6.TextStyle)
}

func TestThatChangingIntoTabWithNoContentDoesNotPanic(t *testing.T) {
	container := NewTabContainer(
		map[string][]fyne.CanvasObject{
			"First tab":  {widget.NewLabel("a1"), widget.NewLabel("b1"), widget.NewLabel("c1")},
			"Second tab": {},
		},
		boldLabelOnSelection,
		unboldLabelOnUnselection,
	)
	container.SelectNextTab()
}

func TestThatThereIsNoPanicWhenCheckingIfTabHasElementsOnContainerWithoutElements(t *testing.T) {
	container := NewTabContainer(
		make(map[string][]fyne.CanvasObject),
		func(element *fyne.CanvasObject) {},
		func(element *fyne.CanvasObject) {})
	container.currentTabHasElements()
}
