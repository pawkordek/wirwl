package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

const columnWidth = 100
const headerHeight = 50
const rowHeight = 141
const widthBetweenColumns = 35

/*
A widget that consists of data displayed like in a table.
It consists of a header with labels displaying the column names and rows below containing the actual data.
*/
type Table struct {
	widget.BaseWidget
	columnData   []TableColumn
	columnLabels []fyne.CanvasObject
	objects      []fyne.CanvasObject
	columnAmount int
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

func NewTable(columnAmount int, columnData []TableColumn, data []fyne.CanvasObject) *Table {
	table := &Table{
		columnData:   columnData,
		columnLabels: createColumnLabels(columnData),
		objects:      data,
		columnAmount: columnAmount,
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

type tableRenderer struct {
	table           Table
	headerRowBorder *canvas.Rectangle
	dataRowsBorders []*canvas.Rectangle
	columnBorders   []*canvas.Rectangle
	borderColor     color.Color
}

func newTableRenderer(table Table) tableRenderer {
	dataRowsBorders := createBorders(len(table.objects) / table.columnAmount)
	return tableRenderer{
		table:           table,
		headerRowBorder: canvas.NewRectangle(color.Black),
		dataRowsBorders: dataRowsBorders,
		columnBorders:   createBorders(table.columnAmount),
		borderColor:     color.Black,
	}
}

func createBorders(amount int) []*canvas.Rectangle {
	borders := []*canvas.Rectangle{}
	for i := 1; i <= amount; i++ {
		borders = append(borders, canvas.NewRectangle(color.Black))
	}
	return borders
}

func (renderer tableRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (renderer tableRenderer) Destroy() {
	//No resources to clear
}

func (renderer tableRenderer) Layout(fyne.Size) {
	renderer.renderHeader()
	renderer.renderData()
}

func (renderer tableRenderer) renderHeader() {
	renderer.renderHeaderData()
	renderer.renderHeaderRowRectangle()
}

func (renderer tableRenderer) renderHeaderData() {
	position := fyne.NewPos(0, 0)
	for _, object := range renderer.table.columnLabels {
		object.Move(position)
		size := fyne.NewSize(columnWidth, headerHeight)
		object.Resize(size)
		position = position.Add(fyne.NewPos(size.Width+widthBetweenColumns, 0))
	}
}

func (renderer tableRenderer) renderHeaderRowRectangle() {
	renderer.headerRowBorder.Move(fyne.NewPos(0, 0))
	headerRowRectangleSize := fyne.NewSize((columnWidth+widthBetweenColumns)*len(renderer.table.columnLabels), headerHeight)
	renderer.headerRowBorder.StrokeWidth = 2
	renderer.headerRowBorder.FillColor = color.Transparent
	renderer.headerRowBorder.StrokeColor = renderer.borderColor
	renderer.headerRowBorder.Resize(headerRowRectangleSize)
}

func (renderer tableRenderer) renderData() {
	renderer.renderDataRows()
	renderer.renderDataRowsBorders()
	renderer.renderColumnBorders()
}

func (renderer tableRenderer) renderDataRows() {
	position := fyne.NewPos(0, headerHeight)
	currentColumnNum := 1
	for _, object := range renderer.table.objects {
		size := fyne.NewSize(columnWidth, rowHeight)
		object.Resize(size)
		object.Move(position)
		position = position.Add(fyne.NewPos(size.Width+widthBetweenColumns, 0))
		if currentColumnNum == renderer.table.columnAmount {
			position = position.Add(fyne.NewPos(0, size.Height))
			position = position.Subtract(fyne.NewPos(position.X, 0))
			currentColumnNum = 0
		}
		currentColumnNum += 1
	}
}

func (renderer tableRenderer) renderDataRowsBorders() {
	size := fyne.NewSize((columnWidth+widthBetweenColumns)*len(renderer.table.columnLabels), rowHeight)
	position := fyne.NewPos(0, headerHeight)
	for _, border := range renderer.dataRowsBorders {
		border.Move(position)
		border.StrokeWidth = 2
		border.FillColor = color.Transparent
		border.StrokeColor = renderer.borderColor
		border.Resize(size)
		position = position.Add(fyne.NewPos(0, rowHeight))
	}
}

func (renderer tableRenderer) renderColumnBorders() {
	size := fyne.NewSize(columnWidth+widthBetweenColumns, headerHeight+rowHeight*len(renderer.table.objects)/renderer.table.columnAmount)
	position := fyne.NewPos(0, 0)
	for _, border := range renderer.columnBorders {
		border.Move(position)
		border.StrokeWidth = 2
		border.FillColor = color.Transparent
		border.StrokeColor = renderer.borderColor
		border.Resize(size)
		position = position.Add(fyne.NewPos(columnWidth+widthBetweenColumns, 0))
	}
}

func (renderer tableRenderer) MinSize() fyne.Size {
	layoutWidth := 0
	layoutHeight := 0
	for i, object := range renderer.table.objects {
		objectMinSize := object.Size()
		layoutWidth += objectMinSize.Width + widthBetweenColumns
		if object.Size().Height > layoutHeight {
			layoutHeight = object.Size().Height
		}
		if i == renderer.table.columnAmount-1 {
			break
		}
	}
	amountOfRows := len(renderer.table.objects) / renderer.table.columnAmount
	layoutHeight = amountOfRows * layoutHeight + headerHeight
	return fyne.NewSize(layoutWidth, layoutHeight)
}

func (renderer tableRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{}
	objects = append(objects, renderer.table.objects...)
	objects = append(objects, renderer.table.columnLabels...)
	objects = append(objects, renderer.headerRowBorder)
	for _, border := range renderer.dataRowsBorders {
		objects = append(objects, border)
	}
	for _, border := range renderer.columnBorders {
		objects = append(objects, border)
	}
	return objects
}

func (renderer tableRenderer) Refresh() {
	//Nothing to refresh due to lack of interactivity
}
