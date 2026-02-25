# hanzi-data-kit

[![CI](https://github.com/xianyudd/hanzi-data-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/xianyudd/hanzi-data-kit/actions/workflows/ci.yml)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/xianyudd/hanzi-data-kit/blob/main/LICENSE)

一个用于**生成中文学生样例数据**并进行 **CSV 写入/解析** 的 Go 项目，适合做数据处理练习、导入导出链路验证和单元测试样本构造。

## 功能特性

- 通过固定随机种子生成可复现的学生数据
- 支持按配置生成姓名、年龄、城市、分数
- 支持将大量数据流式写入 CSV（降低内存压力）
- 支持按中文表头解析 CSV（列顺序可变）
- 支持 BOM 处理、空白裁剪、坏行跳过/严格模式

## 项目结构

```text
.
├── cmd/
│   ├── gen_students/     # 生成学生 CSV 的 CLI
│   └── parse_students/   # 解析学生 CSV 的 CLI
├── generator/            # 数据生成器
├── model/                # 领域模型与 CSV 映射
├── parser/               # CSV 解析与写入
├── main.go               # 端到端示例（先生成再解析）
└── data/                 # 示例数据目录（CSV 默认输出到这里）
```

## 环境要求

- Go `1.24.0`（见 `go.mod`）

## 快速开始

### 1) 生成 CSV 数据

```bash
go run ./cmd/gen_students -n 1000 -seed 42 -out data/students.csv
```

常用参数：

- `-n` 生成条数（默认 `1000`）
- `-seed` 随机种子（默认 `42`，用于复现）
- `-out` 输出路径（默认 `data/students.csv`）
- `-age-min` / `-age-max` 年龄范围（闭区间）
- `-score-min` / `-score-max` 分数范围（闭区间）

### 2) 解析 CSV 数据

```bash
go run ./cmd/parse_students -in data/students.csv -print 5
```

常用参数：

- `-in` 输入 CSV 路径（默认 `data/students.csv`）
- `-print` 打印前 N 条（`0` 表示不打印）
- `-skip-bad-rows` 坏行是否跳过（默认 `true`）
- `-trim-space` 是否去除字段两侧空白（默认 `true`）
- `-allow-bom` 是否允许 UTF-8 BOM（默认 `true`）

### 3) 运行端到端示例

```bash
go run .
```

该示例会：

1. 生成 1000 条学生数据并写入 `data/students.csv`
2. 读取并解析 CSV
3. 打印前 5 条结果

## CSV 约定

- 默认中文表头：`姓名,年龄,城市,得分`
- 解析时按表头名映射字段，因此列顺序可以变化
- `得分` 目前按 1 位小数输出（如 `62.0` / `62.5`）

## 在代码中使用

```go
gen := generator.NewStudentGenerator(generator.StudentGenConfig{
    Seed:     42,
    AgeMin:   18,
    AgeMax:   30,
    ScoreMin: 60,
    ScoreMax: 100,
})

students := make([]model.Student, 0, 100)
for i := 0; i < 100; i++ {
    students = append(students, gen.Next())
}

headers := model.StudentHeadersCN()
_ = parser.WriteLargeCSV("data/students.csv", headers, len(students), func(i int) []string {
    return model.StudentToRowCN(students[i-1])
})
```

## 测试

```bash
go test ./...
```

## 说明

- `.gitignore` 已忽略 `*.csv` / `*.xlsx`，本地生成的数据文件默认不会提交到仓库。
