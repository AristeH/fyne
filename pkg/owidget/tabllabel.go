package owidget

import (
	"context"
	"fmt"
	"image/color"
	"otable/pkg/logger"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/d5/tengo/v2"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logger.GetLog()
}

func (t *OTable) GetToolBar() {

	Log.WithFields(logrus.Fields{"DocumentCreateIcon": "GetToolBar"}).Info("GetToolBar")

	t.Tool = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			Log.WithFields(logrus.Fields{"DocumentCreateIcon": "DocumentCreateIcon"}).Info("GetToolBar")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fd := PutListForm("TableProp", "Tablerop")
			g := t.properties()
			table := fd.NewOTable("invoice", *g)
			w := fd.W
			w.Resize(fyne.NewSize(1200, 400))

			w.SetContent(container.NewMax(table))
			w.SetOnClosed(func() {
				for i := 1; i < len(table.DataV); i++ {
					t.ColumnStyle[table.DataV[i][1]].BGcolor = table.DataV[i][5]
					t.ColumnStyle[table.DataV[i][1]].color = table.DataV[i][4]
					b, _ := strconv.ParseFloat(table.DataV[i][6], 32)
					t.ColumnStyle[table.DataV[i][1]].Width = float32(b)
					t.ColumnStyle[table.DataV[i][1]].formula = table.DataV[i][3]
				}
				for n := 0; n < len(t.DataV[0]); n++ {
					col := t.ColumnStyle[t.DataV[0][n]]
					si := fyne.MeasureText("ш", 12, fyne.TextStyle{})
					t.Table.SetColumnWidth(n, si.Width*col.Width)
					t.Header.SetColumnWidth(n, si.Width*col.Width)
				}
			})
			w.Show()
		}))

}

func (t *OTable) MakeTableLabel() {
	rows := len(t.DataV)
	columns := len(t.DataV[0])
	t.Header = widget.NewTable(
		func() (int, int) { return 1, columns },
		func() fyne.CanvasObject { return canvas.NewText("", color.Black) },
		func(cellID widget.TableCellID, o fyne.CanvasObject) {
			colst := t.ColumnStyle[t.DataV[0][cellID.Col]]
			l := o.(*canvas.Text)
			l.Text = colst.name
			l.Refresh()
		},
	)

	t.Table = widget.NewTable(
		func() (int, int) { return rows, columns },
		func() fyne.CanvasObject {
			fc := CellColor{}
			return t.MakeTappable("", "", &fc)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			FillColor := t.getColorCell(i)
			col := t.ColumnStyle[t.DataV[0][i.Col]]
			tip := col.tip
			if i.Row == 0 {
				tip = "string"
			}
			mystr := []rune(t.DataV[i.Row][i.Col])
			k := int(col.Width)
			if col.Width > 0 {
				if len(mystr) > k-3 {
					k = int(col.Width) - 3
				}
			}
			if tip == "bool" {
				rec := canvas.NewRectangle(FillColor.BGcolor)
				image := canvas.NewImageFromResource(theme.CheckButtonCheckedIcon())
				if t.DataV[i.Row][i.Col] == "1" {
					image = canvas.NewImageFromResource(theme.CheckButtonIcon())
				}
				box.Objects[0] = container.New(layout.NewMaxLayout(), rec, image)
			} else {
				en := string(mystr[0:k])
				if i.Row == 0 {
					switch col.sort {
					case 2:
						en = "\u2191" + en
					case 1:
						en = "\u2193" + en
					}

				}
				entry := t.MakeTappable(en, tip, FillColor)
				box.Objects[0] = entry
			}
			// активная ячейка
			if i == t.Selected {
				t.Form.ActiveWidget.tip = "table"
				if i.Row > 0 && t.Edit {
					if strings.HasPrefix(tip, "id_") { //id другой таблицы
						tip = "id"
					}
					if strings.HasPrefix(tip, "float") {
						tip = "float"
					}

					switch tip {
					case "id":
						FillColor := t.getColorCell(i)
						c := t.MakeTappable(t.DataV[i.Row][i.Col], tip, FillColor)
						box.Objects[0] = container.NewBorder(nil, nil, nil, newTappableIcon(theme.SearchIcon()), c)
					case "float", "int", "string":
						c := NewoEntry()
						c.Text = t.DataV[i.Row][i.Col]
						c.t = t
						t.Form.ActiveWidget.tip = "string"
						t.Form.ActiveWidget.ce = c
						c.t = t
						box.Objects[0] = container.New(layout.NewMaxLayout(), c)
					case "bool":

						ic := newTappableIcon(theme.CheckButtonIcon())
						if t.DataV[i.Row][i.Col] == "0" {
							ic = newTappableIcon(theme.CheckButtonCheckedIcon())
						}
						ic.t = t
						t.Form.ActiveWidget.tip = "bool"
						t.Form.ActiveWidget.ti = ic
						box.Objects[0] = ic
					case "enum":
						c := NewoSelect(t.Enum[col.id])
						c.OnChanged = func(s string) {
							t.DataV[i.Row][i.Col] = s
						}
						c.t = t
						c.Entry.Text = t.DataV[i.Row][i.Col]
						t.Form.ActiveWidget.tip = "enum"
						t.Form.ActiveWidget.sel = c
						box.Objects[0] = container.New(layout.NewMaxLayout(), c)
					case "date":
						c := NewoEntry()
						c.Text = t.DataV[i.Row][i.Col]
						c.t = t
						t.Form.ActiveWidget.tip = "date"
						t.Form.ActiveWidget.ce = c
						c.t = t
						box.Objects[0] = container.New(layout.NewMaxLayout(), c)
					}
					Log.WithFields(logrus.Fields{"T.Form.ActiveWidget": t.Form.ActiveWidget}).Info("t.Selected")
				}
			}
		})
	for n := 0; n < columns; n++ {
		col := t.ColumnStyle[t.DataV[0][n]]
		si := fyne.MeasureText("ш", 12, fyne.TextStyle{})
		t.Table.SetColumnWidth(n, si.Width*col.Width)
		t.Header.SetColumnWidth(n, si.Width*col.Width)
	}
	t.ExtendBaseWidget(t)
	t.Table.OnSelected = func(id widget.TableCellID) {
		Log.WithFields(logrus.Fields{"t.Form": t.Form, "w": id}).Info("OnSelectedMakeTableLabel")
		t.Selected = id
		t.Form.ActiveWidget.t = t

		t.FocusActiveWidget()
	}
}

func (t *OTable) MakeTappable(txt string, tip string, c *CellColor) *fyne.Container {
	entry := canvas.NewText(strings.TrimRight(txt, "\x00"), c.Color)
	if strings.HasPrefix(tip, "float") {
		tip = "float"
	}
	switch tip {
	case "float", "int":
		entry.Alignment = fyne.TextAlignTrailing
		entry.TextStyle.Monospace = true
	default:
		entry.Alignment = fyne.TextAlignLeading
	}
	si := fyne.MeasureText("шii", 24, fyne.TextStyle{})
	rec := canvas.NewRectangle(c.BGcolor)
	rec.SetMinSize(si)
	entry.Resize(si)
	return container.New(layout.NewMaxLayout(), rec, entry)
}

func (t *OTable) ExecuteFormula() {
	col := t.ColumnStyle[t.DataV[0][t.Selected.Col]]
	if col.formula == "" {
		return
	} else {
		script := tengo.NewScript([]byte(col.formula))
		for i := range t.DataV[0] {
			tip := t.ColumnStyle[t.DataV[0][i]].tip

			if strings.HasPrefix(t.ColumnStyle[t.DataV[0][i]].tip, "float") {
				v1, _ := strconv.ParseFloat(t.DataV[t.Selected.Row][i], 32)
				_ = script.Add(t.ColumnStyle[t.DataV[0][i]].id, v1)
			}
			if tip == "int" {
				v1, _ := strconv.Atoi(t.DataV[t.Selected.Row][i])
				_ = script.Add(t.ColumnStyle[t.DataV[0][i]].id, v1)
			}
		}
		// run the script
		compiled, err := script.RunContext(context.Background())
		if err != nil {
			panic(err)
		}
		for i := range t.DataV[0] {
			if strings.HasPrefix(t.ColumnStyle[t.DataV[0][i]].tip, "float") {
				v := t.ColumnStyle[t.DataV[0][i]].id
				t.DataV[t.Selected.Row][i] = fmt.Sprintf("%.2f", compiled.Get(v).Float())
			} else if t.ColumnStyle[t.DataV[0][i]].tip == "int" {
				v := t.ColumnStyle[t.DataV[0][i]].id
				t.DataV[t.Selected.Row][i] = fmt.Sprintf("%d", compiled.Get(v).Int())
			}
		}
	}
}

func (t *OTable) Tapped(ev *fyne.PointEvent) {

}

// TypedKey Implements: fyne.Focusable
func (t *OTable) TypedKey(ev *fyne.KeyEvent) {
	i := t.Selected
	switch ev.Name {
	case "Return":
		if i.Row == 0 {
			if t.Selected.Row == 0 {
				col := t.ColumnStyle[t.DataV[0][t.Selected.Col]]
				switch col.sort {
				case 0, 2:
					col.sort = 1
					t.sortUp()
				case 1:
					col.sort = 2
					t.sortDown()
				}
				for i, _ := range t.DataV[0] {
					if i != t.Selected.Col {
						t.ColumnStyle[t.DataV[0][i]].sort = 0
					}
				}

			}
		} else {
			t.ExecuteFormula()
			if t.Edit {
				switch t.Form.ActiveWidget.tip {
				case "string", "date":
					t.DataV[i.Row][i.Col] = t.Form.ActiveWidget.ce.Text
				}
				if len(t.DataV)-1 > t.Selected.Row {
					t.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
				}
			} else {
				t.Edit = true
				t.Selected = widget.TableCellID{Col: i.Col, Row: i.Row}
			}
		}
	case "Down":
		if len(t.Data) > i.Row {
			t.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
		}
	case "Up":
		if i.Row > 0 {
			t.Selected = widget.TableCellID{Col: i.Col, Row: i.Row - 1}
		}
	case "Left":
		c := i.Col
		for c >= 1 {
			c--
			col := t.ColumnStyle[t.DataV[0][c]]
			if col.Width != 0 {
				t.Selected = widget.TableCellID{Col: c, Row: i.Row}
				break
			}
		}
	case "Escape":
		t.Edit = false
		t.Form.ActiveWidget.tip = "table"
		t.Form.ActiveWidget.t = t
	case "Right":
		c := i.Col
		col := t.ColumnStyle[t.DataV[0][c]]
		for len(t.DataV[0])-1 > c {
			c++
			if col.Width != 0 {
				t.Selected = widget.TableCellID{Col: c, Row: i.Row}
				break
			}
		}

	case fyne.KeySpace:
		if t.Edit && t.Form.ActiveWidget.tip == "bool" {
			if t.DataV[i.Row][i.Col] == "1" {
				t.DataV[i.Row][i.Col] = "0"
			} else {
				t.DataV[i.Row][i.Col] = "1"
			}
		}

	}
	t.FocusActiveWidget()
}

func (t *OTable) TypedRune(r rune) {
	Log.WithFields(logrus.Fields{"entry.text": r}).Info("onEnter ")
}
func (t *OTable) KeyDown(key *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"rows": key}).Info("TappedTappableIcon")
}

// FocusLost Implements: fyne.Focusable
func (t *OTable) FocusLost() {
}

// FocusGained Implements: fyne.Focusable
func (t *OTable) FocusGained() {
}

// FocusActiveWidget - get focus active ceil table
func (t *OTable) FocusActiveWidget() {
	Log.WithFields(logrus.Fields{"selected": t.Selected, "edit": t.Edit, "tip": t.Form.ActiveWidget.tip}).Info("FocusActiveWidget")
	t.Table.ScrollTo(t.Selected)
	t.Table.Refresh()
	tip := t.Form.ActiveWidget.tip
	if strings.HasPrefix(tip, "float") {
		tip = "float"
	}
	if t.Edit {
		switch tip {
		case "string", "float", "date":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.ce)
		case "bool":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.ti)
		case "enum":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.sel)
		case "table":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.t)
		}
	} else {
		t.Form.W.Canvas().Focus(t.Form.ActiveWidget.t)
	}
}
