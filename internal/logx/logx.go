package logx

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	defaultLogDir          = "logs"
	summaryLogName         = "gmcc.log"
	eventLogName           = "gmcc-events.jsonl"
	defaultEventMaxSize    = 5 * 1024 * 1024
	defaultArchiveMaxFiles = 5
)

var (
	mu            sync.Mutex
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	fileWriter    *boundedRotatingWriter
	eventWriter   *structuredWriter
	debugEnabled  bool
)

// Init 初始化日志系统
func Init(logDir string, enableFile bool, maxSize int64, debug bool) error {
	mu.Lock()
	defer mu.Unlock()

	debugEnabled = debug
	_ = closeLocked()

	consoleLogger = log.New(os.Stdout, "", 0)
	consoleLogger.SetOutput(os.Stdout)

	if !enableFile {
		return nil
	}

	if logDir == "" {
		logDir = defaultLogDir
	}

	summaryPath := filepath.Join(logDir, summaryLogName)
	summaryWriter, err := newBoundedRotatingWriter(summaryPath, maxSize, defaultArchiveMaxFiles)
	if err != nil {
		return err
	}

	fileWriter = summaryWriter
	fileLogger = log.New(fileWriter, "", log.LstdFlags)

	structured, err := newStructuredWriter(filepath.Join(logDir, eventLogName), defaultEventMaxSize, defaultArchiveMaxFiles)
	if err != nil {
		Warnf("event log channel disabled: %v", err)
		return nil
	}

	eventWriter = structured
	return nil
}

// Close 关闭日志系统
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	return closeLocked()
}

// Infof 输出信息日志
func Infof(format string, args ...interface{}) {
	Summaryf("info", format, args...)
}

// Warnf 输出警告日志
func Warnf(format string, args ...interface{}) {
	Summaryf("warn", format, args...)
}

// Errorf 输出错误日志
func Errorf(format string, args ...interface{}) {
	Summaryf("error", format, args...)
}

// Summaryf 输出摘要日志到控制台和摘要文件。
func Summaryf(level string, format string, args ...interface{}) {
	logf(normalizeSummaryLevel(level), true, format, args...)
}

// Debugf 输出调试日志（仅在debug模式）
func Debugf(format string, args ...interface{}) {
	mu.Lock()
	debug := debugEnabled
	mu.Unlock()
	if !debug {
		return
	}
	logf("DEBUG", true, format, args...)
}

// PacketLogf 数据包日志（空实现）
func PacketLogf(format string, args ...interface{}) {}

// Emit 写入结构化事件通道。
func Emit(event Event) {
	if err := event.Validate(); err != nil {
		Warnf("event dropped: %v", err)
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		Warnf("event dropped: marshal failed: %v", err)
		return
	}

	data = append(data, '\n')

	mu.Lock()
	writer := eventWriter
	mu.Unlock()

	if writer == nil {
		return
	}

	if _, err := writer.Write(data); err != nil {
		Warnf("event write degraded: %v", err)
	}
}

// LogTokenCache 记录token缓存状态
func LogTokenCache(tokenType string, profileName string, profileID string) {
	if profileName != "" {
		Infof("使用缓存的 %s token: %s (%s)", tokenType, profileName, profileID)
	} else {
		Infof("使用缓存的 %s token", tokenType)
	}
}

// LogTokenExpired 记录token过期状态
func LogTokenExpired(tokenType string, err error) {
	if err != nil {
		Warnf("缓存 %s token 已失效: %v", tokenType, err)
	} else {
		Warnf("缓存 %s token 已失效", tokenType)
	}
}

func logf(level string, writeFile bool, format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	if consoleLogger == nil {
		consoleLogger = log.New(os.Stdout, "", 0)
	}

	now := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	consoleLogger.Printf("%s [%s] %s", now, level, msg)

	if writeFile && fileLogger != nil {
		fileLogger.Printf("[%s] %s", level, msg)
	}
}

func normalizeSummaryLevel(level string) string {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		return "DEBUG"
	case "WARN", "WARNING":
		return "WARN"
	case "ERROR":
		return "ERROR"
	default:
		return "INFO"
	}
}

func closeLocked() error {
	fileLogger = nil
	var errs []error

	if eventWriter != nil {
		if err := eventWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		eventWriter = nil
	}

	if fileWriter != nil {
		if err := fileWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		fileWriter = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}

type boundedRotatingWriter struct {
	mu                 sync.Mutex
	activePath         string
	archiveBaseName    string
	archiveExt         string
	maxSize            int64
	maxFiles           int
	file               *os.File
	size               int64
	rotationWarned     bool
	writeFailureWarned bool
}

func newBoundedRotatingWriter(activePath string, maxSize int64, maxFiles int) (*boundedRotatingWriter, error) {
	if activePath == "" {
		activePath = filepath.Join(defaultLogDir, summaryLogName)
	}

	if err := os.MkdirAll(filepath.Dir(activePath), 0o755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	ext := filepath.Ext(activePath)
	baseName := strings.TrimSuffix(filepath.Base(activePath), ext)
	if ext == "" {
		ext = ".log"
	}

	w := &boundedRotatingWriter{
		activePath:      activePath,
		archiveBaseName: baseName,
		archiveExt:      ext,
		maxSize:         maxSize,
		maxFiles:        maxFiles,
	}

	if err := w.openOrRotateOnInitLocked(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *boundedRotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		if err := w.openOrRotateOnInitLocked(); err != nil {
			return 0, err
		}
	}

	if w.maxSize > 0 && w.size+int64(len(p)) > w.maxSize {
		if err := w.rotateLocked(); err != nil {
			w.warnOnceLocked(&w.rotationWarned, "log rotation disabled: %v", err)
			return w.writeActiveLocked(p)
		}
	}

	return w.writeActiveLocked(p)
}

func (w *boundedRotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}

	err := w.file.Close()
	w.file = nil
	w.size = 0
	return err
}

func (w *boundedRotatingWriter) writeActiveLocked(p []byte) (int, error) {
	if w.file == nil {
		return 0, fmt.Errorf("log writer is closed")
	}

	n, err := w.file.Write(p)
	w.size += int64(n)
	if err != nil {
		w.warnOnceLocked(&w.writeFailureWarned, "log file write degraded to console-only: %v", err)
	}
	return n, err
}

func (w *boundedRotatingWriter) openOrRotateOnInitLocked() error {
	info, err := os.Stat(w.activePath)
	if err == nil {
		if w.maxSize > 0 && info.Size() >= w.maxSize {
			if err := w.rotateActiveFileLocked(); err != nil {
				return err
			}
		} else {
			w.size = info.Size()
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("读取日志文件失败: %w", err)
	}

	f, err := os.OpenFile(w.activePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	st, statErr := f.Stat()
	if statErr != nil {
		_ = f.Close()
		return fmt.Errorf("读取日志文件状态失败: %w", statErr)
	}

	w.file = f
	w.size = st.Size()
	return nil
}

func (w *boundedRotatingWriter) rotateLocked() error {
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return fmt.Errorf("关闭日志文件失败: %w", err)
		}
		w.file = nil
		w.size = 0
	}

	if err := w.rotateActiveFileLocked(); err != nil {
		return err
	}

	if err := w.openOrRotateOnInitLocked(); err != nil {
		return err
	}

	w.rotationWarned = false
	w.writeFailureWarned = false
	return nil
}

func (w *boundedRotatingWriter) rotateActiveFileLocked() error {
	if _, err := os.Stat(w.activePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("读取日志文件失败: %w", err)
	}

	archivePath := filepath.Join(filepath.Dir(w.activePath), archiveName(w.archiveBaseName, w.archiveExt, time.Now()))
	if err := os.Rename(w.activePath, archivePath); err != nil {
		return fmt.Errorf("滚动日志文件失败: %w", err)
	}

	if err := pruneArchives(filepath.Dir(w.activePath), w.archiveBaseName, w.archiveExt, w.maxFiles); err != nil {
		return fmt.Errorf("清理归档日志失败: %w", err)
	}

	return nil
}

func (w *boundedRotatingWriter) warnOnceLocked(flag *bool, format string, args ...interface{}) {
	if *flag {
		return
	}
	*flag = true
	if consoleLogger == nil {
		consoleLogger = log.New(os.Stdout, "", 0)
	}
	consoleLogger.Printf("%s [WARN] %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

type structuredWriter struct {
	mu         sync.Mutex
	activePath string
	baseName   string
	ext        string
	maxSize    int64
	maxFiles   int
	file       *os.File
	size       int64
	closed     bool
}

func newStructuredWriter(activePath string, maxSize int64, maxFiles int) (*structuredWriter, error) {
	if err := os.MkdirAll(filepath.Dir(activePath), 0o755); err != nil {
		return nil, fmt.Errorf("create structured log directory: %w", err)
	}

	ext := filepath.Ext(activePath)
	baseName := strings.TrimSuffix(filepath.Base(activePath), ext)
	if ext == "" {
		ext = ".jsonl"
	}

	w := &structuredWriter{
		activePath: activePath,
		baseName:   baseName,
		ext:        ext,
		maxSize:    maxSize,
		maxFiles:   maxFiles,
	}

	if err := w.openLocked(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *structuredWriter) openLocked() error {
	file, err := os.OpenFile(w.activePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open structured log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("stat structured log file: %w", err)
	}

	w.file = file
	w.size = info.Size()
	return nil
}

func (w *structuredWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w == nil || w.file == nil || w.closed {
		return 0, fmt.Errorf("structured writer is closed")
	}

	if w.maxSize > 0 && w.size+int64(len(p)) > w.maxSize {
		if err := w.rotateLocked(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *structuredWriter) rotateLocked() error {
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return fmt.Errorf("close structured log file: %w", err)
		}
		w.file = nil
		w.size = 0
	}

	archivePath := filepath.Join(filepath.Dir(w.activePath), archiveName(w.baseName, w.ext, time.Now()))
	if err := os.Rename(w.activePath, archivePath); err != nil {
		return fmt.Errorf("rotate structured log file: %w", err)
	}

	if err := pruneArchives(filepath.Dir(w.activePath), w.baseName, w.ext, w.maxFiles); err != nil {
		return fmt.Errorf("prune structured log archives: %w", err)
	}

	return w.openLocked()
}

func (w *structuredWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w == nil || w.file == nil || w.closed {
		return nil
	}

	err := w.file.Close()
	w.file = nil
	w.size = 0
	w.closed = true
	return err
}

func archiveName(base string, ext string, now time.Time) string {
	if ext == "" {
		ext = ".log"
	}
	return fmt.Sprintf("%s-%s%s", base, now.UTC().Format("20060102-150405.000000000"), ext)
}

func pruneArchives(dir string, base string, ext string, maxFiles int) error {
	if maxFiles <= 0 {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	prefix := base + "-"
	var archives []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ext) {
			archives = append(archives, filepath.Join(dir, name))
		}
	}

	if len(archives) <= maxFiles {
		return nil
	}

	sort.Strings(archives)
	for _, archive := range archives[:len(archives)-maxFiles] {
		if err := os.Remove(archive); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
