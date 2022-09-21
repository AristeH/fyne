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
			//top := widget.NewButton("o", nil)
			toolbar := widget.NewToolbar(
				widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
					log.Println("New document")
				}),
				widget.NewToolbarAction(theme.ZoomInIcon(), func() {}),
			)
			middle := entry
			content := container.New(layout.NewBorderLayout(nil, nil, nil, toolbar),
				middle, toolbar)
			return container.New(layout.NewMaxLayout(),
				canvas.NewRectangle(color.Gray{Y: 250}),
				content,
			)

		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			box := o.(*fyne.Container)
			rect := box.Objects[0].(*canvas.Rectangle)
			entry1 := box.Objects[1].(*fyne.Container)
			toolbar := widget.NewToolbar(
				widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
					log.Println("New document")
				}),
				widget.NewToolbarAction(theme.ZoomInIcon(), func() {}),
			)
			entry := newOLabel()
			entry1.Objects[0] = entry
			entry.SetText(t.Data[i.Row][i.Col])
			entry.parent = t
			entry.Alignment = fyne.TextAlignLeading
			t.wol[entry] = i
			entry.TextStyle = fyne.TextStyle{
				Bold:      false,
				Italic:    false,
				Monospace: false,
				TabWidth:  0,
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
			content := container.New(layout.NewBorderLayout(nil, nil, entry, toolbar),
				entry, toolbar)
			box.Objects[1] = content
			if i.Col == 5 {
				toolbar := widget.NewToolbar(
					widget.NewToolbarAction(theme.ZoomInIcon(), func() {}),
				)
				content := container.New(layout.NewBorderLayout(nil, nil, nil, toolbar),
					toolbar)
				box.Objects[1] = content

			}

			if i == t.Selected {
				rect.FillColor = MapColor["tomato"]
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
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}))
}
