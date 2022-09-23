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
)

func (t *TableOtoko) MakeTableLabel() {
	t.Table = widget.NewTable(
		func() (int, int) {
			rows := len(t.Data)
			columns := len(t.Data[0])
			return rows, columns
		},
		func() fyne.CanvasObject {
			entry := newOLabel()
			entry.IDForm = t.IDForm
			entry.IDTable = t.ID
			entry.parent = t
			toolbar := widget.NewToolbar()
			content := container.New(layout.NewBorderLayout(nil, nil, entry, toolbar),
				entry, toolbar)
			return container.New(layout.NewMaxLayout(),
				canvas.NewRectangle(color.Gray{Y: 250}),
				content,
			)

		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			toolbar := widget.NewToolbar()
			open := widget.NewToolbarAction(theme.DocumentIcon(), func() {
				fmt.Println("profile settings clicked in toolbar")
			})
			box := o.(*fyne.Container)
			rect := box.Objects[0].(*canvas.Rectangle)
			entry := newOLabel()
			entry.Ind = &i
			entry.parent = t
			entry.SetText(t.Data[i.Row][i.Col])
			if t.ColumnStyle[i.Col].Width == 0 {
				entry.Hidden = true
			} else {
				entry.Hidden = false
			}
			switch t.ColumnStyle[i.Col].Type {

			case "float":
				entry.Label.Alignment = fyne.TextAlignTrailing
			case "String":
				entry.Label.Alignment = fyne.TextAlignLeading
			default:
				toolbar.Append(open)
				entry.Label.Alignment = fyne.TextAlignLeading
			}
			entry.TextStyle = fyne.TextStyle{
				Bold: false,
			}
			if i.Row == 0 {
				rect.FillColor = MapColor[t.TabStyle.HeaderColor]
				entry.Alignment = fyne.TextAlignCenter
				entry.TextStyle = fyne.TextStyle{
					Bold: true,
				}
			} else if i.Row%2 == 0 {
				rect.FillColor = MapColor[t.TabStyle.RowAlterColor]

			} else {
				rect.FillColor = MapColor[t.TabStyle.RowColor]
			}
			if val, ok := MapColor[t.ColumnStyle[i.Col].BGColor]; ok {
				rect.FillColor = mix(val, rect.FillColor)
			}
			content := container.New(layout.NewBorderLayout(nil, nil, entry, nil), entry)
			box.Objects[1] = content
			if i == t.Selected {
				rect.FillColor = MapColor["Selected"]
				input := newoEntry()
				input.Ind = &i
				input.SetText(t.Data[i.Row][i.Col])
				content := container.New(layout.NewBorderLayout(nil, nil, nil, toolbar),
					toolbar, input)
				box.Objects[1] = content
			}
		})
	for ic, v := range t.ColumnStyle {
		t.Table.SetColumnWidth(ic, v.Width)
	}
	t.Table.OnSelected = func(id widget.TableCellID) {
		t.Selected = id
		activeContainer = t
		fmt.Printf("i.Col: %v\n", id.Col)
		t.Table.Refresh()

	}

	t.Tool = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}))
}
