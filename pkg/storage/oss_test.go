package storage

import (
	"bytes"
	"fmt"
	"github.com/oldbai555/lb/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestOSSStorage_SignURL(t *testing.T) {
	Setup(Config{
		Type:      "qcloud",
		SecretID:  "",
		SecretKey: "",
		BucketURL: "https://baifile-1309918034.cos.ap-guangzhou.myqcloud.com",
	})
	readFile, err := ioutil.ReadFile("C:\\Users\\EDY\\Desktop\\QQ截图20220805110655.png")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	objectKey := `public/link-info/assets/images/` + fmt.Sprintf("%d.%s", rand.Int63(), "png")

	err = FileStorage.Put(objectKey, bytes.NewReader(readFile))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	signedURL, err := FileStorage.SignURL(objectKey, http.MethodGet, 3600*24*365)
	if err != nil {
		t.Failed()
	}
	println(signedURL)

}
