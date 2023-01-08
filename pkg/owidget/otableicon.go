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
	t.t.TypedKey(ev)
}

// FocusGained  fyne.Focusable
func (t *tappableIcon) FocusGained() {
}

func (t *tappableIcon) TypedRune(r rune) {
}

// FocusLost  fyne.Focusable
func (t *tappableIcon) FocusLost() {

}
