package rules

import (
	"strings"

	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleIO struct {
	reporter *reporter.Reporter
}

// NewRuleIO 创建并返回一个新的RuleIO实例
func NewRuleIO(reporter *reporter.Reporter) *RuleIO {
	return &RuleIO{
		reporter: reporter,
	}
}

// Apply 应用规则IO，并将结果添加到报告中
func (r *RuleIO) Apply(ast *parser.ast.Node) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
	r.CheckConditionC(ast)
}

// CheckConditionA 检查所有使用os、subprocess、multiprocessing模块的函数是否传入了命令参数，并报告
func (r *RuleIO) CheckConditionA(ast *parser.ast.Node) {
	// 搜索所有使用os、subprocess、multiprocessing模块的函数调用表达式
	expressions := ast.SearchExpressions(func(node *parser.ast.Node) bool {
		if node.Type == parser.ExpressionTypeCall {
			function := node.Children[0]
			if function.token.TokenLiteral() == "os.system" || function.token.TokenLiteral() == "subprocess.call" || function.token.TokenLiteral() == "multiprocessing.Process" {
				return true
			}
		}
		return false
	})

	// 检查是否传入了命令参数
	for _, expression := range expressions {
		for _, argument := range expression.Children[1:] {
			if argument.Type == parser.ExpressionTypeString {
				if strings.HasPrefix(strings.TrimSpace(argument.token.TokenLiteral()), "-") {
					r.reporter.AddIssue(expression.token.TokenLiteral(), expression.token.Token.Line, "调用命令时未指定命令参数")
				}
			} else if argument.Type == parser.ExpressionTypeIdentifier {
				if strings.HasPrefix(strings.TrimSpace(argument.token.TokenLiteral()), "-") {
					r.reporter.AddIssue(expression.token.TokenLiteral(), expression.token.Token.Line, "调用命令时未指定命令参数")
				}
			}
		}
	}
}

// CheckConditionB 检查所有使用shutil、os模块的函数是否调用了rm、rmdir、remove、unlink等删除文件/目录的函数，并报告
func (r *RuleIO) CheckConditionB(ast *parser.ast.Node) {
	// 搜索所有使用shutil、os模块的函数调用表达式
	expressions := ast.SearchExpressions(func(node *parser.ast.Node) bool {
		if node.Type == parser.ExpressionTypeCall {
			function := node.Children[0]
			if function.token.TokenLiteral() == "shutil.rmtree" || function.token.TokenLiteral() == "os.remove" || function.token.TokenLiteral() == "os.rmdir" || function.token.TokenLiteral() == "os.unlink" {
				return true
			}
		}
		return false
	})

	// 报告所有删除文件/目录的调用
	for _, expression := range expressions {
		r.reporter.AddIssue(expression.token.TokenLiteral(), expression.token.Token.Line, "调用了删除文件/目录的函数："+expression.Children[0].token.TokenLiteral())
	}
}

// CheckConditionC 检查所有使用open函数打开文件的函数是否使用了完整的文件路径，防止访问意外文件，并报告
func (r *RuleIO) CheckConditionC(ast *parser.ast.Node) {
	for _, node := range ast.Children {
		switch n := node.(type) {
		case *parser.FunctionDef:
			// 检查函数中是否包含文件写入操作
			if r.hasFileWriteOperation(n) {
				r.reporter.AddIssue(
					ast.Filename,
					n.LineNumber,
					"Potential security issue: writing to file in function "+n.Name.Value,
				)
			}
		case *parser.ClassDef:
			// 检查类中是否包含文件写入操作
			if r.hasFileWriteOperation(n) {
				r.reporter.AddIssue(
					ast.Filename,
					n.LineNumber,
					"Potential security issue: writing to file in class "+n.Name.Value,
				)
			}
		}
	}
}

// hasFileWriteOperation 检查函数或类中是否包含文件写入操作
func (r *RuleIO) hasFileWriteOperation(node parser.Node) bool {
	for _, child := range node.Children {
		if op, ok := child.(*parser.Call); ok {
			// 检查函数调用中是否包含文件写入操作
			if r.isFileWriteOperation(op) {
				return true
			}
		}
		if r.hasFileWriteOperation(child) {
			return true
		}
	}
	return false
}

// isFileWriteOperation 检查函数调用是否是文件写入操作
func (r *RuleIO) isFileWriteOperation(call *parser.Call) bool {
	if len(call.Args) < 2 {
		return false
	}
	// 检查函数名是否为 open
	if funcName, ok := call.Func.(*parser.Name); ok && funcName.Id == "open" {
		// 检查文件操作模式是否为写入
		if mode, ok := call.Args[1].(*parser.Str); ok {
			if strings.Contains(mode.S, "w") || strings.Contains(mode.S, "a") {
				return true
			}
		}
	}
	return false
}
