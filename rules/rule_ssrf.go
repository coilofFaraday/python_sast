package rules

import (
	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleSSRF struct {
	reporter *reporter.Reporter
}

func NewRuleSSRF(reporter *reporter.Reporter) *RuleSSRF {
	return &RuleSSRF{reporter: reporter}
}

func (r *RuleSSRF) CheckForSSRF(node parser.Node) {
	// 递归地检查所有子节点
	for _, child := range node.GetChildren() {
		r.CheckForSSRF(child)
	}

	// 如果当前节点是一个函数调用
	if callExpr, ok := node.(*parser.CallExpression); ok {
		// 检查函数名称是否与潜在的 SSRF 函数匹配
		if r.isSSRFFunction(callExpr.Function) {
			// 如果匹配，则报告 SSRF
			r.reporter.AddReportItem(reporter.ReportItem{
				Filename:    callExpr.Token.FilePath,
				Line:        callExpr.Token.Line,
				Column:      callExpr.Token.Column,
				RuleName:    "SSRF",
				Description: "Possible SSRF detected",
			})
		}
	}
}

// isSSRFFunction 检查给定的函数名是否与潜在的 SSRF 函数匹配
func (r *RuleSSRF) isSSRFFunction(functionName string) bool {
	// 在这里，我们添加一些可能导致 SSRF 的函数名
	ssrfFunctions := []string{
		"requests.get",
		"requests.post",
		"requests.options",
		// 添加更多潜在的 SSRF 函数
	}

	for _, ssrfFunc := range ssrfFunctions {
		if functionName == ssrfFunc {
			return true
		}
	}

	return false
}
