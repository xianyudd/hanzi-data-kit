package main

import (
	"fmt"
	"go-data-learning/parser"
	"path/filepath"
)

func main() {
	filename := "students.csv"
	outPath := filepath.Join("data", filename)

	fmt.Println(">>> 开始写入csv文件:", outPath)
	headers := []string{"姓名", "年龄", "城市", "得分"}

	// 我们先写死 3 行（让结果可控，便于验收）
	totalRows := 3

	// rowGenerator：根据行号返回每一行的 []string
	rowGenerator := func(i int) []string {
		switch i {
		case 1:
			return []string{"张三", "22", "北京", "95"}
		case 2:
			return []string{"李四", "25", "上海", "88"}
		case 3:
			return []string{"王五", "28", "广州", "92"}
		default:
			// 防御式：不应该发生
			return []string{"", "", "", ""}
		}
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
