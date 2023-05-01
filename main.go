package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coiloffaraday/python_sast/analyzer"
	"github.com/coiloffaraday/python_sast/lexer"
	"github.com/coiloffaraday/python_sast/parser"
	"github.com/coiloffaraday/python_sast/reporter"
)

func main() {
	helpFlag := flag.Bool("h", false, "Display help message")
	helpFlagLong := flag.Bool("help", false, "Display help message")
	dirFlag := flag.String("d", "", "Directory path to analyze")
	dirFlagLong := flag.String("dir", "", "Directory path to analyze")
	fileFlag := flag.String("f", "", "File path to analyze")
	fileFlagLong := flag.String("file", "", "File path to analyze")

	flag.Parse()

	if *helpFlag || *helpFlagLong {
		displayHelp()
		return
	}

	dir := *dirFlag
	if *dirFlagLong != "" {
		dir = *dirFlagLong
	}

	file := *fileFlag
	if *fileFlagLong != "" {
		file = *fileFlagLong
	}

	if dir == "" && file == "" {
		fmt.Println("Error: You must provide either a directory or a file to analyze.")
		displayHelp()
		return
	}

	rep := reporter.NewReporter()

	if dir != "" {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error while accessing path %q: %v\n", path, err)
				return err
			}

			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".py") {
				analyzeFile(path, rep)
			}

			return nil
		})
	} else if file != "" {
		analyzeFile(file, rep)
	}

	rep.PrintReport()
}

func displayHelp() {
	fmt.Println("Usage:")
	fmt.Println("  -h, --help       Display help message")
	fmt.Println("  -d, --dir DIR    Analyze all Python files in the specified directory")
	fmt.Println("  -f, --file FILE  Analyze the specified Python file")
}

func analyzeFile(file string, rep *reporter.Reporter) {
	fmt.Printf("Analyzing file: %s\n", file)

	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error while reading file %q: %v\n", file, err)
		return
	}

	l := lexer.NewLexer(string(content))
	p := parser.New(l)
	program, err := p.ParseProgram()

	if err != nil {
		fmt.Printf("Error while parsing file %q: %v\n", file, err)
		return
	}

	analyzer := analyzer.NewAnalyzer(rep)
	analyzer.Analyze(program)
}
