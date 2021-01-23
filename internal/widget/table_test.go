package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const testLabelWidth = 20
const testLabelHeight = 10

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
		posX += 135
	}
}

func TestThatObjectsInHeaderHaveCorrectSize(t *testing.T) {
	table := NewTable(14, createLabelsForTesting(14), []fyne.CanvasObject{})
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	for i, object := range table.headerObjects {
		assert.Equal(t, 100, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, 50, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectPositions(t *testing.T) {
	const columnAmount = 14
	const rowAmount = 10
	table := NewTable(columnAmount, []fyne.CanvasObject{}, createLabelsForTesting(columnAmount*rowAmount))
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	posX := 0
	posY := 50
	for i, object := range table.objects {
		assert.Equal(t, posX, object.Position().X, "Position x of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, object.Position().Y, "Position y of object num "+strconv.Itoa(i)+" is incorrect")
		posX += 135
		if i != 0 && (i+1)%columnAmount == 0 {
			posX = 0
			posY += 141
		}
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectSize(t *testing.T) {
	const columnAmount = 14
	const rowAmount = 10
	table := NewTable(columnAmount, []fyne.CanvasObject{}, createLabelsForTesting(columnAmount*rowAmount))
	table.CreateRenderer().Layout(fyne.NewSize(1000, 1000))
	for i, object := range table.objects {
		assert.Equal(t, 100, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, 141, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}
