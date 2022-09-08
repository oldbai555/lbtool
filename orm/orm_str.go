package orm

import (
	"fmt"
	"strings"
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
	return strings.Contains(strings.ToUpper(typ), "INT")
}

func escapeMysqlLikeWildcard(val string) string {
	l := len(val)
	dest := make([]byte, 0, 2*len(val))
	var escape byte
	for i := 0; i < l; i++ {
		c := val[i]
		escape = 0
		switch c {
		case '%':
			escape = '%'
		case '_':
			escape = '_'
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
		case '\n': /* Must be escaped for logs */
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': /* Better safe than sorry */
			escape = '"'
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return string(dest)
}

func escapeMysqlLikeWildcardIgnore2End(val string) string {
	l := len(val)
	dest := make([]byte, 0, 2*len(val))
	var escape byte
	for i := 0; i < l; i++ {
		c := val[i]
		if (i == 0 || i == l-1) && c == '%' {
			dest = append(dest, c)
			continue
		}
		escape = 0
		switch c {
		case '%':
			escape = '%'
		case '_':
			escape = '_'
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
		case '\n': /* Must be escaped for logs */
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': /* Better safe than sorry */
			escape = '"'
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return string(dest)
}

func escapeMysqlString(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		escape = 0
		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
		case '\n': /* Must be escaped for logs */
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': /* Better safe than sorry */
			escape = '"'
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return string(dest)
}

func addlrQuotes(val string) string {
	return fmt.Sprintf(" '%s'", escapeMysqlLikeWildcardIgnore2End(val))
}
