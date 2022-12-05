package main

import (
	"fmt"
	"otable/pkg/owidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func choiceFromList(l []string, w *owidget.OTable, e *owidget.CompletionEntry) {

	myWindow := myApp.NewWindow("List Widget")

	list := widget.NewList(
		func() int {
			return len(Names)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(Names[i])
		})
	list.OnSelected = func(id widget.ListItemID) {
		g := w.Selected
		w.Properties.DataV[g.Row][g.Col] = l[id]
		//w.Properties.TabStyle[g.Row].BGColor = l[id]
		e.SetText(l[id])
		fmt.Println(l[id])
		myWindow.Close()
	}
	myWindow.Resize(fyne.NewSize(150, 600))
	myWindow.SetContent(list)
	myWindow.Show()
}
