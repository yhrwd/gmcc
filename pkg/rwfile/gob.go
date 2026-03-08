package rwfile

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// WriteToGobFile 将任意类型的数据序列化为gob格式并写入文件
// 使用示例:
//
//	type Person struct { Name string; Age int }
//	person := Person{Name: "Alice", Age: 30}
//	err := WriteToGobFile("data.gob", person)
func WriteToGobFile(filePath string, data interface{}) error {
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 创建gob编码器
	encoder := gob.NewEncoder(file)

	// 编码数据
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	return nil
}

// ReadFromGobFile 从gob文件中读取数据并反序列化为指定类型的变量
// 使用示例:
//
//	var person Person
//	err := ReadFromGobFile("data.gob", &person)
func ReadFromGobFile(filePath string, data interface{}) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 创建gob解码器
	decoder := gob.NewDecoder(file)

	// 解码数据
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	return nil
}

// EncodeToBytes 将任意类型的数据序列化为gob格式的字节切片
// 使用示例:
//
//	person := Person{Name: "Alice", Age: 30}
//	data, err := EncodeToBytes(person)
func EncodeToBytes(data interface{}) ([]byte, error) {
	// 创建缓冲区
	buf := new(bytes.Buffer)

	// 创建gob编码器
	encoder := gob.NewEncoder(buf)

	// 编码数据
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}

	// 返回字节切片
	return buf.Bytes(), nil
}

// DecodeFromBytes 将gob格式的字节切片反序列化为指定类型的变量
// 使用示例:
//
//	var person Person
//	err := DecodeFromBytes(data, &person)
func DecodeFromBytes(data []byte, target interface{}) error {
	// 检查输入数据
	if len(data) == 0 {
		return fmt.Errorf("input data is empty")
	}

	// 创建缓冲区
	buf := bytes.NewBuffer(data)

	// 创建gob解码器
	decoder := gob.NewDecoder(buf)

	// 解码数据
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	return nil
}

// AppendToGobFile 将数据追加到现有的gob文件中
// 注意: 这会将数据作为独立记录追加，而不是合并数据结构
// 使用示例:
//
//	additionalData := []int{4, 5, 6}
//	err := AppendToGobFile("data.gob", additionalData)
func AppendToGobFile(filePath string, data interface{}) error {
	// 以追加模式打开文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()

	// 创建gob编码器
	encoder := gob.NewEncoder(file)

	// 编码并追加数据
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode and append data: %w", err)
	}

	return nil
}

// WriteToGobFileWithBackup 在写入新数据之前创建备份文件
// 使用示例:
//
//	person := Person{Name: "Bob", Age: 25}
//	err := WriteToGobFileWithBackup("data.gob", person)
func WriteToGobFileWithBackup(filePath string, data interface{}) error {
	// 如果原文件存在，创建备份
	if _, err := os.Stat(filePath); err == nil {
		backupPath := filePath + ".backup"
		if err := copyFile(filePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// 写入新数据
	return WriteToGobFile(filePath, data)
}

// ReadLastFromGobFile 从gob文件中读取最后一条记录
// 适用于包含多个记录的gob文件
// 使用示例:
//
//	var lastRecord Person
//	err := ReadLastFromGobFile("data.gob", &lastRecord)
func ReadLastFromGobFile(filePath string, data interface{}) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 创建gob解码器
	decoder := gob.NewDecoder(file)

	// 逐个解码记录直到最后一个
	var lastErr error
	for {
		err := decoder.Decode(data)
		if err == io.EOF {
			// 到达文件末尾，这是正常的
			break
		}
		if err != nil {
			// 其他错误
			lastErr = err
			break
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to decode last record: %w", lastErr)
	}

	return nil
}

// GetGobFileSize 获取gob文件的大小（字节）
// 使用示例:
//
//	size, err := GetGobFileSize("data.gob")
func GetGobFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	return info.Size(), nil
}

// copyFile 是一个辅助函数，用于复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
