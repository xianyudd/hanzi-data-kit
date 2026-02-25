package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"go-data-learning/model"
	"go-data-learning/parser"
)

// writeTempFile 在临时目录写入文件并返回路径。
// 使用 t.Helper() 让失败栈信息指向调用方，更易定位。
func writeTempFile(t *testing.T, name string, content []byte) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, name)

	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file failed: %v", err)
	}
	return path
}

func TestParseCSVToStudents_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		opts    *parser.CSVParseOptions
		want    []model.Student
		wantErr bool
	}{
		{
			name:    "header_reordered_mapping_ok",
			content: []byte("得分,城市,姓名,年龄\n92.5,上海,李四,25\n74.0,北京,张三,22\n"),
			opts:    nil, // 使用 ParseCSVToStudents 的默认行为
			want: []model.Student{
				{Name: "李四", Age: 25, City: "上海", Score: 92.5},
				{Name: "张三", Age: 22, City: "北京", Score: 74.0},
			},
		},
		{
			name: "bom_header_ok",
			// UTF-8 BOM：EF BB BF，常见于 Excel 导出 CSV。
			content: append([]byte{0xEF, 0xBB, 0xBF},
				[]byte("姓名,年龄,城市,得分\n张三,22,北京,95.0\n")...,
			),
			opts: &parser.CSVParseOptions{
				TrimSpace:   true,
				AllowBOM:    true,
				SkipBadRows: false,
			},
			want: []model.Student{
				{Name: "张三", Age: 22, City: "北京", Score: 95.0},
			},
		},
		{
			name: "skip_bad_rows_ok",
			// 第二行年龄字段非法，SkipBadRows=true 时应跳过该行
			content: []byte(
				"姓名,年龄,城市,得分\n" +
					"张三,22,北京,95.0\n" +
					"李四,notint,上海,88.0\n" +
					"王五,28,广州,92.5\n",
			),
			opts: &parser.CSVParseOptions{
				TrimSpace:   true,
				AllowBOM:    true,
				SkipBadRows: true,
			},
			want: []model.Student{
				{Name: "张三", Age: 22, City: "北京", Score: 95.0},
				{Name: "王五", Age: 28, City: "广州", Score: 92.5},
			},
		},
		{
			name:    "strict_mode_bad_row_returns_error",
			content: []byte("姓名,年龄,城市,得分\n张三,22,北京,95.0\n李四,notint,上海,88.0\n"),
			opts: &parser.CSVParseOptions{
				TrimSpace:   true,
				AllowBOM:    true,
				SkipBadRows: false, // 严格模式：遇到坏行应报错
			},
			wantErr: true,
		},
		{
			name:    "missing_required_header_returns_error",
			content: []byte("姓名,年龄,城市\n张三,22,北京\n"), // 缺“得分”
			opts: &parser.CSVParseOptions{
				TrimSpace:   true,
				AllowBOM:    true,
				SkipBadRows: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // 经典写法：避免闭包捕获循环变量（对旧版本 Go 更安全）
		t.Run(tt.name, func(t *testing.T) {
			path := writeTempFile(t, "in.csv", tt.content)

			var (
				got []model.Student
				err error
			)

			if tt.opts == nil {
				got, err = parser.ParseCSVToStudents(path)
			} else {
				got, err = parser.ParseCSVToStudentsWithOptions(path, *tt.opts)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d students, got %d", len(tt.want), len(got))
			}

			// 不引入第三方断言库的情况下，逐字段对比（可读且稳定）。
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("student[%d] mismatch:\n  got:  %#v\n  want: %#v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
