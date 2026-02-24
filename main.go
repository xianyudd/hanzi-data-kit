package main

import (
	"fmt"
	"go-data-learning/model"
	"go-data-learning/parser"
	"path/filepath"
)

func main() {
	filename := "students.csv"
	outPath := filepath.Join("data", filename)

	fmt.Println(">>> 开始写入csv文件:", outPath)
	headers := model.StudentHeadersCN() // 获取中文表头，未来我们也可以添加 StudentHeadersEN() 来支持英文表头

	students := []model.Student{
		{Name: "张三", Age: 22, City: "北京", Score: 95},
		{Name: "李四", Age: 25, City: "上海", Score: 88},
		{Name: "王五", Age: 28, City: "广州", Score: 92},
	}
	totalRows := len(students)
	rowGenerator := func (i int) []string {
		return model.StudentToRowCN(students[i-1])
	}


	//调用CSV 写入器
	if err := parser.WriteLargeCSV(outPath, headers, totalRows, rowGenerator); err != nil {
		// 先用 fmt 输出即可，后面我们会讲 log / error wrapping
		fmt.Println("写入失败:", err)
		return
	}

	fmt.Println("写入成功:", outPath)

	fmt.Println("<<< 开始读取csv文件:", outPath)
	students, err := parser.ParseCSVToStudents(outPath)
	if err != nil {
		fmt.Println("解析CSV失败:", err)
		return
	}
	fmt.Printf("成功解析了 %d 条学生数据:\n", len(students))
	for _, stu := range students {
		fmt.Printf("- %s (年龄: %d, 分数: %.2f )\n", stu.Name, stu.Age, stu.Score)
	}
	// 未来添加新功能时，只需要在这里继续调用：
	// parser.ParseExcelToStruct("students.xlsx")
	// db.SaveStudents(students)
}
