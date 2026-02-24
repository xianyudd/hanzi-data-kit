package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "data/students.csv"

	fmt.Println(">>> 开始写入csv文件:", filename)
	writeCsv(filename)

	fmt.Println("<<< 开始读取csv文件:", filename)
	readCsv(filename)
}

func writeCsv(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("创建文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := [][]string{
		{"姓名", "年龄", "城市", "得分"},
		{"张三", "22", "北京", "95"},
		{"李四", "25", "上海", "88"},
		{"王五", "28", "广州", "92"},
	}

	err = writer.WriteAll(data)
	if err != nil {
		log.Fatalf("写入数据失败: %v", err)
	}
	fmt.Printf("成功将 %d 行数据写入文件 %s\n", len(data), filename)
}

func readCsv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("读取数据失败: %v", err)
	}

	for rowIndex, row := range records {
		fmt.Printf("第 %d 行: %v\n", rowIndex+1, row)
	}

}
