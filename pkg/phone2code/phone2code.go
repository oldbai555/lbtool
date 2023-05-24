package phone2code

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
)

type CountryInfo struct {
	EnName      string `json:"en_name"`
	CnName      string `json:"cn_name"`
	CountryCode string `json:"country_code"`
	PhoneCode   string `json:"phone_code"`
}

//go:embed phone_rule.json
var countryInfoLists string

var (
	MaxPhoneCodeLength = 4
	countryInfoList    []*CountryInfo
)

func init() {
	// 从配置中心读取
	// apolloNamSpace := "phone_code"
	// // 从 apollo 配置中心取数据
	// apolloCli, err := apollo.New(
	// 	apollosdk.WithDefaultNamespace(apolloNamSpace),
	// 	apollosdk.WithNamespaces(apolloNamSpace),
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// err = apolloCli.GetJsonData("country_list", &countryInfoList)
	// if err != nil {
	// 	panic(err)
	// }
	// time.Sleep(10 * time.Second)
	// 从本地文件读取数据
	// filePath := utils.CurrentFile()
	// dirPath := path.Dir(filePath)
	// fp, err := os.Open(path.Join(dirPath, "phone_rule.json"))
	// if err != nil {
	// 	panic(err)
	// }
	// defer fp.Close()
	// bytes, err := ioutil.ReadAll(fp)
	// if err != nil {
	// 	panic(err)
	// }
	err := json.Unmarshal([]byte(countryInfoLists), &countryInfoList)
	if err != nil {
		panic(err)
	}
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
func getCountryInfo(text string) *CountryInfo {
	for _, countryInfo := range countryInfoList {
		if text == countryInfo.PhoneCode {
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
			phoneCode = countryInfo.PhoneCode
			break
		}
		text = text[:i]
	}
	return
}
