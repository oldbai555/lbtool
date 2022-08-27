package log

import "fmt"

// ================================Format===============================

// Format 格式化类型枚举
type Format int

const (
	FormatText Format = iota
)

// ================================Level===============================

// Level 日志等级
type Level int

const (
	levelDebug Level = iota
	levelInfo
	levelWarn
	levelError
	levelPanic
	levelFatal
)

var levelToStrMap = map[Level]string{
	levelDebug: "DBG",
	levelInfo:  "INFO",
	levelWarn:  "WARN",
	levelError: "ERR",
	levelPanic: "FATAL",
	levelFatal: "PANIC",
}

var levelToStdoutColorMap = map[Level]Color{
	levelDebug: colorBlue,
	levelInfo:  colorGreen,
	levelWarn:  colorYellow,
	levelError: colorRed,
	levelFatal: colorPurple,
	levelPanic: colorRed,
}

// ================================Color===============================

type Color uint8

const (
	colorRed = Color(iota + 91)
	colorGreen
	colorYellow
	colorBlue
	colorPurple
)

const colorEnd = "\x1b[0m"

var colorToStdoutMap = map[Color]string{
	colorRed:    fmt.Sprintf("\x1b[%dm", 91),
	colorGreen:  fmt.Sprintf("\x1b[%dm", 92),
	colorYellow: fmt.Sprintf("\x1b[%dm", 93),
	colorBlue:   fmt.Sprintf("\u001B[%d;1m", 36),
	colorPurple: fmt.Sprintf("\x1b[%dm", 95),
}

func GetColorStdout(color Color) (string, error) {
	if stdout, ok := colorToStdoutMap[color]; ok {
		return stdout, nil
	}
	return "", fmt.Errorf("not support color type %d", color)
}

// ==========================other============================
