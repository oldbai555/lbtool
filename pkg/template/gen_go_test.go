package template

import (
	"log"
	"testing"
)

func Test_GenTemplate(t *testing.T) {
	file, err := GetOsFile("C:\\Users\\baigege\\Desktop\\lb\\internal\\template", "hello")
	if err != nil {
		log.Printf("err is %v", err)
		return
	}
	err = GenTemplate(file, &Function{
		Template:    teplDemo1,
		Package:     "template",
		ModelName:   "lbx",
		Description: "this is a template",
	})
	if err != nil {
		log.Printf("err is %v", err)
		return
	}
}
