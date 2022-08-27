package excel

import "github.com/oldbai555/lb/extrpkg/pie/pie"

const DefaultSheet = "Sheet1"

type ImportExcelLogicFunc func(records []pie.Strings) error

type ImportExcelReq struct {
	Url   string `json:"url"`
	Sheet string `json:"sheet"`
}
