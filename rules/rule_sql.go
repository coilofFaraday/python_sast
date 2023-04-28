package rules

import (
	"regexp"

	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleSQLInjection struct {
	reporter *reporter.Reporter
}

// NewRuleSQLInjection 创建并返回一个新的RuleSQLInjection实例
func NewRuleSQLInjection(reporter *reporter.Reporter) *RuleSQLInjection {
	return &RuleSQLInjection{
		reporter: reporter,
	}
}

// Apply 应用规则SQL注入，并将结果添加到报告中
func (r *RuleSQLInjection) Apply(ast *parser.Node) {
	r.CheckSQLInjection(ast)
}

// CheckSQLInjection 检查SQL注入并报告
func (r *RuleSQLInjection) CheckSQLInjection(node *parser.Node) {
	// 检查字符串拼接、格式化字符串等可能导致SQL注入的用法
	// 遍历AST查找相关语法节点
	// 将发现的问题添加到报告中

	// 正则表达式检测SQL注入风险的模式
	sqlInjectionPattern := regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|CREATE|DROP|ALTER)\s`)

	if node.Type == parser.NodeString {
		if sqlInjectionPattern.MatchString(node.Properties["value"].(string)) {
			r.reporter.AddIssue(node.FilePath, node.Line, "Possible SQL injection vulnerability: "+node.Properties["value"].(string))
		}
	}

	for _, child := range node.Children {
		r.CheckSQLInjection(child)
	}
}
