package countrycode

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
)

type Country struct {
	Cn        string `json:"cn"`
	En        string `json:"en"`
	PhoneCode string `json:"phone_code"`
}

//go:embed countrycode.json
var countryCodeStr string

var countryList []*Country
var MaxPhoneCodeLength = 4

func init() {
	err := json.Unmarshal([]byte(countryCodeStr), &countryList)
	if err != nil {
		panic(err)
	}
}

func GetCountryList() []*Country {
	return countryList
}

// clearText 清洗输入的手机号
// 替换掉空白符号 , 非数字, 前导 00
// 区号最长为 4位，直接返回最长的位数
func clearText(text string) string {
	// 替换空白度
	regBlank := regexp.MustCompile(`\s`)
	text = string(regBlank.ReplaceAll([]byte(text), []byte("")))
	// 替换 非数字
	regNotNumbered := regexp.MustCompile(`[^0-9]`)
	text = string(regNotNumbered.ReplaceAll([]byte(text), []byte("")))
	// 替换掉前导 00
	if strings.HasPrefix(text, "00") {
		text = strings.Replace(text, "00", "", 1)
	}
	if len(text) > MaxPhoneCodeLength {
		text = text[0:MaxPhoneCodeLength]
	}
	return text
}

// getCountryInfo 通过区号获取国家信息
func getCountryInfo(text string) *Country {
	for _, countryInfo := range countryList {
		if strings.TrimPrefix(text, "+") == strings.TrimPrefix(countryInfo.PhoneCode, "+") {
			return countryInfo
		}
	}
	return nil
}

// GetPhoneCode 通过手机号获取 区号
func GetPhoneCode(phone string) (phoneCode string) {
	text := clearText(phone)
	for i := len(text) - 1; i >= 0; i-- {
		countryInfo := getCountryInfo(text)
		if countryInfo != nil {
			phoneCode = strings.TrimPrefix(countryInfo.PhoneCode, "+")
			break
		}
		text = text[:i]
	}
	return
}
