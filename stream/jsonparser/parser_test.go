package jsonparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_Simple 测试基本的JSON解析功能
// 包括字符串、数字、转义字符等基础类型
func TestParser_Simple(t *testing.T) {
	parser := NewParser()
	json := `{"a":"te\n\"st", "b":42}`

	events := []Event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	expectedEvents := []Event{
		{Type: EventObjectStart, Path: "$", Char: '{'},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventKey, Path: "$", Char: 'a'},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventColon, Path: "$.a", Char: ':'},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventString, Path: "$.a", Char: 't'},
		{Type: EventString, Path: "$.a", Char: 'e'},
		{Type: EventStringEscape, Path: "$.a", Char: '\\'},
		{Type: EventString, Path: "$.a", Char: 'n'},
		{Type: EventStringEscape, Path: "$.a", Char: '\\'},
		{Type: EventString, Path: "$.a", Char: '"'},
		{Type: EventString, Path: "$.a", Char: 's'},
		{Type: EventString, Path: "$.a", Char: 't'},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventComma, Path: "$", Char: ','},
		{Type: EventWhitespace, Path: "$", Char: ' '},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventKey, Path: "$", Char: 'b'},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventColon, Path: "$.b", Char: ':'},
		{Type: EventNumber, Path: "$.b", Char: '4'},
		{Type: EventNumber, Path: "$.b", Char: '2'},
		{Type: EventObjectEnd, Path: "$", Char: '}'},
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
	parser := NewParser()
	json := `{"a":{"b":[1,2,"3"],"c":true,"d":{"e":null}},"fake":-1.1}`

	events := []Event{}
	for _, r := range json {
		event := parser.Push(r)
		events = append(events, event)
	}

	expectedEvents := []Event{
		{Type: EventObjectStart, Path: "$", Char: '{'},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventKey, Path: "$", Char: 'a'},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventColon, Path: "$.a", Char: ':'},
		{Type: EventObjectStart, Path: "$.a", Char: '{'},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventKey, Path: "$.a", Char: 'b'},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventColon, Path: "$.a.b", Char: ':'},
		{Type: EventArrayStart, Path: "$.a.b", Char: '['},
		{Type: EventNumber, Path: "$.a.b[0]", Char: '1'},
		{Type: EventComma, Path: "$.a.b[1]", Char: ','},
		{Type: EventNumber, Path: "$.a.b[1]", Char: '2'},
		{Type: EventComma, Path: "$.a.b[2]", Char: ','},
		{Type: EventQuote, Path: "$.a.b[2]", Char: '"'},
		{Type: EventString, Path: "$.a.b[2]", Char: '3'},
		{Type: EventQuote, Path: "$.a.b[2]", Char: '"'},
		{Type: EventArrayEnd, Path: "$.a.b", Char: ']'},
		{Type: EventComma, Path: "$.a", Char: ','},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventKey, Path: "$.a", Char: 'c'},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventColon, Path: "$.a.c", Char: ':'},
		{Type: EventBoolean, Path: "$.a.c", Char: 't'},
		{Type: EventBoolean, Path: "$.a.c", Char: 'r'},
		{Type: EventBoolean, Path: "$.a.c", Char: 'u'},
		{Type: EventBoolean, Path: "$.a.c", Char: 'e'},
		{Type: EventComma, Path: "$.a", Char: ','},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventKey, Path: "$.a", Char: 'd'},
		{Type: EventQuote, Path: "$.a", Char: '"'},
		{Type: EventColon, Path: "$.a.d", Char: ':'},
		{Type: EventObjectStart, Path: "$.a.d", Char: '{'},
		{Type: EventQuote, Path: "$.a.d", Char: '"'},
		{Type: EventKey, Path: "$.a.d", Char: 'e'},
		{Type: EventQuote, Path: "$.a.d", Char: '"'},
		{Type: EventColon, Path: "$.a.d.e", Char: ':'},
		{Type: EventNull, Path: "$.a.d.e", Char: 'n'},
		{Type: EventNull, Path: "$.a.d.e", Char: 'u'},
		{Type: EventNull, Path: "$.a.d.e", Char: 'l'},
		{Type: EventNull, Path: "$.a.d.e", Char: 'l'},
		{Type: EventObjectEnd, Path: "$.a.d", Char: '}'},
		{Type: EventObjectEnd, Path: "$.a", Char: '}'},
		{Type: EventComma, Path: "$", Char: ','},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventKey, Path: "$", Char: 'f'},
		{Type: EventKey, Path: "$", Char: 'a'},
		{Type: EventKey, Path: "$", Char: 'k'},
		{Type: EventKey, Path: "$", Char: 'e'},
		{Type: EventQuote, Path: "$", Char: '"'},
		{Type: EventColon, Path: "$.fake", Char: ':'},
		{Type: EventNumber, Path: "$.fake", Char: '-'},
		{Type: EventNumber, Path: "$.fake", Char: '1'},
		{Type: EventNumber, Path: "$.fake", Char: '.'},
		{Type: EventNumber, Path: "$.fake", Char: '1'},
		{Type: EventObjectEnd, Path: "$", Char: '}'},
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
	parser := NewParser()
	// Test various escaped whitespace characters
	json := `{"tab":"te\t","newline":"li\nne","return":"car\r","backspace":"bs\b","formfeed":"ff\f"}`

	events := []Event{}
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
	parser := NewParser()
	// Test complex string with multiple escaped whitespace characters
	json := `{"mixed":"\t\n\r\b\f","nested":{"inner":"text\twith\nnewlines"}}`

	events := []Event{}
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
	parser := NewParser()
	// Test string with Unicode escape sequences (common in JSON)
	json := `{"unicode":"\u0041\u0042\u0043","with_spaces" : "text\t\nmore"}`

	events := []Event{}
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
	parser := NewParser()
	// Test string with Unicode escape sequences (common in JSON)
	json := `{"key_with_e\n\"":"value_with_escape\u0041"}`

	events := []Event{}
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
