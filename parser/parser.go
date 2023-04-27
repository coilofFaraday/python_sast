package parser

import (
	"fmt"

	"github.com/coiloffaraday/python_sast/lexer"
	"github.com/coiloffaraday/python_sast/token"
)

const LOWEST = -1

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be %s, got %s instead"
	p.errors = append(p.errors, fmt.Sprintf(msg, t, p.peekToken.Type))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Node {
	program := &Node{
		Type:     NodeProgram,
		Children: []*Node{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Children = append(program.Children, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() *Node {
	switch p.curToken.Type {
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.DEF:
		return p.parseFunctionDef()
	case token.CLASS:
		return p.parseClassDef()
	case token.ASSIGN:
		return p.parseAssignStatement()
	// TODO: Add more cases for different statement types as needed.
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIfStatement() *Node {
	// Consume 'if' token
	p.expectPeek(token.IF)

	// Parse condition
	condition := p.parseExpression(LOWEST)

	// Consume ':' and parse the consequent block
	p.expectPeek(token.COLON)
	consequent := p.parseBlockStatement()

	// Check for 'elif' or 'else' clauses
	var alternate *Node
	if p.peekTokenIs(token.ELIF) || p.peekTokenIs(token.ELSE) {
		p.nextToken()
		alternate = p.parseIfStatement()
	}

	return &Node{
		Type:       NodeIfStmt,
		Properties: map[string]*Node{"condition": condition, "consequent": consequent, "alternate": alternate},
	}
}

func (p *Parser) parseWhileStatement() *Node {
	// TODO: Implement the parsing logic for 'while' statements.
	return nil
}

func (p *Parser) parseForStatement() *Node {
	// TODO: Implement the parsing logic for 'for' statements.
	return nil
}

func (p *Parser) parseFunctionDef() *Node {
	// TODO: Implement the parsing logic for function definitions.
	return nil
}

func (p *Parser) parseClassDef() *Node {
	// TODO: Implement the parsing logic for class definitions.
	return nil
}

func (p *Parser) parseAssignStatement() *Node {
	// TODO: Implement the parsing logic for assignment statements.
	return nil
}

func (p *Parser) parseExpressionStatement() *Node {
	expressionStatement := &Node{
		Type: NodeExpressionStatement,
	}

	expressionStatement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expressionStatement
}

func (p *Parser) parseExpression(precedence int) *Node {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// 工具函数
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) parseBlockStatement() *Node {
	block := &Node{
		Type:     NodeBlock,
		Children: []*Node{},
	}

	p.nextToken()

	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.DEDENT) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Children = append(block.Children, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
