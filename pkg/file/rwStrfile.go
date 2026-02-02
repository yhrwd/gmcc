package rwfile

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// 默认文件权限：所有者可读可写，其他用户只读（符合 Unix/Linux 规范）
const defaultFilePerm = 777

// WriteStringToFile 将字符串写入指定文件（覆盖写入模式）
// 若文件不存在则创建，若文件已存在则覆盖原有内容
// 参数：
//   filePath: 目标文件路径
//   content: 要写入的字符串内容
// 返回：
//   error: 错误信息（文件创建失败、写入失败等）
// 使用示例:
//   err := WriteStringToFile("example.txt", "Hello, World!")
func WriteStringToFile(filePath, content string) error {
	// 调用通用写入函数，使用覆盖模式
	return writeString(filePath, content, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, defaultFilePerm)
}

// AppendStringToFile 将字符串追加到指定文件末尾
// 若文件不存在则创建，若文件已存在则在末尾追加内容
// 参数：
//   filePath: 目标文件路径
//   content: 要追加的字符串内容
// 返回：
//   error: 错误信息（文件打开/创建失败、写入失败等）
// 使用示例:
//   err := AppendStringToFile("example.txt", "\nAppend this content.")
func AppendStringToFile(filePath, content string) error {
	// 调用通用写入函数，使用追加模式
	return writeString(filePath, content, os.O_APPEND|os.O_CREATE|os.O_WRONLY, defaultFilePerm)
}

// ReadStringFromFile 从指定文件读取全部字符串内容
// 适用于文本文件，若文件过大可能占用较多内存，建议对大文件使用分段读取
// 参数：
//   filePath: 目标文件路径
// 返回：
//   string: 文件中的字符串内容
//   error: 错误信息（文件打开失败、读取失败等）
// 使用示例:
//   content, err := ReadStringFromFile("example.txt")
func ReadStringFromFile(filePath string) (string, error) {
	// 先读取到字节切片，再转换为字符串
	data, err := ReadBytesFromFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadBytesFromFile 从指定文件读取全部字节内容
// 通用文件读取函数，适用于文本文件和二进制文件
// 参数：
//   filePath: 目标文件路径
// 返回：
//   []byte: 文件中的字节内容
//   error: 错误信息（文件打开失败、读取失败等）
// 使用示例:
//   data, err := ReadBytesFromFile("example.bin")
func ReadBytesFromFile(filePath string) ([]byte, error) {
	// 打开文件（只读模式）
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 获取文件大小，预分配缓冲区（提升大文件读取性能）
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()
	// 空文件直接返回空字节切片
	if fileSize == 0 {
		return []byte{}, nil
	}
	// 预分配缓冲区，避免后续频繁扩容
	buf := make([]byte, 0, fileSize)

	// 读取文件内容到缓冲区
	n, err := io.ReadFull(file, buf[:cap(buf)])
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return buf[:n], nil
}

// FileExists 判断指定文件是否存在（仅判断存在性，不区分文件/目录）
// 若需要区分文件和目录，可结合 os.Stat() 进一步判断
// 参数：
//   filePath: 目标文件路径
// 返回：
//   bool: true（存在）/ false（不存在）
// 使用示例:
//   exists := FileExists("example.txt")
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	// 忽略 "文件不存在" 错误，其他错误（如权限不足）也返回 false
	return !os.IsNotExist(err)
}

// ---------------------- 内部私有辅助函数 ----------------------

// writeString 通用字符串写入函数，封装写入逻辑，供对外接口调用
// flag: 文件打开模式（os.O_TRUNC/os.O_APPEND 等）
// perm: 文件权限（默认 0777）
func writeString(filePath, content string, flag int, perm os.FileMode) error {
	// 打开/创建文件（根据 flag 决定模式）
	file, err := os.OpenFile(filePath, flag, perm)
	if err != nil {
		return fmt.Errorf("failed to open/create file: %w", err)
	}
	defer file.Close()

	// 写入内容（直接使用 io.WriteString，比 bytes.NewBufferString 更高效）
	_, err = io.WriteString(file, content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// 强制刷新缓冲区到磁盘（确保内容立即写入，避免数据丢失）
	return file.Sync()
}