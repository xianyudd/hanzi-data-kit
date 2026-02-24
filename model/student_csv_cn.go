package model

import "strconv"

// StudentHeadersCN 返回中文的CSV表头 （与 Student 结构体字段顺序一致）
func StudentHeadersCN() []string {
	return []string{"姓名", "年龄", "城市", "得分"}
}

func StudentToRowCN(stu Student) []string {
	return []string{
		stu.Name,
		strconv.Itoa(stu.Age),
		stu.City,
		strconv.FormatFloat(stu.Score, 'f', 0, 64),
	}
}
