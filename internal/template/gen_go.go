package template

import (
	"fmt"
	"github.com/oldbai555/lb/log"
	"os"
	"text/template"
)

// GenTemplate 生成模板文件
func GenTemplate(file *os.File, f *Function) (err error) {
	tmpl, _ := template.New(f.ModelName).Parse(f.Template)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	err = tmpl.Execute(file, f)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	return
}

func GetOsFile(filePath string, fileName string) (file *os.File, err error) {
	if filePath != "" {
		file, err = os.Create(fmt.Sprintf("%s\\%s.go", filePath, fileName))
	} else {
		file, err = os.Create(fmt.Sprintf("%s.go", fileName))
	}
	return
}
