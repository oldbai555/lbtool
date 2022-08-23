package comm

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GenUUID() string {
	//Version 1,基于 timestamp 和 MAC address (RFC 4122)
	//Version 2,基于 timestamp, MAC address 和 POSIX UID/GID (DCE 1.1)
	//Version 3, 基于 MD5 hashing (RFC 4122)
	//Version 4, 基于 random numbers (RFC 4122)
	//Version 5, 基于 SHA-1 hashing (RFC 4122)
	u2 := uuid.NewV4()
	str := strings.ReplaceAll(u2.String(), "-", "")
	return str[:16]
}
