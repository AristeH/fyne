package owidget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
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
	// menuItems := make([]*fyne.MenuItem, 0)
	// menuItem := fyne.NewMenuItem(
	// 	"Отбор",
	// 	func() {
	// 		fmt.Println(e.Text)
	// 	},
	// )
	// menuItems = append(menuItems, menuItem)
	// menuItem = fyne.NewMenuItem(
	// 	"Сортировка",
	// 	func() {
	// 		Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("сортировка ")
	// 	},
	// )
	// menuItems = append(menuItems, menuItem)
	// widget.ShowPopUpMenuAtPosition(
	// 	fyne.NewMenu("", menuItems...),
	// 	fyne.CurrentApp().Driver().CanvasForObject(e),
	//
	// )
	//
}

func (e *oEntry) TappedSecondary(ev *fyne.PointEvent) {
	// e.menu()
}

func (e *oEntry) onEnter() {
	t := e.t
	id := t.Selected
	Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("onEnter ")
	if t.Edit {
		t.DataV[id.Row][id.Col] = e.Text
		if len(t.DataV)-1 > t.Selected.Row {
			t.Selected = widget.TableCellID{Col: id.Col, Row: id.Row + 1}
			t.Table.ScrollTo(widget.TableCellID{Col: id.Col, Row: id.Row + 1})
		}
	} else {
		t.Edit = true
		t.Selected = widget.TableCellID{Col: id.Col, Row: id.Row}
	}
	t.FocusActiveWidget()
}

func (e *oEntry) OnChanged(t string) {
	fmt.Println(e.Entry.Text)
	Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("сортировка ")
}

func newoEntry() *oEntry {
	entry := &oEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Entry.OnChanged = func(sText string) {
		fmt.Println(sText)
	}
	return entry
}

func (e *oEntry) KeyDown(key *fyne.KeyEvent) {
	//t := appValues[e.IDForm].Table[e.IDTable]
	Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("KeyDown")
	Log.WithFields(logrus.Fields{"e.t.Selected": e.t.Selected}).Info("KeyDown")
	t := e.t
	id := t.Selected
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	case "KP_Enter":
		e.onEnter()
	case "Down":
		if len(e.t.Data)-1 > e.t.Selected.Row {
			t.Selected = widget.TableCellID{Col: id.Col, Row: id.Row + 1}
			t.FocusActiveWidget()
		}
	case "Menu":
		e.menu()
	case "Up":
		if id.Row > 1 {
			t.Selected = widget.TableCellID{Col: id.Col, Row: id.Row - 1}
			t.FocusActiveWidget()
		}
	case "Escape":
		t.Edit = false
		t.Form.ActiveWidget = nil
		t.FocusActiveWidget()
		fmt.Printf("Key %v pressed\n", key.Name)
	case "Left":
		if !t.Edit {
			c := id.Col
			for c >= 1 {
				c--
				col := t.ColumnStyle[t.DataV[0][c]]
				if col.Width != 0 {
					t.Selected = widget.TableCellID{Col: c, Row: id.Row}
					break
				}
			}
			t.FocusActiveWidget()
		}
	case "Insert":
		t.Edit = true
		t.Selected = widget.TableCellID{Col: id.Col, Row: id.Row}
		t.FocusActiveWidget()
	case "Right":
		if !t.Edit {
			c := id.Col
			col := t.ColumnStyle[t.DataV[0][c]]
			for len(t.DataV[0])-1 > c {
				c++
				if col.Width != 0 {
					t.Selected = widget.TableCellID{Col: c, Row: id.Row}
					break
				}
			}
			t.FocusActiveWidget()
		}
	default:
		fmt.Printf("Key %v pressed\n", key.Name)
	}
}

func (e *oEntry) KeyUp(key *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"entry.text": e.Text}).Info("KeyUp ")
	fmt.Printf("Key %v released\n", key.Name)
}

// / поле ввода
type oSelect struct {
	t *OTable
	widget.SelectEntry
}

func newoSelect() *oSelect {
	entry := &oSelect{}
	entry.ExtendBaseWidget(entry)
	return entry
}
func (e *oSelect) TypedKey(key *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"key": key}).Info("oSelect ")
	switch key.Name {
	case fyne.KeyEscape:
		// e.onEsc()
	default:
		e.SelectEntry.TypedKey(key)
	}
}
