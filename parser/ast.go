package parser

import (
	"strings"

	sasttoken "github.com/coiloffaraday/python_sast/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token sasttoken.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      sasttoken.Token // The token.FUNCTION token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out strings.Builder

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(" ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     sasttoken.Token // The '(' token
	Function  Expression      // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out strings.Builder

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       sasttoken.Token // The token.IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out strings.Builder

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      sasttoken.Token // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ListLiteral struct {
	Token    sasttoken.Token // The '[' token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	var out strings.Builder

	elements := []string{}
	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type ForStatement struct {
	Token    sasttoken.Token // The token.FOR token
	Iterator *Identifier
	Iterable Expression
	Body     *BlockStatement
	ElseBody *BlockStatement // Optional else block
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out strings.Builder

	out.WriteString("for ")
	out.WriteString(fs.Iterator.String())
	out.WriteString(" in ")
	out.WriteString(fs.Iterable.String())
	out.WriteString(": ")
	out.WriteString(fs.Body.String())

	if fs.ElseBody != nil {
		out.WriteString(" else ")
		out.WriteString(fs.ElseBody.String())
	}

	return out.String()
}

type WhileStatement struct {
	Token     sasttoken.Token // The token.WHILE token
	Condition Expression
	Body      *BlockStatement
	ElseBody  *BlockStatement // Optional else block
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out strings.Builder

	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(": ")
	out.WriteString(ws.Body.String())

	if ws.ElseBody != nil {
		out.WriteString(" else ")
		out.WriteString(ws.ElseBody.String())
	}

	return out.String()
}

// ... existing code ...

type TupleLiteral struct {
	Token    sasttoken.Token // The '(' token
	Elements []Expression
}

func (tl *TupleLiteral) expressionNode()      {}
func (tl *TupleLiteral) TokenLiteral() string { return tl.Token.Literal }
func (tl *TupleLiteral) String() string {
	var out strings.Builder

	elements := []string{}
	for _, el := range tl.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(")")

	return out.String()
}

type DictLiteral struct {
	Token sasttoken.Token // The '{' token
	Pairs map[Expression]Expression
}

func (dl *DictLiteral) expressionNode()      {}
func (dl *DictLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl *DictLiteral) String() string {
	var out strings.Builder

	pairs := []string{}
	for key, value := range dl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type ImportStatement struct {
	Token  sasttoken.Token // The token.IMPORT token
	Module *Identifier
	Alias  *Identifier // Optional alias
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out strings.Builder

	out.WriteString("import ")
	out.WriteString(is.Module.String())

	if is.Alias != nil {
		out.WriteString(" as ")
		out.WriteString(is.Alias.String())
	}

	return out.String()
}

type FromImportStatement struct {
	Token      sasttoken.Token // The token.FROM token
	Module     *Identifier
	Level      int // Relative import level, 0 for absolute imports
	ImportList []*ImportSpec
}

type ImportSpec struct {
	Name  *Identifier
	Alias *Identifier // Optional alias
}

func (fis *FromImportStatement) statementNode()       {}
func (fis *FromImportStatement) TokenLiteral() string { return fis.Token.Literal }
func (fis *FromImportStatement) String() string {
	var out strings.Builder

	out.WriteString("from ")
	if fis.Level > 0 {
		out.WriteString(strings.Repeat(".", fis.Level))
	}
	out.WriteString(fis.Module.String())
	out.WriteString(" import ")

	importList := []string{}
	for _, spec := range fis.ImportList {
		if spec.Alias != nil {
			importList = append(importList, spec.Name.String()+" as "+spec.Alias.String())
		} else {
			importList = append(importList, spec.Name.String())
		}
	}

	out.WriteString(strings.Join(importList, ", "))

	return out.String()
}

// ... existing code ...

type SetLiteral struct {
	Token    sasttoken.Token // The '{' token
	Elements []Expression
}

func (sl *SetLiteral) expressionNode()      {}
func (sl *SetLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *SetLiteral) String() string {
	var out strings.Builder

	elements := []string{}
	for _, el := range sl.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

type ListComprehension struct {
	Token      sasttoken.Token // The '[' token
	Expression Expression
	ForClauses []*ForClause
	IfClauses  []*IfClause
}

type ForClause struct {
	Target Expression
	Iter   Expression
}

type IfClause struct {
	Condition Expression
}

func (lc *ListComprehension) expressionNode()      {}
func (lc *ListComprehension) TokenLiteral() string { return lc.Token.Literal }
func (lc *ListComprehension) String() string {
	var out strings.Builder

	out.WriteString("[")
	out.WriteString(lc.Expression.String())

	for _, forClause := range lc.ForClauses {
		out.WriteString(" for ")
		out.WriteString(forClause.Target.String())
		out.WriteString(" in ")
		out.WriteString(forClause.Iter.String())
	}

	for _, ifClause := range lc.IfClauses {
		out.WriteString(" if ")
		out.WriteString(ifClause.Condition.String())
	}

	out.WriteString("]")

	return out.String()
}

// ... Add additional AST nodes here as needed.
type IntegerLiteral struct {
	Token sasttoken.Token
	Value int64
}

type FloatLiteral struct {
	Token sasttoken.Token
	Value float64
}

type StringLiteral struct {
	Token sasttoken.Token
	Value string
}

type BooleanLiteral struct {
	Token sasttoken.Token
	Value bool
}

type NoneLiteral struct {
	Token sasttoken.Token
}

type PrefixExpression struct {
	Token    sasttoken.Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    sasttoken.Token
	Left     Expression
	Operator string
	Right    Expression
}
