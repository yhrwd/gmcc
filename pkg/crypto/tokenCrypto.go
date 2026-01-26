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

func SaveToken(filename string, data interface{}) error {
	fullPath := getSessionPath(filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), dirPerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
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

func LoadToken(filename string, data interface{}) error {
	fullPath := getSessionPath(filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("token文件不存在: %s", filename)
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	if err := gob.NewDecoder(file).Decode(data); err != nil {
		return fmt.Errorf("读取数据失败: %w", err)
	}
	return nil
}
