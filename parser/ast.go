package ast

// Node 接口表示抽象语法树中的一个节点
type Node interface {
	// TokenLiteral 方法返回节点关联的Token的字面值
	TokenLiteral() string
}

// Expression 接口表示一个表达式节点
type Expression interface {
	Node
	expressionNode()
}

// Statement 接口表示一个语句节点
type Statement interface {
	Node
	statementNode()
}

// Program 结构体表示整个程序，由一系列语句组成
type Program struct {
	Statements []Statement
}

// TokenLiteral 方法返回程序的第一个语句的Token字面值
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Identifier 结构体表示一个标识符表达式节点
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// LetStatement 结构体表示一个赋值语句节点
type LetStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// ReturnStatement 结构体表示一个return语句节点
type ReturnStatement struct {
	Token       Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// ExpressionStatement 结构体表示一个表达式语句节点
type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// IntegerLiteral 结构体表示一个整数字面量表达式节点
type IntegerLiteral struct {
	Token Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// StringLiteral 结构体表示一个字符串字面量表达式节点
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// BooleanLiteral 结构体表示一个布尔字面量表达式节点
type BooleanLiteral struct {
	Token Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }

// NoneLiteral 结构体表示一个None字面量表达式节点
type NoneLiteral struct {
	Token Token
}

func (nl *NoneLiteral) expressionNode()      {}
func (nl *NoneLiteral) TokenLiteral() string { return nl.Token.Literal }

// PrefixExpression 结构体表示一个前缀表达式节点
type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// InfixExpression 结构体表示一个中缀表达式节点
type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// IfExpression 结构体表示一个if表达式节点
type IfExpression struct {
	Token       Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// BlockStatement 结构体表示一个代码块语句节点
type BlockStatement struct {
	Token      Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// FunctionLiteral 结构体表示一个函数字面量表达式节点
type FunctionLiteral struct {
	Token      Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// CallExpression 结构体表示一个函数调用表达式节点
type CallExpression struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// ClassDeclaration 结构体表示一个类声明语句节点
type ClassDeclaration struct {
	Token      Token
	Name       *Identifier
	BaseClass  *Identifier
	Statements []Statement
}

func (cd *ClassDeclaration) statementNode()       {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }

// These nodes should cover most of the common cases in a language.
// However, depending on the specific language and its features,
// additional nodes may be necessary.
