package lexer

import (
	sasttoken "github.com/coiloffaraday/python_sast/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	Line         int
	FilePath     string
}

func NewLexer(input string, filePath string) *Lexer {
	l := &Lexer{input: input, Line: 1, FilePath: filePath}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.Line++
	}
}

func (l *Lexer) NextToken() sasttoken.Token {
	var tok sasttoken.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(sasttoken.ASSIGN, l.ch, l.Line, l.FilePath)
	case ';':
		tok = newToken(sasttoken.SEMICOLON, l.ch, l.Line, l.FilePath)
	case '(':
		tok = newToken(sasttoken.LPAREN, l.ch, l.Line, l.FilePath)
	case ')':
		tok = newToken(sasttoken.RPAREN, l.ch, l.Line, l.FilePath)
	case ',':
		tok = newToken(sasttoken.COMMA, l.ch, l.Line, l.FilePath)
	case '+':
		tok = newToken(sasttoken.PLUS, l.ch, l.Line, l.FilePath)
	case '-':
		tok = newToken(sasttoken.MINUS, l.ch, l.Line, l.FilePath)
	case '!':
		tok = newToken(sasttoken.BANG, l.ch, l.Line, l.FilePath)
	case '/':
		tok = newToken(sasttoken.SLASH, l.ch, l.Line, l.FilePath)
	case '*':
		tok = newToken(sasttoken.ASTERISK, l.ch, l.Line, l.FilePath)
	case '<':
		tok = newToken(sasttoken.LT, l.ch, l.Line, l.FilePath)
	case '>':
		tok = newToken(sasttoken.GT, l.ch, l.Line, l.FilePath)
	case ':':
		tok = newToken(sasttoken.COLON, l.ch, l.Line, l.FilePath)
	case '|':
		tok = newToken(sasttoken.PIPE, l.ch, l.Line, l.FilePath)
	case '&':
		tok = newToken(sasttoken.AMPERSAND, l.ch, l.Line, l.FilePath)
	case '%':
		tok = newToken(sasttoken.MOD, l.ch, l.Line, l.FilePath)
	case '^':
		tok = newToken(sasttoken.CARET, l.ch, l.Line, l.FilePath)
	case '~':
		tok = newToken(sasttoken.TILDE, l.ch, l.Line, l.FilePath)
	case 0:
		tok.Literal = ""
		tok.Type = sasttoken.EOF
		tok.Line = l.Line
		tok.FilePath = l.FilePath
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = sasttoken.LookupIdent(tok.Literal)
			tok.Line = l.Line
			tok.FilePath = l.FilePath
			return tok
		} else if isDigit(l.ch) {
			tok.Type = sasttoken.INT
			tok.Literal = l.readNumber()
			tok.Line = l.Line
			tok.FilePath = l.FilePath
			return tok
		} else {
			tok = newToken(sasttoken.ILLEGAL, l.ch, l.Line, l.FilePath)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType sasttoken.TokenType, ch byte, line int, filePath string) sasttoken.Token {
	return sasttoken.Token{Type: tokenType, Literal: string(ch), Line: line, FilePath: filePath}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
