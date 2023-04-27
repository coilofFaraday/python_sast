package parser

import "go/token"

// Expression is an interface for all expression nodes
type Expression interface {
	expressionNode()
}

type Identifier struct {
	Token token.Token // the token.Token.IDENT token.Token
	Value string
}

type IntegerLiteral struct {
	Token token.Token // the token.Token.INT token.Token
	Value int64
}

type FloatLiteral struct {
	Token token.Token // the token.Token.FLOAT token.Token
	Value float64
}

type StringLiteral struct {
	Token token.Token // the token.Token.STRING token.Token
	Value string
}

type BooleanLiteral struct {
	Token token.Token // the token.Token.TRUE or token.Token.FALSE token.Token
	Value bool
}

type NoneLiteral struct {
	Token token.Token // The 'None' token.Token
}

type UnaryExpression struct {
	Token    token.Token // the prefix token.Token; e.g. '!', '-'
	Operator string
	Right    *Node
}

type BinaryExpression struct {
	Token    token.Token // the operator token.Token; e.g. '+'
	Operator string
	Left     *Node
	Right    *Node
}

type CallExpression struct {
	Token     token.Token // the '(' token.Token
	Function  *Node
	Arguments []*Node
}

// Attribute represents an attribute access expression in Python, such as 'a.b'.
type Attribute struct {
	Token token.Token // The token.Token.DOT token.Token
	Value *Node       // The expression that represents the object being accessed
	Name  *Identifier // The attribute being accessed
}

// Subscript represents a subscript expression in Python, such as 'a[0]'.
type Subscript struct {
	Token token.Token // The token.Token.LBRACKET token.Token
	Value *Node       // The expression that represents the object being accessed
	Index *Node       // The expression that represents the index value
}

// List represents a list literal in Python, such as '[1, 2, 3]'.
type List struct {
	Token    token.Token // The token.Token.LBRACKET token.Token
	Elements []*Node     // The expressions that represent the elements of the list
}

// Tuple represents a tuple literal in Python, such as '(1, 2, 3)'.
type Tuple struct {
	Token    token.Token // The token.Token.LPAREN token.Token
	Elements []*Node     // The expressions that represent the elements of the tuple
}

type Dictionary struct {
	Token token.Token // The token.Token.LBRACE token.Token
	Pairs [][2]*Node  // The key-value pairs of the dictionary
}

type Lambda struct {
	Token  token.Token // The token.Token.LAMBDA token.Token
	Params []*Node     // The lambda function parameters
	Body   *Node       // The lambda function body (an expression)
}

type ListComp struct {
	Token token.Token // The token.Token.LBRACKET token.Token
	Expr  *Node       // The main expression of the list comprehension
	Comps []*Node     // The comprehension clauses (for, if)
}

type DictComp struct {
	Token token.Token // The token.Token.LBRACE token.Token
	Key   *Node       // The key expression of the dict comprehension
	Value *Node       // The value expression of the dict comprehension
	Comps []*Node     // The comprehension clauses (for, if)
}

type SetComp struct {
	Token token.Token // The token.Token.LBRACE token.Token
	Expr  *Node       // The main expression of the set comprehension
	Comps []*Node     // The comprehension clauses (for, if)
}

type Generator struct {
	Token token.Token // The token.Token.LPAREN token.Token
	Expr  *Node       // The main expression of the generator
	Comps []*Node     // The comprehension clauses (for, if)
}

type Await struct {
	Token token.Token // The token.Token.AWAIT token.Token
	Value *Node       // The expression representing the value to await
}

type FormattedStr struct {
	Token  token.Token // The token.Token.F_STRING token.Token
	Values []*Node     // The expressions representing the values to format
}

func (i *Identifier) expressionNode()        {}
func (il *IntegerLiteral) expressionNode()   {}
func (fl *FloatLiteral) expressionNode()     {}
func (sl *StringLiteral) expressionNode()    {}
func (bl *BooleanLiteral) expressionNode()   {}
func (nl *NoneLiteral) expressionNode()      {}
func (ue *UnaryExpression) expressionNode()  {}
func (be *BinaryExpression) expressionNode() {}
func (a *Attribute) expressionNode()         {}
func (s *Subscript) expressionNode()         {}
func (l *List) expressionNode()              {}
func (t *Tuple) expressionNode()             {}
func (d *Dictionary) expressionNode()        {}
func (g *Generator) expressionNode()         {}
func (a *Await) expressionNode()             {}
func (f *FormattedStr) expressionNode()      {}
func (l *ListComp) expressionNode()          {}
func (d *DictComp) expressionNode()          {}
func (s *SetComp) expressionNode()           {}
func (l *Lambda) expressionNode()            {}
