package generator

import (
	"github.com/xianyudd/hanzi-data-kit/model"
	"math"
	"math/rand"
)

// StudentGenConfig 定义学生数据生成策略。
// 约定：所有随机数由 Seed 驱动，以便测试场景可复现。
type StudentGenConfig struct {
	// Seed 为随机种子；相同配置+相同 Seed 将生成相同序列的数据。
	Seed int64

	// AgeMin/AgeMax 为闭区间 [min, max]。
	AgeMin int
	AgeMax int

	// ScoreMin/ScoreMax 为闭区间 [min, max]。
	ScoreMin float64
	ScoreMax float64

	// Cities 为候选城市列表；为空时会使用内置默认列表。
	Cities []string

	// Surnames 为候选姓氏列表；为空时会使用内置默认列表。
	Surnames []string

	// GivenNames1 为单字名候选；为空时会使用内置默认列表。
	GivenNames1 []string

	// GivenNames2 用于拼双字名（从中抽两个字拼接）；为空时默认复用 GivenNames1。
	GivenNames2 []string

	// TwoCharNameProb 控制生成双字名的概率（0~1）。默认 0.3。
	TwoCharNameProb float64

	// ScoreStep 为得分步长；0.5 => 只会生成 x.0 或 x.5。
	ScoreStep float64
}

// StudentGenerator 是一个基于 StudentGenConfig 的学生数据生成器。
type StudentGenerator struct {
	rng *rand.Rand
	cfg StudentGenConfig
}

// NewStudentGenerator 构造一个学生数据生成器。
// 注意：该函数会对 cfg 做“默认值补全”，以避免调用方遗漏配置导致 panic。
func NewStudentGenerator(cfg StudentGenConfig) *StudentGenerator {
	applyDefaults(&cfg)

	return &StudentGenerator{
		rng: rand.New(rand.NewSource(cfg.Seed)),
		cfg: cfg,
	}
}

// Next 生成一条新的学生记录。
// 生成的字段满足 cfg 中指定的范围与候选列表约束。
func (g *StudentGenerator) Next() model.Student {
	score := randFloat(g.rng, g.cfg.ScoreMin, g.cfg.ScoreMax)
	score = quantizeStep(score, 0.5) // 关键：步长 0.5 => 只会有 .0 或 .5

	return model.Student{
		Name:  g.genName(),
		Age:   randInt(g.rng, g.cfg.AgeMin, g.cfg.AgeMax),
		City:  pickOne(g.rng, g.cfg.Cities),
		Score: score,
	}
}

// genName 生成中文姓名：姓 +（单字名|双字名）。
// 双字名通过从 GivenNames2 中抽取两个字拼接而成。
func (g *StudentGenerator) genName() string {
	surname := pickOne(g.rng, g.cfg.Surnames)
	if g.rng.Float64() < g.cfg.TwoCharNameProb {
		a := pickOne(g.rng, g.cfg.GivenNames2)
		b := pickOne(g.rng, g.cfg.GivenNames2)
		return surname + a + b
	}
	return surname + pickOne(g.rng, g.cfg.GivenNames1)
}

func applyDefaults(cfg *StudentGenConfig) {
	// 默认年龄范围：更贴近“学生”语义；调用方可自行覆盖。
	if cfg.AgeMin == 0 && cfg.AgeMax == 0 {
		cfg.AgeMin, cfg.AgeMax = 18, 30
	}

	// 默认得分步长为0.5，生成更常见的整数或半整数分数。
	if cfg.ScoreStep == 0 {
		cfg.ScoreStep = 0.5
	}
	// 默认得分范围：0~100。
	if cfg.ScoreMin == 0 && cfg.ScoreMax == 0 {
		cfg.ScoreMin, cfg.ScoreMax = 0, 100
	}

	// 默认双字名概率。
	if cfg.TwoCharNameProb <= 0 || cfg.TwoCharNameProb >= 1 {
		cfg.TwoCharNameProb = 0.30
	}

	if len(cfg.Cities) == 0 {
		cfg.Cities = []string{"北京", "上海", "广州", "深圳", "成都", "杭州", "南京", "武汉", "西安", "重庆"}
	}
	if len(cfg.Surnames) == 0 {
		cfg.Surnames = []string{"赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋", "沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张"}
	}
	if len(cfg.GivenNames1) == 0 {
		cfg.GivenNames1 = []string{"伟", "芳", "娜", "敏", "静", "强", "磊", "军", "洋", "勇", "艳", "杰", "涛", "明", "超", "霞", "平", "刚", "桂英", "欣"}
	}
	if len(cfg.GivenNames2) == 0 {
		cfg.GivenNames2 = cfg.GivenNames1
	}
}

// pickOne 从非空切片中随机取一个元素。
// 该函数假设 slice 非空（由 applyDefaults 保证）。
func pickOne(rng *rand.Rand, xs []string) string {
	return xs[rng.Intn(len(xs))]
}

// randInt 返回闭区间 [min, max] 的随机整数。
// 若 min/max 颠倒则自动交换，避免调用方配置错误导致异常。
func randInt(rng *rand.Rand, min, max int) int {
	if max < min {
		min, max = max, min
	}
	if max == min {
		return min
	}
	return rng.Intn(max-min+1) + min
}

// randFloat 返回闭区间 [min, max] 内的随机浮点。
// 若 min/max 颠倒则自动交换。
func randFloat(rng *rand.Rand, min, max float64) float64 {
	if max < min {
		min, max = max, min
	}
	if max == min {
		return min
	}
	return rng.Float64()*(max-min) + min
}

// quantizeStep 将 x 量化到 step 的整数倍（四舍五入）。
// step=0.5 时，结果的小数部分只可能为 .0 或 .5。
func quantizeStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	return math.Round(x/step) * step
}
