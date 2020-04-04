package widget

import (
	"fyne.io/fyne"
	fyneWidget "fyne.io/fyne/widget"
	"sort"
)

/*
Tab container in which every tab contains a list of labels.
It allows to switch tab to next/previous which is done cyclically, setting next tab when on last tab goes to the first and vice versa.
*/
type TabContainer struct {
	/*Cannot extend fyne's TabContainer right now as there is a bug that prevents tab buttons from updating when tabs change on extended TabContainer.
	  So the only way to do it right now is to to compose fyne's TabContainer and extend a Box which will contain this TabContainer.
	  The bug in question is this: https://github.com/fyne-io/fyne/issues/810
	  TODO: Remove the workaround when fyne 1.2.4 comes out and replace pointers with variables
	*/
	it *fyneWidget.TabContainer
	*fyneWidget.Box
	selectedLabelIndex int
	labels             map[string][]*fyneWidget.Label
}

func NewTabContainer(tabsData map[string][]string) *TabContainer {
	var tabs []*fyneWidget.TabItem
	allLabels := map[string][]*fyneWidget.Label{}
	sortedTabsNames := getAlphabeticallySortedTabsNames(tabsData)
	for _, tabName := range sortedTabsNames {
		labels := createLabels(tabsData[tabName])
		allLabels[tabName] = labels
		labelsAsCanvasObjects := getLabelsAsCanvasObjects(labels)
		formItem := fyneWidget.NewTabItem(tabName, fyneWidget.NewVBox(labelsAsCanvasObjects...))
		tabs = append(tabs, formItem)
	}
	fyneContainer := fyneWidget.NewTabContainer(tabs...)
	container := &TabContainer{
		Box:                fyneWidget.NewVBox(fyneContainer),
		it:                 fyneContainer,
		selectedLabelIndex: 0,
		labels:             allLabels,
	}
	container.selectLabel(0)
	container.it.SelectTabIndex(0)
	return container
}

func getAlphabeticallySortedTabsNames(tabsData map[string][]string) []string {
	sortedNames := make([]string, 0, len(tabsData))
	for tabName, _ := range tabsData {
		sortedNames = append(sortedNames, tabName)
	}
	sort.Strings(sortedNames)
	return sortedNames
}

func createLabels(names []string) []*fyneWidget.Label {
	var labels []*fyneWidget.Label
	for _, name := range names {
		labels = append(labels, fyneWidget.NewLabel(name))
	}
	return labels
}

func getLabelsAsCanvasObjects(labels []*fyneWidget.Label) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(labels))
	for i, _ := range labels {
		objects[i] = labels[i]
	}
	return objects
}

func (container *TabContainer) selectLabel(num int) {
	container.getSelectedLabel().TextStyle = fyne.TextStyle{Bold: false}
	container.getSelectedLabel().Refresh()
	container.selectedLabelIndex = num
	container.getSelectedLabel().TextStyle = fyne.TextStyle{Bold: true}
	container.getSelectedLabel().Refresh()
}

func (container *TabContainer) getSelectedLabel() *fyneWidget.Label {
	currentTabText := container.it.CurrentTab().Text
	return container.labels[currentTabText][container.selectedLabelIndex]
}

func (container *TabContainer) SelectNextTab() {
	currentTabIndex := container.it.CurrentTabIndex()
	if currentTabIndex >= len(container.it.Items)-1 {
		container.setTabTo(0)
	} else {
		container.setTabTo(currentTabIndex + 1)
	}
}

func (container *TabContainer) SelectPreviousTab() {
	currentTabIndex := container.it.CurrentTabIndex()
	if currentTabIndex == 0 {
		container.setTabTo(len(container.it.Items) - 1)
	} else {
		container.setTabTo(currentTabIndex - 1)
	}
}

func (container *TabContainer) setTabTo(index int) {
	container.it.SelectTabIndex(index)
	container.selectLabel(0)
}
