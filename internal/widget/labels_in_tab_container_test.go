package widget

import (
	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createLabelsInTabContainerForTesting() *TabContainer {
	return NewLabelsInTabContainer(
		map[string][]string{
			"First tab":  {"a1", "b1", "c1"},
			"Second tab": {"a2", "b2", "c3"},
			"Third tab":  {"a3", "b2", "c3"},
		})
}

func TestThatLabelsDisplay(t *testing.T) {
	container := createLabelsInTabContainerForTesting()
	assert.Equal(t, "First tab", container.CurrentTab().Text)
	label1 := GetLabelFromContent(container.CurrentTab().Content, "a1")
	assert.NotNil(t, label1)
	label2 := GetLabelFromContent(container.CurrentTab().Content, "b1")
	assert.NotNil(t, label2)
	label3 := GetLabelFromContent(container.CurrentTab().Content, "c1")
	assert.NotNil(t, label3)
}

func TestThatLabelsAreBoldedUnbolded(t *testing.T) {
	container := createLabelsInTabContainerForTesting()
	label1 := GetLabelFromContent(container.CurrentTab().Content, "a1")
	assert.Equal(t, fyne.TextStyle{Bold: true}, label1.TextStyle)
	label2 := GetLabelFromContent(container.CurrentTab().Content, "b1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label2.TextStyle)
	label3 := GetLabelFromContent(container.CurrentTab().Content, "c1")
	assert.Equal(t, fyne.TextStyle{Bold: false}, label3.TextStyle)
}
