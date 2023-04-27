package parser

// Statement is an interface for all statement nodes
type Statement interface {
	statementNode()
}

type AssignStatement struct {
	token.Token    token.Token // the token.Token.ASSIGN token.Token
	Targets  []*Node
	Operator token.Token // the assignment operator, e.g., "=", "+=", "-=", etc.
	Values   []*Node
}

type ExprStatement struct {
	token.Token      token.Token // the first token.Token of the expression
	Expression *Node
}

type IfStatement struct {
	token.Token       token.Token // the token.Token.IF token.Token
	Condition   *Node
	Consequence *Node
	Alternative *Node
}

type WhileStatement struct {
	token.Token     token.Token // the token.Token.WHILE token.Token
	Condition *Node
	Body      *Node
}

type ForStatement struct {
	token.Token  token.Token // the token.Token.FOR token.Token
	Target *Node
	Iter   *Node
	Body   *Node
}

type FunctionDef struct {
	token.Token  token.Token // the token.Token.DEF token.Token
	Name   *Node
	Params []*Node
	Body   *Node
}

type ClassDef struct {
	token.Token       token.Token // the token.Token.CLASS token.Token
	Name        *Node
	BaseClasses []*Node
	Body        *Node
}

type ImportStatement struct {
	token.Token    token.Token // the token.Token.IMPORT token.Token
	Elements []*Node
}

type PassStatement struct {
	token.Token token.Token // the token.Token.PASS token.Token
}

type ReturnStatement struct {
	token.Token token.Token // The token.Token.RETURN token.Token
	Value *Node // The expression representing the value to return (optional)
}

type BreakStatement struct {
	token.Token token.Token // The token.Token.BREAK token.Token
}

type ContinueStatement struct {
	token.Token token.Token // The token.Token.CONTINUE token.Token
}

type YieldStatement struct {
	token.Token token.Token // The token.Token.YIELD token.Token
	Value *Node // The expression representing the value to yield (optional)
}

func (as *AssignStatement) statementNode()  {}
func (es *ExprStatement) statementNode()    {}
func (is *IfStatement) statementNode()      {}
func (ws *WhileStatement) statementNode()   {}
func (fs *ForStatement) statementNode()     {}
func (fd *FunctionDef) statementNode()      {}
func (cd *ClassDef) statementNode()         {}
func (imp *ImportStatement) statementNode() {}
func (ps *PassStatement) statementNode()    {}
func (rs *ReturnStatement) statementNode() {}
func (bs *BreakStatement) statementNode()  {}
func (cs *ContinueStatement) statementNode() {}
func (ys *YieldStatement) statementNode()  {}
