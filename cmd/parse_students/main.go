package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xianyudd/hanzi-data-kit/parser"
)

func main() {
	var (
		in          = flag.String("in", "data/students.csv", "输入CSV路径")
		printN      = flag.Int("print", 5, "打印前N条（0表示不打印）")
		skipBadRows = flag.Bool("skip-bad-rows", true, "遇到坏行是否跳过（否则严格报错）")
		trimSpace   = flag.Bool("trim-space", true, "是否对字段做TrimSpace")
		allowBOM    = flag.Bool("allow-bom", true, "是否剥离UTF-8 BOM")
	)
	flag.Parse()

	students, err := parser.ParseCSVToStudentsWithOptions(*in, parser.CSVParseOptions{
		TrimSpace:   *trimSpace,
		AllowBOM:    *allowBOM,
		SkipBadRows: *skipBadRows,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析CSV失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功解析 %d 条学生数据 <<< %s\n", len(students), *in)

	if *printN <= 0 {
		return
	}
	if *printN > len(students) {
		*printN = len(students)
	}
	fmt.Printf("展示前 %d 条:\n", *printN)
	for i := 0; i < *printN; i++ {
		stu := students[i]
		// 注意：这里打印格式你可以选择 %.1f 以对齐 CSV 的 1 位小数
		fmt.Printf("- %s (年龄: %d, 城市: %s, 分数: %.1f)\n", stu.Name, stu.Age, stu.City, stu.Score)
	}
}
