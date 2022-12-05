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
	Log.WithFields(logrus.Fields{"rows": key}).Info("TappedtappableIcon")
}

// Implements: fyne.Focusable
func (t *tappableIcon) TypedKey(ev *fyne.KeyEvent) {
	Log.WithFields(logrus.Fields{"tappableIcon": ev}).Info("TypedKey")

	otab := t.t
	i := t.t.Selected
	switch ev.Name {
	case "Return":
		if otab.Edit {
			// Log.WithFields(logrus.Fields{"t": otab.DataV[i.Row][i.Col]}).Info("TypedKey")
			//
			// if otab.DataV[i.Row][i.Col] == "1" {
			// 	otab.DataV[i.Row][i.Col] = "0"
			// } else {
			// 	otab.DataV[i.Row][i.Col] = "1"
			//
			// }
			// Log.WithFields(logrus.Fields{"t": otab.DataV[i.Row][i.Col]}).Info("TypedKey")
			//
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
	otab := t.t
	i := otab.Selected
	if otab.Edit {
		if otab.DataV[i.Row][i.Col] == "1" {
			otab.DataV[i.Row][i.Col] = "0"
		} else {
			otab.DataV[i.Row][i.Col] = "1"
		}
	}
	Log.WithFields(logrus.Fields{"tappableIcon": "fl"}).Info("TypedKey")
}
