package rules

import (
	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleXSSSemantic struct {
	reporter *reporter.Reporter
}

// NewRuleXSSSemantic 创建并返回一个新的RuleXSSSemantic实例
func NewRuleXSSSemantic(reporter *reporter.Reporter) *RuleXSSSemantic {
	return &RuleXSSSemantic{
		reporter: reporter,
	}
}

// Apply 应用规则并将结果添加到报告中
func (r *RuleXSSSemantic) Apply(ast *parser.ASTNode) {
	r.checkForXSS(ast)
}

// checkForXSS 使用语义分析检查XSS漏洞并将结果添加到报告中
func (r *RuleXSS) checkForXSS(node *parser.ASTNode) {
	// 找到所有与输出相关的节点，例如函数调用和变量赋值语句
	outputNodes := node.FindAll(parser.NodeTypeFunctionCall, parser.NodeTypeAssignment)

	for _, outputNode := range outputNodes {
		// 获取输出值
		outputValue := outputNode.GetChildByIndex(len(outputNode.Children) - 1)

		// 判断输出值是否需要编码
		isEncoded := false
		switch outputNode.Type {
		case parser.NodeTypeFunctionCall:
			// 判断函数是否是 HTML 编码函数
			functionName := outputNode.GetChildByIndex(0).TokenLiteral()
			if functionName == "html" || functionName == "htmlspecialchars" {
				isEncoded = true
			}
		case parser.NodeTypeAssignment:
			// 判断赋值语句右侧的表达式是否是 HTML 编码函数的调用
			expr := outputNode.GetChildByIndex(1)
			if expr.Type == parser.NodeTypeFunctionCall {
				functionName := expr.GetChildByIndex(0).TokenLiteral()
				if functionName == "html" || functionName == "htmlspecialchars" {
					isEncoded = true
				}
			}
		}

		// 如果输出值未经过编码，则进行进一步检查
		if !isEncoded {
			// 检查输出值是否包含用户输入
			inputNodes := outputNode.FindAll(parser.NodeTypeIdentifier)
			for _, inputNode := range inputNodes {
				inputName := inputNode.TokenLiteral()

				// 根据变量名查找变量定义
				varDecl := node.FindVariableDeclaration(inputName)

				// 如果变量未被定义，或者未被赋值，则忽略
				if varDecl == nil || varDecl.GetChildByIndex(1) == nil {
					continue
				}

				// 获取变量赋值的值，并判断是否是用户输入
				varValue := varDecl.GetChildByIndex(1)
				if varValue.Type == parser.NodeTypeInput {
					// 发现了未编码的用户输入，将其添加到报告中
					r.reporter.AddIssue(node.Filename, outputNode.Line, "Unencoded user input found in output: "+outputValue.TokenLiteral())
				}
			}
		}
	}
}
