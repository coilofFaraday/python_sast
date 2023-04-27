package parser

// Expression is an interface for all expression nodes
type Expression interface {
	expressionNode()
}

type Identifier struct {
	Token Token // the token.IDENT token
	Value string
}

type IntegerLiteral struct {
	Token Token // the token.INT token
	Value int64
}

type FloatLiteral struct {
	Token Token // the token.FLOAT token
	Value float64
}

type StringLiteral struct {
	Token Token // the token.STRING token
	Value string
}

type BooleanLiteral struct {
	Token Token // the token.TRUE or token.FALSE token
	Value bool
}

type NoneLiteral struct {
	Token Token // The 'None' token
}

type UnaryExpression struct {
	Token    Token // the prefix token; e.g. '!', '-'
	Operator string
	Right    *Node
}

type BinaryExpression struct {
	Token    Token // the operator token; e.g. '+'
	Operator string
	Left     *Node
	Right    *Node
}

type CallExpression struct {
	Token     Token // the '(' token
	Function  *Node
	Arguments []*Node
}

// Attribute represents an attribute access expression in Python, such as 'a.b'.
type Attribute struct {
	Token Token       // The token.DOT token
	Value *Node       // The expression that represents the object being accessed
	Name  *Identifier // The attribute being accessed
}

// Subscript represents a subscript expression in Python, such as 'a[0]'.
type Subscript struct {
	Token Token // The token.LBRACKET token
	Value *Node // The expression that represents the object being accessed
	Index *Node // The expression that represents the index value
}

// List represents a list literal in Python, such as '[1, 2, 3]'.
type List struct {
	Token    Token   // The token.LBRACKET token
	Elements []*Node // The expressions that represent the elements of the list
}

// Tuple represents a tuple literal in Python, such as '(1, 2, 3)'.
type Tuple struct {
	Token    Token   // The token.LPAREN token
	Elements []*Node // The expressions that represent the elements of the tuple
}

type Dictionary struct {
	Token Token      // The token.LBRACE token
	Pairs [][2]*Node // The key-value pairs of the dictionary
}

type Lambda struct {
	Token  Token   // The token.LAMBDA token
	Params []*Node // The lambda function parameters
	Body   *Node   // The lambda function body (an expression)
}

type ListComp struct {
	Token Token   // The token.LBRACKET token
	Expr  *Node   // The main expression of the list comprehension
	Comps []*Node // The comprehension clauses (for, if)
}

type DictComp struct {
	Token Token   // The token.LBRACE token
	Key   *Node   // The key expression of the dict comprehension
	Value *Node   // The value expression of the dict comprehension
	Comps []*Node // The comprehension clauses (for, if)
}

type SetComp struct {
	Token Token   // The token.LBRACE token
	Expr  *Node   // The main expression of the set comprehension
	Comps []*Node // The comprehension clauses (for, if)
}

type Generator struct {
	Token Token   // The token.LPAREN token
	Expr  *Node   // The main expression of the generator
	Comps []*Node // The comprehension clauses (for, if)
}

type Await struct {
	Token Token // The token.AWAIT token
	Value *Node // The expression representing the value to await
}

type FormattedStr struct {
	Token  Token   // The token.F_STRING token
	Values []*Node // The expressions representing the values to format
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
