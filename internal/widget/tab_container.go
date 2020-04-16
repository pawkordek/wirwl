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
	fyneWidget.TabContainer
	selectedElementIndex int
	tabsContent          map[string][]fyne.CanvasObject
	onElementSelected    func(element *fyne.CanvasObject)
	onElementUnselected  func(element *fyne.CanvasObject)
}

func NewTabContainer(
	tabsData map[string][]fyne.CanvasObject,
	onElementSelected func(element *fyne.CanvasObject),
	onElementUnselected func(element *fyne.CanvasObject)) *TabContainer {
	container := &TabContainer{
		selectedElementIndex: 0,
		tabsContent:          tabsData,
		onElementSelected:    onElementSelected,
		onElementUnselected:  onElementUnselected,
	}
	container.ExtendBaseWidget(container)
	container.TabContainer.Items = getTabsFromData(tabsData)
	container.selectElement(0)
	container.SelectTabIndex(0)
	return container
}

func getTabsFromData(tabsData map[string][]fyne.CanvasObject) []*fyneWidget.TabItem {
	var tabs []*fyneWidget.TabItem
	sortedTabsNames := getAlphabeticallySortedTabsNames(tabsData)
	for _, tabName := range sortedTabsNames {
		formItem := fyneWidget.NewTabItem(tabName, fyneWidget.NewVBox(tabsData[tabName]...))
		tabs = append(tabs, formItem)
	}
	return tabs
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
	if container.currentTabHasElements() {
		selectedElement := container.selectedElement()
		container.onElementUnselected(&selectedElement)
		selectedElement.Refresh()
		container.selectedElementIndex = num
		selectedElement = container.selectedElement()
		container.onElementSelected(&selectedElement)
		container.selectedElement().Refresh()
	}
}

func (container *TabContainer) currentTabHasElements() bool {
	currentTabText := container.CurrentTab().Text
	return len(container.tabsContent[currentTabText]) > 0
}

func (container *TabContainer) selectedElement() fyne.CanvasObject {
	currentTabText := container.CurrentTab().Text
	return container.tabsContent[currentTabText][container.selectedElementIndex]
}

func (container *TabContainer) SelectNextTab() {
	currentTabIndex := container.CurrentTabIndex()
	if currentTabIndex >= len(container.Items())-1 {
		container.setTabTo(0)
	} else {
		container.setTabTo(currentTabIndex + 1)
	}
}

func (container *TabContainer) SelectPreviousTab() {
	currentTabIndex := container.CurrentTabIndex()
	if currentTabIndex == 0 {
		container.setTabTo(len(container.Items()) - 1)
	} else {
		container.setTabTo(currentTabIndex - 1)
	}
}

func (container *TabContainer) setTabTo(index int) {
	container.SelectTabIndex(index)
	container.selectElement(0)
}

func (container *TabContainer) Items() []*fyneWidget.TabItem {
	return container.TabContainer.Items
}
