package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
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

type ColumnStyle struct {
	ID      string
	Name    string  //Заголовок столбца
	Format  string  //Форматированный вывод
	Width   float32 //Ширина столбца
	BGColor string  // Цвет фона
	Color   string  // Цвет текста
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
	Selected    widget.TableCellID
	wol         map[*oLabel]widget.TableCellID
	wb          map[*widget.Button]int
	wc          map[*widget.Check]widget.TableCellID
	wl          map[*widget.Label]widget.TableCellID
	we          map[*oEntry]widget.TableCellID
}

func getValidator(t string) fyne.StringValidator {
	switch t {
	case "String":
		return dateValidator
	case "email":
		return emailValidator
	case "Time":
		return floatValidator
	case "DOUBLE":
		return intValidator
	case "Перечисление":
		return intValidator
	case "bool":
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
