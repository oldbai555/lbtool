package excel

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
)

// ExportExcel 导出 Excel
func ExportExcel(req *ExportExcelReq, logic ExportExcelLogicFunc) (string, error) {
	f := excelize.NewFile()
	err := logic(f)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	b, err := f.WriteToBuffer()
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	if err = ioutil.WriteFile(req.FileName, []byte(b.String()), 0755); err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	// TODO 上传文件
	return req.FileName, nil
}
