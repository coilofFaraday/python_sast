package parser

import "go/token"

// NodeType is a type to represent different kinds of AST nodes
type NodeType string

const (
	NodeProgram NodeType = "Program"
	// Statement Nodes
	NodeAssignStmt NodeType = "AssignStmt"
	NodeExprStmt   NodeType = "ExprStmt"
	NodeIfStmt     NodeType = "IfStmt"
	NodeWhileStmt  NodeType = "WhileStmt"
	NodeForStmt    NodeType = "ForStmt"
	NodeFuncDef    NodeType = "FuncDef"
	NodeClassDef   NodeType = "ClassDef"
	NodeImport     NodeType = "Import"
	NodePass       NodeType = "Pass"
	NodeBlock      NodeType = "Block"
	// Expression Nodes
	NodeIdentifier NodeType = "Identifier"
	NodeInteger    NodeType = "Integer"
	NodeFloat      NodeType = "Float"
	NodeString     NodeType = "String"

	NodeUnaryOp   NodeType = "UnaryOp"
	NodeBinaryOp  NodeType = "BinaryOp"
	NodeBoolOp    NodeType = "BoolOp"
	NodeCompareOp NodeType = "CompareOp"

	NodeCall         NodeType = "Call"
	NodeAttribute    NodeType = "Attribute"
	NodeSubscript    NodeType = "Subscript"
	NodeList         NodeType = "List"
	NodeTuple        NodeType = "Tuple"
	NodeDictionary   NodeType = "Dictionary"
	NodeLambda       NodeType = "Lambda"
	NodeReturn       NodeType = "Return"
	NodeBreak        NodeType = "Break"
	NodeContinue     NodeType = "Continue"
	NodeYield        NodeType = "Yield"
	NodeListComp     NodeType = "ListComp"
	NodeDictComp     NodeType = "DictComp"
	NodeSetComp      NodeType = "SetComp"
	NodeGenerator    NodeType = "Generator"
	NodeAwait        NodeType = "Await"
	NodeFormattedStr NodeType = "FormattedStr"
)

type Node struct {
	Type       NodeType
	Token      token.Token
	Children   []*Node
	Properties map[string]interface{}
}
