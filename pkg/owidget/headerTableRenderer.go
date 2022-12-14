package owidget

// https://github.com/PaulWaldo/fyne-headertable
import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/sirupsen/logrus"
)

type headerTableRenderer struct {
	headerTable *OTable
	container   *fyne.Container
}

func (t *OTable) CreateRenderer() fyne.WidgetRenderer {
	Log.WithFields(logrus.Fields{"h.Tool": t.Tool}).Info("CreateRenderer")
	ht := headerTableRenderer{}
	ht.headerTable = t
	if t.Tool == nil {
		// container:   container.NewBorder(h.Header, nil, nil, nil, h.Table),
		ht.container = container.NewBorder(nil, nil, nil, nil, t.Table)
	} else {
		ht.container = container.NewBorder(t.Tool, nil, nil, nil, t.Table)

	}
	return ht
}

func (r headerTableRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		float32(math.Max(float64(r.headerTable.Table.MinSize().Width), float64(r.headerTable.Header.MinSize().Width))),
		r.headerTable.Table.MinSize().Height+r.headerTable.Header.MinSize().Height)
}

func (r headerTableRenderer) Layout(s fyne.Size) {
	r.container.Resize(s)
}

func (r headerTableRenderer) Destroy() {
}

func (r headerTableRenderer) Refresh() {
	r.container.Refresh()
}

func (r headerTableRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}
