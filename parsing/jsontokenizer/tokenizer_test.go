package jsontokenizer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_Simple 测试基本的JSON解析功能
// 包括字符串、数字、转义字符等基础类型
func TestParser_Simple(t *testing.T) {
	parser := newInnerTokenizer()
	json := `{"a":"te\n\"st", "b":42}`

	events := []event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	expectedEvents := []event{
		{Type: TokenObjectStart, Path: "$", Char: '{'},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenKey, Path: "$", Char: 'a'},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenColon, Path: "$.a", Char: ':'},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenString, Path: "$.a", Char: 't'},
		{Type: TokenString, Path: "$.a", Char: 'e'},
		{Type: TokenStringEscape, Path: "$.a", Char: '\\'},
		{Type: TokenString, Path: "$.a", Char: 'n'},
		{Type: TokenStringEscape, Path: "$.a", Char: '\\'},
		{Type: TokenString, Path: "$.a", Char: '"'},
		{Type: TokenString, Path: "$.a", Char: 's'},
		{Type: TokenString, Path: "$.a", Char: 't'},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenComma, Path: "$", Char: ','},
		{Type: TokenWhitespace, Path: "$", Char: ' '},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenKey, Path: "$", Char: 'b'},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenColon, Path: "$.b", Char: ':'},
		{Type: TokenNumber, Path: "$.b", Char: '4'},
		{Type: TokenNumber, Path: "$.b", Char: '2'},
		{Type: TokenObjectEnd, Path: "$", Char: '}'},
	}

	require.Len(t, expectedEvents, len(json), "Expected %d events, got %d", len(expectedEvents), len(events))

	accumulatedJSON := strings.Builder{}
	for _, event := range events {
		accumulatedJSON.WriteRune(event.Char)
	}
	require.Equal(t, json, accumulatedJSON.String(), "Accumulated JSON does not match original")

	assert.Len(t, expectedEvents, len(events), "Expected %d events, got %d", len(expectedEvents), len(events))
	assert.Equal(t, expectedEvents, events, "Events do not match expected")
}

// TestPaser_Complex 测试复杂的嵌套JSON结构
// 包括嵌套对象、数组、混合类型等
func TestPaser_Complex(t *testing.T) {
	parser := newInnerTokenizer()
	json := `{"a":{"b":[1,2,"3"],"c":true,"d":{"e":null}},"fake":-1.1}`

	events := []event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	expectedEvents := []event{
		{Type: TokenObjectStart, Path: "$", Char: '{'},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenKey, Path: "$", Char: 'a'},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenColon, Path: "$.a", Char: ':'},
		{Type: TokenObjectStart, Path: "$.a", Char: '{'},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenKey, Path: "$.a", Char: 'b'},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenColon, Path: "$.a.b", Char: ':'},
		{Type: TokenArrayStart, Path: "$.a.b", Char: '['},
		{Type: TokenNumber, Path: "$.a.b[0]", Char: '1'},
		{Type: TokenComma, Path: "$.a.b[1]", Char: ','},
		{Type: TokenNumber, Path: "$.a.b[1]", Char: '2'},
		{Type: TokenComma, Path: "$.a.b[2]", Char: ','},
		{Type: TokenQuote, Path: "$.a.b[2]", Char: '"'},
		{Type: TokenString, Path: "$.a.b[2]", Char: '3'},
		{Type: TokenQuote, Path: "$.a.b[2]", Char: '"'},
		{Type: TokenArrayEnd, Path: "$.a.b", Char: ']'},
		{Type: TokenComma, Path: "$.a", Char: ','},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenKey, Path: "$.a", Char: 'c'},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenColon, Path: "$.a.c", Char: ':'},
		{Type: TokenBoolean, Path: "$.a.c", Char: 't'},
		{Type: TokenBoolean, Path: "$.a.c", Char: 'r'},
		{Type: TokenBoolean, Path: "$.a.c", Char: 'u'},
		{Type: TokenBoolean, Path: "$.a.c", Char: 'e'},
		{Type: TokenComma, Path: "$.a", Char: ','},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenKey, Path: "$.a", Char: 'd'},
		{Type: TokenQuote, Path: "$.a", Char: '"'},
		{Type: TokenColon, Path: "$.a.d", Char: ':'},
		{Type: TokenObjectStart, Path: "$.a.d", Char: '{'},
		{Type: TokenQuote, Path: "$.a.d", Char: '"'},
		{Type: TokenKey, Path: "$.a.d", Char: 'e'},
		{Type: TokenQuote, Path: "$.a.d", Char: '"'},
		{Type: TokenColon, Path: "$.a.d.e", Char: ':'},
		{Type: TokenNull, Path: "$.a.d.e", Char: 'n'},
		{Type: TokenNull, Path: "$.a.d.e", Char: 'u'},
		{Type: TokenNull, Path: "$.a.d.e", Char: 'l'},
		{Type: TokenNull, Path: "$.a.d.e", Char: 'l'},
		{Type: TokenObjectEnd, Path: "$.a.d", Char: '}'},
		{Type: TokenObjectEnd, Path: "$.a", Char: '}'},
		{Type: TokenComma, Path: "$", Char: ','},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenKey, Path: "$", Char: 'f'},
		{Type: TokenKey, Path: "$", Char: 'a'},
		{Type: TokenKey, Path: "$", Char: 'k'},
		{Type: TokenKey, Path: "$", Char: 'e'},
		{Type: TokenQuote, Path: "$", Char: '"'},
		{Type: TokenColon, Path: "$.fake", Char: ':'},
		{Type: TokenNumber, Path: "$.fake", Char: '-'},
		{Type: TokenNumber, Path: "$.fake", Char: '1'},
		{Type: TokenNumber, Path: "$.fake", Char: '.'},
		{Type: TokenNumber, Path: "$.fake", Char: '1'},
		{Type: TokenObjectEnd, Path: "$", Char: '}'},
	}

	require.Len(t, expectedEvents, len(json), "Expected %d events, got %d", len(expectedEvents), len(events))

	accumulatedJSON := strings.Builder{}
	for _, event := range events {
		accumulatedJSON.WriteRune(event.Char)
	}
	require.Equal(t, json, accumulatedJSON.String(), "Accumulated JSON does not match original")

	assert.Len(t, expectedEvents, len(events), "Expected %d events, got %d", len(expectedEvents), len(events))
	assert.Equal(t, expectedEvents, events, "Events do not match expected")
}

// TestParser_EscapedWhitespace 测试转义空白字符的处理
// 包括制表符、换行符、回车符等
func TestParser_EscapedWhitespace(t *testing.T) {
	parser := newInnerTokenizer()
	// Test various escaped whitespace characters
	json := `{"tab":"te\t","newline":"li\nne","return":"car\r","backspace":"bs\b","formfeed":"ff\f"}`

	events := []event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	accumulatedJSON := strings.Builder{}
	for _, event := range events {
		accumulatedJSON.WriteRune(event.Char)
	}
	require.Equal(t, json, accumulatedJSON.String(), "Accumulated JSON does not match original with escaped whitespace")
}

// TestParser_MixedEscapedWhitespace 测试混合转义空白字符的处理
// 包括嵌套结构中的转义字符
func TestParser_MixedEscapedWhitespace(t *testing.T) {
	parser := newInnerTokenizer()
	// Test complex string with multiple escaped whitespace characters
	json := `{"mixed":"\t\n\r\b\f","nested":{"inner":"text\twith\nnewlines"}}`

	events := []event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	accumulatedJSON := strings.Builder{}
	for _, event := range events {
		accumulatedJSON.WriteRune(event.Char)
	}
	require.Equal(t, json, accumulatedJSON.String(), "Accumulated JSON does not match original with mixed escaped whitespace")
}

// TestParser_StringWithUnicodeEscapes 测试Unicode转义序列的处理
// 包括\uXXXX格式的Unicode字符
func TestParser_StringWithUnicodeEscapes(t *testing.T) {
	parser := newInnerTokenizer()
	// Test string with Unicode escape sequences (common in JSON)
	json := `{"unicode":"\u0041\u0042\u0043","with_spaces" : "text\t\nmore"}`

	events := []event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	accumulatedJSON := strings.Builder{}
	for _, event := range events {
		accumulatedJSON.WriteRune(event.Char)
	}
	require.Equal(t, json, accumulatedJSON.String(), "Accumulated JSON does not match original with Unicode escapes")
}

// TestParser_KeyEscapes 测试键名中转义字符的处理
// 包括键名和值中的转义字符
func TestParser_KeyEscapes(t *testing.T) {
	parser := newInnerTokenizer()
	// Test string with Unicode escape sequences (common in JSON)
	json := `{"key_with_e\n\"":"value_with_escape\u0041"}`

	events := []event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	accumulatedJSON := strings.Builder{}
	for _, event := range events {
		accumulatedJSON.WriteRune(event.Char)
	}
	require.Equal(t, json, accumulatedJSON.String(), "Accumulated JSON does not match original with Unicode escapes")
}

func TestParser_AutoEscape(t *testing.T) {
	parser := NewTokenizer()
	parser.AutoEscape()
	json := `{"a":"te\n\"\u0028st", "b":42}`

	events := []Token{}
	for _, r := range json {
		event := parser.Push(r)
		t.Logf("Pushed rune: %c, Event: %+v", r, event)
		if event != nil {
			events = append(events, *event)
		}
	}

	expectedEvents := []Token{
		{Type: TokenObjectStart, Path: "$", Val: "{"},
		{Type: TokenQuote, Path: "$", Val: "\""},
		{Type: TokenKey, Path: "$", Val: "a"},
		{Type: TokenQuote, Path: "$", Val: "\""},
		{Type: TokenColon, Path: "$.a", Val: ":"},
		{Type: TokenQuote, Path: "$.a", Val: "\""},
		{Type: TokenString, Path: "$.a", Val: "t"},
		{Type: TokenString, Path: "$.a", Val: "e"},
		{Type: TokenString, Path: "$.a", Val: "\n"},
		{Type: TokenString, Path: "$.a", Val: "\""},
		{Type: TokenString, Path: "$.a", Val: "("},
		{Type: TokenString, Path: "$.a", Val: "s"},
		{Type: TokenString, Path: "$.a", Val: "t"},
		{Type: TokenQuote, Path: "$.a", Val: "\""},
		{Type: TokenComma, Path: "$", Val: ","},
		{Type: TokenWhitespace, Path: "$", Val: " "},
		{Type: TokenQuote, Path: "$", Val: "\""},
		{Type: TokenKey, Path: "$", Val: "b"},
		{Type: TokenQuote, Path: "$", Val: "\""},
		{Type: TokenColon, Path: "$.b", Val: ":"},
		{Type: TokenNumber, Path: "$.b", Val: "4"},
		{Type: TokenNumber, Path: "$.b", Val: "2"},
		{Type: TokenObjectEnd, Path: "$", Val: "}"},
	}

	for i, event := range events {
		assert.Equal(t, expectedEvents[i], event, "Event at index %d does not match expected", i)
	}
}

func TestNestedCatch(t *testing.T) {
	json := `{"users":[{"id":1,"profile":{"name":"张三"}},{"id":2,"profile":{"name":"李四"}}]}`
	z := NewTokenizer()

	nameBuilder := strings.Builder{}
	for _, r := range json {
		if tk := z.Push(r); tk != nil {
			if tk.Path == "$.users[1].profile.name" &&
				(tk.Type == TokenString || tk.Type == TokenStringEscape) {
				nameBuilder.WriteString(tk.Val)
			}
		}
	}
	assert.Equal(t, "李四", nameBuilder.String())
}
