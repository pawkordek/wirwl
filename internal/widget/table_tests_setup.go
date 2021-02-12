package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"strconv"
)

const (
	testLabelWidth   = 20
	testLabelHeight  = 10
	testColumnAmount = 14
	testRowAmount    = 20
)

const (
	expectedHeaderHeight           = 50
	expectedColumnWidth            = 100
	expectedColumnWidthWithPadding = 135
	expectedRowHeight              = 141
)

func createColumnDataForTesting(amountOfColumns int) []TableColumn {
	data := []TableColumn{}
	for i := 1; i <= amountOfColumns; i++ {
		column := TableColumn{Type: TextColumn, Name: string(i)}
		data = append(data, column)
	}
	return data
}

func createLabelsForTesting(amountOfColumns int, amountOfRows int) []TableRow {
	labels := []TableRow{}
	for j := 1; j <= amountOfRows; j++ {
		row := TableRow{}
		for i := 1; i <= amountOfColumns; i++ {
			label := widget.NewLabel("Test label num " + strconv.Itoa(i))
			label.Resize(fyne.NewSize(testLabelWidth, testLabelHeight))
			row = append(row, label)
		}
		labels = append(labels, row)
	}
	return labels
}

func createTableRendererForTesting(tableColumnAmount int, tableRowAmount int) tableRenderer {
	table := NewTable(testColumnAmount, createColumnDataForTesting(testColumnAmount), createLabelsForTesting(testColumnAmount, testRowAmount))
	renderer := table.CreateRenderer().(tableRenderer)
	//The size is arbitrary but shouldn't be zero as layout with zero size doesn't make any sense
	renderer.Layout(fyne.NewSize(1000, 1000))
	return renderer
}

func createTableForTesting(columnAmount int, rowAmount int) Table {
	renderer := createTableRendererForTesting(columnAmount, rowAmount)
	return renderer.table
}
