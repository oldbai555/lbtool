package excel

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/xuri/excelize/v2"
	"testing"
)

func TestExportExcel(t *testing.T) {
	url, err := ExportExcel(&ExportExcelReq{
		FileName: "列表.xlsx",
	}, func(f *excelize.File) error {

		// Create a new sheet.
		// index := f.NewSheet("Sheet2")
		// f.SetActiveSheet(index)

		f.SetSheetName("Sheet1", "列表")
		excelWriter, err := f.NewStreamWriter("列表")
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 第一行写表头
		err = excelWriter.SetRow("A1", []interface{}{"1", "2"})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 第二行开始写数据
		var rowNum = 2
		for i := 0; i < 5; i++ {
			var record []interface{}
			record = append(record, "123", "456")

			// 写入表格
			cell, err := excelize.CoordinatesToCellName(1, rowNum)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			if err = excelWriter.SetRow(cell, record); err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			rowNum++
		}

		// 刷新一下缓存区
		if err = excelWriter.Flush(); err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("err is %v", err)
	}
	log.Errorf(url)
}
