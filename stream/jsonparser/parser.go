// Package jsonparser provides a low-memory, character-by-character JSON parser
// that generates detailed parsing events with JSON path tracking.
package jsonparser

import (
	"fmt"
	"strconv"
	"strings"
)

// State 表示解析器的当前状态
type State int

// 定义解析器的各种状态
const (
	stateIdle    State = iota // 空闲状态，等待新的值
	stateString               // 处理字符串值
	stateNumber               // 处理数字值
	stateBoolean              // 处理布尔值
	stateNull                 // 处理null值
	stateKey                  // 处理对象键名
)

// EventType 表示解析过程中发生的事件类型
type EventType int

// 定义各种事件类型
const (
	EventUnknown      EventType = iota // 未知事件类型
	EventString                        // 字符串内容字符
	EventStringEscape                  // 字符串中的转义字符
	EventNumber                        // 数字字符
	EventBoolean                       // 布尔值字符
	EventNull                          // null值字符
	EventObjectStart                   // 对象开始 '{'
	EventObjectEnd                     // 对象结束 '}'
	EventArrayStart                    // 数组开始 '['
	EventArrayEnd                      // 数组结束 ']'
	EventKey                           // 对象键名字符
	EventKeyEscape                     // 键名中的转义字符
	EventComma                         // 逗号分隔符 ','
	EventColon                         // 冒号分隔符 ':'
	EventQuote                         // 引号 '"'
	EventWhitespace                    // 空白字符
)

// container 表示JSON中的容器结构（对象或数组）
type container struct {
	Type       containerType // 容器类型（对象或数组）
	ArrayIndex int           // 仅用于数组，表示当前索引
	Key        string        // 仅用于对象，表示当前键名
}

func (c *container) IsArray() bool {
	if c == nil {
		return false
	}
	return c.Type == containerTypeArray
}

func (c *container) IsObject() bool {
	if c == nil {
		return false
	}
	return c.Type == containerTypeObject
}

func (c *container) IsEmpty() bool {
	return c.Type == containerTypeObject && c.Key == "" || c.Type == containerTypeArray && c.ArrayIndex < 0
}

func (c *container) SetKey(s string) {
	if c != nil {
		c.Key = s
	}
}

func (c *container) SetArrayIndex(i int) {
	if c != nil {
		c.ArrayIndex = i
	}
}

type containerType int

const (
	containerTypeObject containerType = iota
	containerTypeArray
)

// event 表示解析过程中的一个事件
// 包含当前处理的字符、事件类型和JSON路径
type event struct {
	Char rune      `json:"char"` // 当前处理的字符
	Type EventType `json:"type"` // 事件类型
	Path string    `json:"path"` // JSON Pointer路径，例如：$.foo.bar, $[0].bar
}

// innerParser 是JSON流式解析器的主要结构
// 使用状态机模式逐个字符解析JSON
type innerParser struct {
	state          State       // 当前解析状态
	stack          []container // 容器栈，用于跟踪嵌套结构
	buffer         []rune      // 临时缓冲区，用于累积字符
	escapeNext     bool        // 标记下一个字符是否为转义字符
	pathCache      string      // 路径缓存，用于性能优化
	pathCacheDirty bool        // 标记路径缓存是否需要更新
}

// newInnerParser 创建一个新的JSON解析器实例
func newInnerParser() *innerParser {
	return &innerParser{
		state: stateIdle,
	}
}

// Push 将单个字符推送到解析器中
// 返回一个事件，如果当前字符不产生事件则返回nil
func (p *innerParser) Push(r rune) event {
	var event event

	// 根据当前状态处理字符
	switch p.state {
	case stateIdle:
		event = p.handleIdleState(r) // 处理空闲状态
	case stateKey:
		event = p.handleStrState(r, true) // 处理键名字符串
	case stateString:
		event = p.handleStrState(r, false) // 处理值字符串
	case stateNumber:
		event = p.handleNumberState(r) // 处理数字
	case stateBoolean, stateNull:
		event = p.handleKeywordState(r) // 处理关键字（true/false/null）
	}

	return event
}

func (p *innerParser) resetState() {
	p.state = stateIdle
}

func (p *innerParser) resetBuffer() {
	p.buffer = p.buffer[:0]
}

func (p *innerParser) popStack() {
	if len(p.stack) > 0 {
		p.stack = p.stack[:len(p.stack)-1]
		p.pathCacheDirty = true
	}
}

func (p *innerParser) peekStack() *container {
	if len(p.stack) == 0 {
		return nil
	}
	p.pathCacheDirty = true
	return &p.stack[len(p.stack)-1]
}

func (p *innerParser) pushStack(c container) {
	p.stack = append(p.stack, c)
	p.pathCacheDirty = true
}

func (p *innerParser) handleIdleState(r rune) event {
	switch r {
	case '{':
		p.pushStack(container{Type: containerTypeObject})
		return event{
			Char: r,
			Type: EventObjectStart,
			Path: p.buildPath(),
		}
	case '}':
		p.popStack()
		p.resetState()
		p.resetBuffer()
		return event{
			Char: r,
			Type: EventObjectEnd,
			Path: p.buildPath(),
		}
	case '[':
		path := p.buildPath()
		p.pushStack(container{Type: containerTypeArray})
		return event{
			Char: r,
			Type: EventArrayStart,
			Path: path,
		}
	case ']':
		p.resetState()
		p.resetBuffer()
		p.popStack()
		return event{
			Char: r,
			Type: EventArrayEnd,
			Path: p.buildPath(),
		}
	case '"':
		p.buffer = []rune{}
		if p.peekStack().IsObject() && p.peekStack().Key == "" {
			p.state = stateKey
		} else {
			p.state = stateString
		}
		return event{
			Char: r,
			Type: EventQuote,
			Path: p.buildPath(),
		}
	case ':':
		p.resetState()
		p.resetBuffer()
		return event{
			Char: r,
			Type: EventColon,
			Path: p.buildPath(),
		}
	case ',':
		p.resetState()
		p.resetBuffer()
		if p.peekStack().IsArray() {
			p.peekStack().ArrayIndex++
		} else if p.peekStack().IsObject() {
			p.peekStack().Key = ""
		}
		return event{
			Char: r,
			Type: EventComma,
			Path: p.buildPath(),
		}
	case ' ', '\t', '\n', '\r':
		return event{
			Char: r,
			Type: EventWhitespace,
			Path: p.buildPath(),
		}
	default:
		return p.handleValueStart(r)
	}
}

func (p *innerParser) handleStrState(r rune, isKey bool) event {
	if p.escapeNext {
		p.escapeNext = false
		p.buffer = append(p.buffer, r)
		var eventType EventType
		if isKey {
			eventType = EventKey
		} else {
			eventType = EventString
		}
		return event{
			Char: r,
			Type: eventType,
			Path: p.getPathCache(),
		}
	}

	switch r {
	case '"':
		path := p.getPathCache()
		if isKey {
			p.peekStack().SetKey(string(p.buffer))
		}
		p.resetState()
		return event{
			Char: r,
			Type: EventQuote,
			Path: path,
		}
	case '\\':
		p.buffer = append(p.buffer, r)
		if p.escapeNext {
			p.escapeNext = false
		} else {
			p.escapeNext = true
		}
		et := EventStringEscape
		if isKey {
			et = EventKeyEscape
		}
		return event{
			Char: r,
			Type: et,
			Path: p.getPathCache(),
		}
	default:
		p.buffer = append(p.buffer, r)
		var eventType EventType
		if isKey {
			eventType = EventKey
		} else {
			eventType = EventString
		}
		return event{
			Char: r,
			Type: eventType,
			Path: p.getPathCache(),
		}
	}
}

func (p *innerParser) handleNumberState(r rune) event {
	if isDigit(r) || r == '.' || r == 'e' || r == 'E' || r == '+' || r == '-' {
		p.buffer = append(p.buffer, r)
		return event{
			Char: r,
			Type: EventNumber,
			Path: p.getPathCache(),
		}
	}
	// Number ended
	p.resetState()
	// Reprocess this character in initial state
	return p.handleIdleState(r)
}

func (p *innerParser) handleKeywordState(r rune) event {
	if isKeywordChar(r) {
		p.buffer = append(p.buffer, r)
		var eventType EventType
		switch p.state {
		case stateBoolean:
			eventType = EventBoolean
		case stateNull:
			eventType = EventNull
		case stateNumber, stateString, stateKey, stateIdle:
			// These states should not occur in keyword state, but handle exhaustively
			eventType = EventUnknown
		}
		return event{
			Char: r,
			Type: eventType,
			Path: p.getPathCache(),
		}
	}
	// Keyword ended
	p.resetState()
	p.resetBuffer()
	// Reprocess this character in initial state
	return p.handleIdleState(r)
}

func (p *innerParser) handleValueStart(r rune) event {
	setBuffer := func(r rune) {
		p.resetBuffer()
		p.buffer = append(p.buffer, r)
	}
	switch {
	case isDigit(r) || r == '-':
		setBuffer(r)
		p.state = stateNumber
		return event{
			Char: r,
			Type: EventNumber,
			Path: p.getPathCache(),
		}
	case r == 't':
		setBuffer(r)
		p.state = stateBoolean
		return event{
			Char: r,
			Type: EventBoolean,
			Path: p.getPathCache(),
		}
	case r == 'f':
		setBuffer(r)
		p.state = stateBoolean
		return event{
			Char: r,
			Type: EventBoolean,
			Path: p.getPathCache(),
		}
	case r == 'n':
		setBuffer(r)
		p.state = stateNull
		return event{
			Char: r,
			Type: EventNull,
			Path: p.getPathCache(),
		}
	default:
		// should not happen, but handle gracefully
		return event{
			Char: r,
			Type: EventUnknown,
			Path: p.getPathCache(),
		}
	}
}

func (p *innerParser) getPathCache() string {
	if !p.pathCacheDirty {
		return p.pathCache
	}
	p.pathCacheDirty = false
	p.pathCache = p.buildPath()
	return p.pathCache
}

// buildPath 根据当前的容器栈构建JSON路径
// 例如：$.foo.bar[0].baz
func (p *innerParser) buildPath() string {
	if len(p.stack) == 0 {
		return "$"
	}
	path := strings.Builder{}
	path.WriteString("$")
	for _, c := range p.stack {
		if c.IsEmpty() {
			continue
		}
		if c.IsObject() {
			path.WriteRune('.')
			path.WriteString(c.Key)
		} else if c.IsArray() {
			path.WriteString(fmt.Sprintf("[%d]", c.ArrayIndex))
		}
	}
	return path.String()
}

// isDigit 检查字符是否为数字
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// isKeywordChar 检查字符是否为关键字字符（字母）
func isKeywordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// Parser is a parser for JSON streams.
type Parser struct {
	buf        []rune
	inner      *innerParser
	autoEscape bool
	escaping   bool // Whether to escape strings automatically
}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{
		buf:   make([]rune, 0, 8),
		inner: newInnerParser(),
	}
}

// AutoEscape enables automatic escaping of string values.
func (p *Parser) AutoEscape() {
	p.autoEscape = true
}

// Event represents a JSON event produced by the parser.
type Event struct {
	Val  string    // The string value of the event
	Type EventType // The type of the event
	Path string    // The JSON Pointer path of the event
}

func fromInnerEvent(e event) *Event {
	return &Event{
		Val:  string(e.Char),
		Type: e.Type,
		Path: e.Path,
	}
}

// Push adds a rune to the parser's buffer and processes it through the inner parser.
func (p *Parser) Push(r rune) *Event {
	e := p.inner.Push(r)
	if !p.autoEscape {
		return fromInnerEvent(e)
	}

	if e.Type == EventStringEscape {
		p.escaping = true
		p.buf = append(p.buf, r)
		return nil
	}

	if e.Type == EventString && p.escaping {
		p.buf = append(p.buf, r)
		unescaped, err := strconv.Unquote(`"` + string(p.buf) + `"`)
		if err != nil {
			return nil
		}
		p.escaping = false
		p.buf = p.buf[:0] // Clear the buffer after unescaping
		return &Event{
			Val:  unescaped,
			Type: e.Type,
			Path: e.Path,
		}
	}
	return fromInnerEvent(e)
}
