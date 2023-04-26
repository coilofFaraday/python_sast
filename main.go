package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/coiloffaraday/python_sast/analyzer"
)

var inputFile string

func init() {
	flag.StringVar(&inputFile, "input", "", "Path to the Python source file to be analyzed.")
}

func main() {
	flag.Parse()

	if inputFile == "" {
		fmt.Println("Error: No input file provided.")
		flag.Usage()
		os.Exit(1)
	}

	configPath := "path/to/config.json" // 设置为你的配置文件路径
	config, err := parseConfigurationFile(configPath)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	// 实例化分析器并分析输入文件
	ana := analyzer.NewAnalyzer(config)
	err = ana.AnalyzeFile(inputFile)
	if err != nil {
		fmt.Printf("Error analyzing file: %v\n", err)
		os.Exit(2)
	}

	// 获取报告并打印
	report := ana.GetReport()
	report.Print()

	// ...其他扩展功能（例如保存报告到文件、转换报告格式等）
}

// parseConfigurationFile 解析配置文件并返回配置
func parseConfigurationFile(configPath string) (*analyzer.Config, error) {
	// 打开配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 解析JSON配置
	decoder := json.NewDecoder(file)
	config := &analyzer.Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// displayHelp 显示帮助信息
func displayHelp() {
	fmt.Println("Usage: python-sast -input <path/to/python/file>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

// ...更多扩展功能函数
