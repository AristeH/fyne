# Table

Program structure

FormData - list of application forms
 
PutListForm(name, header string) name - ID form, registration of the form

NewOTable(name, data) - creating a table

    data - a data structure containing two fields
    type GetData struct {
	Data           [][]string - the first row contains the table header
	DataDesciption [][]string- describes the columns of the table
    }

DataDesciption [][]string 
      1 row of the table contains the column ID
      2 the table row contains the column type
      3 table row column width
      4 row of the table contains the formula
      
Type columns:
	id - ID record
	id_string - Selection from another table
	enum - a list of strings to select a value
	int
	float
	string
	bool
	date
	


You can set the background and font color for each table cell.
Row color, column color               
