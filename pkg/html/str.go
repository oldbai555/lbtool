package html

import "strings"

func trim(str string) string {
	if len(str) == 0 {
		return ""
	}
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\t")
	str = strings.Trim(str, " ")
	return str
}

func shortStr(s string, max int) string {
	sr := []rune(s)
	if len(sr) > max {
		sr = sr[:max]
	}
	return string(sr)
}
