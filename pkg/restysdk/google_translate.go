package restysdk

import (
	"strings"
)

const (
	GoogleTranslateUrl   = "https://translate.googleapis.com/translate_a/single"
	GoogleParamSl        = "sl"
	GoogleParamTl        = "tl"
	GoogleParamQ         = "q"
	GoogleParamClient    = "client"
	GoogleParamClientVal = "gtx"
	GoogleParamDt        = "dt"
	GoogleParamDtVal     = "t"
)

type TranslateReq struct {
	Text string `json:"text"`
	Sl   string `json:"sl"`
	Tl   string `json:"tl"`
}

func (r *RestyClient) RestyTranslate(req *TranslateReq) (string, error) {
	rsp, err := r.R().SetQueryParams(map[string]string{
		GoogleParamSl:     req.Sl,
		GoogleParamTl:     req.Tl,
		GoogleParamQ:      req.Text,
		GoogleParamClient: GoogleParamClientVal,
		GoogleParamDt:     GoogleParamDtVal,
	}).Get(GoogleTranslateUrl)
	if err != nil {
		return "", err
	}

	// 返回的json反序列化比较麻烦, 直接字符串拆解
	ss := string(rsp.Body())
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)

	return ps[0], nil
}
