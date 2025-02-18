/**
 * @Author: zjj
 * @Date: 2025/2/18
 * @Desc:
**/

package restysdk

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	utls "github.com/refraction-networking/utls"
	"net/http"
	"time"
)

func GetRandomJa3Tr() *http.Transport {
	spec, err := utls.UTLSIdToSpec(utls.HelloRandomized)
	if err != nil {
		return nil
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         spec.TLSVersMin,
			MaxVersion:         spec.TLSVersMax,
			CipherSuites:       spec.CipherSuites,
		},
		DisableKeepAlives: false,
	}
	return tr
}

type Ja3Option func(client *resty.Client)

func NewJa3Request(opts ...Ja3Option) (*resty.Request, error) {
	client := resty.New()
	tr := GetRandomJa3Tr()
	client.SetTransport(tr)
	client.SetTimeout(time.Duration(10) * time.Second)
	for _, opt := range opts {
		opt(client)
	}
	return client.NewRequest(), nil
}

func Ja3OptionWithProxy(proxy string) Ja3Option {
	return func(client *resty.Client) {
		client.SetProxy(proxy)
	}
}
