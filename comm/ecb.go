package comm

import (
	"encoding/hex"
	"fmt"
	"github.com/forgoer/openssl"
)

// ECBEncrypt 编码
func ECBEncrypt(src string, key string) (string, error) {
	data, err := openssl.AesECBEncrypt([]byte(src), []byte(key), openssl.PKCS7_PADDING)
	return hex.EncodeToString(data), err
}

// ECBDecrypt 解码
func ECBDecrypt(src string, key string) (rsp string, err error) {
	defer func() {
		if errMsg := recover(); errMsg != any(nil) {
			fmt.Println("捕获异常:", err)
		}
	}()

	data, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	databyte, err := openssl.AesECBDecrypt(data, []byte(key), openssl.PKCS7_PADDING)
	return string(databyte), err
}
