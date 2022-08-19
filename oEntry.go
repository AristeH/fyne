package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type oEntry struct {
	IDForm  string
	IDTable string
	widget.Entry
}

func (m *oEntry) onEnter() {
	fmt.Println(m.Entry.Text)
	i := appValues[m.IDForm].Table[m.IDTable].we[m]
	appValues[m.IDForm].Table[m.IDTable].Data[i.Row][i.Col] = m.Entry.Text

}

func newoEntry() *oEntry {
	entry := &oEntry{}
	entry.ExtendBaseWidget(entry)
	return entry

}

func (m *oEntry) KeyDown(key *fyne.KeyEvent) {
	t := appValues[m.IDForm].Table[m.IDTable]
	switch key.Name {
	case fyne.KeyReturn:
		m.onEnter()
	case "KP_Enter":
		m.onEnter()
		i := t.we[m]
		t.ColumnStyle[i.Col].Width = 40
	case "Down":
		i := t.we[m]
		newTableCellID := widget.TableCellID{Col: i.Col, Row: i.Row + 1}
		t.Table.ScrollTo(newTableCellID)
		for key, value := range appValues[m.IDForm].Table[m.IDTable].we {
			if value == newTableCellID {
				appValues[m.IDForm].W.Canvas().Focus(key)
				break
			}
		}
	case "Up":
		i := t.we[m]
		newTableCellID := widget.TableCellID{Col: i.Col, Row: i.Row - 1}
		t.Table.ScrollTo(newTableCellID)
		for key, value := range t.we {
			if value == newTableCellID {
				appValues[m.IDForm].W.Canvas().Focus(key)
				break
			}
		}
	default:
		m.Entry.KeyDown(key)
		fmt.Printf("Key %v pressed\n", key.Name)
	}
}

func (m *oEntry) KeyUp(key *fyne.KeyEvent) {
	fmt.Printf("Key %v released\n", key.Name)
}

// somehow this catches right click. How?
func (m *oEntry) TappedSecondary(pe *fyne.PointEvent) {
	fmt.Printf("1")
}

func (m *oEntry) DoubleTapped(pe *fyne.PointEvent) {
	ind := appValues["main"].Table["tovar"].we[m]
	if ind.Col == 5 {
		choiceFromList(Names, appValues["main"].Table["tovar"], m)
	}
	if ind.Col == 3 {
		var Types = []string{
			"label",
			"string",
			"bool",
			"float",
		}
		choiceFromList(Types, appValues["main"].Table["tovar"], m)
	}
	if ind.Col == 4 {
		//if s, err := strconv.ParseFloat(mc.Text, 32); err == nil {
		appValues["main"].Table["tovar"].ColumnStyle[ind.Row].Width = float32(20) // 3.1415927410125732
		appValues["main"].Table["tovar"].Table.SetColumnWidth(ind.Row, 20)
		//	}
	}

}
func (m *oEntry) OnChanged() {

	fmt.Printf("1")

}
