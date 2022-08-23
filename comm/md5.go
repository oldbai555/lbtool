package comm

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func GetFileMd5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	m := md5.New()
	_, err = io.Copy(m, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", m.Sum(nil)), nil
}
func StrMd5(s string) string {
	m := md5.New()
	_, _ = m.Write([]byte(s))
	return fmt.Sprintf("%x", m.Sum(nil))
}
