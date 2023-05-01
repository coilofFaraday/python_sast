package parser

import (
	"fmt"

	sasttoken "github.com/coiloffaraday/python_sast/token"
)

func (p *Parser) parseStatement() (Node, error) {
	switch p.curToken.Type {
	case sasttoken.DEF:
		return p.parseFunction()
	case sasttoken.RETURN:
		return p.parseReturnStatement()
	case sasttoken.IMPORT:
		return p.parseImportStatement()
	case sasttoken.FROM:
		return p.parseImportFromStatement()
	case sasttoken.IDENT:
		return p.parseIdentifierStatement()
	case sasttoken.IF:
		return p.parseIfStatement()
	case sasttoken.ELIF:
		return p.parseElifStatement()
	case sasttoken.ELSE:
		return p.parseElseStatement()
	case sasttoken.FOR:
		return p.parseForStatement()
	case sasttoken.WHILE:
		return p.parseWhileStatement()
	case sasttoken.PASS:
		return p.parsePassStatement()
	case sasttoken.BREAK:
		return p.parseBreakStatement()
	case sasttoken.CONTINUE:
		return p.parseContinueStatement()
	default:
		expression, err := p.parseExpressionStatement()
		if err != nil {
			return nil, err
		}
		return expression, nil
	}
}

func (p *Parser) parseFunction() (Node, error) {
	p.nextToken()
	functionName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.LPAREN)
	params, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}
	p.expectPeek(sasttoken.RPAREN)
	p.expectPeek(sasttoken.COLON)

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &Function{
		Name:       functionName.(*Identifier),
		Parameters: params,
		Body:       body,
	}, nil
}

func (p *Parser) parseReturnStatement() (Node, error) {
	p.nextToken()

	retValue, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.SEMICOLON)

	return &ReturnStatement{ReturnValue: retValue}, nil
}

func (p *Parser) parseImportStatement() (Node, error) {
	p.nextToken()

	module, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ImportStatement{Module: module.(*Identifier)}, nil
}

func (p *Parser) parseImportFromStatement() (Node, error) {
	p.nextToken()

	module, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.IMPORT)
	imports, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}

	return &ImportFromStatement{Module: module.(*Identifier), Imports: imports}, nil
}

func (p *Parser) parseIdentifierStatement() (Node, error) {
	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if p.peekTokenIs(sasttoken.ASSIGN) {
		p.nextToken()
		p.nextToken()
		value, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		return &AssignmentStatement{Left: ident, Value: value}, nil
	}

	return ident, nil
}

func (p *Parser) parseIfStatement() (Node, error) {
	p.nextToken()
	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.COLON)
	consequence, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	var alternatives []*ElifStatement
	var elseStmt *ElseStatement

	for p.peekTokenIs(sasttoken.ELIF) {
		p.nextToken()
		elifCondition, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		p.expectPeek(sasttoken.COLON)
		elifConsequence, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		alternatives = append(alternatives, &ElifStatement{Condition: elifCondition, Consequence: elifConsequence})
	}

	if p.peekTokenIs(sasttoken.ELSE) {
		p.nextToken()
		p.expectPeek(sasttoken.COLON)
		elseBody, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		elseStmt = &ElseStatement{Body: elseBody}
	}

	return &IfStatement{
		Condition:   condition,
		Consequence: consequence,
		ElifClauses: alternatives,
		ElseClause:  elseStmt,
	}, nil
}

func (p *Parser) parseForStatement() (Node, error) {
	p.nextToken()
	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.IN)
	iterable, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.COLON)
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ForStatement{Identifier: ident, Iterable: iterable, Body: body}, nil
}

func (p *Parser) parseWhileStatement() (Node, error) {
	p.nextToken()
	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.expectPeek(sasttoken.COLON)
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) parseTryStatement() (Node, error) {
	p.nextToken()
	p.expectPeek(sasttoken.COLON)
	tryBlock, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	var exceptClauses []*ExceptStatement
	var elseStmt *ElseStatement
	var finallyStmt *FinallyStatement

	for p.peekTokenIs(sasttoken.EXCEPT) {
		p.nextToken()
		var exceptType Node
		if !p.peekTokenIs(sasttoken.COLON) {
			exceptType, err = p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
		}
		p.expectPeek(sasttoken.COLON)
		exceptBlock, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		exceptClauses = append(exceptClauses, &ExceptStatement{ExceptionType: exceptType, Body: exceptBlock})
	}

	if p.peekTokenIs(sasttoken.ELSE) {
		p.nextToken()
		p.expectPeek(sasttoken.COLON)
		elseBody, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		elseStmt = &ElseStatement{Body: elseBody}
	}

	if p.peekTokenIs(sasttoken.FINALLY) {
		p.nextToken()
		p.expectPeek(sasttoken.COLON)
		finallyBody, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		finallyStmt = &FinallyStatement{Body: finallyBody}
	}

	return &TryStatement{
		TryBlock:      tryBlock,
		ExceptClauses: exceptClauses,
		ElseClause:    elseStmt,
		FinallyClause: finallyStmt,
	}, nil
}

func (p *Parser) parseExpressionStatement() (*ExpressionStatement, error) {
	stmt := &ExpressionStatement{Token: p.curToken}

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Expression = exp

	if p.peekTokenIs(sasttoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseElseStatement() (Node, error) {
	stmt := &ElseStatement{}

	if !p.expectPeek(sasttoken.COLON) {
		return nil, fmt.Errorf("expected ':' after 'else', got %s", p.peekToken)
	}

	p.nextToken()

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	stmt.Block = block

	return stmt, nil
}

func (p *Parser) parseElifStatement() (Node, error) {
	stmt := &ElifStatement{}

	p.nextToken()

	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Condition = condition

	if !p.expectPeek(sasttoken.COLON) {
		return nil, fmt.Errorf("expected ':' after 'elif' condition, got %s", p.peekToken)
	}

	p.nextToken()

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	stmt.Block = block

	return stmt, nil
}

func (p *Parser) parsePassStatement() (Node, error) {
	stmt := &PassStatement{}

	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseBreakStatement() (Node, error) {
	stmt := &BreakStatement{}

	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseContinueStatement() (Node, error) {
	stmt := &ContinueStatement{}

	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseFunctionParameters() ([]*Identifier, error) {
	var identifiers []*Identifier

	if p.peekTokenIs(sasttoken.RPAREN) {
		p.nextToken()
		return identifiers, nil
	}

	p.nextToken()
	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	identifiers = append(identifiers, ident.(*Identifier))

	for p.peekTokenIs(sasttoken.COMMA) {
		p.nextToken()
		p.nextToken()

		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		identifiers = append(identifiers, ident.(*Identifier))
	}

	return identifiers, nil
}

func (p *Parser) parseIdentifierList() ([]*Identifier, error) {
	var identifiers []*Identifier

	p.nextToken()
	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	identifiers = append(identifiers, ident.(*Identifier))

	for p.peekTokenIs(sasttoken.COMMA) {
		p.nextToken()
		p.nextToken()

		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		identifiers = append(identifiers, ident.(*Identifier))
	}

	return identifiers, nil
}

func (p *Parser) parseBlockStatement() (*BlockStatement, error) {
	block := &BlockStatement{}
	block.Statements = []Node{}

	p.nextToken()

	for !p.curTokenIs(sasttoken.DEDENT) && !p.curTokenIs(sasttoken.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block, nil
}
