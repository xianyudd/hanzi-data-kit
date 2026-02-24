package model

//注意：结构体的名字和字段的名字首字母必须大写，否则在其他包中无法访问
type Student struct {
	Name string
	Age int
	City string
	Score float64
}