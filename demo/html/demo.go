package main

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/oldbai555/lb/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func main() {

}

func GetHtmlInfoByUrl(u string) error {
	u = Trim(u)
	p, err := url.Parse(u)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	host := net.ParseIP(p.Host)
	if host != nil {
		log.Errorf("err: %v ip is not allow", host)
		return fmt.Errorf("网页不允许请求IP")
	}
	httpRsp, err := http.Get(u)
	if err != nil {
		// if err.Error() == -2 {
		// 	return fmt.Errorf("网页请求失败")
		// }
		log.Errorf("err:%v", err)
		return fmt.Errorf("网页请求失败")
	}

	rspHeader := httpRsp.Header
	body, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		// if err.Error() == -2 {
		// 	return fmt.Errorf("网页请求失败")
		// }
		log.Errorf("err:%v", err)
		return fmt.Errorf("网页请求失败")
	}

	body, err = ConvertCharset(rspHeader.Get("Content-Type"), body)

	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		log.Errorf("err:%v", err)
		return fmt.Errorf("网页请求失败")
	}
	nodes, err := htmlquery.QueryAll(doc, "/html/head/title")
	if err == nil {
		if len(nodes) > 0 {
			title := strings.TrimSpace(htmlquery.InnerText(nodes[0]))
			fmt.Println(title)
		}
	} else {
		log.Errorf("err:%v", err)
	}
	picUrl := fmt.Sprintf("%v://%v/favicon.ico", p.Scheme, p.Hostname())
	fmt.Println(picUrl)
	nodes, err = htmlquery.QueryAll(doc, "//p")
	if err == nil {
		var desc string
		for _, node := range nodes {
			desc += strings.Replace(htmlquery.InnerText(node), "\n", "", -1)
			if len(desc) > 600 {
				break
			}
		}
		desc = ShortStr(desc, 512)
		fmt.Println(desc)
	} else {
		log.Errorf("err:%v", err)
	}
	return nil
}

func ShortStr(s string, max int) string {
	sr := []rune(s)
	if len(sr) > max {
		sr = sr[:max]
	}
	return string(sr)
}

func Trim(str string) string {
	if len(str) == 0 {
		return ""
	}
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\t")
	str = strings.Trim(str, " ")
	return str
}

func ConvertCharset(contentType string, byte []byte) ([]byte, error) {
	charset := GetCharset(contentType, string(byte))
	switch charset {
	case GB18030:
		return GB18030ToUtf8(byte)
	case GBK:
		return GbkToUtf8(byte)
	case UTF8:
		fallthrough
	default:
	}
	return byte, nil
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Errorf("err:%+v", e)
		return nil, e
	}
	return d, nil
}

func GB18030ToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Errorf("err:%+v", e)
		return nil, e
	}
	return d, nil
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
)

func GetCharset(contentType, content string) Charset {
	reg, _ := regexp.Compile("charset=[^\\w]?([-\\w]+)")
	if len(reg.FindStringSubmatch(contentType)) >= 2 {
		return Charset(strings.ToUpper(reg.FindStringSubmatch(contentType)[1]))
	}
	if len(reg.FindStringSubmatch(content)) >= 2 {
		return Charset(strings.ToUpper(reg.FindStringSubmatch(content)[1]))
	}
	return ""
}
