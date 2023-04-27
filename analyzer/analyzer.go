package analyzer

import (
	"fmt"

	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
	"github.com/coiloffaraday/python_sast/rules"
	"github.com/coiloffaraday/python_sast/utils"
)

type Config struct {
	Rules []string `json:"rules"`
}

type Analyzer struct {
	reporter *reporter.Reporter
	config   *Config
}

// NewAnalyzer 创建并返回一个新的Analyzer实例
func NewAnalyzer(config *Config) *Analyzer {
	return &Analyzer{
		reporter: reporter.NewReporter(),
		config:   config,
	}
}

// AnalyzeFile 分析指定文件并将结果添加到报告中
func (a *Analyzer) AnalyzeFile(file string) error {
	source, err := utils.ReadFile(file)
	if err != nil {
		return err
	}

	lexer := parser.NewLexer(source)
	parser := parser.NewParser(lexer)
	ast, err := parser.Parse()
	if err != nil {
		return err
	}

	for _, ruleName := range a.config.Rules {
		switch ruleName {
		case "RuleSQL":
			ruleSQL := rules.NewRuleSQL(a.reporter)
			ruleSQL.Apply(ast)
		case "RuleIO":
			ruleIO := rules.NewRuleIO(a.reporter)
			ruleIO.Apply(ast)
		case "RuleCSRF":
			ruleCSRF := rules.NewRuleCSRF(a.reporter)
			ruleCSRF.Apply(ast)
		default:
			fmt.Printf("Warning: Unknown rule '%s' specified in config.\n", ruleName)
		}
	}
	return nil
}

// GetReport 返回分析器的报告
func (a *Analyzer) GetReport() *reporter.Reporter {
	return a.reporter
}
