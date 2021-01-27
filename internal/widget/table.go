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
First a header is displayed for which data is passed separately, then the rest of the data.
*/
type Table struct {
	widget.BaseWidget
	headerObjects []fyne.CanvasObject
	objects       []fyne.CanvasObject
	columnAmount  int
}

func NewTable(columnAmount int, headerData []fyne.CanvasObject, data []fyne.CanvasObject) *Table {
	table := &Table{
		headerObjects: headerData,
		objects:       data,
		columnAmount:  columnAmount,
	}
	table.ExtendBaseWidget(table)
	return table
}

func (table Table) HeaderColumns() []fyne.CanvasObject {
	return table.headerObjects
}

func (table Table) CreateRenderer() fyne.WidgetRenderer {
	renderer := tableRenderer{
		table:           table,
		headerRowBorder: canvas.NewRectangle(color.Black),
	}
	return renderer
}

type tableRenderer struct {
	table           Table
	headerRowBorder *canvas.Rectangle
}

func (renderer tableRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (renderer tableRenderer) Destroy() {
	//No resources to clear
}

func (renderer tableRenderer) Layout(fyne.Size) {
	renderer.renderHeader()
	renderer.renderDataRows()
}

func (renderer tableRenderer) renderHeader() {
	renderer.renderHeaderData()
	renderer.renderHeaderRowRectangle()
}

func (renderer tableRenderer) renderHeaderData() {
	position := fyne.NewPos(0, 0)
	for _, object := range renderer.table.headerObjects {
		object.Move(position)
		size := fyne.NewSize(columnWidth, headerHeight)
		object.Resize(size)
		position = position.Add(fyne.NewPos(size.Width+widthBetweenColumns, 0))
	}
}

func (renderer tableRenderer) renderHeaderRowRectangle() {
	renderer.headerRowBorder.Move(fyne.NewPos(0, 0))
	headerRowRectangleSize := fyne.NewSize((columnWidth+widthBetweenColumns)*len(renderer.table.headerObjects), headerHeight)
	renderer.headerRowBorder.StrokeWidth = 2
	renderer.headerRowBorder.FillColor = color.Transparent
	renderer.headerRowBorder.StrokeColor = color.Black
	renderer.headerRowBorder.Resize(headerRowRectangleSize)
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

func (renderer tableRenderer) MinSize() fyne.Size {
	layoutWidth := 0
	layoutHeight := 0
	for i, object := range renderer.table.objects {
		objectMinSize := object.Size()
		layoutWidth += objectMinSize.Width
		if object.Size().Height > layoutHeight {
			layoutHeight = object.Size().Height
		}
		if i == renderer.table.columnAmount-1 {
			break
		}
	}
	amountOfRows := len(renderer.table.objects) / renderer.table.columnAmount
	layoutHeight = amountOfRows * layoutHeight
	return fyne.NewSize(layoutWidth, layoutHeight)
}

func (renderer tableRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{}
	objects = append(objects, renderer.table.objects...)
	objects = append(objects, renderer.table.headerObjects...)
	objects = append(objects, renderer.headerRowBorder)
	return objects
}

func (renderer tableRenderer) Refresh() {
	//Nothing to refresh due to lack of interactivity
}
