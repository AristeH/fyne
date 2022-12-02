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
	os.Setenv("FYNE_FONT", "C:\\проект\\Table\\ttf\\FiraCode-Regular.ttf")

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
	fd.ActiveContainer = &table
	Log.WithFields(logrus.Fields{"3table.Form ": table.Form}).Info("tableLabel")
	// content := container.NewBorder(
	// 	nil,
	// 	nil,
	// 	nil,
	// 	nil,
	// 	table.Table,
	// )

	w1.Resize(fyne.NewSize(1200, 400))

	w1.SetContent(container.NewMax(&table))
	w1.Show()

	w1.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {

		fd1 := owidget.GetApp()

		fd = fd1["Table"]

		i := fd.ActiveContainer.Selected
		otab := fd.ActiveContainer
		Log.WithFields(logrus.Fields{"otab.selected": otab.Selected, "key": k.Name}).Info("Form")
		switch k.Name {
		case "Insert":
			otab.Edit = true
			otab.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}

			otab.FocusActiveWidget()
		case "Return":
			if otab.Edit {
				otab.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
				otab.FocusActiveWidget()

			} else {

				otab.Edit = true
				otab.Selected = widget.TableCellID{Col: i.Col, Row: i.Row}
				otab.FocusActiveWidget()

			}
		case "Down":
			if len(otab.Data)-1 > i.Row {
				otab.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
			}
		case "Up":
			if i.Row > 0 {
				tc := widget.TableCellID{Col: i.Col, Row: i.Row - 1}
				otab.Selected = tc
			}
		case "Left":
			c := i.Col
			for c >= 1 {
				c--
				col := otab.ColumnStyle[otab.DataV[0][c]]
				if col.Width != 0 {
					otab.Selected = widget.TableCellID{Col: c, Row: i.Row}
					break
				}
			}
		case "Right":
			c := i.Col
			col := otab.ColumnStyle[otab.DataV[0][c]]

			for len(otab.DataV[0])-1 > c {
				c++
				if col.Width != 0 {
					otab.Selected = widget.TableCellID{Col: c, Row: i.Row}
					break
				}
			}
		}

		if i != fd.ActiveContainer.Selected {
			fd.ActiveContainer.FocusActiveWidget()
		}
	})

}
