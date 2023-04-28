package analyzer

import (
	"fmt"
	"io/ioutil"

	"github.com/coiloffaraday/python_sast/lexer"
	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
	"github.com/coiloffaraday/python_sast/rules"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules []string `yaml:"rules"`
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

// LoadConfigFromFile 从YAML文件中加载配置
func LoadConfigFromFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// AnalyzeFile 分析指定文件并将结果添加到报告中
func (a *Analyzer) AnalyzeFile(file string) error {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	l := lexer.NewLexer(string(source))
	parser := parser.New(l)
	program := parser.ParseProgram()

	for _, ruleName := range a.config.Rules {
		switch ruleName {
		case "RuleSQL":
			ruleSQL := rules.NewRuleSQLInjection(a.reporter)
			ruleSQL.Apply(program)
		case "RuleIO":
			ruleIO := rules.NewRuleIO(a.reporter)
			ruleIO.Apply(program)
		case "RuleCSRF":
			ruleCSRF := rules.NewRuleCSRF(a.reporter)
			ruleCSRF.Apply(program)
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
