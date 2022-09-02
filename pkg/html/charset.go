package html

import (
	"bytes"
	"github.com/oldbai555/lb/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"regexp"
	"strings"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
)

func convertCharset(contentType string, byte []byte) ([]byte, error) {
	charset := getCharset(contentType, string(byte))
	switch charset {
	case GB18030:
		return gB18030ToUtf8(byte)
	case GBK:
		return gbkToUtf8(byte)
	case UTF8:
		fallthrough
	default:
	}
	return byte, nil
}

func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Errorf("err:%+v", e)
		return nil, e
	}
	return d, nil
}

func gB18030ToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Errorf("err:%+v", e)
		return nil, e
	}
	return d, nil
}

func getCharset(contentType, content string) Charset {
	reg, _ := regexp.Compile("charset=[^\\w]?([-\\w]+)")
	if len(reg.FindStringSubmatch(contentType)) >= 2 {
		return Charset(strings.ToUpper(reg.FindStringSubmatch(contentType)[1]))
	}
	if len(reg.FindStringSubmatch(content)) >= 2 {
		return Charset(strings.ToUpper(reg.FindStringSubmatch(content)[1]))
	}
	return ""
}
