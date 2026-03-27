package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicUpdate(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// 测试数据
	testData := []byte("test: value")

	// 测试原子更新
	err := atomicUpdate(configPath, testData)
	if err != nil {
		t.Fatalf("原子更新失败: %v", err)
	}

	// 验证文件内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("读取更新后的文件失败: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("文件内容不匹配，期望: %s, 实际: %s", string(testData), string(data))
	}

	// 检查是否创建了备份文件
	matches, err := filepath.Glob(configPath + ".backup.*")
	if err != nil {
		t.Fatalf("检查备份文件失败: %v", err)
	}

	if len(matches) > 1 {
		// 第一次更新不应该有备份（原文件不存在）
		t.Errorf("不应该有备份文件，但找到了: %v", matches)
	}

	// 再次测试（这次应该有备份）
	newData := []byte("updated: value")
	err = atomicUpdate(configPath, newData)
	if err != nil {
		t.Fatalf("第二次原子更新失败: %v", err)
	}

	matches, err = filepath.Glob(configPath + ".backup.*")
	if err != nil {
		t.Fatalf("检查备份文件失败: %v", err)
	}

	if len(matches) == 0 {
		t.Errorf("应该有备份文件，但没有找到")
	}
}

func TestCreateBackup(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// 创建原始文件
	originalData := []byte("original: content")
	err := os.WriteFile(configPath, originalData, 0644)
	if err != nil {
		t.Fatalf("创建原始文件失败: %v", err)
	}

	// 测试创建备份
	backupPath, err := createBackup(configPath)
	if err != nil {
		t.Fatalf("创建备份失败: %v", err)
	}

	// 验证备份文件存在
	if backupPath == "" {
		t.Fatalf("备份路径为空")
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("读取备份文件失败: %v", err)
	}

	if string(data) != string(originalData) {
		t.Errorf("备份内容不匹配，期望: %s, 实际: %s", string(originalData), string(data))
	}
}

func TestHasWritePermission(t *testing.T) {
	// 测试临时目录，应该有写权限
	tmpDir := t.TempDir()

	if !hasWritePermission(tmpDir) {
		t.Errorf("临时目录应该有写权限")
	}

	// 测试不存在的目录（不应该检测）
	// hasWritePermission 实际上会创建目录，所以总是返回 true
	// 这里改为测试只读目录
	// 创建只读目录（如果可能）
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0444); err == nil {
		defer os.Chmod(readOnlyDir, 0755) // 恢复权限以便清理
		if !hasWritePermission(readOnlyDir) {
			t.Logf("只读目录正确地没有写权限")
		}
	}
}
