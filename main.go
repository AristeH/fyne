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

var myApp fyne.App

func main() {
	l := logger.GetLog()
	l.WithFields(logrus.Fields{
		"form":  "main",
		"event": "start",
		"Out":   os.Stderr,
	}).Info("Начало")
	os.Setenv("FYNE_FONT", "C:/goproject/otable/ttf/Go Mono Nerd Font Complete Mono.ttf")

	myApp = app.New()

	myWindow := myApp.NewWindow("Test table")
	fd := owidget.InitFormData("Main")

	fd.W = myWindow

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
	Log := logger.GetLog()

	w1 := myApp.NewWindow("Table test")
	fd := owidget.InitFormData("Table")
	fd.W = w1
	table := owidget.OTable{}
	table.Form = *fd
	table.Edit = true
	Log.WithFields(logrus.Fields{"1table.Form ": table.Form}).Info("tableLabel")
	fd.Table["invoice"] = &table
	table.MakeTable(*data.TestData())

	Log.WithFields(logrus.Fields{"3table.Form ": table.Form}).Info("tableLabel")

	w1.Resize(fyne.NewSize(1200, 400))

	w1.SetContent(container.NewMax(&table))
	w1.Show()

}
