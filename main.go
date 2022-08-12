package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FormData struct {
	Table map[string]*TableOtoko //Entry - список таблиц формы
	W     fyne.Window
}

var appValues = make(map[string]FormData)
var myApp = app.New()

func main() {

	myWindow := myApp.NewWindow("TabContainer Widget")

	table := TableInit()
	t := make(map[string]*TableOtoko)
	t["tovar"] = table
	appValues["main"] = FormData{Table: t, W: myWindow}
	table.MakeTable()
	table.Properties = TableInitProperties(table)

	t1 := make(map[string]*TableOtoko)
	t1["prop"] = table.Properties
	appValues["mainprop"] = FormData{Table: t1, W: myWindow}
	table.Properties.MakeTable()
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

func TableInit() *TableOtoko {
	colColumns := 12
	colRows := 20
	data := make([][]string, colRows)
	for i := 0; i < colRows; i++ {
		data[i] = make([]string, colColumns)
		for i1 := 0; i1 < colColumns; i1++ {
			data[i][i1] = "row " + strconv.Itoa(i) + "_" + strconv.Itoa(i1) + ","
			if i1 > 5 && i1 < 10 {
				if i1%3 == 0 {
					data[i][i1] = "0"
				} else {
					data[i][i1] = "1"
				}
			}
		}
		data[i][0] = fmt.Sprintf("%d", i)
	}
	var TO = TableOtoko{}
	for i1 := 0; i1 < colColumns; i1++ {
		cs := ColumnStyle{}
		cs.Name = "Col " + strconv.Itoa(i1)
		cs.Width = 150

		if i1 > 5 && i1 < 10 {
			cs.Type = "bool"
			cs.Width = 30
			cs.Name = strconv.Itoa(i1)
		} else {
			cs.Type = "string"
		}
		if i1 < 5 {
			cs.Width = 90
			cs.Name = "label" + strconv.Itoa(i1)
			cs.Type = "label"
			if i1 == 3 {
				cs.BGColor = "skyblue"
			}
		}
		if i1 == 0 {
			cs.Width = 40
			cs.Name = "N"
			cs.Format = "int"
		}
		TO.ColumnStyle = append(TO.ColumnStyle, cs)
	}
	ts := TabStyle{}

	ts.RowAlterColor = "darkgray"
	ts.HeaderColor = "darkslategray"
	ts.RowColor = "lightgray"
	TO.TabStyle = ts
	TO.Data = data
	TO.ID = "tovar"
	TO.IDForm = "main"
	TO.wb = make(map[*widget.Button]int)
	TO.wc = make(map[*widget.Check]widget.TableCellID)
	TO.we = make(map[*enterEntry]widget.TableCellID)
	TO.wl = make(map[*widget.Label]widget.TableCellID)
	return &TO
}
