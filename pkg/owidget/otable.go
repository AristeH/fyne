package owidget

import (
	"fmt"
	"otable/data"
	"strings"

	"fyne.io/fyne/v2"
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
	Tool        *widget.Toolbar
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

func (t *OTable) TableInitProperties() *OTable {
	colColumns := 6
	colRows := len(t.ColumnStyle)
	data := make([][]string, colRows)
	cs := t.ColumnStyle
	i := 0
	for _, v := range cs {

		data[i] = make([]string, colColumns)
		data[i][0] = v.id
		data[i][1] = v.name
		data[i][2] = v.tip
		data[i][3] = v.formula
		data[i][4] = v.color
		data[i][5] = fmt.Sprintf("%v", v.Width)
		i++
	}
	var TO = OTable{}
	TO.ColumnStyle = make(map[string]*ColumnStyle)
	csn := ColumnStyle{}
	csn.name = "Col ID"
	csn.tip = "label"
	csn.Width = 50
	TO.ColumnStyle[csn.name] = &csn

	csn = ColumnStyle{}
	csn.name = "Name"
	csn.tip = "string"
	csn.Width = 100
	TO.ColumnStyle[csn.name] = &csn

	csn = ColumnStyle{}
	csn.name = "Type"
	csn.tip = "string"
	csn.Width = 80
	TO.ColumnStyle[csn.name] = &csn

	csn.name = "Format"
	csn.tip = "string"
	csn.Width = 100
	TO.ColumnStyle[csn.name] = &csn

	csn = ColumnStyle{}
	csn.name = "Width"
	csn.tip = "float"
	csn.Width = 80
	TO.ColumnStyle[csn.name] = &csn

	csn = ColumnStyle{}
	csn.name = "Color"
	csn.Width = 150
	csn.tip = "string"
	TO.ColumnStyle[csn.name] = &csn

	ts := TableStyle{}

	ts.RowAlterColor = "lightgray"
	ts.HeaderColor = "lightslategray"
	ts.RowColor = "slategray"
	TO.TabStyle = ts
	TO.DataV = data
	TO.ID = "prop"
	Log.WithFields(logrus.Fields{"TO.TabStyle": TO.DataV}).Info("prop")

	return &TO
}
