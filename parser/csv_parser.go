// Package parser 提供了项目中所有文件格式（如 CSV, Excel）的解析与导出工具。
// 它可以将各种非结构化文件映射为 model 层的强类型结构体。
package parser

import (
	"encoding/csv"
	"go-data-learning/model"
	"log"
	"os"
	"strconv"
)

// ParseCSVToStruct 读取指定路径的 CSV 文件并将其映射为结构体切片。
//
// 示例用法:
//
//	students, err := parser.ParseCSVToStruct("data.csv")
//	if err != nil { ... }
//
// 注意：该函数会自动忽略长度不足 4 列的脏数据行。
func ParseCSVToStudents(filename string) ([]model.Student, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("读取CSV失败: %v", err)
	}
	var students []model.Student

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 4 {
			log.Printf("列数不足，跳过第 %d 行: %v\n", i+1, row)
			continue
		}

		age, _ := strconv.Atoi(row[1])
		score, _ := strconv.ParseFloat(row[3], 64)
		stu := model.Student{
			Name:  row[0],
			Age:   age,
			City:  row[2],
			Score: score,
		}
		students = append(students, stu)
	}
	return students, nil // 返回解析好的结构体切片和 nil（表示没有错误）
}
