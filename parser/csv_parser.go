package parser

import (
	"encoding/csv"
	"go-data-learning/model"
	"log"
	"os"
	"strconv"
)

// ParseCSVToStudents 读取 CSV 文件并解析为 Student 切片。
//
// 期望每行至少包含 4 列：姓名、年龄、城市、得分（按列顺序）。
// 若某行列数不足 4 列，将被跳过。
// 若年龄/得分无法解析为数字，将返回 error。
//
// 返回的切片顺序与 CSV 文件行顺序一致（不包含表头行）。
//
// 示例用法:
//
//	students, err := parser.ParseCSVToStudents("data.csv")
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
