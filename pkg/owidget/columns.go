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
	sort    int8    // 0 - нет сортировки, 1 - возрастание, 2-убывание
}

// fillcolumns - filling in columns from incoming data
func (t *OTable) fillColumns(d data.GetData) {
	columns := len(d.DataDescription[0]) // количесто колонок таблицы
	Log.WithFields(logrus.Fields{
		"form":    t.ID,
		"columns": columns,
		"event":   "fillColumns()",
	}).Info("Columns")
	//инициализация стиля
	t.ColumnStyle = make(map[string]*ColumnStyle)
	//ширина символа

	for i := 0; i < columns; i++ {
		// Log.WithFields(logrus.Fields{"columns": d.Data[0][i]}).Info("columns")
		cs := ColumnStyle{}
		cs.name = d.Data[0][i]
		cs.id = d.DataDescription[0][i]
		cs.BGcolor = "rowcolor" // индивидуальный цвет столбца фон
		cs.color = ""
		cs.formula = d.DataDescription[3][i] // индивидуальный цвет текста столбца
		cs.tip = d.DataDescription[1][i]
		p, _ := strconv.Atoi(d.DataDescription[2][i]) //ширина столбца в символах
		cs.Width = float32(p)                         // ширина колонки
		cs.visible = true                             // видимость столбца
		cs.edit = true                                // редактируемость столбца
		t.ColumnStyle[cs.name] = &cs
	}
	defer Log.WithFields(logrus.Fields{
		"columns": t.ColumnStyle,
		"event":   "Finish",
	}).Info("Columns")
}
