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
	l := logger.GetLog()
	l.WithFields(logrus.Fields{"DocumentCreateIcon": "GetToolBar"}).Info("GetToolBar")

	t.Tool = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			l.WithFields(logrus.Fields{"DocumentCreateIcon": "DocumentCreateIcon"}).Info("GetToolBar")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fd := PutListForm("TableProp", "Tablerop")
			g := t.properties()
			l.WithFields(logrus.Fields{"Properties": len(g.Data)}).Info("GetToolBar")

			table := fd.NewOTable("invoice", *g)

			l.WithFields(logrus.Fields{"Properties len  dv": len(table.DataV)}).Info("GetToolBar")
			w := fd.W
			w.Resize(fyne.NewSize(1200, 400))

			w.SetContent(container.NewMax(table))
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
					c := NewCompletionEntry([]string{})

					t.Form.ActiveWidget.tip = "string"
					t.Form.ActiveWidget.ce = c
					if strings.HasPrefix(tip, "id_") { //id другой таблицы
						en := string(mystr[0 : k-4])
						entry := t.MakeTappable(en, tip, FillColor)
						c := container.NewHSplit(entry, newTappableIcon(theme.SearchIcon()))
						box.Objects[0] = c
					} else {
						if strings.HasPrefix(tip, "float") {
							tip = "float"
						}
						switch tip {
						case "float", "int", "string":
							c.Text = t.DataV[i.Row][i.Col]
							c.t = t
							box.Objects[0] = container.New(layout.NewMaxLayout(), c)
						case "bool":
							t.Form.ActiveWidget.tip = "bool"
							ic := newTappableIcon(theme.CheckButtonIcon())
							if t.DataV[i.Row][i.Col] == "0" {
								ic = newTappableIcon(theme.CheckButtonCheckedIcon())
							}
							ic.t = t
							t.Form.ActiveWidget.ti = ic
							box.Objects[0] = container.New(layout.NewMaxLayout(), ic)
						case "enum":
							c.SetOptions([]string{"product", "service"})
							c.OnChanged = func(s string) {
								if len(s) < 3 {
									c.HideCompletion()
									return
								}
								c.ShowCompletion()
							}
							c.t = t
							c.Text = t.DataV[i.Row][i.Col]
							box.Objects[0] = container.New(layout.NewMaxLayout(), c)
						}
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
		for i, _ := range t.DataV[0] {
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
		for i, _ := range t.DataV[0] {
			if strings.HasPrefix(t.ColumnStyle[t.DataV[0][i]].tip, "float") || t.ColumnStyle[t.DataV[0][i]].tip == "int" {
				v := t.ColumnStyle[t.DataV[0][i]].id
				t.DataV[t.Selected.Row][i] = fmt.Sprintf("%.2f", compiled.Get(v).Float())
				fmt.Println(v, t.DataV[t.Selected.Row][i])
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
				case 0:
					col.sort = 1
					t.sortUp()
				case 1:
					col.sort = 2
					t.sortDown()
				case 2:
					col.sort = 0
				}
			}
		} else {
			t.ExecuteFormula()
			if t.Edit {
				t.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
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
			tc := widget.TableCellID{Col: i.Col, Row: i.Row - 1}
			t.Selected = tc
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
		case "string", "float":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.ce)
		case "bool":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.ti)
		case "table":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.t)
		}
	} else {
		t.Form.W.Canvas().Focus(t.Form.ActiveWidget.t)
	}
}
