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
	tableData := []fyne.CanvasObject{}
	headerData := createHeader()
	for i, entry := range entries {
		tableData = append(tableData, newSpreadsheetLabelWithNumber(i))
		tableData = append(tableData, newSpreadsheetLabelWithText("This will be an image"))
		tableData = append(tableData, newSpreadsheetLabelWithText(string(entry.Status)))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.Title))
		tableData = append(tableData, newSpreadsheetLabelWithNumber(entry.ElementsCompleted))
		tableData = append(tableData, newSpreadsheetLabelWithNumber(entry.TotalAmountOfElementsToComplete))
		tableData = append(tableData, newSpreadsheetLabelWithNumber(entry.Score))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.StartDate))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.FinishDate))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.Link))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.Description))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.Comment))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.Tags))
		tableData = append(tableData, newSpreadsheetLabelWithText(entry.ImageQuery))
	}
	table := widget.NewTable(columnAmount, headerData, tableData)
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

func createHeader() []fyne.CanvasObject {
	headerColumns := []string{
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
	headerData := []fyne.CanvasObject{}
	for _, columnName := range headerColumns {
		label := fyneWidget.NewLabel(columnName)
		label.TextStyle = fyne.TextStyle{
			Bold:      true,
			Italic:    false,
			Monospace: false,
		}
		label.Alignment = fyne.TextAlignCenter
		headerData = append(headerData, label)
	}
	return headerData
}
