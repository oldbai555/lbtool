package main

import (
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lb/log"
	"net/url"
)
import "github.com/go-resty/resty/v2"

// https://github.com/go-resty/resty 官方文档

type genNua struct {
	Code uint32   `json:"code"`
	Data []string `json:"data"`
	Msg  string   `json:"msg"`
}

func main() {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().EnableTrace().Get("https://httpbin.org/get")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
	fmt.Println("  Trace Info    :", resp.Request.TraceInfo())

	val := url.Values{}
	val.Set("count", fmt.Sprintf("%d", 1))
	val.Set("type", "windows")
	resp, err = client.R().SetFormDataFromValues(val).Post("https://www.bejson.com/Bejson/Api/Common/ge_nua")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	var gen genNua
	err = json.Unmarshal(resp.Body(), &gen)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	fmt.Printf("  Body :%+v\n", gen)
}
