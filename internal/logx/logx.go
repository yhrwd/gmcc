package logx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	mu            sync.Mutex
	consoleLogger = log.New(os.Stdout, "", 0)
	fileLogger    *log.Logger
	fileWriter    *rotatingFileWriter
	packetLogger  *log.Logger
	packetWriter  *rotatingFileWriter
	debugEnabled  bool
)

func Init(logDir string, enableFile bool, maxSize int64, debug bool) error {
	mu.Lock()
	defer mu.Unlock()

	debugEnabled = debug
	_ = closeLocked()

	if enableFile {
		w, err := newRotatingFileWriter(logDir, maxSize)
		if err != nil {
			return err
		}
		fileWriter = w
		fileLogger = log.New(fileWriter, "", log.LstdFlags)

		pw, err := newRotatingFileWriter(filepath.Join(logDir, "packets"), maxSize)
		if err != nil {
			return err
		}
		packetWriter = pw
		packetLogger = log.New(packetWriter, "", log.LstdFlags)
	}
	return nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()
	return closeLocked()
}

func Infof(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if len(args) == 0 {
		consoleLogger.Printf("%s", format)
	} else {
		consoleLogger.Printf(format, args...)
	}
	if fileLogger != nil {
		fileLogger.Printf("[INFO] "+format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	consoleLogger.Printf("\033[33m[WARN] "+format+"\033[0m", args...)
	if fileLogger != nil {
		fileLogger.Printf("[WARN] "+format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	consoleLogger.Printf("\033[31m[ERROR] "+format+"\033[0m", args...)
	if fileLogger != nil {
		fileLogger.Printf("[ERROR] "+format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if debugEnabled {
		now := time.Now().Format("15:04:05")
		consoleLogger.Printf("%s [DEBUG] %s", now, fmt.Sprintf(format, args...))
		if fileLogger != nil {
			fileLogger.Printf("[DEBUG] "+format, args...)
		}
	}
}

func PacketLogf(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if packetLogger != nil {
		packetLogger.Printf("[PACKET] "+format, args...)
	}
}

func closeLocked() error {
	fileLogger = nil
	packetLogger = nil
	var errs []error
	if fileWriter != nil {
		if err := fileWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		fileWriter = nil
	}
	if packetWriter != nil {
		if err := packetWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		packetWriter = nil
	}
	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}

type rotatingFileWriter struct {
	mu         sync.Mutex
	logDir     string
	activePath string
	maxSize    int64
	file       *os.File
	size       int64
}

func newRotatingFileWriter(logDir string, maxSize int64) (*rotatingFileWriter, error) {
	if logDir == "" {
		logDir = "logs"
	}
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	w := &rotatingFileWriter{
		logDir:     logDir,
		activePath: filepath.Join(logDir, "gmcc.log"),
		maxSize:    maxSize,
	}

	if err := w.openOrRotateOnInitLocked(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *rotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		if err := w.openOrRotateOnInitLocked(); err != nil {
			return 0, err
		}
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

func (w *rotatingFileWriter) Close() error {
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

func (w *rotatingFileWriter) openOrRotateOnInitLocked() error {
	info, err := os.Stat(w.activePath)
	if err == nil && w.maxSize > 0 && info.Size() >= w.maxSize {
		if err := w.rotateActiveFileLocked(); err != nil {
			return err
		}
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

func (w *rotatingFileWriter) rotateLocked() error {
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
		w.size = 0
	}
	if err := w.rotateActiveFileLocked(); err != nil {
		return err
	}
	return w.openOrRotateOnInitLocked()
}

func (w *rotatingFileWriter) rotateActiveFileLocked() error {
	if _, err := os.Stat(w.activePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("读取日志文件失败: %w", err)
	}

	base := "gmcc-" + time.Now().Format("20060102-150405")
	target := filepath.Join(w.logDir, base+".log")
	for i := 1; ; i++ {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			break
		}
		target = filepath.Join(w.logDir, fmt.Sprintf("%s-%03d.log", base, i))
	}

	if err := os.Rename(w.activePath, target); err != nil {
		return fmt.Errorf("滚动日志文件失败: %w", err)
	}
	return nil
}
