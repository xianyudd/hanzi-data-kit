package parser

import (
	"encoding/csv"
	"fmt"
	"os"
)

// WriteLargeCSV 以流式方式将大量行写入 CSV 文件 filename。
//
// 行写入顺序：可选表头 headers（若非空）→ 共 totalRows 行数据。
// rowGenerator 会被调用 totalRows 次，参数 rowIndex 为 1-based（范围 [1, totalRows]），
// 返回值为该行的列数据（[]string）。
//
// 函数会定期 Flush 缓冲区；若发生底层 I/O 错误（如磁盘写满、权限不足），会返回 error。
//
// 参数说明:
//   - filename: 文件名
//   - headers: CSV 表头
//   - totalRows: 需要写入的总行数
//   - rowGenerator: 行生成器（1-based 行号）
func WriteLargeCSV(filename string, headers []string, totalRows int, rowGenerator func(rowNum int) []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败：%w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("写入表头失败：%w", err)
		}
	}
	flushInterval := 100_000

	for i := 1; i <= totalRows; i++ {
		row := rowGenerator(i)

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入第 %d 行失败: %w", i, err)
		}
		if i%flushInterval == 0 {
			writer.Flush()
			if err := writer.Error(); err != nil {
				return fmt.Errorf("刷新缓冲区到磁盘失败: %w", err)
			}
		}
		fmt.Printf("已写入 %d/%d 行\r", i, totalRows)
	}
	return nil
}
