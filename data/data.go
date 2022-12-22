package data

import (
	"math/rand"
	"otable/pkg/logger"
	"strconv"

	"github.com/sirupsen/logrus"
)

type GetData struct {
	Data            [][]string
	DataDescription [][]string
	Enum            map[string][]string
}

// TestData - тестовые данные
func TestData() *GetData {
	l := logger.GetLog()
	colColumns := 10

	colRows := 500

	// инициализация данных таблицы
	data := make([][]string, colRows)
	for i := 0; i < colRows; i++ {
		data[i] = make([]string, colColumns)
	}
	// Header
	// data[0][0] = "\u2713 №"
	data[0][0] = "id"
	data[0][1] = "id_product"
	data[0][2] = "Product"
	data[0][3] = "Type"
	data[0][4] = "Quantity"
	data[0][5] = "Price"
	data[0][6] = "Amount"
	data[0][7] = "Check"
	data[0][8] = "DateCheck"
	data[0][9] = "Comment"
	//Dstart",

	l.WithFields(logrus.Fields{"columns": len(data[0]), "rows": len(data), "event": "Fill Data"}).Info("Data")

	// инициализация описания данных таблицы
	datadescription := make([][]string, 4)
	for i := 0; i < 4; i++ {
		datadescription[i] = make([]string, colColumns)
	}
	// Name column
	datadescription[0][0] = "id"
	datadescription[0][1] = "id_product"
	datadescription[0][2] = "product"
	datadescription[0][3] = "Type"
	datadescription[0][4] = "Quantity"
	datadescription[0][5] = "Price"
	datadescription[0][6] = "Amount"
	datadescription[0][7] = "Check"
	datadescription[0][8] = "DateCheck"
	datadescription[0][9] = "Comment"

	//  Type column
	datadescription[1][0] = "id"
	datadescription[1][1] = "id"
	datadescription[1][2] = "id_string"
	datadescription[1][3] = "enum"
	datadescription[1][4] = "int"
	datadescription[1][5] = "float,15,2"
	datadescription[1][6] = "float,15,2"
	datadescription[1][7] = "bool"
	datadescription[1][8] = "date"
	datadescription[1][9] = "string"

	// Width column
	datadescription[2][0] = "0"
	datadescription[2][1] = "0"
	datadescription[2][2] = "20"
	datadescription[2][3] = "10"
	datadescription[2][4] = "10"
	datadescription[2][5] = "18"
	datadescription[2][6] = "18"
	datadescription[2][7] = "3"
	datadescription[2][8] = "10"
	datadescription[2][9] = "25"

	//Formula column
	datadescription[3][0] = ""
	datadescription[3][1] = "Amount =Quantity*Price"
	datadescription[3][2] = "Amount =Quantity*Price"
	datadescription[3][3] = "Amount =Quantity*Price"
	datadescription[3][4] = "Amount =Quantity*Price"
	datadescription[3][5] = ""
	datadescription[3][6] = ""
	datadescription[3][7] = ""
	datadescription[3][8] = ""
	datadescription[3][9] = ""

	a := 10
	b := 100
	for i := 1; i < colRows; i++ {
		data[i][0] = "id" + strconv.Itoa(i)
		data[i][1] = "id_product" + strconv.Itoa(i)
		data[i][2] = "product " + strconv.Itoa(i)
		data[i][3] = "product"
		data[i][4] = strconv.Itoa(a + rand.Intn(b-a+1))
		data[i][5] = strconv.FormatFloat(rand.Float64()*200, 'f', 2, 64)
		data[i][6] = "0"
		data[i][7] = "0"
		data[i][8] = ""
		data[i][9] = "Comment"
	}

	l.WithFields(logrus.Fields{
		"form":         "Test data",
		"columns ddes": len(datadescription[0]),
		"rows data":    len(data),
		"event":        "Init data",
	}).Info("Data")
	e := map[string][]string{
		"Type": {"product", "service"},
	}
	return &GetData{data, datadescription, e}
}
