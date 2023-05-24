package countrycode

import (
	"fmt"
	"testing"
)

func TestGetCountryList(t *testing.T) {
	countries := GetCountryList()
	for _, country := range countries {
		fmt.Println(country)
	}
}

func TestGetPhoneCode(t *testing.T) {
	fmt.Println(GetPhoneCode("8407747877"))
}
