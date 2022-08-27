package log

import (
	"fmt"
)

type formatter interface {
	Sprintf(level Level, color Color, buf string) (string, error)
}

func transferLevelToStr(level Level) (string, error) {
	if str, ok := levelToStrMap[level]; ok {
		return str, nil
	} else {
		return "", fmt.Errorf("unknow level %d", level)
	}
}
