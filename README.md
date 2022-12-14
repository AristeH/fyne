# fyne

Program structure

FormData - list of application forms
 
PutListForm(name, header string) name - ID form

NewOTable(name, data)

data - a data structure containing two fields
 
    Data [][]string - the first row contains the table header
	DataDesciption [][]string - describes the columns of the table
  
Type columns - id, string, float, enum, date

You can set the background and font color for each table cell.
 Row color, column color               
