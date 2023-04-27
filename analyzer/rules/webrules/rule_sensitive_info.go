package rules

import (
	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleSensitiveInfo struct {
	reporter *reporter.Reporter
}

func NewRuleSensitiveInfo(reporter *reporter.Reporter) *RuleSensitiveInfo {
	return &RuleSensitiveInfo{
		reporter: reporter,
	}
}

func (r *RuleSensitiveInfo) Apply(ast *parser.ASTNode) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
	r.CheckConditionC(ast)
}

func (r *RuleSensitiveInfo) CheckConditionA(ast *parser.ASTNode) {
	// 检查敏感信息泄露
	// 根据你的语言和场景设计相关正则表达式和关键字
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
}

func (r *RuleSensitiveInfo) CheckConditionB(ast *parser.ASTNode) {
	// 检查敏感信息写入
	// 根据你的语言和场景设计相关正则表达式和关键字
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
}

func (r *RuleSensitiveInfo) CheckConditionC(ast *parser.ASTNode) {
	// 检查敏感信息加密
	// 根据你的语言和场景设计相关正则表达式和关键字
	// 使用ast.FindAll函数查找相关语法节点
	// 将发现的问题添加到报告中
}

// ...更多检查敏感信息的函数
