package lexer

import (
	"fmt"
	"strings"
)

// TokenType 表示Token的类型
type TokenType int

const (
	// Keyword 关键字
	Keyword TokenType = iota
	// Identifier 标识符
	Identifier
	// Operator 运算符
	Operator
	// Delimiter 分隔符
	Delimiter
	// StringLiteral 字符串字面量
	StringLiteral
	// NumberLiteral 数字字面量
	NumberLiteral
	// BooleanLiteral 布尔字面量
	BooleanLiteral
	// NoneLiteral None字面量
	NoneLiteral
	// EOF 文件结束标志
	EOF
)

// Token 代表代码中的单个词法单元
type Token struct {
	Type  TokenType // 类型
	Value string    // 值
	Pos   int       // 位置
}

// Lexer 代表词法分析器
type Lexer struct {
	input  string // 输入的代码
	tokens []Token
	pos    int // 当前分析的位置
}

// NewLexer 创建一个新的Lexer实例
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		tokens: []Token{},
		pos:    0,
	}
}

// NextToken 读取并返回下一个Token
func (l *Lexer) NextToken() Token {
	for l.pos < len(l.input) {
		c := l.input[l.pos]

		if isSpace(c) {
			l.skipSpace()
			continue
		}

		if isLetter(c) {
			return l.readIdentifierOrKeyword()
		}

		if isDigit(c) {
			return l.readNumberLiteral()
		}

		if c == '"' {
			return l.readStringLiteral()
		}

		if isOperator(c) {
			return l.readOperator()
		}

		if isDelimiter(c) {
			return l.readDelimiter()
		}

		panic(fmt.Sprintf("unexpected character: %c", c))
	}

	return Token{Type: EOF}
}

// 以下是辅助方法

func (l *Lexer) skipSpace() {
	for l.pos < len(l.input) && isSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *Lexer) readIdentifierOrKeyword() Token {
	start := l.pos
	for l.pos < len(l.input) && (isLetter(l.input[l.pos]) || isDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
		l.pos++
	}

	value := l.input[start:l.pos]
	tokenType := Identifier
	if isKeyword(value) {
		tokenType = Keyword
	}

	return Token{
		Type:  tokenType,
		Value: value,
		Pos:   start,
	}
}

func (l *Lexer) readNumberLiteral() Token {
	start := l.pos
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	value := l.input[start:l.pos]
	return Token{
		Type:  NumberLiteral,
		Value: value,
		Pos:   start,
	}
}

func (l *Lexer) readStringLiteral() Token {
	start := l.pos
	l.pos++ // 跳过第一个双引号

	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\\' {
			l.pos++
		}
		l.pos++
	}

	if l.pos >= len(l.input) {
		panic("unexpected end of input")
	}

	// 跳过最后一个双引号
	l.pos++

	value := l.input[start:l.pos]
	return Token{
		Type:  StringLiteral,
		Value: value,
		Pos:   start,
	}
}

func (l *Lexer) readOperator() Token {
	start := l.pos
	l.pos++

	// 连续读取运算符
	for l.pos < len(l.input) && isOperator(l.input[l.pos]) {
		l.pos++
	}

	value := l.input[start:l.pos]
	return Token{
		Type:  Operator,
		Value: value,
		Pos:   start,
	}
}

func (l *Lexer) readDelimiter() Token {
	start := l.pos
	l.pos++

	value := l.input[start:l.pos]
	return Token{
		Type:  Delimiter,
		Value: value,
		Pos:   start,
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isOperator(c byte) bool {
	return strings.ContainsAny(string(c), "+-*/%=&|<>!") || c == ':' || c == '.'
}

func isDelimiter(c byte) bool {
	return strings.ContainsAny(string(c), "()[]{},;")
}

func isKeyword(value string) bool {
	switch value {
	case "def", "if", "else", "for", "while", "in", "return", "True", "False", "None", "and", "or", "not", "is", "class":
		return true
	default:
		return false
	}
}
