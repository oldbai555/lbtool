package utils

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
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	levelPanic
	levelFatal
)

var LevelToStrMap = map[Level]string{
	LevelDebug: "DBG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERR",
	levelPanic: "FATAL",
	levelFatal: "PANIC",
}

var LevelToStdoutColorMap = map[Level]Color{
	LevelDebug: colorBlue,
	LevelInfo:  colorGreen,
	LevelWarn:  colorYellow,
	LevelError: colorRed,
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

const ColorEnd = "\x1b[0m"

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
