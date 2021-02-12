package wirwl

import (
	"fyne.io/fyne"
	fyneWidget "fyne.io/fyne/widget"
	"strconv"
	"wirwl/internal/data"
	widget "wirwl/internal/widget"
)

//Should be equal to amount of fields Entry type has, minus the id field
const columnAmount = 14

func (app *App) createEntriesTable(entries []data.Entry) {
	rowData := []widget.TableRow{}
	columnData := createColumnData()
	for i, entry := range entries {
		row := widget.TableRow{}
		row = append(row, newSpreadsheetLabelWithNumber(i))
		row = append(row, newSpreadsheetLabelWithText("This will be an image"))
		row = append(row, newSpreadsheetLabelWithText(string(entry.Status)))
		row = append(row, newSpreadsheetLabelWithText(entry.Title))
		row = append(row, newSpreadsheetLabelWithNumber(entry.ElementsCompleted))
		row = append(row, newSpreadsheetLabelWithNumber(entry.TotalAmountOfElementsToComplete))
		row = append(row, newSpreadsheetLabelWithNumber(entry.Score))
		row = append(row, newSpreadsheetLabelWithText(entry.StartDate))
		row = append(row, newSpreadsheetLabelWithText(entry.FinishDate))
		row = append(row, newSpreadsheetLabelWithText(entry.Link))
		row = append(row, newSpreadsheetLabelWithText(entry.Description))
		row = append(row, newSpreadsheetLabelWithText(entry.Comment))
		row = append(row, newSpreadsheetLabelWithText(entry.Tags))
		row = append(row, newSpreadsheetLabelWithText(entry.ImageQuery))
		rowData = append(rowData, row)
	}
	table := widget.NewTable(columnData, rowData)
	app.entriesTable = table
}

func newSpreadsheetLabelWithText(text string) *fyneWidget.Label {
	label := fyneWidget.NewLabel(text)
	label.Alignment = fyne.TextAlignCenter
	return label
}

func newSpreadsheetLabelWithNumber(number int) *fyneWidget.Label {
	return newSpreadsheetLabelWithText(strconv.Itoa(number))
}

func createColumnData() []widget.TableColumn {
	columnsNames := []string{
		"Num",
		"Image",
		"Status",
		"Title",
		"Elements completed",
		"Total amount",
		"Score",
		"Start date",
		"Finish date",
		"Link",
		"Description",
		"Comment",
		"Tags",
		"Image query",
	}
	columnData := []widget.TableColumn{}
	for _, columnName := range columnsNames {
		column := widget.TableColumn{Type: widget.TextColumn, Name: columnName}
		columnData = append(columnData, column)
	}
	return columnData
}
