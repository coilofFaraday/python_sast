package rules

import (
	"fmt"
	"strings"

	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleFileInclude struct {
	reporter *reporter.Reporter
}

// NewRuleFileInclude 创建并返回一个新的RuleFileInclude实例
func NewRuleFileInclude(reporter *reporter.Reporter) *RuleFileInclude {
	return &RuleFileInclude{
		reporter: reporter,
	}
}

// Apply 应用规则并将结果添加到报告中
func (r *RuleFileInclude) Apply(ast *parser.ast.Node) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
}

// CheckConditionA 检查本地文件包含
func (r *RuleFileInclude) CheckConditionA(ast *parser.ast.Node) {
	findLocalFileInclude := func(node *parser.ast.Node) {
		if node.Type == parser.ImportStmt && node.token.Token.Literal == "include" {
			arg := node.Children[0].token.Token.Literal
			if strings.HasPrefix(arg, "./") || strings.HasPrefix(arg, "../") {
				r.reporter.AddIssue(node.token.Token.SourceFile, node.token.Token.Line, fmt.Sprintf("Local file inclusion detected: %s", arg))
			}
		}
	}
	ast.Traverse(findLocalFileInclude)
}

// CheckConditionB 检查远程文件包含
func (r *RuleFileInclude) CheckConditionB(ast *parser.ast.Node) {
	findRemoteFileInclude := func(node *parser.ast.Node) {
		if node.Type == parser.CallExpr && node.FunctionName() == "file_get_contents" {
			arg := node.Arguments[0].token.Token.Literal
			if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
				r.reporter.AddIssue(node.token.Token.SourceFile, node.token.Token.Line, fmt.Sprintf("Remote file inclusion detected: %s", arg))
			}
		}
	}
	ast.Traverse(findRemoteFileInclude)
}
