package analyzer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/coiloffaraday/python_sast/parser"

	"gopkg.in/yaml.v2"
)

type Rule interface {
	Check(*parser.Program) []ReportItem
}

type RuleConfig struct {
	Name      string `yaml:"name"`
	FilePath  string `yaml:"file_path"`
	ClassName string `yaml:"class_name"`
}

type Analyzer struct {
	rules      []Rule
	configFile string
}

func NewAnalyzer(configFile string) *Analyzer {
	return &Analyzer{
		configFile: configFile,
	}
}

func (a *Analyzer) LoadRules() error {
	ruleConfigs, err := a.loadRuleConfigs()
	if err != nil {
		return err
	}

	for _, ruleConfig := range ruleConfigs {
		rule, err := a.loadRule(ruleConfig)
		if err != nil {
			return err
		}
		a.rules = append(a.rules, rule)
	}

	return nil
}

func (a *Analyzer) Analyze(program *parser.Program) []ReportItem {
	reportItems := make([]ReportItem, 0)
	var wg sync.WaitGroup
	reportItemChan := make(chan []ReportItem, len(a.rules))

	for _, rule := range a.rules {
		wg.Add(1)
		go func(rule Rule) {
			defer wg.Done()
			reportItemChan <- rule.Check(program)
		}(rule)
	}

	wg.Wait()
	close(reportItemChan)

	for items := range reportItemChan {
		reportItems = append(reportItems, items...)
	}

	return reportItems
}

func (a *Analyzer) loadRuleConfigs() ([]RuleConfig, error) {
	configBytes, err := ioutil.ReadFile(a.configFile)
	if err != nil {
		return nil, err
	}

	var ruleConfigs []RuleConfig
	if err := yaml.Unmarshal(configBytes, &ruleConfigs); err != nil {
		return nil, err
	}

	return ruleConfigs, nil
}

func (a *Analyzer) loadRule(ruleConfig RuleConfig) (Rule, error) {
	pluginPath := filepath.Join(".", ruleConfig.FilePath)
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %s, error: %v", pluginPath, err)
	}

	symRule, err := plug.Lookup(ruleConfig.ClassName)
	if err != nil {
		return nil, fmt.Errorf("failed to find symbol %s, error: %v", ruleConfig.ClassName, err)
	}

	rule, ok := symRule.(Rule)
	if !ok {
		return nil, errors.New("loaded symbol is not of type Rule")
	}

	return rule, nil
}
