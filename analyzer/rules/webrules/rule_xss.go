package rules

import (
	"github.com/yourusername/yourproject/parser"
	"github.com/yourusername/yourproject/reporter"
	"regexp"
	"strings"
)

type RuleXSS struct {
	reporter *reporter.Reporter
}

// NewRuleXSS 创建并返回一个新的RuleXSS实例
func NewRuleXSS(reporter *reporter.Reporter) *RuleXSS {
	return &RuleXSS{
		reporter: reporter,
	}
}

// Apply 应用规则XSS，并将结果添加到报告中
func (r *RuleXSS) Apply(ast *parser.ASTNode) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
}

// CheckConditionA 检查敏感信息泄露并报告
func (r *RuleXSS) CheckConditionA(ast *parser.ASTNode) {
	// 查找HTML标签和属性值，检查是否存在危险的HTML属性（如onclick、onmouseover等）
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
	re := regexp.MustCompile(`(?i)<[^>]*?\s(on\w+)=["']?([^"'>]+)["']?`)
	nodes := ast.FindAll(parser.NodeTypeElement)
	for _, node := range nodes {
		if !strings.HasPrefix(node.TokenLiteral(), "<script") && !strings.HasPrefix(node.TokenLiteral(), "<style") {
			for _, attr := range node.Attributes {
				if re.MatchString(attr.Value) {
					r.reporter.AddIssue(ast.FilePath, node.Line, "Possible XSS vulnerability in attribute '"+attr.Key+"'")
				}
			}
		}
	}
}

// CheckConditionB 检查敏感信息泄露并报告
func (r *RuleXSS) CheckConditionB(ast *parser.ASTNode) {
	// 检查JavaScript中是否存在危险的函数（如eval、setInterval等）
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
	re := regexp.MustCompile(`(?i)\b(eval|setInterval|setTimeout|document\.write|document\.writeln|document\.innerhtml|window\.location)\b`)
	nodes := ast.FindAll(parser.NodeTypeScript)
	for _, node := range nodes {
		if re.MatchString(node.TokenLiteral()) {
			r.reporter.AddIssue(ast.FilePath, node.Line, "Possible XSS vulnerability in script: "+node.TokenLiteral())
		}
	}
}
