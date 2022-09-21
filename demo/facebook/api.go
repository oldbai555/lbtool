package main

import (
	"fmt"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/pkg/resty_utils"
	"net/http"
)

func genToken() (string, error) {
	client := resty_utils.NewRestyClient()
	var prefix = "https://graph.facebook.com/v15.0/"
	response, err := client.R().
		SetQueryParams(map[string]string{
			"access_token": "EAAHK1WzvV4gBAOCnTYF0OlZBcw99TBHm41xZBXe8a5NhmiqcKCDKJK86cYINTD8IBBgeSUPtLdSjWZCcNE76WL6O6wt0ZBifnGJwmQT2VWo8czv2mqDaZAsGMED2wSfcD4FJKgQuQ1rzZCDxB6eCuEUT3EYZBeGHR3XmPfGPqZAQJhYZC1fKApAiBapJqQfHPMk90bdeKDl0KEdIVnvlu4pVYS73FOqilJvJ9ozi3hYnHamZAGhz2lx6cvgHlZAHGyUw4YZD",
		}).
		Get(fmt.Sprintf("%s/%s", prefix, "me"))
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	log.Infof("response is %v", string(response.Body()))
	return "", nil
}

// curl -i -X GET \
// "https://graph.facebook.com/v15.0/me?access_token=EAAHK1WzvV4gBAKmhjmMzZCLozoBV2mV4iA3fGVZCXnhLLEULhnZC1vW10aFtklLQBtnwDdADQl2BxCa9K09KMIZAZCUpIJtAU8f5KaIYb0ZA788Pd9A7gZBfrn30OimZBivipNQPBjyaQhHyN2lbZATFLdgyasOpjlspx17j1bcUJQSA4l6PtijaAPQbOk80RMSD2Oq2lqKWLWgZDZD"
func httpget() {
	rsp, err := http.Get("https://graph.facebook.com/v15.0/me?access_token=EAAHK1WzvV4gBAKmhjmMzZCLozoBV2mV4iA3fGVZCXnhLLEULhnZC1vW10aFtklLQBtnwDdADQl2BxCa9K09KMIZAZCUpIJtAU8f5KaIYb0ZA788Pd9A7gZBfrn30OimZBivipNQPBjyaQhHyN2lbZATFLdgyasOpjlspx17j1bcUJQSA4l6PtijaAPQbOk80RMSD2Oq2lqKWLWgZDZD")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("rsp is %v", rsp)
}

func main() {
	// httpget()
	token, err := genToken()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("token : %s", token)
	// updatedAtMin := time.Unix(int64(1663664723), 0)
	// fmt.Println(time.Unix(1663733835, 0).Format("2006-01-02T15:04:05-00:00"))
}
