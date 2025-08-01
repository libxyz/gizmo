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
)

func main() {
    // 创建解析器
    parser := NewParser()
    
    // JSON输入
    jsonInput := `{"name":"张三","age":25,"items":[1,2,3]}`
    
    // 累积特定路径的字符
    acc := strings.Builder{}
    
    for _, r := range jsonInput {
        event := parser.Push(r)
        if event != nil && event.Path == "$.name" && 
           (event.Type == EventString || event.Type == EventStringEscape) {
            acc.WriteRune(event.Char)
        }
    }
    
    fmt.Println("提取的name值:", acc.String()) // 输出: 张三
}
```

## 事件类型

解析器为每个字符生成以下类型的事件：

| 事件类型 | 描述 | 示例 |
|---------|------|------|
| `EventString` | 字符串内容字符 | `"hello"` 中的 `h` |
| `EventStringEscape` | 转义字符 | `"\n"` 中的 `\` |
| `EventNumber` | 数字字符 | `42.5` 中的 `4`、`2`、`.`、`5` |
| `EventBoolean` | 布尔值字符 | `true` 中的 `t`、`r`、`u`、`e` |
| `EventNull` | null值字符 | `null` 中的每个字符 |
| `EventObjectStart` | 对象开始 | `{` |
| `EventObjectEnd` | 对象结束 | `}` |
| `EventArrayStart` | 数组开始 | `[` |
| `EventArrayEnd` | 数组结束 | `]` |
| `EventKey` | 键名字符 | `{"key":1}` 中的 `k`、`e`、`y` |
| `EventKeyEscape` | 键名转义字符 | `{"k\ey":1}` 中的 `\` |
| `EventComma` | 逗号分隔符 | `,` |
| `EventColon` | 冒号分隔符 | `:` |
| `EventQuote` | 引号 | `"` |
| `EventWhitespace` | 空白字符 | 空格、制表符、换行符等 |

## JSON路径格式

解析器使用JSON Pointer格式来标识当前处理的位置：

- `$` - 根节点
- `$.key` - 对象中的字段
- `$[0]` - 数组中的索引
- `$.users[0].name` - 嵌套结构

## 示例用法

### 基本对象解析

```go
json := `{"name":"张三","age":30,"active":true}`
parser := NewParser()

for _, r := range json {
    if event := parser.Push(r); event != nil {
        fmt.Printf("字符: %c, 类型: %d, 路径: %s\n", 
                  event.Char, event.Type, event.Path)
    }
}
```

### 数组处理

```go
json := `[1,2,3,{"name":"测试"}]`
parser := NewParser()

for _, r := range json {
    if event := parser.Push(r); event != nil {
        if event.Path == "$[3].name" && event.Type == EventString {
            fmt.Printf("找到名称: %c\n", event.Char)
        }
    }
}
```

### 复杂嵌套结构

```go
json := `{"users":[{"id":1,"profile":{"name":"张三"}}]}`
parser := NewParser()

nameBuilder := strings.Builder{}
for _, r := range json {
    if event := parser.Push(r); event != nil {
        if event.Path == "$.users[0].profile.name" && 
           (event.Type == EventString || event.Type == EventStringEscape) {
            nameBuilder.WriteRune(event.Char)
        }
    }
}
fmt.Println("用户名:", nameBuilder.String()) // 输出: 张三
