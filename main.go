package main

import (
	"os"
	"otable/data"
	"otable/pkg/logger"
	"otable/pkg/owidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	table := fd.NewOTable("invoice", *data.TestData(), GetToolBar())
	w := fd.W
	w.Resize(fyne.NewSize(1200, 400))
	w.SetContent(container.NewMax(table))
	w.Show()

}

func GetToolBar() *widget.Toolbar {
	l := logger.GetLog()

	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			l.WithFields(logrus.Fields{"DocumentCreateIcon": "DocumentCreateIcon"}).Info("OnSelectedMakeTableLabel")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			// fd := PutListForm("Property", "Table test")
			// w1 := fd.W
			/* table := t.TableInitProperties()
			table.MakeTable(*data.TestData())
			table.Form = *fd
			table.Edit = true
			l.WithFields(logrus.Fields{"1table.Form ": table.Form}).Info("tableLabel")
			// fd.Table["Property"] = table

			w1.Resize(fyne.NewSize(1200, 400))
			w1.SetContent(container.NewMax(table))

			w1.Show() */
		}))

}
