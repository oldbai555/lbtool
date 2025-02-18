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
	request, err := NewJa3Request()
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
	response, err := request.Get("https://surrit.com/f17fa4f3-4e70-428e-b7ad-441455a56027/playlist.m3u8")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(string(response.Body()))
	response, err = request.Get("https://surrit.com/f17fa4f3-4e70-428e-b7ad-441455a56027/720p/video.m3u8")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(string(response.Body()))
}
