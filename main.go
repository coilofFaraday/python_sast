package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/coiloffaraday/python_sast/analyzer"
)

func printHelp() {
	fmt.Println("Usage: python_sast [options] [files...]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help        Show this help message")
	fmt.Println("  -c, --config FILE Specify the path to the config.yaml file (default: 'config.yaml')")
	fmt.Println("  -d, --directory DIR Process all files in the specified directory")
	fmt.Println("  -f, --file FILE Process a single file")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  python_sast -c myconfig.yaml -d myproject/")
}

func processDirectory(analyzerInstance *analyzer.Analyzer, dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err := analyzerInstance.AnalyzeFile(path)
			if err != nil {
				log.Printf("Error analyzing file '%s': %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error processing directory '%s': %v\n", dir, err)
	}
}

func main() {
	help := flag.Bool("h", false, "Show help message")
	helpLong := flag.Bool("help", false, "Show help message")
	configPath := flag.String("c", "config.yaml", "Path to the config.yaml file")
	configPathLong := flag.String("config", "config.yaml", "Path to the config.yaml file")
	directory := flag.String("d", "", "Process all files in the specified directory")
	directoryLong := flag.String("directory", "", "Process all files in the specified directory")
	file := flag.String("f", "", "Process a single file")
	fileLong := flag.String("file", "", "Process a single file")

	flag.Parse()

	if *help || *helpLong {
		printHelp()
		return
	}

	config, err := analyzer.LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if *configPath != "config.yaml" && *configPathLong != "config.yaml" {
		log.Fatalf("Error: Both -c and --config options are provided. Please use only one of them.")
	}

	if *configPathLong != "config.yaml" {
		config, err = analyzer.LoadConfigFromFile(*configPathLong)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
	}

	analyzerInstance := analyzer.NewAnalyzer(config)

	if *directory != "" || *directoryLong != "" {
		dir := *directory
		if *directoryLong != "" {
			dir = *directoryLong
		}
		processDirectory(analyzerInstance, dir)
	} else if *file != "" || *fileLong != "" {
		filename := *file
		if *fileLong != "" {
			filename = *fileLong
		}
		err := analyzerInstance.AnalyzeFile(filename)
		if err != nil {
			log.Printf("Error analyzing file '%s': %v\n", filename, err)
		}
	} else {
		printHelp()
		return
	}

	report := analyzerInstance.GetReport()
	report.Print()
}
