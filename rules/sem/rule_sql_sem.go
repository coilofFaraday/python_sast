package rules

import (
	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleSQLInjection struct {
	reporter *reporter.Reporter
}

func NewRuleSQLInjection(reporter *reporter.Reporter) *RuleSQLInjection {
	return &RuleSQLInjection{
		reporter: reporter,
	}
}

func (r *RuleSQLInjection) Apply(ast *parser.ast.Node) {
	r.CheckSQLInjection(ast)
}

func (r *RuleSQLInjection) CheckSQLInjection(ast *parser.ast.Node) {
	for _, node := range ast.Children {
		if node.Type == parser.NodeTypeFunctionCall {
			functionNameNode := node.Children[0]
			if functionNameNode.Type == parser.NodeTypeIdentifier {
				functionName := functionNameNode.token.TokenLiteral()
				if isSQLInjectionFunction(functionName) {
					argumentNode := node.Children[1]
					if argumentNode.Type == parser.NodeTypeStringLiteral {
						if isUserInput(argumentNode.token.TokenLiteral()) {
							r.reporter.AddIssue(ast.FilePath, argumentNode.Line, "Possible SQL injection vulnerability")
						}
					}
				}
			}
		}
	}

	for _, child := range ast.Children {
		r.CheckSQLInjection(child)
	}
}

func isSQLInjectionFunction(functionName string) bool {
	sqlFunctions := []string{"mysql_query", "mysqli_query", "pg_query", "sqlite_query", "pg_send_query", "pg_query_params"}
	for _, f := range sqlFunctions {
		if f == functionName {
			return true
		}
	}
	return false
}

func isUserInput(node *parser.ast.Node) bool {
	// 检查节点类型是否为字符串字面量
	if node.Type == parser.StringLiteral {
		// 检查节点的父节点是否为函数调用
		parent := node.Parent
		if parent.Type == parser.CallExpression {
			// 检查函数名是否为预定义的输入函数，如get/post/cookie等
			function := parent.Children[0]
			if function.token.TokenLiteral() == "get" || function.token.TokenLiteral() == "post" || function.token.TokenLiteral() == "cookie" {
				return true
			}
		}
	}

	// 遍历子节点进行递归检查
	for _, child := range node.Children {
		if isUserInput(child) {
			return true
		}
	}

	return false
}
