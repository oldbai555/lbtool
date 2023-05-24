package excel

import (
	"fmt"
	"github.com/oldbai555/lbtool/extpkg/pie/pie"
	"strings"
)

const (
	FieldName      = "name"
	FieldErrReason = "错误原因"
	MaxTitleNumber = 4 // 表头字段数量
	MaxRow         = 3
)

// Header 定位标题表头的结构体
type Header struct {
	TitleRowIdx int // 表头所在的行
	NameIdx     int
}

// RowCheckResult Excel 行检查结果
type RowCheckResult struct {
	ErrorMsgList []string
	Name         string
	IsEmpty      bool
}

// ParseFileExt 截取文件后缀
func ParseFileExt(fileUrl string) (string, error) {
	var fileExt string
	if strings.HasSuffix(fileUrl, ".xlsx") {
		fileExt = "xlsx"
	} else if strings.HasSuffix(fileUrl, ".xls") {
		fileExt = "xls"
	} else {
		return "", fmt.Errorf("invalid file url missed file ext")
	}
	return fileExt, nil
}

// GetHeader 解析文件表头
func GetHeader(record [][]string) (*Header, error) {
	header := &Header{
		NameIdx: -1,
	}

	// 查找 excel 表头字段的位置下标
	var getExcelTitle = func(record pie.Strings) int {
		nameIdx := -1
		// 校验表头的字段数量
		if len(record) != MaxTitleNumber {
			return nameIdx
		}
		for idx, value := range record {
			value = strings.Trim(value, " ")
			value = strings.ReplaceAll(value, "*", "")
			if value == FieldName {
				nameIdx = idx
				continue
			}
		}
		return nameIdx
	}

	// 在此处定义读取文件的对象
	// 一行行读数据，MaxRow 行内如果找不到标题就格式错误处理
	var maxRow = MaxRow
	if len(record) <= 0 {
		return nil, fmt.Errorf("invalid file format fail")
	} else if len(record) < 3 {
		maxRow = len(record)
	}
	for header.TitleRowIdx = 0; header.TitleRowIdx < maxRow; header.TitleRowIdx++ {
		excelRecordLine := record[header.TitleRowIdx]

		header.NameIdx = getExcelTitle(excelRecordLine)
		if header.NameIdx >= 0 {
			break
		}
	}

	// 找不到对应标题的列
	if header.NameIdx < 0 {
		return header, fmt.Errorf("invalid file format fail")
	}

	return header, nil
}
