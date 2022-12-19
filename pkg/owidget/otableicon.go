package owidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

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
	Log.WithFields(logrus.Fields{"Tapped": ev}).Info("1tappableIcon")
}
func (t *tappableIcon) KeyDown(key *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"rows": key}).Info("TappedTappableIcon")
}

// TypedKey  fyne.Focusable
func (t *tappableIcon) TypedKey(ev *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"tappableIcon": ev}).Info("TypedKey")
	otab := t.t
	i := t.t.Selected
	switch ev.Name {
	case "Return":
		if otab.Edit {
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
		if i.Row > 1 {
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
	case fyne.KeySpace:
		if otab.Edit {
			if otab.DataV[i.Row][i.Col] == "1" {
				otab.DataV[i.Row][i.Col] = "0"
			} else {
				otab.DataV[i.Row][i.Col] = "1"
			}
		}
	}
	otab.FocusActiveWidget()
}

// FocusGained  fyne.Focusable
func (t *tappableIcon) FocusGained() {
}

func (t *tappableIcon) TypedRune(r rune) {
}

// FocusLost  fyne.Focusable
func (t *tappableIcon) FocusLost() {

}
