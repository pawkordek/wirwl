package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
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
	expectedPadding                = 35
	expectedColumnWidthWithPadding = expectedColumnWidth + expectedPadding
	expectedRowHeight              = 141
	expectedTableWidth             = 741
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

func createTableRendererForTesting(table Table) tableRenderer {
	renderer := table.CreateRenderer().(tableRenderer)
	//The size is arbitrary but shouldn't be zero as layout with zero size doesn't make any sense
	renderer.Layout(fyne.NewSize(1000, 1000))
	return renderer
}

func createTableForTesting(canvas fyne.Canvas, columnAmount int, rowAmount int) Table {
	table := NewTable(canvas, createColumnDataForTesting(testColumnAmount), createLabelsForTesting(testColumnAmount, testRowAmount))
	renderer := createTableRendererForTesting(*table)
	return renderer.table
}

func createDefaultTableForTesting() Table {
	return createTableForTesting(test.Canvas(), testColumnAmount, testRowAmount)
}

func createDefaultTableForTestingWithCustomCanvas(canvas fyne.Canvas) Table {
	return createTableForTesting(canvas, testColumnAmount, testRowAmount)
}
