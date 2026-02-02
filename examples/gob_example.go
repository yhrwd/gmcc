// Package main demonstrates the usage of the enhanced gob utility functions.
package main

import (
	"fmt"
	"log"
	"os"

	"gmcc/pkg/file"
)

// Person 是一个示例结构体，用于演示gob编码/解码
type Person struct {
	Name string
	Age  int
	City string
}

// DeleteFile 是一个辅助函数，用于删除文件
func DeleteFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
}

func main() {
	fmt.Println("Gob Utility Functions 示例")

	// 清理之前的文件
	DeleteFile("person.gob")
	DeleteFile("person.gob.backup")
	DeleteFile("multiple.gob")

	// 示例1: 基本的写入和读取
	fmt.Println("\n=== 示例1: 基本的写入和读取 ===")
	person1 := Person{Name: "Alice", Age: 30, City: "Beijing"}

	// 写入到文件
	err := rwfile.WriteToGobFile("person.gob", person1)
	if err != nil {
		log.Fatal("写入失败:", err)
	}
	fmt.Println("成功将数据写入 person.gob")

	// 从文件读取
	var loadedPerson Person
	err = rwfile.ReadFromGobFile("person.gob", &loadedPerson)
	if err != nil {
		log.Fatal("读取失败:", err)
	}
	fmt.Printf("从文件读取的数据: %+v\n", loadedPerson)

	// 示例2: 编码为字节切片和解码
	fmt.Println("\n=== 示例2: 编码为字节切片和解码 ===")
	person2 := Person{Name: "Bob", Age: 25, City: "Shanghai"}

	// 编码为字节切片
	data, err := rwfile.EncodeToBytes(person2)
	if err != nil {
		log.Fatal("编码失败:", err)
	}
	fmt.Printf("编码后的字节长度: %d\n", len(data))

	// 从字节切片解码
	var decodedPerson Person
	err = rwfile.DecodeFromBytes(data, &decodedPerson)
	if err != nil {
		log.Fatal("解码失败:", err)
	}
	fmt.Printf("解码后的数据: %+v\n", decodedPerson)

	// 示例3: 追加数据到文件
	fmt.Println("\n=== 示例3: 追加数据到文件 ===")
	// 使用相同类型的数据进行演示
	numbers := []int{1, 2, 3}
	err = rwfile.WriteToGobFile("multiple.gob", numbers)
	if err != nil {
		log.Fatal("写入失败:", err)
	}

	// 追加更多相同类型的数据
	err = rwfile.AppendToGobFile("multiple.gob", []int{4, 5, 6})
	if err != nil {
		log.Fatal("追加失败:", err)
	}
	fmt.Println("成功追加数据到 multiple.gob")

	// 读取最后一条记录
	var lastRecord []int
	err = rwfile.ReadLastFromGobFile("multiple.gob", &lastRecord)
	if err != nil {
		log.Fatal("读取最后记录失败:", err)
	}
	fmt.Printf("文件中的最后一条记录: %v\n", lastRecord)

	// 示例4: 带备份的写入
	fmt.Println("\n=== 示例4: 带备份的写入 ===")
	person3 := Person{Name: "Charlie", Age: 35, City: "Guangzhou"}

	// 首先写入一些数据
	err = rwfile.WriteToGobFile("person.gob", person1)
	if err != nil {
		log.Fatal("写入失败:", err)
	}

	// 然后使用带备份的写入
	err = rwfile.WriteToGobFileWithBackup("person.gob", person3)
	if err != nil {
		log.Fatal("带备份写入失败:", err)
	}
	fmt.Println("成功使用备份方式写入数据")

	// 检查备份文件是否存在
	if _, err := os.Stat("person.gob.backup"); err == nil {
		fmt.Println("备份文件已创建: person.gob.backup")
	}

	// 示例5: 获取文件大小
	fmt.Println("\n=== 示例5: 获取文件大小 ===")
	size, err := rwfile.GetGobFileSize("person.gob")
	if err != nil {
		log.Fatal("获取文件大小失败:", err)
	}
	fmt.Printf("person.gob 文件大小: %d 字节\n", size)

	fmt.Println("\n所有示例完成!")
}
