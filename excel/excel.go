package excel

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

type Excel struct {
	path string
	file *xlsx.File
}

func New(path string) *Excel {
	return &Excel{
		file: xlsx.NewFile(),
		path: path,
	}
}

func (e *Excel) NewSheet(name string, heads []string) (*xlsx.Sheet, error) {
	sheet, err := e.file.AddSheet(name)
	if err != nil {
		return nil, err
	}
	r := sheet.AddRow()
	for _, s := range heads {
		r.AddCell().SetString(s)
	}
	return sheet, nil
}

func (e *Excel) Save() error {
	return e.file.Save(e.path)
}

func (e *Excel) AddRow(sheet *xlsx.Sheet, data []interface{}) {
	row := sheet.AddRow()
	for _, v := range data {
		row.AddCell().SetString(fmt.Sprintf("%v", v))
	}
}

func (e *Excel) Path() string {
	return e.path
}
