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
)

var (
	floatValidator = validation.NewRegexp("[+-]?([0-9]*[.])?[0-9]+", "Не правильное число")
	emailValidator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "Не правильный email")
	intValidator   = validation.NewRegexp(`[+-]?[0-9]*$`, "Не целое число")
	dateValidator  = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "Не правильная дата")
	emptyValidator = validation.NewRegexp(`^.+$`, "Поле не может быть пустым")
	yearValidator  = validation.NewRegexp(`^[0-9]{4}$`, "Год содержит только 4 цифры.")
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
	Name    string      //Заголовок столбца
	Format  string      //Форматированный вывод
	Width   float32     //Ширина столбца
	BGColor color.Color // Цвет фона
	Type    string
}

type TabStyle struct {
	ID            string
	Name          string
	BGColor       color.Color // Цвет фона
	RowAlterColor color.Color // Цвет строки четной
	HeaderColor   color.Color // Цвет текста
	RowColor      color.Color // Цвет строки нечетной
	Font          string      // Шрифт

}

type TableOtoko struct {
	ID          string
	IDForm      string
	ColumnStyle []ColumnStyle
	TabStyle    TabStyle
	Data        [][]string
	Edit        bool
	Tool        *widget.Toolbar
	Table       *widget.Table
	Header      *fyne.Container
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

func (t *TableOtoko) LoadTable(mes []byte) {
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
				rect.FillColor = t.TabStyle.RowAlterColor
			} else {
				rect.FillColor = t.TabStyle.RowColor
			}
			if t.ColumnStyle[i.Col].BGColor != nil {
				rect.FillColor = mix(t.ColumnStyle[i.Col].BGColor, rect.FillColor)
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
				entry.SetText(t.Data[i.Row][i.Col])

				appValues[t.IDForm].Table[t.ID].we[entry] = i
				entry.Hidden = false
			default:
				label.SetText(t.Data[i.Row][i.Col])
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
			log.Println("Display help")
		}))

	t.CreateHeader()

}

//func (e *enterEntry) BackgroundColor() color.Color {
//	return color.RGBA{R: 255, G: 20, B: 147, A: 255}
//}
//func (e *enterEntry) CreateRenderer() fyne.WidgetRenderer {
//	r := e.Entry.CreateRenderer()
//	return &enterEntry{r}
//}

func (m *enterEntry) KeyDown(key *fyne.KeyEvent) {
	t := appValues[m.IDForm].Table[m.IDTable]
	switch key.Name {
	case fyne.KeyReturn:
		m.onEnter()
	case "KP_Enter":
		m.onEnter()
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
