package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"image/color"
	"strconv"
	"testing"
)

func TestThatTableHasCorrectMinSize(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	assert.Equal(t, testColumnAmount*expectedColumnWidthWithPadding, table.MinSize().Width, "Table has incorrect minimum width")
	assert.Equal(t, testRowAmount*expectedRowHeight+expectedHeaderHeight, table.MinSize().Height, "Table has incorrect minimum height")
}

func TestThatObjectsInHeaderHaveCorrectPositions(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	posX := expectedPadding / 2
	posY := 0
	for i, columnLabel := range table.columnLabels {
		assert.Equal(t, posX, columnLabel.Position().X, "Position x of columnLabel num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, columnLabel.Position().Y, "Position y of columnLabel num "+strconv.Itoa(i)+" is incorrect")
		posX += columnLabel.Size().Width + expectedPadding
	}
}

func TestThatObjectsInHeaderHaveCorrectSize(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	for i, object := range table.columnLabels {
		assert.Equal(t, object.MinSize().Width, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, expectedHeaderHeight, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatColumnLabelsAreBolded(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	for i, object := range table.columnLabels {
		label := object.(*widget.Label)
		assert.Equal(t, true, label.TextStyle.Bold, "Column label num "+strconv.Itoa(i)+" is not bolded")
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectPositions(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	posX := 0
	posY := expectedHeaderHeight
	for _, row := range table.rowData {
		for i, cell := range row {
			assert.Equal(t, posX, cell.Position().X, "Position x of cell num "+strconv.Itoa(i)+" is incorrect")
			assert.Equal(t, posY, cell.Position().Y, "Position y of cell num "+strconv.Itoa(i)+" is incorrect")
			posX += expectedColumnWidthWithPadding
			if i != 0 && (i+1)%testColumnAmount == 0 {
				posX = 0
				posY += expectedRowHeight
			}
		}
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectSize(t *testing.T) {
	table := createTableForTesting(testColumnAmount, testRowAmount)
	for _, row := range table.rowData {
		for i, cell := range row {
			assert.Equal(t, expectedColumnWidth, cell.Size().Width, "Width of cell num "+strconv.Itoa(i)+" is incorrect")
			assert.Equal(t, expectedRowHeight, cell.Size().Height, "Height of cell num "+strconv.Itoa(i)+" is incorrect")
		}
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
		assert.Equal(t, renderer.table.columnLabels[i].Size().Width+expectedPadding, border.Size().Width, "Border with number "+strconv.Itoa(i)+" does not have the correct width")
		assert.Equal(t, expectedHeaderHeight+expectedRowHeight*testRowAmount, border.Size().Height, "Border with number "+strconv.Itoa(i)+" does not have the correct height")
	}
}

func TestThatAllColumnBordersHaveCorrectPosition(t *testing.T) {
	renderer := createTableRendererForTesting(testColumnAmount, testRowAmount)
	expectedPosition := fyne.NewPos(0, 0)
	for i, border := range renderer.columnBorders {
		assert.Equal(t, expectedPosition, border.Position(), "Border with number "+strconv.Itoa(i)+" does not have correct position")
		expectedPosition = expectedPosition.Add(fyne.NewPos(renderer.table.columnLabels[i].Size().Width+expectedPadding, 0))
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
