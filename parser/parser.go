package parser

import "github.com/yourusername/sast-python/lexer"

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token
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

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Node {
	program := &Node{
		Type:     NodeProgram,
		Children: []*Node{},
	}

	// TODO: Add logic to parse more complex programs.
	// This is just a starting point for you to build upon.
	for p.curToken.Type != lexer.EOF {
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
	case lexer.LET:
		return p.parseExpressionStatement()
	default:
		return nil
	}
}

func (p *Parser) parseExpressionStatement() *Node {
	stmt := &ExpressionStatement{Token: p.curToken}

	// TODO: Add logic to parse more complex expressions.
	// This is just a starting point for you to build upon.
	stmt.Expression = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return &Node{
		Type:     NodeStatementList,
		Token:    stmt.Token,
		Children: []*Node{stmt.Expression.(*Identifier)},
	}
}
