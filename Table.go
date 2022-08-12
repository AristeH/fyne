package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"sort"
	"strconv"
)

var (
	floatValidator = validation.NewRegexp("[+-]?([0-9]*[.])?[0-9]+", "Не правильное число")
	emailValidator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "Не правильный email")
	intValidator   = validation.NewRegexp(`[+-]?[0-9]*$`, "Не целое число")
	dateValidator  = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "Не правильная дата")
	//emptyValidator = validation.NewRegexp(`^.+$`, "Поле не может быть пустым")
	//yearValidator  = validation.NewRegexp(`^[0-9]{4}$`, "Год содержит только 4 цифры.")
)

type enterEntry struct {
	IDForm  string
	IDTable string
	widget.Entry
	//	fyne.WidgetRenderer
	//	canvas.Text
}

type ColumnStyle struct {
	ID      string
	Name    string  //Заголовок столбца
	Format  string  //Форматированный вывод
	Width   float32 //Ширина столбца
	BGColor string  // Цвет фона
	Type    string
}

type TabStyle struct {
	ID            string
	Name          string
	BGColor       string // Цвет фона
	RowAlterColor string // Цвет строки четной
	HeaderColor   string // Цвет текста
	RowColor      string // Цвет строки нечетной
	Font          string // Шрифт

}

type TableOtoko struct {
	ID          string
	IDForm      string
	ColumnStyle []ColumnStyle
	TabStyle    TabStyle
	Data        [][]string
	Tool        *widget.Toolbar
	Table       *widget.Table
	Header      *fyne.Container
	Properties  *TableOtoko
	wl          map[*widget.Label]widget.TableCellID
	we          map[*enterEntry]widget.TableCellID
	wc          map[*widget.Check]widget.TableCellID
	wb          map[*widget.Button]int
}

func (t *TableOtoko) CreateHeader() {
	t.Header = container.New(&ToolButton{IDForm: t.IDForm, IDTable: t.ID})
	for col, value := range t.ColumnStyle {
		d := widget.NewButtonWithIcon(value.Name, nil, nil)
		d.OnTapped = func() {
			c := appValues[t.IDForm].Table[t.ID].wb[d]
			sort.Slice(t.Data, func(i, j int) bool { return t.Data[i][c] < t.Data[j][c] })
			t.Table.Refresh()
		}

		t.Header.Add(d)
		appValues[t.IDForm].Table[t.ID].wb[d] = col
	}

}

// ////////////////////////////////////////////

func (m *enterEntry) onEnter() {
	fmt.Println(m.Entry.Text)
	i := appValues[m.IDForm].Table[m.IDTable].we[m]
	appValues[m.IDForm].Table[m.IDTable].Data[i.Row][i.Col] = m.Entry.Text

}

func newEnterEntry() *enterEntry {
	entry := &enterEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func getValidator(t string) fyne.StringValidator {
	switch t {
	case "date":
		return dateValidator
	case "email":
		return emailValidator
	case "float":
		return floatValidator
	case "int":
		return intValidator
	default:
		return nil
	}
}

func getValue(s string, t string) string {
	switch t {
	case "float":
		if s, err := strconv.ParseFloat(s, 32); err == nil {
			return fmt.Sprintf("%15f", s)
		}
	case "int":
		if s, err := strconv.Atoi(s); err == nil {
			return fmt.Sprintf("%3d", s)
		}
	default:
		return s
	}
	return ""
}

func (t *TableOtoko) MakeTable() {

	t.Table = widget.NewTable(
		func() (int, int) {
			rows := len(t.Data)
			columns := len(t.Data[0])
			return rows, columns
		},
		func() fyne.CanvasObject {
			con := container.NewHBox()
			con.Layout = layout.NewMaxLayout()
			con.Add(widget.NewLabel(""))

			check := widget.NewCheck("", nil)
			check.OnChanged = func(b bool) {
				i := appValues[t.IDForm].Table[t.ID].wc[check]
				if check.Checked {
					t.Data[i.Row][i.Col] = "1"
				} else {
					t.Data[i.Row][i.Col] = "0"
				}

				newTableCellID := widget.TableCellID{Col: i.Col, Row: i.Row + 1}
				t.Table.ScrollTo(newTableCellID)
				println(i.Row)
				for key, value := range appValues[t.IDForm].Table[t.ID].wc {
					if value == newTableCellID {
						appValues[t.IDForm].W.Canvas().Focus(key)
						break
					}
				}
			}
			con.Add(check)
			entry := newEnterEntry()

			entry.IDForm = t.IDForm
			entry.IDTable = t.ID
			con.Add(entry)

			return container.New(layout.NewMaxLayout(),
				canvas.NewRectangle(color.Gray{Y: 250}),
				con,
			)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			var label *widget.Label
			var ic *widget.Check
			var entry *enterEntry

			box := o.(*fyne.Container)
			rect := box.Objects[0].(*canvas.Rectangle)
			if i.Row%2 == 0 {
				rect.FillColor = MapColor[t.TabStyle.RowAlterColor]
			} else {
				rect.FillColor = MapColor[t.TabStyle.RowColor]
			}
			if val, ok := MapColor[t.ColumnStyle[i.Col].BGColor]; ok {
				rect.FillColor = mix(val, rect.FillColor)
			}

			cont := box.Objects[1].(*fyne.Container)

			label = cont.Objects[0].(*widget.Label)

			ic = cont.Objects[1].(*widget.Check)

			entry = cont.Objects[2].(*enterEntry)
			label.Hidden = true
			ic.Hidden = true
			entry.Hidden = true
			switch t.ColumnStyle[i.Col].Type {
			case "bool":
				appValues[t.IDForm].Table[t.ID].wc[ic] = i
				if t.Data[i.Row][i.Col] == "1" {
					ic.Checked = true
				} else {
					ic.Checked = false
				}
				ic.Refresh()
				ic.Hidden = false
			case "string":
				entry.SetText(getValue(t.Data[i.Row][i.Col], t.ColumnStyle[i.Col].Format))
				entry.Validator = getValidator(t.ColumnStyle[i.Col].Format)
				appValues[t.IDForm].Table[t.ID].we[entry] = i
				entry.Hidden = false
			case "float":
				entry.SetText(getValue(t.Data[i.Row][i.Col], t.ColumnStyle[i.Col].Format))
				entry.Validator = getValidator(t.ColumnStyle[i.Col].Format)
				appValues[t.IDForm].Table[t.ID].we[entry] = i
				entry.Hidden = false
			case "int":
				entry.SetText(getValue(t.Data[i.Row][i.Col], t.ColumnStyle[i.Col].Format))
				entry.Validator = getValidator(t.ColumnStyle[i.Col].Format)
				appValues[t.IDForm].Table[t.ID].we[entry] = i
				entry.Hidden = false
			default:
				label.SetText(getValue(t.Data[i.Row][i.Col], t.ColumnStyle[i.Col].Format))
				appValues[t.IDForm].Table[t.ID].wl[label] = i
				label.Hidden = false
			}

		})
	for ic, v := range t.ColumnStyle {
		t.Table.SetColumnWidth(ic, v.Width)
	}
	t.Table.OnSelected = func(id widget.TableCellID) {
		fmt.Printf("i.Col: %v\n", id.Col)
	}

	t.Tool = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {

			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			w1 := myApp.NewWindow("TabContainer Widget")
			table := t.Properties
			if table == nil {
				return
			}
			content := container.NewBorder(
				container.NewVBox(
					table.Tool,
					widget.NewSeparator(),
					table.Header,
					widget.NewSeparator(),
				),
				nil,
				nil,
				nil,
				table.Table,
			)

			w1.Resize(fyne.NewSize(600, 800))
			w1.SetContent(container.NewMax(content))
			w1.Show()
		}))

	t.CreateHeader()

}

func (m *enterEntry) KeyDown(key *fyne.KeyEvent) {
	t := appValues[m.IDForm].Table[m.IDTable]
	switch key.Name {
	case fyne.KeyReturn:
		m.onEnter()
	case "KP_Enter":
		m.onEnter()
		i := t.we[m]
		t.ColumnStyle[i.Col].Width = 40
	case "Down":
		i := t.we[m]
		newTableCellID := widget.TableCellID{Col: i.Col, Row: i.Row + 1}
		t.Table.ScrollTo(newTableCellID)
		for key, value := range appValues[m.IDForm].Table[m.IDTable].we {
			if value == newTableCellID {
				appValues[m.IDForm].W.Canvas().Focus(key)
				break
			}
		}
	case "Up":
		i := t.we[m]
		newTableCellID := widget.TableCellID{Col: i.Col, Row: i.Row - 1}
		t.Table.ScrollTo(newTableCellID)
		for key, value := range t.we {
			if value == newTableCellID {
				appValues[m.IDForm].W.Canvas().Focus(key)
				break
			}
		}
	default:
		m.Entry.KeyDown(key)
		fmt.Printf("Key %v pressed\n", key.Name)
	}
}

func (m *enterEntry) KeyUp(key *fyne.KeyEvent) {
	fmt.Printf("Key %v released\n", key.Name)
}

// somehow this catches right click. How?
func (mc *enterEntry) TappedSecondary(pe *fyne.PointEvent) {
	fmt.Printf("1")
}

func (mc *enterEntry) DoubleTapped(pe *fyne.PointEvent) {
	ind := appValues["main"].Table["tovar"].we[mc]
	if ind.Col == 5 {
		choiceFromList(Names, appValues["main"].Table["tovar"], mc)
	}
	if ind.Col == 3 {
		var Types = []string{
			"label",
			"string",
			"bool",
			"float",
		}
		choiceFromList(Types, appValues["main"].Table["tovar"], mc)
	}
	if ind.Col == 4 {
		//if s, err := strconv.ParseFloat(mc.Text, 32); err == nil {
		appValues["main"].Table["tovar"].ColumnStyle[ind.Row].Width = float32(20) // 3.1415927410125732
		appValues["main"].Table["tovar"].Table.SetColumnWidth(ind.Row, 20)
		//	}
	}

}
func (mc *enterEntry) OnChanged() {

	fmt.Printf("1")

}

type ToolButton struct {
	IDForm  string
	IDTable string
}

func (d *ToolButton) MinSize(objects []fyne.CanvasObject) fyne.Size {
	TO := appValues[d.IDForm].Table[d.IDTable]
	w, h := float32(0), float32(0)
	for i, o := range objects {
		childSize := o.MinSize()
		o.Resize(fyne.NewSize(TO.ColumnStyle[i].Width, childSize.Height))
		w += TO.ColumnStyle[i].Width
		h = childSize.Height
	}
	return fyne.NewSize(w, h)
}

func (d *ToolButton) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	TO := appValues[d.IDForm].Table[d.IDTable]
	for i, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(TO.ColumnStyle[i].Width, 0))

	}
}

func TableInitProperties(t *TableOtoko) *TableOtoko {
	colColumns := 6
	colRows := len(t.ColumnStyle)
	data := make([][]string, colRows)
	cs := t.ColumnStyle
	for i := 0; i < colRows; i++ {
		data[i] = make([]string, colColumns)
		data[i][0] = cs[i].ID
		data[i][1] = cs[i].Name
		data[i][2] = cs[i].Type
		data[i][3] = cs[i].Format
		data[i][4] = fmt.Sprintf("%5.0f", cs[i].Width) // 3.1415927410125732
		data[i][5] = "blanchedalmond"
	}
	var TO = TableOtoko{}
	csn := ColumnStyle{}
	csn.Name = "Col ID"
	csn.Type = "label"
	csn.Width = 50
	TO.ColumnStyle = append(TO.ColumnStyle, csn)
	csn = ColumnStyle{}
	csn.Name = "Name"
	csn.Type = "string"
	csn.Width = 100
	TO.ColumnStyle = append(TO.ColumnStyle, csn)
	csn = ColumnStyle{}
	csn.Name = "Type"
	csn.Type = "string"
	csn.Width = 80
	TO.ColumnStyle = append(TO.ColumnStyle, csn)
	csn = ColumnStyle{}
	csn.Name = "Format"
	csn.Type = "string"
	csn.Width = 100
	TO.ColumnStyle = append(TO.ColumnStyle, csn)
	csn = ColumnStyle{}
	csn.Name = "Width"
	csn.Type = "float"
	csn.Width = 80
	TO.ColumnStyle = append(TO.ColumnStyle, csn)
	csn = ColumnStyle{}
	csn.Name = "Color"
	csn.Width = 150
	csn.Type = "string"
	TO.ColumnStyle = append(TO.ColumnStyle, csn)

	ts := TabStyle{}

	ts.RowAlterColor = "lightgray"
	ts.HeaderColor = "lightslategray"
	ts.RowColor = "slategray"
	TO.TabStyle = ts
	TO.Data = data
	TO.ID = "prop"
	TO.IDForm = "mainprop"

	TO.wb = make(map[*widget.Button]int)
	TO.wc = make(map[*widget.Check]widget.TableCellID)
	TO.we = make(map[*enterEntry]widget.TableCellID)
	TO.wl = make(map[*widget.Label]widget.TableCellID)

	return &TO
}
