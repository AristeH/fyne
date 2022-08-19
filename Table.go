package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"sort"
)

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

func (t *TableOtoko) MakeTableEntry() {

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
			entry := newoEntry()

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
			var entry *oEntry

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
			if i == t.Selected {
				rect.FillColor = MapColor["tomato"]
			}

			cont := box.Objects[1].(*fyne.Container)

			label = cont.Objects[0].(*widget.Label)

			ic = cont.Objects[1].(*widget.Check)

			entry = cont.Objects[2].(*oEntry)
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
		t.Selected = id
		activeContainer = t
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
	TO.we = make(map[*oEntry]widget.TableCellID)
	TO.wl = make(map[*widget.Label]widget.TableCellID)
	TO.wol = make(map[*oLabel]widget.TableCellID)

	return &TO
}
