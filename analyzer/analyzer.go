package analyzer

import (
	"encoding/json"
	"io/ioutil"

	"github.com/yourusername/yourproject/parser"
	"github.com/yourusername/yourproject/reporter"
	"github.com/yourusername/yourproject/utils"
)

type Config struct {
	Rules       []string `json:"rules"`
	MaxFileSize int      `json:"max_file_size"`
}

type Analyzer struct {
	reporter *reporter.Reporter
	config   *Config
}

// Rule 接口表示一个静态分析规则
type Rule interface {
	Apply(file string, ast *parser.ASTNode, reporter *reporter.Reporter)
}

// NewAnalyzer 创建并返回一个新的Analyzer实例
func NewAnalyzer(configFile string) (*Analyzer, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	rules, err := loadRules(config.Rules)
	if err != nil {
		return nil, err
	}

	return &Analyzer{
		reporter: reporter.NewReporter(),
		rules:    rules,
	}, nil
}

// loadRules 加载并初始化指定的规则
func loadRules(ruleConfigs []RuleConfig) ([]Rule, error) {
	rules := []Rule{}

	for _, ruleConfig := range ruleConfigs {
		// TODO: 根据 ruleConfig.File 加载并初始化相应的规则实现
		// rule := ...
		// rules = append(rules, rule)
	}

	return rules, nil
}

// AnalyzeFiles 分析指定的文件列表并将结果添加到报告中
func (a *Analyzer) AnalyzeFiles(files []string) error {
	for _, file := range files {
		err := a.analyzeFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// analyzeFile 分析指定文件并将结果添加到报告中
func (a *Analyzer) analyzeFile(file string) error {
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

	for _, rule := range a.rules {
		rule.Apply(file, ast, a.reporter)
	}

	return nil
}

// GetReport 返回分析器的报告
func (a *Analyzer) GetReport() *reporter.Reporter {
	return a.reporter
}
