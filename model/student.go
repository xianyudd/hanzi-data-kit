package model

// Student 代表系统中的一个标准学生信息对象。
// 它用于在 CSV/Excel 解析器与业务逻辑层之间传递数据。
type Student struct {
	Name  string  // Name 学生的真实姓名
	Age   int     // Age 学生的年龄
	City  string  // City 所在城市
	Score float64 // Score 考试得分
}
