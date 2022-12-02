package owidget

import (
	"fmt"
	"math"
	"otable/data"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

// TableStyle - стиль таблицы
type TableStyle struct {
	BGColor       string // Цвет фона
	RowAlterColor string // Цвет строки четной
	HeaderColor   string // Цвет шапки
	RowColor      string // Цвет строки нечетной
}

// oTable - таблица
type OTable struct {
	widget.BaseWidget

	ID          string   // имя таблицы уникальное в пределах формы
	Form        FormData //  формa владелец таблицы
	ColumnStyle map[string]*ColumnStyle
	TabStyle    TableStyle          // стиль таблицы
	Data        map[string][]string // данные таблицы
	DataV       [][]string          // отображаемые данные(сортировка, фильтр)
	Table       *widget.Table
	Header      *widget.Table
	Properties  *OTable
	Selected    widget.TableCellID
	Edit        bool //
	wb          map[*widget.Button]int
}

func (t *OTable) MakeTable(d data.GetData) {
	Log.WithFields(logrus.Fields{"form": t.Form, "event": "MakeTable"}).Info("table")
	colColumns := len(d.DataDesciption[0])
	t.fillcolumns(d)
	Log.WithFields(logrus.Fields{"form": t.ID, "event": "MakeTableFillColumns"}).Info("table")
	colV := 0
	for i := 0; i < colColumns; i++ {
		b := strings.HasPrefix(d.Data[0][i], "id_")
		if !b {
			colV++
		}
	}
	t.Data = make(map[string][]string)
	t.DataV = make([][]string, len(d.Data))

	for i := 0; i < len(d.Data); i++ {
		data := make([]string, colColumns)
		datav := make([]string, colV)
		v := 0
		// l.WithFields(logrus.Fields{"i": i}).Info("mt")
		for j := 0; j < colColumns; j++ {
			b := strings.HasPrefix(d.DataDesciption[0][j], "id_")
			if !b {

				datav[v] = d.Data[i][j]
				v++
			} else {
				// Log.WithFields(logrus.Fields{"column j": j, "datadescription": d.DataDesciption[0][j]}).Info("filldata")

			}
			data[j] = d.Data[i][j]
		}

		t.Data[d.Data[i][0]] = data

		// l.WithFields(logrus.Fields{"j": len(t.DataV), "i": i}).Info("mt")
		t.DataV[i] = datav
		// Log.WithFields(logrus.Fields{"v": v}).Info("mt")
	}

	t.TabStyle.RowAlterColor = "RowAlterColor"
	t.TabStyle.HeaderColor = "HeaderColor"
	t.TabStyle.RowColor = "RowColor"
	t.Selected = widget.TableCellID{}
	Log.WithFields(logrus.Fields{
		"form":     "Table",
		"event":    "Fill data",
		"selected": t.Selected,
	}).Info("table")
	t.MakeTableLabel()
}

// not work
func (h *OTable) Scrolled(event *fyne.ScrollEvent) {
	fmt.Println(event.Position, event.AbsolutePosition)
	Log.WithFields(logrus.Fields{"rows": event}).Info("ScrollEvent")

}

// https://github.com/PaulWaldo/fyne-headertable
type headerTableRenderer struct {
	headerTable *OTable
	container   *fyne.Container
}

func (h *OTable) CreateRenderer() fyne.WidgetRenderer {
	return headerTableRenderer{
		headerTable: h,
		// container:   container.NewBorder(h.Header, nil, nil, nil, h.Table),
		container: container.NewBorder(nil, nil, nil, nil, h.Table),
	}
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
