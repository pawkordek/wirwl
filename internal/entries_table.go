package wirwl

import (
	fyneWidget "fyne.io/fyne/v2/widget"
	"strconv"
	"wirwl/internal/data"
	widget "wirwl/internal/widget"
)

//Should be equal to amount of fields Entry type has, minus the id field
const columnAmount = 14

func (app *App) createEntriesTable(entryType data.EntryType, entries []data.Entry) {
	rowData := []widget.TableRow{}
	columnData := createColumnData()
	for i, entry := range entries {
		row := createEntriesTableRow(i, entry)
		rowData = append(rowData, row)
	}
	table := widget.NewTable(app.mainWindow.Canvas(), app.inputHandler, columnData, rowData)
	app.entriesTables[entryType] = table
}

func createEntriesTableRow(rowNum int, entry data.Entry) widget.TableRow {
	row := widget.TableRow{}
	row = append(row, newSpreadsheetLabelWithNumber(rowNum))
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
	return row
}

func newSpreadsheetLabelWithText(text string) *fyneWidget.Label {
	label := fyneWidget.NewLabel(text)
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
