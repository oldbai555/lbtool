package restysdk

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/oldbai555/lbtool/log"
	"net/url"
)

const (
	HeaderUserAgent = "user-agent"
)

type UserAgentResult struct {
	Code uint32   `json:"code"`
	Data []string `json:"data"`
	Msg  string   `json:"msg"`
}

func (r *RestyClient) GetRandomUserAgent() (string, error) {
	val := url.Values{}
	val.Set("count", fmt.Sprintf("%d", 1))
	val.Set("type", "windows")
	resp, err := resty.New().R().SetFormDataFromValues(val).Post("https://www.bejson.com/Bejson/Api/Common/ge_nua")
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	var result UserAgentResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return result.Data[0], nil
}
