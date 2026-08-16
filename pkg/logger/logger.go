// Package logger 提供一个简单、线程安全的分级日志器。
package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Level 日志级别。
type Level int

// 日志级别常量，数值越大越严重。
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ParseLevel 将字符串解析为日志级别。
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// String 返回级别的字符串表示。
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "info"
	}
}

// Logger 分级日志器。
type Logger struct {
	mu    sync.Mutex
	level Level
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

// New 创建默认级别（Info）的日志器，输出到标准错误。
func New() *Logger {
	return NewLevel(LevelInfo)
}

// NewLevel 创建指定级别的日志器。
func NewLevel(level Level) *Logger {
	flag := log.LstdFlags | log.Lmicroseconds
	return &Logger{
		level: level,
		debug: log.New(os.Stderr, "[DEBUG] ", flag),
		info:  log.New(os.Stderr, "[INFO ] ", flag),
		warn:  log.New(os.Stderr, "[WARN ] ", flag),
		err:   log.New(os.Stderr, "[ERROR] ", flag),
	}
}

// SetLevel 动态调整日志级别。
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Level 返回当前日志级别。
func (l *Logger) Level() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// Debugf 输出调试日志。
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelDebug {
		l.debug.Printf(format, args...)
	}
}

// Infof 输出信息日志。
func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelInfo {
		l.info.Printf(format, args...)
	}
}

// Warnf 输出警告日志。
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelWarn {
		l.warn.Printf(format, args...)
	}
}

// Errorf 输出错误日志。
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelError {
		l.err.Printf(format, args...)
	}
}

// Printf 兼容标准库，以 Info 级别输出。
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}

// Sprint 将参数拼接为字符串（供测试使用）。
func Sprint(args ...interface{}) string {
	return fmt.Sprint(args...)
}
