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

func createLabelsForTesting(amountOfLabels int) []fyne.CanvasObject {
	labels := []fyne.CanvasObject{}
	for i := 1; i <= amountOfLabels; i++ {
		label := widget.NewLabel("Test label num " + strconv.Itoa(i))
		label.Resize(fyne.NewSize(testLabelWidth, testLabelHeight))
		labels = append(labels, label)
	}
	return labels
}

func TestThatTableHasCorrectMinSize(t *testing.T) {
	const columnAmount = 14
	const rowAmount = 100
	table := NewTable(columnAmount, []fyne.CanvasObject{}, createLabelsForTesting(columnAmount*rowAmount))
	minSize := table.MinSize()
	assert.Equal(t, columnAmount*testLabelWidth, minSize.Width, "Table has incorrect minimum width")
	assert.Equal(t, rowAmount*testLabelHeight, minSize.Height, "Table has incorrect minimum height")
}

func TestThatObjectsInHeaderHaveCorrectPositions(t *testing.T) {
	table := NewTable(14, createLabelsForTesting(14), []fyne.CanvasObject{})
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	posX := 0
	posY := 0
	for i, object := range table.headerObjects {
		assert.Equal(t, posX, object.Position().X, "Position x of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, object.Position().Y, "Position y of object num "+strconv.Itoa(i)+" is incorrect")
		posX += expectedColumnWidthWithPadding
	}
}

func TestThatObjectsInHeaderHaveCorrectSize(t *testing.T) {
	table := NewTable(14, createLabelsForTesting(14), []fyne.CanvasObject{})
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	for i, object := range table.headerObjects {
		assert.Equal(t, expectedColumnWidth, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, expectedHeaderHeight, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectPositions(t *testing.T) {
	const columnAmount = 14
	const rowAmount = 10
	table := NewTable(columnAmount, []fyne.CanvasObject{}, createLabelsForTesting(columnAmount*rowAmount))
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	posX := 0
	posY := expectedHeaderHeight
	for i, object := range table.objects {
		assert.Equal(t, posX, object.Position().X, "Position x of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, object.Position().Y, "Position y of object num "+strconv.Itoa(i)+" is incorrect")
		posX += expectedColumnWidthWithPadding
		if i != 0 && (i+1)%columnAmount == 0 {
			posX = 0
			posY += expectedRowHeight
		}
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectSize(t *testing.T) {
	const columnAmount = 14
	const rowAmount = 10
	table := NewTable(columnAmount, []fyne.CanvasObject{}, createLabelsForTesting(columnAmount*rowAmount))
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	for i, object := range table.objects {
		assert.Equal(t, expectedColumnWidth, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, expectedRowHeight, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatHeaderRowBorderIsDrawnCorrectly(t *testing.T) {
	const columnAmount = 14
	table := NewTable(columnAmount, createLabelsForTesting(14), []fyne.CanvasObject{})
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	rectangle := renderer.headerRowBorder
	assert.Equal(t, columnAmount*expectedColumnWidthWithPadding, rectangle.Size().Width)
	assert.Equal(t, expectedHeaderHeight, rectangle.Size().Height)
	assert.Equal(t, float32(2), rectangle.StrokeWidth)
	assert.Equal(t, color.Black, rectangle.StrokeColor)
	assert.Equal(t, color.Transparent, rectangle.FillColor)
}

func TestThatThereIsCorrectAmountOfDataRowBorders(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	assert.Equal(t, amountOfRows, len(renderer.dataRowsBorders))
}

func TestThatAllDataRowBordersHaveCorrectSize(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	for i, border := range renderer.dataRowsBorders {
		assert.Equal(t, columnAmount*expectedColumnWidthWithPadding, border.Size().Width, "Border with number "+strconv.Itoa(i)+" does not have the correct width")
		assert.Equal(t, expectedRowHeight, border.Size().Height, "Border with number "+strconv.Itoa(i)+" does not have the correct height")
	}
}

func TestThatAllDataBordersHaveCorrectPosition(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	expectedPosition := fyne.NewPos(0, expectedHeaderHeight)
	for i, border := range renderer.dataRowsBorders {
		assert.Equal(t, expectedPosition, border.Position(), "Border with number "+strconv.Itoa(i)+" does not have correct position")
		expectedPosition = expectedPosition.Add(fyne.NewPos(0, expectedRowHeight))
	}
}

func TestThatAllDataBordersHaveCorrectProperties(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	for i, border := range renderer.dataRowsBorders {
		assert.Equal(t, float32(2), border.StrokeWidth, "Border with number "+strconv.Itoa(i)+" does not have correct stroke width")
		assert.Equal(t, color.Black, border.StrokeColor, "Border with number "+strconv.Itoa(i)+" does not have correct stroke color")
		assert.Equal(t, color.Transparent, border.FillColor, "Border with number "+strconv.Itoa(i)+" does not have correct fill color")
	}
}

func TestThatThereIsCorrectAmountOfColumnBorders(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	assert.Equal(t, columnAmount, len(renderer.columnBorders))
}

func TestThatAllColumnBordersHaveCorrectSize(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	for i, border := range renderer.columnBorders {
		assert.Equal(t, expectedColumnWidthWithPadding, border.Size().Width, "Border with number "+strconv.Itoa(i)+" does not have the correct width")
		assert.Equal(t, expectedHeaderHeight+expectedRowHeight*amountOfRows, border.Size().Height, "Border with number "+strconv.Itoa(i)+" does not have the correct height")
	}
}

func TestThatAllColumnBordersHaveCorrectPosition(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	expectedPosition := fyne.NewPos(0, 0)
	for i, border := range renderer.columnBorders {
		assert.Equal(t, expectedPosition, border.Position(), "Border with number "+strconv.Itoa(i)+" does not have correct position")
		expectedPosition = expectedPosition.Add(fyne.NewPos(expectedColumnWidthWithPadding, 0))
	}
}

func TestThatAllColumnBordersHaveCorrectProperties(t *testing.T) {
	const columnAmount = 14
	const amountOfRows = 20
	table := NewTable(columnAmount, createLabelsForTesting(columnAmount), createLabelsForTesting(columnAmount*amountOfRows))
	renderer := table.CreateRenderer().(tableRenderer)
	renderer.Layout(fyne.NewSize(1000, 1000))
	for i, border := range renderer.columnBorders {
		assert.Equal(t, float32(2), border.StrokeWidth, "Border with number "+strconv.Itoa(i)+" does not have correct stroke width")
		assert.Equal(t, color.Black, border.StrokeColor, "Border with number "+strconv.Itoa(i)+" does not have correct stroke color")
		assert.Equal(t, color.Transparent, border.FillColor, "Border with number "+strconv.Itoa(i)+" does not have correct fill color")
	}
}
