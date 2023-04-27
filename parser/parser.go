package parser

//修改：
// 统一处理错误
import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/coiloffaraday/python_sast/lexer"
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

func (p *Parser) Parse() (*ast.Node, error) {
	// 实现递归下降解析的入口点
	ast, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func (p *Parser) parseStatement() (*ast.Node, error) {
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
		case "True":
			return p.parseTrueStatement()
		case "False":
			return p.parseFalseStatement()
		case "None":
			return p.parseNoneStatement()
		case "and":
			return p.parseAndStatement()
		case "or":
			return p.parseOrStatement()
		case "not":
			return p.parseNotStatement()
		case "is":
			return p.parseIsStatement()
		default:
			return nil, fmt.Errorf("unexpected keyword: %s", token.Value)
		}

	case lexer.Identifier:
		return p.parseIdentifierStatement()
	case lexer.StringLiteral:
		return p.parseStringLiteralStatement()
	case lexer.NumberLiteral:
		return p.parseNumberLiteralStatement()
	case lexer.BooleanLiteral:
		return p.parseBooleanLiteralStatement()
	case lexer.NoneLiteral:
		return p.parseNoneLiteralStatement()
	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
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

func (p *Parser) parseBlockStatement() (*ast.Node, error) {
	// 解析代码块
	err := p.expectToken(LBRACE)
	if err != nil {
		return nil, err
	}

	node := &ast.Node{
		Type:     NodeStatementList,
		Token:    p.currentToken,
		Children: []*ast.Node{},
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
func (p *Parser) parseFunctionDefinition() (*ast.Node, error) {
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

	return &ast.Node{
		Type: FunctionDefinition,
		Value: &FunctionDefinitionNode{
			Name:       nameToken.Value,
			Parameters: params,
			Body:       body,
		},
	}, nil
}

func (p *Parser) parseClassDefinition() (*ast.Node, error) {
	// 解析类定义
	err := p.expectToken(CLASS)
	if err != nil {
		return nil, err
	}

	className := &ast.Node{
		Type:  NodeIdentifier,
		Token: p.currentToken,
	}
	p.expectToken(IDENT)

	err = p.expectToken(LBRACE)
	if err != nil {
		return nil, err
	}

	methods := []*ast.Node{}
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

	return &ast.Node{
		Type:     NodeClassDefinition,
		Token:    className.Token,
		Children: methods,
	}, nil
}

func (p *Parser) parseIfStatement() (*ast.Node, error) {
	// 解析 if 语句
	node := &ast.Node{
		Type:     NodeCallExpression,
		Token:    p.currentToken,
		Children: []*ast.Node{},
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

func (p *Parser) parseElseStatement() (*ast.Node, error) {
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

func (p *Parser) parseForStatement() (*ast.Node, error) {
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
	node := &ast.Node{
		Type: NodeStatementList,
		Children: []*ast.Node{
			initStmt,
			condition,
			updateStmt,
			body,
		},
	}

	return node, nil
}

func (p *Parser) parseWhileStatement() (*ast.Node, error) {
	// 解析while循环语句
	node := &ast.Node{
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

func (p *Parser) parseReturnStatement() (*ast.Node, error) {
	// 跳过 `return` 关键字
	p.current++

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.Node{
		Type:  ReturnStatement,
		Value: "return",
		Children: []*ast.Node{
			expr,
		},
	}, nil
}

func (p *Parser) parseInStatement(left *ast.Node) (*ast.Node, error) {
	p.consumeToken(lexer.Operator, "in")

	right, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.Node{
		Type:     ASTIn,
		Value:    "in",
		Children: []*ast.Node{left, right},
	}, nil
}

func (p *Parser) parseTrueStatement() (*ast.Node, error) {
	if p.tokens[p.current].Value != "True" {
		return nil, errors.New("unexpected token")
	}

	p.current++
	return &ast.Node{
		Type:  True,
		Value: "True",
	}, nil
}

func (p *Parser) parseFalseStatement() (*ast.Node, error) {
	// 解析 `False`
	token := p.tokens[p.current]
	if token.Type != lexer.Keyword || token.Value != "False" {
		return nil, fmt.Errorf("expected False, but got %s", token.Value)
	}
	p.current++

	return &ast.Node{
		Type:  BooleanLiteral,
		Value: "False",
	}, nil
}

func (p *Parser) parseNoneStatement() (*ast.Node, error) {
	// 检查当前Token是否是NoneLiteral
	token := p.tokens[p.current]
	if token.Type != lexer.NoneLiteral {
		return nil, fmt.Errorf("unexpected token: %v", token)
	}

	p.current++ // 跳过NoneLiteral

	// 创建ast.Node
	return Newast.Node(NoneLiteral, "None"), nil
}

func (p *Parser) parseAndStatement() (*ast.Node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	if p.peek().Type != Operator || p.peek().Value != "and" {
		return left, nil
	}

	// 读取and
	p.advance()

	// 解析右侧表达式
	right, err := p.parseAndStatement()
	if err != nil {
		return nil, err
	}

	// 构建AST节点
	return &ast.Node{
		Type:  AndStatement,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseOrStatement() (*ast.Node, error) {
	leftNode, err := p.parseAndStatement()
	if err != nil {
		return nil, err
	}

	for p.matchToken(TokenType{Type: Operator, Value: "or"}) {
		operator := p.lastMatchedToken()

		rightNode, err := p.parseAndStatement()
		if err != nil {
			return nil, err
		}

		leftNode = Newast.Node(OrExpression, operator, leftNode, rightNode)
	}

	return leftNode, nil
}

func (p *Parser) parseNotStatement() (*ast.Node, error) {
	token := p.tokens[p.current]
	if token.Type != lexer.Operator || token.Value != "not" {
		return nil, fmt.Errorf("unexpected token: %s", token.Value)
	}
	p.current++

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.Node{
		Type:     ASTNotNode,
		Children: []*ast.Node{expr},
	}, nil
}

func (p *Parser) parseIsStatement() (*ast.Node, error) {
	// 读取 is 左边的表达式
	expr1, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// 读取 is
	token := p.peek()
	if token.Type != Operator || token.Value != "is" {
		return nil, errors.New("unexpected token")
	}
	p.consume()

	// 读取 is 右边的表达式
	expr2, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// 生成 AST 节点
	return &ast.Node{
		Type:     IsStatement,
		Operator: "is",
		Children: []*ast.Node{expr1, expr2},
	}, nil
}

func (p *Parser) parseIdentifierStatement() (*ast.Node, error) {
	token := p.tokens[p.current]
	p.current++

	if p.current >= len(p.tokens) {
		// 如果是最后一个Token，直接返回该Token对应的ast.Node
		return Newast.Node(Identifier, token.Value, token.Pos), nil
	}

	switch p.tokens[p.current].Type {
	case lexer.Operator:
		// 如果下一个Token是运算符，表示这是一个赋值语句
		p.current++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return Newast.Node(AssignmentStatement, token.Value, token.Pos, expr), nil
	case lexer.Delimiter:
		switch p.tokens[p.current].Value {
		case "(":
			// 如果下一个Token是左括号，表示这是一个函数调用
			p.current++
			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}
			return Newast.Node(FunctionCall, token.Value, token.Pos, args), nil
		default:
			// 其他情况暂不处理
			break
		}
	default:
		// 其他情况暂不处理
		break
	}

	// 如果没有下一个Token，或者下一个Token不是运算符或左括号，那么表示这是一个变量引用
	return Newast.Node(Identifier, token.Value, token.Pos), nil
}

func (p *Parser) parseStringLiteralStatement() (*ast.Node, error) {
	token := p.tokens[p.current]
	p.current++

	return Newast.Node(StringLiteral, token.Value, token.Pos), nil
}

func (p *Parser) parseNumberLiteralStatement() (*ast.Node, error) {
	token := p.tokens[p.current]
	p.current++

	return Newast.Node(NumberLiteral, token.Value, token.Pos), nil
}

func (p *Parser) parseBooleanLiteralStatement() (*ast.Node, error) {
	token := p.tokens[p.current]
	p.current++

	return Newast.Node(BooleanLiteral, token.Value, token.Pos), nil
}

func (p *Parser) parseNoneLiteralStatement() (*ast.Node, error) {
	token := p.tokens[p.current]
	p.current++

	return Newast.Node(NoneLiteral, token.Value, token.Pos), nil
}
