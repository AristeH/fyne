package main

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)
type TableData struct {
	Data   [][]string
	Widget *widget.Table
}

type FormData struct {
	Table TableOtoko
	W     fyne.Window
}

var app_values = make(map[string]FormData)
var myApp = app.New()

func main() {

	myWindow := myApp.NewWindow("TabContainer Widget")

	table := TableInit()
	app_values["main"] = FormData{Table: table, W: myWindow}
	table.makeTable()
	content := container.NewBorder(
		container.NewVBox(
			table.Tool,
			widget.NewSeparator(),
			table.Header,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		table.Table,
	)

	myWindow.Resize(fyne.NewSize(1200, 400))
	myWindow.SetContent(container.NewMax(content))
	myWindow.ShowAndRun()

}

func TableInit() TableOtoko {
	col_columns := 15
	col_rows := 100
	columns := make([]string, col_columns)
	columnstype := make([]string, col_columns)
	data := make([][]string, col_rows)
	columnswidth := make([]float32, col_columns)

	for i := 0; i < col_rows; i++ {
		data[i] = make([]string, col_columns)
		for i1 := 0; i1 < col_columns; i1++ {
			data[i][i1] = "row " + strconv.Itoa(i) + "_" + strconv.Itoa(i1) + ","
			if i1 > 5 && i1 < 10 {
				if i1%3 == 0 {
					data[i][i1] = "0"
				} else {
					data[i][i1] = "1"
				}
			}
		}
	}


		for i1 := 0; i1 < col_columns; i1++ {
			columns[i1] = "Col " + strconv.Itoa(i1)  + ","
			columnswidth[i1] = 150
			if i1 > 5 && i1 < 10 {
				columnstype[i1] = "bool"
				columnswidth[i1] = 30
				columns[i1] =strconv.Itoa(i1) 
			} else {
				columnstype[i1] = "string"
			}
			if i1 < 5 {
				columnswidth[i1] = 90
				columnstype[i1] = "label"
			}
		}
	
	var TO = TableOtoko{}
	TO.ColumnsName = columns
	TO.ColumnsType = columnstype
	TO.ColumnsWidth = columnswidth
	TO.AlterRowColor = color.Gray{250}
	TO.HeaderColor = color.Gray{150}
	TO.RowColor = color.Gray{200}
	TO.Data = data
	TO.Edit = true
	
	return TO
}
