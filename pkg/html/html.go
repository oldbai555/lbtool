package html

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Result struct {
	PicUrl string `json:"pic_url"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
}

func GetHtmlInfoByUrl(u string) error {
	u = trim(u)
	p, err := url.Parse(u)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	host := net.ParseIP(p.Host)
	if host != nil {
		log.Errorf("err: %v ip is not allow", host)
		return lberr.NewCustomErr("网页不允许请求IP")
	}
	httpRsp, err := http.Get(u)
	if err != nil {
		// if err.Error() == -2 {
		// 	return fmt.Errorf("网页请求失败")
		// }
		log.Errorf("err:%v", err)
		return lberr.NewCustomErr("网页请求失败")
	}

	rspHeader := httpRsp.Header
	body, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		// if err.Error() == -2 {
		// 	return fmt.Errorf("网页请求失败")
		// }
		log.Errorf("err:%v", err)
		return lberr.NewCustomErr("网页请求失败")
	}

	body, err = convertCharset(rspHeader.Get("Content-Type"), body)

	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		log.Errorf("err:%v", err)
		return lberr.NewCustomErr("网页请求失败")
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
		desc = shortStr(desc, 512)
		fmt.Println(desc)
	} else {
		log.Errorf("err:%v", err)
	}
	return nil
}

func GetHtmlResultByUrl(u string) (*Result, error) {
	var res Result
	u = trim(u)
	p, err := url.Parse(u)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	host := net.ParseIP(p.Host)
	if host != nil {
		log.Errorf("err: %v ip is not allow", host)
		return nil, lberr.NewCustomErr("网页不允许请求IP")
	}
	httpRsp, err := http.Get(u)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, lberr.NewCustomErr("网页请求失败")
	}

	rspHeader := httpRsp.Header
	body, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, lberr.NewCustomErr("网页请求失败")
	}

	body, err = convertCharset(rspHeader.Get("Content-Type"), body)

	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, lberr.NewCustomErr("网页请求失败")
	}
	nodes, err := htmlquery.QueryAll(doc, "/html/head/title")
	if err == nil {
		if len(nodes) > 0 {
			res.PicUrl = strings.TrimSpace(htmlquery.InnerText(nodes[0]))
		}
	} else {
		log.Errorf("err:%v", err)
	}
	res.PicUrl = fmt.Sprintf("%v://%v/favicon.ico", p.Scheme, p.Hostname())
	nodes, err = htmlquery.QueryAll(doc, "//p")
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	var desc string
	for _, node := range nodes {
		desc += strings.Replace(htmlquery.InnerText(node), "\n", "", -1)
		if len(desc) > 600 {
			break
		}
	}
	desc = shortStr(desc, 512)
	res.Desc = desc
	log.Infof("res: %v", res)
	return &res, nil
}
