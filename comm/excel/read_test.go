package excel

import (
	"github.com/oldbai555/lb/extrpkg/pie/pie"
	"github.com/oldbai555/lb/log"
	"testing"
)

func TestImportExcel(t *testing.T) {
	err := ImportExcel(&ImportExcelReq{
		Url:   "https://quan-peak-1259287960.cos.ap-guangzhou.myqcloud.com/21/af8e01b03ab14fddae1c59259a09501c.xlsx",
		Sheet: "成员列表",
	}, func(records []pie.Strings) error {
		for _, record := range records {
			log.Infof("record: %v", record)
		}
		return nil
	})
	if err != nil {
		log.Errorf("err is %v", err)
	}
}
