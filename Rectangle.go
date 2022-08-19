package main

import (
	"fmt"
	"fyne.io/fyne/v2"
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

	t := appValues[e.IDForm].Table[e.IDTable]
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
	t := appValues[e.IDForm].Table[e.IDTable]
	n := len(t.Data)
	row := 0
	for i := 1; i < n; i++ {
		if t.Data[i][0] == e.ID {
			row = i
			break
		}
	}

	if row == 0 {
		sortDown(t.Data, e.col)
		n := len(t.Data)
		for i := 1; i < n; i++ {
			t.Data[i][1] = strconv.Itoa(i)
		}
		t.Table.Refresh()
	}
}

func (e *oLabel) TappedSecondary(ev *fyne.PointEvent) {
	fmt.Println(e.Text)
}

func (e *oLabel) onEnter() {
	fmt.Println(e.Text)

}

func (e *oLabel) Focusable(key *fyne.KeyEvent) {
	fmt.Printf("Key %v released\n", key.Name)
}

func newOLabel() *oLabel {
	entry := &oLabel{}
	entry.ExtendBaseWidget(entry)
	return entry
}
func (e *oLabel) TypedShortcut(s fyne.Shortcut) {
	if _, ok := s.(*desktop.CustomShortcut); !ok {
		println(s)
		return
	}

	log.Println("Shortcut typed:", s)
}
