package excel

import (
	"fmt"
	"github.com/oldbai555/lbtool/extrpkg/pie/pie"
	"net/http"

	"github.com/xuri/excelize/v2"
)

// ImportExcel 导入 Excel
func ImportExcel(req *ImportExcelReq, logic ImportExcelLogicFunc) error {

	response, err := http.Get(req.Url)
	if err != nil {
		fmt.Println(err)
	}

	f, err := excelize.OpenReader(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}

		if err = response.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var sheet = DefaultSheet
	if req.Sheet != "" {
		sheet = req.Sheet
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println(err)
	}

	var records []pie.Strings
	for _, row := range rows {
		records = append(records, row)
	}

	return logic(records)
}
