package parser

import (
	"encoding/csv"
	"fmt"
	"os"
)

// WriteLargeCSV 通用流式写入工具
// 参数说明:
// filename: 文件名
// headers: CSV 表头
// totalRows: 需要写入的总行数
// rowGenerator: 一个函数，调用方通过它来告诉底层每一行该写什么数据
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
