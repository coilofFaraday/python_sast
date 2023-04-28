package rules

import (
	"regexp"
	"strings"

	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
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
func (r *RuleXSS) Apply(ast *parser.Node) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
}

// CheckConditionA 检查敏感信息泄露并报告
func (r *RuleXSS) CheckConditionA(ast *parser.Node) {
	// 查找HTML标签和属性值，检查是否存在危险的HTML属性（如onclick、onmouseover等）
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
	re := regexp.MustCompile(`(?i)<[^>]*?\s(on\w+)=["']?([^"'>]+)["']?`)
	nodes := ast.FindAll(parser.NodeTypeElement)
	for _, node := range nodes {
		if !strings.HasPrefix(node.Properties["token"].(parser.Token).Literal(), "<script") && !strings.HasPrefix(node.Properties["token"].(parser.Token).Literal(), "<style") {
			for _, attr := range node.Properties["attributes"].([]parser.Attribute) {
				if re.MatchString(attr.Value) {
					r.reporter.AddIssue(ast.FilePath, node.Line, "Possible XSS vulnerability in attribute '"+attr.Key+"'")
				}
			}
		}
	}
}

// CheckConditionB 检查敏感信息泄露并报告
func (r *RuleXSS) CheckConditionB(ast *parser.Node) {
	// 检查JavaScript中是否存在危险的函数（如eval、setInterval等）
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
	re := regexp.MustCompile(`(?i)\b(eval|setInterval|setTimeout|document\.write|document\.writeln|document\.innerhtml|window\.location)\b`)
	nodes := ast.FindAll(parser.NodeTypeScript)
	for _, node := range nodes {
		if re.MatchString(node.Properties["token"].(parser.Token).Literal()) {
			r.reporter.AddIssue(ast.FilePath, node.Line, "Possible XSS vulnerability in script: "+node.Properties["token"].(parser.Token).Literal())
		}
	}
}
