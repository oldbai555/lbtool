package excel

import "github.com/oldbai555/lbtool/extpkg/pie/pie"

const DefaultSheet = "Sheet1"

type ImportExcelLogicFunc func(records []pie.Strings) error

type ImportExcelReq struct {
	Url string `json:"url"`
	// Sheet 要操作的 sheet
	Sheet string `json:"sheet"`
}
