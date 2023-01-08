package main

import (
	"os"
	"otable/data"
	"otable/pkg/logger"
	"otable/pkg/owidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

func main() {

	l := logger.GetLog()
	l.WithFields(logrus.Fields{
		"form":  "main",
		"event": "start",
		"Out":   os.Stderr,
	}).Info("Начало")
	os.Setenv("FYNE_FONT", "C:\\goproject\\otable/ttf/Go Mono Nerd Font Complete Mono.ttf")
	app.New()
	fd := owidget.PutListForm("Main", "MainForm")
	myWindow := fd.W
	bTable := widget.NewButton("Open table", nil)
	bTable.OnTapped = func() {
		tableLabel()
	}
	content := container.NewVBox()
	content.Add(bTable)
	myWindow.Resize(fyne.NewSize(600, 200))
	myWindow.SetContent(container.NewMax(content))
	myWindow.ShowAndRun()

}

func tableLabel() {
	fd := owidget.PutListForm("Table", "Table test")
	table := fd.NewOTable("invoice", *data.TestData())
	table.CellColor["3;3"] = &owidget.CellColor{
		Color:   owidget.MapColor["aliceblue"],
		BGcolor: owidget.MapColor["darkgreen"]}
	table.ColumnStyle["Amount"].BGcolor = "darkgreen"
	l := logger.GetLog()
	l.WithFields(logrus.Fields{
		"form":  "main",
		"event": "start",
		"Out":   os.Stderr,
	}).Info("GetToolBar")
	table.GetToolBar()
	w := fd.W
	w.Resize(fyne.NewSize(1200, 400))
	w.SetContent(container.NewMax(table))

	w.Show()

}
