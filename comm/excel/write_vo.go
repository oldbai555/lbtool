package excel

import "github.com/xuri/excelize/v2"

type ExportExcelLogicFunc func(f *excelize.File) error

type ExportExcelReq struct {
	// FileName 导出的文件名
	FileName string `json:"file_name"`
}
