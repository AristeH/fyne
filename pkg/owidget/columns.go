package owidget

import (
	"otable/data"
	"strconv"

	"github.com/sirupsen/logrus"
)

// columnstyle - стиль колонки
type ColumnStyle struct {
	id      string  // id column
	name    string  // header
	formula string  // calculated expressions
	Width   float32 // ширина столбца
	BGcolor string  // цвет фона
	color   string  // цвет текста
	tip     string  // тип колонки(id, float, string, enum, date)
	visible bool    // видимость
	edit    bool    // редактирование колонки
	order   int16   // column output order
}

// fillcolumns - filling in columns from incoming data
func (t *OTable) fillcolumns(d data.GetData) {
	columns := len(d.DataDesciption[0]) // количесто колонок таблицы
	Log.WithFields(logrus.Fields{
		"form":    t.ID,
		"columns": columns,
		"event":   "fillcolumns()",
	}).Info("Columns")
	//инициализация стиля
	t.ColumnStyle = make(map[string]*ColumnStyle)
	//ширина символа

	for i := 0; i < columns; i++ {
		// Log.WithFields(logrus.Fields{"columns": d.Data[0][i]}).Info("columns")
		cs := ColumnStyle{}
		cs.name = d.Data[0][i]
		cs.id = d.DataDesciption[0][i]
		cs.BGcolor = "rowcolor" // индивидуальный цвет столбца фон
		cs.color = ""
		cs.formula = d.DataDesciption[3][i] // индивидуальный цвет текста столбца
		cs.tip = d.DataDesciption[1][i]
		p, _ := strconv.Atoi(d.DataDesciption[2][i]) //ширина столбца в символах
		cs.Width = float32(p)                        // ширина колонки
		cs.visible = true                            // видимость столбца
		cs.edit = true                               // редактируемость столбца
		t.ColumnStyle[cs.name] = &cs
	}
	defer Log.WithFields(logrus.Fields{
		"columns": t.ColumnStyle,
		"event":   "Finish",
	}).Info("Columns")
}

func (t *OTable) getcolumnid(id string) *ColumnStyle {
	return t.ColumnStyle[id]
}
