package rules

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

type RuleSQL struct {
	reporter *reporter.Reporter
}

// NewRuleSQL 创建并返回一个新的RuleSQL实例
func NewRuleSQL(reporter *reporter.Reporter) *RuleSQL {
	return &RuleSQL{
		reporter: reporter,
	}
}

// Apply 应用规则SQL，并将结果添加到报告中
func (r *RuleSQL) Apply(ast *parser.ast.Node) {
	r.CheckConditionA(ast)
	r.CheckConditionB(ast)
	r.CheckConditionC(ast)
	r.CheckConditionD(ast)
}

// CheckConditionA 检查是否存在 SQL 注入风险
func (r *RuleSQL) CheckConditionA(ast *parser.ast.Node) {
	for _, node := range ast.Nodes {
		switch n := node.(type) {
		case *parser.CallExpression:
			if strings.ToLower(n.Function.TokenLiteral()) == "execute" && len(n.Arguments) > 0 {
				arg := n.Arguments[0]
				if strLit, ok := arg.(*parser.StringLiteral); ok {
					if strings.Contains(strings.ToLower(strLit.Value), "select ") ||
						strings.Contains(strings.ToLower(strLit.Value), "drop ") ||
						strings.Contains(strings.ToLower(strLit.Value), "delete ") ||
						strings.Contains(strings.ToLower(strLit.Value), "insert ") ||
						strings.Contains(strings.ToLower(strLit.Value), "update ") {
						r.reporter.AddIssue(ast.TokenLiteral(), arg.Line(), "Potential SQL injection vulnerability")
					}
				}
			}
		}
	}
}

// CheckConditionC 检查弱口令漏洞
func (r *RuleSQL) CheckConditionB(ast *parser.ast.Node) {
	if ast.Type == parser.SQLNode {
		sqlStatement := ast.Value

		// 从 weak.txt 中读取弱口令
		file, err := os.Open("weak.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// 逐行读取并匹配密码
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			weakPassword := scanner.Text()
			if strings.Contains(sqlStatement, weakPassword) {
				r.reporter.AddIssue(ast.File, ast.Line, "使用弱口令 '"+weakPassword+"'")
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, child := range ast.Children {
		r.CheckConditionB(child)
	}
}

// CheckConditionC 检查常见的 SQL 注入漏洞
func (r *RuleSQL) CheckConditionC(ast *parser.ast.Node) {
	// 检查 Select 语句中的注入漏洞
	selectStmts := ast.FindAll(parser.NodeTypeSelectStmt)
	for _, stmt := range selectStmts {
		// 检查 where 子句中是否包含类似 'or'=' 的字符串
		whereClause := stmt.FindFirst(parser.NodeTypeWhereClause)
		if whereClause != nil && strings.Contains(whereClause.TokenLiteral(), "'or'='") {
			r.reporter.AddIssue(ast.TokenLiteral(), stmt.Token.Line, "SQL Injection Vulnerability: 'or'=' found in WHERE clause")
		}
		// 检查 select 子句中是否包含类似 'union' 的字符串
		selectClause := stmt.FindFirst(parser.NodeTypeSelectClause)
		if selectClause != nil && strings.Contains(selectClause.TokenLiteral(), "union") {
			r.reporter.AddIssue(ast.TokenLiteral(), stmt.Token.Line, "SQL Injection Vulnerability: 'union' found in SELECT clause")
		}
	}

	// 检查 Insert 语句中的注入漏洞
	insertStmts := ast.FindAll(parser.NodeTypeInsertStmt)
	for _, stmt := range insertStmts {
		// 检查 values 子句中是否包含未过滤的用户输入
		valuesClause := stmt.FindFirst(parser.NodeTypeValuesClause)
		if valuesClause != nil && strings.Contains(valuesClause.TokenLiteral(), "'") {
			r.reporter.AddIssue(ast.TokenLiteral(), stmt.Token.Line, "SQL Injection Vulnerability: Unsanitized user input in VALUES clause")
		}
	}

	// 检查 Update 语句中的注入漏洞
	updateStmts := ast.FindAll(parser.NodeTypeUpdateStmt)
	for _, stmt := range updateStmts {
		// 检查 set 子句中是否包含未过滤的用户输入
		setClause := stmt.FindFirst(parser.NodeTypeSetClause)
		if setClause != nil && strings.Contains(setClause.TokenLiteral(), "'") {
			r.reporter.AddIssue(ast.TokenLiteral(), stmt.Token.Line, "SQL Injection Vulnerability: Unsanitized user input in SET clause")
		}
	}
}

// CheckConditionD 检查 SQL 注入防御措施
func (r *RuleSQL) CheckConditionD(ast *parser.ast.Node) {
	// 检查是否使用了参数化查询，如使用了，即视为通过了 SQL 注入防御
	hasParameterizedQuery := false
	ast.Walk(func(n parser.Node, depth int) bool {
		if callExpr, ok := n.(*parser.CallExpression); ok {
			if funcName, ok := callExpr.Function.(*parser.Identifier); ok {
				if funcName.Value == "prepare" {
					hasParameterizedQuery = true
					return false
				}
			}
		}
		return true
	})

	if !hasParameterizedQuery {
		r.reporter.AddIssue(ast.TokenLiteral(), 0, "SQL Injection Vulnerability: Parameterized query not used")
	}
}
