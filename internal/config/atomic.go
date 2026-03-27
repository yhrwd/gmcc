package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// atomicUpdate 原子性更新配置文件
func atomicUpdate(path string, data []byte) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建临时文件
	tmpPath := path + ".tmp." + time.Now().Format("20060102150405")
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}

	// 备份原文件（如果存在）
	var backupPath string
	if _, err := os.Stat(path); err == nil {
		backupPath = path + ".backup." + time.Now().Format("20060102150405")
		if err := os.Rename(path, backupPath); err != nil {
			os.Remove(tmpPath) // 清理临时文件
			return fmt.Errorf("备份原文件失败: %w", err)
		}
	}

	// 原子重命名
	if err := os.Rename(tmpPath, path); err != nil {
		// 尝试恢复备份
		if backupPath != "" {
			os.Rename(backupPath, path)
		}
		os.Remove(tmpPath) // 清理临时文件
		return fmt.Errorf("原子更新失败: %w", err)
	}

	return nil
}

// createBackup 创建配置文件备份
func createBackup(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil // 文件不存在，无需备份
	}

	backupPath := path + ".backup." + time.Now().Format("20060102150405")

	// 读取原文件
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取原文件失败: %w", err)
	}

	// 写入备份
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", fmt.Errorf("创建备份失败: %w", err)
	}

	return backupPath, nil
}

// hasWritePermission 检查是否有写入权限
func hasWritePermission(path string) bool {
	dir := filepath.Dir(path)

	// 尝试创建测试文件
	testFile := filepath.Join(dir, ".gmcc_config_test")

	// 尝试创建文件
	file, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	file.Close()

	// 清理测试文件
	os.Remove(testFile)
	return true
}

// getBackupDir 获取备份目录路径
func getBackupDir() string {
	return filepath.Join(".", ".config", "backups")
}

// cleanOldBackups 清理旧备份文件，保留最近的5个备份
func cleanOldBackups() error {
	backupDir := getBackupDir()
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nil // 备份目录不存在，无需清理
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("读取备份目录失败: %w", err)
	}

	// 按修改时间排序
	type FileInfo struct {
		Name    string
		ModTime time.Time
	}

	var fileInfos []FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, FileInfo{
			Name:    file.Name(),
			ModTime: info.ModTime(),
		})
	}

	// 按修改时间降序排序
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime.After(fileInfos[j].ModTime)
	})

	// 保留最新的5个，删除其余
	if len(fileInfos) > 5 {
		for i := 5; i < len(fileInfos); i++ {
			backupPath := filepath.Join(backupDir, fileInfos[i].Name)
			if err := os.Remove(backupPath); err != nil {
				// 删除失败时记录警告但不中断
				continue
			}
		}
	}

	return nil
}
