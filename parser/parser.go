package parser

import (
	"fmt"

	"github.com/coiloffaraday/python_sast/lexer"
	sasttoken "github.com/coiloffaraday/python_sast/token"
)

type Parser struct {
	lexer          *lexer.Lexer
	curToken       sasttoken.Token
	peekToken      sasttoken.Token
	errors         []string
	prefixParseFns map[sasttoken.TokenType]prefixParseFn
	infixParseFns  map[sasttoken.TokenType]infixParseFn
}

type (
	prefixParseFn func() (Node, error)
	infixParseFn  func(Node) (Node, error)
)

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.prefixParseFns = make(map[sasttoken.TokenType]prefixParseFn)
	p.registerPrefix(sasttoken.IDENT, p.parseIdentifier)
	p.registerPrefix(sasttoken.INT, p.parseIntegerLiteral)
	p.registerPrefix(sasttoken.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(sasttoken.STRING, p.parseStringLiteral)
	p.registerPrefix(sasttoken.LBRACKET, p.parseListLiteral)
	p.registerPrefix(sasttoken.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(sasttoken.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(sasttoken.NONE, p.parseNoneLiteral)
	p.registerPrefix(sasttoken.LAMBDA, p.parseLambdaExpression)
	p.registerPrefix(sasttoken.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[sasttoken.TokenType]infixParseFn)
	p.registerInfix(sasttoken.PLUS, p.parseInfixExpression)
	p.registerInfix(sasttoken.MINUS, p.parseInfixExpression)
	p.registerInfix(sasttoken.ASTERISK, p.parseInfixExpression)
	p.registerInfix(sasttoken.SLASH, p.parseInfixExpression)
	p.registerInfix(sasttoken.EQ, p.parseInfixExpression)

	// Read two tokens so that both curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() (*Node, error) {
	return p.parseStatement()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t sasttoken.TokenType) error {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
	return fmt.Errorf(msg)
}

func (p *Parser) peekTokenIs(t sasttoken.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t sasttoken.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) registerPrefix(tokenType sasttoken.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType sasttoken.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedenceTable[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedenceTable[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// All the statement and expression parsing functions should be included here as they are already defined in expression.go and statement.go files.
