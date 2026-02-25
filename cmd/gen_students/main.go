package main

import (
	"flag"
	"fmt"
	"go-data-learning/generator"
	"go-data-learning/model"
	"go-data-learning/parser"
	"os"
	"path/filepath"
)

func main() {
	var (
		n        = flag.Int("n", 1000, "生成学生数量")
		seed     = flag.Int64("seed", 42, "随机种子(用于复现)")
		out      = flag.String("out", "data/students.csv", "输出CSV文件路径")
		ageMin   = flag.Int("age-min", 18, "年龄下限(闭区间) ")
		ageMax   = flag.Int("age-max", 30, "年龄上限(闭区间)")
		scoreMin = flag.Float64("score-min", 60, "得分下限(闭区间)")
		scoreMax = flag.Float64("score-max", 100, "得分上限(闭区间)")
	)
	flag.Parse()

	if *n <= 0 {
		fmt.Fprintf(os.Stderr, "参数错误: -n 必须 > 0")
		os.Exit(2)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	gen := generator.NewStudentGenerator(generator.StudentGenConfig{
		Seed:     *seed,
		AgeMin:   *ageMin,
		AgeMax:   *ageMax,
		ScoreMin: *scoreMin,
		ScoreMax: *scoreMax,
	})

	students := make([]model.Student, 0, *n)
	for i := 0; i < *n; i++ {
		students = append(students, gen.Next())
	}

	headers := model.StudentHeadersCN()
	totalRows := len(students)
	rowGenerator := func(i int) []string {
		return model.StudentToRowCN(students[i-1])
	}

	if err := parser.WriteLargeCSV(*out, headers, totalRows, rowGenerator); err != nil {
		fmt.Fprintf(os.Stderr, "写入CSV失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功生成并写入 %d 条学生数据到 >>> %s\n", totalRows, *out)
}
