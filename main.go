package main

import (
	"fmt"
	"github.com/xianyudd/hanzi-data-kit/generator"
	"github.com/xianyudd/hanzi-data-kit/model"
	"github.com/xianyudd/hanzi-data-kit/parser"
	"path/filepath"
)

func main() {
	filename := "students.csv"
	outPath := filepath.Join("data", filename)
	const (
		n    = 1000 // 生成1000条数据
		seed = 42   // 固定随机种子，确保每次生成的结果一致，便于测试和调试
	)

	fmt.Println(">>> 开始写入csv文件:", outPath)
	headers := model.StudentHeadersCN() // 获取中文表头，未来我们也可以添加 StudentHeadersEN() 来支持英文表头

	gen := generator.NewStudentGenerator(generator.StudentGenConfig{
		Seed:     seed,
		AgeMin:   18,
		AgeMax:   30,
		ScoreMin: 60,
		ScoreMax: 100,
		// Cities/Surnames/GivenNames* 留空则使用 generator 内置默认值
	})

	// 先生成 model.Student，再由 model 层统一负责“结构体 -> CSV 行”的映射。
	students := make([]model.Student, 0, n)
	for i := 0; i < n; i++ {
		students = append(students, gen.Next())
	}
	totalRows := len(students)
	rowGenerator := func(i int) []string {
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
	parsed, err := parser.ParseCSVToStudents(outPath)
	if err != nil {
		fmt.Println("解析CSV失败:", err)
		return
	}

	fmt.Printf("成功解析了 %d 条学生数据 (展示前 5 条):\n", len(parsed))
	limit := 5
	if len(parsed) < limit {
		limit = len(parsed)
	}
	for i := 0; i < limit; i++ {
		stu := parsed[i]
		fmt.Printf("- %s (年龄: %d, 分数: %.1f)\n", stu.Name, stu.Age, stu.Score)
	}
	// 未来添加新功能时，只需要在这里继续调用：
	// parser.ParseExcelToStruct("students.xlsx")
	// db.SaveStudents(students)
}
