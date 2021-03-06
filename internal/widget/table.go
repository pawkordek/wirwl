package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

/*
A widget that consists of data displayed like in a table.
It consists of a header with labels displaying the column names and rows below containing the actual data.
*/
type Table struct {
	widget.BaseWidget
	columnData   []TableColumn
	columnLabels []fyne.CanvasObject
	rowData      []TableRow
	canvas       fyne.Canvas
	focused      bool
}

type TableColumn struct {
	Type ColumnType
	Name string
}

type ColumnType string

const (
	TextColumn  ColumnType = "TEXT_COLUMN"
	ImageColumn ColumnType = "IMAGE_COLUMN"
)

type TableRow []fyne.CanvasObject

func NewTable(canvas fyne.Canvas, columnData []TableColumn, rowData []TableRow) *Table {
	table := &Table{
		columnData:   columnData,
		columnLabels: createColumnLabels(columnData),
		rowData:      rowData,
		canvas:       canvas,
		focused:      false,
	}
	table.ExtendBaseWidget(table)
	return table
}

func createColumnLabels(columnData []TableColumn) []fyne.CanvasObject {
	labels := []fyne.CanvasObject{}
	for _, column := range columnData {
		labels = append(labels, widget.NewLabel(column.Name))
	}
	return labels
}

func (table Table) HeaderColumns() []fyne.CanvasObject {
	return table.columnLabels
}

func (table Table) CreateRenderer() fyne.WidgetRenderer {
	return newTableRenderer(table)
}

func (table Table) columnAmount() int {
	return len(table.columnData)
}

func (table Table) FocusGained() {
	table.focused = true
}

func (table Table) FocusLost() {
	table.focused = false
}

func (table Table) Focused() bool {
	return table.focused
}

func (table Table) TypedRune(rune) {
	//Table will not support any sort of typing therefore no implementation is needed
}

func (table Table) TypedKey(*fyne.KeyEvent) {
}

func (table Table) EnterInputMode() {
	table.canvas.Focus(table)
}

func (table Table) ExitInputMode() {
	table.canvas.Unfocus()
}
