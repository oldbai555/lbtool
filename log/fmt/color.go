package fmt

import "fmt"

type Color int

const (
	ColorNil Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorPurple
)

var ColorToStdoutMap = map[Color]string{
	ColorNil:    fmt.Sprintf("\u001B[0m"),
	ColorRed:    fmt.Sprintf("\x1b[%dm", 91),
	ColorGreen:  fmt.Sprintf("\x1b[%dm", 92),
	ColorYellow: fmt.Sprintf("\x1b[%dm", 93),
	ColorBlue:   fmt.Sprintf("\u001B[%d;1m", 36),
	ColorPurple: fmt.Sprintf("\x1b[%dm", 95),
}

var LevelToStdoutColorMap = map[Level]Color{
	LevelDebug: ColorBlue,
	LevelInfo:  ColorGreen,
	LevelWarn:  ColorYellow,
	LevelError: ColorRed,
	LevelFatal: ColorPurple,
	LevelPanic: ColorRed,
}

func GetColorStdout(color Color) (string, error) {
	if stdout, ok := ColorToStdoutMap[color]; ok {
		return stdout, nil
	}
	return "", fmt.Errorf("not support color type %d", color)
}
