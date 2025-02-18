/**
 * @Author: zjj
 * @Date: 2025/2/18
 * @Desc:
**/

package restysdk

import (
	"testing"
)

func TestGetRandomJa3(t *testing.T) {
	request, err := NewJa3Request(Ja3OptionWithProxy("http://127.0.0.1:7890"))
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	agent, err := GetRandomUserAgent()
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	request.SetHeader("User-Agent", agent)
	response, err := request.Get("")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(string(response.Body()))
}
