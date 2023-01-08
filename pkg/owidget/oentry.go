package owidget

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"strings"
)

// / поле ввода
type oEntry struct {
	t *OTable
	widget.Entry
}

func (e *oEntry) Tapped(ev *fyne.PointEvent) {
	//t := appValues[e.IDForm].Table[e.IDTable]
	//n := len(t.Data)
	//row := 0
	//for i := 1; i < n; i++ {
	//	if t.Data[i][0] == e.ID {
	//		row = i
	//		break
	//	}
	//}
	//
	//if row == 0 {
	//	sortS(t.Data, e.col)
	//	for i := 1; i < n; i++ {
	//		t.Data[i][1] = strconv.Itoa(i)
	//	}
	//	t.Table.Refresh()
	//}
}
func (e *oEntry) DoubleTapped(ev *fyne.PointEvent) {
	//t := appValues[e.IDForm].Table[e.IDTable]
	//n := len(t.Data)
	//row := 0
	//for i := 1; i < n; i++ {
	//	if t.Data[i][0] == e.ID {
	//		row = i
	//		break
	//	}
	//}
	//
	//if row == 0 {
	//	sortDown(t.Data, e.col)
	//	n := len(t.Data)
	//	for i := 1; i < n; i++ {
	//		t.Data[i][1] = strconv.Itoa(i)
	//	}
	//	t.Table.Refresh()
	//}
}

func (e *oEntry) menu() {
	menuItems := make([]*fyne.MenuItem, 0)
	menuItem := fyne.NewMenuItem(
		"Отбор",
		func() {
			fmt.Println(e.Text)
		},
	)
	menuItems = append(menuItems, menuItem)
	menuItem = fyne.NewMenuItem(
		"Сортировка",
		func() {
			Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("сортировка ")
		},
	)
	menuItems = append(menuItems, menuItem)
	widget.ShowPopUpMenuAtPosition(
		fyne.NewMenu("", menuItems...),
		fyne.CurrentApp().Driver().CanvasForObject(e),
		e.Position(),
	)

}

func (e *oEntry) TappedSecondary(ev *fyne.PointEvent) {
	e.menu()
}

func (e *oEntry) OnChanged(t string) {
	fmt.Println(e.Entry.Text)
	Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("сортировка ")
}

func NewoEntry() *oEntry {
	entry := &oEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Entry.OnChanged = func(sText string) {
		fmt.Println(sText)
	}
	return entry
}

func (e *oEntry) TypedRune(r rune) {
	if e.t.Form.ActiveWidget.tip == "date" {
		buf := []rune(e.Text)
		switch e.CursorColumn {
		case 0, 1, 2, 3:
			if strings.Contains("0123456789", string(r)) {
				buf[e.CursorColumn] = r
				e.Text = string(buf)

				e.CursorColumn++
				e.Refresh()
			}
		case 4:
			if strings.Contains("012", string(r)) {
				buf[e.CursorColumn+1] = r
				e.Text = string(buf)
				e.CursorColumn++
				e.CursorColumn++
				e.Refresh()
			}
		case 5, 6:
			if strings.Contains("012", string(r)) {
				buf[e.CursorColumn] = r
				e.Text = string(buf)
				e.CursorColumn++
				e.Refresh()
			}
		case 7:
			if strings.Contains("0123", string(r)) {
				buf[e.CursorColumn+1] = r
				e.Text = string(buf)
				e.CursorColumn++
				e.CursorColumn++
				e.Refresh()
			}
		case 8:
			if strings.Contains("0123", string(r)) {
				buf[e.CursorColumn] = r
				e.Text = string(buf)
				e.CursorColumn++
				e.Refresh()
			}
		case 9:
			if strings.Contains("0123456789", string(r)) {
				buf[e.CursorColumn] = r
				e.Text = string(buf)
				e.Refresh()
			}
		}

	}

	//
	//if (r >= '0' && r <= '9') || r == '.' || r == ',' {
	//	e.Entry.TypedRune(r)
	//}
}

func (e *oEntry) TypedKey(key *fyne.KeyEvent) {
	t := e.t
	switch key.Name {
	case "Left":
		if e.CursorColumn == 0 {
			e.t.TypedKey(key)
		} else {
			e.CursorColumn--
			e.Refresh()
		}
	case "Right":
		if e.CursorColumn == len(e.Text) {
			e.t.TypedKey(key)
		} else {
			e.CursorColumn++
			e.Refresh()
		}

	case "Escape":
		t.Edit = false
		t.Form.ActiveWidget.tip = "table"
		t.Form.ActiveWidget.t = t
		t.FocusActiveWidget()

	case "Up", "Down", "Return":
		e.t.TypedKey(key)
	default:

	}
}

func (e *oEntry) KeyDown(key *fyne.KeyEvent) {

}

func (e *oEntry) KeyUp(key *fyne.KeyEvent) {
	//Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("KeyUp ")
	//fmt.Printf("Key %v released\n", key.Name)
}

// oSelect поле выбора
type oSelect struct {
	t *OTable
	widget.SelectEntry
}

func NewoSelect(s []string) *oSelect {
	entry := &oSelect{}
	entry.ExtendBaseWidget(entry)
	entry.SetOptions(s)
	return entry
}

func (e *oSelect) KeyDown(key *fyne.KeyEvent) {
	e.t.TypedKey(key)
}

func (e *oSelect) TypedKey(key *fyne.KeyEvent) {

}
