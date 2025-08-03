# JSON流式解析器

这是一个用Go语言编写的JSON流式解析器，采用状态机模式逐个字符解析JSON数据。该解析器可以在内存占用极小的情况下处理大型JSON文件，并提供详细的解析事件。

## 特性

- 🚀 **流式处理**: 逐个字符解析，内存占用极低
- 📊 **详细事件**: 为每个字符提供解析事件和JSON路径
- 🔍 **完整支持**: 支持所有JSON数据类型（对象、数组、字符串、数字、布尔值、null）
- 🛡️ **转义处理**: 正确处理字符串中的转义字符
- 🎯 **路径跟踪**: 实时跟踪当前JSON路径（如 `$.users[0].name`）

## 快速开始

```go
package main

import (
    "fmt"
    "strings"
    "github.com/Crescent617/x/parsing/jsontokenizer"
)

func main() {
    // 创建解析器
    t := jsontokenizer.NewTokenizer()
    
    // JSON输入
    jsonInput := `{"name":"张三","age":25,"items":[1,2,3]}`
    
    // 累积特定路径的字符
    acc := strings.Builder{}
    
    for _, r := range jsonInput {
        tk := t.Push(r)
        if tk != nil && tk.Path == "$.name" && 
           (tk.Type == jsontokenizer.TokenString || tk.Type == jsontokenizer.TokenStringEscape) {
            acc.WriteString(tk.Val)
        }
    }
    
    fmt.Println("提取的name值:", acc.String()) // 输出: 张三
}
```

## Token类型

解析器为每个字符生成以下类型的Token：

| Token类型 | 描述 | 示例 |
|---------|------|------|
| `jsontokenizer.TokenString` | 字符串内容字符 | `"hello"` 中的 `h` |
| `jsontokenizer.TokenStringEscape` | 转义字符 | `"\n"` 中的 `\` |
| `jsontokenizer.TokenNumber` | 数字字符 | `42.5` 中的 `4`、`2`、`.`、`5` |
| `jsontokenizer.TokenBoolean` | 布尔值字符 | `true` 中的 `t`、`r`、`u`、`e` |
| `jsontokenizer.TokenNull` | null值字符 | `null` 中的每个字符 |
| `jsontokenizer.TokenObjectStart` | 对象开始 | `{` |
| `jsontokenizer.TokenObjectEnd` | 对象结束 | `}` |
| `jsontokenizer.TokenArrayStart` | 数组开始 | `[` |
| `jsontokenizer.TokenArrayEnd` | 数组结束 | `]` |
| `jsontokenizer.TokenKey` | 键名字符 | `{"key":1}` 中的 `k`、`e`、`y` |
| `jsontokenizer.TokenKeyEscape` | 键名转义字符 | `{"k\ey":1}` 中的 `\` |
| `jsontokenizer.TokenComma` | 逗号分隔符 | `,` |
| `jsontokenizer.TokenColon` | 冒号分隔符 | `:` |
| `jsontokenizer.TokenQuote` | 引号 | `"` |
| `jsontokenizer.TokenWhitespace` | 空白字符 | 空格、制表符、换行符等 |

## JSON路径格式

解析器使用JSON Pointer格式来标识当前处理的位置：

- `$` - 根节点
- `$.key` - 对象中的字段
- `$[0]` - 数组中的索引
- `$.users[0].name` - 嵌套结构

## 示例用法

### 基本对象解析

```go
import "github.com/Crescent617/x/parsing/jsontokenizer"

json := `{"name":"张三","age":30,"active":true}`
t := jsontokenizer.NewTokenizer()

for _, r := range json {
    if tk := t.Push(r); tk != nil {
        fmt.Printf("字符: %s, 类型: %d, 路径: %s\n", 
                  tk.Val, tk.Type, tk.Path)
    }
}
```

### 数组处理

```go
import "github.com/Crescent617/x/parsing/jsontokenizer"

json := `[1,2,3,{"name":"测试"}]`
t := jsontokenizer.NewTokenizer()

for _, r := range json {
    if tk := t.Push(r); tk != nil {
        if tk.Path == "$[3].name" && tk.Type == jsontokenizer.TokenString {
            fmt.Printf("找到名称: %s\n", tk.Val)
        }
    }
}
```

### 复杂嵌套结构

```go
json := `{"users":[{"id":1,"profile":{"name":"张三"}}]}`
t := jsontokenizer.NewTokenizer()

nameBuilder := strings.Builder{}
for _, r := range json {
    if tk := t.Push(r); tk != nil {
        if tk.Path == "$.users[0].profile.name" && 
           (tk.Type == jsontokenizer.TokenString || tk.Type == jsontokenizer.TokenStringEscape) {
            nameBuilder.WriteString(tk.Val)
        }
    }
}
fmt.Println("用户名:", nameBuilder.String()) // 输出: 张三
