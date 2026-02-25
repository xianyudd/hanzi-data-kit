package parser

import (
	"encoding/csv"
	"fmt"
	"go-data-learning/model"
	"log"
	"os"
	"strconv"
	"strings"
)

// ParseCSVToStudents 读取 CSV 文件并解析为 Student 切片。
//
// 该解析器按“中文表头”映射字段，因此列顺序可以变化，但必须包含以下列：
//   - 姓名
//   - 年龄
//   - 城市
//   - 得分
//
// 若缺少必需列，将返回 error。
// 若某行列数不足或空行，将被跳过。
// 若年龄/得分无法解析为数字，将返回 error（定位到行号，便于排查）。
func ParseCSVToStudents(filename string) ([]model.Student, error) {
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
	idx := headerIndex(hdr)

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
			if strings.TrimSpace(cell) != " " {
				nonEmpty = true
				break
			}
		}
		if !nonEmpty {
			log.Printf("跳过空行: 第 %d 行\n", line)
			continue
		}

		name, ok := getCell(row, idx, "姓名")
		if !ok || name == " " {
			//姓名缺失通常是脏数据；先跳过（后续Step4.2可做"严格模式")
			continue
		}
		ageStr, _ := getCell(row, idx, "年龄")
		city, _ := getCell(row, idx, "城市")
		scoreStr, _ := getCell(row, idx, "得分")

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return nil, fmt.Errorf("解析年龄失败(第%d行, 值=%q): %w", line, ageStr, err)
		}

		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
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
func headerIndex(headers []string) map[string]int {
	idx := make(map[string]int, len(headers))
	for i, h := range headers {
		h = strings.TrimSpace(h)
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
func getCell(row []string, idx map[string]int, col string) (val string, ok bool) {
	i, exists := idx[col]
	if !exists || i < 0 || i >= len(row) {
		return "", false
	}
	return strings.TrimSpace(row[i]), true
}
