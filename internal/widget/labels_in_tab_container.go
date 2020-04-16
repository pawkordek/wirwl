package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

/*
Tab container which contains labels as elements stored in every tab.
It will mark selected label using bolded text.
*/
type LabelsInTabContainer struct {
	widget.TabContainer
}

func NewLabelsInTabContainer(tabsData map[string][]string) *TabContainer {
	labelsData := getTabsDataAsLabelsMap(tabsData)
	container := NewTabContainer(labelsData, boldSelectedLabel, unboldSelectedLabel)
	return container
}

func getTabsDataAsLabelsMap(tabsData map[string][]string) map[string][]fyne.CanvasObject {
	labelsMap := make(map[string][]fyne.CanvasObject)
	for tabName, labelsNames := range tabsData {
		labels := getLabelsWithNames(labelsNames)
		labelsMap[tabName] = labels
	}
	return labelsMap
}

func getLabelsWithNames(names []string) []fyne.CanvasObject {
	labels := make([]fyne.CanvasObject, 0, len(names))
	for _, name := range names {
		labels = append(labels, widget.NewLabel(name))
	}
	return labels
}

func boldSelectedLabel(element *fyne.CanvasObject) {
	label := *element
	label.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
}
func unboldSelectedLabel(element *fyne.CanvasObject) {
	label := *element
	label.(*widget.Label).TextStyle = fyne.TextStyle{Bold: false}
}
