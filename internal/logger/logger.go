package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ----------------------------------------------------------------------------
// 日志级别
// ----------------------------------------------------------------------------

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ----------------------------------------------------------------------------
// 样式
// ----------------------------------------------------------------------------

var (
	debugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#B0BEC5"))
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#C8E6C9"))
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF59D"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCDD2"))
)

// ----------------------------------------------------------------------------
// Logger
// ----------------------------------------------------------------------------

type Logger struct {
	mu sync.Mutex

	// config
	logDir      string
	maxSize     int64
	enableDebug bool
	enableUI    bool
	enableFile  bool
	withTime    bool

	// file
	file     *os.File
	filePath string

	// ui callback
	uiLogFunc func(string)
}

// 全局实例
var defaultLogger = &Logger{
	withTime: true,
	enableUI: true,
}

// ----------------------------------------------------------------------------
// 全局 API
// ----------------------------------------------------------------------------

// Debug / Debugf
func Debug(args ...any) {
	defaultLogger.log(LevelDebug, "%s", fmt.Sprint(args...))
}

func Debugf(format string, args ...any) {
	defaultLogger.log(LevelDebug, format, args...)
}

// Info / Infof
func Info(args ...any) {
	defaultLogger.log(LevelInfo, "%s", fmt.Sprint(args...))
}

func Infof(format string, args ...any) {
	defaultLogger.log(LevelInfo, format, args...)
}

// Warn / Warnf
func Warn(args ...any) {
	defaultLogger.log(LevelWarn, "%s", fmt.Sprint(args...))
}

func Warnf(format string, args ...any) {
	defaultLogger.log(LevelWarn, format, args...)
}

// Error / Errorf
func Error(args ...any) {
	defaultLogger.log(LevelError, "%s", fmt.Sprint(args...))
}

func Errorf(format string, args ...any) {
	defaultLogger.log(LevelError, format, args...)
}

// Fatal / Fatalf  ← 新增致命级别（会退出程序）
func Fatal(args ...any) {
    defaultLogger.log(LevelError, "%s", fmt.Sprint(args...))
}

func Fatalf(format string, args ...any) {
    defaultLogger.log(LevelError, format, args...)
}

// ----------------------------------------------------------------------------
// 初始化
// ----------------------------------------------------------------------------

func InitLogger(logDir string, maxSize int64, enableDebug bool, enableFile bool) error {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()

	defaultLogger.logDir = logDir
	defaultLogger.maxSize = maxSize
	defaultLogger.enableDebug = enableDebug
	defaultLogger.enableFile = enableFile

	if !enableFile {
		return nil
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	defaultLogger.filePath = filepath.Join(logDir, "app.log")
	return defaultLogger.openFile()
}

func (l *Logger) openFile() error {
	f, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.file = f
	return nil
}

// ----------------------------------------------------------------------------
// 日志切割
// ----------------------------------------------------------------------------

func (l *Logger) rotateIfNeeded() {
	if l.file == nil || l.maxSize <= 0 {
		return
	}

	info, err := l.file.Stat()
	if err != nil {
		return
	}

	if info.Size() < l.maxSize {
		return
	}

	_ = l.file.Close()

	for i := 9; i >= 1; i-- {
		old := fmt.Sprintf("%s.%d", l.filePath, i)
		next := fmt.Sprintf("%s.%d", l.filePath, i+1)
		if _, err := os.Stat(old); err == nil {
			_ = os.Rename(old, next)
		}
	}
	_ = os.Rename(l.filePath, l.filePath+".1")

	_ = l.openFile()
}

// ----------------------------------------------------------------------------
// 核心输出
// ----------------------------------------------------------------------------

func (l *Logger) log(level LogLevel, format string, args ...any) {
    if level == LevelDebug && !l.enableDebug {
        return
    }

    l.mu.Lock()
    defer l.mu.Unlock()

    text := fmt.Sprintf(format, args...)
    ts := ""
    if l.withTime {
        ts = time.Now().Format("15:04:05")
    }

    isFatal := false
    levelStr := level.String()
    style := errorStyle // 默认 error 的样式

    // 如果是 Fatalf/Fatal，强制改成 FATAL 前缀 + 加粗
    if strings.Contains(format, "fatal") || strings.Contains(format, "Fatal") { // 简单的判断
        isFatal = true
        levelStr = "FATAL"
        style = errorStyle.Bold(true).Foreground(lipgloss.Color("#FF5252"))
    } else {
        switch level {
        case LevelDebug:
            style = debugStyle
        case LevelInfo:
            style = infoStyle
        case LevelWarn:
            style = warnStyle
        case LevelError:
            style = errorStyle
        }
    }

    colored := style.Render(text)

    // TUI 输出
    if l.enableUI && l.uiLogFunc != nil {
        var msg string
        if l.withTime {
            msg = fmt.Sprintf("[%s] [%s] %s", ts, levelStr, colored)
        } else {
            msg = fmt.Sprintf("[%s] %s", levelStr, colored)
        }
        l.uiLogFunc(msg)
    }

    if l.enableFile && l.file != nil {
        l.rotateIfNeeded()

        line := text
        if l.withTime {
            line = fmt.Sprintf("[%s] [%s] %s\n", ts, levelStr, text)
        } else {
            line = fmt.Sprintf("[%s] %s\n", levelStr, text)
        }
        _, _ = l.file.WriteString(line)
    }

    if isFatal {
        os.Exit(1)
    }
}
// ----------------------------------------------------------------------------
// 供 TUI 注册回调
// ----------------------------------------------------------------------------

func SetUILogFunc(fn func(string)) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.uiLogFunc = fn
}
