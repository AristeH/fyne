package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"strings"
)

// / поле ввода
type oLabel struct {
	IDForm  string
	IDTable string
	ID      string
	col     int
	widget.Label
	parent *TableOtoko
}

func sortS(x [][]string, k int) {
	var temp []string
	n := len(x)
	for i := 1; i < n; i++ {
		for j := i; j < n; j++ {
			if strings.ToUpper(x[i][k]) > strings.ToUpper(x[j][k]) {
				temp = x[i]
				x[i] = x[j]
				x[j] = temp
			}
		}
	}
}

func sortDown(x [][]string, k int) {
	var temp []string
	n := len(x)
	for i := 1; i < n; i++ {
		for j := i; j < n; j++ {
			if strings.ToUpper(x[i][k]) < strings.ToUpper(x[j][k]) {
				temp = x[i]
				x[i] = x[j]
				x[j] = temp
			}
		}
	}
}

func (e *oLabel) Tapped(ev *fyne.PointEvent) {

	t := e.parent
	id := t.wol[e]
	t.Selected = id

	if id.Row == 0 {
		sortS(t.Data, e.col)
		for i := 1; i < len(t.Data); i++ {
			t.Data[i][1] = strconv.Itoa(i)
		}

	}
	t.Table.Refresh()
	activeContainer = t
}

func (e *oLabel) DoubleTapped(ev *fyne.PointEvent) {
	ind := activeContainer.wol[e]
	if ind.Row == 0 {
		sortDown(activeContainer.Data, e.col)
		n := len(activeContainer.Data)
		for i := 1; i < n; i++ {
			activeContainer.Data[i][1] = strconv.Itoa(i)
		}
		activeContainer.Table.Refresh()
	}

	items := make([]*widget.FormItem, 0)
	for col, style := range activeContainer.ColumnStyle {
		if style.Width != 0 {
			Entry := widget.NewEntry()
			Entry.Validator = getValidator(style.Type)
			Entry.Text = activeContainer.Data[ind.Row][col]
			items = append(items, widget.NewFormItem(style.Name, Entry))

		}
	}
	dialog.ShowForm("введите", "", "cancel", items, func(b bool) {
		if !b {
			return
		}

	}, appValues["main"].W)

}

func (e *oLabel) TappedSecondary(ev *fyne.PointEvent) {
	ind := activeContainer.wol[e]
	Entry := widget.NewEntry()
	Entry.Validator = getValidator(activeContainer.ColumnStyle[ind.Row].Type)
	Entry.Text = activeContainer.Data[ind.Row][ind.Col]
	items := []*widget.FormItem{
		widget.NewFormItem(activeContainer.ColumnStyle[ind.Row].Name, Entry),
	}
	dialog.ShowForm("введите", "", "cancel", items, func(b bool) {
		if !b {
			return
		}
		fmt.Println("KP_Enter", Entry.Text)
		activeContainer.Data[ind.Row][ind.Col] = Entry.Text
	}, appValues["main"].W)

}

func (e *oLabel) onEnter() {
	fmt.Println(e.Text)

}

func (e *oLabel) Focusable(key *fyne.KeyEvent) {
	fmt.Printf("Key %v released\n", key.Name)
}

func (e *oLabel) TypedShortcut(s fyne.Shortcut) {
	if _, ok := s.(*desktop.CustomShortcut); !ok {
		println(s)
		return
	}
	log.Println("Shortcut typed:", s)
}

func newOLabel() *oLabel {
	entry := &oLabel{}
	entry.ExtendBaseWidget(entry)
	return entry
}
