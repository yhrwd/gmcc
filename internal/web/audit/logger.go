package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gmcc/internal/logx"
	"gmcc/internal/webtypes"
)

// Logger 审计日志管理器
type Logger struct {
	logDir    string
	retention time.Duration
	mu        sync.Mutex
}

// NewLogger 创建审计日志管理器
func NewLogger(logDir string, retentionDays int) (*Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit log directory: %w", err)
	}

	return &Logger{
		logDir:    logDir,
		retention: time.Duration(retentionDays) * 24 * time.Hour,
	}, nil
}

// Log 记录操作日志
func (l *Logger) Log(log *webtypes.OperationLog) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 设置默认值
	if log.ID == "" {
		log.ID = generateLogID()
	}
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	// 生成文件名：YYYY-MM-DD.jsonl
	filename := filepath.Join(l.logDir, log.Timestamp.Format("2006-01-02")+".jsonl")

	// 序列化日志条目
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// 追加写入文件
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	logx.Debugf("审计日志已记录: action=%s, password_id=%s, success=%v", log.Action, log.PasswordID, log.Success)
	return nil
}

// Query 查询日志
func (l *Logger) Query(start, end time.Time, passwordID string) ([]webtypes.OperationLog, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var logs []webtypes.OperationLog

	// 遍历日期范围内的所有日志文件
	for date := start; !date.After(end); date = date.Add(24 * time.Hour) {
		filename := filepath.Join(l.logDir, date.Format("2006-01-02")+".jsonl")

		data, err := os.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue // 文件不存在，跳过
			}
			return nil, fmt.Errorf("failed to read log file %s: %w", filename, err)
		}

		// 解析每一行
		lines := splitLines(string(data))
		for _, line := range lines {
			if line == "" {
				continue
			}

			var log webtypes.OperationLog
			if err := json.Unmarshal([]byte(line), &log); err != nil {
				logx.Warnf("解析日志条目失败: %v", err)
				continue
			}

			// 过滤条件
			if log.Timestamp.Before(start) || log.Timestamp.After(end) {
				continue
			}
			if passwordID != "" && log.PasswordID != passwordID {
				continue
			}

			logs = append(logs, log)
		}
	}

	return logs, nil
}

// Rotate 日志轮转
// 删除超过保留期限的旧日志文件
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-l.retention)

	entries, err := os.ReadDir(l.logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 解析文件名中的日期
		name := entry.Name()
		if filepath.Ext(name) != ".jsonl" {
			continue
		}

		dateStr := name[:len(name)-6] // 去掉 ".jsonl"
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // 跳过无法解析的文件
		}

		if date.Before(cutoff) {
			path := filepath.Join(l.logDir, name)
			if err := os.Remove(path); err != nil {
				logx.Warnf("删除旧日志文件失败: %s, %v", path, err)
			} else {
				logx.Infof("已删除旧日志文件: %s", path)
			}
		}
	}

	return nil
}

// GetLogFilePath 获取指定日期的日志文件路径
func (l *Logger) GetLogFilePath(date time.Time) string {
	return filepath.Join(l.logDir, date.Format("2006-01-02")+".jsonl")
}

// generateLogID 生成日志ID
func generateLogID() string {
	return fmt.Sprintf("log-%d-%d", time.Now().UnixNano(), os.Getpid())
}

// splitLines 分割字符串为行
func splitLines(s string) []string {
	var lines []string
	var start int
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// GetRetention 获取日志保留期限
func (l *Logger) GetRetention() time.Duration {
	return l.retention
}

// GetLogDir 获取日志目录
func (l *Logger) GetLogDir() string {
	return l.logDir
}
