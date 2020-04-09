package widget

import (
	"fyne.io/fyne"
	fyneWidget "fyne.io/fyne/widget"
	"sort"
)

/*
Tab container in which every tab contains a list of CanvasObjects. Tabs are displayed alphabetically.
It allows to switch tab to next/previous which is done cyclically, setting next tab when on last tab goes to the first and vice versa.
The way selected items display graphically should be controlled by onElementSelected and onElementUnselected functions e.g. labels becoming bold on selection.
*/
type TabContainer struct {
	/*Cannot extend fyne's TabContainer right now as there is a bug that prevents tab buttons from updating when tabs change on extended TabContainer.
	  So the only way to do it right now is to to compose fyne's TabContainer and extend a Box which will contain this TabContainer.
	  The bug in question is this: https://github.com/fyne-io/fyne/issues/810
	  TODO: Remove the workaround when fyne 1.2.4 comes out and replace pointers with variables
	*/
	it *fyneWidget.TabContainer
	*fyneWidget.Box
	selectedElementIndex int
	tabsContent          map[string][]fyne.CanvasObject
	onElementSelected    func(element *fyne.CanvasObject)
	onElementUnselected  func(element *fyne.CanvasObject)
}

func NewTabContainer(
	tabsData map[string][]fyne.CanvasObject,
	onElementSelected func(element *fyne.CanvasObject),
	onElementUnselected func(element *fyne.CanvasObject)) *TabContainer {
	var tabs []*fyneWidget.TabItem
	sortedTabsNames := getAlphabeticallySortedTabsNames(tabsData)
	for _, tabName := range sortedTabsNames {
		formItem := fyneWidget.NewTabItem(tabName, fyneWidget.NewVBox(tabsData[tabName]...))
		tabs = append(tabs, formItem)
	}
	fyneContainer := fyneWidget.NewTabContainer(tabs...)
	container := &TabContainer{
		Box:                  fyneWidget.NewVBox(fyneContainer),
		it:                   fyneContainer,
		selectedElementIndex: 0,
		tabsContent:          tabsData,
		onElementSelected:    onElementSelected,
		onElementUnselected:  onElementUnselected,
	}
	container.selectElement(0)
	container.it.SelectTabIndex(0)
	return container
}

func getAlphabeticallySortedTabsNames(tabsData map[string][]fyne.CanvasObject) []string {
	sortedNames := make([]string, 0, len(tabsData))
	for tabName, _ := range tabsData {
		sortedNames = append(sortedNames, tabName)
	}
	sort.Strings(sortedNames)
	return sortedNames
}

func (container *TabContainer) selectElement(num int) {
	selectedElement := container.selectedElement()
	container.onElementUnselected(&selectedElement)
	selectedElement.Refresh()
	container.selectedElementIndex = num
	selectedElement = container.selectedElement()
	container.onElementSelected(&selectedElement)
	container.selectedElement().Refresh()
}

func (container *TabContainer) selectedElement() fyne.CanvasObject {
	currentTabText := container.it.CurrentTab().Text
	return container.tabsContent[currentTabText][container.selectedElementIndex]
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
	container.selectElement(0)
}
