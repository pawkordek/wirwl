package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"strconv"
)

const testLabelWidth = 20
const testLabelHeight = 10
const expectedHeaderHeight = 50
const expectedColumnWidth = 100
const expectedColumnWidthWithPadding = 135
const expectedRowHeight = 141

const testColumnAmount = 14
const testRowAmount = 20

func createLabelsForTesting(amountOfLabels int) []fyne.CanvasObject {
	labels := []fyne.CanvasObject{}
	for i := 1; i <= amountOfLabels; i++ {
		label := widget.NewLabel("Test label num " + strconv.Itoa(i))
		label.Resize(fyne.NewSize(testLabelWidth, testLabelHeight))
		labels = append(labels, label)
	}
	return labels
}

func createTableRendererForTesting(tableColumnAmount int, tableRowAmount int) tableRenderer {
	table := NewTable(testColumnAmount, createLabelsForTesting(testColumnAmount), createLabelsForTesting(testColumnAmount*testRowAmount))
	renderer := table.CreateRenderer().(tableRenderer)
	//The size is arbitrary but shouldn't be zero as layout with zero size doesn't make any sense
	renderer.Layout(fyne.NewSize(1000, 1000))
	return renderer
}

func createTableForTesting(columnAmount int, rowAmount int) Table {
	renderer := createTableRendererForTesting(columnAmount, rowAmount)
	return renderer.table
}
