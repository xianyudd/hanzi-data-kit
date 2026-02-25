package generator_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/xianyudd/hanzi-data-kit/generator"
	"github.com/xianyudd/hanzi-data-kit/model"
)

func TestStudentGenerator_ReproducibleWithSameSeed(t *testing.T) {
	cfg := generator.StudentGenConfig{
		Seed:     42,
		AgeMin:   18,
		AgeMax:   30,
		ScoreMin: 60,
		ScoreMax: 100,
		// 你如果是通过 cfg.ScoreStep 控制步长，这里也可以显式写：ScoreStep: 0.5
	}

	g1 := generator.NewStudentGenerator(cfg)
	g2 := generator.NewStudentGenerator(cfg)

	const n = 50
	for i := 0; i < n; i++ {
		a := g1.Next()
		b := g2.Next()

		// model.Student 是可比较的 struct（字段都是可比较类型），可直接 != 对比
		if a != b {
			t.Fatalf("not reproducible at i=%d:\n  a=%#v\n  b=%#v", i, a, b)
		}
	}
}

func TestStudentGenerator_ScoreIsStepOfHalf(t *testing.T) {
	cfg := generator.StudentGenConfig{
		Seed:     7,
		ScoreMin: 60,
		ScoreMax: 100,
		// 若你的实现允许配置步长，建议写：ScoreStep: 0.5
	}

	g := generator.NewStudentGenerator(cfg)

	const n = 1000
	for i := 0; i < n; i++ {
		stu := g.Next()

		// 核心断言：Score 必须是 0.5 的整数倍
		// 用 Score*2 变成整数，考虑浮点误差，用 Round-like 的方式判断
		v := stu.Score * 2
		iv := int64(v + 0.5) // 适用于非负值；你的分数范围本身是正数
		if diff := v - float64(iv); diff < -1e-9 || diff > 1e-9 {
			t.Fatalf("score not multiple of 0.5 at i=%d: score=%v", i, stu.Score)
		}

		// 同时确保在范围内（避免量化后越界）
		if stu.Score < cfg.ScoreMin-1e-9 || stu.Score > cfg.ScoreMax+1e-9 {
			t.Fatalf("score out of range at i=%d: score=%v", i, stu.Score)
		}
	}
}

func TestStudentToRowCN_ScoreOneDecimalAndHalfStep(t *testing.T) {
	students := []model.Student{
		{Name: "张三", Age: 22, City: "北京", Score: 62.0},
		{Name: "李四", Age: 25, City: "上海", Score: 62.5},
		{Name: "王五", Age: 28, City: "广州", Score: 75.0},
		{Name: "赵六", Age: 21, City: "深圳", Score: 75.5},
	}

	for i, stu := range students {
		row := model.StudentToRowCN(stu)
		if len(row) != 4 {
			t.Fatalf("row length mismatch at i=%d: got=%d", i, len(row))
		}

		scoreStr := row[3]

		// 断言：必须是 1 位小数（形如 "62.0"）
		if strings.Count(scoreStr, ".") != 1 {
			t.Fatalf("score format must have one dot at i=%d: %q", i, scoreStr)
		}
		parts := strings.Split(scoreStr, ".")
		if len(parts) != 2 || len(parts[1]) != 1 {
			t.Fatalf("score must have exactly 1 decimal digit at i=%d: %q", i, scoreStr)
		}
		if parts[1] != "0" && parts[1] != "5" {
			t.Fatalf("score decimal must be 0 or 5 at i=%d: %q", i, scoreStr)
		}

		// 进一步：能解析回 float，且仍是 0.5 步长
		f, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			t.Fatalf("score string not parseable at i=%d: %q err=%v", i, scoreStr, err)
		}
		v := f * 2
		iv := int64(v + 0.5)
		if diff := v - float64(iv); diff < -1e-9 || diff > 1e-9 {
			t.Fatalf("parsed score not multiple of 0.5 at i=%d: %v", i, f)
		}
	}
}
