package comm

import (
	"crypto/sha256"
	"fmt"
)

func StrSha256(dst string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(dst))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func StrSha256WithSalt(dst string, salt string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(dst + salt))
	return fmt.Sprintf("%x", h.Sum(nil))
}
