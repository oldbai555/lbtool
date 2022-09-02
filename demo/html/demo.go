package main

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/pkg/html"
)

func main() {
	err := html.GetHtmlInfoByUrl("https://www.msn.cn/zh-cn/health/other/%E7%97%85%E4%BA%BA%E5%B7%B2%E9%BA%BB%E9%86%89-%E8%BF%98%E5%BE%97%E7%AD%89%E5%9B%BD%E5%A4%96%E8%A7%A3%E9%94%81-%E8%BF%99%E4%BD%8D%E5%8C%BB%E7%94%9F%E7%9C%BC%E7%9D%9B%E6%B9%BF%E6%B6%A6-%E6%88%91%E4%BB%AC%E6%9C%89%E8%87%AA%E5%B7%B1%E7%9A%84%E5%9B%BD%E4%BA%A7%E6%9C%BA%E5%99%A8%E4%BA%BA%E4%BA%86/ar-AA11lIjL?ocid=msedgntp&cvid=9cc4ec7052344e2c9670652ceaa35958")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
