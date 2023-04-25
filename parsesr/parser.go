package parser

import (
	"errors"
	"fmt"
	"sast_python/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (*ASTNode, error) {
	// 实现递归下降解析的入口点
	ast, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func (p *Parser) parseStatement() (*ASTNode, error) {
	token := p.tokens[p.current]

	switch token.Type {
	case lexer.Keyword:
		switch token.Value {
		case "def":
			return p.parseFunctionDefinition()
		case "class":
			return p.parseClassDefinition()
		case "if":
			return p.parseIfStatement()
		case "else":
			return p.parseElseStatement()
		case "for":
			return p.parseForStatement()
		case "while":
			return p.parseWhileStatement()
		case "return":
			return p.parseReturnStatement()
		case "in":
			return p.parseInStatement()
		// ...其他关键字处理
		default:
			return nil, fmt.Errorf("unexpected keyword: %s", token.Value)
		}
	// ...其他词素类型处理
	default:
		return nil, errors.New("unexpected token")
	}
}

// 工具函数
func (p *Parser) consume(tokenType lexer.TokenType, tokenValue string) (lexer.Token, error) {
	token := p.tokens[p.current]
	if token.Type == tokenType && (tokenValue == "" || token.Value == tokenValue) {
		p.current++
		return token, nil
	}
	return lexer.Token{}, fmt.Errorf("expected %s with value '%s', got %s with value '%s'", tokenType, tokenValue, token.Type, token.Value)
}

func (p *Parser) parseBlockStatement() (*ASTNode, error) {
	// 解析代码块
	err := p.expectToken(LBRACE)
	if err != nil {
		return nil, err
	}

	node := &ASTNode{
		Type:     NodeStatementList,
		Token:    p.currentToken,
		Children: []*ASTNode{},
	}

	for p.currentToken.Type != RBRACE {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, statement)
	}

	err = p.expectToken(RBRACE)
	if err != nil {
		return nil, err
	}

	return node, nil
}


// 关键字函数
func (p *Parser) parseFunctionDefinition() (*ASTNode, error) {
	// 消耗 "def" 关键字
	p.consume(lexer.Keyword, "def")

	// 解析函数名
	nameToken, err := p.consume(lexer.Identifier, "")
	if err != nil {
		return nil, fmt.Errorf("expected function name, got %v", err)
	}

	// 消耗 "("
	p.consume(lexer.Delimiter, "(")

	// 解析参数列表
	params, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}

	// 消耗 ")"
	p.consume(lexer.Delimiter, ")")

	// 消耗 ":"
	p.consume(lexer.Delimiter, ":")

	// 解析函数体
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ASTNode{
		Type: FunctionDefinition,
		Value: &FunctionDefinitionNode{
			Name:       nameToken.Value,
			Parameters: params,
			Body:       body,
		},
	}, nil
}

func (p *Parser) parseClassDefinition() (*ASTNode, error) {
	// 解析类定义
	err := p.expectToken(CLASS)
	if err != nil {
		return nil, err
	}

	className := &ASTNode{
		Type:  NodeIdentifier,
		Token: p.currentToken,
	}
	p.expectToken(IDENT)

	err = p.expectToken(LBRACE)
	if err != nil {
		return nil, err
	}

	methods := []*ASTNode{}
	for p.currentToken.Type != RBRACE {
		method, err := p.parseFunctionDefinition()
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	err = p.expectToken(RBRACE)
	if err != nil {
		return nil, err
	}

	return &ASTNode{
		Type:     NodeClassDefinition,
		Token:    className.Token,
		Children: methods,
	}, nil
}

func (p *Parser) parseIfStatement() (*ASTNode, error) {
	// 解析 if 语句
	node := &ASTNode{
		Type:     NodeCallExpression,
		Token:    p.currentToken,
		Children: []*ASTNode{},
	}

	// 解析 if 子句
	err := p.expectToken(IF)
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, condition)

	err = p.expectToken(THEN)
	if err != nil {
		return nil, err
	}

	consequence, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, consequence)

	// 解析 else 子句
	if p.currentToken.Type == ELSE {
		p.nextToken()

		alternative, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, alternative)
	}

	return node, nil
}

func (p *Parser) parseElseStatement() (*ASTNode, error) {
	// 解析else语句
	err := p.expectToken(ELSE)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (p *Parser) parseForStatement() (*ASTNode, error) {
	// 解析for循环语句
	err := p.expectToken(FOR)
	if err != nil {
		return nil, err
	}

	// 解析循环变量
	initStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// 解析循环条件
	err = p.expectToken(SEMICOLON)
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// 解析循环变量更新语句
	err = p.expectToken(SEMICOLON)
	if err != nil {
		return nil, err
	}
	updateStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// 解析循环体
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// 创建for循环语句节点
	node := &ASTNode{
		Type: NodeStatementList,
		Children: []*ASTNode{
			initStmt,
			condition,
			updateStmt,
			body,
		},
	}

	return node, nil
}

func (p *Parser) parseWhileStatement() (*ASTNode, error) {
	// 解析while循环语句
	node := &ASTNode{
		Type:  NodeWhileStatement,
		Token: p.currentToken,
	}

	// 解析条件表达式
	err := p.expectToken(LPAREN)
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, condition)

	err = p.expectToken(RPAREN)
	if err != nil {
		return nil, err
	}

	// 解析循环体
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, body)

	return node, nil
}

func (p *Parser) parseReturnStatement() (*ASTNode, error) {
	// 跳过 `return` 关键字
	p.current++

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ASTNode{
		Type:  ReturnStatement,
		Value: "return",
		Children: []*ASTNode{
			expr,
		},
	}, nil
}

func (p *Parser) parseInStatement(left *ASTNode) (*ASTNode, error) {
	p.consumeToken(lexer.Operator, "in")

	right, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ASTNode{
		Type:     ASTIn,
		Value:    "in",
		Children: []*ASTNode{left, right},
	}, nil
}

// ...其他解析函数
