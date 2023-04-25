package parser

type NodeType int

const (
	// 类型常量
	FunctionDefinition NodeType = iota
	ClassDefinition
	IfStatement
	// ...其他节点类型
)

type ASTNode struct {
	Type  NodeType
	Value interface{}
}

type FunctionDefinitionNode struct {
	Name       string
	Parameters []ParameterNode
	Body       *BlockNode
}

type ParameterNode struct {
	Name         string
	DefaultValue *ExpressionNode // 可以为nil，表示没有默认值
}

type ClassDefinitionNode struct {
	Name       string
	BaseClasses []*ExpressionNode // 可以为空，表示没有基类
	Body       *BlockNode
}

type IfStatementNode struct {
	Condition *ExpressionNode
	ThenBlock *BlockNode
	ElseBlock *BlockNode // 可以为nil，表示没有else部分
}

type BlockNode struct {
	Statements []ASTNode
}

type ExpressionNode struct {
	// ...表达式的结构
}

// ...其他节点类型定义
