package parser

import (
	"fmt"
	"strconv"

	sasttoken "github.com/coiloffaraday/python_sast/token"
)

LOWEST := 0
PREFIX := 1

func (p *Parser) parseExpression() (Expression, error) {
	return p.parsePrefix()
}

func (p *Parser) parsePrefix() (Expression, error) {
	switch p.curToken.Type {
	case sasttoken.INT:
		return p.parseIntegerLiteral()
	case sasttoken.FLOAT:
		return p.parseFloatLiteral()
	case sasttoken.IDENT:
		return p.parseIdentifier()
	case sasttoken.STRING:
		return p.parseStringLiteral()
	case sasttoken.TRUE, sasttoken.FALSE:
		return p.parseBooleanLiteral()
	case sasttoken.NONE:
		return p.parseNoneLiteral()
	case sasttoken.LAMBDA:
		return p.parseLambdaExpression()
	case sasttoken.LPAREN:
		return p.parseGroupedExpression()
	case sasttoken.LBRACKET:
		return p.parseListLiteral()
	case sasttoken.PLUS, sasttoken.MINUS, sasttoken.BANG:
		return p.parsePrefixExpression()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curToken)
	}
}

// ... other functions such as parseIntegerLiteral, parseFloatLiteral, etc.

func (p *Parser) parseInfixExpression(left Expression) (Expression, error) {
	infixExp := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	infixExp.Right = right

	return infixExp, nil
}

func (p *Parser) parseGroupedExpression() (Expression, error) {
	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if !p.expectPeek(sasttoken.RPAREN) {
		return nil, fmt.Errorf("expected right parenthesis, got %s", p.peekToken.Type)
	}

	return exp, nil
}

func (p *Parser) parseListLiteral() (Expression, error) {
	// Consume the opening bracket.
	p.nextToken()

	// Parse the list elements.
	elements, err := p.parseListElements()
	if err != nil {
		return nil, err
	}

	// Consume the closing bracket.
	if !p.expectPeek(sasttoken.RBRACKET) {
		return nil, fmt.Errorf("missing closing bracket")
	}

	return &ListLiteral{Elements: elements}, nil
}

func (p *Parser) parseIntegerLiteral() (Expression, error) {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	p.nextToken()

	return &IntegerLiteral{Value: value}, nil
}

func (p *Parser) parseFloatLiteral() (Expression, error) {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, err
	}

	p.nextToken()

	return &FloatLiteral{Value: value}, nil
}

func (p *Parser) parseIdentifier() (Expression, error) {
	ident := &Identifier{Value: p.curToken.Literal}
	p.nextToken()
	return ident, nil
}

func (p *Parser) parseStringLiteral() (Expression, error) {
	str := &StringLiteral{Value: p.curToken.Literal}
	p.nextToken()
	return str, nil
}

func (p *Parser) parseBooleanLiteral() (Expression, error) {
	value := p.curToken.Type == sasttoken.TRUE
	p.nextToken()
	return &BooleanLiteral{Value: value}, nil
}

func (p *Parser) parseNoneLiteral() (Expression, error) {
	p.nextToken()
	return &NoneLiteral{}, nil
}

func (p *Parser) parseLambdaExpression() (Expression, error) {
	// TODO: Implement lambda expression parsing
	return nil, fmt.Errorf("lambda expression parsing not implemented")
}

func (p *Parser) parseListElements() ([]Expression, error) {
	elements := []Expression{}

	// If the next token is a closing bracket, the list is empty.
	if p.peekToken.Type == sasttoken.RBRACKET {
		return elements, nil
	}

	// Parse the first element.
	element, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	elements = append(elements, element)

	// Continue parsing the rest of the elements.
	for p.peekToken.Type == sasttoken.COMMA {
		p.nextToken() // Consume the comma.
		p.nextToken() // Move to the next element.

		element, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}

	return elements, nil
}

func (p *Parser) expectPeek(t sasttoken.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) parsePrefixExpression() (*Expression, error) {
	operator := p.curToken.Literal
	p.nextToken()
	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	prefix := &PrefixExpression{Operator: operator, Right: right}
	node := Expression(prefix)
	return &node, nil
}
