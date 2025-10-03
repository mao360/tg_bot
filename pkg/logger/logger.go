package logger

import (
	"log"
	"os"
)

type Logger struct {
	level string
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
	warn  *log.Logger
}

func New(level string) *Logger {
	return &Logger{
		level: level,
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.shouldLog("info") {
		l.info.Printf(msg, fields...)
	}
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.shouldLog("error") {
		l.error.Printf(msg, fields...)
	}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.shouldLog("debug") {
		l.debug.Printf(msg, fields...)
	}
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.shouldLog("warn") {
		l.warn.Printf(msg, fields...)
	}
}

func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}
	
	currentLevel, exists := levels[l.level]
	if !exists {
		currentLevel = 1 // default to info
	}
	
	requestedLevel, exists := levels[level]
	if !exists {
		return false
	}
	
	return requestedLevel >= currentLevel
}
