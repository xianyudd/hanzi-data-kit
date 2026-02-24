package parser

import (
	"encoding/csv"
	"go-data-learning/model"
	"log"
	"os"
	"strconv"
)

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

		age, _ := strconv.Atoi(row[i])
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
