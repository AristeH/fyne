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
var activeContainer *TableOtoko

func main() {

	myWindow := myApp.NewWindow("Table test")
	contentT := widget.NewButton("Table enter", nil)
	contentT.OnTapped = func() {
		tableEntry()
	}
	contentL := widget.NewButton("Table label", nil)
	contentL.OnTapped = func() {
		tableLabel()
	}
	content := container.NewVBox()
	content.Add(contentT)
	content.Add(contentL)
	myWindow.Resize(fyne.NewSize(600, 200))
	myWindow.SetContent(container.NewMax(content))
	myWindow.ShowAndRun()

}

func tableEntry() {
	wTable := myApp.NewWindow("Some ")
	table := TableInit()
	t := make(map[string]*TableOtoko)
	t["tovar"] = table
	appValues["main"] = FormData{Table: t, W: wTable}
	table.MakeTableEntry()
	activeContainer = table
	table.Properties = TableInitProperties(table)

	t1 := make(map[string]*TableOtoko)
	t1["prop"] = table.Properties
	appValues["mainprop"] = FormData{Table: t1, W: wTable}
	table.Properties.MakeTableEntry()
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
	wTable.Resize(fyne.NewSize(1200, 400))
	wTable.SetContent(container.NewMax(content))
	wTable.Show()

	wTable.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		i := activeContainer.Selected
		switch k.Name {
		case "Down":
			if len(activeContainer.Data)-1 > i.Row {
				activeContainer.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
			}
		case "Up":
			if i.Row > 1 {
				activeContainer.Selected = widget.TableCellID{Col: i.Col, Row: i.Row - 1}
			}
		case "Left":
			if i.Col >= 1 {
				activeContainer.Selected = widget.TableCellID{Col: i.Col - 1, Row: i.Row}

			}
		case "Right":
			if len(activeContainer.Data[0])-1 > i.Col {

				activeContainer.Selected = widget.TableCellID{Col: i.Col + 1, Row: i.Row}
			}
		}
		activeContainer.Table.ScrollTo(activeContainer.Selected)
		activeContainer.Table.Refresh()
	})
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
	TO.we = make(map[*oEntry]widget.TableCellID)
	TO.wl = make(map[*widget.Label]widget.TableCellID)
	TO.wol = make(map[*oLabel]widget.TableCellID)
	return &TO
}

func tableLabel() {
	w1 := myApp.NewWindow("Some ")
	table := TableInit()
	t := make(map[string]*TableOtoko)
	t["tovar"] = table
	appValues["main"] = FormData{Table: t, W: w1}
	table.MakeTableLabel()

	activeContainer = table
	table.Properties = TableInitProperties(table)

	t1 := make(map[string]*TableOtoko)
	t1["prop"] = table.Properties
	appValues["mainprop"] = FormData{Table: t1, W: w1}
	table.Properties.MakeTableEntry()
	activeContainer = table
	//table.Properties.MakeTable()
	content := container.NewBorder(
		table.Tool,
		nil,
		nil,
		nil,
		table.Table,
	)

	w1.Resize(fyne.NewSize(1200, 400))
	w1.SetContent(container.NewMax(content))
	w1.Show()
	w1.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		i := activeContainer.Selected
		switch k.Name {
		case "Down":
			if len(activeContainer.Data)-1 > i.Row {
				activeContainer.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
			}
		case "Up":
			if i.Row > 1 {
				activeContainer.Selected = widget.TableCellID{Col: i.Col, Row: i.Row - 1}
			}
		case "Left":
			if i.Col >= 1 {
				activeContainer.Selected = widget.TableCellID{Col: i.Col - 1, Row: i.Row}

			}
		case "Right":
			if len(activeContainer.Data[0])-1 > i.Col {

				activeContainer.Selected = widget.TableCellID{Col: i.Col + 1, Row: i.Row}
			}
		}
		activeContainer.Table.ScrollTo(activeContainer.Selected)
		activeContainer.Table.Refresh()
	})
}
