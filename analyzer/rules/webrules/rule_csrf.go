package rules

import (
	"github.com/yourusername/yourproject/parser"
	"github.com/yourusername/yourproject/reporter"
)

type RuleCSRF struct {
	reporter *reporter.Reporter
}

// NewRuleCSRF 创建并返回一个新的RuleCSRF实例
func NewRuleCSRF(reporter *reporter.Reporter) *RuleCSRF {
	return &RuleCSRF{
		reporter: reporter,
	}
}

// Apply 应用规则并将结果添加到报告中
func (r *RuleCSRF) Apply(ast *parser.ASTNode) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
}

// CheckConditionA 检查是否存在未进行CSRF保护的HTTP请求，并报告
func (r *RuleCSRF) CheckConditionA(ast *parser.ASTNode) {
	// 检查所有HTTP请求
	httpRequests := ast.FindAll(parser.NodeTypeHTTPRequest)
	for _, httpRequest := range httpRequests {
		// 检查是否存在CSRF保护
		csrfProtection := httpRequest.FindFirst(parser.NodeTypeCSRFProtection)
		if csrfProtection == nil {
			// 发现未进行CSRF保护的HTTP请求
			r.reporter.AddIssue(httpRequest.File, httpRequest.Line, "HTTP request is not protected against CSRF attacks")
		}
	}
}

// CheckConditionB 检查是否存在未进行CSRF保护的HTTP请求，但同时具有危险操作（例如删除，修改等），并报告
func (r *RuleCSRF) CheckConditionB(ast *parser.ASTNode) {
	// 检查所有HTTP请求
	httpRequests := ast.FindAll(parser.NodeTypeHTTPRequest)
	for _, httpRequest := range httpRequests {
		// 检查是否存在危险操作
		if httpRequest.Contains(parser.NodeTypeDelete) || httpRequest.Contains(parser.NodeTypeUpdate) {
			// 检查是否存在CSRF保护
			csrfProtection := httpRequest.FindFirst(parser.NodeTypeCSRFProtection)
			if csrfProtection == nil {
				// 发现未进行CSRF保护的HTTP请求
				r.reporter.AddIssue(httpRequest.File, httpRequest.Line, "HTTP request with dangerous operation is not protected against CSRF attacks")
			}
		}
	}
}
