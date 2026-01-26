package crypto

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

const (
	sessionDir = ".session"
	filePerm   = 0600
	dirPerm    = 0700
)

func getSessionPath(filename string) string {
	execDir, err := os.Getwd()
	if err != nil {
		execDir = "."
	}
	return filepath.Join(execDir, sessionDir, filename)
}

// func getSessionPath(filename string) string {
	// exePath, err := os.Executable()
	// if err != nil {
		// exePath, _ = os.Getwd()
	// }

	// execDir := filepath.Dir(exePath)

	// return filepath.Join(execDir, sessionDir, filename)
// }

// SaveToken 保存任意数据到 .session 目录
func SaveToken(filename string, data interface{}) error {
	fullPath := getSessionPath(filename)

	// 自动创建目录
	if err := os.MkdirAll(filepath.Dir(fullPath), dirPerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 安全创建文件（覆盖模式）
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(data); err != nil {
		return fmt.Errorf("保存数据失败: %w", err)
	}

	return nil
}

// LoadToken 从 .session 目录加载数据。
// data 必须是一个指向目标类型的指针，函数会将读取到的数据直接填充到该指针指向的对象中
func LoadToken(filename string, data interface{}) error {
    fullPath := getSessionPath(filename)

    // 检查文件是否存在
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return fmt.Errorf("token文件不存在: %s", filename)
    }

    file, err := os.Open(fullPath)
    if err != nil {
        return fmt.Errorf("打开文件失败: %w", err)
    }
    defer file.Close()

    // 数据将写入 data 指向的内存地址
    if err := gob.NewDecoder(file).Decode(data); err != nil {
        return fmt.Errorf("读取数据失败: %w", err)
    }

    return nil
}
