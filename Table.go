package main

import (
	"fmt"
	"image/color"
	"log"

	//	"sort"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TableOtoko struct {
	ColumnsName   [][]string
	ColumnsFormat []string
	ColumnsType   []string
	ColumnsWidth  []float32
	HeaderColor   color.Color
	AlterRowColor color.Color
	RowColor      color.Color
	Data          [][]string
	Edit          bool
	Tool         *widget.Toolbar
	Table         *widget.Table
	Header        *widget.Table
	we            map[*widget.Entry]widget.TableCellID
	wc            map[*widget.Check]widget.TableCellID
}

func (t *TableOtoko) CreateHeader() {
	t.Header = widget.NewTable(func() (int, int) {
		rows := len(t.ColumnsName)
		columns := len(t.ColumnsName[0])
		return rows, columns
	},
		func() fyne.CanvasObject {
			con := container.NewHBox()
			con.Layout = layout.NewMaxLayout()
			con.Add(widget.NewLabel(""))
			return container.New(layout.NewMaxLayout(),
				canvas.NewRectangle(color.Gray{250}),
				con,
			)
		},

		func(i widget.TableCellID, o fyne.CanvasObject) {

			box := o.(*fyne.Container)

			rect := box.Objects[0].(*canvas.Rectangle)
			rect.FillColor = t.HeaderColor

			cont := box.Objects[1].(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)

			label.SetText(t.ColumnsName[i.Row][i.Col])

		})

	for ic, v := range t.ColumnsWidth {
		t.Header.SetColumnWidth(ic, v)
	}

}

func (t *TableOtoko) LoarTable(mes []byte) {
}

func (t *TableOtoko) makeTable() {
	t.we = make(map[*widget.Entry]widget.TableCellID)
	t.wc = make(map[*widget.Check]widget.TableCellID)
	b := -1
	s := -1
	l := -1
	col := 0
	for _, v := range t.ColumnsType {
		switch v {
		case "bool":
			if b == -1 {
				b = col
				col = col + 1
			}

		case "string":
			if s == -1 && t.Edit {
				s = col
				col = col + 1
			}

		default:
			if l == -1 {
				l = col
				col = col + 1
			}

		}
	}

	t.Table = widget.NewTable(
		func() (int, int) {
			rows := len(t.Data)
			columns := len(t.Data[0])
			return rows, columns
		},
		func() fyne.CanvasObject {
			con := container.NewHBox()
			con.Layout = layout.NewMaxLayout()
			if l != -1 {
				con.Add(widget.NewLabel(""))
			}

			if b != -1 {
				check := widget.NewCheck("Optional", func(value bool) {
					log.Println("Check set to", value)
				})
				con.Add(check)
			}
			if s != -1 && t.Edit {
				entry := newEnterEntry()
				

				con.Add(entry)
			}

			return container.New(layout.NewMaxLayout(),
				canvas.NewRectangle(color.Gray{250}),
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
				rect.FillColor = t.AlterRowColor
			} else {
				rect.FillColor = t.RowColor
			}
			cont := box.Objects[1].(*fyne.Container)
			if l >= 0 {
				label = cont.Objects[l].(*widget.Label)
			}

			if b >= 0 {
				ic = cont.Objects[b].(*widget.Check)
			}
			if s >= 0 {
				entry = cont.Objects[s].(*enterEntry)
			}
			switch t.ColumnsType[i.Col] {
			case "bool":
				t.wc[ic] = i
				if l >= 0 {
					label.Hidden = true
				}
				if s >= 0 {
					entry.Hidden = true
				}
				if b >= 0 {
					ic.Hidden = false
				}
			case "string":
				entry.SetText(t.Data[i.Row][i.Col])
				//t.we[entry] = i
				if l >= 0 {
					label.Hidden = true
				}
				if s >= 0 {
					entry.Hidden = false
				}
				if b >= 0 {
					ic.Hidden = true
				}
			default:
				if l != -1 {
					label.SetText(t.Data[i.Row][i.Col])
				}
				if l >= 0 {
					label.Hidden = false
				}
				if s >= 0 {
					entry.Hidden = true
				}
				if b >= 0 {
					ic.Hidden = true
				}
			}
		})
	
	for ic, v := range t.ColumnsWidth {
		t.Table.SetColumnWidth(ic, v)
	}
	t.Table.OnSelected = func(id widget.TableCellID) {
		fmt.Printf("i.Col: %v\n", id.Col)

	}

t.Tool = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),)

	//t.CreateHeader()
t.Header = widget.NewTable(
		func() (int, int) {
			rows := len(t.ColumnsName)
			columns := len(t.ColumnsName[0])
			return rows, columns
		},
		func() fyne.CanvasObject {
			con := container.NewHBox()
			con.Layout = layout.NewMaxLayout()
			con.Add(widget.NewLabel(""))
			return container.New(layout.NewMaxLayout(),
				canvas.NewRectangle(color.Gray{250}),
				con,
			)
		},

		func(i widget.TableCellID, o fyne.CanvasObject) {
			var label *widget.Label
			box := o.(*fyne.Container)
			rect := box.Objects[0].(*canvas.Rectangle)
			rect.FillColor = t.HeaderColor

			cont := box.Objects[1].(*fyne.Container)
			label = cont.Objects[l].(*widget.Label)
			label.SetText(t.Data[i.Row][i.Col])
		})
	for ic, v := range t.ColumnsWidth {
		t.Header.SetColumnWidth(ic, v)
	}
	t.Header.OnSelected = func(id widget.TableCellID) {
		fmt.Printf("i.Col: %v\n", id.Col)

	}



}
//////////////////////////////////////////////
type enterEntry struct {
	widget.Entry
}

func (e *enterEntry) onEnter() {
	fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
}

func newEnterEntry() *enterEntry {
	entry := &enterEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *enterEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	default:
		e.Entry.KeyDown(key)
		fmt.Printf("Key %v pressed\n", key.Name)
	}
}

func (e *enterEntry) KeyUp(key *fyne.KeyEvent) {
	fmt.Printf("Key %v released\n", key.Name)
}



