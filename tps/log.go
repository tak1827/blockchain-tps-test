package tps

import (
	"log"
)

type LogLevel int

const (
	DEBUG_LEVEL LogLevel = iota
	INFO_LEVEL  LogLevel = iota
	WARN_LEVEL  LogLevel = iota
	FATAL_LEVEL LogLevel = iota
)

type Logger struct {
	level LogLevel
}

func NewLogger(level LogLevel) Logger {
	return Logger{level: level}
}

func (l Logger) Info(msg string) {
	if l.level > INFO_LEVEL {
		return
	}
	l.print("[INFO] " + msg)
}

func (l Logger) Warn(msg string) {
	if l.level > WARN_LEVEL {
		return
	}
	l.print("[WARN] " + msg)
}

func (l Logger) Fatal(args ...interface{}) {
	if l.level > FATAL_LEVEL {
		return
	}
	log.Fatal(args...)
}

func (l Logger) print(msg string) {
	log.Print(msg)
}
