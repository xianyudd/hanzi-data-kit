package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/xianyudd/hanzi-data-kit/model"
	"log"
	"os"
	"strconv"
	"strings"
)

// CSVParseOptions 控制 CSV 解析行为。
// 解析默认策略由 defaultCSVParseOptions() 给出；若你需要宽松解析，可将 SkipBadRows 设为 true。
type CSVParseOptions struct {
	// TrimSpace 是否对表头和每个单元格做 strings.TrimSpace。
	TrimSpace bool

	// AllowBOM 是否允许并自动剥离 UTF-8 BOM（常见于 Excel 导出的 CSV）。
	AllowBOM bool

	// SkipBadRows 为 true 时，遇到脏数据行会跳过继续解析；
	// 为 false 时，遇到第一条脏数据就返回 error（严格模式）。
	SkipBadRows bool
}

func defaultCSVParseOptions() CSVParseOptions {
	return CSVParseOptions{
		TrimSpace:   true,
		AllowBOM:    true,
		SkipBadRows: false, // 默认保持你当前行为：严格
	}
}

// ParseCSVToStudents 读取 CSV 文件并解析为 Student 切片（使用默认解析选项）。
func ParseCSVToStudents(filename string) ([]model.Student, error) {
	return ParseCSVToStudentsWithOptions(filename, defaultCSVParseOptions())
}

// ParseCSVToStudentsWithOptions 读取 CSV 文件并解析为 Student 切片（可配置解析策略）。
//
// 该解析器按“中文表头”映射字段，因此列顺序可以变化，但必须包含以下列：
//   - 姓名
//   - 年龄
//   - 城市
//   - 得分
//
// 行级策略：
//   - 空行/全空字段行：跳过
//   - 姓名缺失：严格模式返回 error；宽松模式（SkipBadRows=true）跳过该行
//   - 年龄/得分解析失败：严格模式返回 error；宽松模式跳过该行
//
// 兼容性：
//   - AllowBOM=true 时会自动剥离 UTF-8 BOM（常见于 Excel 导出的 CSV 表头）
//   - TrimSpace=true 时会对表头与单元格做 TrimSpace
//
// 返回值：
//   - 成功时返回解析得到的 students（顺序与文件行顺序一致，不含表头行）
//   - 失败时返回 error，并尽可能包含行号与原始值，便于定位数据问题
func ParseCSVToStudentsWithOptions(filename string, opts CSVParseOptions) ([]model.Student, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV失败: %s, 错误: %v", filename, err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("CSV文件为空: %s", filename)
	}

	hdr := records[0]
	idx := headerIndex(hdr, opts.TrimSpace, opts.AllowBOM)

	required := []string{"姓名", "年龄", "城市", "得分"}
	for _, col := range required {
		if _, ok := idx[col]; !ok {
			return nil, fmt.Errorf("缺少必需列: %s", col)
		}
	}

	// 逐行解析 (跳过表头)
	student := make([]model.Student, 0, len(records)-1)
	for line := 2; line <= len(records); line++ {
		row := records[line-1]

		//空行/全空字段：跳过
		nonEmpty := false
		for _, cell := range row {
			if opts.TrimSpace {
				cell = strings.TrimSpace(cell)
			}
			if cell != "" {
				nonEmpty = true
				break
			}
		}
		if !nonEmpty {
			log.Printf("跳过空行: 第 %d 行\n", line)
			continue
		}

		name, ok := getCell(row, idx, "姓名", opts.TrimSpace)
		if !ok || name == " " {
			if opts.SkipBadRows {
				continue
			}
			return nil, fmt.Errorf("缺少姓名(第%d行)", line)
		}
		ageStr, _ := getCell(row, idx, "年龄", opts.TrimSpace)
		city, _ := getCell(row, idx, "城市", opts.TrimSpace)
		scoreStr, _ := getCell(row, idx, "得分", opts.TrimSpace)

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			if opts.SkipBadRows {
				continue
			}
			return nil, fmt.Errorf("解析年龄失败(第%d行, 值=%q): %w", line, ageStr, err)
		}

		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			if opts.SkipBadRows {
				continue
			}
			return nil, fmt.Errorf("解析得分失败(第%d行, 值=%q): %w", line, scoreStr, err)
		}

		student = append(student, model.Student{
			Name:  name,
			Age:   age,
			City:  city,
			Score: score,
		})
	}

	return student, nil

}

// headerIndex 根据表头行构建 “列名 -> 下标” 的索引。
// 会对列名做 TrimSpace，并忽略空列名。
func headerIndex(headers []string, trimSpace bool, allowBOM bool) map[string]int {
	idx := make(map[string]int, len(headers))
	for i, h := range headers {
		if allowBOM && i == 0 {
			h = stripBOM(h)
		}

		if trimSpace {
			h = strings.TrimSpace(h)
		}

		if h == " " {
			continue
		}
		// 只记录第一个出现的表头位置，后续重复的表头会被忽略
		if _, exists := idx[h]; !exists {
			idx[h] = i
		}
	}
	return idx
}

// getCell 按列名安全取值；若列不存在或下标越界，返回 ok=false。
func getCell(row []string, idx map[string]int, col string, trimSpace bool) (val string, ok bool) {
	i, exists := idx[col]
	if !exists || i < 0 || i >= len(row) {
		return "", false
	}
	v := row[i]
	if trimSpace {
		v = strings.TrimSpace(v)
	}
	return v, true
}

// stripBOM 去掉可能出现在 UTF-8 文本开头的 BOM 字符（\ufeff）。
func stripBOM(s string) string {
	return strings.TrimPrefix(s, "\ufeff")
}
