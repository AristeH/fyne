package owidget

import (
	"image/color"
	"otable/pkg/logger"
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
		func() fyne.CanvasObject { return t.MakeTappable("", "", color.Black, color.Opaque) },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			var FillColor color.Color
			box := o.(*fyne.Container)
			//цвет фона ячейки
			if i.Row == 0 {
				FillColor = MapColor[t.TabStyle.HeaderColor]
			} else if i.Row%2 == 0 {
				FillColor = MapColor[t.TabStyle.RowAlterColor]
			} else {
				FillColor = MapColor[t.TabStyle.RowColor]
			}
			// выделенный столбец
			col := t.ColumnStyle[t.DataV[0][i.Col]]
			if val, ok := MapColor[col.bgcolor]; ok {
				FillColor = mix(val, FillColor)
			}
			// выделенная ячейка
			if i == t.Selected {
				FillColor = MapColor["Selected"]
			}
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
				rec := canvas.NewRectangle(FillColor)
				image := canvas.NewImageFromResource(theme.CheckButtonCheckedIcon())
				if t.DataV[i.Row][i.Col] == "1" {
					image = canvas.NewImageFromResource(theme.CheckButtonIcon())
				}

				box.Objects[0] = container.New(layout.NewMaxLayout(), rec, image)

			} else {
				box.Objects[0] = t.MakeTappable(en, tip, color.Black, FillColor)
			}
			// активная ячейка
			if i == t.Selected {
				FillColor = MapColor["Selected"]
				t.Form.ActiveWidget.tip = ""
				if i.Row > 0 && t.Edit {
					c := NewCompletionEntry([]string{})
					t.Form.ActiveWidget.tip = "string"
					t.Form.ActiveWidget.ce = c
					if strings.HasPrefix(tip, "id_") { //id  другой таблицы
						en := string(mystr[0 : k-3])
						entry := t.MakeTappable(en, tip, color.Black, FillColor)
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
							box.Objects[0] = container.New(layout.NewMaxLayout(),
								ic)

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
		t.Form.ActiveContainer = t
		t.FocusActiveWidget()
	}

}

func (t *OTable) MakeTappable(txt string, tip string, c color.Color, b color.Color) *fyne.Container {
	entry := canvas.NewText(txt, c)
	if strings.HasPrefix(tip, "float") {
		tip = "float"
	}

	switch tip {
	case "float", "int":
		entry.Alignment = fyne.TextAlignTrailing
	default:
		entry.Alignment = fyne.TextAlignLeading
	}

	//	entry.Objects[0].(*oLabel).CellID = &i

	//entry.parent = t
	si := fyne.MeasureText("шii", 24, fyne.TextStyle{})
	rec := canvas.NewRectangle(b)
	rec.SetMinSize(si)

	entry.Resize(si)

	return container.New(layout.NewMaxLayout(), rec, entry)
}

type tappableIcon struct {
	widget.Icon
	t *OTable
}

func newTappableIcon(res fyne.Resource) *tappableIcon {
	icon := &tappableIcon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)

	return icon
}

func (t *tappableIcon) Tapped(ev *fyne.PointEvent) {
	Log.WithFields(logrus.Fields{"Tapped": ev}).Info("tappableIcon")

}
func (t *tappableIcon) KeyDown(key *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"rows": key}).Info("TappedtappableIcon")
}

// Implements: fyne.Focusable
func (t *tappableIcon) TypedKey(ev *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"tappableIcon": ev}).Info("TypedKey")
	Log.WithFields(logrus.Fields{"s": t.t, "i": t.t.Selected}).Info("TypedKey")
	otab := t.t
	i := t.t.Selected
	switch ev.Name {
	case "Return":
		if otab.Edit {
			otab.DataV[i.Row][i.Col] = "1"
			otab.Selected = widget.TableCellID{Col: i.Col, Row: i.Row + 1}
		} else {
			otab.Edit = true
			otab.Selected = widget.TableCellID{Col: i.Col, Row: i.Row}
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
	case "Escape":
		otab.Edit = false
		otab.Form.ActiveWidget.tip = "table"
		otab.Form.ActiveWidget.t = otab
		otab.FocusActiveWidget()
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
	otab.FocusActiveWidget()

}

// Implements: fyne.Focusable
func (t *tappableIcon) FocusGained() {
	Log.WithFields(logrus.Fields{"tappableIcon": "fg"}).Info("TypedKey")
}

func (n *tappableIcon) TypedRune(r rune) {
	Log.WithFields(logrus.Fields{"entry.text": r}).Info("onEnter ")
}

// Implements: fyne.Focusable
func (t *tappableIcon) FocusLost() {
	Log.WithFields(logrus.Fields{"tappableIcon": "fl"}).Info("TypedKey")
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
			// Log.WithFields(logrus.Fields{"ActiveWidget": t.Form.ActiveWidget.ce}).Info("FocusActiveWidget")
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.ce)
		case "bool":
			Log.WithFields(logrus.Fields{"bool": t.Form.ActiveWidget.ti}).Info("FocusActiveWidget")
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.ti)
		case "table":
			t.Form.W.Canvas().Focus(t.Form.ActiveWidget.t)
		}
	}
}
