package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	RandomStringModNumberPlusLetter           = 1
	RandomStringModNumberPlusLetterPlusSymbol = 2
	RandomStringModNumber                     = 3
)

func GenRandomStr() string {
	rndStr := fmt.Sprint(
		os.Getpid(), time.Now().UnixNano(), rand.Float64())
	h := md5.New()
	_, _ = io.WriteString(h, rndStr)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetRandomString(length int64, mod uint32) string {
	var strKey string
	if mod == RandomStringModNumberPlusLetter {
		strKey = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	} else if mod == RandomStringModNumberPlusLetterPlusSymbol {
		strKey = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*(){}|\\[]?/"
	} else if mod == RandomStringModNumber {
		strKey = "0123456789"
	} else {
		strKey = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = strKey[r.Intn(len(strKey))]
	}

	return string(bytes)
}

func GetRandomStrNoExisted(length int64, existedMap map[string]bool) string {
	var str string
	for {
		str = GetRandomString(length, RandomStringModNumberPlusLetter)
		if !existedMap[str] {
			break
		}
	}
	return str
}
