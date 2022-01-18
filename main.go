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

var myApp = app.New()
var TO = TableOtoko{}
func main() {

	myWindow := myApp.NewWindow("TabContainer Widget")

	table := TableInit()
	//	h :=container.New(layout.NewMaxLayout(), table.Table)
	//
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
	var t = make([][]string, 5000)
	for i := range t {
		t[i] = []string{strconv.Itoa(i), "45", "t2", "1", "bottom right"}

	}
	//var TO = TableOtoko{}
	TO.ColumnsName = []string{"N", "n_1", "node_2", "V", "node_4"}
	TO.ColumnsType = []string{"label", "label", "label", "bool", "string"}
	TO.ColumnsWidth = []float32{40, 30, 200, 30, 190}
	TO.AlterRowColor = color.Gray{250}
	TO.HeaderColor = color.Gray{150}
	TO.RowColor = color.Gray{200}
	TO.Data = t
	TO.Edit = true
	TO.makeTable()
	return TO
}
