package owidget

import (
	"fmt"
	"image/color"
	"otable/data"
	"otable/pkg/logger"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

// CellColor - цвета для ячейки
type CellColor struct {
	Color   color.RGBA
	BGcolor color.RGBA
}

// oTable - таблица
type OTable struct {
	widget.BaseWidget
	ID          string                  // имя таблицы уникальное в пределах формы
	Form        FormData                //  формa владелец таблицы
	ColumnStyle map[string]*ColumnStyle // описание колонок
	TabStyle    TableStyle              // стиль таблицы
	Data        map[string][]string     // данные таблицы
	DataV       [][]string              // отображаемые данные(сортировка, фильтр)
	Table       *widget.Table
	Header      *widget.Table
	//	Footer      *widget.Table
	//	left      *widget.Table
	Properties *OTable               // таблица для редактирования описания  колонок
	Tool       *widget.Toolbar       // командная панель  таблицы
	Selected   widget.TableCellID    // выделенная ячейка таблицы
	Edit       bool                  // редактируемость таблицы
	CellColor  map[string]*CellColor //  индивидуальный массив отображения ячеек
	// wb         map[*widget.Button]int
}

func (t *OTable) MakeTableData(d data.GetData) {

	colColumns := len(d.DataDesciption[0])
	t.fillcolumns(d)
	Log.WithFields(logrus.Fields{"form": t.Form.ID, "event": "fillcolumns"}).Info("MakeTableData")
	colV := 0 //количество видимых столбцов для пользователя
	for i := 0; i < colColumns; i++ {
		b := strings.HasPrefix(d.Data[0][i], "id_") //исключим столбцы с типом ID
		if !b {
			colV++
		}
	}
	t.Data = make(map[string][]string)
	t.DataV = make([][]string, len(d.Data))

	for i := 0; i < len(d.Data); i++ {
		datad := make([]string, colColumns)
		datav := make([]string, colV)
		v := 0
		// l.WithFields(logrus.Fields{"i": i}).Info("mt")
		for j := 0; j < colColumns; j++ {
			b := strings.HasPrefix(d.DataDesciption[0][j], "id_")
			if !b {
				datav[v] = d.Data[i][j]
				v++
			}
			datad[j] = d.Data[i][j]
		}
		t.Data[d.Data[i][0]] = datad
		// l.WithFields(logrus.Fields{"j": len(t.DataV), "i": i}).Info("mt")
		t.DataV[i] = datav
		// Log.WithFields(logrus.Fields{"v": v}).Info("mt")
	}
	t.TabStyle.RowAlterColor = "RowAlterColor"
	t.TabStyle.HeaderColor = "HeaderColor"
	t.TabStyle.RowColor = "RowColor"
	t.Selected = widget.TableCellID{}
	Log.WithFields(logrus.Fields{"form": t.Form.ID, "event": "Fill data"}).Info("MakeTable")
	t.MakeTableLabel()
}

// not work
func (t *OTable) Scrolled(event *fyne.ScrollEvent) {
	fmt.Println(event.Position, event.AbsolutePosition)
	Log.WithFields(logrus.Fields{"rows": event}).Info("ScrollEvent")

}

func (t *OTable) GetToolBar() {
	l := logger.GetLog()
	l.WithFields(logrus.Fields{"DocumentCreateIcon": "GetToolBar"}).Info("GetToolBar")

	t.Tool = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			l.WithFields(logrus.Fields{"DocumentCreateIcon": "DocumentCreateIcon"}).Info("GetToolBar")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fd := PutListForm("TableProp", "Tablerop")
			g := t.TableInitProperties()
			l.WithFields(logrus.Fields{"Properties": len(g.Data)}).Info("GetToolBar")

			table := fd.NewOTable("invoice", *g)

			l.WithFields(logrus.Fields{"Properties len  dv": len(table.DataV)}).Info("GetToolBar")
			w := fd.W
			w.Resize(fyne.NewSize(1200, 400))
			w.SetContent(container.NewMax(table))
			w.Show()
		}))

}

func (t *OTable) TableInitProperties() *data.GetData {
	colColumns := 10
	colRows := len(t.ColumnStyle)
	datag := make([][]string, colRows)
	cs := t.ColumnStyle
	i := 0
	for _, v := range cs {
		datag[i] = make([]string, colColumns)
		datag[i][0] = v.id
		datag[i][1] = v.name                     // заголовок
		datag[i][2] = v.tip                      // тип столбца
		datag[i][3] = v.formula                  // Формула
		datag[i][4] = v.color                    // цвет теста столбца
		datag[i][5] = v.BGcolor                  // цвет фона столбца
		datag[i][6] = fmt.Sprintf("%v", v.Width) // ширина столбца в символах
		if v.visible {                           // видимость столбца
			datag[i][7] = "1"
		} else {
			datag[i][7] = "0"
		}
		if v.edit { // видимость столбца
			datag[i][8] = "1"
		} else {
			datag[i][8] = "0"
		}
		datag[i][9] = strconv.FormatInt(int64(v.order), 2) // порядок вывода
		i++
	}

	datag[0][0] = "id"
	datag[0][1] = "Header"
	datag[0][2] = "Type"
	datag[0][3] = "formula"
	datag[0][4] = "Color"
	datag[0][5] = "BGcolor"
	datag[0][6] = "Width"
	datag[0][7] = "visible"
	datag[0][8] = "edit"
	datag[0][9] = "order"

	Log.WithFields(logrus.Fields{"form1": t.ID, "datag": len(datag)}).Info("TableInitProperties")

	// инициализация описания данных таблицы
	datadescription := make([][]string, 4)
	for i := 0; i < 4; i++ {
		datadescription[i] = make([]string, colColumns)
	}

	// Name column
	datadescription[0][0] = "id"
	datadescription[0][1] = "Header"
	datadescription[0][2] = "formula"
	datadescription[0][3] = "Type"
	datadescription[0][4] = "Color"
	datadescription[0][5] = "BGColor"
	datadescription[0][6] = "Width"
	datadescription[0][7] = "visible"
	datadescription[0][8] = "edit"
	datadescription[0][9] = "order"

	//  Type column
	datadescription[1][0] = "string"
	datadescription[1][1] = "string"
	datadescription[1][2] = "string"
	datadescription[1][3] = "enum"
	datadescription[1][4] = "id_color"
	datadescription[1][5] = "id_color"
	datadescription[1][6] = "int"
	datadescription[1][7] = "bool"
	datadescription[1][8] = "bool"
	datadescription[1][9] = "int"

	// Width column
	datadescription[2][0] = "15"
	datadescription[2][1] = "15"
	datadescription[2][2] = "20"
	datadescription[2][3] = "10"
	datadescription[2][4] = "10"
	datadescription[2][5] = "10"
	datadescription[2][6] = "6"
	datadescription[2][7] = "4"
	datadescription[2][8] = "4"
	datadescription[2][9] = "3"

	//Formula column
	datadescription[3][0] = ""
	datadescription[3][1] = ""
	datadescription[3][2] = ""
	datadescription[3][3] = ""
	datadescription[3][4] = ""
	datadescription[3][5] = ""
	datadescription[3][6] = ""
	datadescription[3][7] = ""
	datadescription[3][8] = ""
	datadescription[3][9] = ""

	f := data.GetData{}
	f.Data = datag
	f.DataDesciption = datadescription
	return &f

}
