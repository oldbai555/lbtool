package orm

import (
	"fmt"
)

func quoteName(name string) string {
	if name != "" {
		if name[0] != '`' {
			q := true
			l := len(name)
			for i := 0; i < l; i++ {
				c := name[i]
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
					q = false
					break
				}
			}
			if q {
				name = fmt.Sprintf("`%s`", name)
			}
		}
	}
	return name
}

func inSliceStr(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}

func isIntType(typ string) bool {
	switch typ {
	case "uint32", "int32", "uint64", "int64":
		return true
	}
	return false
}
