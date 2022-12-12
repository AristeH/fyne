package owidget

import (
	"otable/pkg/logger"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logger.GetLog()
}

func (t *OTable) getColorCell(i widget.TableCellID) *CellColor {
	c := CellColor{}
	c.Color = MapColor["black"]
	//цвет фона ячейки
	if i.Row == 0 {
		c.BGcolor = MapColor[t.TabStyle.HeaderColor]
	} else if i.Row%2 == 0 {
		c.BGcolor = MapColor[t.TabStyle.RowAlterColor]

	} else {
		c.BGcolor = MapColor[t.TabStyle.RowColor]
	}
	// выделенный столбец
	col := t.ColumnStyle[t.DataV[0][i.Col]]
	if val, ok := MapColor[col.bgcolor]; ok {
		c.BGcolor = mix(val, c.BGcolor)
	}
	// individual
	id, ok := t.CellColor[strconv.Itoa(i.Row)+";"+strconv.Itoa(i.Col)]
	if ok {
		c = *id
	}

	// выделенная ячейка
	if i == t.Selected {
		c.BGcolor = MapColor["Selected"]
	}

	return &c
}

func (t *OTable) MakeTableLabel() {

	rows := len(t.DataV)
	columns := len(t.DataV[0])

	Log.WithFields(logrus.Fields{"rows": rows, "columns": columns}).Info("MakeTableLabel")
	t.Header = widget.NewTable(
		// Dimensions (rows, cols)
		func() (int, int) { return 1, columns },
		// Default value
		func() fyne.CanvasObject { return widget.NewLabel("the content") },
		// Cell values
		func(cellID widget.TableCellID, o fyne.CanvasObject) {
			colst := t.ColumnStyle[t.DataV[0][cellID.Col]]
			l := o.(*widget.Label)
			l.SetText(colst.name)
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
				k = int(col.Width) - 1
			}
			en := string(mystr[0:k])
			if tip == "bool" {
				rec := canvas.NewRectangle(FillColor.BGcolor)
				image := canvas.NewImageFromResource(theme.CheckButtonCheckedIcon())
				if t.DataV[i.Row][i.Col] == "1" {
					image = canvas.NewImageFromResource(theme.CheckButtonIcon())
				}
				box.Objects[0] = container.New(layout.NewMaxLayout(), rec, image)
			} else {
				box.Objects[0] = t.MakeTappable(en, tip, FillColor)
			}
			// активная ячейка
			if i == t.Selected {
				t.Form.ActiveWidget.tip = ""
				if i.Row > 0 && t.Edit {
					c := NewCompletionEntry([]string{})
					t.Form.ActiveWidget.tip = "string"
					t.Form.ActiveWidget.ce = c
					if strings.HasPrefix(tip, "id_") { //id  другой таблицы
						en := string(mystr[0 : k-3])
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
							if t.DataV[i.Row][i.Col] == "1" {
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

	//	entry.Objects[0].(*oLabel).CellID = &i

	//entry.parent = t
	si := fyne.MeasureText("шii", 24, fyne.TextStyle{})
	rec := canvas.NewRectangle(c.BGcolor)
	rec.SetMinSize(si)

	entry.Resize(si)

	return container.New(layout.NewMaxLayout(), rec, entry)
}

// Implements: fyne.Focusable
func (t *OTable) TypedKey(ev *fyne.KeyEvent) {
}

func (n *OTable) TypedRune(r rune) {
	Log.WithFields(logrus.Fields{"entry.text": r}).Info("onEnter ")
}

// Implements: fyne.Focusable
func (t *OTable) FocusLost() {
}

// Implements: fyne.Focusable
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
	}
}
