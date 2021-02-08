package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"image/color"
	"strconv"
	"testing"
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

func TestThatTableHasCorrectMinSize(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	assert.Equal(t, testColumnAmount*expectedColumnWidthWithPadding, table.MinSize().Width, "Table has incorrect minimum width")
	assert.Equal(t, testRowAmount*expectedRowHeight+expectedHeaderHeight, table.MinSize().Height, "Table has incorrect minimum height")
}

func TestThatObjectsInHeaderHaveCorrectPositions(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	posX := 0
	posY := 0
	for i, object := range table.headerObjects {
		assert.Equal(t, posX, object.Position().X, "Position x of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, object.Position().Y, "Position y of object num "+strconv.Itoa(i)+" is incorrect")
		posX += expectedColumnWidthWithPadding
	}
}

func TestThatObjectsInHeaderHaveCorrectSize(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	for i, object := range table.headerObjects {
		assert.Equal(t, expectedColumnWidth, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, expectedHeaderHeight, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectPositions(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	posX := 0
	posY := expectedHeaderHeight
	for i, object := range table.objects {
		assert.Equal(t, posX, object.Position().X, "Position x of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, object.Position().Y, "Position y of object num "+strconv.Itoa(i)+" is incorrect")
		posX += expectedColumnWidthWithPadding
		if i != 0 && (i+1)%testColumnAmount == 0 {
			posX = 0
			posY += expectedRowHeight
		}
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectSize(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	for i, object := range table.objects {
		assert.Equal(t, expectedColumnWidth, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, expectedRowHeight, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatHeaderRowBorderIsDrawnCorrectly(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	rectangle := renderer.headerRowBorder
	assert.Equal(t, testColumnAmount*expectedColumnWidthWithPadding, rectangle.Size().Width)
	assert.Equal(t, expectedHeaderHeight, rectangle.Size().Height)
	assert.Equal(t, float32(2), rectangle.StrokeWidth)
	assert.Equal(t, color.Black, rectangle.StrokeColor)
	assert.Equal(t, color.Transparent, rectangle.FillColor)
}

func TestThatThereIsCorrectAmountOfDataRowBorders(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	assert.Equal(t, testRowAmount, len(renderer.dataRowsBorders))
}

func TestThatAllDataRowBordersHaveCorrectSize(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	for i, border := range renderer.dataRowsBorders {
		assert.Equal(t, testColumnAmount*expectedColumnWidthWithPadding, border.Size().Width, "Border with number "+strconv.Itoa(i)+" does not have the correct width")
		assert.Equal(t, expectedRowHeight, border.Size().Height, "Border with number "+strconv.Itoa(i)+" does not have the correct height")
	}
}

func TestThatAllDataBordersHaveCorrectPosition(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	expectedPosition := fyne.NewPos(0, expectedHeaderHeight)
	for i, border := range renderer.dataRowsBorders {
		assert.Equal(t, expectedPosition, border.Position(), "Border with number "+strconv.Itoa(i)+" does not have correct position")
		expectedPosition = expectedPosition.Add(fyne.NewPos(0, expectedRowHeight))
	}
}

func TestThatAllDataBordersHaveCorrectProperties(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	for i, border := range renderer.dataRowsBorders {
		assert.Equal(t, float32(2), border.StrokeWidth, "Border with number "+strconv.Itoa(i)+" does not have correct stroke width")
		assert.Equal(t, color.Black, border.StrokeColor, "Border with number "+strconv.Itoa(i)+" does not have correct stroke color")
		assert.Equal(t, color.Transparent, border.FillColor, "Border with number "+strconv.Itoa(i)+" does not have correct fill color")
	}
}

func TestThatThereIsCorrectAmountOfColumnBorders(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	assert.Equal(t, testColumnAmount, len(renderer.columnBorders))
}

func TestThatAllColumnBordersHaveCorrectSize(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	for i, border := range renderer.columnBorders {
		assert.Equal(t, expectedColumnWidthWithPadding, border.Size().Width, "Border with number "+strconv.Itoa(i)+" does not have the correct width")
		assert.Equal(t, expectedHeaderHeight+expectedRowHeight*testRowAmount, border.Size().Height, "Border with number "+strconv.Itoa(i)+" does not have the correct height")
	}
}

func TestThatAllColumnBordersHaveCorrectPosition(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	expectedPosition := fyne.NewPos(0, 0)
	for i, border := range renderer.columnBorders {
		assert.Equal(t, expectedPosition, border.Position(), "Border with number "+strconv.Itoa(i)+" does not have correct position")
		expectedPosition = expectedPosition.Add(fyne.NewPos(expectedColumnWidthWithPadding, 0))
	}
}

func TestThatAllColumnBordersHaveCorrectProperties(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	for i, border := range renderer.columnBorders {
		assert.Equal(t, float32(2), border.StrokeWidth, "Border with number "+strconv.Itoa(i)+" does not have correct stroke width")
		assert.Equal(t, color.Black, border.StrokeColor, "Border with number "+strconv.Itoa(i)+" does not have correct stroke color")
		assert.Equal(t, color.Transparent, border.FillColor, "Border with number "+strconv.Itoa(i)+" does not have correct fill color")
	}
}
