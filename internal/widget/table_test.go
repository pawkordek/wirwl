package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestThatTableHasCorrectMinSize(t *testing.T) {
	table := createDefaultTableForTesting()
	assert.Equal(t, expectedTableWidth, table.MinSize().Width, "Table has incorrect minimum width")
	assert.Equal(t, testRowAmount*expectedRowHeight+expectedHeaderHeight, table.MinSize().Height, "Table has incorrect minimum height")
}

func TestThatObjectsInHeaderHaveCorrectPositions(t *testing.T) {
	table := createDefaultTableForTesting()
	posX := expectedPadding / 2
	posY := 0
	for i, columnLabel := range table.columnLabels {
		assert.Equal(t, posX, columnLabel.Position().X, "Position x of columnLabel num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, posY, columnLabel.Position().Y, "Position y of columnLabel num "+strconv.Itoa(i)+" is incorrect")
		posX += columnLabel.Size().Width + expectedPadding
	}
}

func TestThatObjectsInHeaderHaveCorrectSize(t *testing.T) {
	table := createDefaultTableForTesting()
	for i, object := range table.columnLabels {
		assert.Equal(t, object.MinSize().Width, object.Size().Width, "Width of object num "+strconv.Itoa(i)+" is incorrect")
		assert.Equal(t, expectedHeaderHeight, object.Size().Height, "Height of object num "+strconv.Itoa(i)+" is incorrect")
	}
}

func TestThatColumnLabelsAreBolded(t *testing.T) {
	table := createDefaultTableForTesting()
	for i, object := range table.columnLabels {
		label := object.(*widget.Label)
		assert.Equal(t, true, label.TextStyle.Bold, "Column label num "+strconv.Itoa(i)+" is not bolded")
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectPositions(t *testing.T) {
	table := createDefaultTableForTesting()
	posX := expectedPadding / 2
	posY := expectedHeaderHeight
	for _, row := range table.rowData {
		for i, cell := range row {
			assert.Equal(t, posX, cell.Position().X, "Position x of cell num "+strconv.Itoa(i)+" is incorrect")
			assert.Equal(t, posY, cell.Position().Y, "Position y of cell num "+strconv.Itoa(i)+" is incorrect")
			posX += table.columnLabels[i].Size().Width + expectedPadding
			if i != 0 && (i+1)%testColumnAmount == 0 {
				posX = expectedPadding / 2
				posY += expectedRowHeight
			}
		}
	}
}

func TestThatObjectsThatCreateDataRowsHaveCorrectSize(t *testing.T) {
	table := createDefaultTableForTesting()
	for _, row := range table.rowData {
		for i, cell := range row {
			assert.Equal(t, table.columnLabels[i].Size().Width, cell.Size().Width, "Width of cell num "+strconv.Itoa(i)+" is incorrect")
			assert.Equal(t, expectedRowHeight, cell.Size().Height, "Height of cell num "+strconv.Itoa(i)+" is incorrect")
		}
	}
}

func TestThatTableGetsFocusWhenEnteringInputMode(t *testing.T) {
	testWindow := test.NewApp().NewWindow("")
	table := createDefaultTableForTestingWithCustomCanvas(testWindow.Canvas())
	testWindow.SetContent(table)
	table.EnterInputMode()
	assert.Equal(t, table, testWindow.Canvas().Focused())
}

func TestThatTableIsNotFocusedAfterExitingInputMode(t *testing.T) {
	testWindow := test.NewApp().NewWindow("")
	table := createDefaultTableForTestingWithCustomCanvas(testWindow.Canvas())
	testWindow.SetContent(table)
	table.EnterInputMode()
	table.ExitInputMode()
	assert.NotEqual(t, table, testWindow.Canvas().Focused())
}

func TestThatAfterAddingARowPreviousRowsAndNewRowsDataDisplaysOnCorrectPositions(t *testing.T) {
	table := createDefaultTableForTesting()
	table.AddRow(createTestTableRow(testColumnAmount))
	posX := expectedPadding / 2
	posY := expectedHeaderHeight
	for _, row := range table.rowData {
		for i, cell := range row {
			assert.Equal(t, posX, cell.Position().X, "Position x of cell num "+strconv.Itoa(i)+" is incorrect")
			assert.Equal(t, posY, cell.Position().Y, "Position y of cell num "+strconv.Itoa(i)+" is incorrect")
			posX += table.columnLabels[i].Size().Width + expectedPadding
			if i != 0 && (i+1)%testColumnAmount == 0 {
				posX = expectedPadding / 2
				posY += expectedRowHeight
			}
		}
	}
}

func TestThatTableCallsOnExitCallbackFunction(t *testing.T) {
	functionExecuted := false
	function := func() { functionExecuted = true }
	table := createDefaultTableForTesting()
	table.SetOnExitCallbackFunction(function)
	table.EnterInputMode()
	SimulateKeyPressOnTestCanvas(fyne.KeySpace)
	assert.True(t, functionExecuted)
}
